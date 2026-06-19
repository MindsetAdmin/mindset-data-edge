package transforms

import (
	"fmt"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// FilterConfig configuration du filtre
type FilterConfig struct {
	Field    string      `json:"field"`    // Champ à filtrer (ex: "data_type", "value")
	Operator string      `json:"operator"` // eq, ne, gt, lt, contains
	Value    interface{} `json:"value"`    // Valeur de comparaison
}

// FilterResult résultat du filtrage
type FilterResult struct {
	Passed bool        `json:"passed"`
	Data   interface{} `json:"data"`
	Reason string      `json:"reason,omitempty"`
}

// FilterHandler handler
type FilterHandler struct{}

// NewFilterHandler crée un nouveau handler
func NewFilterHandler() *FilterHandler {
	return &FilterHandler{}
}

// GetFunction retourne la définition
func (h *FilterHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "filter",
		Type:        functions.TypeTransform,
		Description: "Filtre les données selon une condition",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			data, _ := params["data"].(map[string]interface{})
			if data == nil {
				data = params
			}
			return h.Execute(data, params)
		},
	}
}

// Execute filtre les données
func (h *FilterHandler) Execute(data map[string]interface{}, config map[string]interface{}) (*FilterResult, error) {
	field, ok := config["field"].(string)
	if !ok {
		return nil, fmt.Errorf("missing field in config")
	}

	operator, ok := config["operator"].(string)
	if !ok {
		return nil, fmt.Errorf("missing operator in config")
	}

	expected, ok := config["value"]
	if !ok {
		return nil, fmt.Errorf("missing value in config")
	}

	// Récupérer la valeur réelle
	actual, exists := data[field]
	if !exists {
		return &FilterResult{
			Passed: false,
			Data:   data,
			Reason: fmt.Sprintf("field %s not found", field),
		}, nil
	}

	// Appliquer l'opérateur
	passed := false
	switch operator {
	case "eq":
		passed = fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected)
	case "ne":
		passed = fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected)
	case "gt":
		passed = compareNumbers(actual, expected) > 0
	case "lt":
		passed = compareNumbers(actual, expected) < 0
	case "contains":
		passed = containsString(actual, expected)
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
	}

	return &FilterResult{
		Passed: passed,
		Data:   data,
	}, nil
}

func compareNumbers(a, b interface{}) int {
	aFloat, aOK := toFloat64(a)
	bFloat, bOK := toFloat64(b)
	if aOK && bOK {
		if aFloat > bFloat {
			return 1
		} else if aFloat < bFloat {
			return -1
		}
		return 0
	}
	return 0
}

func toFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func containsString(actual, expected interface{}) bool {
	actualStr := fmt.Sprintf("%v", actual)
	expectedStr := fmt.Sprintf("%v", expected)
	return contains(actualStr, expectedStr)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
