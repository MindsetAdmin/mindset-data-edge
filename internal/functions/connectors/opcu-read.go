package connectors

import (
	"context"
	"fmt"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

// OPCUAReadConfig configuration pour la lecture OPC-UA
type OPCUAReadConfig struct {
	NodeID  string        `json:"node_id"`
	Timeout time.Duration `json:"timeout"`
}

// OPCUAReadResult résultat de la lecture
type OPCUAReadResult struct {
	Value     interface{} `json:"value"`
	DataType  string      `json:"data_type"`
	Timestamp int64       `json:"timestamp_ms"`
	NodeID    string      `json:"node_id"`
}

// OPCUAReadHandler lit une valeur OPC-UA
type OPCUAReadHandler struct {
	client *opcua.Client
}

// NewOPCUAReadHandler crée un nouveau handler
func NewOPCUAReadHandler(client *opcua.Client) *OPCUAReadHandler {
	return &OPCUAReadHandler{client: client}
}

// GetFunction retourne la définition de la fonction
func (h *OPCUAReadHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "opcua_read",
		Type:        functions.TypeConnector,
		Description: "Lit une valeur depuis un nœud OPC-UA",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			return h.Execute(context.Background(), params)
		},
	}
}

// Execute exécute la lecture OPC-UA
func (h *OPCUAReadHandler) Execute(ctx context.Context, config map[string]interface{}) (*OPCUAReadResult, error) {
	nodeIDStr, ok := config["node_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing node_id in config")
	}

	nodeID, err := ua.ParseNodeID(nodeIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid node_id: %w", err)
	}

	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: nodeID, AttributeID: ua.AttributeIDValue},
		},
	}

	resp, err := h.client.Read(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Results == nil || len(resp.Results) == 0 {
		return nil, fmt.Errorf("no results")
	}

	result := resp.Results[0]
	if result.Status != ua.StatusOK {
		return nil, fmt.Errorf("read failed: %v", result.Status)
	}

	return &OPCUAReadResult{
		Value:     result.Value.Value(),
		DataType:  fmt.Sprintf("%T", result.Value.Value()),
		Timestamp: time.Now().UnixMilli(),
		NodeID:    nodeIDStr,
	}, nil
}
