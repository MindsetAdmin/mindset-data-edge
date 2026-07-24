// internal/kg/bootstrap.go
// v0 structural bootstrap — docs/analysis_log.md Entries 87/89/90/92/94/95/96,
// confidence-gated as of Entry 107. Auto-generates the OT Equipment/Area/Site
// skeleton from discovered structure (already ISA-95-mapped by the caller —
// kg stays protocol-agnostic) instead of waiting for a micro-stop event.
// Nodes with confidence >= AutoAcceptThreshold are written already validated
// (pending: false) — a human only reviews the ones the heuristic itself is
// unsure about, not every single generated node. No separate staging table
// yet (that refinement, from the Synapt comparison in Entry 89, can come
// later without losing this work).
package kg

import (
	"encoding/json"
	"fmt"
)

// AutoAcceptThreshold — nodes scoring at or above this confidence are written
// already validated; below it, they're flagged pending for human review.
// A named constant, not a config option, for v0 — see docs/analysis_log.md
// Entry 107 for the confidence heuristic itself (computed by the caller,
// cmd/server/opcua.go's seedKG — kg stays protocol-agnostic and just applies
// the gate).
const AutoAcceptThreshold = 0.7

// HierarchyEntry is one discovered OT structural element, already mapped to
// ISA-95 (Site/Area/WorkCenter/WorkUnit/Tag) by the caller.
//
// Depth is the raw tag name's dot-segment count (before mapping) — needed
// because internal/uns.Mapper gives WorkCenter/WorkUnit different real-world
// meanings at different depths: at depth 3, WorkCenter IS the machine and
// WorkUnit is a sub-component of it; at depth 4+, WorkCenter is a grouping
// level (e.g. a line) ABOVE the machine, and WorkUnit is the machine itself.
// See internal/uns/mapper.go's doc comment and docs/analysis_log.md Entry 98.
//
// Confidence is a 0.0-1.0 score computed by the caller across the whole
// discovered batch (depth-consistency + no naming collisions within the same
// equipment — see cmd/server/opcua.go's seedKG). Not ML — an explainable
// heuristic appropriate for a naming-convention mapper.
type HierarchyEntry struct {
	Site       string
	Area       string
	WorkCenter string
	WorkUnit   string
	Depth      int
	Confidence float64
	TagName    string
	NodeID     string // original OPC-UA node id, kept for traceability
	DataType   string
}

