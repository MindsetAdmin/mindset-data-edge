// internal/storage/sqlite.go
package storage

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// SQLiteStore gère la base de données locale
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore crée une nouvelle connexion SQLite
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// modernc.org/sqlite utilise "sqlite" comme driver (pas "sqlite3")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Tester la connexion
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// cmd/server and cmd/agent each open their own *sql.DB against this same
	// file (server: TagRegistry + its own KG; agent: Rules Engine's KG
	// subscriber) — concurrent writers from separate OS processes hit
	// SQLITE_BUSY immediately without this, since SQLite's default is to fail
	// fast rather than wait for a lock to clear (Entry 130/131: this surfaced
	// once the KG subscribers' client-ID collision was fixed and both
	// processes' subscribers started actually writing at the same time).
	if _, err := db.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
		log.Printf("[STORAGE] Failed to set busy_timeout: %v", err)
	}

	// Créer les tables
	if err := initTables(db); err != nil {
		return nil, err
	}

	log.Printf("[STORAGE] SQLite initialized at %s", dbPath)
	return &SQLiteStore{db: db}, nil
}

// initTables crée le schéma. Unified KG (2026-07-02 refactor): nodes + edges
// carry a `category` column ("business" for site-fingerprint entities like
// Equipment/Event/Cause/Cost; "platform" for pipeline-topology entities like
// Pipeline/Function/Topic/Connection/Dashboard). Old schemas auto-migrate.
//
// Sequencing matters: for pre-refactor DBs, CREATE TABLE IF NOT EXISTS is a
// no-op (the old schema without `category` is preserved). We must therefore
// ALTER TABLE ADD COLUMN BEFORE any CREATE INDEX that references `category`.
func initTables(db *sql.DB) error {
	// Step 1 — CREATE TABLE (new DBs get the unified schema; existing DBs keep their old one).
	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS kg_nodes (
			id TEXT PRIMARY KEY,
			category TEXT NOT NULL DEFAULT 'business',
			type TEXT NOT NULL,
			label TEXT NOT NULL,
			properties TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS kg_edges (
			id TEXT PRIMARY KEY,
			category TEXT NOT NULL DEFAULT 'business',
			from_id TEXT NOT NULL,
			to_id TEXT NOT NULL,
			relation TEXT NOT NULL,
			weight REAL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			work_center TEXT NOT NULL,
			duration_seconds REAL,
			cause TEXT,
			cost_eur REAL,
			payload TEXT,
			timestamp DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// connections — SQL connection definitions created via /api/connections
		// (internal/connections.Registry is the runtime pool; this table is the
		// persisted source of truth so they survive a restart). Never stores a
		// password — only the env var name that holds it.
		`CREATE TABLE IF NOT EXISTS connections (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			driver TEXT NOT NULL,
			host TEXT NOT NULL,
			port INTEGER NOT NULL,
			database TEXT NOT NULL,
			username TEXT NOT NULL,
			password_env TEXT NOT NULL,
			tls TEXT NOT NULL,
			read_timeout_seconds INTEGER NOT NULL DEFAULT 30,
			write_timeout_seconds INTEGER NOT NULL DEFAULT 10,
			max_open_conns INTEGER NOT NULL DEFAULT 5,
			max_idle_conns INTEGER NOT NULL DEFAULT 2,
			conn_max_lifetime_seconds INTEGER NOT NULL DEFAULT 300,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, q := range tableQueries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("[STORAGE] Error creating table: %v", err)
			return err
		}
	}

	// Step 2 — Migration: for pre-refactor DBs, add the `category` column BEFORE
	// creating indexes on it. SQLite's PRAGMA table_info() probe makes ALTER TABLE
	// idempotent-safe.
	if err := migrateAddCategoryColumn(db); err != nil {
		log.Printf("[STORAGE] Migration warning: %v", err)
		return err
	}

	// Step 3 — Indexes (now safe to reference `category` — column exists everywhere).
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_events_work_center ON events(work_center)`,
		`CREATE INDEX IF NOT EXISTS idx_kg_nodes_type ON kg_nodes(type)`,
		`CREATE INDEX IF NOT EXISTS idx_kg_nodes_category ON kg_nodes(category)`,
		`CREATE INDEX IF NOT EXISTS idx_kg_edges_relation ON kg_edges(relation)`,
		`CREATE INDEX IF NOT EXISTS idx_kg_edges_category ON kg_edges(category)`,
	}
	for _, q := range indexQueries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("[STORAGE] Error creating index: %v", err)
			return err
		}
	}
	return nil
}

// migrateAddCategoryColumn adds the `category` column to legacy kg_nodes and
// kg_edges tables. Existing rows default to 'business' which is safe because
// pre-refactor DBs only contained business (Domain) KG entries.
func migrateAddCategoryColumn(db *sql.DB) error {
	for _, table := range []string{"kg_nodes", "kg_edges"} {
		if hasCategoryColumn(db, table) {
			continue
		}
		q := "ALTER TABLE " + table + " ADD COLUMN category TEXT NOT NULL DEFAULT 'business'"
		if _, err := db.Exec(q); err != nil {
			log.Printf("[STORAGE] Failed to add category to %s: %v", table, err)
			return err
		}
		log.Printf("[STORAGE] Migrated %s: added category column", table)
	}
	return nil
}

func hasCategoryColumn(db *sql.DB, table string) bool {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			continue
		}
		if name == "category" {
			return true
		}
	}
	return false
}

// Close ferme la base
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// DB retourne la connexion SQLite
func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}