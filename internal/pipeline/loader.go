// internal/pipeline/loader.go
package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Loader charge les pipelines depuis des fichiers YAML
type Loader struct {
	pipelinesDir string
}

// NewLoader crée un nouveau loader
func NewLoader(pipelinesDir string) *Loader {
	return &Loader{
		pipelinesDir: pipelinesDir,
	}
}

// LoadAll charge tous les pipelines YAML du dossier
func (l *Loader) LoadAll() ([]*Pipeline, error) {
	var pipelines []*Pipeline

	files, err := filepath.Glob(filepath.Join(l.pipelinesDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p, err := l.LoadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", file, err)
		}
		if p != nil {
			pipelines = append(pipelines, p)
		}
	}

	return pipelines, nil
}

// LoadFile charge un pipeline depuis un fichier YAML
func (l *Loader) LoadFile(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	// Vérifier si le pipeline est activé
	enabled := true
	// Pour YAML, on peut avoir un champ enabled directement
	// On va le parser manuellement depuis le raw YAML pour éviter de modifier la structure
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err == nil {
		if enabledVal, ok := raw["enabled"]; ok {
			if enabledBool, ok := enabledVal.(bool); ok {
				enabled = enabledBool
			}
		}
	}

	if !enabled {
		return nil, nil // Ignorer les pipelines désactivés
	}

	// Valider le pipeline
	if p.ID == "" {
		return nil, fmt.Errorf("pipeline ID is required")
	}
	if p.Name == "" {
		return nil, fmt.Errorf("pipeline name is required")
	}

	// Timestamps
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	return &p, nil
}
