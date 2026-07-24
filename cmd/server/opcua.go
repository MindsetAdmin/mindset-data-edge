package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/config"
	"github.com/MindsetAdmin/mindset-data-edge/internal/discovery"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
	"github.com/MindsetAdmin/mindset-data-edge/internal/mqtt"
	"github.com/MindsetAdmin/mindset-data-edge/internal/uns"
)

// TagSelection is one row of the user's data-flow choice: a discovered node,
// how its data should be routed, and — new in Entry 124 — an optional
// correction to its auto-computed ISA-95 mapping. Any override field left
// empty falls back to the mapper's own guess for that field; this isn't
// all-or-nothing per tag.
type TagSelection struct {
	NodeID     string `json:"node_id"`
	Mode       string `json:"mode"` // "raw" | "isa95" | "both"
	Area       string `json:"area,omitempty"`
	WorkCenter string `json:"work_center,omitempty"`
	WorkUnit   string `json:"work_unit,omitempty"`
	TagName    string `json:"tag_name,omitempty"`
}

// discoveredTag is the API shape returned by /api/opcua/discover. Entry 124
// added the ISA-95 preview + confidence fields — previously the mapping was
// computed (for the KG-seed side effect) but never reached the frontend, so
// there was nothing to review or correct before a tag was routed.
type discoveredTag struct {
	NodeID     string      `json:"node_id"`
	Name       string      `json:"name"`
	DataType   string      `json:"data_type"`
	Value      interface{} `json:"value"`
	Site       string      `json:"site"`
	Area       string      `json:"area"`
	WorkCenter string      `json:"work_center"`
	WorkUnit   string      `json:"work_unit,omitempty"`
	TagName    string      `json:"tag_name"`
	Confidence float64     `json:"confidence"`
	Pending    bool        `json:"pending"` // confidence < kg.AutoAcceptThreshold — same gate the KG bootstrap uses
}

