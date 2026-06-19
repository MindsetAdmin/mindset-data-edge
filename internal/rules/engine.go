package rules

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// EnrichedMessage est la structure du message UNS reçu
type EnrichedMessage struct {
	TimestampMs int64       `json:"timestamp_ms"`
	Value       interface{} `json:"value"`
	Unit        string      `json:"unit"`
	DataType    string      `json:"data_type"`
	Metadata    struct {
		SourceProtocol string `json:"source_protocol"`
		OriginalNodeID string `json:"original_node_id"`
		OriginalName   string `json:"original_name"`
		SiteID         string `json:"site_id"`
		Area           string `json:"area"`
		WorkCenter     string `json:"work_center"`
		WorkUnit       string `json:"work_unit"`
		TagName        string `json:"tag_name"`
		UNSTopic       string `json:"uns_topic"`
	} `json:"metadata"`
}

// Engine est le moteur de règles - il gère l'état et déclenche les pipelines
type Engine struct {
	client     paho.Client
	publisher  paho.Client
	stateStore *StateStore
	brokerURL  string
	siteID     string

	// stopStarts stocke l'heure de début des arrêts pour le calcul de durée
	stopStarts map[string]time.Time
	stopMu     sync.Mutex
}

// NewEngine crée un nouveau moteur
func NewEngine(brokerURL, siteID string) (*Engine, error) {
	e := &Engine{
		stateStore: NewStateStore(),
		brokerURL:  brokerURL,
		siteID:     siteID,
		stopStarts: make(map[string]time.Time),
	}

	// Client pour souscrire aux topics UNS
	subOpts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("mindset-rules-engine-sub").
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second)

	e.client = paho.NewClient(subOpts)
	if token := e.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	log.Printf("[RULES] Subscriber connected to broker: %s", brokerURL)

	// Client pour publier les événements
	pubOpts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("mindset-rules-engine-pub").
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second)

	e.publisher = paho.NewClient(pubOpts)
	if token := e.publisher.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	log.Printf("[RULES] Publisher connected to broker: %s", brokerURL)

	return e, nil
}

// Start démarre le moteur et s'abonne aux topics UNS
func (e *Engine) Start() error {
	token := e.client.Subscribe("mindset/site/#", 1, e.onMessage)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	log.Printf("[RULES] Engine started — subscribed to mindset/site/#")
	return nil
}

// GetStateStore retourne le store pour les autres composants (energy, causality, etc.)
func (e *Engine) GetStateStore() *StateStore {
	return e.stateStore
}

// GetPublisher retourne le client MQTT pour la publication
func (e *Engine) GetPublisher() paho.Client {
	return e.publisher
}

// onMessage est le callback pour chaque message UNS reçu
func (e *Engine) onMessage(_ paho.Client, msg paho.Message) {
	var enriched EnrichedMessage

	if err := json.Unmarshal(msg.Payload(), &enriched); err != nil {
		log.Printf("[RULES] Failed to unmarshal: %v", err)
		return
	}

	timestamp := time.UnixMilli(enriched.TimestampMs)
	topic := msg.Topic()

	// Mettre à jour le state store
	oldState := e.stateStore.Set(topic, enriched.Value, timestamp)

	// Ajouter à l'historique pour la corrélation (causality)
	e.stateStore.AddHistory(topic, &TagState{
		Value:     enriched.Value,
		Timestamp: timestamp,
	})

	// Détecter les transitions d'état pour les tags de type "status"
	if enriched.Metadata.TagName == "status" {
		e.handleStatusChange(
			topic,
			enriched.Value,
			enriched.Metadata.WorkCenter,
			enriched.Metadata.UNSTopic,
			oldState,
			timestamp,
		)
	}
}

// handleStatusChange gère les transitions d'état (Run ↔ Stop)
func (e *Engine) handleStatusChange(topic string, value interface{}, workCenter, unsTopic string, oldState *TagState, currentTime time.Time) {
	// Convertir la valeur en booléen (état Run/Stop)
	currentRunning, ok := value.(bool)
	if !ok {
		return
	}

	// Premier message reçu pour ce tag — juste enregistrer l'état
	if oldState == nil {
		log.Printf("[RULES] Initial status for %s: running=%v", workCenter, currentRunning)
		return
	}

	// Convertir l'ancienne valeur en booléen
	oldRunning, ok := oldState.Value.(bool)
	if !ok {
		return
	}

	// Pas de changement d'état
	if oldRunning == currentRunning {
		return
	}

	// Transition Run → Stop : enregistrer le début de l'arrêt
	if oldRunning && !currentRunning {
		log.Printf("[RULES] 🛑 %s stopped at %v", workCenter, currentTime)
		e.stopMu.Lock()
		e.stopStarts[topic] = currentTime
		e.stopMu.Unlock()
		return
	}

	// Transition Stop → Run : calculer la durée et publier l'événement
	if !oldRunning && currentRunning {
		e.stopMu.Lock()
		startTime, exists := e.stopStarts[topic]
		delete(e.stopStarts, topic)
		e.stopMu.Unlock()

		if !exists {
			return
		}

		duration := currentTime.Sub(startTime).Seconds()
		log.Printf("[RULES] ▶️ %s started after %.0f seconds", workCenter, duration)

		// Publier un événement de changement de statut
		// → Le pipeline YAML (microstop_detection) s'abonne à ce topic
		// → Il vérifie la durée et publie un micro-stop si 30s < durée < 3min
		e.publishStatusEvent(topic, workCenter, unsTopic, startTime, duration, oldRunning, currentRunning)
	}
}

// publishStatusEvent publie un événement de changement de statut sur MQTT
// Ce sera consommé par le pipeline YAML microstop_detection
func (e *Engine) publishStatusEvent(topic, workCenter, unsTopic string, startTime time.Time, duration float64, oldState, newState bool) {
	event := map[string]interface{}{
		"topic":            topic,
		"work_center":      workCenter,
		"uns_topic":        unsTopic,
		"start_time":       startTime.Format(time.RFC3339),
		"duration_seconds": duration,
		"previous_state":   oldState,
		"current_state":    newState,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("[RULES] Failed to marshal status event: %v", err)
		return
	}

	token := e.publisher.Publish("mindset/events/status-change", 1, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Printf("[RULES] Failed to publish status event: %v", token.Error())
	} else {
		log.Printf("[RULES] 📤 Published status change for %s (duration: %.0fs)", workCenter, duration)
	}
}

// Stop arrête proprement le moteur
func (e *Engine) Stop() {
	log.Printf("[RULES] Stopping engine...")
	if e.client != nil {
		e.client.Unsubscribe("mindset/site/#")
		e.client.Disconnect(500)
	}
	if e.publisher != nil {
		e.publisher.Disconnect(500)
	}
	log.Printf("[RULES] Engine stopped")
}
