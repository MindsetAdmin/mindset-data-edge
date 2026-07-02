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
	"github.com/MindsetAdmin/mindset-data-edge/internal/mqtt"
	"github.com/MindsetAdmin/mindset-data-edge/internal/uns"
)

// TagSelection is one row of the user's data-flow choice: a discovered node and
// how its data should be routed.
type TagSelection struct {
	NodeID string `json:"node_id"`
	Mode   string `json:"mode"` // "raw" | "isa95" | "both"
}

// discoveredTag is the API shape returned by /api/opcua/discover.
type discoveredTag struct {
	NodeID   string      `json:"node_id"`
	Name     string      `json:"name"`
	DataType string      `json:"data_type"`
	Value    interface{} `json:"value"`
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

	disc       *discovery.OPCUADiscovery
	pub        *mqtt.Publisher
	status     string // "disconnected" | "connecting" | "connected" | "error"
	endpoint   string
	lastErr    string
	tags       []discovery.Tag   // last Discover() result
	selections map[string]string // nodeID -> mode
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewOPCUAManager builds a manager bound to the given MQTT broker and site.
// It does not connect to anything until Connect is called.
func NewOPCUAManager(broker string, cfg *config.Config) *OPCUAManager {
	site := "local-test"
	if cfg != nil && cfg.Site.ID != "" {
		site = cfg.Site.ID
	}
	return &OPCUAManager{
		broker:     broker,
		siteID:     site,
		mapper:     uns.NewMapper(site),
		status:     "disconnected",
		selections: map[string]string{},
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

	out := make([]discoveredTag, 0, len(tags))
	for _, t := range tags {
		out = append(out, discoveredTag{NodeID: t.NodeID, Name: t.Name, DataType: t.DataType, Value: t.Value})
	}
	log.Printf("[OPCUA-MGR] Discovered %d tags", len(out))
	return out, nil
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
	var selected []discovery.Tag
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
	}

	if len(selected) == 0 {
		return 0, fmt.Errorf("no valid tag selections (discover tags first, then select with mode raw|isa95|both)")
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
// mode allows it.
func (m *OPCUAManager) route(tag discovery.Tag) {
	m.mu.RLock()
	mode := m.selections[tag.NodeID]
	pub := m.pub
	m.mu.RUnlock()

	if pub == nil || (mode != "isa95" && mode != "both") {
		return
	}

	node := m.mapper.MapTag(tag.Name, tag.DataType)
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
			WorkCenter:     node.WorkCenter,
			WorkUnit:       node.WorkUnit,
			TagName:        node.TagName,
			UNSTopic:       node.FullTopic(),
		},
	}
	pub.PublishJSON(node.FullTopic(), msg)
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
		node := m.mapper.MapTag(name, "")
		out = append(out, SelectionDetail{
			NodeID:     id,
			Name:       name,
			Mode:       mode,
			WorkCenter: node.WorkCenter,
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
