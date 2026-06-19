// cmd/server/main.go
// Standalone HTTP API server for the MINDSET web UI.
// Serves functions, connectors, pipelines (list + save-as-YAML), the technical
// graph and the domain knowledge graph — with NO OPC-UA and NO MQTT required.
// This is the backend the React UI (frontend/pipeline-builder) talks to in dev,
// proxied via Vite at /api -> http://localhost:8080.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"

	"github.com/MindsetAdmin/mindset-data-edge/internal/config"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/calculates"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/conditions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/connectors"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/outputs"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/transforms"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
)

// server holds the shared dependencies for the HTTP handlers.
type server struct {
	funcRegistry *functions.Registry
	kg           *kg.KnowledgeGraph
	pipelinesDir string
	hourlyRate   float64
}

func main() {
	cfgPath := flag.String("config", "config/agent.yaml", "path to the agent config")
	dbPath := flag.String("db", "./data/mindset.db", "path to the mindset SQLite database")
	pipelinesDir := flag.String("pipelines", "config/pipelines", "directory holding pipeline YAML files")
	addr := flag.String("addr", ":8080", "HTTP listen address (matches the Vite dev proxy)")
	flag.Parse()

	// Config is optional — fall back to a sane default hourly rate if absent.
	hourlyRate := 85.0
	if cfg, err := config.LoadConfig(*cfgPath); err != nil {
		log.Printf("[API] ⚠️ Could not load %s (%v); using defaults", *cfgPath, err)
	} else if cfg.Cost.HourlyCost > 0 {
		hourlyRate = cfg.Cost.HourlyCost
	}

	kgInstance, err := kg.NewKnowledgeGraph(*dbPath)
	if err != nil {
		log.Fatalf("[API] Failed to open KG at %s: %v", *dbPath, err)
	}
	defer kgInstance.Close()

	// Best-effort MQTT connection so "Run" can actually execute MQTT handlers.
	// If no broker is reachable, those handlers fail gracefully at run time.
	mqttClient := connectMQTT("tcp://localhost:1883")
	if mqttClient != nil {
		log.Printf("[API] MQTT connected — real execution enabled for MQTT handlers")
		defer mqttClient.Disconnect(250)
	} else {
		log.Printf("[API] No MQTT broker — MQTT handlers will error if a pipeline is run")
	}

	srv := &server{
		funcRegistry: buildRegistry(hourlyRate, kgInstance, mqttClient),
		kg:           kgInstance,
		pipelinesDir: *pipelinesDir,
		hourlyRate:   hourlyRate,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/functions", srv.handleFunctions)
	mux.HandleFunc("/api/connectors", srv.handleConnectors)
	mux.HandleFunc("/api/pipelines", srv.handlePipelines)          // GET list, POST save
	mux.HandleFunc("/api/pipelines/{id}/run", srv.handleRunPipeline) // POST execute
	mux.HandleFunc("/api/kg/technical", srv.handleTechnicalGraph)
	mux.HandleFunc("/api/kg/domain", srv.handleDomainGraph)
	mux.HandleFunc("/api/stats", srv.handleStats)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})

	log.Printf("[API] Config: %s | DB: %s | Pipelines: %s", *cfgPath, *dbPath, *pipelinesDir)
	log.Printf("[API] Registered %d functions", len(srv.funcRegistry.List()))
	log.Printf("[API] Listening on http://localhost%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, withCORS(mux)))
}

