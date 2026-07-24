// internal/connections/dsn.go
package connections

import "fmt"

// BuildMySQLDSN builds a go-sql-driver/mysql DSN from a connection config
// and a resolved password. The password is passed in rather than read from
// the environment here, so this stays a pure function that's trivial to
// unit test.
func BuildMySQLDSN(cfg ConnectionConfig, password string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC&charset=utf8mb4&interpolateParams=false&readTimeout=%ds&writeTimeout=%ds&tls=%s",
		cfg.Username, password, cfg.Host, cfg.Port, cfg.Database,
		cfg.ReadTimeoutSeconds, cfg.WriteTimeoutSeconds, cfg.TLS,
	)
}
