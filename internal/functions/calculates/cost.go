package calculates

import (
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// CostConfig configuration du calcul de coût
type CostConfig struct {
	HourlyRate float64 `json:"hourly_rate"` // Coût horaire en €/h
	Currency   string  `json:"currency"`    // Devise (EUR)
}

// CostResult résultat du calcul.
//
// WorkCenter/UNSTopic/Cause/Timestamp are passed through from the trigger
// event (they arrive already-merged into the handler's params — see
// internal/pipeline/engine.go's executeNode) rather than computed here.
// They exist so this node's output, when it's a pipeline's terminal/output
// node, matches the shape internal/kg/subscriber.go's onMicroStop expects on
// mindset/events/micro-stop (work_center + timestamp are required for the KG
// to accept the event; cost_eur is what makes it show up in kg_cost_summary).
type CostResult struct {
	DurationSeconds float64   `json:"duration_seconds"`
	DurationMinutes float64   `json:"duration_minutes"`
	CostPerMinute   float64   `json:"cost_per_minute"`
	TotalCost       float64   `json:"total_cost_eur"`
	CostEur         float64   `json:"cost_eur"`
	Currency        string    `json:"currency"`
	WorkCenter      string    `json:"work_center,omitempty"`
	UNSTopic        string    `json:"uns_topic,omitempty"`
	Cause           string    `json:"cause,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

// CostHandler handler
type CostHandler struct {
	hourlyRate float64
}

// NewCostHandler crée un nouveau handler
func NewCostHandler(hourlyRate float64) *CostHandler {
	return &CostHandler{
		hourlyRate: hourlyRate,
	}
}

// GetFunction retourne la définition
func (h *CostHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "calculate_cost",
		Type:        functions.TypeCalculate,
		Description: "Calcule le coût d'un micro-stop en euros",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			durationSeconds, _ := params["duration_seconds"].(float64)
			return h.Execute(durationSeconds, params)
		},
	}
}

// asFloat64 handles both float64 and int — a pipeline YAML's node.Config is
// decoded by yaml.v3 into map[string]interface{}, and yaml.v3 unmarshals a
// plain integer literal (e.g. "hourly_rate: 400", no decimal point) as Go
// `int`, not `float64`. A bare `.(float64)` type assertion silently fails on
// that and falls through to the default rate with no error anywhere — found
// live when a seeded hourly_rate: 400 produced ~85€/h totals instead
// (Entry 135).
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

// Execute calcule le coût
func (h *CostHandler) Execute(durationSeconds float64, config map[string]interface{}) (*CostResult, error) {
	hourlyRate := h.hourlyRate

	// Surcharge par config si présente
	if rate, ok := asFloat64(config["hourly_rate"]); ok {
		hourlyRate = rate
	}

	// Tarif par produit (table chargée depuis CSV/Excel). Si l'événement porte un
	// "product" présent dans la table, on utilise son coût horaire.
	if product, ok := config["product"].(string); ok && product != "" {
		if rates, ok := config["rates"].(map[string]interface{}); ok {
			if row, ok := rates[product].(map[string]interface{}); ok {
				if hr, ok := asFloat64(row["hourly_rate"]); ok && hr > 0 {
					hourlyRate = hr
				}
			}
		}
	}

	currency := "EUR"
	if curr, ok := config["currency"].(string); ok {
		currency = curr
	}

	durationMinutes := durationSeconds / 60
	costPerMinute := hourlyRate / 60
	totalCost := durationMinutes * costPerMinute

	workCenter, _ := config["work_center"].(string)
	unsTopic, _ := config["uns_topic"].(string)
	cause, _ := config["cause"].(string)

	// The KG subscriber requires a non-zero timestamp. status-change events
	// (internal/rules/engine.go) carry start_time as RFC3339; fall back to
	// now() when it's absent or unparseable (e.g. a manual/seeded run).
	ts := time.Now()
	if startTimeStr, ok := config["start_time"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			ts = parsed
		}
	}

	return &CostResult{
		DurationSeconds: durationSeconds,
		DurationMinutes: durationMinutes,
		CostPerMinute:   costPerMinute,
		TotalCost:       totalCost,
		CostEur:         totalCost,
		Currency:        currency,
		WorkCenter:      workCenter,
		UNSTopic:        unsTopic,
		Cause:           cause,
		Timestamp:       ts,
	}, nil
}

// SetHourlyRate met à jour le taux horaire
func (h *CostHandler) SetHourlyRate(rate float64) {
	h.hourlyRate = rate
}
