// internal/kg/query.go
// Read-only, joined queries over the business-category graph — built for the
// MCP tool layer (cmd/server/mcp_server.go, Track A of the proposal in
// docs/analysis_log.md Entry 113) but kept transport-agnostic so they're
// independently testable and reusable elsewhere (e.g. a future dashboard
// widget could use CostSummary directly).
package kg

import (
	"sort"
	"time"
)

// EventDetail is a flattened, joined view of a micro-stop Event node plus its
// linked Cause (via the caused_by edge) and Cost (via the costs edge) — the
// shape a caller actually wants, instead of walking nodes/edges by hand.
type EventDetail struct {
	EventID         string    `json:"event_id"`
	WorkCenter      string    `json:"work_center"`
	Timestamp       time.Time `json:"timestamp"`
	DurationSeconds float64   `json:"duration_seconds"`
	Cause           string    `json:"cause,omitempty"`
	CauseConfidence float64   `json:"cause_confidence,omitempty"`
	CostEUR         float64   `json:"cost_eur,omitempty"`
}

// QueryEvents returns business Event nodes (micro-stops) with a timestamp in
// [from, to], joined with their Cause and Cost. workCenter/cause filters are
// exact-match and skipped when empty. Events without a parseable timestamp
// are skipped rather than guessed into the window.
func (kg *KnowledgeGraph) QueryEvents(workCenter, cause string, from, to time.Time) ([]EventDetail, error) {
	graph, err := kg.GetGraph(string(CategoryBusiness))
	if err != nil {
		return nil, err
	}

	nodeByID := make(map[string]Node, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodeByID[n.ID] = n
	}

	// Edges out of a given node, grouped by relation, so each Event can look
	// up its Cause/Cost without re-scanning the whole edge list per event.
	type fromRelation struct{ from, relation string }
	edgesFrom := make(map[fromRelation][]Edge)
	for _, e := range graph.Edges {
		k := fromRelation{e.FromID, e.Relation}
		edgesFrom[k] = append(edgesFrom[k], e)
	}

	var out []EventDetail
	for _, n := range graph.Nodes {
		if n.Type != TypeEvent {
			continue
		}
		tsRaw, _ := n.Properties["timestamp"].(string)
		ts, err := time.Parse(time.RFC3339, tsRaw)
		if err != nil {
			continue
		}
		if ts.Before(from) || ts.After(to) {
			continue
		}
		wc, _ := n.Properties["work_center"].(string)
		if workCenter != "" && wc != workCenter {
			continue
		}

		detail := EventDetail{EventID: n.ID, WorkCenter: wc, Timestamp: ts}
		if d, ok := n.Properties["duration_seconds"].(float64); ok {
			detail.DurationSeconds = d
		}
		for _, e := range edgesFrom[fromRelation{n.ID, "caused_by"}] {
			if causeNode, ok := nodeByID[e.ToID]; ok {
				detail.Cause = causeNode.Label
				detail.CauseConfidence = e.Weight
			}
		}
		for _, e := range edgesFrom[fromRelation{n.ID, "costs"}] {
			if costNode, ok := nodeByID[e.ToID]; ok {
				if amt, ok := costNode.Properties["amount_eur"].(float64); ok {
					detail.CostEUR = amt
				}
			}
		}

		if cause != "" && detail.Cause != cause {
			continue
		}
		out = append(out, detail)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Timestamp.Before(out[j].Timestamp) })
	return out, nil
}

// CostSummaryEntry is one grouped row of CostSummary's output.
type CostSummaryEntry struct {
	Group          string  `json:"group"`
	EventCount     int     `json:"event_count"`
	TotalDurationS float64 `json:"total_duration_seconds"`
	TotalCostEUR   float64 `json:"total_cost_eur"`
}

// CostSummary aggregates events in [from, to], grouped by "cause" or
// "work_center" (any other value falls back to "work_center"). Sorted
// descending by total cost — highest-impact group first, matching the
// "Top 3 actions" framing in docs/impact_engine.md.
func (kg *KnowledgeGraph) CostSummary(from, to time.Time, groupBy string) ([]CostSummaryEntry, error) {
	events, err := kg.QueryEvents("", "", from, to)
	if err != nil {
		return nil, err
	}

	groups := make(map[string]*CostSummaryEntry)
	var order []string
	for _, e := range events {
		key := e.WorkCenter
		if groupBy == "cause" {
			key = e.Cause
			if key == "" {
				key = "(unknown)"
			}
		}
		g, ok := groups[key]
		if !ok {
			g = &CostSummaryEntry{Group: key}
			groups[key] = g
			order = append(order, key)
		}
		g.EventCount++
		g.TotalDurationS += e.DurationSeconds
		g.TotalCostEUR += e.CostEUR
	}

	out := make([]CostSummaryEntry, 0, len(groups))
	for _, k := range order {
		out = append(out, *groups[k])
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TotalCostEUR > out[j].TotalCostEUR })
	return out, nil
}
