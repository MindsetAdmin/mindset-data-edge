// internal/pipeline/registry.go
package pipeline

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// Registry gère l'enregistrement des pipelines
type Registry struct {
	mu        sync.RWMutex
	pipelines map[string]*Pipeline
}

// NewRegistry crée un nouveau registre
func NewRegistry() *Registry {
	return &Registry{
		pipelines: make(map[string]*Pipeline),
	}
}

// Register enregistre un pipeline
func (r *Registry) Register(p *Pipeline) error {
	if p.ID == "" {
		return fmt.Errorf("pipeline ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.pipelines[p.ID] = p
	return nil
}

// Get retourne un pipeline par son ID
func (r *Registry) Get(id string) (*Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if p, ok := r.pipelines[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("pipeline %s not found", id)
}

// List retourne tous les pipelines
func (r *Registry) List() []*Pipeline {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*Pipeline, 0, len(r.pipelines))
	for _, p := range r.pipelines {
		list = append(list, p)
	}
	return list
}

// GetHash retourne un hash de l'état actuel des pipelines
func (r *Registry) GetHash() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1. Récupérer tous les IDs de pipelines
	ids := make([]string, 0, len(r.pipelines))
	for id := range r.pipelines {
		ids = append(ids, id)
	}
	sort.Strings(ids) // Pour garantir la stabilité

	// 2. Construire une représentation du state
	var state string
	for _, id := range ids {
		p := r.pipelines[id]
		state += id + p.Version + fmt.Sprintf("%d", len(p.Nodes))
		for _, node := range p.Nodes {
			state += node.ID + node.Function
		}
	}

	// 3. Hasher
	hash := md5.Sum([]byte(state))
	return hex.EncodeToString(hash[:])
}

// Remove supprime un pipeline
func (r *Registry) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.pipelines[id]; !ok {
		return fmt.Errorf("pipeline %s not found", id)
	}
	delete(r.pipelines, id)
	return nil
}
