// cmd/server/mcp_server.go
// Track A of the MCP agent-query proposal (docs/analysis_log.md Entry 113):
// 4 read-only tools over data that already worked end to end from day one —
// KG events/cause/cost and live machine state. Entry 116 added a 5th tool,
// kg_active_production (Track B Phase 4), answering "what's running now"
// from a human-validated ERP work-order mapping.
//
// Still deliberately NOT here: anything answering a HISTORICAL
// product-scoped question ("how long did product B take yesterday"). That
// needs retroactive event-tagging, which doesn't exist yet — adding a tool
// for it now would let an agent silently fabricate an answer instead of
// truthfully saying "I don't know".
//
// Lives in cmd/server (not internal/mcp) because it's a transport/wiring
// concern, same as ws.go and live.go in this same package — the actual query
// logic it calls lives in internal/kg/query.go and active_production.go,
// transport-agnostic and independently testable.
package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
)

// --- kg_query_events ---------------------------------------------------------

type queryEventsArgs struct {
	WorkCenter string `json:"work_center,omitempty" jsonschema:"filter to one machine/work center (e.g. machine1); omit for all"`
	Cause      string `json:"cause,omitempty" jsonschema:"filter to one cause label (e.g. Jam); omit for all"`
	FromTime   string `json:"from_time" jsonschema:"start of the time window, RFC3339 (e.g. 2026-07-21T00:00:00Z)"`
	ToTime     string `json:"to_time" jsonschema:"end of the time window, RFC3339"`
}

type queryEventsOut struct {
	Events []kg.EventDetail `json:"events"`
	Total  int              `json:"total"`
}

func (s *server) toolQueryEvents(_ context.Context, _ *mcp.CallToolRequest, args queryEventsArgs) (*mcp.CallToolResult, queryEventsOut, error) {
	from, err := time.Parse(time.RFC3339, args.FromTime)
	if err != nil {
		return nil, queryEventsOut{}, fmt.Errorf("invalid from_time (want RFC3339): %w", err)
	}
	to, err := time.Parse(time.RFC3339, args.ToTime)
	if err != nil {
		return nil, queryEventsOut{}, fmt.Errorf("invalid to_time (want RFC3339): %w", err)
	}

	events, err := s.kg.QueryEvents(args.WorkCenter, args.Cause, from, to)
	if err != nil {
		return nil, queryEventsOut{}, err
	}
	if events == nil {
		events = []kg.EventDetail{}
	}
	return nil, queryEventsOut{Events: events, Total: len(events)}, nil
}

// --- kg_cost_summary ----------------------------------------------------------

type costSummaryArgs struct {
	FromTime string `json:"from_time,omitempty" jsonschema:"start of the time window, RFC3339 — omit to default to 30 days ago"`
	ToTime   string `json:"to_time,omitempty" jsonschema:"end of the time window, RFC3339 — omit to default to now"`
	GroupBy  string `json:"group_by,omitempty" jsonschema:"'cause' or 'work_center' (default work_center)"`
}

