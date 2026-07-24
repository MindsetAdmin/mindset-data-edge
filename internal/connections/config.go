// internal/connections/config.go
package connections

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the parsed shape of config/connections.yaml.
type Config struct {
	Connections []ConnectionConfig `yaml:"connections"`
}

// ConnectionConfig describes one SQL connection MindSet Data can query.
// V1a supports only Driver == "mysql". JSON tags mirror the YAML tags
// (snake_case) so the same shape round-trips through the /api/connections
// REST endpoints (Day 6) without a separate DTO.
type ConnectionConfig struct {
	ID       string `yaml:"id" json:"id"`
	Name     string `yaml:"name" json:"name"`
	Driver   string `yaml:"driver" json:"driver"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Database string `yaml:"database" json:"database"`
	Username string `yaml:"username" json:"username"`
	// PasswordEnv names the environment variable holding the password.
	// Never inlined in the YAML or returned by the REST API.
	PasswordEnv string `yaml:"password_env" json:"password_env"`
	// TLS is "true", "false", or "skip-verify" (self-signed certs, dev only).
	TLS                    string `yaml:"tls" json:"tls"`
	ReadTimeoutSeconds     int    `yaml:"read_timeout_seconds" json:"read_timeout_seconds"`
	WriteTimeoutSeconds    int    `yaml:"write_timeout_seconds" json:"write_timeout_seconds"`
	MaxOpenConns           int    `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns           int    `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_seconds" json:"conn_max_lifetime_seconds"`
}

// LoadConfig reads and validates config/connections.yaml, applying defaults
// for any pool/timeout field left at zero.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	for i := range cfg.Connections {
		cfg.Connections[i].applyDefaults()
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *ConnectionConfig) applyDefaults() {
	if c.TLS == "" {
		c.TLS = "false"
	}
	if c.ReadTimeoutSeconds == 0 {
		c.ReadTimeoutSeconds = 30
	}
	if c.WriteTimeoutSeconds == 0 {
		c.WriteTimeoutSeconds = 10
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 5
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 2
	}
	if c.ConnMaxLifetimeSeconds == 0 {
		c.ConnMaxLifetimeSeconds = 300
	}
}

func (c *Config) validate() error {
	seen := make(map[string]bool, len(c.Connections))
	for _, conn := range c.Connections {
		if seen[conn.ID] {
			return fmt.Errorf("connections %q: duplicate id", conn.ID)
		}
		seen[conn.ID] = true
		if err := validateConnection(conn); err != nil {
			return err
		}
	}
	return nil
}

// validateConnection checks one entry in isolation (no duplicate-id check —
// that's only meaningful across a whole Config). Shared by Config.validate
// and Registry.Add so YAML-loaded and REST-created connections follow the
// same rules.
func validateConnection(conn ConnectionConfig) error {
	if conn.ID == "" {
		return fmt.Errorf("connections: entry with empty id")
	}
	if conn.Driver != "mysql" {
		return fmt.Errorf("connections %q: unsupported driver %q (V1a supports only mysql)", conn.ID, conn.Driver)
	}
	if conn.PasswordEnv == "" {
		return fmt.Errorf("connections %q: password_env is required (never inline a password)", conn.ID)
	}
	switch conn.TLS {
	case "true", "false", "skip-verify":
	default:
		return fmt.Errorf("connections %q: tls must be \"true\", \"false\", or \"skip-verify\", got %q", conn.ID, conn.TLS)
	}
	return nil
}

// Get returns the connection config with the given id.
func (c *Config) Get(id string) (ConnectionConfig, bool) {
	for _, conn := range c.Connections {
		if conn.ID == id {
			return conn, true
		}
	}
	return ConnectionConfig{}, false
}
