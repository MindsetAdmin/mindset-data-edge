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

// asFloat64 handles both float64 and int — see the identical helper's
// comment in internal/functions/calculates/cost.go for why: yaml.v3 decodes
// a plain integer literal in a pipeline YAML's config as Go `int`, not
// `float64`, so a bare `.(float64)` assertion silently misses it.
func asFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	}
	return 0, false
}

// Execute vérifie le seuil
func (h *ThresholdHandler) Execute(durationSeconds float64, config map[string]interface{}) (*ThresholdResult, error) {
	min := 30.0  // défaut 30 secondes
	max := 180.0 // défaut 3 minutes

	if minVal, ok := asFloat64(config["min"]); ok {
		min = minVal
	}
	if maxVal, ok := asFloat64(config["max"]); ok {
		max = maxVal
	}

	isInRange := durationSeconds >= min && durationSeconds <= max

	return &ThresholdResult{
		Value:       durationSeconds,
		IsInRange:   isInRange,
		IsMicroStop: isInRange,
	}, nil
}