// SeedFromDiscovery writes the Site/Area/[WorkCenter]/Equipment/Tag skeleton
// for newly discovered OT structure into the business category. Nodes derived
// from a low-confidence mapping are flagged pending human validation; nodes
// at or above AutoAcceptThreshold are written already confirmed — see Entry
// 107. Idempotent (AddNodeCat/AddEdgeCat use INSERT OR IGNORE), so calling
// this again after every Discover() is safe and cheap.
//
// Equipment identity (fixed in Entry 98 after a live-test finding — see
// docs/analysis_log.md Entry 97/98): the node typed Equipment must match what
// the rest of the system means by "work center" — AddMicroStop, ERP
// work_orders.work_center, of_enrichment.yaml — which is the physical
// machine, not a grouping level. At depth 3, the mapper's WorkCenter field
// already IS the machine. At depth 4+, WorkCenter is a grouping level (e.g. a
// line) ABOVE the machine, and WorkUnit is the actual machine — so at that
// depth an extra "WorkCenter" node is created for the grouping, and Equipment
// uses WorkUnit instead. Getting this wrong was the real bug a live test
// against a real Prosys server found: two different machines' tags were
// landing on one Equipment node because WorkUnit was dropped entirely.
//
// Equipment nodes reuse AddMicroStop's ID scheme ("machine_<name>") so the
// two paths converge on one node instead of duplicating once real operational
// data starts flowing for the same machine. Because AddNodeCat is INSERT OR
// IGNORE, whichever path writes first wins the row — if this seed runs
// first, the node stays pending:true even after a real micro-stop
// references it, until someone explicitly validates it. Acceptable for v0.
func (kg *KnowledgeGraph) SeedFromDiscovery(entries []HierarchyEntry) (int, error) {
	seeded := 0
	for _, e := range entries {
		if e.Site == "" || e.WorkCenter == "" {
			continue // nothing structural to seed without at least a site + work center
		}
		pending := e.Confidence < AutoAcceptThreshold

		siteID := fmt.Sprintf("site_%s", e.Site)
		if err := kg.AddNodeCat(CategoryBusiness, siteID, "Site", e.Site, map[string]interface{}{
			"pending":    pending,
			"confidence": e.Confidence,
		}); err != nil {
			return seeded, err
		}

		areaID := siteID
		if e.Area != "" {
			areaID = fmt.Sprintf("area_%s_%s", e.Site, e.Area)
			if err := kg.AddNodeCat(CategoryBusiness, areaID, "Area", e.Area, map[string]interface{}{
				"pending":    pending,
				"confidence": e.Confidence,
				"site":       e.Site,
			}); err != nil {
				return seeded, err
			}
			if err := kg.AddEdgeCat(CategoryBusiness, fmt.Sprintf("edge_%s_%s", siteID, areaID), siteID, areaID, "contains", 1.0); err != nil {
				return seeded, err
			}
		}

		// Resolve the Equipment identity + its parent, branching on depth
		// (see the doc comment above — WorkCenter/WorkUnit swap roles).
		equipmentName := e.WorkCenter
		equipmentParentID := areaID
		if e.Depth >= 4 && e.WorkUnit != "" {
			lineID := fmt.Sprintf("wc_%s_%s_%s", e.Site, e.Area, e.WorkCenter)
			if err := kg.AddNodeCat(CategoryBusiness, lineID, "WorkCenter", e.WorkCenter, map[string]interface{}{
				"pending":    pending,
				"confidence": e.Confidence,
				"area":       e.Area,
				"site":       e.Site,
			}); err != nil {
				return seeded, err
			}
			if err := kg.AddEdgeCat(CategoryBusiness, fmt.Sprintf("edge_%s_%s", areaID, lineID), areaID, lineID, "contains", 1.0); err != nil {
				return seeded, err
			}
			equipmentName = e.WorkUnit
			equipmentParentID = lineID
		}

		eqID := fmt.Sprintf("machine_%s", equipmentName) // matches AddMicroStop's ID scheme — see doc comment above
		if err := kg.AddNodeCat(CategoryBusiness, eqID, TypeEquipment, equipmentName, map[string]interface{}{
			"pending":     pending,
			"confidence":  e.Confidence,
			"work_center": equipmentName,
			"area":        e.Area,
			"site":        e.Site,
		}); err != nil {
			return seeded, err
		}
		if err := kg.AddEdgeCat(CategoryBusiness, fmt.Sprintf("edge_%s_%s", equipmentParentID, eqID), equipmentParentID, eqID, "contains", 1.0); err != nil {
			return seeded, err
		}

		if e.TagName != "" {
			tagID := fmt.Sprintf("tag_%s_%s_%s_%s", e.Site, e.Area, equipmentName, e.TagName)
			if err := kg.AddNodeCat(CategoryBusiness, tagID, "Tag", e.TagName, map[string]interface{}{
				"pending":    pending,
				"confidence": e.Confidence,
				"node_id":    e.NodeID,
				"data_type":  e.DataType,
			}); err != nil {
				return seeded, err
			}
			if err := kg.AddEdgeCat(CategoryBusiness, fmt.Sprintf("edge_%s_%s", eqID, tagID), eqID, tagID, "has_tag", 1.0); err != nil {
				return seeded, err
			}
		}

		seeded++
	}
	return seeded, nil
}

// ListPending returns every business-category node still flagged pending:true.
// Filtered in Go, not SQL — consistent with the rest of this package, which
// treats `properties` as an opaque JSON blob rather than relying on SQLite's
// JSON1 extension.
func (kg *KnowledgeGraph) ListPending() ([]Node, error) {
	graph, err := kg.GetGraph(string(CategoryBusiness))
	if err != nil {
		return nil, err
	}
	pending := make([]Node, 0)
	for _, n := range graph.Nodes {
		if isPending(n.Properties) {
			pending = append(pending, n)
		}
	}
	return pending, nil
}

func isPending(props map[string]interface{}) bool {
	if props == nil {
		return false
	}
	v, ok := props["pending"].(bool)
	return ok && v
}

// ValidateNode clears the pending flag on a node — it's now confirmed correct
// and behaves like any other business-category node.
func (kg *KnowledgeGraph) ValidateNode(id string) error {
	return kg.setPending(id, false)
}

// RejectNode deletes a pending node outright, along with any edges touching
// it. Used when the auto-generated structure was wrong.
func (kg *KnowledgeGraph) RejectNode(id string) error {
	if _, err := kg.store.DB().Exec(`DELETE FROM kg_edges WHERE from_id = ? OR to_id = ?`, id, id); err != nil {
		return err
	}
	_, err := kg.store.DB().Exec(`DELETE FROM kg_nodes WHERE id = ?`, id)
	return err
}

func (kg *KnowledgeGraph) setPending(id string, pending bool) error {
	row := kg.store.DB().QueryRow(`SELECT properties FROM kg_nodes WHERE id = ?`, id)
	var propsJSON string
	if err := row.Scan(&propsJSON); err != nil {
		return fmt.Errorf("node %q not found: %w", id, err)
	}

	props := map[string]interface{}{}
	if propsJSON != "" {
		if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
			return err
		}
	}
	props["pending"] = pending

	propsBytes, err := json.Marshal(props)
	if err != nil {
		return err
	}
	_, err = kg.store.DB().Exec(`UPDATE kg_nodes SET properties = ? WHERE id = ?`, string(propsBytes), id)
	return err
}
