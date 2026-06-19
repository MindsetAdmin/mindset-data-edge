// internal/kg/subscriber.go
package kg

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// KGSubscriber souscrit aux événements et enrichit le KG
type KGSubscriber struct {
	client paho.Client
	kg     *KnowledgeGraph
}

// NewKGSubscriber crée un nouveau subscriber KG
func NewKGSubscriber(brokerURL string, kg *KnowledgeGraph) (*KGSubscriber, error) {
	opts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("mindset-kg-subscriber").
		SetAutoReconnect(true)

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	log.Printf("[KG] Subscriber connected to broker: %s", brokerURL)
	return &KGSubscriber{
		client: client,
		kg:     kg,
	}, nil
}

// Start démarre la souscription
func (s *KGSubscriber) Start() error {
	token := s.client.Subscribe("mindset/events/micro-stop", 1, s.onMicroStop)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	log.Printf("[KG] Subscribed to mindset/events/micro-stop")
	return nil
}

// onMicroStop traite les micro-stops et enrichit le KG
func (s *KGSubscriber) onMicroStop(_ paho.Client, msg paho.Message) {
	var event struct {
		Timestamp  time.Time `json:"timestamp"`
		Duration   float64   `json:"duration_seconds"`
		WorkCenter string    `json:"work_center"`
		UNSTopic   string    `json:"uns_topic,omitempty"`
		Cause      string    `json:"cause,omitempty"`
		Confidence float64   `json:"confidence,omitempty"`
		CostEur    float64   `json:"cost_eur,omitempty"`
	}

	if err := json.Unmarshal(msg.Payload(), &event); err != nil {
		log.Printf("[KG] Failed to unmarshal micro-stop: %v", err)
		return
	}

	log.Printf("[KG] Processing micro-stop: %s, duration: %.0fs, cost: %.2f€",
		event.WorkCenter, event.Duration, event.CostEur)

	// 1. Nœud Équipement (machine) - existe déjà si créé par discovery
	equipmentID := fmt.Sprintf("equipment_%s", event.WorkCenter)
	err := s.kg.AddNode(equipmentID, "Equipment", event.WorkCenter, map[string]interface{}{
		"work_center": event.WorkCenter,
		"uns_topic":   event.UNSTopic,
		"type":        "machine",
	})
	if err != nil {
		log.Printf("[KG] Failed to add equipment node: %v", err)
	}

	// 2. Nœud Événement (micro-stop)
	eventID := fmt.Sprintf("event_%d_%s", event.Timestamp.UnixNano(), event.WorkCenter)
	err = s.kg.AddNode(eventID, "Event", "Micro-stop", map[string]interface{}{
		"duration_seconds": event.Duration,
		"timestamp":        event.Timestamp.Format(time.RFC3339),
		"type":             "microstop",
		"work_center":      event.WorkCenter,
	})
	if err != nil {
		log.Printf("[KG] Failed to add event node: %v", err)
	}

	// 3. Relation: Événement → Équipement (occurred_at)
	edgeID := fmt.Sprintf("edge_%s_%s", eventID, equipmentID)
	err = s.kg.AddEdge(edgeID, eventID, equipmentID, "occurred_at", 1.0)
	if err != nil {
		log.Printf("[KG] Failed to add occurred_at edge: %v", err)
	}

	// 4. Si une cause est présente, ajouter le nœud Cause
	if event.Cause != "" {
		causeID := fmt.Sprintf("cause_%s", sanitizeID(event.Cause))
		err = s.kg.AddNode(causeID, "Cause", event.Cause, map[string]interface{}{
			"confidence": event.Confidence,
		})
		if err != nil {
			log.Printf("[KG] Failed to add cause node: %v", err)
		}

		// Relation: Événement → Cause (caused_by)
		causeEdgeID := fmt.Sprintf("edge_%s_%s", eventID, causeID)
		err = s.kg.AddEdge(causeEdgeID, eventID, causeID, "caused_by", event.Confidence)
		if err != nil {
			log.Printf("[KG] Failed to add caused_by edge: %v", err)
		}
	}

	// 5. Si un coût est présent, ajouter le nœud Coût
	if event.CostEur > 0 {
		costID := fmt.Sprintf("cost_%s", eventID)
		err = s.kg.AddNode(costID, "Cost", fmt.Sprintf("%.2f€", event.CostEur), map[string]interface{}{
			"amount_eur": event.CostEur,
		})
		if err != nil {
			log.Printf("[KG] Failed to add cost node: %v", err)
		}

		// Relation: Événement → Coût (costs)
		costEdgeID := fmt.Sprintf("edge_cost_%s", eventID)
		err = s.kg.AddEdge(costEdgeID, eventID, costID, "costs", 1.0)
		if err != nil {
			log.Printf("[KG] Failed to add costs edge: %v", err)
		}
	}
}

// Stop arrête le subscriber
func (s *KGSubscriber) Stop() {
	if s.client != nil {
		s.client.Disconnect(250)
		log.Printf("[KG] Subscriber stopped")
	}
}
