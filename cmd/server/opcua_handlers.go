package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MindsetAdmin/mindset-data-edge/internal/discovery"
)

// opcuaConnectRequest is the body of POST /api/opcua/connect.
type opcuaConnectRequest struct {
	Endpoint       string `json:"endpoint"`
	SecurityMode   string `json:"security_mode"`
	SecurityPolicy string `json:"security_policy"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	SessionTimeout int    `json:"session_timeout"`
}

// handleOpcuaConnect connects to a user-specified OPC-UA server.
func (s *server) handleOpcuaConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req opcuaConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Endpoint) == "" {
		http.Error(w, "endpoint is required", http.StatusBadRequest)
		return
	}

	err := s.opcua.Connect(discovery.ConnectionConfig{
		Endpoint:          req.Endpoint,
		SecurityMode:      req.SecurityMode,
		SecurityPolicy:    req.SecurityPolicy,
		Username:          req.Username,
		Password:          req.Password,
		SessionTimeoutSec: req.SessionTimeout,
	})
	if err != nil {
		writeJSONStatus(w, http.StatusBadGateway, map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	writeJSON(w, s.opcua.Status())
}

// handleOpcuaDiscover browses the connected server and returns its tags.
func (s *server) handleOpcuaDiscover(w http.ResponseWriter, r *http.Request) {
	tags, err := s.opcua.Discover()
	if err != nil {
		writeJSONStatus(w, http.StatusBadGateway, map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	writeJSON(w, map[string]interface{}{"tags": tags, "total": len(tags)})
}

// opcuaSubscribeRequest is the body of POST /api/opcua/subscribe.
type opcuaSubscribeRequest struct {
	Selections []TagSelection `json:"selections"`
}

// handleOpcuaSubscribe starts monitoring the selected tags with their modes.
func (s *server) handleOpcuaSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req opcuaSubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	n, err := s.opcua.Subscribe(req.Selections)
	if err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	writeJSON(w, map[string]interface{}{"status": "subscribed", "count": n})
}

// handleOpcuaDisconnect closes the current OPC-UA session.
func (s *server) handleOpcuaDisconnect(w http.ResponseWriter, r *http.Request) {
	s.opcua.Disconnect()
	writeJSON(w, s.opcua.Status())
}

// handleOpcuaStatus returns the current connection status.
func (s *server) handleOpcuaStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.opcua.Status())
}

// handleOpcuaSelections returns the current per-tag routing with ISA-95 mapping.
// The builder uses this to restrict function field pickers to isa95/both tags.
func (s *server) handleOpcuaSelections(w http.ResponseWriter, r *http.Request) {
	list := s.opcua.SelectionsDetailed()
	writeJSON(w, map[string]interface{}{"selections": list, "total": len(list)})
}

// writeJSONStatus writes a JSON body with an explicit HTTP status code.
func writeJSONStatus(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
