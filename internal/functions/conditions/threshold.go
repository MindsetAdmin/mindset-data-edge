package conditions

import (
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// ThresholdConfig configuration
type ThresholdConfig struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// ThresholdResult résultat
type ThresholdResult struct {
	Value       float64 `json:"value"`
	IsInRange   bool    `json:"is_in_range"`
	IsMicroStop bool    `json:"is_micro_stop"`
}

// ThresholdHandler handler
type ThresholdHandler struct{}

// NewThresholdHandler crée un nouveau handler
func NewThresholdHandler() *ThresholdHandler {
	return &ThresholdHandler{}
}

// GetFunction retourne la définition
func (h *ThresholdHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "threshold",
		Type:        functions.TypeCondition,
		Description: "Vérifie si une valeur est entre un min et un max",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			durationSeconds, _ := params["duration_seconds"].(float64)
			return h.Execute(durationSeconds, params)
		},
	}
}

// Execute vérifie le seuil
func (h *ThresholdHandler) Execute(durationSeconds float64, config map[string]interface{}) (*ThresholdResult, error) {
	min := 30.0  // défaut 30 secondes
	max := 180.0 // défaut 3 minutes

	if minVal, ok := config["min"].(float64); ok {
		min = minVal
	}
	if maxVal, ok := config["max"].(float64); ok {
		max = maxVal
	}

	isInRange := durationSeconds >= min && durationSeconds <= max

	return &ThresholdResult{
		Value:       durationSeconds,
		IsInRange:   isInRange,
		IsMicroStop: isInRange,
	}, nil
}
