//go:build integration

// internal/e2e/sql_pipeline_test.go
//
// Integration tests against a real, disposable MySQL container
// (testcontainers-go). These are NOT run by a plain `go test ./...` — they
// need Docker, so they're gated behind the "integration" build tag:
//
//	go test -tags=integration ./internal/e2e/...
//
// Per docs/mysql_connector.md §12 "Integration (testcontainers)".
package e2e

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/MindsetAdmin/mindset-data-edge/internal/connections"
	"github.com/MindsetAdmin/mindset-data-edge/internal/functions/connectors"
)

const (
	writerUser     = "mindset_writer_test"
	writerPassword = "writer_test_pw"
	readonlyUser   = "mindset_readonly_test"
	readonlyPass   = "readonly_test_pw"
	testDatabase   = "mindset_test"
)

// initSQL seeds a work_orders table (a tiny slice of sim/erp/schema.mysql.sql)
// and grants a genuinely read-only user, so TestHealthCheck_readonly_user
// exercises the same read-only/writable distinction as the real dev-erp stack
// (sim/erp/grant.mysql.sql).
const initSQL = `
CREATE TABLE work_orders (
  of_number   VARCHAR(64) PRIMARY KEY,
  work_center VARCHAR(64) NOT NULL,
  status      VARCHAR(16) NOT NULL
);
INSERT INTO work_orders (of_number, work_center, status) VALUES
  ('WO-1', 'machine1', 'RUNNING'),
  ('WO-2', 'machine2', 'PLANNED');

CREATE USER '` + readonlyUser + `'@'%' IDENTIFIED BY '` + readonlyPass + `';
GRANT SELECT ON ` + testDatabase + `.* TO '` + readonlyUser + `'@'%';
FLUSH PRIVILEGES;
`

// testEnv holds everything a test needs against one shared MySQL container:
// a writer connection (used for the happy-path/timeout/injection tests, which
// only need SELECT) and a read-only connection (for the health-check test).
type testEnv struct {
	writerConn   connections.ConnectionConfig
	readonlyConn connections.ConnectionConfig
	registry     *connections.Registry
	handler      *connectors.SQLQueryHandler
}

// setupMySQL starts a disposable MySQL 8 container. If Docker isn't
// reachable, it skips (not fails) the test — these tests are meant to be run
// deliberately (go test -tags=integration), and a missing Docker daemon in a
// given environment isn't a code failure.
func setupMySQL(t *testing.T) *testEnv {
	t.Helper()
	ctx := context.Background()

	initScript := writeInitScript(t)

	ctr, err := mysql.Run(ctx, "mysql:8",
		mysql.WithDatabase(testDatabase),
		mysql.WithUsername(writerUser),
		mysql.WithPassword(writerPassword),
		mysql.WithScripts(initScript),
	)
	if err != nil {
		t.Skipf("skipping — could not start MySQL testcontainer (is Docker running?): %v", err)
	}
	t.Cleanup(func() {
		if err := ctr.Terminate(context.Background()); err != nil {
			t.Logf("terminate container: %v", err)
		}
	})

	host, err := ctr.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}
	mappedPort, err := ctr.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatalf("container port: %v", err)
	}
	port, err := strconv.Atoi(mappedPort.Port())
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}

	base := connections.ConnectionConfig{
		ID:                     "e2e_writer",
		Name:                   "e2e writer",
		Driver:                 "mysql",
		Host:                   host,
		Port:                   port,
		Database:               testDatabase,
		Username:               writerUser,
		PasswordEnv:            "MINDSET_E2E_WRITER_PASSWORD",
		TLS:                    "false",
		ReadTimeoutSeconds:     30,
		WriteTimeoutSeconds:    10,
		MaxOpenConns:           5,
		MaxIdleConns:           2,
		ConnMaxLifetimeSeconds: 300,
	}
	readonly := base
	readonly.ID = "e2e_readonly"
	readonly.Name = "e2e readonly"
	readonly.Username = readonlyUser
	readonly.PasswordEnv = "MINDSET_E2E_READONLY_PASSWORD"

	t.Setenv("MINDSET_E2E_WRITER_PASSWORD", writerPassword)
	t.Setenv("MINDSET_E2E_READONLY_PASSWORD", readonlyPass)

	registry := connections.NewRegistry(&connections.Config{
		Connections: []connections.ConnectionConfig{base, readonly},
	})
	t.Cleanup(registry.CloseAll)

	// Fail fast (with a clear message) rather than let every test time out
	// individually if the container came up but isn't actually reachable yet.
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := waitForConnection(waitCtx, registry, base.ID); err != nil {
		t.Fatalf("mysql container did not become ready: %v", err)
	}

	return &testEnv{
		writerConn:   base,
		readonlyConn: readonly,
		registry:     registry,
		handler:      connectors.NewSQLQueryHandler(registry),
	}
}

