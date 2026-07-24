// internal/storage/connections.go
package storage

import (
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/connections"
)

// ConnectionRecord is a persisted connections table row: a ConnectionConfig
// plus the bookkeeping timestamps the REST API surfaces.
type ConnectionRecord struct {
	connections.ConnectionConfig
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ListConnections returns every persisted connection, ordered by id.
func (s *SQLiteStore) ListConnections() ([]ConnectionRecord, error) {
	rows, err := s.db.Query(`SELECT id, name, driver, host, port, database, username,
		password_env, tls, read_timeout_seconds, write_timeout_seconds,
		max_open_conns, max_idle_conns, conn_max_lifetime_seconds, created_at, updated_at
		FROM connections ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ConnectionRecord
	for rows.Next() {
		var r ConnectionRecord
		if err := rows.Scan(&r.ID, &r.Name, &r.Driver, &r.Host, &r.Port, &r.Database, &r.Username,
			&r.PasswordEnv, &r.TLS, &r.ReadTimeoutSeconds, &r.WriteTimeoutSeconds,
			&r.MaxOpenConns, &r.MaxIdleConns, &r.ConnMaxLifetimeSeconds, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpsertConnection inserts or replaces a connection row, refreshing updated_at.
func (s *SQLiteStore) UpsertConnection(cfg connections.ConnectionConfig) error {
	_, err := s.db.Exec(`INSERT INTO connections
		(id, name, driver, host, port, database, username, password_env, tls,
		 read_timeout_seconds, write_timeout_seconds, max_open_conns, max_idle_conns,
		 conn_max_lifetime_seconds, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name, driver=excluded.driver, host=excluded.host, port=excluded.port,
			database=excluded.database, username=excluded.username, password_env=excluded.password_env,
			tls=excluded.tls, read_timeout_seconds=excluded.read_timeout_seconds,
			write_timeout_seconds=excluded.write_timeout_seconds, max_open_conns=excluded.max_open_conns,
			max_idle_conns=excluded.max_idle_conns, conn_max_lifetime_seconds=excluded.conn_max_lifetime_seconds,
			updated_at=CURRENT_TIMESTAMP`,
		cfg.ID, cfg.Name, cfg.Driver, cfg.Host, cfg.Port, cfg.Database, cfg.Username, cfg.PasswordEnv, cfg.TLS,
		cfg.ReadTimeoutSeconds, cfg.WriteTimeoutSeconds, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetimeSeconds)
	return err
}

// DeleteConnection removes a connection row. Not an error if it doesn't exist.
func (s *SQLiteStore) DeleteConnection(id string) error {
	_, err := s.db.Exec(`DELETE FROM connections WHERE id = ?`, id)
	return err
}
