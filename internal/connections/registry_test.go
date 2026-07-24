// internal/connections/registry_test.go
package connections

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func writeTempConfig(t *testing.T, yamlContent string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "connections.yaml")
	if err := os.WriteFile(path, []byte(yamlContent), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestLoadConfig_appliesDefaultsAndParses(t *testing.T) {
	path := writeTempConfig(t, `
connections:
  - id: dev_erp
    name: "Dev ERP"
    driver: mysql
    host: localhost
    port: 3307
    database: fake_erp
    username: mindset_readonly
    password_env: MINDSET_ERP_PASSWORD
    tls: "false"
`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	conn, ok := cfg.Get("dev_erp")
	if !ok {
		t.Fatalf("expected connection %q to be found", "dev_erp")
	}
	if conn.ReadTimeoutSeconds != 30 || conn.WriteTimeoutSeconds != 10 {
		t.Errorf("expected default timeouts 30/10, got %d/%d", conn.ReadTimeoutSeconds, conn.WriteTimeoutSeconds)
	}
	if conn.MaxOpenConns != 5 || conn.MaxIdleConns != 2 || conn.ConnMaxLifetimeSeconds != 300 {
		t.Errorf("expected default pool limits 5/2/300, got %d/%d/%d", conn.MaxOpenConns, conn.MaxIdleConns, conn.ConnMaxLifetimeSeconds)
	}
}

func TestLoadConfig_rejectsDuplicateID(t *testing.T) {
	path := writeTempConfig(t, `
connections:
  - id: dup
    driver: mysql
    password_env: X
  - id: dup
    driver: mysql
    password_env: X
`)
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected error for duplicate id, got nil")
	}
}

func TestLoadConfig_rejectsNonMySQLDriver(t *testing.T) {
	path := writeTempConfig(t, `
connections:
  - id: pg
    driver: postgres
    password_env: X
`)
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected error for unsupported driver, got nil")
	}
}

func TestLoadConfig_rejectsMissingPasswordEnv(t *testing.T) {
	path := writeTempConfig(t, `
connections:
  - id: noenv
    driver: mysql
`)
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected error for missing password_env, got nil")
	}
}

func TestBuildMySQLDSN_format(t *testing.T) {
	cfg := ConnectionConfig{
		Username: "mindset_readonly", Host: "10.42.1.15", Port: 3306, Database: "erp_prod",
		ReadTimeoutSeconds: 30, WriteTimeoutSeconds: 10, TLS: "true",
	}
	got := BuildMySQLDSN(cfg, "s3cret")
	want := "mindset_readonly:s3cret@tcp(10.42.1.15:3306)/erp_prod?parseTime=true&loc=UTC&charset=utf8mb4&interpolateParams=false&readTimeout=30s&writeTimeout=10s&tls=true"
	if got != want {
		t.Errorf("DSN mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestRegistry_pool_reuse(t *testing.T) {
	cfg := &Config{Connections: []ConnectionConfig{{ID: "fake", Driver: "mysql", PasswordEnv: "UNUSED"}}}
	reg := NewRegistry(cfg)

	opened := 0
	reg.dial = func(ConnectionConfig) (*sql.DB, error) {
		opened++
		return sql.Open("sqlite", ":memory:")
	}
	reg.verify = func(*sql.DB) (bool, error) { return true, nil }

	db1, err := reg.Get("fake")
	if err != nil {
		t.Fatalf("first Get: %v", err)
	}
	db2, err := reg.Get("fake")
	if err != nil {
		t.Fatalf("second Get: %v", err)
	}
	if db1 != db2 {
		t.Error("expected the same *sql.DB on repeated Get calls")
	}
	if opened != 1 {
		t.Errorf("expected dial to run once, ran %d times", opened)
	}

	readOnly, known := reg.ReadOnly("fake")
	if !known || !readOnly {
		t.Errorf("expected ReadOnly to report known=true readOnly=true, got known=%v readOnly=%v", known, readOnly)
	}

	reg.CloseAll()
	if _, known := reg.ReadOnly("fake"); known {
		t.Error("expected ReadOnly to be unknown after CloseAll")
	}
}

func TestRegistry_unknownID(t *testing.T) {
	reg := NewRegistry(&Config{})
	if _, err := reg.Get("nope"); err == nil {
		t.Fatal("expected error for unknown connection id")
	}
}

func TestRegistry_AddListRemove(t *testing.T) {
	reg := NewRegistry(&Config{})
	reg.dial = func(ConnectionConfig) (*sql.DB, error) { return sql.Open("sqlite", ":memory:") }
	reg.verify = func(*sql.DB) (bool, error) { return true, nil }

	if err := reg.Add(ConnectionConfig{ID: "a", Driver: "mysql", PasswordEnv: "X"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if got := len(reg.List()); got != 1 {
		t.Fatalf("expected 1 connection after Add, got %d", got)
	}
	if _, err := reg.Get("a"); err != nil {
		t.Fatalf("Get after Add: %v", err)
	}

	reg.Remove("a")
	if got := len(reg.List()); got != 0 {
		t.Errorf("expected 0 connections after Remove, got %d", got)
	}
	if _, err := reg.Get("a"); err == nil {
		t.Error("expected error getting a removed connection")
	}
}

func TestRegistry_Add_rejectsInvalid(t *testing.T) {
	reg := NewRegistry(&Config{})
	if err := reg.Add(ConnectionConfig{ID: "bad", Driver: "postgres", PasswordEnv: "X"}); err == nil {
		t.Error("expected error adding a non-mysql connection in V1a")
	}
}

func TestRegistry_Add_closesExistingPoolOnReplace(t *testing.T) {
	reg := NewRegistry(&Config{})
	dialed := 0
	reg.dial = func(ConnectionConfig) (*sql.DB, error) {
		dialed++
		return sql.Open("sqlite", ":memory:")
	}
	reg.verify = func(*sql.DB) (bool, error) { return true, nil }

	if err := reg.Add(ConnectionConfig{ID: "a", Driver: "mysql", PasswordEnv: "X"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if _, err := reg.Get("a"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	// Re-Add the same id — should drop the cached pool so the next Get re-dials.
	if err := reg.Add(ConnectionConfig{ID: "a", Driver: "mysql", PasswordEnv: "X", Host: "changed"}); err != nil {
		t.Fatalf("re-Add: %v", err)
	}
	if _, err := reg.Get("a"); err != nil {
		t.Fatalf("Get after re-Add: %v", err)
	}
	if dialed != 2 {
		t.Errorf("expected dial to run twice (pool reset on re-Add), ran %d times", dialed)
	}
}

func TestRegistry_Test_reverifies(t *testing.T) {
	cfg := &Config{Connections: []ConnectionConfig{{ID: "fake", Driver: "mysql", PasswordEnv: "UNUSED"}}}
	reg := NewRegistry(cfg)

	dialed := 0
	reg.dial = func(ConnectionConfig) (*sql.DB, error) {
		dialed++
		return sql.Open("sqlite", ":memory:")
	}
	verified := 0
	reg.verify = func(*sql.DB) (bool, error) {
		verified++
		return verified%2 == 0, nil // alternates, to prove each Test call is fresh
	}

	if _, err := reg.Test("fake"); err != nil {
		t.Fatalf("first Test: %v", err)
	}
	if _, err := reg.Test("fake"); err != nil {
		t.Fatalf("second Test: %v", err)
	}
	if dialed != 1 {
		t.Errorf("expected dial to run once (pool reused across Test calls), ran %d times", dialed)
	}
	if verified != 2 {
		t.Errorf("expected verify to run on every Test call, ran %d times", verified)
	}
}
