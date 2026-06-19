// internal/kg/builder.go
package kg

import (
	"fmt"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
)

// TechnicalBuilder construit le graphe technique automatiquement
type TechnicalBuilder struct {
	pipelineReg *pipeline.Registry
	nodes       map[string]TechnicalNode
	edges       []TechnicalEdge
}

// NewTechnicalBuilder crée un nouveau builder
func NewTechnicalBuilder(pipelineReg *pipeline.Registry) *TechnicalBuilder {
	return &TechnicalBuilder{
		pipelineReg: pipelineReg,
		nodes:       make(map[string]TechnicalNode),
		edges:       make([]TechnicalEdge, 0),
	}
}

// Build construit le graphe technique complet
func (b *TechnicalBuilder) Build() *TechnicalGraph {
	// 1. Ajouter tous les pipelines
	b.addPipelines()

	// 2. Ajouter toutes les fonctions
	b.addFunctions()

	// 3. Ajouter les topics
	b.addTopics()

	// 4. Ajouter les connexions
	b.addConnections()

	// 5. Ajouter les relations
	b.addRelations()

	return &TechnicalGraph{
		Nodes: b.nodesToList(),
		Edges: b.edges,
	}
}

// addPipelines ajoute les pipelines comme nœuds
func (b *TechnicalBuilder) addPipelines() {
	for _, p := range b.pipelineReg.List() {
		node := TechnicalNode{
			ID:   fmt.Sprintf("pipeline_%s", p.ID),
			Type: TechNodePipeline,
			Name: p.Name,
			Properties: map[string]interface{}{
				"id":          p.ID,
				"description": p.Description,
				"version":     p.Version,
				"enabled":     true,
			},
			CreatedAt: time.Now(),
		}
		b.nodes[node.ID] = node
	}
}

// addFunctions ajoute les fonctions utilisées
func (b *TechnicalBuilder) addFunctions() {
	functionsSeen := make(map[string]bool)

	for _, p := range b.pipelineReg.List() {
		for _, node := range p.Nodes {
			fnName := node.Function
			if functionsSeen[fnName] {
				continue
			}
			functionsSeen[fnName] = true

			nodeID := fmt.Sprintf("function_%s", fnName)
			b.nodes[nodeID] = TechnicalNode{
				ID:   nodeID,
				Type: TechNodeFunction,
				Name: fnName,
				Properties: map[string]interface{}{
					"type": node.Type,
				},
				CreatedAt: time.Now(),
			}
		}
	}
}

// addTopics ajoute les topics MQTT
func (b *TechnicalBuilder) addTopics() {
	topics := map[string]bool{
		"mindset/raw/#":                  true,
		"mindset/site/#":                 true,
		"mindset/events/#":               true,
		"mindset/events/micro-stop":      true,
		"mindset/events/micro-stop-cost": true,
	}

	for topic := range topics {
		nodeID := fmt.Sprintf("topic_%s", sanitizeID(topic))
		b.nodes[nodeID] = TechnicalNode{
			ID:   nodeID,
			Type: TechNodeTopic,
			Name: topic,
			Properties: map[string]interface{}{
				"qos_default": 1,
			},
			CreatedAt: time.Now(),
		}
	}
}

// addConnections ajoute les connexions
func (b *TechnicalBuilder) addConnections() {
	connections := []struct {
		id    string
		name  string
		props map[string]interface{}
	}{
		{
			id:   "conn_opcua",
			name: "OPC-UA Connection",
			props: map[string]interface{}{
				"protocol": "opcua",
				"endpoint": "opc.tcp://localhost:4840",
			},
		},
		{
			id:   "conn_mqtt",
			name: "MQTT Broker",
			props: map[string]interface{}{
				"protocol": "mqtt",
				"broker":   "tcp://localhost:1883",
			},
		},
	}

	for _, conn := range connections {
		b.nodes[conn.id] = TechnicalNode{
			ID:         conn.id,
			Type:       TechNodeConnection,
			Name:       conn.name,
			Properties: conn.props,
			CreatedAt:  time.Now(),
		}
	}
}

// addRelations ajoute toutes les relations entre nœuds
func (b *TechnicalBuilder) addRelations() {
	// Relation: OPC-UA Connection → Topic raw
	b.addEdge("conn_opcua", "topic_mindset_raw_", EdgePublishesTo, 1.0)

	// Relation: MQTT Connection → Topic raw
	b.addEdge("conn_mqtt", "topic_mindset_raw_", EdgeSubscribesTo, 1.0)

	for _, p := range b.pipelineReg.List() {
		pipelineID := fmt.Sprintf("pipeline_%s", p.ID)

		// Relation: Pipeline → Fonctions
		for _, node := range p.Nodes {
			functionID := fmt.Sprintf("function_%s", node.Function)
			b.addEdge(pipelineID, functionID, EdgeDependsOn, 1.0)
		}

		// Relation: Trigger → Pipeline
		triggerTopic := "mindset/raw/#"
		if topic, ok := p.Trigger.Config["topic"].(string); ok {
			triggerTopic = topic
		}
		topicID := fmt.Sprintf("topic_%s", sanitizeID(triggerTopic))
		b.addEdge(topicID, pipelineID, EdgeTriggers, 1.0)

		// Relation: Pipeline → Output Topic
		for _, node := range p.Nodes {
			if node.Type == "output" {
				if topic, ok := node.Config["topic"].(string); ok {
					topicID := fmt.Sprintf("topic_%s", sanitizeID(topic))
					b.addEdge(pipelineID, topicID, EdgeProduces, 1.0)
				}
				if prefix, ok := node.Config["topic_prefix"].(string); ok {
					topicID := fmt.Sprintf("topic_%s", sanitizeID(prefix+"/#"))
					b.addEdge(pipelineID, topicID, EdgeProduces, 1.0)
				}
			}
		}
	}
}

// addEdge ajoute une relation (évite les doublons)
func (b *TechnicalBuilder) addEdge(from, to string, edgeType TechnicalEdgeType, weight float64) {
	// Vérifier si la relation existe déjà
	for _, e := range b.edges {
		if e.From == from && e.To == to && e.Type == edgeType {
			return
		}
	}

	b.edges = append(b.edges, TechnicalEdge{
		ID:     fmt.Sprintf("edge_%s_%s", from, to),
		From:   from,
		To:     to,
		Type:   edgeType,
		Weight: weight,
	})
}

// nodesToList convertit la map en slice
func (b *TechnicalBuilder) nodesToList() []TechnicalNode {
	nodes := make([]TechnicalNode, 0, len(b.nodes))
	for _, node := range b.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// sanitizeID nettoie un string pour l'utiliser comme ID
func sanitizeID(s string) string {
	result := ""
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result += string(c)
		} else if c == '/' || c == '_' || c == '-' {
			result += "_"
		}
	}
	return result
}
