// internal/uns/mapper.go
package uns

import (
	"fmt"
	"strings"
)

// UNSNode represents a fully contextualized node in the ISA-95 hierarchy.
// Every signal in the factory gets exactly one UNSNode.
type UNSNode struct {
	// ISA-95 hierarchy levels
	Site       string // e.g. "local-test" or "usine-paris-nord"
	Area       string // e.g. "area1", "packaging-line"
	WorkCenter string // e.g. "machine1", "reactor-3"
	WorkUnit   string // e.g. "ligne1" — optional, only if depth >= 3

	// Tag metadata
	TagName     string // normalized: "temperature", "pressure", "status"
	OriginalName string // raw OPC-UA name: "machine1.temp"
	DataType    string // "Float", "Boolean", "Int32"...
	Unit        string // inferred: "celsius", "bar", "rpm", ""
	Description string // human-readable description
}

// Topic returns the full MQTT topic string for this UNS node.
// Format: <site>/<area>/<workcenter>[/<workunit>]/<tagname>
// Example: local-test/area1/machine1/temperature
func (n UNSNode) Topic() string {
	parts := []string{n.Site}

	if n.Area != "" {
		parts = append(parts, n.Area)
	}
	if n.WorkCenter != "" {
		parts = append(parts, n.WorkCenter)
	}
	// WorkUnit is optional — only add if non-empty
	if n.WorkUnit != "" {
		parts = append(parts, n.WorkUnit)
	}
	if n.TagName != "" {
		parts = append(parts, n.TagName)
	}

	return strings.Join(parts, "/")
}

// FullTopic returns the MQTT topic prefixed with "mindset/site/"
// This is what gets published on the broker.
// Example: mindset/site/local-test/area1/machine1/temperature
func (n UNSNode) FullTopic() string {
	return "mindset/site/" + n.Topic()
}

// Mapper converts raw OPC-UA tag names into ISA-95 UNS nodes.
// It uses dot-separated naming conventions to infer the hierarchy.
type Mapper struct {
	SiteID string
}

// NewMapper creates a new mapper for the given site.
func NewMapper(siteID string) *Mapper {
	return &Mapper{SiteID: siteID}
}

// MapTag is the core function. It takes a raw OPC-UA tag name and dataType,
// and returns a fully populated UNSNode ready for publication.
//
// Naming convention assumed (dot-separated):
//   "temp"                    → root-level tag (no machine context)
//   "machine1.temp"           → 2 levels: workcenter + tag
//   "machine2.ligne1.presion" → 3 levels: workcenter + workunit + tag
//   "a.b.c.d"                 → 4+ levels: area + workcenter + workunit + tag
func (m *Mapper) MapTag(tagName, dataType string) UNSNode {
	parts := strings.Split(tagName, ".")

	node := UNSNode{
		Site:        m.SiteID,
		OriginalName: tagName,
		DataType:    dataType,
	}

	switch len(parts) {
	case 1:
		// Root-level tag — no machine context
		// e.g. "temp", "stat"
		node.Area = "general"
		node.WorkCenter = "site"
		node.TagName = normalizeTagName(parts[0])

	case 2:
		// Most common case: machine + tag
		// e.g. "machine1.temp" → area1 / machine1 / temperature
		node.Area = "area1"
		node.WorkCenter = parts[0]
		node.TagName = normalizeTagName(parts[1])

	case 3:
		// 3 levels: machine + subunit + tag
		// e.g. "machine2.ligne1.presion" → area1 / machine2 / ligne1 / pressure
		node.Area = "area1"
		node.WorkCenter = parts[0]
		node.WorkUnit = parts[1]
		node.TagName = normalizeTagName(parts[2])

	default:
		// 4+ levels: first part = area, second = workcenter,
		// middle parts = workunit (joined), last = tag
		node.Area = parts[0]
		node.WorkCenter = parts[1]
		node.WorkUnit = strings.Join(parts[2:len(parts)-1], "/")
		node.TagName = normalizeTagName(parts[len(parts)-1])
	}

	// Infer unit from the normalized tag name
	node.Unit = inferUnit(node.TagName, dataType)

	// Build a human-readable description
	node.Description = buildDescription(node)

	return node
}

// MapTags maps a slice of (name, dataType) pairs to UNS nodes.
// Useful for bulk mapping at startup.
func (m *Mapper) MapTags(tags []struct{ Name, DataType string }) []UNSNode {
	nodes := make([]UNSNode, 0, len(tags))
	for _, t := range tags {
		nodes = append(nodes, m.MapTag(t.Name, t.DataType))
	}
	return nodes
}

