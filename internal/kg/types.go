// internal/kg/types.go (à compléter)
package kg

import "time"

// TechnicalNodeType représente les types de nœuds techniques
type TechnicalNodeType string

const (
	TechNodeConnection TechnicalNodeType = "connection"
	TechNodeTopic      TechnicalNodeType = "topic"
	TechNodeFunction   TechnicalNodeType = "function"
	TechNodePipeline   TechnicalNodeType = "pipeline"
	TechNodeDashboard  TechnicalNodeType = "dashboard"
)

// TechnicalEdgeType représente les types de relations techniques
type TechnicalEdgeType string

const (
	EdgePublishesTo  TechnicalEdgeType = "publishes_to"
	EdgeSubscribesTo TechnicalEdgeType = "subscribes_to"
	EdgeDependsOn    TechnicalEdgeType = "depends_on"
	EdgeTriggers     TechnicalEdgeType = "triggers"
	EdgeConsumes     TechnicalEdgeType = "consumes"
	EdgeProduces     TechnicalEdgeType = "produces"
)

// TechnicalNode nœud technique du KG
type TechnicalNode struct {
	ID         string                 `json:"id"`
	Type       TechnicalNodeType      `json:"type"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
}

// TechnicalEdge relation technique
type TechnicalEdge struct {
	ID     string            `json:"id"`
	From   string            `json:"from"`
	To     string            `json:"to"`
	Type   TechnicalEdgeType `json:"type"`
	Weight float64           `json:"weight"`
}

// TechnicalGraph graphe technique complet
type TechnicalGraph struct {
	Nodes []TechnicalNode `json:"nodes"`
	Edges []TechnicalEdge `json:"edges"`
}
