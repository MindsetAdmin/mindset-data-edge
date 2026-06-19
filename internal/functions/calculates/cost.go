package calculates

import (
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// CostConfig configuration du calcul de coût
type CostConfig struct {
	HourlyRate float64 `json:"hourly_rate"` // Coût horaire en €/h
	Currency   string  `json:"currency"`    // Devise (EUR)
}

// CostResult résultat du calcul
type CostResult struct {
	DurationSeconds float64 `json:"duration_seconds"`
	DurationMinutes float64 `json:"duration_minutes"`
	CostPerMinute   float64 `json:"cost_per_minute"`
	TotalCost       float64 `json:"total_cost_eur"`
	Currency        string  `json:"currency"`
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

// Execute calcule le coût
func (h *CostHandler) Execute(durationSeconds float64, config map[string]interface{}) (*CostResult, error) {
	hourlyRate := h.hourlyRate

	// Surcharge par config si présente
	if rate, ok := config["hourly_rate"].(float64); ok {
		hourlyRate = rate
	}

	currency := "EUR"
	if curr, ok := config["currency"].(string); ok {
		currency = curr
	}

	durationMinutes := durationSeconds / 60
	costPerMinute := hourlyRate / 60
	totalCost := durationMinutes * costPerMinute

	return &CostResult{
		DurationSeconds: durationSeconds,
		DurationMinutes: durationMinutes,
		CostPerMinute:   costPerMinute,
		TotalCost:       totalCost,
		Currency:        currency,
	}, nil
}

// SetHourlyRate met à jour le taux horaire
func (h *CostHandler) SetHourlyRate(rate float64) {
	h.hourlyRate = rate
}
