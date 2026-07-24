// cmd/server/active_production.go
// Track B Phase 4 (docs/analysis_log.md Entries 115/116): answers "what
// product is running right now" by querying a human-validated work_order
// SchemaMapping. Deliberately does NOT answer "how long did product B take
// yesterday" — that needs retroactive event-tagging, out of scope here.
//
// Lives in cmd/server (not internal/kg or internal/connections) for the same
// reason OPCUAManager does: it needs both the KG (to find the validated
// mapping) and the connection registry (to query it), and neither package
// should depend on the other.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ActiveProductionFact is one row of "what's actively running" — one active
// work order on one machine, from one IT connection's validated mapping.
type ActiveProductionFact struct {
	ConnectionID string `json:"connection_id"`
	WorkCenter   string `json:"work_center"`
	OrderID      string `json:"order_id"`
	ProductID    string `json:"product_id"`
	Status       string `json:"status"`
	// EquipmentID is the matched OT Equipment node's id (e.g. "machine_Machine1"),
	// or "" if no OT node's work_center matched (Entry 120 entity resolution) —
	// left empty rather than guessed, so a caller can tell "linked" from "not."
	EquipmentID string `json:"equipment_id,omitempty"`
	// DueDate/CustomerID/DaysUntilDue (Entry 133) are the deadline-urgency
	// axis alongside kg_cost_summary's cost axis — "what should we fix
	// first" isn't only "what's expensive," it's also "what's due soon."
	// Both are optional: due_date/customer_id are bonus fields in the
	// work_order canonical mapping (internal/connections/canonical_suggest.go),
	// not every ERP exposes them, so an empty DueDate/nil DaysUntilDue means
	// "not available," never a guessed 0.
	DueDate      string `json:"due_date,omitempty"`
	CustomerID   string `json:"customer_id,omitempty"`
	DaysUntilDue *int   `json:"days_until_due,omitempty"`
}

// inProgressTokens — status values (case-insensitive substring match)
// considered "actively running" across ERPs with different vocabularies.
// Deliberately broad and hardcoded for v0, not configurable — a stated
// limitation of this phase, same as Entry 115 flagged in advance.
var inProgressTokens = []string{"running", "in progress", "in_progress", "released", "active", "started", "open"}

// validIdentifier — table/column names come from the canonical-mapping
// heuristic over information_schema, not directly from user input, but
// they're still spliced into a raw SQL string below. Cheap defensive check
// before doing that.
var validIdentifier = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// ActiveProduction queries every validated (pending:false) work_order
// SchemaMapping for the currently active order per work center, optionally
// filtered to one work center. A mapping missing a required field, or one
// whose table/column names fail the identifier check, is skipped rather than
// guessed — one bad mapping shouldn't silently fabricate an answer, and one
// bad connection shouldn't fail every other connection's query.
func (s *server) ActiveProduction(ctx context.Context, workCenter string) ([]ActiveProductionFact, error) {
	graph, err := s.kg.GetGraph("business")
	if err != nil {
		return nil, err
	}
	equipmentIdx := equipmentByWorkCenter(graph)

	var facts []ActiveProductionFact
	for _, n := range graph.Nodes {
		if n.Type != "SchemaMapping" {
			continue
		}
		if pending, _ := n.Properties["pending"].(bool); pending {
			continue // not yet human-validated
		}
		canonicalType, _ := n.Properties["canonical_type"].(string)
		if canonicalType != "work_order" {
			continue
		}

		connectionID, _ := n.Properties["connection_id"].(string)
		table, _ := n.Properties["table"].(string)
		fieldMap, _ := n.Properties["field_map"].(map[string]interface{})

		orderCol, _ := fieldMap["order_id"].(string)
		statusCol, _ := fieldMap["status"].(string)
		productCol, _ := fieldMap["product_id"].(string)
		wcCol, _ := fieldMap["work_center"].(string)
		if orderCol == "" || statusCol == "" || productCol == "" || wcCol == "" {
			continue // incomplete mapping — nothing reliable to query
		}
		// Bonus fields (Entry 133) — optional, unlike the core 4 above. Many
		// ERPs don't expose these (docs/impact_engine.md's pricing rule:
		// flag, don't fabricate), so an empty column name just means this
		// mapping won't surface deadline urgency, not that the query fails.
		dueDateCol, _ := fieldMap["due_date"].(string)
		customerCol, _ := fieldMap["customer_id"].(string)

		identCheck := []string{table, orderCol, statusCol, productCol, wcCol}
		if dueDateCol != "" {
			identCheck = append(identCheck, dueDateCol)
		}
		if customerCol != "" {
			identCheck = append(identCheck, customerCol)
		}
		identsOK := true
		for _, ident := range identCheck {
			if !validIdentifier.MatchString(ident) {
				identsOK = false
				break
			}
		}
		if !identsOK {
			continue
		}

		rows, err := s.queryActiveOrders(ctx, connectionID, table, orderCol, statusCol, productCol, wcCol, dueDateCol, customerCol, workCenter)
		if err != nil {
			continue
		}
		for i := range rows {
			if eq, ok := equipmentIdx[normalizeWorkCenter(rows[i].WorkCenter)]; ok {
				rows[i].EquipmentID = eq.ID
			}
		}
		facts = append(facts, rows...)
	}
	return facts, nil
}

