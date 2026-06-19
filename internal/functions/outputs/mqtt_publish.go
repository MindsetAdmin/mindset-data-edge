package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTPublishConfig configuration
type MQTTPublishConfig struct {
	Topic    string `json:"topic"`
	QoS      byte   `json:"qos"`
	Retained bool   `json:"retained"`
}

// MQTTPublishHandler handler
type MQTTPublishHandler struct {
	client mqtt.Client
}

// NewMQTTPublishHandler crée un nouveau handler
func NewMQTTPublishHandler(client mqtt.Client) *MQTTPublishHandler {
	return &MQTTPublishHandler{client: client}
}

// GetFunction retourne la définition
func (h *MQTTPublishHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "mqtt_publish",
		Type:        functions.TypeOutput,
		Description: "Publie un message sur un topic MQTT",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			payload := params["payload"]
			if payload == nil {
				payload = params
			}
			return nil, h.Execute(payload, params)
		},
	}
}

// Execute publie le message
func (h *MQTTPublishHandler) Execute(payload interface{}, config map[string]interface{}) error {
	topic, ok := config["topic"].(string)
	if !ok {
		return fmt.Errorf("missing topic in config")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	qos := byte(1)
	if qosVal, ok := config["qos"].(float64); ok {
		qos = byte(qosVal)
	}

	retained := false
	if retainedVal, ok := config["retained"].(bool); ok {
		retained = retainedVal
	}

	token := h.client.Publish(topic, qos, retained, data)
	token.Wait()
	return token.Error()
}
