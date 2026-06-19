// internal/pipeline/types.go
package pipeline

import (
	"time"
)

// NodeType définit le type de nœud dans un pipeline
type NodeType string

const (
	NodeTypeTrigger   NodeType = "trigger"
	NodeTypeConnector NodeType = "connector"
	NodeTypeTransform NodeType = "transform"
	NodeTypeCalculate NodeType = "calculate"
	NodeTypeCondition NodeType = "condition"
	NodeTypeOutput    NodeType = "output"
)

// Node représente une étape dans le pipeline
type Node struct {
	ID        string                 `json:"id" yaml:"id"`
	Name      string                 `json:"name" yaml:"name"`
	Type      NodeType               `json:"type" yaml:"type"`
	Function  string                 `json:"function" yaml:"function"` // nom de la Function à appeler
	Config    map[string]interface{} `json:"config" yaml:"config"`
	DependsOn []string               `json:"depends_on" yaml:"depends_on"` // IDs des nœuds précédents
}

// Trigger représente le déclencheur du pipeline
type Trigger struct {
	Type     string                 `json:"type" yaml:"type"`         // timer, mqtt, webhook
	Function string                 `json:"function" yaml:"function"` // fonction de déclenchement
	Config   map[string]interface{} `json:"config" yaml:"config"`
}

// Pipeline représente un pipeline complet
type Pipeline struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Version     string                 `json:"version" yaml:"version"`
	Enabled     bool                   `json:"enabled" yaml:"enabled"`
	Properties  map[string]interface{} `json:"properties" yaml:"properties"`
	Trigger     Trigger                `json:"trigger" yaml:"trigger"`
	Nodes       []Node                 `json:"nodes" yaml:"nodes"`
	Output      string                 `json:"output" yaml:"output"` // ID du nœud de sortie
	CreatedAt   time.Time              `json:"created_at" yaml:"-"`  // runtime-only, never persisted to YAML
	UpdatedAt   time.Time              `json:"updated_at" yaml:"-"`  // runtime-only, never persisted to YAML
}

// ExecutionStatus représente le statut d'exécution
type ExecutionStatus string

const (
	StatusPending ExecutionStatus = "pending"
	StatusRunning ExecutionStatus = "running"
	StatusSuccess ExecutionStatus = "success"
	StatusFailed  ExecutionStatus = "failed"
	StatusSkipped ExecutionStatus = "skipped"
)

// NodeExecutionResult résultat d'exécution d'un nœud
type NodeExecutionResult struct {
	NodeID    string                 `json:"node_id"`
	Status    ExecutionStatus        `json:"status"`
	Input     map[string]interface{} `json:"input,omitempty"`
	Output    interface{}            `json:"output,omitempty"`
	Error     string                 `json:"error,omitempty"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Duration  time.Duration          `json:"duration_ms"`
}

// ExecutionResult résultat complet d'un pipeline
type ExecutionResult struct {
	PipelineID    string                 `json:"pipeline_id"`
	PipelineName  string                 `json:"pipeline_name"`
	Status        ExecutionStatus        `json:"status"`
	Trigger       map[string]interface{} `json:"trigger,omitempty"`
	Nodes         []NodeExecutionResult  `json:"nodes"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	TotalDuration time.Duration          `json:"total_duration_ms"`
	Error         string                 `json:"error,omitempty"`
}

// PipelineEvent événement émis pendant l'exécution
type PipelineEvent struct {
	PipelineID string      `json:"pipeline_id"`
	NodeID     string      `json:"node_id,omitempty"`
	Type       string      `json:"type"` // started, node_started, node_completed, completed, failed
	Data       interface{} `json:"data,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}