func waitForConnection(ctx context.Context, registry *connections.Registry, id string) error {
	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return lastErr
		default:
		}
		if _, err := registry.Get(id); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func writeInitScript(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "init.sql")
	if err := os.WriteFile(path, []byte(initSQL), 0o644); err != nil {
		t.Fatalf("write init script: %v", err)
	}
	return path
}

// ---------------------------------------------------------------------------
// TestHappyPath — real MySQL + seed -> SELECT returns rows with correct types
// ---------------------------------------------------------------------------

func TestHappyPath(t *testing.T) {
	env := setupMySQL(t)

	out, err := env.handler.Execute(context.Background(), map[string]interface{}{
		"connection_id": env.writerConn.ID,
		"query":         "SELECT of_number, work_center, status FROM work_orders WHERE work_center = :wc",
		"params":        map[string]interface{}{"wc": "machine1"},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	result := out.(map[string]interface{})
	rows := result["rows"].([]map[string]interface{})
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["of_number"] != "WO-1" || rows[0]["status"] != "RUNNING" {
		t.Errorf("unexpected row: %+v", rows[0])
	}
}

// ---------------------------------------------------------------------------
// TestTimeoutKicksIn — SELECT SLEEP(5) with a 1s timeout -> error
// ---------------------------------------------------------------------------

func TestTimeoutKicksIn(t *testing.T) {
	env := setupMySQL(t)

	started := time.Now()
	_, err := env.handler.Execute(context.Background(), map[string]interface{}{
		"connection_id":   env.writerConn.ID,
		"query":           "SELECT SLEEP(5)",
		"timeout_seconds": float64(1),
	})
	elapsed := time.Since(started)

	if err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
	if elapsed > 4*time.Second {
		t.Errorf("expected the 1s timeout to cut this short, took %s", elapsed)
	}
}

// ---------------------------------------------------------------------------
// TestInjectionAttempt — a parameter value containing SQL is bound as data,
// never executed as a second statement.
// ---------------------------------------------------------------------------

func TestInjectionAttempt(t *testing.T) {
	env := setupMySQL(t)

	out, err := env.handler.Execute(context.Background(), map[string]interface{}{
		"connection_id": env.writerConn.ID,
		"query":         "SELECT of_number FROM work_orders WHERE of_number = :id",
		"params":        map[string]interface{}{"id": "WO-1; DROP TABLE work_orders"},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	result := out.(map[string]interface{})
	rows := result["rows"].([]map[string]interface{})
	if len(rows) != 0 {
		t.Errorf("expected 0 rows (no of_number equals the literal injected string), got %d", len(rows))
	}

	// The table must still exist and still hold both seeded rows.
	db, err := env.registry.Get(env.writerConn.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM work_orders").Scan(&count); err != nil {
		t.Fatalf("work_orders table appears to be gone: %v", err)
	}
	if count != 2 {
		t.Errorf("expected work_orders to still have 2 rows, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// TestReadOnlyEnforcement — a non-SELECT query is rejected before the
// connection is even resolved.
// ---------------------------------------------------------------------------

func TestReadOnlyEnforcement(t *testing.T) {
	env := setupMySQL(t)

	_, err := env.handler.Execute(context.Background(), map[string]interface{}{
		// A connection id that doesn't exist — if the guard fires before the
		// registry lookup (as designed), the error is the SELECT-only guard,
		// never "unknown connection id".
		"connection_id": "does_not_exist",
		"query":         "INSERT INTO work_orders (of_number) VALUES ('WO-X')",
	})
	if err == nil {
		t.Fatal("expected an error for a non-SELECT query, got nil")
	}
	if got := err.Error(); !strings.Contains(got, "only SELECT statements are allowed") {
		t.Errorf("expected the SELECT-only guard to fire before connection resolution, got: %v", got)
	}
}

// ---------------------------------------------------------------------------
// TestHealthCheck_readonly_user — mindset_readonly_test cannot CREATE
// TEMPORARY TABLE, so VerifyReadOnly (via Registry.Test) reports read_only=true.
// ---------------------------------------------------------------------------

func TestHealthCheck_readonly_user(t *testing.T) {
	env := setupMySQL(t)

	readOnly, err := env.registry.Test(env.readonlyConn.ID)
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if !readOnly {
		t.Error("expected the read-only user to be reported read_only=true")
	}

	// Sanity check the writer user is NOT reported read-only, so this test
	// is actually distinguishing the two, not just always returning true.
	writerReadOnly, err := env.registry.Test(env.writerConn.ID)
	if err != nil {
		t.Fatalf("Test (writer): %v", err)
	}
	if writerReadOnly {
		t.Error("expected the writer user (full grants on the test DB) to be reported read_only=false")
	}
}
