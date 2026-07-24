// cmd/agent/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/config"
	"github.com/MindsetAdmin/mindset-data-edge/internal/connections"
	"github.com/MindsetAdmin/mindset-data-edge/internal/discovery"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/calculates"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/conditions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/connectors"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/outputs"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/transforms"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
	"github.com/MindsetAdmin/mindset-data-edge/internal/mqtt"
	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
	"github.com/MindsetAdmin/mindset-data-edge/internal/rules"
	"github.com/MindsetAdmin/mindset-data-edge/internal/uns"
)

func main() {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║     MINDSET DATA — Edge Agent v0     ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()
	time.Sleep(time.Second)

	// Load config
	cfg, err := config.LoadConfig("config/agent.yaml")
	if err != nil {
		log.Fatalf("[CONFIG] Failed to load config: %v", err)
	}
	log.Printf("[CONFIG] Site: %s (%s)", cfg.Site.Name, cfg.Site.ID)
	log.Printf("[CONFIG] OPC-UA endpoint: %s", cfg.OpcUA.Endpoint)
	log.Printf("[CONFIG] Cost: %.2f €/h", cfg.Cost.HourlyCost)

	// Context with graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n[AGENT] Shutting down gracefully...")
		cancel()
	}()

	// ── STEP 0: MQTT Publisher ──────────────────────────────────────────
	fmt.Println("\n[MQTT] Initializing MQTT publisher...")

	var mqttPub *mqtt.Publisher = nil
	mqttPub, err = mqtt.NewPublisher("tcp://localhost:1883", cfg.Site.ID)
	if err != nil {
		log.Printf("[MQTT] ⚠️ Warning: failed to connect to broker: %v", err)
		log.Printf("[MQTT] Continuing without MQTT (raw data won't be published)")
	} else {
		log.Printf("[MQTT] ✅ Connected to broker")
		defer mqttPub.Disconnect()
	}

	// ── STEP 0.5: UNS Contextualizer ────────────────────────────────────
	fmt.Println("\n[UNS] Initializing data contextualizer...")

	// In UI-controlled mode (opcua.auto_connect=false) the API server owns the
	// OPC-UA session and does the ISA-95 normalization, so the agent must NOT also
	// run a contextualizer — otherwise mindset/site/# would be published twice.
	var contextualizer *uns.Contextualizer = nil
	if mqttPub != nil && cfg.OpcUA.AutoConnect {
		mapper := uns.NewMapper(cfg.Site.ID)
		contextualizer, err = uns.NewContextualizer("tcp://localhost:1883", cfg.Site.ID, mapper)
		if err != nil {
			log.Printf("[UNS] ⚠️ Warning: failed to start contextualizer: %v", err)
			log.Printf("[UNS] Continuing without UNS enrichment")
		} else {
			if err := contextualizer.Start(); err != nil {
				log.Printf("[UNS] ⚠️ Warning: failed to subscribe: %v", err)
			} else {
				log.Printf("[UNS] ✅ Contextualizer started — publishing to mindset/site/#")
				defer contextualizer.Stop()
			}
		}
	}

	// ── STEP 0.6: Rules Engine ──────────────────────────────────────────
	fmt.Println("\n[RULES] Initializing rules engine...")

	var rulesEngine *rules.Engine = nil
	if mqttPub != nil {
		rulesEngine, err = rules.NewEngine("tcp://localhost:1883", cfg.Site.ID)
		if err != nil {
			log.Printf("[RULES] ⚠️ Warning: failed to create engine: %v", err)
		} else {
			if err := rulesEngine.Start(); err != nil {
				log.Printf("[RULES] ⚠️ Warning: failed to start engine: %v", err)
			} else {
				log.Printf("[RULES] ✅ Engine started — detecting micro-stops")
				defer rulesEngine.Stop()
			}
		}
	}

	// ── STEP 0.7: Knowledge Graph ───────────────────────────────────────
	fmt.Println("\n[KG] Initializing Knowledge Graph...")

	kgInstance, err := kg.NewKnowledgeGraph("./data/mindset.db")
	if err != nil {
		log.Printf("[KG] ⚠️ Warning: failed to create KG: %v", err)
	} else {
		defer kgInstance.Close()
		log.Printf("[KG] ✅ Knowledge Graph ready")

		kgSub, err := kg.NewKGSubscriber("tcp://localhost:1883", "mindset-agent-kg", kgInstance)
		if err != nil {
			log.Printf("[KG] ⚠️ Warning: failed to create subscriber: %v", err)
		} else {
			if err := kgSub.Start(); err != nil {
				log.Printf("[KG] ⚠️ Warning: failed to start subscriber: %v", err)
			} else {
				defer kgSub.Stop()
				log.Printf("[KG] ✅ Subscriber started")
			}
		}
	}

	// ── STEP 0.75: SQL Connections Registry ─────────────────────────────
	fmt.Println("\n[CONNECTIONS] Loading SQL connection definitions...")

	connCfg, err := connections.LoadConfig("config/connections.yaml")
	if err != nil {
		log.Printf("[CONNECTIONS] No config/connections.yaml (%v); starting with an empty connection set", err)
		connCfg = &connections.Config{}
	}
	connReg := connections.NewRegistry(connCfg)
	defer connReg.CloseAll()
	if kgInstance != nil {
		if records, err := kgInstance.Store().ListConnections(); err != nil {
			log.Printf("[CONNECTIONS] Could not load persisted connections: %v", err)
		} else {
			for _, rec := range records {
				if err := connReg.Add(rec.ConnectionConfig); err != nil {
					log.Printf("[CONNECTIONS] Skipping invalid persisted connection %q: %v", rec.ID, err)
				}
			}
		}
	}

	// ── STEP 0.8: Functions Registry ────────────────────────────────────
	fmt.Println("\n[FUNCTIONS] Registering functions...")

	funcRegistry := functions.NewRegistry()

	// Connectors (nécessitent un client MQTT)
	if mqttPub != nil {
		// Note: mqtt_publish et mqtt_subscribe nécessitent un client MQTT
		// À implémenter si tu as un client MQTT séparé
		log.Printf("[FUNCTIONS] MQTT functions available")
	}
	funcRegistry.Register(connectors.NewSQLQueryHandler(connReg).GetFunction())

	// Transforms
	stateMachineHandler := transforms.NewStateMachineHandler()
	funcRegistry.Register(stateMachineHandler.GetFunction())

	filterHandler := transforms.NewFilterHandler()
	funcRegistry.Register(filterHandler.GetFunction())

	// Calculates
	durationHandler := calculates.NewDurationHandler()
	funcRegistry.Register(durationHandler.GetFunction())

	hourlyRate := 85.0
	if cfg.Cost.HourlyCost > 0 {
		hourlyRate = cfg.Cost.HourlyCost
	}
	costHandler := calculates.NewCostHandler(hourlyRate)
	funcRegistry.Register(costHandler.GetFunction())

	// Conditions
	thresholdHandler := conditions.NewThresholdHandler()
	funcRegistry.Register(thresholdHandler.GetFunction())

	// Outputs
	if kgInstance != nil {
		kgSaveHandler := outputs.NewKGSaveHandler(kgInstance)
		funcRegistry.Register(kgSaveHandler.GetFunction())
	}

	// Afficher toutes les fonctions enregistrées
	fmt.Printf("\n[FUNCTIONS] Registered %d functions:\n", len(funcRegistry.List()))
	for _, fn := range funcRegistry.List() {
		fmt.Printf("  ✓ %s (%s) - %s\n", fn.Name, fn.Type, fn.Description)
	}

	// ── STEP 0.9: Pipeline Engine ──────────────────────────────────────
	fmt.Println("\n[PIPELINE] Initializing pipeline engine...")

	pipelineEngine := pipeline.NewEngine(funcRegistry)

	// Charger les pipelines depuis config/pipelines/
	loader := pipeline.NewLoader("config/pipelines")
	pipelines, err := loader.LoadAll()
	if err != nil {
		log.Printf("[PIPELINE] ⚠️ Warning: failed to load pipelines: %v", err)
	} else {
		for _, p := range pipelines {
			if p == nil {
				continue
			}
			if err := pipelineEngine.RegisterPipeline(p); err != nil {
				log.Printf("[PIPELINE] ⚠️ Warning: failed to register pipeline %s: %v", p.Name, err)
			} else {
				log.Printf("[PIPELINE] ✅ Registered pipeline: %s (%s)", p.Name, p.ID)
			}
		}
	}

	// ── STEP 0.10: Technical Knowledge Graph ──────────────────────────
	fmt.Println("\n[KG] Building technical knowledge graph...")

	if pipelineEngine != nil && kgInstance != nil {
		techGraph, err := kgInstance.GetTechnicalGraph(pipelineEngine.GetRegistry())
		if err != nil {
			log.Printf("[KG] ⚠️ Warning: failed to build technical graph: %v", err)
		} else {
			log.Printf("[KG] ✅ Technical graph built with %d nodes and %d edges",
				len(techGraph.Nodes), len(techGraph.Edges))

			// Optionnel: sauvegarder le graphe dans un fichier
			// data, _ := json.MarshalIndent(techGraph, "", "  ")
			// os.WriteFile("data/technical_graph.json", data, 0644)
		}
	}

	// HTTP is owned by cmd/server (the standalone API the web UI talks to).
	// The agent is purely the edge runtime; run `go run ./cmd/server` for the UI.

	// ── UI-controlled mode gate ─────────────────────────────────────────
	// When opcua.auto_connect is false (the default), the OPC-UA connection,
	// discovery and subscription are driven from the web UI via cmd/server. The
	// agent stays alive as the edge runtime (MQTT, rules engine, KG enrichment)
	// but does NOT auto-connect to OPC-UA. Deferred cleanups still run on return.
	if !cfg.OpcUA.AutoConnect {
		log.Printf("[DISCOVERY] OPC-UA auto-connect disabled (opcua.auto_connect=false)")
		log.Printf("[DISCOVERY] Connection is driven from the web UI (cmd/server). Rules + KG keep running.")
		fmt.Println("\n[AGENT] Running in UI-controlled mode. Press Ctrl+C to stop.")
		<-ctx.Done()
		fmt.Println("\n[AGENT] Stopped by user.")
		return
	}

	// ── STEP 1: OPC-UA Discovery ────────────────────────────────────────
	fmt.Println("\n[DISCOVERY] Starting OPC-UA auto-discovery...")
	fmt.Printf("[DISCOVERY] Connecting to: %s\n\n", cfg.OpcUA.Endpoint)

	opcua := discovery.NewOPCUADiscovery(cfg.OpcUA.Endpoint, mqttPub)

	if err := opcua.Connect(ctx); err != nil {
		log.Printf("[DISCOVERY] ⚠️ Connection failed: %v", err)
		log.Printf("[DISCOVERY] Make sure the OPC-UA server is running. Agent will idle (rules + KG keep running).")
		fmt.Println("\n[AGENT] Running (OPC-UA unavailable). Press Ctrl+C to stop.")
		<-ctx.Done()
		fmt.Println("\n[AGENT] Stopped by user.")
		return
	}
	defer opcua.Disconnect(ctx)

	// ── STEP 2: Browse node tree ────────────────────────────────────────
	fmt.Println("\n[DISCOVERY] Browsing node tree...\n")
	tags, err := opcua.BrowseNodeTree(ctx)
	if err != nil {
		log.Printf("[DISCOVERY] ⚠️ Browse failed: %v", err)
		fmt.Println("\n[AGENT] Running (browse failed). Press Ctrl+C to stop.")
		<-ctx.Done()
		fmt.Println("\n[AGENT] Stopped by user.")
		return
	}

	fmt.Printf("\n[DISCOVERY] ✅ Found %d tags total\n", len(tags))

	// ── STEP 3: Print summary ───────────────────────────────────────────
	boolCount, floatCount, intCount, otherCount := 0, 0, 0, 0
	for _, tag := range tags {
		switch tag.DataType {
		case "Boolean":
			boolCount++
		case "Float", "Double":
			floatCount++
		case "Int16", "Int32", "Int64", "UInt16", "UInt32", "UInt64", "SByte", "Byte":
			intCount++
		default:
			otherCount++
		}
	}

	fmt.Printf("\n[SUMMARY] Tag types found:\n")
	fmt.Printf("  Boolean  : %d  (→ candidate machine states, alarms)\n", boolCount)
	fmt.Printf("  Float    : %d  (→ candidate speeds, pressures, temperatures)\n", floatCount)
	fmt.Printf("  Integer  : %d  (→ candidate counters, setpoints)\n", intCount)
	fmt.Printf("  Other    : %d\n", otherCount)

	// ── STEP 4: Live subscription ───────────────────────────────────────
	fmt.Println("\n[SUBSCRIBE] Starting live data subscription...")
	fmt.Println("[SUBSCRIBE] Watching for value changes:\n")

	err = opcua.Subscribe(ctx, tags, 500*time.Millisecond, func(tag discovery.Tag) {
		fmt.Printf("  📡 CHANGE | %-40s | %v (%s)\n",
			tag.Name,
			tag.Value,
			tag.DataType,
		)
	})
	if err != nil {
		log.Printf("[SUBSCRIBE] Warning: subscription failed: %v", err)
	}

	// Start the watcher in a goroutine — non-blocking
	go opcua.WatchForChanges(
		ctx,
		tags,
		30*time.Second,
		func(change discovery.TagChange, allTags []discovery.Tag) {
			fmt.Printf("\n[WATCH] ⚡ Tag change detected!\n")
			fmt.Printf("  Added   : %d new tags\n", len(change.Added))
			fmt.Printf("  Removed : %d tags gone\n", len(change.Removed))

			for _, t := range change.Added {
				fmt.Printf("  + %s (%s)\n", t.Name, t.DataType)
			}
			for _, t := range change.Removed {
				fmt.Printf("  - %s\n", t.Name)
			}
			fmt.Printf("[WATCH] Rebuild subscription with %d tags\n", len(allTags))
		},
	)

	// ── STEP 5: Wait for demo ───────────────────────────────────────────
	fmt.Println("\n[AGENT] Demo running...")
	fmt.Println("[AGENT] Press Ctrl+C to stop\n")

	select {
	case <-ctx.Done():
		fmt.Println("\n[AGENT] Stopped by user.")
	case <-time.After(20 * time.Minute):
		fmt.Println("\n[AGENT] 20-minute demo complete.")
	}
}

func countNodesByType(graph *kg.TechnicalGraph) map[string]int {
	counts := make(map[string]int)
	for _, node := range graph.Nodes {
		counts[string(node.Type)]++
	}
	return counts
}
