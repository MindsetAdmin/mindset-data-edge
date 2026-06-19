package connectors

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTSubscribeConfig configuration
type MQTTSubscribeConfig struct {
	Topic string `json:"topic"`
	QoS   byte   `json:"qos"`
}

// MQTTSubscribeResult résultat
type MQTTSubscribeResult struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// MQTTSubscribeHandler handler
type MQTTSubscribeHandler struct {
	client    mqtt.Client
	messageCh chan MQTTSubscribeResult
}

// NewMQTTSubscribeHandler crée un nouveau handler
func NewMQTTSubscribeHandler(client mqtt.Client) *MQTTSubscribeHandler {
	return &MQTTSubscribeHandler{
		client:    client,
		messageCh: make(chan MQTTSubscribeResult, 100),
	}
}

// GetFunction retourne la définition
func (h *MQTTSubscribeHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "mqtt_subscribe",
		Type:        functions.TypeConnector,
		Description: "Souscrit à un topic MQTT et reçoit les messages",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			return h.Execute(context.Background(), params)
		},
	}
}

// Execute démarre la souscription
func (h *MQTTSubscribeHandler) Execute(ctx context.Context, config map[string]interface{}) (<-chan MQTTSubscribeResult, error) {
	topic, ok := config["topic"].(string)
	if !ok {
		return nil, fmt.Errorf("missing topic in config")
	}

	qos := byte(1)
	if qosVal, ok := config["qos"].(float64); ok {
		qos = byte(qosVal)
	}

	token := h.client.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		h.messageCh <- MQTTSubscribeResult{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}
	})

	token.Wait()
	if token.Error() != nil {
		return nil, token.Error()
	}

	return h.messageCh, nil
}