// buildRegistry registers every function. MQTT handlers get a real client when
// one is available (for real execution); OPC-UA stays nil and errors gracefully
// if run without a live server. Modbus/SQL are demo stubs.
func buildRegistry(hourlyRate float64, kgInstance *kg.KnowledgeGraph, mqttClient mqtt.Client) *functions.Registry {
	reg := functions.NewRegistry()

	// Connectors
	reg.Register(connectors.NewOPCUAReadHandler(nil).GetFunction())
	reg.Register(connectors.NewMQTTSubscribeHandler(mqttClient).GetFunction())
	reg.Register(stubConnector("modbus_read", "Lire depuis un équipement Modbus (démo)"))
	reg.Register(stubConnector("sql_query", "Interroger une base SQL (démo)"))

	// Transforms
	reg.Register(transforms.NewStateMachineHandler().GetFunction())
	reg.Register(transforms.NewUNSMapperHandler("local").GetFunction())
	reg.Register(transforms.NewFilterHandler().GetFunction())
	// Passthrough transforms referenced by the predefined pipelines.
	reg.Register(passthroughTransform("parse_json", "Parse un payload JSON (démo)"))
	reg.Register(passthroughTransform("enrich_context", "Ajoute du contexte (démo)"))

	// Calculates
	reg.Register(calculates.NewDurationHandler().GetFunction())
	reg.Register(calculates.NewCostHandler(hourlyRate).GetFunction())

	// Conditions
	reg.Register(conditions.NewThresholdHandler().GetFunction())

	// Outputs
	reg.Register(outputs.NewMQTTPublishHandler(mqttClient).GetFunction())
	if kgInstance != nil {
		reg.Register(outputs.NewKGSaveHandler(kgInstance).GetFunction())
	}

	return reg
}

// passthroughTransform returns a transform that passes its params through —
// used for demo functions not yet implemented as real handlers.
func passthroughTransform(name, desc string) *functions.Function {
	return &functions.Function{
		Name:        name,
		Type:        functions.TypeTransform,
		Description: desc,
		Handler: func(p map[string]interface{}) (interface{}, error) {
			return p, nil
		},
	}
}

// stubConnector returns a metadata-only connector whose handler errors if run.
func stubConnector(name, desc string) *functions.Function {
	return &functions.Function{
		Name:        name,
		Type:        functions.TypeConnector,
		Description: desc,
		Handler: func(map[string]interface{}) (interface{}, error) {
			return nil, fmt.Errorf("%s is a demo stub (not implemented)", name)
		},
	}
}

// connectMQTT returns a connected client, or nil if no broker is reachable.
func connectMQTT(broker string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("mindset-api-server").
		SetConnectTimeout(2 * time.Second)
	c := mqtt.NewClient(opts)
	tok := c.Connect()
	if tok.WaitTimeout(3*time.Second) && tok.Error() == nil {
		return c
	}
	return nil
}

// --- Handlers ---------------------------------------------------------------

func (s *server) handleFunctions(w http.ResponseWriter, r *http.Request) {
	typeParam := r.URL.Query().Get("type")
	var list []*functions.FunctionInfo

	if typeParam != "" {
		fnType, ok := parseFunctionType(typeParam)
		if !ok {
			http.Error(w, "Invalid type parameter", http.StatusBadRequest)
			return
		}
		list = s.funcRegistry.ListFunctionsByType(fnType)
	} else {
		list = s.funcRegistry.ListFunctions()
	}

	writeJSON(w, map[string]interface{}{"functions": list, "total": len(list)})
}

func (s *server) handleConnectors(w http.ResponseWriter, r *http.Request) {
	connList := s.funcRegistry.ListFunctionsByType(functions.TypeConnector)
	writeJSON(w, map[string]interface{}{"connectors": connList, "total": len(connList)})
}

func (s *server) handlePipelines(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listPipelines(w, r)
	case http.MethodPost:
		s.savePipeline(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *server) listPipelines(w http.ResponseWriter, r *http.Request) {
	loader := pipeline.NewLoader(s.pipelinesDir)
	pipelines, err := loader.LoadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"pipelines": pipelines, "total": len(pipelines)})
}

