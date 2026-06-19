// internal/kg/graph.go
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


// Node représente un nœud dans le Knowledge Graph
type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// Edge représente une relation entre deux nœuds
type Edge struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Relation  string    `json:"relation"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}

// GraphJSON pour Cytoscape
type GraphJSON struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// CachedTechnicalGraph stocke le graphe avec un timestamp
type CachedTechnicalGraph struct {
	Graph     *TechnicalGraph
	Generated time.Time
	CacheKey  string // Hash des pipelines pour détecter les changements
}

// KnowledgeGraph gère le graphe de connaissances
type KnowledgeGraph struct {
	store       *storage.SQLiteStore
	cachedGraph *CachedTechnicalGraph
	cacheMu     sync.RWMutex
}

// NewKnowledgeGraph crée un nouveau KG
func NewKnowledgeGraph(dbPath string) (*KnowledgeGraph, error) {
	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}
	return &KnowledgeGraph{store: store}, nil
}

// AddNode ajoute un nœud (ignore si existe déjà)
func (kg *KnowledgeGraph) AddNode(id, nodeType, label string, props map[string]interface{}) error {
	var propsJSON string
	if props != nil {
		jsonBytes, err := json.Marshal(props)
		if err != nil {
			return err
		}
		propsJSON = string(jsonBytes)
	}

	query := `INSERT OR IGNORE INTO kg_nodes (id, type, label, properties) VALUES (?, ?, ?, ?)`
	_, err := kg.store.DB().Exec(query, id, nodeType, label, propsJSON)
	if err != nil {
		return err
	}
	log.Printf("[KG] Added node: %s (%s)", label, nodeType)
	return nil
}

// AddEdge ajoute une relation
func (kg *KnowledgeGraph) AddEdge(id, fromID, toID, relation string, weight float64) error {
	query := `INSERT INTO kg_edges (id, from_id, to_id, relation, weight) VALUES (?, ?, ?, ?, ?)`
	_, err := kg.store.DB().Exec(query, id, fromID, toID, relation, weight)
	if err != nil {
		return err
	}
	log.Printf("[KG] Added edge: %s → %s (%s)", fromID, toID, relation)
	return nil
}

// AddMicroStop enrichit le KG avec un micro-stop
func (kg *KnowledgeGraph) AddMicroStop(eventID, workCenter string, duration float64) error {
	// Nœud événement
	eventNodeID := fmt.Sprintf("event_%s", eventID)
	err := kg.AddNode(eventNodeID, "Event", "Micro-stop", map[string]interface{}{
		"duration_seconds": duration,
		"type":             "microstop",
		"work_center":      workCenter,
	})
	if err != nil {
		return err
	}

	// Nœud machine (si n'existe pas déjà)
	machineNodeID := fmt.Sprintf("machine_%s", workCenter)
	err = kg.AddNode(machineNodeID, "Equipment", workCenter, map[string]interface{}{
		"work_center": workCenter,
	})
	if err != nil {
		return err
	}

	// Relation: événement → machine
	edgeID := fmt.Sprintf("edge_%s_%s", eventID, workCenter)
	err = kg.AddEdge(edgeID, eventNodeID, machineNodeID, "occurred_at", 1.0)
	if err != nil {
		return err
	}

	return nil
}

// AddCause ajoute une cause à un micro-stop
func (kg *KnowledgeGraph) AddCause(eventID, cause string, confidence float64) error {
	// Nœud cause
	causeNodeID := fmt.Sprintf("cause_%s", cause)
	err := kg.AddNode(causeNodeID, "Cause", cause, map[string]interface{}{
		"confidence": confidence,
	})
	if err != nil {
		return err
	}

	// Relation: événement → cause
	eventNodeID := fmt.Sprintf("event_%s", eventID)
	edgeID := fmt.Sprintf("edge_%s_%s", eventID, cause)
	err = kg.AddEdge(edgeID, eventNodeID, causeNodeID, "caused_by", confidence)
	if err != nil {
		return err
	}

	return nil
}

// AddCost ajoute un coût à un événement
func (kg *KnowledgeGraph) AddCost(eventID string, costEUR float64) error {
	// Nœud coût
	costNodeID := fmt.Sprintf("cost_%s", eventID)
	err := kg.AddNode(costNodeID, "Cost", fmt.Sprintf("%.2f€", costEUR), map[string]interface{}{
		"amount_eur": costEUR,
	})
	if err != nil {
		return err
	}

	// Relation: événement → coût
	eventNodeID := fmt.Sprintf("event_%s", eventID)
	edgeID := fmt.Sprintf("edge_cost_%s", eventID)
	err = kg.AddEdge(edgeID, eventNodeID, costNodeID, "costs", 1.0)
	if err != nil {
		return err
	}

	return nil
}

// GetFullGraph retourne tout le graphe pour Cytoscape
func (kg *KnowledgeGraph) GetFullGraph() (*GraphJSON, error) {
	// Récupérer tous les nœuds
	rows, err := kg.store.DB().Query(`
		SELECT id, type, label, properties, created_at FROM kg_nodes
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []Node
	for rows.Next() {
		var n Node
		var propsJSON string
		err := rows.Scan(&n.ID, &n.Type, &n.Label, &propsJSON, &n.CreatedAt)
		if err != nil {
			continue
		}
		if propsJSON != "" {
			if err := json.Unmarshal([]byte(propsJSON), &n.Properties); err != nil {
				log.Printf("[KG] Warning: failed to unmarshal properties: %v", err)
			}
		}
		nodes = append(nodes, n)
	}

	// Récupérer toutes les relations
	edgeRows, err := kg.store.DB().Query(`
		SELECT id, from_id, to_id, relation, weight, created_at FROM kg_edges
	`)
	if err != nil {
		return nil, err
	}
	defer edgeRows.Close()

	var edges []Edge
	for edgeRows.Next() {
		var e Edge
		err := edgeRows.Scan(&e.ID, &e.FromID, &e.ToID, &e.Relation, &e.Weight, &e.CreatedAt)
		if err != nil {
			continue
		}
		edges = append(edges, e)
	}

	return &GraphJSON{Nodes: nodes, Edges: edges}, nil
}

