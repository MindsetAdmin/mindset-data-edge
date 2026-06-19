// internal/pipeline/engine.go
package pipeline

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

// Engine est le moteur d'exécution des pipelines
type Engine struct {
	funcRegistry *functions.Registry
	pipelineReg  *Registry
	eventCh      chan PipelineEvent
	mu           sync.RWMutex
	stopCh       chan struct{}
	stopped      bool
	stopMu       sync.Mutex
}

// NewEngine crée un nouveau moteur
func NewEngine(funcRegistry *functions.Registry) *Engine {
	return &Engine{
		funcRegistry: funcRegistry,
		pipelineReg:  NewRegistry(),
		eventCh:      make(chan PipelineEvent, 100),
		stopCh:       make(chan struct{}),
	}
}

// RegisterPipeline enregistre un pipeline
func (e *Engine) RegisterPipeline(p *Pipeline) error {
	return e.pipelineReg.Register(p)
}

// GetPipeline retourne un pipeline
func (e *Engine) GetPipeline(id string) (*Pipeline, error) {
	return e.pipelineReg.Get(id)
}

// ListPipelines retourne tous les pipelines
func (e *Engine) ListPipelines() []*Pipeline {
	return e.pipelineReg.List()
}

// Events retourne le canal des événements
func (e *Engine) Events() <-chan PipelineEvent {
	return e.eventCh
}

// Execute exécute un pipeline
func (e *Engine) Execute(ctx context.Context, pipelineID string, triggerData map[string]interface{}) (*ExecutionResult, error) {
	pipeline, err := e.pipelineReg.Get(pipelineID)
	if err != nil {
		return nil, err
	}

	result := &ExecutionResult{
		PipelineID:   pipeline.ID,
		PipelineName: pipeline.Name,
		Status:       StatusRunning,
		StartTime:    time.Now(),
		Nodes:        make([]NodeExecutionResult, 0),
		Trigger:      triggerData,
	}

	e.emitEvent(PipelineEvent{
		PipelineID: pipeline.ID,
		Type:       "started",
		Timestamp:  time.Now(),
	})

	// Construire le graphe de dépendances
	nodeMap := make(map[string]Node)
	for _, node := range pipeline.Nodes {
		nodeMap[node.ID] = node
	}

	// Exécuter les nœuds dans l'ordre (tri topologique simple)
	executed := make(map[string]bool)
	results := make(map[string]interface{})

	for len(executed) < len(pipeline.Nodes) {
		progress := false

		for _, node := range pipeline.Nodes {
			if executed[node.ID] {
				continue
			}

			// Vérifier que toutes les dépendances sont satisfaites.
			// Une dépendance qui n'est pas un nœud du pipeline (ex: "trigger")
			// est considérée comme déjà satisfaite — c'est le déclencheur.
			depsMet := true
			for _, depID := range node.DependsOn {
				if _, isNode := nodeMap[depID]; isNode && !executed[depID] {
					depsMet = false
					break
				}
			}

			if !depsMet {
				continue
			}

			// Exécuter le nœud
			nodeResult := e.executeNode(ctx, node, results, triggerData)
			result.Nodes = append(result.Nodes, nodeResult)

			if nodeResult.Status == StatusSuccess {
				results[node.ID] = nodeResult.Output
				executed[node.ID] = true
				progress = true

				e.emitEvent(PipelineEvent{
					PipelineID: pipeline.ID,
					NodeID:     node.ID,
					Type:       "node_completed",
					Data:       nodeResult.Output,
					Timestamp:  time.Now(),
				})
			} else if nodeResult.Status == StatusFailed {
				result.Status = StatusFailed
				result.Error = nodeResult.Error
				result.EndTime = time.Now()
				result.TotalDuration = result.EndTime.Sub(result.StartTime)

				e.emitEvent(PipelineEvent{
					PipelineID: pipeline.ID,
					Type:       "failed",
					Data:       nodeResult.Error,
					Timestamp:  time.Now(),
				})

				return result, nil
			}
		}

		if !progress {
			break
		}
	}

	// Vérifier que tous les nœuds ont été exécutés
	if len(executed) != len(pipeline.Nodes) {
		result.Status = StatusFailed
		result.Error = "circular dependency or missing nodes"
	} else {
		result.Status = StatusSuccess
	}

	result.EndTime = time.Now()
	result.TotalDuration = result.EndTime.Sub(result.StartTime)

	e.emitEvent(PipelineEvent{
		PipelineID: pipeline.ID,
		Type:       "completed",
		Data:       result.Status,
		Timestamp:  time.Now(),
	})

	return result, nil
}