// OPCUAManager owns the single live OPC-UA session driven from the frontend.
// It connects on demand, browses tags, and — for the tags the user selects —
// publishes raw values (always) and ISA-95 contextualized values (isa95/both).
// Governance is enforced here at the source: raw-only tags never reach
// mindset/site/#, so functions (which consume site/#) physically cannot use them.
type OPCUAManager struct {
	mu     sync.RWMutex
	broker string
	siteID string
	mapper *uns.Mapper
	kg     *kg.KnowledgeGraph // v0 structural bootstrap target — Entry 95/96, may be nil in tests

	disc       *discovery.OPCUADiscovery
	pub        *mqtt.Publisher
	status     string // "disconnected" | "connecting" | "connected" | "error"
	endpoint   string
	lastErr    string
	tags       []discovery.Tag   // last Discover() result
	selections map[string]string // nodeID -> mode
	// overrides holds a user-corrected ISA-95 mapping per node (Entry 124) —
	// checked first in route() on every live value, ahead of the mapper's own
	// guess, so a correction actually changes where a tag's data is published,
	// not just what gets written into the KG once.
	overrides map[string]uns.UNSNode
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewOPCUAManager builds a manager bound to the given MQTT broker and site.
// It does not connect to anything until Connect is called. kgInst is the v0
// structural-bootstrap target (Entry 95/96) — pass nil to disable seeding
// (e.g. in tests) without touching the rest of the manager.
func NewOPCUAManager(broker string, cfg *config.Config, kgInst *kg.KnowledgeGraph) *OPCUAManager {
	site := "local-test"
	if cfg != nil && cfg.Site.ID != "" {
		site = cfg.Site.ID
	}
	return &OPCUAManager{
		broker:     broker,
		siteID:     site,
		mapper:     uns.NewMapper(site),
		kg:         kgInst,
		status:     "disconnected",
		selections: map[string]string{},
		overrides:  map[string]uns.UNSNode{},
	}
}

// Connect opens (or replaces) the OPC-UA session using the provided config.
func (m *OPCUAManager) Connect(cfg discovery.ConnectionConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Tear down any existing session first — one session at a time.
	m.disconnectLocked()

	// Lazily create the MQTT publisher used for raw + site routing.
	if m.pub == nil {
		pub, err := mqtt.NewPublisher(m.broker, m.siteID)
		if err != nil {
			m.status = "error"
			m.lastErr = fmt.Sprintf("mqtt publisher: %v", err)
			return fmt.Errorf("connect to MQTT broker %s: %w", m.broker, err)
		}
		m.pub = pub
	}

	m.status = "connecting"
	m.endpoint = cfg.Endpoint
	m.lastErr = ""

	// Use a long-lived context tied to the session: gopcua keeps the connection
	// alive for the lifetime of this context, so it must outlive Connect().
	m.ctx, m.cancel = context.WithCancel(context.Background())

	disc := discovery.NewOPCUADiscoveryWithConfig(cfg, m.pub)
	if err := disc.Connect(m.ctx); err != nil {
		m.status = "error"
		m.lastErr = err.Error()
		m.cancel()
		m.cancel = nil
		return err
	}

	m.disc = disc
	m.status = "connected"
	log.Printf("[OPCUA-MGR] Connected to %s", cfg.Endpoint)
	return nil
}

// Discover browses the connected server's node tree and caches the result.
// Entry 124: the ISA-95 mapping + confidence, previously computed only for
// the KG-seed side effect and never returned, is now attached to every tag
// so the UI can show it — and let the user correct it — before anything is
// subscribed.
func (m *OPCUAManager) Discover() ([]discoveredTag, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.disc == nil || m.status != "connected" {
		return nil, fmt.Errorf("not connected to an OPC-UA server")
	}

	tags, err := m.disc.BrowseNodeTree(m.ctx)
	if err != nil {
		m.lastErr = err.Error()
		return nil, err
	}
	m.tags = tags

	entries := m.computeMappings(tags)
	entryByID := make(map[string]kg.HierarchyEntry, len(entries))
	for _, e := range entries {
		entryByID[e.NodeID] = e
	}

	out := make([]discoveredTag, 0, len(tags))
	for _, t := range tags {
		e := entryByID[t.NodeID]
		out = append(out, discoveredTag{
			NodeID: t.NodeID, Name: t.Name, DataType: t.DataType, Value: t.Value,
			Site: e.Site, Area: e.Area, WorkCenter: e.WorkCenter, WorkUnit: e.WorkUnit, TagName: e.TagName,
			Confidence: e.Confidence, Pending: e.Confidence < kg.AutoAcceptThreshold,
		})
	}
	log.Printf("[OPCUA-MGR] Discovered %d tags", len(out))

	m.seedKG(entries)

	return out, nil
}

// computeMappings ISA-95-maps every discovered tag and scores it with the
// same confidence heuristic seedKG has always persisted (Entry 107) — factored
// out so Discover() can return it for display/editing (Entry 124) without
// duplicating the heuristic logic.
//
// Confidence signals computed across the whole discovered batch, not per tag
// in isolation — a naming-convention heuristic is only as trustworthy as it
// is consistent. Two cheap, explainable signals (not ML):
//  1. Does this tag's dot-depth match the server's modal (most common)
//     depth? A one-off depth among a batch of otherwise-uniform tags is a
//     sign the naming convention assumption doesn't hold for it.
//  2. Does its normalized tag name collide with another tag already mapped
//     to the same equipment? A collision means the mapper folded two
//     distinct signals onto one name — ambiguous, not trustworthy.
func (m *OPCUAManager) computeMappings(tags []discovery.Tag) []kg.HierarchyEntry {
	depthCounts := map[int]int{}
	for _, t := range tags {
		depthCounts[len(strings.Split(t.Name, "."))]++
	}
	modalDepth, modalCount := 0, 0
	for d, c := range depthCounts {
		if c > modalCount {
			modalDepth, modalCount = d, c
		}
	}
	seenPerEquipment := map[string]map[string]bool{} // equipment identity -> tag names already claimed

	entries := make([]kg.HierarchyEntry, 0, len(tags))
	for _, t := range tags {
		node := m.mapper.MapTag(t.Name, t.DataType)
		depth := len(strings.Split(t.Name, "."))

		equipmentKey := node.EquipmentIdentity() // matches kg.SeedFromDiscovery's own depth branching
		if seenPerEquipment[equipmentKey] == nil {
			seenPerEquipment[equipmentKey] = map[string]bool{}
		}
		collision := seenPerEquipment[equipmentKey][node.TagName]
		seenPerEquipment[equipmentKey][node.TagName] = true

		confidence := 1.0
		if depth != modalDepth {
			confidence -= 0.5
		}
		if collision {
			confidence -= 0.5
		}

		entries = append(entries, kg.HierarchyEntry{
			Site:       node.Site,
			Area:       node.Area,
			WorkCenter: node.WorkCenter,
			WorkUnit:   node.WorkUnit,
			Depth:      depth, // WorkCenter/WorkUnit mean different things at different depths — see kg.SeedFromDiscovery
			Confidence: confidence,
			TagName:    node.TagName,
			NodeID:     t.NodeID,
			DataType:   t.DataType,
		})
	}
	return entries
}

// seedKG writes ISA-95-mapped structural entries into the Knowledge Graph
// (Entry 95/96 — v0, OT only). Entries below kg.AutoAcceptThreshold are
// flagged pending — a human validates the uncertain ones, not everything.
// Best-effort: a seeding failure is logged, not returned, so a KG hiccup
// never breaks the discover/subscribe flow the UI depends on.
func (m *OPCUAManager) seedKG(entries []kg.HierarchyEntry) {
	if m.kg == nil || len(entries) == 0 {
		return
	}
	seeded, err := m.kg.SeedFromDiscovery(entries)
	if err != nil {
		log.Printf("[OPCUA-MGR] KG structural seed failed: %v", err)
		return
	}
	log.Printf("[OPCUA-MGR] KG structural seed: %d work centers seeded/confirmed, pending validation", seeded)
}

// Subscribe starts monitoring the selected tags with their routing modes.
// discovery.Subscribe publishes raw for every monitored tag; the route callback
// adds ISA-95 publication for isa95/both.
//
// NOTE (v1): re-applying selections stacks a new subscription on top of the old
// one. Apply once per session; to change the selection, Disconnect then Connect.
func (m *OPCUAManager) Subscribe(sel []TagSelection) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.disc == nil || m.status != "connected" {
		return 0, fmt.Errorf("not connected to an OPC-UA server")
	}

	byID := make(map[string]discovery.Tag, len(m.tags))
	for _, t := range m.tags {
		byID[t.NodeID] = t
	}

	m.selections = map[string]string{}
	m.overrides = map[string]uns.UNSNode{}
	var selected []discovery.Tag
	var correctedEntries []kg.HierarchyEntry
	for _, s := range sel {
		mode := normalizeMode(s.Mode)
		if mode == "" {
			continue
		}
		t, ok := byID[s.NodeID]
		if !ok {
			continue // unknown node id — skip rather than fail the whole batch
		}
		m.selections[s.NodeID] = mode
		selected = append(selected, t)

		// Entry 124: a user-supplied ISA-95 correction. Any field left blank
		// falls back to the mapper's own guess — not all-or-nothing per tag.
		if s.Area != "" || s.WorkCenter != "" || s.WorkUnit != "" || s.TagName != "" {
			node := m.mapper.MapTag(t.Name, t.DataType)
			if s.Area != "" {
				node.Area = s.Area
			}
			if s.WorkCenter != "" {
				node.WorkCenter = s.WorkCenter
			}
			if s.WorkUnit != "" {
				node.WorkUnit = s.WorkUnit
			}
			if s.TagName != "" {
				node.TagName = s.TagName
			}
			// Unit was inferred from the auto tag name — stale if the user
			// renamed it, and there's no way to re-infer a corrected one
			// reliably, so drop it rather than publish a wrong unit.
			if s.TagName != "" {
				node.Unit = ""
			}
			m.overrides[s.NodeID] = node

			if mode == "isa95" || mode == "both" {
				correctedEntries = append(correctedEntries, kg.HierarchyEntry{
					Site: node.Site, Area: node.Area, WorkCenter: node.WorkCenter, WorkUnit: node.WorkUnit,
					// A human just supplied this — treat it as confirmed, not
					// another guess needing its own review.
					Confidence: 1.0,
					TagName:    node.TagName, NodeID: t.NodeID, DataType: t.DataType,
					Depth: len(strings.Split(t.Name, ".")),
				})
			}
		}
	}

	if len(selected) == 0 {
		return 0, fmt.Errorf("no valid tag selections (discover tags first, then select with mode raw|isa95|both)")
	}

	// Write corrections before subscribing — additive (idempotent INSERT OR
	// IGNORE), doesn't retract whatever the original auto-guess already wrote;
	// if that guess is still sitting pending, reject it from the KG page same
	// as any other bad auto-guess (known v0 limitation, same as elsewhere in
	// the structural bootstrap).
	if len(correctedEntries) > 0 {
		m.seedKG(correctedEntries)
	}

	if err := m.disc.Subscribe(m.ctx, selected, 500*time.Millisecond, m.route); err != nil {
		m.lastErr = err.Error()
		return 0, err
	}
	log.Printf("[OPCUA-MGR] Subscribed %d tags (raw always; site for isa95/both)", len(selected))
	return len(selected), nil
}

