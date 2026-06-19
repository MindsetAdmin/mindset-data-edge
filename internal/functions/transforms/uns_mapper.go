package transforms

import (
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/uns"
)

// UNSMapperConfig configuration
type UNSMapperConfig struct {
	SiteID string `json:"site_id"`
	Area   string `json:"area"`
}

// UNSMapperResult résultat
type UNSMapperResult struct {
	Topic     string      `json:"topic"`
	FullTopic string      `json:"full_topic"`
	Node      uns.UNSNode `json:"node"`
}

// UNSMapperHandler handler
type UNSMapperHandler struct {
	mapper *uns.Mapper
}

// NewUNSMapperHandler crée un nouveau handler
func NewUNSMapperHandler(siteID string) *UNSMapperHandler {
	return &UNSMapperHandler{
		mapper: uns.NewMapper(siteID),
	}
}

// GetFunction retourne la définition
func (h *UNSMapperHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "uns_mapper",
		Type:        functions.TypeTransform,
		Description: "Transforme un tag OPC-UA en nœud UNS ISA-95",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			tagName, _ := params["tag_name"].(string)
			dataType, _ := params["data_type"].(string)
			return h.Execute(tagName, dataType)
		},
	}
}

// Execute exécute le mapping
func (h *UNSMapperHandler) Execute(tagName, dataType string) (*UNSMapperResult, error) {
	node := h.mapper.MapTag(tagName, dataType)

	return &UNSMapperResult{
		Topic:     node.Topic(),
		FullTopic: node.FullTopic(),
		Node:      node,
	}, nil
}