// handleActiveProduction exposes ActiveProduction over REST (Entry 120) —
// previously only reachable via the kg_active_production MCP tool. Same data,
// same limitation restated here: current state only, no historical duration.
func (s *server) handleActiveProduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	workCenter := r.URL.Query().Get("work_center")

	facts, err := s.ActiveProduction(r.Context(), workCenter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if facts == nil {
		facts = []ActiveProductionFact{}
	}
	writeJSON(w, map[string]interface{}{"active": facts, "total": len(facts)})
}

// queryActiveOrders runs the actual read-only SELECT against one mapped
// table, filtering the status column against inProgressTokens and, if given,
// the work_center column against workCenter. dueDateCol/customerCol are
// optional (Entry 133) — appended to the SELECT list only when the mapping
// resolved them, so an ERP without those columns still gets the core facts.
func (s *server) queryActiveOrders(ctx context.Context, connectionID, table, orderCol, statusCol, productCol, wcCol, dueDateCol, customerCol, workCenter string) ([]ActiveProductionFact, error) {
	db, err := s.connReg.Get(connectionID)
	if err != nil {
		return nil, err
	}

	statusClauses := make([]string, len(inProgressTokens))
	args := make([]interface{}, 0, len(inProgressTokens)+1)
	for i, tok := range inProgressTokens {
		statusClauses[i] = fmt.Sprintf("LOWER(%s) LIKE ?", statusCol)
		args = append(args, "%"+tok+"%")
	}

	cols := []string{orderCol, productCol, wcCol, statusCol}
	if dueDateCol != "" {
		cols = append(cols, dueDateCol)
	}
	if customerCol != "" {
		cols = append(cols, customerCol)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE (%s)",
		strings.Join(cols, ", "), table, strings.Join(statusClauses, " OR "),
	)
	if workCenter != "" {
		query += fmt.Sprintf(" AND %s = ?", wcCol)
		args = append(args, workCenter)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ActiveProductionFact
	for rows.Next() {
		var orderID, productID, wc, status string
		var dueDate, customerID sql.NullString
		dest := []interface{}{&orderID, &productID, &wc, &status}
		if dueDateCol != "" {
			dest = append(dest, &dueDate)
		}
		if customerCol != "" {
			dest = append(dest, &customerID)
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		fact := ActiveProductionFact{
			ConnectionID: connectionID, WorkCenter: wc, OrderID: orderID, ProductID: productID, Status: status,
			CustomerID: customerID.String,
		}
		if dueDate.Valid {
			fact.DueDate = dueDate.String
			// Truncate to whole days before diffing — a due_date is a
			// calendar day, not a timestamp; comparing raw durations against
			// "now" would round down/flicker across a day depending on the
			// time of day this runs.
			if parsed, err := time.Parse("2006-01-02", dueDate.String[:10]); err == nil {
				today := time.Now().Truncate(24 * time.Hour)
				days := int(parsed.Sub(today).Hours() / 24)
				fact.DaysUntilDue = &days
			}
		}
		out = append(out, fact)
	}
	return out, rows.Err()
}
