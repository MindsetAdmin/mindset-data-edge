// internal/uns/contextualizer.go
package uns

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// RawTagMessage is what we receive from the raw MQTT topic.
// It matches the struct published by mqtt/publisher.go PublishRaw().
type RawTagMessage struct {
	NodeID    string      `json:"node_id"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	DataType  string      `json:"data_type"`
	Timestamp int64   `json:"timestamp_ms"`
}

// ContextualizedMessage is the enriched, ISA-95 structured payload.
// This is what every downstream consumer (Rules Engine, KG, Dashboard) receives.
type ContextualizedMessage struct {
	TimestampMs int64       `json:"timestamp_ms"` // Unix milliseconds — precise
	Value       interface{} `json:"value"`
	Unit        string      `json:"unit"`     // inferred: "celsius", "bar", ""
	DataType    string      `json:"data_type"` // "Float", "Boolean", "Int32"
	Metadata    Metadata    `json:"metadata"`
}

// Metadata carries the full context of a signal's origin and position.
type Metadata struct {
	SourceProtocol string `json:"source_protocol"` // "OPC-UA"
	OriginalNodeID string `json:"original_node_id"` // "ns=3;i=1011"
	OriginalName   string `json:"original_name"`    // "machine1.temp"
	SiteID         string `json:"site_id"`          // "local-test"
	Area           string `json:"area"`             // "area1"
	WorkCenter     string `json:"work_center"`      // "machine1"
	WorkUnit       string `json:"work_unit"`        // "ligne1" or ""
	TagName        string `json:"tag_name"`         // "temperature"
	UNSTopic       string `json:"uns_topic"`        // full ISA-95 topic
}

// Contextualizer subscribes to raw MQTT topics, enriches the data,
// and republishes on standardized ISA-95 topics.
//
// It uses TWO separate MQTT clients to avoid receiving its own messages:
//   - subscriber: listens on mindset/raw/#
//   - publisher:  publishes on mindset/site/#
type Contextualizer struct {
	subscriber paho.Client
	publisher  paho.Client
	mapper     *Mapper
	siteID     string
}

// NewContextualizer creates a contextualizer connected to the given broker.
// It initializes both subscriber and publisher clients but does not start yet.
func NewContextualizer(brokerURL, siteID string, mapper *Mapper) (*Contextualizer, error) {
	c := &Contextualizer{
		mapper: mapper,
		siteID: siteID,
	}

	// ── Subscriber client ─────────────────────────────────────────────────
	subOpts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("mindset-uns-contextualizer-sub"). // unique client ID
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second).
		SetOnConnectHandler(func(_ paho.Client) {
			log.Printf("[UNS] Subscriber connected to broker %s", brokerURL)
		}).
		SetConnectionLostHandler(func(_ paho.Client, err error) {
			log.Printf("[UNS] Subscriber connection lost: %v", err)
		})

	c.subscriber = paho.NewClient(subOpts)
	if token := c.subscriber.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("subscriber connect failed: %w", token.Error())
	}

	// ── Publisher client ──────────────────────────────────────────────────
	pubOpts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("mindset-uns-contextualizer-pub"). // different client ID
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second).
		SetOnConnectHandler(func(_ paho.Client) {
			log.Printf("[UNS] Publisher connected to broker %s", brokerURL)
		})

	c.publisher = paho.NewClient(pubOpts)
	if token := c.publisher.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("publisher connect failed: %w", token.Error())
	}

	return c, nil
}

// Start subscribes to mindset/raw/# and begins processing messages.
// This is non-blocking — it returns immediately after subscribing.
func (c *Contextualizer) Start() error {
	// QoS 0 — at most once delivery
	// Sufficient for high-frequency sensor data where occasional loss is acceptable
	// Use QoS 1 for events that must not be lost (alarms, stops)
	token := c.subscriber.Subscribe("mindset/raw/#", 0, c.onRawMessage)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("subscribe to mindset/raw/# failed: %w", token.Error())
	}

	log.Printf("[UNS] Contextualizer started — subscribed to mindset/raw/#")
	log.Printf("[UNS] Will publish enriched data to mindset/site/#")
	return nil
}

// onRawMessage is the callback for every raw MQTT message received.
// It is called in a goroutine by the paho library — it must be thread-safe.
func (c *Contextualizer) onRawMessage(_ paho.Client, msg paho.Message) {
	// ── Step 1: Deserialize the raw message ───────────────────────────────
	log.Printf("[UNS] Received raw message on %s", msg.Topic())
	var raw RawTagMessage
	if err := json.Unmarshal(msg.Payload(), &raw); err != nil {
		log.Printf("[UNS] Warning: failed to unmarshal raw message on %s: %v",
			msg.Topic(), err)
		return
	}

	// Guard: skip if name is empty (shouldn't happen, but defensive)
	if raw.Name == "" {
		log.Printf("[UNS] Warning: received raw message with empty name on %s", msg.Topic())
		return
	}

	// ── Step 2: Map to ISA-95 UNS structure ───────────────────────────────
	node := c.mapper.MapTag(raw.Name, raw.DataType)

	// ── Step 3: Build enriched contextualized message ─────────────────────
	contextualized := ContextualizedMessage{
		TimestampMs: raw.Timestamp, // precise Unix milliseconds
		Value:       raw.Value,
		Unit:        node.Unit,
		DataType:    raw.DataType,
		Metadata: Metadata{
			SourceProtocol: "OPC-UA",
			OriginalNodeID: raw.NodeID,
			OriginalName:   raw.Name,
			SiteID:         node.Site,
			Area:           node.Area,
			WorkCenter:     node.WorkCenter,
			WorkUnit:       node.WorkUnit,
			TagName:        node.TagName,
			UNSTopic:       node.FullTopic(),
		},
	}

	// ── Step 4: Serialize and publish on ISA-95 topic ─────────────────────
	payload, err := json.Marshal(contextualized)
	if err != nil {
		log.Printf("[UNS] Warning: failed to marshal contextualized message: %v", err)
		return
	}

	// Publish with QoS 0 (non-retained)
	// The topic is the full ISA-95 path:
	// mindset/site/local-test/area1/machine1/temperature
	token := c.publisher.Publish(node.FullTopic(), 0, false, payload)
	token.Wait()

	if token.Error() != nil {
		log.Printf("[UNS] Warning: failed to publish to %s: %v",
			node.FullTopic(), token.Error())
		return
	}

	// Debug log — comment out in production to reduce noise
	log.Printf("[UNS] %s → %s = %v (%s) [%s]",
		raw.Name,
		node.FullTopic(),
		raw.Value,
		node.DataType,
		node.Unit,
	)
}

// Stop gracefully disconnects both MQTT clients.
func (c *Contextualizer) Stop() {
	log.Printf("[UNS] Contextualizer stopping...")
	c.subscriber.Unsubscribe("mindset/raw/#")
	c.subscriber.Disconnect(500) // wait 500ms for in-flight messages
	c.publisher.Disconnect(500)
	log.Printf("[UNS] Contextualizer stopped")
}