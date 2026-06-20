package main

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// rateWindow is the sliding window used to compute msg/s per topic.
const rateWindow = 5 * time.Second

// --- Topic stats -----------------------------------------------------------

type topicStat struct {
	count  int64
	recent []int64 // unix-milli timestamps, capped
}

type TopicRegistry struct {
	mu     sync.RWMutex
	topics map[string]*topicStat
}

func NewTopicRegistry() *TopicRegistry {
	return &TopicRegistry{topics: make(map[string]*topicStat)}
}

func (r *TopicRegistry) record(topic string, ms int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.topics[topic]
	if s == nil {
		s = &topicStat{}
		r.topics[topic] = s
	}
	s.count++
	s.recent = append(s.recent, ms)
	if len(s.recent) > 200 {
		s.recent = s.recent[len(s.recent)-200:]
	}
}

// TopicView is the API shape for a topic.
type TopicView struct {
	Topic    string  `json:"topic"`
	Category string  `json:"category"` // raw | site | events | other
	Count    int64   `json:"count"`
	Rate     float64 `json:"rate_per_sec"`
}

func categoryOf(topic string) string {
	switch {
	case strings.HasPrefix(topic, "mindset/raw/"):
		return "raw"
	case strings.HasPrefix(topic, "mindset/site/"):
		return "site"
	case strings.HasPrefix(topic, "mindset/events"):
		return "events"
	default:
		return "other"
	}
}

func (r *TopicRegistry) list() []TopicView {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cutoff := time.Now().Add(-rateWindow).UnixMilli()
	out := make([]TopicView, 0, len(r.topics))
	for topic, s := range r.topics {
		n := 0
		for _, ts := range s.recent {
			if ts >= cutoff {
				n++
			}
		}
		out = append(out, TopicView{
			Topic:    topic,
			Category: categoryOf(topic),
			Count:    s.count,
			Rate:     float64(n) / rateWindow.Seconds(),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Topic < out[j].Topic })
	return out
}

// --- Machine state tracking ------------------------------------------------

type Transition struct {
	From        bool      `json:"from"`
	To          bool      `json:"to"`
	At          time.Time `json:"at"`
	DurationSec float64   `json:"duration_seconds"`
}

type MachineState struct {
	Running bool         `json:"running"`
	Since   time.Time    `json:"since"`
	History []Transition `json:"history"`
}

type StateTracker struct {
	mu     sync.RWMutex
	states map[string]*MachineState
}

func NewStateTracker() *StateTracker {
	return &StateTracker{states: make(map[string]*MachineState)}
}

// observe records a status value for a work center and detects transitions.
func (t *StateTracker) observe(workCenter string, running bool, ts time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	st := t.states[workCenter]
	if st == nil {
		t.states[workCenter] = &MachineState{Running: running, Since: ts}
		return
	}
	if st.Running == running {
		return
	}
	dur := ts.Sub(st.Since).Seconds()
	st.History = append(st.History, Transition{From: st.Running, To: running, At: ts, DurationSec: dur})
	if len(st.History) > 50 {
		st.History = st.History[len(st.History)-50:]
	}
	st.Running = running
	st.Since = ts
}

func (t *StateTracker) get(workCenter string) *MachineState {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.states[workCenter]
}

// --- Live hub: one subscription feeding tags, topics and state -------------

type LiveHub struct {
	tags   *TagRegistry
	topics *TopicRegistry
	states *StateTracker
}

func NewLiveHub(tags *TagRegistry, topics *TopicRegistry, states *StateTracker) *LiveHub {
	return &LiveHub{tags: tags, topics: topics, states: states}
}

// Start subscribes to mindset/# and dispatches every message.
func (h *LiveHub) Start(client mqtt.Client) error {
	tok := client.Subscribe("mindset/#", 0, func(_ mqtt.Client, m mqtt.Message) {
		now := time.Now()
		h.topics.record(m.Topic(), now.UnixMilli())

		if strings.HasPrefix(m.Topic(), "mindset/raw/") {
			var t Tag
			if err := json.Unmarshal(m.Payload(), &t); err != nil || t.NodeID == "" {
				return
			}
			h.tags.upsert(t)

			// Status tags (e.g. machine1.status, boolean) drive state tracking.
			if wc, ok := workCenterOf(t.Name); ok && isStatusTag(t.Name) {
				if b, ok := toBool(t.Value); ok {
					h.states.observe(wc, b, now)
				}
			}
		}
	})
	tok.Wait()
	return tok.Error()
}

// workCenterOf extracts the machine name from a tag name like "machine1.status".
func workCenterOf(name string) (string, bool) {
	if i := strings.Index(name, "."); i > 0 {
		return name[:i], true
	}
	return "", false
}

func isStatusTag(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".status")
}

func toBool(v interface{}) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case float64:
		return b != 0, true
	default:
		return false, false
	}
}
