// internal/functions/connectors/sql_query_test.go
package connectors

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// ---------------------------------------------------------------------------
// ensureSelectOnly
// ---------------------------------------------------------------------------

func TestEnsureSelectOnly_accepts_SELECT(t *testing.T) {
	cases := []string{
		"SELECT 1",
		"  select 1",
		"-- a comment\nSELECT 1",
		"/* block */ SELECT 1",
		"# mysql line comment\nSELECT 1",
		"SELECT 1;",
	}
	for _, q := range cases {
		if err := ensureSelectOnly(q); err != nil {
			t.Errorf("query %q: expected no error, got %v", q, err)
		}
	}
}

func TestEnsureSelectOnly_rejects_INSERT_UPDATE_DELETE_DROP_ALTER_TRUNCATE_CALL(t *testing.T) {
	cases := []string{
		"INSERT INTO t VALUES (1)",
		"UPDATE t SET x=1",
		"DELETE FROM t",
		"DROP TABLE t",
		"ALTER TABLE t ADD COLUMN x INT",
		"TRUNCATE TABLE t",
		"CALL some_proc()",
	}
	for _, q := range cases {
		if err := ensureSelectOnly(q); err == nil {
			t.Errorf("query %q: expected error, got nil", q)
		}
	}
}

func TestEnsureSelectOnly_rejects_multi_statement(t *testing.T) {
	if err := ensureSelectOnly("SELECT 1; DROP TABLE x;"); err == nil {
		t.Error("expected error for multi-statement query, got nil")
	}
}

// ---------------------------------------------------------------------------
// bindPositional
// ---------------------------------------------------------------------------

