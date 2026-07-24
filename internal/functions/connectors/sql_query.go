// internal/functions/connectors/sql_query.go
package connectors

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
)

const (
	defaultTimeout = 30 * time.Second
	maxTimeout     = 300 * time.Second
	defaultLimit   = 1000
	maxLimit       = 10000
)

// ConnectionGetter is satisfied by *connections.Registry. Accepting the
// narrow interface instead of the concrete type keeps SQLQueryHandler
// testable against a fake pool (e.g. an in-memory SQLite DB) without
// depending on internal/connections in tests.
type ConnectionGetter interface {
	Get(id string) (*sql.DB, error)
}

// SQLQueryHandler executes sql_query nodes: a parameterized, read-only
// SELECT against a connection resolved via ConnectionGetter.
type SQLQueryHandler struct {
	connections ConnectionGetter
}

// NewSQLQueryHandler creates a handler backed by the given connection pool.
func NewSQLQueryHandler(connections ConnectionGetter) *SQLQueryHandler {
	return &SQLQueryHandler{connections: connections}
}

// GetFunction retourne la définition de la fonction.
func (h *SQLQueryHandler) GetFunction() *functions.Function {
	return &functions.Function{
		Name:        "sql_query",
		Type:        functions.TypeConnector,
		Description: "Exécute une requête SELECT paramétrée en lecture seule sur une connexion SQL configurée.",
		Handler: func(params map[string]interface{}) (interface{}, error) {
			return h.Execute(context.Background(), params)
		},
	}
}