func (s *server) savePipeline(w http.ResponseWriter, r *http.Request) {
	var p pipeline.Pipeline
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation mirrors the loader's rules (loader.go).
	if strings.TrimSpace(p.ID) == "" {
		http.Error(w, "pipeline id is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(p.Name) == "" {
		http.Error(w, "pipeline name is required", http.StatusBadRequest)
		return
	}

	filename := sanitizeFilename(p.ID)
	if filename == "" {
		http.Error(w, "pipeline id produces an empty filename", http.StatusBadRequest)
		return
	}

	// A pipeline saved from the builder is enabled by default so the agent loads it.
	p.Enabled = true

	data, err := yaml.Marshal(&p)
	if err != nil {
		http.Error(w, "failed to serialize pipeline: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(s.pipelinesDir, 0o755); err != nil {
		http.Error(w, "failed to create pipelines dir: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outPath := filepath.Join(s.pipelinesDir, filename+".yaml")
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		http.Error(w, "failed to write pipeline: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[API] Saved pipeline %q -> %s", p.ID, outPath)
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]interface{}{"status": "saved", "id": p.ID, "file": outPath})
}

// handleRunPipeline executes a pipeline by id and returns the ExecutionResult
// (per-node status + timing). Real execution: handlers actually run.
func (s *server) handleRunPipeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	loader := pipeline.NewLoader(s.pipelinesDir)
	pipelines, err := loader.LoadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var target *pipeline.Pipeline
	for _, p := range pipelines {
		if p != nil && p.ID == id {
			target = p
			break
		}
	}
	if target == nil {
		http.Error(w, fmt.Sprintf("pipeline %s not found", id), http.StatusNotFound)
		return
	}

	eng := pipeline.NewEngine(s.funcRegistry)
	if err := eng.RegisterPipeline(target); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := eng.Execute(context.Background(), target.ID, map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[API] Ran pipeline %q -> %s", target.ID, result.Status)
	writeJSON(w, result)
}

func (s *server) handleTechnicalGraph(w http.ResponseWriter, r *http.Request) {
	loader := pipeline.NewLoader(s.pipelinesDir)
	pipelines, err := loader.LoadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reg := pipeline.NewRegistry()
	for _, p := range pipelines {
		if p != nil {
			_ = reg.Register(p)
		}
	}

	graph, err := s.kg.GetTechnicalGraph(reg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, graph)
}

func (s *server) handleDomainGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := s.kg.GetFullGraph()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, graph)
}

func (s *server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"functions":              len(s.funcRegistry.List()),
		"pipelines":              0,
		"kg_nodes":               0,
		"kg_edges":               0,
		"micro_stops":            0,
		"total_downtime_seconds": 0.0,
		"estimated_cost_eur":     0.0,
		"hourly_cost":            s.hourlyRate,
	}

	loader := pipeline.NewLoader(s.pipelinesDir)
	if pipelines, err := loader.LoadAll(); err == nil {
		stats["pipelines"] = len(pipelines)
	}

	if graph, err := s.kg.GetFullGraph(); err == nil {
		stats["kg_nodes"] = len(graph.Nodes)
		stats["kg_edges"] = len(graph.Edges)

		microStops := 0
		downtime := 0.0
		for _, n := range graph.Nodes {
			if n.Type == "Event" {
				microStops++
				if d, ok := n.Properties["duration_seconds"].(float64); ok {
					downtime += d
				}
			}
		}
		stats["micro_stops"] = microStops
		stats["total_downtime_seconds"] = downtime
		stats["estimated_cost_eur"] = (downtime / 3600.0) * s.hourlyRate
	}

	writeJSON(w, stats)
}

// --- Helpers ----------------------------------------------------------------

func parseFunctionType(s string) (functions.FunctionType, bool) {
	switch s {
	case "connector":
		return functions.TypeConnector, true
	case "transform":
		return functions.TypeTransform, true
	case "calculate":
		return functions.TypeCalculate, true
	case "condition":
		return functions.TypeCondition, true
	case "output":
		return functions.TypeOutput, true
	default:
		return "", false
	}
}

// sanitizeFilename keeps only safe characters, preventing path traversal.
func sanitizeFilename(s string) string {
	var b strings.Builder
	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '-', c == '_':
			b.WriteRune(c)
		default:
			b.WriteRune('_')
		}
	}
	return strings.Trim(b.String(), "_")
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[API] encode error: %v", err)
	}
}

// withCORS allows the Vite dev server (or any origin) to call the API directly,
// in case requests aren't going through the proxy.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