// route runs on every monitored value change. Raw is already published by
// discovery.Subscribe; here we add the ISA-95 site publication when the tag's
// mode allows it. Entry 124: a user correction from Subscribe (m.overrides)
// takes priority over the mapper's own guess — otherwise a corrected mapping
// would only ever affect the KG, not what topic the live data actually lands on.
func (m *OPCUAManager) route(tag discovery.Tag) {
	m.mu.RLock()
	mode := m.selections[tag.NodeID]
	pub := m.pub
	override, hasOverride := m.overrides[tag.NodeID]
	m.mu.RUnlock()

	if pub == nil || (mode != "isa95" && mode != "both") {
		return
	}

	node := override
	if !hasOverride {
		node = m.mapper.MapTag(tag.Name, tag.DataType)
	}
	msg := uns.ContextualizedMessage{
		TimestampMs: time.Now().UnixMilli(),
		Value:       tag.Value,
		Unit:        node.Unit,
		DataType:    tag.DataType,
		Metadata: uns.Metadata{
			SourceProtocol: "OPC-UA",
			OriginalNodeID: tag.NodeID,
			OriginalName:   tag.Name,
			SiteID:         node.Site,
			Area:           node.Area,
			// EquipmentIdentity, not the raw WorkCenter field (Entry 127):
			// at 4+ level tag names, node.WorkCenter is a grouping level (a
			// line) ABOVE the actual machine, and node.WorkUnit is the real
			// machine — using WorkCenter directly here silently merged every
			// machine on the same line into one shared StateTracker entry.
			WorkCenter: node.EquipmentIdentity(),
			WorkUnit:   node.WorkUnit,
			TagName:    node.TagName,
			UNSTopic:   node.FullTopic(),
		},
	}
	pub.PublishJSONRetained(node.FullTopic(), msg)
}

