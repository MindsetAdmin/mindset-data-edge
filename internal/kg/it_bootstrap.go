// internal/kg/it_bootstrap.go
// IT-side structural bootstrap — Phase 2 of the Track B plan
// (docs/analysis_log.md Entry 115), mirroring internal/kg/bootstrap.go's OT
// pattern: write auto-generated structure immediately, gated by the same
// AutoAcceptThreshold, reviewed through the same pending/validate/reject flow.
// SchemaMapping is a new node type, but it needs none of ListPending/
// ValidateNode/RejectNode to change — they're already generic over any
// business-category node carrying a `pending` property.
package kg

import "fmt"

// SchemaMappingCandidate is one table's best-guess canonical mapping, ready
// to persist. Mirrors connections.MappingCandidate without kg importing the
// connections package — the caller (cmd/server) converts between the two.
type SchemaMappingCandidate struct {
	Table         string
	CanonicalType string
	Confidence    float64
	FieldMap      map[string]interface{} // JSON-friendly copy for node properties
}

// SeedSchemaMappings writes one SchemaMapping node per candidate under the
// given connection. Idempotent (AddNodeCat is INSERT OR IGNORE), so calling
// this again after every /discover is safe and cheap — same as
// SeedFromDiscovery on the OT side.
func (kg *KnowledgeGraph) SeedSchemaMappings(connectionID string, candidates []SchemaMappingCandidate) (int, error) {
	seeded := 0
	for _, c := range candidates {
		pending := c.Confidence < AutoAcceptThreshold
		id := fmt.Sprintf("mapping_%s_%s", connectionID, c.Table)
		label := fmt.Sprintf("%s (%s)", c.Table, c.CanonicalType)
		props := map[string]interface{}{
			"pending":        pending,
			"confidence":     c.Confidence,
			"connection_id":  connectionID,
			"table":          c.Table,
			"canonical_type": c.CanonicalType,
			"field_map":      c.FieldMap,
		}
		if err := kg.AddNodeCat(CategoryBusiness, id, "SchemaMapping", label, props); err != nil {
			return seeded, err
		}
		seeded++
	}
	return seeded, nil
}
