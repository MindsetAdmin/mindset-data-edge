// cmd/server/kg_handlers.go
// v0 structural bootstrap validation flow — docs/analysis_log.md Entries 95/96.
// Flat accept/reject over nodes SeedFromDiscovery wrote as pending:true.
package main

import "net/http"

// handleKGPending returns every business-category node still awaiting human
// validation after auto-generation from OPC-UA discovery.
func (s *server) handleKGPending(w http.ResponseWriter, r *http.Request) {
	nodes, err := s.kg.ListPending()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"nodes": nodes, "total": len(nodes)})
}

// handleKGValidate confirms an auto-generated node is correct.
func (s *server) handleKGValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if err := s.kg.ValidateNode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"status": "validated", "id": id})
}

// handleKGReject discards an auto-generated node that was wrong.
func (s *server) handleKGReject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if err := s.kg.RejectNode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"status": "rejected", "id": id})
}
