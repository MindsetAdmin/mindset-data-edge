// internal/connections/registry.go
package connections

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type entry struct {
	db       *sql.DB
	readOnly bool
}

// Registry holds one pooled *sql.DB per known connection id, opened lazily
// on first Get and cached for reuse. Connections known to the registry come
// from the YAML config it's constructed with, plus whatever Add/Remove adds
// or drops at runtime (e.g. via the /api/connections REST endpoints).
type Registry struct {
	// dial and verify are swapped out in tests so pool-caching logic can be
	// exercised without a live MySQL server.
	dial   func(ConnectionConfig) (*sql.DB, error)
	verify func(*sql.DB) (bool, error)

	mu    sync.Mutex
	conns map[string]ConnectionConfig
	dbs   map[string]*entry
}

// NewRegistry creates a Registry backed by the given connection definitions.
func NewRegistry(cfg *Config) *Registry {
	conns := make(map[string]ConnectionConfig, len(cfg.Connections))
	for _, c := range cfg.Connections {
		conns[c.ID] = c
	}
	return &Registry{
		dial:   dialMySQL,
		verify: VerifyReadOnly,
		conns:  conns,
		dbs:    make(map[string]*entry),
	}
}

// Add registers or replaces a connection definition, applying defaults and
// validating it first. If a pool already exists for this id (the config
// changed), it's closed so the next Get reopens with the new settings.
func (r *Registry) Add(cfg ConnectionConfig) error {
	cfg.applyDefaults()
	if err := validateConnection(cfg); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[cfg.ID] = cfg
	if e, ok := r.dbs[cfg.ID]; ok {
		e.db.Close()
		delete(r.dbs, cfg.ID)
	}
	return nil
}

// Remove forgets a connection definition and closes its pool if open.
func (r *Registry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, id)
	if e, ok := r.dbs[id]; ok {
		e.db.Close()
		delete(r.dbs, id)
	}
}

// List returns every known connection definition, ordered by id.
func (r *Registry) List() []ConnectionConfig {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]ConnectionConfig, 0, len(r.conns))
	for _, c := range r.conns {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get returns the pooled *sql.DB for the given connection id, opening and
// health-checking it on first use. Subsequent calls return the cached pool
// without re-verifying — use Test to force a fresh check.
func (r *Registry) Get(id string) (*sql.DB, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if e, ok := r.dbs[id]; ok {
		return e.db, nil
	}

	connCfg, ok := r.conns[id]
	if !ok {
		return nil, fmt.Errorf("connections: unknown connection id %q", id)
	}

	db, err := r.dial(connCfg)
	if err != nil {
		return nil, fmt.Errorf("connections %q: %w", id, err)
	}

	readOnly, err := r.verify(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("connections %q: %w", id, err)
	}

	r.dbs[id] = &entry{db: db, readOnly: readOnly}
	return db, nil
}

// Test re-runs the health check against connection id, reusing the pool if
// already open (or opening it if not). Unlike Get, it always re-verifies
// rather than trusting a cached result — this is what backs the "Test"
// button in the Pipeline Studio.
func (r *Registry) Test(id string) (readOnly bool, err error) {
	r.mu.Lock()
	connCfg, ok := r.conns[id]
	existing, hasExisting := r.dbs[id]
	r.mu.Unlock()
	if !ok {
		return false, fmt.Errorf("connections: unknown connection id %q", id)
	}

	var db *sql.DB
	if hasExisting {
		db = existing.db
	} else {
		db, err = r.dial(connCfg)
		if err != nil {
			return false, fmt.Errorf("connections %q: %w", id, err)
		}
	}

	readOnly, err = r.verify(db)
	if err != nil {
		if !hasExisting {
			db.Close()
		}
		return false, fmt.Errorf("connections %q: %w", id, err)
	}

	r.mu.Lock()
	r.dbs[id] = &entry{db: db, readOnly: readOnly}
	r.mu.Unlock()

	return readOnly, nil
}

// ReadOnly reports whether the connection's account was verified read-only
// the last time it was opened or tested. known is false if the connection
// hasn't been opened yet.
func (r *Registry) ReadOnly(id string) (readOnly bool, known bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.dbs[id]
	if !ok {
		return false, false
	}
	return e.readOnly, true
}

// CloseAll closes every pooled connection. Call once on shutdown.
func (r *Registry) CloseAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, e := range r.dbs {
		if err := e.db.Close(); err != nil {
			log.Printf("[connections] close %s: %v", id, err)
		}
	}
	r.dbs = make(map[string]*entry)
}

// dialMySQL opens a connection pool for cfg, resolving the password from its
// configured environment variable and applying the pool-size/lifetime limits.
func dialMySQL(cfg ConnectionConfig) (*sql.DB, error) {
	password := os.Getenv(cfg.PasswordEnv)
	if password == "" {
		return nil, fmt.Errorf("env var %s is empty or unset", cfg.PasswordEnv)
	}

	db, err := sql.Open("mysql", BuildMySQLDSN(cfg, password))
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)

	return db, nil
}
