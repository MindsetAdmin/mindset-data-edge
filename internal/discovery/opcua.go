// internal/discovery/opcua.go
package discovery

import (
    "context"
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/gopcua/opcua"
    "github.com/gopcua/opcua/ua"
    "github.com/MindsetAdmin/mindset-data-edge/internal/mqtt"
)

// Tag represents a discovered OPC-UA tag
type Tag struct {
    NodeID   string
    Name     string
    DataType string
    Value    interface{}
}

// TagChange represents a diff between two tag snapshots
type TagChange struct {
    Added   []Tag
    Removed []Tag
}

// ConnectionConfig holds everything needed to open an OPC-UA session. It lets the
// caller (e.g. the dynamic, frontend-driven flow in cmd/server) override what used
// to be hardcoded here. Empty fields fall back to safe defaults (None security,
// anonymous auth, 60s session timeout).
type ConnectionConfig struct {
    Endpoint          string
    SecurityMode      string // "None" | "Sign" | "SignAndEncrypt"
    SecurityPolicy    string // "None" | "Basic256Sha256" | "Basic256" | "Basic128Rsa15"
    Username          string
    Password          string
    SessionTimeoutSec int // default 60
}

// OPCUADiscovery handles OPC-UA connection and node tree browsing
type OPCUADiscovery struct {
    endpoint string
    connCfg  ConnectionConfig
    client   *opcua.Client
    mqttPub   *mqtt.Publisher
}

// NewOPCUADiscovery creates a discovery instance with the legacy defaults
// (no security, anonymous). Kept so cmd/agent's existing call site is untouched.
func NewOPCUADiscovery(endpoint string, mqttPub *mqtt.Publisher) *OPCUADiscovery {
    return NewOPCUADiscoveryWithConfig(ConnectionConfig{
        Endpoint:       endpoint,
        SecurityMode:   "None",
        SecurityPolicy: "None",
    }, mqttPub)
}

// NewOPCUADiscoveryWithConfig creates a discovery instance from a full
// ConnectionConfig, applying defaults for any empty field.
func NewOPCUADiscoveryWithConfig(cfg ConnectionConfig, mqttPub *mqtt.Publisher) *OPCUADiscovery {
    if cfg.SecurityMode == "" {
        cfg.SecurityMode = "None"
    }
    if cfg.SecurityPolicy == "" {
        cfg.SecurityPolicy = "None"
    }
    if cfg.SessionTimeoutSec <= 0 {
        cfg.SessionTimeoutSec = 300
    }
    return &OPCUADiscovery{
        endpoint: cfg.Endpoint,
        connCfg:  cfg,
        mqttPub:  mqttPub,
    }
}

// parseSecurityMode maps the config string onto the gopcua enum.
func parseSecurityMode(s string) ua.MessageSecurityMode {
    switch strings.ToLower(strings.TrimSpace(s)) {
    case "sign":
        return ua.MessageSecurityModeSign
    case "signandencrypt", "sign_and_encrypt", "sign-and-encrypt":
        return ua.MessageSecurityModeSignAndEncrypt
    default:
        return ua.MessageSecurityModeNone
    }
}

// Connect establishes connection to the OPC-UA server using the ConnectionConfig.
func (d *OPCUADiscovery) Connect(ctx context.Context) error {
    secMode := parseSecurityMode(d.connCfg.SecurityMode)
    secPolicy := d.connCfg.SecurityPolicy
    if secPolicy == "" {
        secPolicy = "None"
    }
    log.Printf("[OPC-UA] Connecting to %s (security=%s/%s)", d.endpoint, d.connCfg.SecurityMode, secPolicy)

    // Sign / SignAndEncrypt require a client certificate + key pair, which isn't
    // wired yet. Fail fast with a clear message rather than a cryptic handshake error.
    if secMode != ua.MessageSecurityModeNone {
        return fmt.Errorf("security mode %q requires a client certificate, which is not yet configured — use \"None\" for now", d.connCfg.SecurityMode)
    }

    opts := []opcua.Option{
        opcua.SecurityMode(secMode),
        opcua.SecurityPolicy(secPolicy),
        opcua.ApplicationName("MindsetData"),
        opcua.RequestTimeout(30 * time.Second),
        opcua.SessionTimeout(time.Duration(d.connCfg.SessionTimeoutSec) * time.Second),
    }

    // Username/password auth when provided, otherwise anonymous.
    if d.connCfg.Username != "" {
        opts = append(opts, opcua.AuthUsername(d.connCfg.Username, d.connCfg.Password))
    } else {
        opts = append(opts, opcua.AuthAnonymous())
    }

    client, err := opcua.NewClient(d.endpoint, opts...)
    if err != nil {
        return fmt.Errorf("failed to create client: %w", err)
    }

    if err := client.Connect(ctx); err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }

    d.client = client
    log.Printf("[OPC-UA] Connected successfully to %s", d.endpoint)
    return nil
}


