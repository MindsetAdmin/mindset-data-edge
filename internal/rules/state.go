package rules

import (
	"sync"
	"time"
)

// TagState représente l'état d'un tag à un moment donné
type TagState struct {
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
}

// StateStore est un store thread-safe pour les états des tags
type StateStore struct {
	mu      sync.RWMutex
	states  map[string]*TagState
	history map[string][]*TagState
}

// NewStateStore crée un nouveau store
func NewStateStore() *StateStore {
	return &StateStore{
		states:  make(map[string]*TagState),
		history: make(map[string][]*TagState),
	}
}

// Set met à jour un état et retourne l'ancien
func (s *StateStore) Set(topic string, value interface{}, timestamp time.Time) *TagState {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldState := s.states[topic]
	s.states[topic] = &TagState{Value: value, Timestamp: timestamp}
	return oldState
}

// Get récupère un état
func (s *StateStore) Get(topic string) (*TagState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.states[topic]
	return state, ok
}

// GetAll retourne tous les états
func (s *StateStore) GetAll() map[string]*TagState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*TagState)
	for k, v := range s.states {
		result[k] = v
	}
	return result
}

// GetHistory retourne l'historique pour la corrélation
func (s *StateStore) GetHistory(topic string, since time.Time) []*TagState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, ok := s.history[topic]
	if !ok {
		return nil
	}

	var result []*TagState
	for _, h := range history {
		if h.Timestamp.After(since) {
			result = append(result, h)
		}
	}
	return result
}

// AddHistory ajoute un état à l'historique
func (s *StateStore) AddHistory(topic string, state *TagState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.history[topic]; !ok {
		s.history[topic] = make([]*TagState, 0)
	}
	s.history[topic] = append(s.history[topic], state)

	// Garder les 100 derniers pour éviter la mémoire infinie
	if len(s.history[topic]) > 100 {
		s.history[topic] = s.history[topic][len(s.history[topic])-100:]
	}
}
