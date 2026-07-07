// internal/kg/graph.go
// Unified Knowledge Graph (2026-07-02 merge — see docs/analysis_log.md Entry 50).
// One graph, one storage layer. Nodes/edges tagged with Category (business/platform).
// Platform sub-graph (Pipeline/Function/Topic/Connection/Dashboard) is rebuilt
// from the pipeline registry — this replaces the previous in-memory "Technical KG".
package kg

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
	"github.com/MindsetAdmin/mindset-data-edge/internal/storage"
)

// Node — unified KG node. Category = business | platform.
type Node struct {
	ID         string                 `json:"id"`
	Category   string                 `json:"category"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// Edge — unified KG relation. Category inherited from context (business events vs
// platform wiring). Cross-category edges (e.g., Dashboard→Event) are legal and
// take the more specific side's category — by convention, edges written during
// platform rebuild are category=platform even if they land on a business node.
type Edge struct {
	ID        string    `json:"id"`
	Category  string    `json:"category"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Relation  string    `json:"relation"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}

// GraphJSON — the shape returned to the frontend (Cytoscape-friendly).
type GraphJSON struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// KnowledgeGraph — the single source of truth for site + platform data.
type KnowledgeGraph struct {
	store        *storage.SQLiteStore
	platformMu   sync.Mutex // serialize platform rebuilds
	lastRebuilt  time.Time
	lastRegHash  string
}

// NewKnowledgeGraph opens (or creates) the unified KG.
func NewKnowledgeGraph(dbPath string) (*KnowledgeGraph, error) {
	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}
	return &KnowledgeGraph{store: store}, nil
}

// AddNode inserts or ignores. Uses category="business" as the safe default when
// legacy callers don't specify one; new code should use AddNodeCat.
func (kg *KnowledgeGraph) AddNode(id, nodeType, label string, props map[string]interface{}) error {
	return kg.AddNodeCat(CategoryBusiness, id, nodeType, label, props)
}

// AddNodeCat inserts or ignores a node under an explicit category.
func (kg *KnowledgeGraph) AddNodeCat(cat Category, id, nodeType, label string, props map[string]interface{}) error {
	var propsJSON string
	if props != nil {
		jsonBytes, err := json.Marshal(props)
		if err != nil {
			return err
		}
		propsJSON = string(jsonBytes)
	}
	query := `INSERT OR IGNORE INTO kg_nodes (id, category, type, label, properties) VALUES (?, ?, ?, ?, ?)`
	if _, err := kg.store.DB().Exec(query, id, string(cat), nodeType, label, propsJSON); err != nil {
		return err
	}
	return nil
}

// AddEdge — business default.
func (kg *KnowledgeGraph) AddEdge(id, fromID, toID, relation string, weight float64) error {
	return kg.AddEdgeCat(CategoryBusiness, id, fromID, toID, relation, weight)
}

// AddEdgeCat — explicit category.
func (kg *KnowledgeGraph) AddEdgeCat(cat Category, id, fromID, toID, relation string, weight float64) error {
	query := `INSERT OR IGNORE INTO kg_edges (id, category, from_id, to_id, relation, weight) VALUES (?, ?, ?, ?, ?, ?)`
	if _, err := kg.store.DB().Exec(query, id, string(cat), fromID, toID, relation, weight); err != nil {
		return err
	}
	return nil
}

// AddMicroStop enrichit le KG (business) avec un micro-stop
func (kg *KnowledgeGraph) AddMicroStop(eventID, workCenter string, duration float64) error {
	eventNodeID := fmt.Sprintf("event_%s", eventID)
	if err := kg.AddNodeCat(CategoryBusiness, eventNodeID, "Event", "Micro-stop", map[string]interface{}{
		"duration_seconds": duration,
		"type":             "microstop",
		"work_center":      workCenter,
	}); err != nil {
		return err
	}
	machineNodeID := fmt.Sprintf("machine_%s", workCenter)
	if err := kg.AddNodeCat(CategoryBusiness, machineNodeID, "Equipment", workCenter, map[string]interface{}{
		"work_center": workCenter,
	}); err != nil {
		return err
	}
	edgeID := fmt.Sprintf("edge_%s_%s", eventID, workCenter)
	return kg.AddEdgeCat(CategoryBusiness, edgeID, eventNodeID, machineNodeID, "occurred_at", 1.0)
}

// AddCause — business enrichment
func (kg *KnowledgeGraph) AddCause(eventID, cause string, confidence float64) error {
	causeNodeID := fmt.Sprintf("cause_%s", cause)
	if err := kg.AddNodeCat(CategoryBusiness, causeNodeID, "Cause", cause, map[string]interface{}{
		"confidence": confidence,
	}); err != nil {
		return err
	}
	eventNodeID := fmt.Sprintf("event_%s", eventID)
	edgeID := fmt.Sprintf("edge_%s_%s", eventID, cause)
	return kg.AddEdgeCat(CategoryBusiness, edgeID, eventNodeID, causeNodeID, "caused_by", confidence)
}

// AddCost — business enrichment
func (kg *KnowledgeGraph) AddCost(eventID string, costEUR float64) error {
	costNodeID := fmt.Sprintf("cost_%s", eventID)
	if err := kg.AddNodeCat(CategoryBusiness, costNodeID, "Cost", fmt.Sprintf("%.2f€", costEUR), map[string]interface{}{
		"amount_eur": costEUR,
	}); err != nil {
		return err
	}
	eventNodeID := fmt.Sprintf("event_%s", eventID)
	edgeID := fmt.Sprintf("edge_cost_%s", eventID)
	return kg.AddEdgeCat(CategoryBusiness, edgeID, eventNodeID, costNodeID, "costs", 1.0)
}

// GetFullGraph returns the entire graph (both categories). Legacy alias.
func (kg *KnowledgeGraph) GetFullGraph() (*GraphJSON, error) {
	return kg.GetGraph("all")
}

// GetGraph — unified read. category = "business" | "platform" | "all" | "" (=all).
func (kg *KnowledgeGraph) GetGraph(category string) (*GraphJSON, error) {
	if category == "" {
		category = "all"
	}
	nodeQuery := `SELECT id, category, type, label, properties, created_at FROM kg_nodes`
	edgeQuery := `SELECT id, category, from_id, to_id, relation, weight, created_at FROM kg_edges`
	var nodeArgs, edgeArgs []interface{}
	if category != "all" {
		nodeQuery += ` WHERE category = ?`
		edgeQuery += ` WHERE category = ?`
		nodeArgs = append(nodeArgs, category)
		edgeArgs = append(edgeArgs, category)
	}

	// Nodes
	rows, err := kg.store.DB().Query(nodeQuery, nodeArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nodes []Node
	for rows.Next() {
		var n Node
		var propsJSON string
		if err := rows.Scan(&n.ID, &n.Category, &n.Type, &n.Label, &propsJSON, &n.CreatedAt); err != nil {
			continue
		}
		if propsJSON != "" {
			if err := json.Unmarshal([]byte(propsJSON), &n.Properties); err != nil {
				log.Printf("[KG] Warning: bad properties JSON on node %s: %v", n.ID, err)
			}
		}
		nodes = append(nodes, n)
	}

	// Edges
	edgeRows, err := kg.store.DB().Query(edgeQuery, edgeArgs...)
	if err != nil {
		return nil, err
	}
	defer edgeRows.Close()
	var edges []Edge
	for edgeRows.Next() {
		var e Edge
		if err := edgeRows.Scan(&e.ID, &e.Category, &e.FromID, &e.ToID, &e.Relation, &e.Weight, &e.CreatedAt); err != nil {
			continue
		}
		edges = append(edges, e)
	}

	return &GraphJSON{Nodes: nodes, Edges: edges}, nil
}

// RepopulatePlatform wipes existing platform nodes/edges and rebuilds them from
// the current pipeline registry. Call this whenever pipelines change.
// This replaces the old 5-min-cached in-memory Technical KG.
func (kg *KnowledgeGraph) RepopulatePlatform(reg *pipeline.Registry) error {
	kg.platformMu.Lock()
	defer kg.platformMu.Unlock()

	// If nothing changed since last rebuild, no-op (fast path).
	currentHash := reg.GetHash()
	if currentHash == kg.lastRegHash && !kg.lastRebuilt.IsZero() {
		return nil
	}

	// Wipe existing platform rows.
	if _, err := kg.store.DB().Exec(`DELETE FROM kg_edges WHERE category = 'platform'`); err != nil {
		return err
	}
	if _, err := kg.store.DB().Exec(`DELETE FROM kg_nodes WHERE category = 'platform'`); err != nil {
		return err
	}

	// Build in-memory using the existing algorithm, then persist.
	builder := NewTechnicalBuilder(reg)
	graph := builder.Build()

	for _, n := range graph.Nodes {
		var propsJSON string
		if n.Properties != nil {
			b, _ := json.Marshal(n.Properties)
			propsJSON = string(b)
		}
		_, err := kg.store.DB().Exec(
			`INSERT OR REPLACE INTO kg_nodes (id, category, type, label, properties) VALUES (?, ?, ?, ?, ?)`,
			n.ID, string(CategoryPlatform), string(n.Type), n.Name, propsJSON,
		)
		if err != nil {
			log.Printf("[KG] Failed to insert platform node %s: %v", n.ID, err)
		}
	}
	for _, e := range graph.Edges {
		_, err := kg.store.DB().Exec(
			`INSERT OR REPLACE INTO kg_edges (id, category, from_id, to_id, relation, weight) VALUES (?, ?, ?, ?, ?, ?)`,
			e.ID, string(CategoryPlatform), e.From, e.To, string(e.Type), e.Weight,
		)
		if err != nil {
			log.Printf("[KG] Failed to insert platform edge %s: %v", e.ID, err)
		}
	}

	kg.lastRebuilt = time.Now()
	kg.lastRegHash = currentHash
	log.Printf("[KG] Platform sub-graph rebuilt (%d nodes, %d edges)", len(graph.Nodes), len(graph.Edges))
	return nil
}

// ─── Legacy API preserved for existing callers ──────────────────────────────

// Store returns the underlying storage.
func (kg *KnowledgeGraph) Store() *storage.SQLiteStore {
	return kg.store
}

// GetTechnicalGraph — legacy alias. Now sourced from the unified store.
// Callers should migrate to GetGraph("platform"). This wrapper adapts to the
// old TechnicalGraph shape for backward compat with cmd/agent code.
func (kg *KnowledgeGraph) GetTechnicalGraph(reg *pipeline.Registry) (*TechnicalGraph, error) {
	if reg != nil {
		if err := kg.RepopulatePlatform(reg); err != nil {
			log.Printf("[KG] Warning: platform repopulate failed: %v", err)
		}
	}
	g, err := kg.GetGraph("platform")
	if err != nil {
		return nil, err
	}
	tg := &TechnicalGraph{}
	for _, n := range g.Nodes {
		tg.Nodes = append(tg.Nodes, TechnicalNode{
			ID:         n.ID,
			Type:       TechnicalNodeType(n.Type),
			Name:       n.Label,
			Properties: n.Properties,
			CreatedAt:  n.CreatedAt,
		})
	}
	for _, e := range g.Edges {
		tg.Edges = append(tg.Edges, TechnicalEdge{
			ID:     e.ID,
			From:   e.FromID,
			To:     e.ToID,
			Type:   TechnicalEdgeType(e.Relation),
			Weight: e.Weight,
		})
	}
	return tg, nil
}

// GetTechnicalGraphWithCache — legacy alias (cache logic removed since the
// unified store IS the source of truth; RepopulatePlatform is idempotent).
func (kg *KnowledgeGraph) GetTechnicalGraphWithCache(reg *pipeline.Registry) (*TechnicalGraph, error) {
	return kg.GetTechnicalGraph(reg)
}

// PurgeCache — legacy no-op. Cache is gone; use RepopulatePlatform instead.
func (kg *KnowledgeGraph) PurgeCache() {
	kg.lastRegHash = ""
	kg.lastRebuilt = time.Time{}
	log.Printf("[KG] Platform rebuild flag reset (next request will rebuild)")
}

// Close ferme la base
func (kg *KnowledgeGraph) Close() error {
	return kg.store.Close()
}
