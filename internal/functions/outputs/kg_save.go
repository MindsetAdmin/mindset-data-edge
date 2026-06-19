package outputs

import (
	"fmt"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
)

// KGSaveConfig configuration
type KGSaveConfig struct {
	NodeType   string                 `json:"node_type"`
	Labels     map[string]string      `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

// KGSaveHandler handler
type KGSaveHandler struct {
	kg *kg.KnowledgeGraph
}

// NewKGSaveHandler crée un nouveau handler
func NewKGSaveHandler(kgInstance *kg.KnowledgeGraph) *KGSaveHandler {
	return &KGSaveHandler{
		kg: kgInstance,
	}
}

// GetFunction retourne la définition
func (h *KGSaveHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "kg_save",
		Type:        functions.TypeOutput,
		Description: "Sauvegarde un événement dans le Knowledge Graph",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			eventID, _ := params["event_id"].(string)
			workCenter, _ := params["work_center"].(string)
			durationSeconds, _ := params["duration_seconds"].(float64)
			cause, _ := params["cause"].(string)
			costEUR, _ := params["cost_eur"].(float64)
			return nil, h.Execute(eventID, workCenter, durationSeconds, cause, costEUR, params)
		},
	}
}

// Execute sauvegarde dans le KG
func (h *KGSaveHandler) Execute(eventID, workCenter string, durationSeconds float64, cause string, costEUR float64, config map[string]interface{}) error {
	if h.kg == nil {
		return fmt.Errorf("knowledge graph not initialized")
	}

	// Rien à enregistrer si l'événement est vide (ex: exécution manuelle sans
	// données de trigger) — évite de polluer le graphe avec des nœuds vides.
	if eventID == "" {
		return nil
	}

	// Ajouter le micro-stop au KG
	if err := h.kg.AddMicroStop(eventID, workCenter, durationSeconds); err != nil {
		return fmt.Errorf("failed to add micro-stop: %w", err)
	}

	// Ajouter la cause si présente
	if cause != "" {
		confidence := 0.9
		if conf, ok := config["confidence"].(float64); ok {
			confidence = conf
		}
		if err := h.kg.AddCause(eventID, cause, confidence); err != nil {
			return fmt.Errorf("failed to add cause: %w", err)
		}
	}

	// Ajouter le coût
	if costEUR > 0 {
		if err := h.kg.AddCost(eventID, costEUR); err != nil {
			return fmt.Errorf("failed to add cost: %w", err)
		}
	}

	return nil
}