// Disconnect closes the current OPC-UA session, if any.
func (m *OPCUAManager) Disconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.disconnectLocked()
	log.Printf("[OPCUA-MGR] Disconnected")
}

// disconnectLocked tears down the session. Caller must hold m.mu.
func (m *OPCUAManager) disconnectLocked() {
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
	if m.disc != nil {
		m.disc.Disconnect(context.Background())
		m.disc = nil
	}
	m.tags = nil
	m.selections = map[string]string{}
	m.overrides = map[string]uns.UNSNode{}
	m.status = "disconnected"
	m.endpoint = ""
	m.lastErr = ""
}

// Status is the snapshot behind GET /api/opcua/status.
func (m *OPCUAManager) Status() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return map[string]interface{}{
		"status":    m.status,
		"endpoint":  m.endpoint,
		"connected": m.status == "connected",
		"tag_count": len(m.tags),
		"selected":  len(m.selections),
		"error":     m.lastErr,
	}
}

// Selections returns the current per-tag routing modes.
func (m *OPCUAManager) Selections() []TagSelection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]TagSelection, 0, len(m.selections))
	for id, mode := range m.selections {
		out = append(out, TagSelection{NodeID: id, Mode: mode})
	}
	return out
}

// SelectionDetail enriches a selection with its ISA-95 mapping so the builder
// can offer only function-eligible (isa95/both) work centers and site topics.
type SelectionDetail struct {
	NodeID     string `json:"node_id"`
	Name       string `json:"name"`
	Mode       string `json:"mode"`
	WorkCenter string `json:"work_center"`
	TagName    string `json:"tag_name"`
	UNSTopic   string `json:"uns_topic"`
}

// SelectionsDetailed maps each current selection through the UNS mapper.
func (m *OPCUAManager) SelectionsDetailed() []SelectionDetail {
	m.mu.RLock()
	defer m.mu.RUnlock()
	nameByID := make(map[string]string, len(m.tags))
	for _, t := range m.tags {
		nameByID[t.NodeID] = t.Name
	}
	out := make([]SelectionDetail, 0, len(m.selections))
	for id, mode := range m.selections {
		name := nameByID[id]
		// Same precedence as route(): a user correction (Entry 124) wins over
		// the mapper's own guess, so this reflects what's actually being
		// published, not a stale recomputation from scratch.
		node, hasOverride := m.overrides[id]
		if !hasOverride {
			node = m.mapper.MapTag(name, "")
		}
		out = append(out, SelectionDetail{
			NodeID: id,
			Name:   name,
			Mode:   mode,
			// EquipmentIdentity, not the raw WorkCenter field (Entry 127) —
			// see route()'s identical fix for why.
			WorkCenter: node.EquipmentIdentity(),
			TagName:    node.TagName,
			UNSTopic:   node.FullTopic(),
		})
	}
	return out
}

// normalizeMode canonicalizes the data-flow mode string.
func normalizeMode(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "raw":
		return "raw"
	case "isa95", "isa-95", "normalized", "site":
		return "isa95"
	case "both":
		return "both"
	default:
		return ""
	}
}
