// cmd/server/connections_handlers.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/connections"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/connectors"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
)

// connectionView is what GET /api/connections returns — never a password,
// only the env var name it's resolved from.
type connectionView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	TLS      string `json:"tls"`
	ReadOnly *bool  `json:"read_only"` // nil = not checked yet this run
	Status   string `json:"status"`    // "ok" | "unknown"
}

func (s *server) handleConnections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listConnections(w, r)
	case http.MethodPost:
		s.createConnection(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listConnections reads from the registry, not the SQLite table directly —
// the registry is the merged view (YAML-seeded connections + persisted ones,
// see main.go's startup wiring). Reading SQLite alone would hide any
// connection defined only in config/connections.yaml (e.g. the shipped
// dev_erp entry), even though it's fully usable by sql_query.
func (s *server) listConnections(w http.ResponseWriter, r *http.Request) {
	records := s.connReg.List()

	views := make([]connectionView, 0, len(records))
	for _, rec := range records {
		v := connectionView{
			ID: rec.ID, Name: rec.Name, Driver: rec.Driver, Host: rec.Host,
			Port: rec.Port, Database: rec.Database, Username: rec.Username, TLS: rec.TLS,
			Status: "unknown",
		}
		if ro, known := s.connReg.ReadOnly(rec.ID); known {
			roCopy := ro
			v.ReadOnly = &roCopy
			v.Status = "ok"
		}
		views = append(views, v)
	}
	writeJSON(w, map[string]interface{}{"connections": views, "total": len(views)})
}

// createConnection validates + registers the connection for immediate use
// (connReg.Add) and persists it so it survives a restart.
func (s *server) createConnection(w http.ResponseWriter, r *http.Request) {
	var cfg connections.ConnectionConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.connReg.Add(cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.kg.Store().UpsertConnection(cfg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[API] Saved connection %q (%s)", cfg.ID, cfg.Driver)
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]interface{}{"status": "saved", "id": cfg.ID})
}

// handleConnectionTest runs (or re-runs) the read-only health check.
func (s *server) handleConnectionTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	started := time.Now()
	readOnly, err := s.connReg.Test(id)
	latencyMs := time.Since(started).Milliseconds()

	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "latency_ms": latencyMs, "error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "latency_ms": latencyMs, "read_only": readOnly})
}

// handleConnectionPreview runs a query through the exact same sql_query
// guards (ensureSelectOnly / bindPositional / ensureLimit) by reusing
// SQLQueryHandler.Execute, capped at 5 rows for the Pipeline Studio preview.
func (s *server) handleConnectionPreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	var body struct {
		Query  string                 `json:"query"`
		Params map[string]interface{} `json:"params"`
		Limit  int                    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	limit := body.Limit
	if limit <= 0 || limit > 5 {
		limit = 5
	}

	handler := connectors.NewSQLQueryHandler(s.connReg)
	out, err := handler.Execute(r.Context(), map[string]interface{}{
		"connection_id": id,
		"query":         body.Query,
		"params":        body.Params,
		"limit":         float64(limit),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, out)
}

// handleConnectionDiscover browses the connection's schema (information_schema)
// and, as a side effect, runs the canonical-type heuristic and seeds any
// resulting SchemaMapping nodes into the KG — the IT-side analog of
// /api/opcua/discover triggering seedKG. See docs/analysis_log.md Entry 115
// for the full plan (Phase 1+2) this implements.
func (s *server) handleConnectionDiscover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	db, err := s.connReg.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tables, err := connections.DiscoverSchema(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	candidates := connections.SuggestMappings(tables)
	kgCandidates := make([]kg.SchemaMappingCandidate, 0, len(candidates))
	for _, c := range candidates {
		fieldMap := make(map[string]interface{}, len(c.FieldMap))
		for k, v := range c.FieldMap {
			fieldMap[k] = v
		}
		kgCandidates = append(kgCandidates, kg.SchemaMappingCandidate{
			Table: c.Table, CanonicalType: c.CanonicalType, Confidence: c.Confidence, FieldMap: fieldMap,
		})
	}

	seeded, err := s.kg.SeedSchemaMappings(id, kgCandidates)
	if err != nil {
		log.Printf("[API] Failed to seed schema mappings for %q: %v", id, err)
	}

	// Entity resolution (Entry 120): link any now-validated work_order
	// mapping to the real OT Equipment nodes it references. Runs on every
	// discover, same "browse triggers everything downstream" pattern as the
	// OT bootstrap — cheap and idempotent, so no separate trigger needed.
	resolved, err := s.ResolveWorkCenters(r.Context())
	if err != nil {
		log.Printf("[API] Entity resolution failed for %q: %v", id, err)
	}

	writeJSON(w, map[string]interface{}{
		"tables":            tables,
		"suggested":         candidates,
		"seeded":            seeded,
		"equipment_resolved": resolved,
	})
}

// handleConnectionDatabases browses every database and table visible to this
// connection's user in one call — "connect and see everything," separate
// from handleConnectionDiscover's canonical-mapping side effect, which stays
// scoped to the connection's own configured database. Purely read-only
// visibility; nothing here writes to the KG.
func (s *server) handleConnectionDatabases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	db, err := s.connReg.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	databases, err := connections.ListDatabasesAndTables(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"databases": databases, "total": len(databases)})
}

func (s *server) handleConnectionDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	s.connReg.Remove(id)
	if err := s.kg.Store().DeleteConnection(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"status": "deleted", "id": id})
}
