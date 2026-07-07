// internal/kg/types.go
// Unified KG type definitions (2026-07-02 merge — see docs/analysis_log.md Entry 50).
// The old TechnicalNode/TechnicalEdge/TechnicalGraph structs are preserved for
// backward compatibility with builder.go's transient in-memory shape, but the
// canonical persistence model is now the unified Node/Edge with Category (in graph.go).
package kg

import "time"

// Category tags every KG node/edge as either operational-site data or platform-wiring data.
type Category string

const (
	CategoryBusiness Category = "business" // Equipment, Event, Cause, Cost, Operator, OF, Product, etc.
	CategoryPlatform Category = "platform" // Pipeline, Function, Connection, Topic, Dashboard
)

// ─── Business node type strings (used in the free-form Type field) ──────────
const (
	TypeEquipment = "Equipment"
	TypeEvent     = "Event"
	TypeCause     = "Cause"
	TypeCost      = "Cost"
	TypeOperator  = "Operator"
	TypeProduct   = "Product"
	TypeOF        = "OF"
)

// ─── Platform node type constants (also used as free-form Type strings) ─────
type TechnicalNodeType string

const (
	TechNodeConnection TechnicalNodeType = "connection"
	TechNodeTopic      TechnicalNodeType = "topic"
	TechNodeFunction   TechnicalNodeType = "function"
	TechNodePipeline   TechnicalNodeType = "pipeline"
	TechNodeDashboard  TechnicalNodeType = "dashboard"
)

// TechnicalEdgeType — labels for platform edges.
type TechnicalEdgeType string

const (
	EdgePublishesTo  TechnicalEdgeType = "publishes_to"
	EdgeSubscribesTo TechnicalEdgeType = "subscribes_to"
	EdgeDependsOn    TechnicalEdgeType = "depends_on"
	EdgeTriggers     TechnicalEdgeType = "triggers"
	EdgeConsumes     TechnicalEdgeType = "consumes"
	EdgeProduces     TechnicalEdgeType = "produces"
)

// ─── Legacy in-memory shapes (only used by the builder as an intermediate) ──
// These are kept because the builder.go algorithm still works with them as
// scratch space, then persists via the KnowledgeGraph.AddNode/AddEdge API.
type TechnicalNode struct {
	ID         string                 `json:"id"`
	Type       TechnicalNodeType      `json:"type"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
}

type TechnicalEdge struct {
	ID     string            `json:"id"`
	From   string            `json:"from"`
	To     string            `json:"to"`
	Type   TechnicalEdgeType `json:"type"`
	Weight float64           `json:"weight"`
}

type TechnicalGraph struct {
	Nodes []TechnicalNode `json:"nodes"`
	Edges []TechnicalEdge `json:"edges"`
}
