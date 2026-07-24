// cmd/server/entity_resolution.go
// OT<->IT entity resolution — the gap docs/analysis_log.md Entry 109
// explicitly flagged as missing entirely ("nothing computes the same_as
// OT<->IT match from Entry 102/103"). Links a validated work_order
// SchemaMapping to the real OT Equipment node(s) it references, instead of
// relying on both sides happening to use matching work_center strings by
// coincidence (confirmed live: the OT bootstrap seeded "Machine1"/"Machine3",
// the fake ERP uses "machine1"/"machine2"/"machine3" — different case, and
// no OT node for machine2 at all).
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
)

// normalizeWorkCenter makes an OT/IT work-center label comparable across
// naming conventions. Deliberately just lowercase + trim, not fuzzy/similarity
// matching — a match is always an exact, explainable identity claim, never a
// guess with a confidence score attached.
func normalizeWorkCenter(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// equipmentByWorkCenter indexes every OT Equipment node in a business graph
// by its normalized work_center property, for O(1) lookup during resolution
// and during live ActiveProduction queries.
func equipmentByWorkCenter(graph *kg.GraphJSON) map[string]kg.Node {
	idx := make(map[string]kg.Node)
	for _, n := range graph.Nodes {
		if n.Type != "Equipment" {
			continue
		}
		wc, _ := n.Properties["work_center"].(string)
		if wc == "" {
			continue
		}
		idx[normalizeWorkCenter(wc)] = n
	}
	return idx
}

// ResolveWorkCenters links IT-side work-center references to their OT
// Equipment node, where one exists. For every validated (pending:false)
// work_order SchemaMapping, queries the real distinct work_center values
// actually present in that ERP table (not assumed from the mapping alone),
// matches each against known OT Equipment nodes, and writes a persisted
// same_as edge where they resolve. Idempotent (AddEdgeCat is INSERT OR
// IGNORE) — safe to call after every /discover, same as the OT bootstrap.
func (s *server) ResolveWorkCenters(ctx context.Context) (int, error) {
	graph, err := s.kg.GetGraph("business")
	if err != nil {
		return 0, err
	}
	equipmentIdx := equipmentByWorkCenter(graph)

	resolved := 0
	for _, n := range graph.Nodes {
		if n.Type != "SchemaMapping" {
			continue
		}
		if pending, _ := n.Properties["pending"].(bool); pending {
			continue // not yet human-validated
		}
		canonicalType, _ := n.Properties["canonical_type"].(string)
		if canonicalType != "work_order" {
			continue
		}

		connectionID, _ := n.Properties["connection_id"].(string)
		table, _ := n.Properties["table"].(string)
		fieldMap, _ := n.Properties["field_map"].(map[string]interface{})
		wcCol, _ := fieldMap["work_center"].(string)
		if wcCol == "" || !validIdentifier.MatchString(table) || !validIdentifier.MatchString(wcCol) {
			continue
		}

		values, err := s.distinctWorkCenters(ctx, connectionID, table, wcCol)
		if err != nil {
			continue // one bad connection shouldn't fail resolution for the others
		}

		for _, wc := range values {
			eq, ok := equipmentIdx[normalizeWorkCenter(wc)]
			if !ok {
				continue // no matching OT Equipment — correctly left unresolved, not guessed
			}
			edgeID := fmt.Sprintf("edge_same_as_%s_%s", n.ID, eq.ID)
			if err := s.kg.AddEdgeCat(kg.CategoryBusiness, edgeID, eq.ID, n.ID, "same_as", 1.0); err != nil {
				continue
			}
			resolved++
		}
	}
	return resolved, nil
}

// distinctWorkCenters reads the real distinct work-center values present in
// an ERP table.
func (s *server) distinctWorkCenters(ctx context.Context, connectionID, table, wcCol string) ([]string, error) {
	db, err := s.connReg.Get(connectionID)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s IS NOT NULL", wcCol, table, wcCol)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
