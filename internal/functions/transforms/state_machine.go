package transforms

import (
	"sync"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// StateMachineConfig configuration
type StateMachineConfig struct {
	InitialState bool `json:"initial_state"`
}

// StateTransition représente une transition d'état
type StateTransition struct {
	From      bool      `json:"from"`
	To        bool      `json:"to"`
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration_seconds,omitempty"`
}

// StateMachineHandler handler
type StateMachineHandler struct {
	states map[string]bool
	mu     sync.RWMutex
}

// NewStateMachineHandler crée un nouveau handler
func NewStateMachineHandler() *StateMachineHandler {
	return &StateMachineHandler{
		states: make(map[string]bool),
	}
}

// GetFunction retourne la définition
func (h *StateMachineHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "state_machine",
		Type:        functions.TypeTransform,
		Description: "Détecte les transitions d'état (Run/Stop)",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			currentValue, _ := params["current_value"].(bool)
			machineID, _ := params["machine_id"].(string)
			return h.Execute(currentValue, machineID)
		},
	}
}

// Execute détecte la transition
func (h *StateMachineHandler) Execute(currentValue bool, machineID string) (*StateTransition, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	oldState, exists := h.states[machineID]

	// Premier état, juste enregistrer
	if !exists {
		h.states[machineID] = currentValue
		return nil, nil
	}

	// Pas de changement
	if oldState == currentValue {
		return nil, nil
	}

	// Transition détectée
	transition := &StateTransition{
		From:      oldState,
		To:        currentValue,
		Timestamp: time.Now(),
	}

	// Si c'est un redémarrage (Stop → Run), calculer la durée
	if !oldState && currentValue {
		// La durée sera calculée par la fonction duration
	}

	// Mettre à jour l'état
	h.states[machineID] = currentValue

	return transition, nil
}