// BrowseNodeTree recursively browses the OPC-UA node tree
// starting from the Objects folder
func (d *OPCUADiscovery) BrowseNodeTree(ctx context.Context) ([]Tag, error) {
    log.Printf("[OPC-UA] Starting node tree browse...")

    // Start from Objects folder (standard OPC-UA entry point)
    objectsNodeID := ua.NewNumericNodeID(0, 85)

    var tags []Tag
    err := d.browseNode(ctx, objectsNodeID, "", &tags, 0)
    if err != nil {
        return nil, fmt.Errorf("browse failed: %w", err)
    }

    log.Printf("[OPC-UA] Browse complete — found %d nodes", len(tags))
    return tags, nil
}

// browseNode recursively browses a node and its children

func (d *OPCUADiscovery) browseNode(
    ctx context.Context,
    nodeID *ua.NodeID,
    parentName string,
    tags *[]Tag,
    depth int,
) error {
    if depth > 10 {
        return nil
    }

    // CRITICAL: Breathe to prevent flooding the server
    time.Sleep(50 * time.Millisecond)
    
    req := &ua.BrowseRequest{
        RequestedMaxReferencesPerNode: 100, // Lowered from 1000 to avoid message size limits
        NodesToBrowse: []*ua.BrowseDescription{
            {
                NodeID:          nodeID,
                BrowseDirection: ua.BrowseDirectionForward,
                ReferenceTypeID: ua.NewNumericNodeID(0, 31), // HierarchicalReferences
                IncludeSubtypes: true,
                NodeClassMask:   uint32(ua.NodeClassObject | ua.NodeClassVariable),
                ResultMask:      uint32(ua.BrowseResultMaskAll),
            },
        },
    }
    time.Sleep(1 *time.Second)
    resp, err := d.client.Browse(ctx, req)
    
    if err != nil {
        time.Sleep(100 * time.Millisecond)
        return fmt.Errorf("browse failed for %s: %w", parentName, err)
    }

    if resp.Results == nil || len(resp.Results) == 0 {
        return nil
    }

    result := resp.Results[0]
    if result.StatusCode != ua.StatusOK {
        return fmt.Errorf("browse returned status %v", result.StatusCode)
    }

    // Process the references we received
    for _, ref := range result.References {
        name := ref.BrowseName.Name

        // Skip internal server nodes that crash Prosys
        if name == "Server" || name == "ServerCapabilities" || name == "Aliases" || name == "Views" || name == "Types" || name == "Locations" || name == "MyObjects" || name == "StaticData"{
            continue
        }

        fullName := name
        if parentName != "" {
            fullName = parentName + "." + name
        }
    
        // Parse NodeID cleanly
        nodeID, err := ua.ParseNodeID(ref.NodeID.NodeID.String())
        if err != nil {
            continue
        }

        // Read value AND type BEFORE creating the tag
        _, dataType := d.readNodeValue(ctx, nodeID)
    
        if ref.NodeClass == ua.NodeClassVariable {
            tag := Tag{
                NodeID:   ref.NodeID.NodeID.String(),
                Name:     fullName,
                DataType: dataType,
            }
            *tags = append(*tags, tag)
            fmt.Printf("  ✓ [%s] %s   %s\n  ", tag.NodeID, tag.Name,tag.DataType)
        }

        if ref.NodeClass == ua.NodeClassObject {
            if ref.NodeID == nil || ref.NodeID.NodeID == nil {
                continue
            }
            time.Sleep(10 * time.Millisecond)
            // Parse the NodeID string to get a clean *ua.NodeID
            childNodeID, err := ua.ParseNodeID(ref.NodeID.NodeID.String())
            if err != nil {
                log.Printf("[OPC-UA] Warning: cannot parse NodeID for %s: %v", name, err)
                continue
            }

            if err := d.browseNode(ctx, childNodeID, fullName, tags, depth+1); err != nil {
                log.Printf("[OPC-UA] Warning: failed to browse child %s: %v", name, err)
            }
        }
    }

    // CRITICAL FIX: Release the Continuation Point!
    // If Prosys returns a continuation point, we MUST release it, otherwise it 
    // allocates memory, gets angry, and drops the connection (EOF).
    if len(result.ContinuationPoint) > 0 {
        _, _ = d.client.BrowseNext(ctx, &ua.BrowseNextRequest{
            ContinuationPoints:        [][]byte{result.ContinuationPoint},
            ReleaseContinuationPoints: true, // Tell server we are done, free the memory
        })
    }

    return nil
}
    
