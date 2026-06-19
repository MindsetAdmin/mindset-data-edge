package functions

import (
	"fmt"
	"log"
	"sync"
)

// Registry enregistre toutes les fonctions disponibles
type Registry struct {
	mu        sync.RWMutex
	functions map[string]*Function
}

// HandlerFunc est la signature attendue pour tous les handlers de fonction.
type HandlerFunc func(params map[string]interface{}) (interface{}, error)

// Function définit une fonction
type Function struct {
	Name        string
	Type        FunctionType
	Description string
	Handler     HandlerFunc
}

// NewRegistry crée un nouveau registre
func NewRegistry() *Registry {
	return &Registry{
		functions: make(map[string]*Function),
	}
}

// Register enregistre une fonction
func (r *Registry) Register(fn *Function) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.functions[fn.Name] = fn
	log.Printf("[FUNCTIONS] Registered: %s (%s)", fn.Name, fn.Type)
}

// Get retourne une fonction par son nom
func (r *Registry) Get(name string) (*Function, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if fn, ok := r.functions[name]; ok {
		return fn, nil
	}
	return nil, fmt.Errorf("function %s not found", name)
}

// List retourne toutes les fonctions
func (r *Registry) List() []*Function {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*Function, 0, len(r.functions))
	for _, fn := range r.functions {
		list = append(list, fn)
	}
	return list
}

// GetFunctionInfo retourne les infos détaillées d'une fonction
func (r *Registry) GetFunctionInfo(name string) (*FunctionInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, ok := r.functions[name]
	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}

	return functionToInfo(fn), nil
}

// functionToInfo converts a Function to FunctionInfo without acquiring any lock.
// Callers must hold at least a read lock.
func functionToInfo(fn *Function) *FunctionInfo {
	return &FunctionInfo{
		Name:        fn.Name,
		Type:        string(fn.Type),
		Category:    getCategory(fn.Type),
		Description: fn.Description,
		Inputs:      getInputs(fn.Name),
		Outputs:     getOutputs(fn.Name),
		Icon:        getIcon(fn.Type),
		Color:       getColor(fn.Type),
	}
}

// ListFunctions retourne toutes les fonctions disponibles
func (r *Registry) ListFunctions() []*FunctionInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*FunctionInfo, 0, len(r.functions))
	for _, fn := range r.functions {
		result = append(result, functionToInfo(fn))
	}
	return result
}

// ListFunctionsByType retourne les fonctions d'un type donné
func (r *Registry) ListFunctionsByType(fnType FunctionType) []*FunctionInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*FunctionInfo, 0)
	for _, fn := range r.functions {
		if fn.Type == fnType {
			result = append(result, functionToInfo(fn))
		}
	}
	return result
}

// --- Helpers ---

func getCategory(t FunctionType) string {
	switch t {
	case TypeConnector:
		return "Connecteurs"
	case TypeTransform:
		return "Transformations"
	case TypeCalculate:
		return "Calculs"
	case TypeCondition:
		return "Conditions"
	case TypeOutput:
		return "Sorties"
	default:
		return "Autres"
	}
}

func getIcon(t FunctionType) string {
	switch t {
	case TypeConnector:
		return "Plug"
	case TypeTransform:
		return "GitMerge"
	case TypeCalculate:
		return "Calculator"
	case TypeCondition:
		return "GitBranch"
	case TypeOutput:
		return "Send"
	default:
		return "Box"
	}
}

func getColor(t FunctionType) string {
	switch t {
	case TypeConnector:
		return "#3b82f6"
	case TypeTransform:
		return "#f59e0b"
	case TypeCalculate:
		return "#10b981"
	case TypeCondition:
		return "#8b5cf6"
	case TypeOutput:
		return "#ef4444"
	default:
		return "#6b7280"
	}
}

func getInputs(name string) []Param {
	// À implémenter selon chaque fonction
	return []Param{}
}

func getOutputs(name string) []Param {
	// À implémenter selon chaque fonction
	return []Param{}
}
