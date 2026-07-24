// internal/connections/canonical_suggest.go
// Heuristic canonical-type mapping — Phase 2 of the Track B plan
// (docs/analysis_log.md Entry 115), the IT-side analog of the OT structural
// bootstrap's depth/collision confidence heuristic (internal/kg/bootstrap.go,
// Entry 107). Column-name synonym matching, not ML — appropriate for the same
// reason the OT heuristic is: explainable, auditable, and this is a naming
// convention problem, not a learning problem.
package connections

import "strings"

// MappingCandidate is one table's best-guess canonical mapping, before human
// validation via the KG's existing pending/validate/reject flow.
type MappingCandidate struct {
	Table         string            `json:"table"`
	CanonicalType string            `json:"canonical_type"` // "work_order" | "product"
	Confidence    float64           `json:"confidence"`
	FieldMap      map[string]string `json:"field_map"` // canonical field -> actual column name
}

// suggestionFloor — candidates scoring below this against every canonical
// type aren't suggested at all, so a schema's unrelated tables (operators,
// batches, quality logs, ...) don't get forced into a guess. Mirrors
// SeedFromDiscovery skipping OT entries with no site+work_center.
const suggestionFloor = 0.5

// category is one thing to look for in a table's columns: a canonical field
// name plus a set of synonym substrings, any one of which counts as a match.
type category struct {
	field    string
	synonyms []string
}

// Scoped to 2 canonical types for v0 (deliberately, per Entry 115) — the two
// that unblock "which product is running": work_order and product. The
// fuller 9-object ISA-95-aligned set from Entry 92 is future work.
//
// Core categories are the fields Track B actually needs to be useful at all
// (id, status, product/work-center reference) — worth 80% of the score.
// Bonus categories (customer, due date, margin) sweeten confidence but aren't
// required, since many mid-market ERPs don't expose them
// (docs/it_connectors.md's own agrifood note: often no MES, ERP-only).
var (
	workOrderCore = []category{
		{"order_id", []string{"of_number", "order_number", "order_id", "work_order_id", "wo_number"}},
		{"status", []string{"status", "state"}},
		{"product_id", []string{"product_code", "product_id", "material_id"}},
		{"work_center", []string{"work_center", "workcenter", "machine_id", "line_id"}},
	}
	workOrderBonus = []category{
		{"customer_id", []string{"customer_id", "customer_code", "client_id"}},
		{"due_date", []string{"due_date", "delivery_date", "requested_delivery"}},
	}

	productCore = []category{
		{"product_id", []string{"product_code", "product_id", "material_id", "sku"}},
		{"name", []string{"name", "label", "description"}},
	}
	productBonus = []category{
		{"margin", []string{"margin", "price", "cost_per_unit", "unit_cost"}},
	}
)

// score matches core+bonus categories against a table's columns and returns
// a 0.0-1.0 confidence plus the field_map it would produce. Each category
// matches at most one column (the first one found), and a column can satisfy
// more than one category if its name is genuinely ambiguous — acceptable for
// a v0 heuristic since the result is gated by human validation regardless.
func score(table TableSchema, core, bonus []category) (float64, map[string]string) {
	fieldMap := map[string]string{}
	matchCount := func(cats []category) int {
		n := 0
		for _, cat := range cats {
			for _, col := range table.Columns {
				colLower := strings.ToLower(col.Name)
				matched := false
				for _, syn := range cat.synonyms {
					if strings.Contains(colLower, syn) {
						matched = true
						break
					}
				}
				if matched {
					fieldMap[cat.field] = col.Name
					n++
					break
				}
			}
		}
		return n
	}

	coreN := matchCount(core)
	coreScore := 0.0
	if len(core) > 0 {
		coreScore = float64(coreN) / float64(len(core))
	}

	bonusScore := 0.0
	if len(bonus) > 0 {
		bonusScore = float64(matchCount(bonus)) / float64(len(bonus))
	}

	return coreScore*0.8 + bonusScore*0.2, fieldMap
}

// SuggestMappings scores every table against every known canonical type and
// returns the single best-scoring candidate per table, for tables whose best
// score clears suggestionFloor. A table is never suggested as two things.
func SuggestMappings(tables []TableSchema) []MappingCandidate {
	var out []MappingCandidate
	for _, t := range tables {
		woScore, woFields := score(t, workOrderCore, workOrderBonus)
		prScore, prFields := score(t, productCore, productBonus)

		bestType, bestScore, bestFields := "work_order", woScore, woFields
		if prScore > bestScore {
			bestType, bestScore, bestFields = "product", prScore, prFields
		}

		if bestScore >= suggestionFloor {
			out = append(out, MappingCandidate{
				Table:         t.Name,
				CanonicalType: bestType,
				Confidence:    bestScore,
				FieldMap:      bestFields,
			})
		}
	}
	return out
}