/*
// readNodeValue reads the current value of a node
func (d *OPCUADiscovery) readNodeValue(ctx context.Context, nodeID *ua.NodeID) (interface{}, string) {
    req := &ua.ReadRequest{
        NodesToRead: []*ua.ReadValueID{
            {
                NodeID:      nodeID,
                AttributeID: ua.AttributeIDValue,
            },
        },
    }

    resp, err := d.client.Read(ctx, req)
    if err != nil || resp.Results == nil || len(resp.Results) == 0 {
        return nil, "unknown"
    }

    result := resp.Results[0]
    if result.Status != ua.StatusOK {
        return nil, "unknown"
    }

    value := result.Value.Value()
    dataType := fmt.Sprintf("%T", value)

    return value, dataType
}
*/

// readNodeValue reads both the current value AND the OPC-UA DataType in one request
func (d *OPCUADiscovery) readNodeValue(ctx context.Context, nodeID *ua.NodeID) (interface{}, string) {
	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{
				NodeID:      nodeID,
				AttributeID: ua.AttributeIDValue, // ← the actual value
			},
			{
				NodeID:      nodeID,
				AttributeID: ua.AttributeIDDataType, // ← the OPC-UA type NodeID
			},
		},
	}

	resp, err := d.client.Read(ctx, req)
	if err != nil || resp.Results == nil || len(resp.Results) < 2 {
		return nil, "unknown"
	}

	// ── Result 0: the value ───────────────────────────────────────────────
	valueResult := resp.Results[0]
	if valueResult.Status != ua.StatusOK || valueResult.Value == nil {
		return nil, "unknown"
	}

	value := valueResult.Value.Value()
	if value == nil {
		return nil, "unknown"
	}

	// ── Result 1: the DataType NodeID ─────────────────────────────────────
	// The DataType attribute returns a NodeID that points to the type definition
	// Standard OPC-UA type NodeIDs live in namespace 0:
	//   0:1  = Boolean
	//   0:2  = SByte      0:3  = Byte
	//   0:4  = Int16      0:5  = UInt16
	//   0:6  = Int32      0:7  = UInt32
	//   0:8  = Int64      0:9  = UInt64
	//   0:10 = Float      0:11 = Double
	//   0:12 = String     0:13 = DateTime
	//   0:15 = ByteString
	dataType := "unknown"
	typeResult := resp.Results[1]
	if typeResult.Status == ua.StatusOK && typeResult.Value != nil {
		if typeNodeID, ok := typeResult.Value.Value().(*ua.NodeID); ok {
			dataType = opcuaTypeNodeIDToString(typeNodeID)
		}
	}

	// Fallback: if DataType attribute read failed, infer from Go type
	if dataType == "unknown" {
		dataType = mapDataType(value)
	}

	return value, dataType
}

// opcuaTypeNodeIDToString converts an OPC-UA DataType NodeID to a readable name
// These are the standard numeric NodeIDs for built-in types (namespace 0)
func opcuaTypeNodeIDToString(nodeID *ua.NodeID) string {
	if nodeID == nil {
		return "unknown"
	}

	// Only handle namespace 0 (standard OPC-UA built-in types)
	if nodeID.Namespace() != 0 {
		return fmt.Sprintf("CustomType(ns=%d;i=%d)", nodeID.Namespace(), nodeID.IntID())
	}

	switch nodeID.IntID() {
	case 1:
		return "Boolean"
	case 2:
		return "SByte"
	case 3:
		return "Byte"
	case 4:
		return "Int16"
	case 5:
		return "UInt16"
	case 6:
		return "Int32"
	case 7:
		return "UInt32"
	case 8:
		return "Int64"
	case 9:
		return "UInt64"
	case 10:
		return "Float"
	case 11:
		return "Double"
	case 12:
		return "String"
	case 13:
		return "DateTime"
	case 14:
		return "Guid"
	case 15:
		return "ByteString"
	case 21:
		return "Number"
	case 22:
		return "Integer"
	case 23:
		return "UInteger"
	default:
		return fmt.Sprintf("Type(i=%d)", nodeID.IntID())
	}
}