// PrintTree prints the ISA-95 hierarchy to stdout in a readable tree format.
// Use this at startup to verify the mapping is correct.
func PrintTree(nodes []UNSNode) {
	fmt.Println("\n[UNS] ISA-95 Hierarchy:")
	fmt.Printf("└── Site\n")

	// Track what we've already printed to avoid duplicates
	seenArea := map[string]bool{}
	seenWC := map[string]bool{}

	for _, n := range nodes {
		areaKey := n.Area
		if !seenArea[areaKey] {
			fmt.Printf("    └── %s\n", n.Area)
			seenArea[areaKey] = true
		}

		wcKey := n.Area + "/" + n.WorkCenter
		if n.WorkCenter != "" && !seenWC[wcKey] {
			fmt.Printf("        └── %s\n", n.WorkCenter)
			seenWC[wcKey] = true
		}

		if n.WorkUnit != "" {
			fmt.Printf("            └── %s\n", n.WorkUnit)
			fmt.Printf("                └── %-20s (%s) [%s] → %s\n",
				n.TagName, n.DataType, n.Unit, n.FullTopic())
		} else {
			fmt.Printf("            └── %-20s (%s) [%s] → %s\n",
				n.TagName, n.DataType, n.Unit, n.FullTopic())
		}
	}
	fmt.Println()
}

// normalizeTagName converts common OPC-UA abbreviations to full names.
// This is crucial for the Knowledge Graph and dashboard readability.
// Add entries here as you discover new abbreviations in the field.
func normalizeTagName(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))

	// Map of abbreviations → full names
	// Ordered by specificity — more specific matches first
	replacements := map[string]string{
		// Temperature variants
		"temp":        "temperature",
		"tmp":         "temperature",
		"temperature": "temperature",

		// Pressure variants
		"presion":  "pressure", // typo common in French SCADA config
		"pression": "pressure",
		"press":    "pressure",
		"pres":     "pressure",
		"pressure": "pressure",

		// Status / state
		"stat":   "status",
		"state":  "status",
		"status": "status",
		"etat":   "status", // French
		"run":    "status",

		// Speed / velocity
		"spd":    "speed",
		"speed":  "speed",
		"vitesse": "speed", // French
		"vit":    "speed",

		// Counters
		"ctr":     "counter",
		"cnt":     "counter",
		"counter": "counter",
		"count":   "counter",
		"compteur": "counter", // French

		// Power / energy
		"pwr":    "power",
		"power":  "power",
		"puissance": "power", // French
		"energy": "energy",
		"conso":  "energy", // French abbrev for consommation

		// Flow
		"flw":  "flow",
		"flow": "flow",
		"debit": "flow", // French

		// Level
		"lvl":   "level",
		"level": "level",
		"niveau": "level", // French

		// Current (electrical)
		"curr":    "current",
		"current": "current",
		"courant": "current", // French

		// Voltage
		"volt":    "voltage",
		"voltage": "voltage",
		"tension": "voltage", // French

		// Alarm / fault
		"alm":    "alarm",
		"alarm":  "alarm",
		"fault":  "fault",
		"defaut": "fault", // French défaut

		// Setpoint
		"sp":       "setpoint",
		"setpoint": "setpoint",
		"consigne": "setpoint", // French

		// Vibration
		"vib":       "vibration",
		"vibration": "vibration",
	}

	if normalized, ok := replacements[raw]; ok {
		return normalized
	}

	// If no match found, return the raw name as-is
	// (better than guessing wrong)
	return raw
}

// inferUnit deduces the physical unit from the normalized tag name and data type.
// Returns empty string if the unit cannot be inferred (e.g. for counters, status).
func inferUnit(normalizedName, dataType string) string {
	unitMap := map[string]string{
		"temperature": "celsius",
		"pressure":    "bar",
		"speed":       "rpm",
		"flow":        "m3h",   // m³/h
		"level":       "pct",   // percentage
		"power":       "kw",
		"energy":      "kwh",
		"current":     "a",     // amperes
		"voltage":     "v",     // volts
		"vibration":   "mms",   // mm/s
		// Intentionally no unit for: status, counter, alarm, fault, setpoint
	}

	if unit, ok := unitMap[normalizedName]; ok {
		return unit
	}

	// For boolean types, unit is never meaningful
	if dataType == "Boolean" {
		return ""
	}

	return ""
}

// buildDescription creates a human-readable description for the tag.
func buildDescription(n UNSNode) string {
	location := n.WorkCenter
	if n.WorkUnit != "" {
		location = n.WorkCenter + " / " + n.WorkUnit
	}
	return fmt.Sprintf("%s sensor on %s", n.TagName, location)
}