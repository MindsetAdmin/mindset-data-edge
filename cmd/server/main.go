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
	"sort"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"

	"github.com/MindsetAdmin/mindset-data-edge/internal/config"
	"github.com/MindsetAdmin/mindset-data-edge/internal/connections"
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
	tags         *TagRegistry
	topics       *TopicRegistry
	states       *StateTracker
	live         *LiveHub
	opcua        *OPCUAManager
	cfg          *config.Config
	mqttClient   mqtt.Client
	connReg      *connections.Registry
	startTime    time.Time
}

func main() {
	cfgPath := flag.String("config", "config/agent.yaml", "path to the agent config")
	dbPath := flag.String("db", "./data/mindset.db", "path to the mindset SQLite database")
	pipelinesDir := flag.String("pipelines", "config/pipelines", "directory holding pipeline YAML files")
	addr := flag.String("addr", ":8080", "HTTP listen address (matches the Vite dev proxy)")
	mcpStdio := flag.Bool("mcp-stdio", false, "run only the MCP server over stdio (for local MCP clients like Claude Desktop's mcpServers config) instead of the HTTP API — no port is bound in this mode")
	connectionsPath := flag.String("connections", "config/connections.yaml", "path to the SQL connections config")
	flag.Parse()

	// Config is optional — fall back to defaults if absent.
	cfg, cfgErr := config.LoadConfig(*cfgPath)
	if cfgErr != nil {
		log.Printf("[API] ⚠️ Could not load %s (%v); using defaults", *cfgPath, cfgErr)
		cfg = &config.Config{}
	}
	hourlyRate := 85.0
	if cfg.Cost.HourlyCost > 0 {
		hourlyRate = cfg.Cost.HourlyCost
	}
	broker := cfg.Mqtt.Broker
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	kgInstance, err := kg.NewKnowledgeGraph(*dbPath)
	if err != nil {
		log.Fatalf("[API] Failed to open KG at %s: %v", *dbPath, err)
	}
	defer kgInstance.Close()

	// Best-effort MQTT connection so "Run" can actually execute MQTT handlers.
	// A distinct client ID in -mcp-stdio mode: the paho MQTT spec disconnects
	// whichever connection loses when two clients share a client ID, which
	// would otherwise silently kick the already-running HTTP instance the
	// moment a Claude Desktop-spawned stdio process connects.
	mqttClientID := "mindset-api-server"
	if *mcpStdio {
		mqttClientID = "mindset-mcp-stdio"
	}
	mqttClient := connectMQTT(broker, mqttClientID)
	if mqttClient != nil {
		log.Printf("[API] MQTT connected (%s)", broker)
		defer mqttClient.Disconnect(250)
	} else {
		log.Printf("[API] No MQTT broker at %s — live data and Run will be limited", broker)
	}

	// Automatic Knowledge Graph enrichment from events.
	if mqttClient != nil {
		if kgSub, err := kg.NewKGSubscriber(broker, mqttClientID+"-kg", kgInstance); err != nil {
			log.Printf("[API] KG auto-enrichment unavailable: %v", err)
		} else if err := kgSub.Start(); err != nil {
			log.Printf("[API] KG subscriber failed to start: %v", err)
		} else {
			log.Printf("[API] KG auto-enrichment active (mindset/events/micro-stop)")
			defer kgSub.Stop()
		}
	}

	// Live data hub: one mindset/# subscription feeding the tag registry, topic
	// stats (rates), and machine state tracking (Running/Stopped + transitions).
	tagReg := NewTagRegistry(kgInstance.Store().DB())
	topicReg := NewTopicRegistry()
	stateTracker := NewStateTracker()
	wsHubInstance := newWSHub()
	hub := NewLiveHub(tagReg, topicReg, stateTracker)
	hub.broadcast = wsHubInstance.broadcast // push live updates to WebSocket clients
	if mqttClient != nil {
		if err := hub.Start(mqttClient); err != nil {
			log.Printf("[API] Live data hub failed: %v", err)
		} else {
			log.Printf("[API] Live data active (tags, topics, machine state, dashboard pins from mindset/#)")
		}
	}

	// Dynamic, frontend-driven OPC-UA control plane. Idle until the UI calls
	// /api/opcua/connect; publishes selected tags to the same broker the LiveHub
	// already watches, so discovered tags surface via /api/tags + WebSocket.
	opcuaMgr := NewOPCUAManager(broker, cfg, kgInstance)

	// SQL connection registry (Day 6 of docs/mysql_connector.md). YAML-seeded
	// connections load first, then persisted ones from data/mindset.db
	// (created via POST /api/connections) — the latter win on id conflicts
	// since they reflect the most recent user edit.
	connCfg, connCfgErr := connections.LoadConfig(*connectionsPath)
	if connCfgErr != nil {
		log.Printf("[API] No %s (%v); starting with an empty connection set", *connectionsPath, connCfgErr)
		connCfg = &connections.Config{}
	}
	connReg := connections.NewRegistry(connCfg)
	defer connReg.CloseAll()
	if records, err := kgInstance.Store().ListConnections(); err != nil {
		log.Printf("[API] Could not load persisted connections: %v", err)
	} else {
		for _, rec := range records {
			if err := connReg.Add(rec.ConnectionConfig); err != nil {
				log.Printf("[API] Skipping invalid persisted connection %q: %v", rec.ID, err)
			}
		}
	}

	srv := &server{
		funcRegistry: buildRegistry(hourlyRate, mqttClient, connReg),
		kg:           kgInstance,
		pipelinesDir: *pipelinesDir,
		hourlyRate:   hourlyRate,
		tags:         tagReg,
		topics:       topicReg,
		states:       stateTracker,
		live:         hub,
		opcua:        opcuaMgr,
		cfg:          cfg,
		mqttClient:   mqttClient,
		connReg:      connReg,
		startTime:    time.Now(),
	}

	// -mcp-stdio: no HTTP server at all — Claude Desktop (or any local MCP
	// client that launches a subprocess rather than connecting to a URL)
	// owns this process's lifecycle via stdin/stdout. Everything above (KG,
	// MQTT/LiveHub, connections registry) stays wired up exactly as in HTTP
	// mode, so the same tools return real, live data either way.
	if *mcpStdio {
		log.Printf("[API] Running MCP server over stdio (no HTTP listener, no port bound)")
		if err := runMCPStdio(srv); err != nil {
			log.Fatalf("[API] MCP stdio server exited with error: %v", err)
		}
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/functions", srv.handleFunctions)
	mux.HandleFunc("/api/connectors", srv.handleConnectors)
	mux.HandleFunc("/api/pipelines", srv.handlePipelines)           // GET list, POST save
	mux.HandleFunc("/api/pipelines/examples", srv.handleExamplePipelines) // GET templates
	mux.HandleFunc("/api/pipelines/{id}/run", srv.handleRunPipeline) // POST execute
	mux.HandleFunc("/api/pipelines/{id}", srv.handleDeletePipeline)  // DELETE
	mux.HandleFunc("/api/tags", srv.handleTags)
	mux.HandleFunc("/api/machines", srv.handleMachines)
	mux.HandleFunc("/api/topics", srv.handleTopics)
	mux.HandleFunc("/api/config", srv.handleConfig)
	mux.HandleFunc("/api/opcua/connect", srv.handleOpcuaConnect)       // POST: dynamic connect
	mux.HandleFunc("/api/opcua/discover", srv.handleOpcuaDiscover)     // GET: browse tags
	mux.HandleFunc("/api/opcua/subscribe", srv.handleOpcuaSubscribe)   // POST: select tags + modes
	mux.HandleFunc("/api/opcua/disconnect", srv.handleOpcuaDisconnect) // POST: close session
	mux.HandleFunc("/api/opcua/status", srv.handleOpcuaStatus)         // GET: connection status
	mux.HandleFunc("/api/opcua/selections", srv.handleOpcuaSelections) // GET: per-tag routing (governance)
	mux.HandleFunc("/api/dashboard/pins", srv.handleDashboardPins)
	mux.HandleFunc("/api/connections", srv.handleConnections)                   // GET list, POST create
	mux.HandleFunc("/api/connections/{id}/test", srv.handleConnectionTest)      // POST: re-run health check
	mux.HandleFunc("/api/connections/{id}/preview", srv.handleConnectionPreview) // POST: preview a query, capped at 5 rows
	mux.HandleFunc("/api/connections/{id}/discover", srv.handleConnectionDiscover)   // GET: schema browse + canonical-mapping auto-suggest (Track B Phase 1+2)
	mux.HandleFunc("/api/connections/{id}/databases", srv.handleConnectionDatabases) // GET: browse every database + table visible to this connection
	mux.HandleFunc("/api/production/active", srv.handleActiveProduction)          // GET: live active production per machine, OT-linked where resolved (Entry 120)
	mux.HandleFunc("/api/connections/{id}", srv.handleConnectionDelete)         // DELETE
	mux.HandleFunc("/api/kg", srv.handleKG)                        // unified — ?category=business|platform|all
	mux.HandleFunc("/api/kg/technical", srv.handleTechnicalGraph)   // legacy alias → category=platform
	mux.HandleFunc("/api/kg/domain", srv.handleDomainGraph)         // legacy alias → category=business
	mux.HandleFunc("/api/kg/pending", srv.handleKGPending)                   // GET: nodes awaiting validation (v0 structural bootstrap, Entry 95/96)
	mux.HandleFunc("/api/kg/pending/{id}/validate", srv.handleKGValidate)    // POST: confirm an auto-generated node
	mux.HandleFunc("/api/kg/pending/{id}/reject", srv.handleKGReject)        // POST: discard an auto-generated node
	mux.HandleFunc("/api/stats", srv.handleStats)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/ws", wsHubInstance.handle) // WebSocket: live push to the UI
	mountMCP(mux, srv)                              // /mcp — read-only agent tools (Track A, see mcp_server.go)

	log.Printf("[API] Config: %s | DB: %s | Pipelines: %s", *cfgPath, *dbPath, *pipelinesDir)
	log.Printf("[API] Registered %d functions", len(srv.funcRegistry.List()))
	log.Printf("[API] Listening on http://localhost%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, withCORS(mux)))
}

// buildRegistry registers every function. MQTT handlers get a real client when
// one is available (for real execution); OPC-UA stays nil and errors gracefully
// if run without a live server. Modbus is still a demo stub; sql_query runs for
// real against connReg (empty registry = every query errors "unknown connection").
func buildRegistry(hourlyRate float64, mqttClient mqtt.Client, connReg *connections.Registry) *functions.Registry {
	reg := functions.NewRegistry()

	// Connectors
	reg.Register(connectors.NewOPCUAReadHandler(nil).GetFunction())
	reg.Register(connectors.NewMQTTSubscribeHandler(mqttClient).GetFunction())
	reg.Register(stubConnector("modbus_read", "Lire depuis un équipement Modbus (démo)"))
	reg.Register(connectors.NewSQLQueryHandler(connReg).GetFunction())

	// Transforms
	reg.Register(transforms.NewStateMachineHandler().GetFunction())
	reg.Register(transforms.NewFilterHandler().GetFunction())

	// Calculates
	reg.Register(calculates.NewDurationHandler().GetFunction())
	reg.Register(calculates.NewCostHandler(hourlyRate).GetFunction())

	// Conditions
	reg.Register(conditions.NewThresholdHandler().GetFunction())

	// Outputs — kg_save is intentionally NOT registered: the Knowledge Graph
	// enriches itself automatically via the KG subscriber. mqtt_publish is
	// intentionally NOT registered either (removed Entry 119): a pipeline's
	// terminal node result now auto-publishes without an explicit node —
	// see cmd/server/pipeline_output.go.
	reg.Register(outputs.NewDashboardHandler(mqttClient).GetFunction())

	return reg
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
func connectMQTT(broker, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
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

// handleExamplePipelines lists the shipped template pipelines (config/pipelines/
// examples). These are NOT loaded into the engine or the KG — they're starting
// points the user can load into Compose.
func (s *server) handleExamplePipelines(w http.ResponseWriter, r *http.Request) {
	loader := pipeline.NewLoader(filepath.Join(s.pipelinesDir, "examples"))
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

// handleDeletePipeline removes a saved pipeline's YAML file. Uses the same
// sanitizeFilename(id) mapping savePipeline uses to derive the file, so a
// pipeline always deletes the exact file it would have been saved to.
func (s *server) handleDeletePipeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, "pipeline id is required", http.StatusBadRequest)
		return
	}
	filename := sanitizeFilename(id)
	if filename == "" {
		http.Error(w, "pipeline id produces an empty filename", http.StatusBadRequest)
		return
	}

	path := filepath.Join(s.pipelinesDir, filename+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("pipeline %s not found", id), http.StatusNotFound)
		return
	}
	if err := os.Remove(path); err != nil {
		http.Error(w, "failed to delete pipeline: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[API] Deleted pipeline %q -> %s", id, path)
	writeJSON(w, map[string]interface{}{"status": "deleted", "id": id})
}

// handleRunPipeline executes a pipeline by id and returns the ExecutionResult
// (per-node status + timing). Real execution: handlers actually run.
//
// No pipeline is ever auto-triggered by live MQTT messages today — a
// pipeline's trigger.function/config in its YAML is declarative only
// (shown in the Compose UI's ENTRÉE zone), never actually subscribed by the
// engine. This is the only path that runs a pipeline at all. An optional
// JSON body ({"trigger_data": {...}}) lets a caller supply what a real
// trigger message would have carried (e.g. work_center/duration_seconds),
// which is how a manual test/demo run exercises the same params a live
// mqtt_subscribe trigger would eventually provide — an empty/absent body
// keeps the previous no-trigger-data behavior unchanged.
func (s *server) handleRunPipeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	var body struct {
		TriggerData map[string]interface{} `json:"trigger_data"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body) // best-effort — empty/absent body is valid
	}

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

	triggerData := body.TriggerData
	if triggerData == nil {
		triggerData = map[string]interface{}{}
	}
	result, err := eng.Execute(context.Background(), target.ID, triggerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[API] Ran pipeline %q -> %s", target.ID, result.Status)
	s.publishPipelineOutput(target, result)
	writeJSON(w, result)
}

func (s *server) handleDashboardPins(w http.ResponseWriter, r *http.Request) {
	pins := s.live.Pins()
	writeJSON(w, map[string]interface{}{"pins": pins, "total": len(pins)})
}

func (s *server) handleTags(w http.ResponseWriter, r *http.Request) {
	list := s.tags.list()
	writeJSON(w, map[string]interface{}{"tags": list, "total": len(list)})
}

// handleMachines groups discovered tags by work center. When the user has
// configured ISA-95 routing in Connect, the mapping in OPCUAManager is
// authoritative and used first (keyed by NodeID). Otherwise the tag name is
// parsed (segment before the leaf attribute).
func (s *server) handleMachines(w http.ResponseWriter, r *http.Request) {
	wcByNodeID := map[string]string{}
	if s.opcua != nil {
		for _, sel := range s.opcua.SelectionsDetailed() {
			if sel.WorkCenter != "" {
				wcByNodeID[sel.NodeID] = sel.WorkCenter
			}
		}
	}

	groups := map[string][]Tag{}
	for _, t := range s.tags.list() {
		wc := wcByNodeID[t.NodeID]
		if wc == "" {
			var ok bool
			wc, ok = workCenterOf(t.Name)
			if !ok {
				wc = "(autres)"
			}
		}
		groups[wc] = append(groups[wc], t)
	}

	type machine struct {
		WorkCenter string        `json:"work_center"`
		Tags       []Tag         `json:"tags"`
		State      *MachineState `json:"state,omitempty"`
	}
	machines := make([]machine, 0, len(groups))
	for wc, tags := range groups {
		machines = append(machines, machine{WorkCenter: wc, Tags: tags, State: s.states.get(wc)})
	}
	sort.Slice(machines, func(i, j int) bool { return machines[i].WorkCenter < machines[j].WorkCenter })
	writeJSON(w, map[string]interface{}{"machines": machines, "total": len(machines)})
}

func (s *server) handleTopics(w http.ResponseWriter, r *http.Request) {
	connected := s.mqttClient != nil && s.mqttClient.IsConnected()
	list := s.topics.list()
	writeJSON(w, map[string]interface{}{
		"topics":           list,
		"total":            len(list),
		"broker_connected": connected,
		"broker":           brokerOrDefault(s.cfg),
	})
}

// handleConfig exposes a safe subset of agent.yaml so the UI can pre-fill fields.
func (s *server) handleConfig(w http.ResponseWriter, r *http.Request) {
	currency := "EUR"
	writeJSON(w, map[string]interface{}{
		"opcua": map[string]interface{}{
			"endpoint":        s.cfg.OpcUA.Endpoint,
			"security_mode":   orDefault(s.cfg.OpcUA.SecurityMode, "None"),
			"security_policy": orDefault(s.cfg.OpcUA.SecurityPolicy, "None"),
			"timeout":         5000,
		},
		"mqtt": map[string]interface{}{"broker": brokerOrDefault(s.cfg)},
		"cost": map[string]interface{}{"hourly_rate": s.hourlyRate, "currency": currency},
		"site": map[string]interface{}{"id": orDefault(s.cfg.Site.ID, "local-test"), "name": s.cfg.Site.Name, "area": "area1"},
	})
}

func brokerOrDefault(cfg *config.Config) string {
	if cfg != nil && cfg.Mqtt.Broker != "" {
		return cfg.Mqtt.Broker
	}
	return "tcp://localhost:1883"
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// handleKG — unified KG endpoint. Query params:
//   ?category=business  → site fingerprint (Equipment/Event/Cause/Cost/...)
//   ?category=platform  → pipeline topology (Pipeline/Function/Topic/Connection/Dashboard)
//   ?category=all       → both (default)
// For category=platform or category=all, we ensure the platform sub-graph is
// up-to-date by loading pipelines and calling RepopulatePlatform.
func (s *server) handleKG(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "all"
	}
	if category != "business" && category != "platform" && category != "all" {
		http.Error(w, `invalid category — use "business", "platform", or "all"`, http.StatusBadRequest)
		return
	}

	// Refresh the platform sub-graph when the request touches it.
	if category == "platform" || category == "all" {
		reg := s.loadPipelineRegistry()
		if err := s.kg.RepopulatePlatform(reg); err != nil {
			log.Printf("[API] handleKG: RepopulatePlatform failed: %v", err)
		}
	}

	graph, err := s.kg.GetGraph(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, graph)
}

// handleTechnicalGraph — legacy alias for /api/kg?category=platform.
// Returns the old TechnicalGraph shape (nodes/edges without category field) for
// backward compat with any external consumer.
func (s *server) handleTechnicalGraph(w http.ResponseWriter, r *http.Request) {
	reg := s.loadPipelineRegistry()
	graph, err := s.kg.GetTechnicalGraph(reg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, graph)
}

// handleDomainGraph — legacy alias for /api/kg?category=business.
func (s *server) handleDomainGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := s.kg.GetGraph("business")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, graph)
}

// loadPipelineRegistry builds a fresh pipeline registry from disk. Used by KG
// handlers that need the current pipeline topology.
func (s *server) loadPipelineRegistry() *pipeline.Registry {
	reg := pipeline.NewRegistry()
	loader := pipeline.NewLoader(s.pipelinesDir)
	pipelines, err := loader.LoadAll()
	if err != nil {
		log.Printf("[API] loadPipelineRegistry: %v", err)
		return reg
	}
	for _, p := range pipelines {
		if p != nil {
			_ = reg.Register(p)
		}
	}
	return reg
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
		"uptime_seconds":         time.Since(s.startTime).Seconds(),
		"broker_connected":       s.mqttClient != nil && s.mqttClient.IsConnected(),
	}

	loader := pipeline.NewLoader(s.pipelinesDir)
	if pipelines, err := loader.LoadAll(); err == nil {
		stats["pipelines"] = len(pipelines)
	}

	// Stats count business (site fingerprint) nodes/edges only. Platform
	// topology counts are volatile (rebuilt from pipelines) and would inflate
	// the "your factory has X data points" story.
	if graph, err := s.kg.GetGraph("business"); err == nil {
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