// Store retourne le storage sous-jacent
func (kg *KnowledgeGraph) Store() *storage.SQLiteStore {
	return kg.store
}

// GetTechnicalGraph retourne le graphe technique complet
func (kg *KnowledgeGraph) GetTechnicalGraph(pipelineReg *pipeline.Registry) (*TechnicalGraph, error) {
	builder := NewTechnicalBuilder(pipelineReg)
	return builder.Build(), nil
}

// GetTechnicalGraphWithCache retourne le graphe depuis le cache ou le recalcule
func (kg *KnowledgeGraph) GetTechnicalGraphWithCache(reg *pipeline.Registry) (*TechnicalGraph, error) {
	currentKey := reg.GetHash()

	kg.cacheMu.RLock()
	cached := kg.cachedGraph
	kg.cacheMu.RUnlock()

	if cached != nil && cached.CacheKey == currentKey && time.Since(cached.Generated) < 5*time.Minute {
		log.Printf("[KG] Returning cached graph (%d nodes, %d edges)",
			len(cached.Graph.Nodes), len(cached.Graph.Edges))
		return cached.Graph, nil
	}
	log.Printf("[KG] Cache miss, regenerating...")

	builder := NewTechnicalBuilder(reg)
	graph := builder.Build()

	kg.cacheMu.Lock()
	kg.cachedGraph = &CachedTechnicalGraph{
		Graph:     graph,
		Generated: time.Now(),
		CacheKey:  currentKey,
	}
	kg.cacheMu.Unlock()

	log.Printf("[KG] Cache updated (%d nodes, %d edges)", len(graph.Nodes), len(graph.Edges))
	return graph, nil
}

// PurgeCache vide le cache du graphe technique
func (kg *KnowledgeGraph) PurgeCache() {
	kg.cacheMu.Lock()
	kg.cachedGraph = nil
	kg.cacheMu.Unlock()
	log.Printf("[KG] Cache purged")
}

// Close ferme la base
func (kg *KnowledgeGraph) Close() error {
	return kg.store.Close()
}
