// cmd/agent/init.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/config"
	"github.com/MindsetAdmin/mindset-data-edge/internal/discovery"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/calculates"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/conditions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/outputs"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/transforms"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
	"github.com/MindsetAdmin/mindset-data-edge/internal/mqtt"
	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
	"github.com/MindsetAdmin/mindset-data-edge/internal/rules"
	"github.com/MindsetAdmin/mindset-data-edge/internal/uns"
)

// initMQTT initialise le publisher MQTT
func initMQTT(cfg *config.Config) *mqtt.Publisher {
	fmt.Println("\n[MQTT] Initializing MQTT publisher...")

	mqttPub, err := mqtt.NewPublisher("tcp://localhost:1883", cfg.Site.ID)
	if err != nil {
		log.Printf("[MQTT] ⚠️ Warning: failed to connect to broker: %v", err)
		log.Printf("[MQTT] Continuing without MQTT (raw data won't be published)")
		return nil
	}
	log.Printf("[MQTT] ✅ Connected to broker")
	return mqttPub
}

// initUNS initialise le contextualizer UNS
func initUNS(cfg *config.Config, mqttPub *mqtt.Publisher) *uns.Contextualizer {
	fmt.Println("\n[UNS] Initializing data contextualizer...")

	if mqttPub == nil {
		log.Printf("[UNS] Skipping: MQTT not available")
		return nil
	}

	mapper := uns.NewMapper(cfg.Site.ID)
	contextualizer, err := uns.NewContextualizer("tcp://localhost:1883", cfg.Site.ID, mapper)
	if err != nil {
		log.Printf("[UNS] ⚠️ Warning: failed to start contextualizer: %v", err)
		return nil
	}

	if err := contextualizer.Start(); err != nil {
		log.Printf("[UNS] ⚠️ Warning: failed to subscribe: %v", err)
		return nil
	}

	log.Printf("[UNS] ✅ Contextualizer started — publishing to mindset/site/#")
	return contextualizer
}

// initRulesEngine initialise le moteur de règles
func initRulesEngine(cfg *config.Config, mqttPub *mqtt.Publisher) *rules.Engine {
	fmt.Println("\n[RULES] Initializing rules engine...")

	if mqttPub == nil {
		log.Printf("[RULES] Skipping: MQTT not available")
		return nil
	}

	engine, err := rules.NewEngine("tcp://localhost:1883", cfg.Site.ID)
	if err != nil {
		log.Printf("[RULES] ⚠️ Warning: failed to create engine: %v", err)
		return nil
	}

	if err := engine.Start(); err != nil {
		log.Printf("[RULES] ⚠️ Warning: failed to start engine: %v", err)
		return nil
	}

	log.Printf("[RULES] ✅ Engine started — detecting micro-stops")
	return engine
}

// initKnowledgeGraph initialise le Knowledge Graph
func initKnowledgeGraph() *kg.KnowledgeGraph {
	fmt.Println("\n[KG] Initializing Knowledge Graph...")

	kgInstance, err := kg.NewKnowledgeGraph("./data/mindset.db")
	if err != nil {
		log.Printf("[KG] ⚠️ Warning: failed to create KG: %v", err)
		return nil
	}

	log.Printf("[KG] ✅ Knowledge Graph ready")

	// Démarrer le subscriber KG
	kgSub, err := kg.NewKGSubscriber("tcp://localhost:1883", "mindset-agent-kg", kgInstance)
	if err != nil {
		log.Printf("[KG] ⚠️ Warning: failed to create subscriber: %v", err)
	} else {
		if err := kgSub.Start(); err != nil {
			log.Printf("[KG] ⚠️ Warning: failed to start subscriber: %v", err)
		} else {
			log.Printf("[KG] ✅ Subscriber started")
		}
	}

	return kgInstance
}

// initFunctionsRegistry initialise le registre des fonctions
func initFunctionsRegistry(cfg *config.Config, mqttPub *mqtt.Publisher, kgInstance *kg.KnowledgeGraph) *functions.Registry {
	fmt.Println("\n[FUNCTIONS] Registering functions...")

	funcRegistry := functions.NewRegistry()

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

	// Afficher les fonctions enregistrées
	fmt.Printf("\n[FUNCTIONS] Registered %d functions:\n", len(funcRegistry.List()))
	for _, fn := range funcRegistry.List() {
		fmt.Printf("  ✓ %s (%s) - %s\n", fn.Name, fn.Type, fn.Description)
	}

	return funcRegistry
}

// initPipelines initialise le moteur de pipelines
func initPipelines(funcRegistry *functions.Registry) *pipeline.Engine {
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

	pipelineEngine.Start()
	return pipelineEngine
}

// initOPCUA initialise la découverte OPC-UA
func initOPCUA(cfg *config.Config, mqttPub *mqtt.Publisher) *discovery.OPCUADiscovery {
	fmt.Println("\n[DISCOVERY] Starting OPC-UA auto-discovery...")
	fmt.Printf("[DISCOVERY] Connecting to: %s\n\n", cfg.OpcUA.Endpoint)

	opcua := discovery.NewOPCUADiscovery(cfg.OpcUA.Endpoint, mqttPub)

	ctx := context.Background()
	if err := opcua.Connect(ctx); err != nil {
		log.Fatalf("[DISCOVERY] Connection failed: %v\n\nMake sure Prosys OPC-UA Simulator is running!", err)
	}

	return opcua
}

// browseTags browse l'arbre des tags
func browseTags(ctx context.Context, opcua *discovery.OPCUADiscovery) []discovery.Tag {
	fmt.Println("\n[DISCOVERY] Browsing node tree...\n")
	tags, err := opcua.BrowseNodeTree(ctx)
	if err != nil {
		log.Fatalf("[DISCOVERY] Browse failed: %v", err)
	}

	fmt.Printf("\n[DISCOVERY] ✅ Found %d tags total\n", len(tags))
	return tags
}

// printSummary affiche le résumé des tags
func printSummary(tags []discovery.Tag) {
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
}

// startSubscription démarre la subscription live
func startSubscription(ctx context.Context, opcua *discovery.OPCUADiscovery, tags []discovery.Tag) {
	fmt.Println("\n[SUBSCRIBE] Starting live data subscription...")
	fmt.Println("[SUBSCRIBE] Watching for value changes:\n")

	err := opcua.Subscribe(ctx, tags, 500*time.Millisecond, func(tag discovery.Tag) {
		fmt.Printf("  📡 CHANGE | %-40s | %v (%s)\n",
			tag.Name,
			tag.Value,
			tag.DataType,
		)
	})
	if err != nil {
		log.Printf("[SUBSCRIBE] Warning: subscription failed: %v", err)
	}
}

// startWatcher démarre le watcher de changements
func startWatcher(ctx context.Context, opcua *discovery.OPCUADiscovery, tags []discovery.Tag) {
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
}