// Execute runs the configured query and returns rows (raw) plus, when
// field_map is present, a canonical-mapped copy of each row.
func (h *SQLQueryHandler) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	connID, _ := params["connection_id"].(string)
	query, _ := params["query"].(string)
	if connID == "" {
		return nil, fmt.Errorf("sql_query: missing connection_id")
	}
	if query == "" {
		return nil, fmt.Errorf("sql_query: missing query")
	}

	queryParams, _ := params["params"].(map[string]interface{})
	timeout := durationOr(params["timeout_seconds"], defaultTimeout, maxTimeout)
	limit := intOr(params["limit"], defaultLimit, maxLimit)

	if err := ensureSelectOnly(query); err != nil {
		return nil, err
	}
	boundQuery, args, err := bindPositional(query, queryParams)
	if err != nil {
		return nil, err
	}
	boundQuery = ensureLimit(boundQuery, limit)

	db, err := h.connections.Get(connID)
	if err != nil {
		return nil, fmt.Errorf("sql_query: %w", err)
	}

	qctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	started := time.Now()
	rows, err := db.QueryContext(qctx, boundQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("sql_query: %w", err)
	}
	defer rows.Close()

	out, err := mapRows(rows)
	if err != nil {
		return nil, fmt.Errorf("sql_query: %w", err)
	}

	result := map[string]interface{}{
		"rows":      out,
		"row_count": len(out),
		"query_ms":  time.Since(started).Milliseconds(),
	}

	fieldMap, _ := params["field_map"].(map[string]interface{})
	if len(fieldMap) > 0 {
		canonical, err := applyFieldMap(out, fieldMap)
		if err != nil {
			return nil, fmt.Errorf("sql_query: %w", err)
		}
		canonicalType, _ := params["canonical"].(string)
		result["canonical"] = canonical
		result["canonical_type"] = canonicalType
	} else {
		// No field_map configured — degrade gracefully rather than error.
		result["canonical"] = []map[string]interface{}{}
		result["canonical_type"] = nil
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Guards
// ---------------------------------------------------------------------------

var leadingKeywordRe = regexp.MustCompile(`^[A-Za-z]+`)

// ensureSelectOnly strips leading comments/whitespace and rejects anything
// but a single SELECT statement.
func ensureSelectOnly(query string) error {
	body := stripLeadingNoise(query)
	kw := leadingKeywordRe.FindString(body)
	if !strings.EqualFold(kw, "SELECT") {
		if kw == "" {
			return fmt.Errorf("sql_query: only SELECT statements are allowed")
		}
		return fmt.Errorf("sql_query: only SELECT statements are allowed, got %q", strings.ToUpper(kw))
	}
	if hasMultipleStatements(query) {
		return fmt.Errorf("sql_query: multi-statement queries are not allowed")
	}
	return nil
}

// stripLeadingNoise repeatedly trims leading whitespace and SQL comments
// (--, #, /* */) so the real leading keyword can be inspected.
func stripLeadingNoise(s string) string {
	for {
		before := s
		s = strings.TrimLeft(s, " \t\r\n")
		switch {
		case strings.HasPrefix(s, "--"):
			if i := strings.IndexByte(s, '\n'); i >= 0 {
				s = s[i+1:]
			} else {
				s = ""
			}
		case strings.HasPrefix(s, "#"):
			if i := strings.IndexByte(s, '\n'); i >= 0 {
				s = s[i+1:]
			} else {
				s = ""
			}
		case strings.HasPrefix(s, "/*"):
			if i := strings.Index(s, "*/"); i >= 0 {
				s = s[i+2:]
			} else {
				s = ""
			}
		}
		if s == before {
			return s
		}
	}
}

// hasMultipleStatements reports whether query contains a semicolon outside
// a string literal, other than a single trailing terminator. This is a
// lightweight heuristic guard, not a full SQL parser.
func hasMultipleStatements(query string) bool {
	q := strings.TrimRight(query, " \t\r\n")
	q = strings.TrimSuffix(q, ";")
	return containsUnquotedSemicolon(q)
}

func containsUnquotedSemicolon(s string) bool {
	var quote rune
	for i := 0; i < len(s); i++ {
		c := rune(s[i])
		switch {
		case quote != 0:
			if c == '\\' && quote != '`' {
				i++
				continue
			}
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"' || c == '`':
			quote = c
		case c == ';':
			return true
		}
	}
	return false
}

var namedParamRe = regexp.MustCompile(`:([A-Za-z_][A-Za-z0-9_]*)`)

// bindPositional converts MindSet's :name placeholders into the driver's
// positional "?" placeholders, building the args slice in occurrence order.
// Repeated uses of the same name are each bound to the same value.
func bindPositional(query string, params map[string]interface{}) (string, []interface{}, error) {
	var missing string
	args := []interface{}{}
	bound := namedParamRe.ReplaceAllStringFunc(query, func(match string) string {
		name := match[1:]
		val, ok := params[name]
		if !ok {
			missing = name
			return match
		}
		args = append(args, val)
		return "?"
	})
	if missing != "" {
		return "", nil, fmt.Errorf("sql_query: missing value for parameter %q", missing)
	}
	return bound, args, nil
}

var limitRe = regexp.MustCompile(`(?i)\bLIMIT\s+(\d+)\b`)

// ensureLimit caps an existing LIMIT clause at max, or appends one if the
// query doesn't have one.
func ensureLimit(query string, max int) string {
	if loc := limitRe.FindStringSubmatchIndex(query); loc != nil {
		n, err := strconv.Atoi(query[loc[2]:loc[3]])
		if err == nil && n <= max {
			return query
		}
		return query[:loc[2]] + strconv.Itoa(max) + query[loc[3]:]
	}
	trimmed := strings.TrimRight(query, " \t\r\n")
	trimmed = strings.TrimSuffix(trimmed, ";")
	return fmt.Sprintf("%s LIMIT %d", trimmed, max)
}

func intOr(v interface{}, def, max int) int {
	n := 0
	switch x := v.(type) {
	case float64:
		n = int(x)
	case int:
		n = x
	}
	if n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

func durationOr(v interface{}, def, max time.Duration) time.Duration {
	secs := 0
	switch x := v.(type) {
	case float64:
		secs = int(x)
	case int:
		secs = x
	}
	if secs <= 0 {
		return def
	}
	d := time.Duration(secs) * time.Second
	if d > max {
		return max
	}
	return d
}

// ---------------------------------------------------------------------------
// Row mapping — type layer (§6): MySQL type -> Go -> JSON-clean value
// ---------------------------------------------------------------------------

// columnMeta is a minimal extraction from *sql.ColumnType so coerce can be
// unit tested without a live database connection.
type columnMeta struct {
	dbType    string
	length    int64
	hasLength bool
}

func columnMetaFrom(ct *sql.ColumnType) columnMeta {
	length, ok := ct.Length()
	return columnMeta{
		dbType:    strings.ToUpper(ct.DatabaseTypeName()),
		length:    length,
		hasLength: ok,
	}
}

func mapRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	out := []map[string]interface{}{}
	for rows.Next() {
		raw := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			row[col] = coerce(raw[i], columnMetaFrom(colTypes[i]))
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// coerce converts a driver-scanned value into something that JSON-encodes
// cleanly, per the type table in docs/mysql_connector.md §6.
func coerce(v interface{}, cm columnMeta) interface{} {
	if v == nil {
		return nil
	}

	switch x := v.(type) {
	case []byte:
		switch cm.dbType {
		case "JSON":
			var j interface{}
			if err := json.Unmarshal(x, &j); err == nil {
				return j
			}
			return string(x)
		case "BLOB", "TINYBLOB", "MEDIUMBLOB", "LONGBLOB", "BINARY", "VARBINARY":
			return base64.StdEncoding.EncodeToString(x)
		default:
			// DECIMAL/NUMERIC, TIME, VARCHAR/TEXT/CHAR, ENUM, SET all arrive
			// as raw bytes and decode cleanly as strings.
			return string(x)
		}
	case int64:
		if cm.dbType == "TINYINT" && cm.hasLength && cm.length == 1 {
			return x != 0
		}
		return x
	case time.Time:
		return x.UTC().Format(time.RFC3339)
	default:
		// float64, float32, bool, string pass through as-is.
		return v
	}
}

// ---------------------------------------------------------------------------
// Semantic mapping — the other layer (§6b): canonical model + field_map/value_map
// ---------------------------------------------------------------------------

// applyFieldMap builds one canonical row per raw row. fieldMap maps a
// canonical field name to either a raw column name (string) or an enum
// translation spec: {"from": "<column>", "values": {"<raw>": "<canonical>"}}.
func applyFieldMap(rows []map[string]interface{}, fieldMap map[string]interface{}) ([]map[string]interface{}, error) {
	canonical := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		out := make(map[string]interface{}, len(fieldMap))
		for canonicalField, spec := range fieldMap {
			switch s := spec.(type) {
			case string:
				out[canonicalField] = row[s]
			case map[string]interface{}:
				from, _ := s["from"].(string)
				if from == "" {
					return nil, fmt.Errorf("field_map %q: missing \"from\"", canonicalField)
				}
				values, _ := s["values"].(map[string]interface{})
				out[canonicalField] = applyValueMap(row[from], values)
			default:
				return nil, fmt.Errorf("field_map %q: unsupported spec type %T", canonicalField, spec)
			}
		}
		canonical = append(canonical, out)
	}
	return canonical, nil
}

// applyValueMap translates a raw enum value to its canonical label. A raw
// value with no matching entry passes through unchanged — an incomplete
// enum map shouldn't break the pipeline.
func applyValueMap(raw interface{}, values map[string]interface{}) interface{} {
	if len(values) == 0 || raw == nil {
		return raw
	}
	key := fmt.Sprintf("%v", raw)
	if mapped, ok := values[key]; ok {
		return mapped
	}
	return raw
}