// executeNode exécute un nœud individuel
func (e *Engine) executeNode(ctx context.Context, node Node, previousResults map[string]interface{}, triggerData map[string]interface{}) NodeExecutionResult {
	startTime := time.Now()

	result := NodeExecutionResult{
		NodeID:    node.ID,
		Status:    StatusRunning,
		StartTime: startTime,
	}

	e.emitEvent(PipelineEvent{
		NodeID:    node.ID,
		Type:      "node_started",
		Timestamp: startTime,
	})

	// Récupérer la fonction
	fn, err := e.funcRegistry.Get(node.Function)
	if err != nil {
		result.Status = StatusFailed
		result.Error = fmt.Sprintf("function not found: %s", node.Function)
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(startTime)
		return result
	}

	// Construire les paramètres d'entrée
	params := make(map[string]interface{})

	// Ajouter les résultats des nœuds précédents
	for k, v := range previousResults {
		params[k] = v
	}

	// Ajouter les données du trigger
	for k, v := range triggerData {
		params[k] = v
	}

	// Ajouter la configuration du nœud
	for k, v := range node.Config {
		params[k] = v
	}

	// Appeler la fonction
	// Note: Cette partie dépend du type de fonction
	// Une implémentation plus sophistiquée utiliserait la réflexion
	output, err := e.callFunction(fn, params)
	if err != nil {
		result.Status = StatusFailed
		result.Error = err.Error()
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(startTime)
		return result
	}

	result.Status = StatusSuccess
	result.Output = output
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(startTime)

	return result
}

// callFunction appelle une fonction avec ses paramètres.
// Un recover protège le serveur : un handler sans client (ex: mqtt_publish/
// opcua_read sans connexion) renvoie une erreur au lieu de planter le process.
func (e *Engine) callFunction(fn *functions.Function, params map[string]interface{}) (output interface{}, err error) {
	if fn.Handler == nil {
		return nil, fmt.Errorf("function %s has no handler", fn.Name)
	}
	defer func() {
		if r := recover(); r != nil {
			output = nil
			err = fmt.Errorf("function %s panicked: %v", fn.Name, r)
		}
	}()
	log.Printf("[PIPELINE] Calling function: %s", fn.Name)
	return fn.Handler(params)
}

// emitEvent émet un événement
func (e *Engine) emitEvent(event PipelineEvent) {
	e.stopMu.Lock()
	stopped := e.stopped
	e.stopMu.Unlock()
	if stopped {
		return
	}
	select {
	case e.eventCh <- event:
	default:
		log.Printf("[PIPELINE] Event channel full, dropping event")
	}
}

// Start démarre le moteur (goroutines de monitoring)
func (e *Engine) Start() {
	log.Printf("[PIPELINE] Engine started with %d pipelines", len(e.pipelineReg.List()))
}

// Stop arrête le moteur
func (e *Engine) Stop() {
	e.stopMu.Lock()
	e.stopped = true
	e.stopMu.Unlock()
	close(e.stopCh)
	log.Printf("[PIPELINE] Engine stopped")
}

// GetRegistry retourne le registre des pipelines
func (e *Engine) GetRegistry() *Registry {
	return e.pipelineReg
}