// costPriorityEntry is kg.CostSummaryEntry plus deadline-urgency context
// (Entry 134) — merged server-side rather than left for the caller to
// combine two separate tool calls, so the ranking is deterministic and
// doesn't depend on the model reliably synthesizing kg_cost_summary +
// kg_active_production correctly every time (Entry 132 already showed
// phrasing alone can make a model skip a tool it should use). DaysUntilDue/
// CustomerID are only populated when grouping by work_center — a "cause"
// grouping has no single machine to join a due date against.
//
// Urgent, not a fabricated blended €: per docs/impact_engine.md's locked
// pricing rule (Entry 71), a missed deadline's real cost lives in a contract
// penalty clause most ERPs don't expose — inventing a euro figure for it
// would be a number nobody could audit. So this flags and re-ranks instead
// of pricing: an urgent group can outrank a costlier-but-not-urgent one,
// same "boost" concept as the doc's Enrichment #2, without asserting an
// unaudited number.
type costPriorityEntry struct {
	kg.CostSummaryEntry
	DaysUntilDue *int   `json:"days_until_due,omitempty"`
	CustomerID   string `json:"customer_id,omitempty"`
	Urgent       bool   `json:"urgent,omitempty"`
	// ProductID/ProductName (Entry 135) — the machine grouping this entry is
	// keyed on is an internal identifier; the product it's actively running
	// is what a plant manager actually thinks in terms of, so it's surfaced
	// here too rather than making the caller cross-reference
	// kg_active_production separately just to name what's at stake.
	ProductID   string `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	// Reason (Entry 135) is a ready-made one-line explanation of why this
	// entry ranks where it does — so a caller can lead with the "why" before
	// the euro figure instead of just restating the raw numbers.
	Reason string `json:"reason"`
}

type costSummaryOut struct {
	Groups []costPriorityEntry `json:"groups"`
}

// defaultCostWindowDays is the fallback lookback window when the caller
// omits from_time/to_time — chosen so a vague question ("what are our
// financial priorities") never has to be met with a clarifying question
// about dates before it can even query the data.
const defaultCostWindowDays = 30

// urgentWithinDays mirrors docs/impact_engine.md's Enrichment #2 default
// (customer_commitment.due_date_window_days: 7) — the same threshold that
// document already locked in, reused here rather than inventing a new one.
const urgentWithinDays = 7

func (s *server) toolCostSummary(ctx context.Context, _ *mcp.CallToolRequest, args costSummaryArgs) (*mcp.CallToolResult, costSummaryOut, error) {
	to := time.Now()
	if args.ToTime != "" {
		var err error
		to, err = time.Parse(time.RFC3339, args.ToTime)
		if err != nil {
			return nil, costSummaryOut{}, fmt.Errorf("invalid to_time (want RFC3339): %w", err)
		}
	}
	from := to.AddDate(0, 0, -defaultCostWindowDays)
	if args.FromTime != "" {
		var err error
		from, err = time.Parse(time.RFC3339, args.FromTime)
		if err != nil {
			return nil, costSummaryOut{}, fmt.Errorf("invalid from_time (want RFC3339): %w", err)
		}
	}

	rawGroups, err := s.kg.CostSummary(from, to, args.GroupBy)
	if err != nil {
		return nil, costSummaryOut{}, err
	}

	groups := make([]costPriorityEntry, len(rawGroups))
	for i, g := range rawGroups {
		groups[i] = costPriorityEntry{CostSummaryEntry: g}
	}

	// Only work_center grouping has a machine to join deadline/product data
	// against ("cause" groups span every machine that had that cause).
	// Best-effort throughout: a failure here shouldn't hide the cost data
	// that already succeeded.
	if args.GroupBy == "" || args.GroupBy == "work_center" {
		factByWC := map[string]ActiveProductionFact{}
		if facts, err := s.ActiveProduction(ctx, ""); err == nil {
			for _, f := range facts {
				factByWC[normalizeWorkCenter(f.WorkCenter)] = f
			}
		}
		names, _ := s.productNames(ctx) // product_code -> display name; nil map on failure, lookups just miss

		for i := range groups {
			f, hasFact := factByWC[normalizeWorkCenter(groups[i].Group)]
			if hasFact {
				groups[i].ProductID = f.ProductID
				groups[i].ProductName = names[f.ProductID]
				if f.DaysUntilDue != nil {
					days := *f.DaysUntilDue
					groups[i].DaysUntilDue = &days
					groups[i].CustomerID = f.CustomerID
					groups[i].Urgent = days <= urgentWithinDays
				}
			}
			groups[i].Reason = costPriorityReason(groups[i])
		}
		sort.SliceStable(groups, func(i, j int) bool {
			if groups[i].Urgent != groups[j].Urgent {
				return groups[i].Urgent // urgent groups first, regardless of cost rank
			}
			return groups[i].TotalCostEUR > groups[j].TotalCostEUR // then cost descending, as before
		})
	}
	return nil, costSummaryOut{Groups: groups}, nil
}

// costPriorityReason renders the one-line "why" a caller should lead with
// before the euro figure (Entry 135) — kept server-side so it's the same
// explanation regardless of which model is reading the tool result.
func costPriorityReason(e costPriorityEntry) string {
	product := e.ProductName
	if product == "" {
		product = e.ProductID
	}
	if e.Urgent && product != "" {
		return fmt.Sprintf("%s on %s is due in %d day(s) for %s — missing that deadline costs regardless of this machine's downtime total",
			product, e.Group, *e.DaysUntilDue, e.CustomerID)
	}
	if e.Urgent {
		return fmt.Sprintf("%s has a delivery due in %d day(s)", e.Group, *e.DaysUntilDue)
	}
	if product != "" {
		return fmt.Sprintf("%s (running %s) has the highest micro-stop cost with no near-term deadline pressure", e.Group, product)
	}
	return fmt.Sprintf("%s has the highest micro-stop cost in this window", e.Group)
}

// productNames resolves product_code -> display name via the validated
// 'product' SchemaMapping (internal/connections/canonical_suggest.go),
// mirroring how ActiveProduction resolves 'work_order' mappings. Returns nil
// (not an error) when no product mapping is validated yet — callers treat a
// nil map as "no names available" and fall back to the raw code.
func (s *server) productNames(ctx context.Context) (map[string]string, error) {
	graph, err := s.kg.GetGraph("business")
	if err != nil {
		return nil, err
	}
	for _, n := range graph.Nodes {
		if n.Type != "SchemaMapping" {
			continue
		}
		if pending, _ := n.Properties["pending"].(bool); pending {
			continue
		}
		if ct, _ := n.Properties["canonical_type"].(string); ct != "product" {
			continue
		}
		connectionID, _ := n.Properties["connection_id"].(string)
		table, _ := n.Properties["table"].(string)
		fieldMap, _ := n.Properties["field_map"].(map[string]interface{})
		idCol, _ := fieldMap["product_id"].(string)
		nameCol, _ := fieldMap["name"].(string)
		if idCol == "" || nameCol == "" {
			continue
		}
		for _, ident := range []string{table, idCol, nameCol} {
			if !validIdentifier.MatchString(ident) {
				continue
			}
		}

		db, err := s.connReg.Get(connectionID)
		if err != nil {
			continue
		}
		rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT %s, %s FROM %s", idCol, nameCol, table))
		if err != nil {
			continue
		}
		names := map[string]string{}
		for rows.Next() {
			var code, name string
			if err := rows.Scan(&code, &name); err != nil {
				rows.Close()
				return nil, err
			}
			names[code] = name
		}
		rows.Close()
		return names, rows.Err()
	}
	return nil, nil
}

// --- kg_current_state ---------------------------------------------------------

type currentStateArgs struct {
	WorkCenter string `json:"work_center,omitempty" jsonschema:"filter to one machine; omit to list every known machine"`
}

type machineStateOut struct {
	WorkCenter string    `json:"work_center"`
	Running    bool      `json:"running"`
	Since      time.Time `json:"since"`
}

type currentStateOut struct {
	Machines []machineStateOut `json:"machines"`
}

func (s *server) toolCurrentState(_ context.Context, _ *mcp.CallToolRequest, args currentStateArgs) (*mcp.CallToolResult, currentStateOut, error) {
	if args.WorkCenter != "" {
		st := s.states.get(args.WorkCenter)
		if st == nil {
			return nil, currentStateOut{}, fmt.Errorf("no known state for work center %q", args.WorkCenter)
		}
		return nil, currentStateOut{Machines: []machineStateOut{
			{WorkCenter: args.WorkCenter, Running: st.Running, Since: st.Since},
		}}, nil
	}

	out := currentStateOut{Machines: []machineStateOut{}}
	for wc, st := range s.states.snapshot() {
		out.Machines = append(out.Machines, machineStateOut{WorkCenter: wc, Running: st.Running, Since: st.Since})
	}
	return nil, out, nil
}

// --- kg_describe_node -----------------------------------------------------

type describeNodeArgs struct {
	NodeID string `json:"node_id" jsonschema:"the KG node id to inspect (e.g. an event, equipment, or cause id)"`
}

type kgNodeOut struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type relatedNode struct {
	NodeID   string  `json:"node_id"`
	Label    string  `json:"label"`
	Type     string  `json:"type"`
	Relation string  `json:"relation"`
	Weight   float64 `json:"weight"`
}

type describeNodeOut struct {
	Node     *kgNodeOut    `json:"node"`
	Outgoing []relatedNode `json:"outgoing"`
	Incoming []relatedNode `json:"incoming"`
}

func (s *server) toolDescribeNode(_ context.Context, _ *mcp.CallToolRequest, args describeNodeArgs) (*mcp.CallToolResult, describeNodeOut, error) {
	graph, err := s.kg.GetGraph("all")
	if err != nil {
		return nil, describeNodeOut{}, err
	}

	nodeByID := make(map[string]kg.Node, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodeByID[n.ID] = n
	}

	target, ok := nodeByID[args.NodeID]
	if !ok {
		return nil, describeNodeOut{}, fmt.Errorf("no node found with id %q", args.NodeID)
	}

	out := describeNodeOut{
		Node:     &kgNodeOut{ID: target.ID, Type: target.Type, Label: target.Label, Properties: target.Properties},
		Outgoing: []relatedNode{},
		Incoming: []relatedNode{},
	}
	for _, e := range graph.Edges {
		if e.FromID == args.NodeID {
			if to, ok := nodeByID[e.ToID]; ok {
				out.Outgoing = append(out.Outgoing, relatedNode{NodeID: to.ID, Label: to.Label, Type: to.Type, Relation: e.Relation, Weight: e.Weight})
			}
		}
		if e.ToID == args.NodeID {
			if from, ok := nodeByID[e.FromID]; ok {
				out.Incoming = append(out.Incoming, relatedNode{NodeID: from.ID, Label: from.Label, Type: from.Type, Relation: e.Relation, Weight: e.Weight})
			}
		}
	}
	return nil, out, nil
}

// --- kg_active_production (Track B Phase 4) ---------------------------------

type activeProductionArgs struct {
	WorkCenter string `json:"work_center,omitempty" jsonschema:"filter to one machine; omit for every machine with a validated work-order mapping"`
}

type activeProductionOut struct {
	Active []ActiveProductionFact `json:"active"`
}

func (s *server) toolActiveProduction(ctx context.Context, _ *mcp.CallToolRequest, args activeProductionArgs) (*mcp.CallToolResult, activeProductionOut, error) {
	facts, err := s.ActiveProduction(ctx, args.WorkCenter)
	if err != nil {
		return nil, activeProductionOut{}, err
	}
	if facts == nil {
		facts = []ActiveProductionFact{}
	}
	return nil, activeProductionOut{Active: facts}, nil
}

// --- server construction + mounting -----------------------------------------

// newMCPServer builds the Track A MCP server described at the top of this file.
func newMCPServer(s *server) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{Name: "mindset-data", Version: "0.1.0"}, &mcp.ServerOptions{
		Instructions: "Read-only tools over MindSet Data's Knowledge Graph, live machine state, and (where a " +
			"work-order mapping has been human-validated) live ERP production context. kg_active_production " +
			"answers what's running right now; it does NOT answer historical product-scoped questions like " +
			"how long a product ran yesterday — no tool here can, so say that rather than guessing. Any phrasing " +
			"about cost, spend, financial priorities, biggest problem, or what to fix first is a kg_cost_summary " +
			"question — call it (grouped by work_center is the default; if the user names causes, group by cause) " +
			"before asking the user to clarify. Phrasing about deadlines, delivery dates, or what's urgent is a " +
			"kg_active_production question instead (its days_until_due field) — cost and urgency are separate " +
			"axes that can point at different machines; don't silently pick one when the question could mean " +
			"either without saying so.",
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "kg_query_events",
		Description: "List micro-stop events (with cause and cost, if known) for a machine and/or time window.",
	}, s.toolQueryEvents)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kg_cost_summary",
		Description: "Aggregate micro-stop cost/duration/count over a time window, grouped by cause or by " +
			"machine. This is the tool for 'what are our financial priorities', 'which machine/cause is costing " +
			"us the most', 'top cost drivers', or 'what should we fix first' — the first N groups in the result " +
			"ARE the top-N financial priorities, no extra ranking needed. When grouped by machine (the default), " +
			"each group is already merged with that machine's delivery-deadline urgency (from kg_active_production), " +
			"the product it's currently running (product_id/product_name), and re-ranked: a group with urgent:true " +
			"jumps ahead of costlier-but-not-urgent groups, because a missed deadline has real cost too. This tool " +
			"alone already answers a combined cost-and-deadline priorities question — don't separately call " +
			"kg_active_production and merge it yourself. Each group's `reason` field is a ready-made one-line " +
			"explanation of why it's ranked there; when presenting priorities, lead with that reason (and the " +
			"product/customer it names, not the internal machine id, when reason names one) and put the euro " +
			"figure second, not first — the number alone doesn't explain the decision. If the caller doesn't give " +
			"a time window, use a reasonable default (e.g. the last 30 days) instead of asking what window they " +
			"meant.",
	}, s.toolCostSummary)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "kg_current_state",
		Description: "Get the current Running/Stopped state of one machine, or every known machine.",
	}, s.toolCurrentState)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "kg_describe_node",
		Description: "Look up a single Knowledge Graph node by id and list everything directly connected to it.",
	}, s.toolDescribeNode)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kg_active_production",
		Description: "List the currently active production order + product per machine, read live from a " +
			"human-validated ERP work-order mapping. Answers 'what's running now' only — not historical " +
			"duration questions like 'how long did product B run yesterday'. Also carries the deadline-urgency " +
			"axis, when the ERP mapping has it: due_date, customer_id, and days_until_due per order — use this " +
			"(not kg_cost_summary) for 'should we prioritize product X because its delivery is coming up', " +
			"'what's due soon', or 'which order is most urgent'. Cost priority and deadline priority are two " +
			"different axes and can disagree — a cheap stop on a today-due order can matter more than an " +
			"expensive one with no deadline pressure; say so if the caller seems to conflate them. " +
			"days_until_due absent/null means this ERP mapping doesn't expose a due date — say that rather than " +
			"guessing an urgency.",
	}, s.toolActiveProduction)

	return srv
}

// mountMCP registers the MCP server at /mcp (Streamable HTTP transport,
// stateless — each request is self-contained since every tool here is a
// read-only query, so there's no session state worth keeping across calls).
func mountMCP(mux *http.ServeMux, s *server) {
	mcpServer := newMCPServer(s)
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return mcpServer }, &mcp.StreamableHTTPOptions{
		Stateless: true,
	})
	mux.Handle("/mcp", handler)
}

// runMCPStdio runs the same MCP server as mountMCP, but over stdio instead of
// HTTP — for local MCP clients (Claude Desktop's mcpServers config) that
// launch the server as a subprocess and speak the protocol over its
// stdin/stdout rather than connecting to a URL. No TLS/HTTPS requirement,
// since nothing goes over the network. Blocks until the client disconnects
// (stdin closes); the caller (main) treats that as normal exit.
func runMCPStdio(s *server) error {
	mcpServer := newMCPServer(s)
	return mcpServer.Run(context.Background(), &mcp.StdioTransport{})
}