// mapDataType fallback — infers type from Go native type
func mapDataType(value interface{}) string {
	switch value.(type) {
	case bool:
		return "Boolean"
	case int8:
		return "SByte"
	case uint8:
		return "Byte"
	case int16:
		return "Int16"
	case uint16:
		return "UInt16"
	case int32:
		return "Int32"
	case uint32:
		return "UInt32"
	case int64:
		return "Int64"
	case uint64:
		return "UInt64"
	case float32:
		return "Float"
	case float64:
		return "Double"
	case string:
		return "String"
	case []byte:
		return "ByteString"
	default:
		return fmt.Sprintf("%T", value)
	}
}












// Disconnect cleanly closes the OPC-UA connection
func (d *OPCUADiscovery) Disconnect(ctx context.Context) {
    if d.client != nil {
        d.client.Close(ctx) // Fixed: updated API method
        log.Printf("[OPC-UA] Disconnected")
    }
}

// Subscribe monitors a list of tags and calls the callback on change
func (d *OPCUADiscovery) Subscribe(
    ctx context.Context,
    tags []Tag,
    interval time.Duration,
    callback func(tag Tag),
) error {
    log.Printf("[OPC-UA] Starting subscription for %d tags", len(tags))

    notifyCh := make(chan *opcua.PublishNotificationData)

    sub, err := d.client.Subscribe(ctx, &opcua.SubscriptionParameters{
        Interval: interval,
    }, notifyCh)
    if err != nil {
        return fmt.Errorf("failed to create subscription: %w", err)
    }

    // Create monitored items for each tag
    var miCreateRequests []*ua.MonitoredItemCreateRequest
    for i, tag := range tags {
        id, err := ua.ParseNodeID(tag.NodeID)
        if err != nil {
            log.Printf("[OPC-UA] Warning: invalid node ID %s: %v", tag.NodeID, err)
            continue
        }

        miCreateRequests = append(miCreateRequests, opcua.NewMonitoredItemCreateRequestWithDefaults(
            id,
            ua.AttributeIDValue,
            uint32(i+1), // ClientHandle must be unique
        ))
    }

    if _, err := sub.Monitor(ctx, ua.TimestampsToReturnBoth, miCreateRequests...); err != nil {
        return fmt.Errorf("failed to monitor items: %w", err)
    }

    // Listen for notifications (gopcua handles the publish loop automatically in v0.8+)
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case notify := <-notifyCh:
                if notify.Error != nil {
                    log.Printf("[OPC-UA] Notification error: %v", notify.Error)
                    continue
                }
                data, ok := notify.Value.(*ua.DataChangeNotification)
                if !ok {
                    continue
                }
                for _, item := range data.MonitoredItems {
                    // Find matching tag and call callback
                    idx := int(item.ClientHandle) - 1
                    if idx >= 0 && idx < len(tags) {
                        tags[idx].Value = item.Value.Value.Value()
                        d.mqttPub.PublishRaw(tags[idx].Name, tags[idx].NodeID, tags[idx].DataType, tags[idx].Value)
                        callback(tags[idx])
                    }
                }
            }
        }
    }()

    return nil
}

// WatchForChanges polls the node tree every interval and calls onChange
// when tags are added or removed. Only rebuilds on actual change.
func (d *OPCUADiscovery) WatchForChanges(
    ctx context.Context,
    currentTags []Tag,
    interval time.Duration,
    onChange func(change TagChange, allTags []Tag),
) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    // Build a map of current tags for fast lookup
    known := tagsToMap(currentTags)

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // Browse the node tree again — silent, no stdout
            freshTags, err := d.browseNodeTreeSilent(ctx)
            if err != nil {
                log.Printf("[WATCH] Browse failed: %v — skipping this cycle", err)
                continue
            }

            fresh := tagsToMap(freshTags)

            // Diff: what's new?
            var added []Tag
            for id, tag := range fresh {
                if _, exists := known[id]; !exists {
                    added = append(added, tag)
                    log.Printf("[WATCH] 🆕 New tag detected: %s (%s)", tag.Name, tag.DataType)
                }
            }

            // Diff: what's gone?
            var removed []Tag
            for id, tag := range known {
                if _, exists := fresh[id]; !exists {
                    removed = append(removed, tag)
                    log.Printf("[WATCH] 🗑  Tag removed: %s", tag.Name)
                }
            }

            // No change — do nothing
            if len(added) == 0 && len(removed) == 0 {
                log.Printf("[WATCH] No changes detected (%d tags stable)", len(freshTags))
                continue
            }

            // Update known map
            known = fresh

            // Notify caller with the diff and full fresh list
            onChange(TagChange{Added: added, Removed: removed}, freshTags)
        }
    }
}


