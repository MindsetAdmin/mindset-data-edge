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

	// Créer les tables
	if err := initTables(db); err != nil {
		return nil, err
	}

	log.Printf("[STORAGE] SQLite initialized at %s", dbPath)
	return &SQLiteStore{db: db}, nil
}

// initTables crée le schéma
func initTables(db *sql.DB) error {
	queries := []string{
		// Table des nœuds du Knowledge Graph
		`CREATE TABLE IF NOT EXISTS kg_nodes (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			label TEXT NOT NULL,
			properties TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Table des relations (edges)
		`CREATE TABLE IF NOT EXISTS kg_edges (
			id TEXT PRIMARY KEY,
			from_id TEXT NOT NULL,
			to_id TEXT NOT NULL,
			relation TEXT NOT NULL,
			weight REAL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Table des événements
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

		// Index pour les recherches rapides
		`CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_events_work_center ON events(work_center)`,
		`CREATE INDEX IF NOT EXISTS idx_kg_nodes_type ON kg_nodes(type)`,
		`CREATE INDEX IF NOT EXISTS idx_kg_edges_relation ON kg_edges(relation)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("[STORAGE] Error creating table: %v", err)
			return err
		}
	}
	return nil
}

// Close ferme la base
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// DB retourne la connexion SQLite
func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}