package outputs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DashboardHandler is an output that pins the incoming data/event onto the
// dashboard. It publishes to mindset/dashboard/<label> (retained) — the API
// server's LiveHub forwards those to the dashboard over WebSocket.
type DashboardHandler struct {
	client mqtt.Client
}

func NewDashboardHandler(client mqtt.Client) *DashboardHandler {
	return &DashboardHandler{client: client}
}

// GetFunction returns the function definition for the registry.
func (h *DashboardHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "add_to_dashboard",
		Type:        functions.TypeOutput,
		Description: "Affiche la donnée ou l'événement dans le dashboard",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			return nil, h.Execute(params)
		},
	}
}

// Execute publishes a dashboard widget message.
//
//	config:
//	  label  string  widget name (e.g. "Température machine1")
//	  kind   string  "value" | "event"  (how the dashboard renders it)
func (h *DashboardHandler) Execute(params map[string]interface{}) error {
	if h.client == nil {
		return fmt.Errorf("add_to_dashboard: no MQTT client available")
	}

	label, _ := params["label"].(string)
	if label == "" {
		label = "widget"
	}
	kind, _ := params["kind"].(string)
	if kind == "" {
		kind = "value"
	}

	data := params["payload"]
	if data == nil {
		data = params
	}

	msg := map[string]interface{}{
		"label":        label,
		"kind":         kind,
		"data":         data,
		"timestamp_ms": time.Now().UnixMilli(),
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	topic := "mindset/dashboard/" + sanitizeSegment(label)
	token := h.client.Publish(topic, 1, true, payload) // retained: last value persists
	token.Wait()
	return token.Error()
}

// sanitizeSegment makes a safe MQTT topic segment from a label.
func sanitizeSegment(s string) string {
	var b strings.Builder
	for _, c := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9', c == '-', c == '_':
			b.WriteRune(c)
		case c == ' ':
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "widget"
	}
	return b.String()
}
