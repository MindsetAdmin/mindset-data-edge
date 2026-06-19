// internal/pipeline/builder.go
package pipeline

import (
	"fmt"
	"time"
)

// Builder permet de construire des pipelines programmatiquement
type Builder struct {
	pipeline *Pipeline
}

// NewBuilder crée un nouveau builder
func NewBuilder(id, name string) *Builder {
	return &Builder{
		pipeline: &Pipeline{
			ID:        id,
			Name:      name,
			Version:   "1.0",
			Nodes:     make([]Node, 0),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// WithDescription ajoute une description
func (b *Builder) WithDescription(desc string) *Builder {
	b.pipeline.Description = desc
	return b
}

// WithVersion ajoute une version
func (b *Builder) WithVersion(version string) *Builder {
	b.pipeline.Version = version
	return b
}

// WithTrigger ajoute un trigger
func (b *Builder) WithTrigger(triggerType, function string, config map[string]interface{}) *Builder {
	b.pipeline.Trigger = Trigger{
		Type:     triggerType,
		Function: function,
		Config:   config,
	}
	return b
}

// AddNode ajoute un nœud
func (b *Builder) AddNode(id, name string, nodeType NodeType, function string, config map[string]interface{}, dependsOn []string) *Builder {
	node := Node{
		ID:        id,
		Name:      name,
		Type:      nodeType,
		Function:  function,
		Config:    config,
		DependsOn: dependsOn,
	}
	b.pipeline.Nodes = append(b.pipeline.Nodes, node)
	return b
}

// WithOutput définit le nœud de sortie
func (b *Builder) WithOutput(nodeID string) *Builder {
	b.pipeline.Output = nodeID
	return b
}

// Build retourne le pipeline construit
func (b *Builder) Build() (*Pipeline, error) {
	if b.pipeline.ID == "" {
		return nil, fmt.Errorf("pipeline ID is required")
	}
	if b.pipeline.Name == "" {
		return nil, fmt.Errorf("pipeline name is required")
	}
	if len(b.pipeline.Nodes) == 0 {
		return nil, fmt.Errorf("pipeline must have at least one node")
	}
	return b.pipeline, nil
}