func TestBindPositional_named_to_question_mark(t *testing.T) {
	q, args, err := bindPositional("SELECT * FROM t WHERE a=:foo AND b=:bar", map[string]interface{}{"foo": 1, "bar": 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q != "SELECT * FROM t WHERE a=? AND b=?" {
		t.Errorf("unexpected query: %s", q)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != 2 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestBindPositional_repeated_placeholder(t *testing.T) {
	q, args, err := bindPositional("SELECT :x + :x", map[string]interface{}{"x": 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q != "SELECT ? + ?" {
		t.Errorf("unexpected query: %s", q)
	}
	if len(args) != 2 || args[0] != 5 || args[1] != 5 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestBindPositional_missing_param(t *testing.T) {
	if _, _, err := bindPositional("SELECT :missing", map[string]interface{}{}); err == nil {
		t.Error("expected error for missing param, got nil")
	}
}

// ---------------------------------------------------------------------------
// ensureLimit
// ---------------------------------------------------------------------------

func TestEnsureLimit_appends_when_missing(t *testing.T) {
	got := ensureLimit("SELECT *", 1000)
	if got != "SELECT * LIMIT 1000" {
		t.Errorf("unexpected query: %s", got)
	}
}

func TestEnsureLimit_respects_smaller_user_limit(t *testing.T) {
	got := ensureLimit("SELECT * LIMIT 10", 1000)
	if got != "SELECT * LIMIT 10" {
		t.Errorf("unexpected query: %s", got)
	}
}

func TestEnsureLimit_caps_larger_user_limit(t *testing.T) {
	got := ensureLimit("SELECT * LIMIT 99999", 10000)
	if got != "SELECT * LIMIT 10000" {
		t.Errorf("unexpected query: %s", got)
	}
}

// ---------------------------------------------------------------------------
// coerce
// ---------------------------------------------------------------------------

func TestCoerce_TINYINT_1_becomes_bool(t *testing.T) {
	got := coerce(int64(1), columnMeta{dbType: "TINYINT", length: 1, hasLength: true})
	if got != true {
		t.Errorf("expected true, got %v (%T)", got, got)
	}
}

func TestCoerce_JSON_parses(t *testing.T) {
	got := coerce([]byte(`{"a":1}`), columnMeta{dbType: "JSON"})
	m, ok := got.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map[string]interface{}, got %T", got)
	}
	if m["a"] != float64(1) {
		t.Errorf("expected a=1, got %v", m["a"])
	}
}

func TestCoerce_time_RFC3339(t *testing.T) {
	ts := time.Date(2026, 7, 6, 8, 0, 0, 0, time.UTC)
	got := coerce(ts, columnMeta{dbType: "DATETIME"})
	want := "2026-07-06T08:00:00Z"
	if got != want {
		t.Errorf("expected %s, got %v", want, got)
	}
}

func TestCoerce_DECIMAL_stays_string(t *testing.T) {
	got := coerce([]byte("123.45"), columnMeta{dbType: "DECIMAL"})
	if got != "123.45" {
		t.Errorf("expected \"123.45\", got %v (%T)", got, got)
	}
}

func TestCoerce_null(t *testing.T) {
	got := coerce(nil, columnMeta{dbType: "VARCHAR"})
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// applyFieldMap / applyValueMap
// ---------------------------------------------------------------------------

func TestApplyFieldMap_simpleAndEnum(t *testing.T) {
	rows := []map[string]interface{}{
		{"wo_no": "WO-1", "st": "R"},
	}
	fieldMap := map[string]interface{}{
		"of_number": "wo_no",
		"status": map[string]interface{}{
			"from": "st",
			"values": map[string]interface{}{
				"R": "RUNNING",
			},
		},
	}
	got, err := applyFieldMap(rows, fieldMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0]["of_number"] != "WO-1" {
		t.Errorf("expected of_number WO-1, got %v", got[0]["of_number"])
	}
	if got[0]["status"] != "RUNNING" {
		t.Errorf("expected status RUNNING, got %v", got[0]["status"])
	}
}

func TestApplyValueMap_passesThroughUnmapped(t *testing.T) {
	got := applyValueMap("X", map[string]interface{}{"R": "RUNNING"})
	if got != "X" {
		t.Errorf("expected raw value X to pass through unchanged, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Execute — end-to-end against an in-memory SQLite DB (no docker, no network)
// ---------------------------------------------------------------------------

type fakeGetter struct{ db *sql.DB }

func (f fakeGetter) Get(id string) (*sql.DB, error) {
	if id != "test_conn" {
		return nil, fmt.Errorf("unknown connection %q", id)
	}
	return f.db, nil
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec(`CREATE TABLE work_orders (wo_no TEXT, st TEXT, qty INTEGER)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO work_orders (wo_no, st, qty) VALUES (?, ?, ?)`, "WO-1", "R", 42); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return db
}

func TestExecute_rowsPath(t *testing.T) {
	h := NewSQLQueryHandler(fakeGetter{db: newTestDB(t)})
	out, err := h.Execute(context.Background(), map[string]interface{}{
		"connection_id": "test_conn",
		"query":         "SELECT wo_no, st, qty FROM work_orders",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	result := out.(map[string]interface{})
	rows := result["rows"].([]map[string]interface{})
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["wo_no"] != "WO-1" {
		t.Errorf("expected wo_no WO-1, got %v", rows[0]["wo_no"])
	}
	if result["canonical_type"] != nil {
		t.Errorf("expected nil canonical_type without field_map, got %v", result["canonical_type"])
	}
}

func TestExecute_canonicalPath(t *testing.T) {
	h := NewSQLQueryHandler(fakeGetter{db: newTestDB(t)})
	out, err := h.Execute(context.Background(), map[string]interface{}{
		"connection_id": "test_conn",
		"query":         "SELECT wo_no, st, qty FROM work_orders",
		"canonical":     "work_order",
		"field_map": map[string]interface{}{
			"of_number":  "wo_no",
			"actual_qty": "qty",
			"status": map[string]interface{}{
				"from": "st",
				"values": map[string]interface{}{
					"R": "RUNNING",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	result := out.(map[string]interface{})
	canonical := result["canonical"].([]map[string]interface{})
	if len(canonical) != 1 {
		t.Fatalf("expected 1 canonical row, got %d", len(canonical))
	}
	if canonical[0]["of_number"] != "WO-1" {
		t.Errorf("expected of_number WO-1, got %v", canonical[0]["of_number"])
	}
	if canonical[0]["status"] != "RUNNING" {
		t.Errorf("expected status RUNNING, got %v", canonical[0]["status"])
	}
	if result["canonical_type"] != "work_order" {
		t.Errorf("expected canonical_type work_order, got %v", result["canonical_type"])
	}
}

func TestExecute_rejectsNonSelect(t *testing.T) {
	h := NewSQLQueryHandler(fakeGetter{db: newTestDB(t)})
	_, err := h.Execute(context.Background(), map[string]interface{}{
		"connection_id": "test_conn",
		"query":         "DELETE FROM work_orders",
	})
	if err == nil {
		t.Error("expected error for non-SELECT query, got nil")
	}
}

func TestExecute_missingConnectionID(t *testing.T) {
	h := NewSQLQueryHandler(fakeGetter{db: newTestDB(t)})
	_, err := h.Execute(context.Background(), map[string]interface{}{
		"query": "SELECT 1",
	})
	if err == nil {
		t.Error("expected error for missing connection_id, got nil")
	}
}
