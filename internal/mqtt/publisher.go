package mqtt

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TagMessage structure du message brut OPC-UA
type TagMessage struct {
    NodeID    string      `json:"node_id"`
    Name      string      `json:"name"`
    Value     interface{} `json:"value"`
    DataType  string      `json:"data_type"`
    Timestamp int64       `json:"timestamp_ms"`
}

// Publisher client MQTT pour publier les données
type Publisher struct {
    client  mqtt.Client
    siteID  string
    broker  string
}

// NewPublisher crée un nouveau publisher
func NewPublisher(brokerURL, siteID string) (*Publisher, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(brokerURL)
    opts.SetClientID(fmt.Sprintf("mindset-edge-%d", time.Now().Unix()))
    opts.SetCleanSession(true)
    opts.SetAutoReconnect(true)

    client := mqtt.NewClient(opts)
    token := client.Connect()
    if token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    log.Printf("[MQTT] Connected to broker: %s", brokerURL)
    return &Publisher{
        client: client,
        siteID: siteID,
        broker: brokerURL,
    }, nil
}

// PublishRaw publie une donnée brute OPC-UA
func (p *Publisher) PublishRaw(tagName, nodeID, dataType string, value interface{}) {
    msg := TagMessage{
        NodeID:    nodeID,
        Name:      tagName,
        Value:     value,
        DataType:  dataType,
        Timestamp: time.Now().UnixMilli(),
    }

    payload, err := json.Marshal(msg)
    if err != nil {
        log.Printf("[MQTT] Failed to marshal: %v", err)
        return
    }

    topic := fmt.Sprintf("mindset/raw/%s", nodeID)
    token := p.client.Publish(topic, 1, false, payload)
    token.Wait()

    if token.Error() != nil {
        log.Printf("[MQTT] Failed to publish: %v", token.Error())
    } else {
        log.Printf("[MQTT] Published: %s = %v", tagName, value)
    }
}

// PublishJSON publie un payload arbitraire sur un topic explicite (QoS 0, non retenu).
// Utilisé par le routage dynamique côté serveur pour publier les messages ISA-95
// contextualisés sur mindset/site/#.
func (p *Publisher) PublishJSON(topic string, payload interface{}) {
    p.publishJSON(topic, payload, false)
}

// PublishJSONRetained is PublishJSON with the MQTT retained flag set, so any
// client that subscribes after this call (a restarted MCP stdio process, a
// freshly opened dashboard tab) immediately receives the last known value
// instead of waiting for the next change — a non-retained publish only
// reaches subscribers that were already listening at the moment it fired.
func (p *Publisher) PublishJSONRetained(topic string, payload interface{}) {
    p.publishJSON(topic, payload, true)
}

func (p *Publisher) publishJSON(topic string, payload interface{}, retained bool) {
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("[MQTT] Failed to marshal payload for %s: %v", topic, err)
        return
    }
    token := p.client.Publish(topic, 0, retained, data)
    token.Wait()
    if token.Error() != nil {
        log.Printf("[MQTT] Failed to publish to %s: %v", topic, token.Error())
    }
}

// PublishEvent publie un événement traité
func (p *Publisher) PublishEvent(eventType string, payload interface{}) {
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("[MQTT] Failed to marshal event: %v", err)
        return
    }

    topic := fmt.Sprintf("mindset/events/%s", eventType)
    token := p.client.Publish(topic, 1, false, data)
    token.Wait()
}

// Disconnect ferme la connexion
func (p *Publisher) Disconnect() {
    if p.client != nil {
        p.client.Disconnect(250)
        log.Printf("[MQTT] Disconnected")
    }
}