// browseNodeTreeSilent browses without printing — used for periodic checks
func (d *OPCUADiscovery) browseNodeTreeSilent(ctx context.Context) ([]Tag, error) {
    objectsNodeID := ua.NewNumericNodeID(0, 85)
    var tags []Tag
    err := d.browseNodeSilent(ctx, objectsNodeID, "", &tags, 0)
    if err != nil {
        return nil, err
    }
    return tags, nil
}


// browseNodeSilent is identical to browseNode but without fmt.Printf
func (d *OPCUADiscovery) browseNodeSilent(
    ctx context.Context,
    nodeID *ua.NodeID,
    parentName string,
    tags *[]Tag,
    depth int,
) error {
    if depth > 10 {
        return nil
    }

    time.Sleep(50 * time.Millisecond)

    req := &ua.BrowseRequest{
        RequestedMaxReferencesPerNode: 100,
        NodesToBrowse: []*ua.BrowseDescription{
            {
                NodeID:          nodeID,
                BrowseDirection: ua.BrowseDirectionForward,
                ReferenceTypeID: ua.NewNumericNodeID(0, 33),
                IncludeSubtypes: true,
                NodeClassMask:   uint32(ua.NodeClassObject | ua.NodeClassVariable),
                ResultMask:      uint32(ua.BrowseResultMaskAll),
            },
        },
    }

    resp, err := d.client.Browse(ctx, req)
    if err != nil {
        time.Sleep(200 * time.Millisecond)
        resp, err = d.client.Browse(ctx, req)
        if err != nil {
            return fmt.Errorf("browse failed for %s: %w", parentName, err)
        }
    }

    if resp.Results == nil || len(resp.Results) == 0 {
        return nil
    }

    result := resp.Results[0]
    if result.StatusCode != ua.StatusOK {
        return nil
    }

    skipNodes := map[string]bool{
        "Server": true, "ServerCapabilities": true,
        "Aliases": true, "Views": true, "Types": true,
        "Locations": true, "MyObjects": true, "StaticData": true,
    }

    for _, ref := range result.References {
        name := ref.BrowseName.Name
        if skipNodes[name] {
            continue
        }

        fullName := name
        if parentName != "" {
            fullName = parentName + "." + name
        }

        if ref.NodeClass == ua.NodeClassVariable {
            childNodeID, err := ua.ParseNodeID(ref.NodeID.NodeID.String())
            if err != nil {
                continue
            }
            value, dataType := d.readNodeValue(ctx, childNodeID)
            *tags = append(*tags, Tag{
                NodeID:   ref.NodeID.NodeID.String(),
                Name:     fullName,
                DataType: dataType,
                Value:    value,
            })
        }

        if ref.NodeClass == ua.NodeClassObject {
            if ref.NodeID == nil {
                continue
            }
            childNodeID, err := ua.ParseNodeID(ref.NodeID.NodeID.String())
            if err != nil {
                continue
            }
            _ = d.browseNodeSilent(ctx, childNodeID, fullName, tags, depth+1)
        }
    }

    if len(result.ContinuationPoint) > 0 {
        _, _ = d.client.BrowseNext(ctx, &ua.BrowseNextRequest{
            ContinuationPoints:        [][]byte{result.ContinuationPoint},
            ReleaseContinuationPoints: true,
        })
    }

    return nil
}

// tagsToMap converts a tag slice to a map keyed by NodeID for fast diff
func tagsToMap(tags []Tag) map[string]Tag {
    m := make(map[string]Tag, len(tags))
    for _, t := range tags {
        m[t.NodeID] = t
    }
    return m
}