package calculates

import (
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// DurationResult résultat
type DurationResult struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Seconds   float64   `json:"duration_seconds"`
	Minutes   float64   `json:"duration_minutes"`
}

// DurationHandler handler
type DurationHandler struct {
	startTimes map[string]time.Time
}

// NewDurationHandler crée un nouveau handler
func NewDurationHandler() *DurationHandler {
	return &DurationHandler{
		startTimes: make(map[string]time.Time),
	}
}

// GetFunction retourne la définition
func (h *DurationHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "calculate_duration",
		Type:        functions.TypeCalculate,
		Description: "Calcule la durée entre un début et une fin",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			eventID, _ := params["event_id"].(string)
			endTime, _ := params["end_time"].(time.Time)
			if endTime.IsZero() {
				endTime = time.Now()
			}
			return h.Execute(eventID, endTime)
		},
	}
}

// Execute calcule la durée
func (h *DurationHandler) Execute(eventID string, endTime time.Time) (*DurationResult, error) {
	startTime, exists := h.startTimes[eventID]
	if !exists {
		// C'est le début, enregistrer
		h.startTimes[eventID] = endTime
		return nil, nil
	}

	// Calculer la durée
	seconds := endTime.Sub(startTime).Seconds()

	// Nettoyer
	delete(h.startTimes, eventID)

	return &DurationResult{
		StartTime: startTime,
		EndTime:   endTime,
		Seconds:   seconds,
		Minutes:   seconds / 60,
	}, nil
}

// SetStartTime enregistre un temps de début (pour usage externe)
func (h *DurationHandler) SetStartTime(eventID string, startTime time.Time) {
	h.startTimes[eventID] = startTime
}
