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
	"github.com/MindsetAdmin/mindset-data-edge/internal/discovery"
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

	// ── STEP 1: OPC-UA Discovery ──────────────────────────────────────
	fmt.Println("\n[DISCOVERY] Starting OPC-UA auto-discovery...")
	fmt.Printf("[DISCOVERY] Connecting to: %s\n\n", cfg.OpcUA.Endpoint)

	opcua := discovery.NewOPCUADiscovery(cfg.OpcUA.Endpoint)

	if err := opcua.Connect(ctx); err != nil {
		log.Fatalf("[DISCOVERY] Connection failed: %v\n\nMake sure Prosys OPC-UA Simulator is running!", err)
	}
	defer opcua.Disconnect(ctx)

	// ── STEP 2: Browse node tree ──────────────────────────────────────
	fmt.Println("\n[DISCOVERY] Browsing node tree...\n")
	tags, err := opcua.BrowseNodeTree(ctx)
	if err != nil {
		log.Fatalf("[DISCOVERY] Browse failed: %v", err)
	}

	fmt.Printf("\n[DISCOVERY] ✅ Found %d tags total\n", len(tags))

	// ── STEP 3: Print summary ─────────────────────────────────────────
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




	// ── STEP 4: Live subscription (120 seconds demo) ───────────────────
	fmt.Println("\n[SUBSCRIBE] Starting live data subscription (10 seconds)...")
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

			// TODO Session 3: rebuild subscription with allTags
			// For now just log — subscription rebuild comes next
			fmt.Printf("[WATCH] Rebuild subscription with %d tags\n", len(allTags))
		},
	)

	// Wait 10 seconds then exit (or Ctrl+C)
	select {
	case <-ctx.Done():
		fmt.Println("\n[AGENT] Stopped by user.")
	case <-time.After(120 * time.Second):
		fmt.Println("\n[AGENT] 120-second demo complete.")
	}
}