// internal/connections/health.go
package connections

import (
	"database/sql"
	"fmt"
	"log"
)

// VerifyReadOnly pings db for reachability/credentials, then probes whether
// the account can write. A read-only account fails the CREATE — that's the
// expected, healthy outcome. An account that CAN write is logged as a
// warning, not refused: enterprise IT sometimes provisions broader accounts
// than requested, and the connection should still work.
func VerifyReadOnly(db *sql.DB) (readOnly bool, err error) {
	if err := db.Ping(); err != nil {
		return false, fmt.Errorf("ping: %w", err)
	}

	if _, err := db.Exec("CREATE TEMPORARY TABLE mindset_probe (id INT)"); err != nil {
		return true, nil
	}

	log.Printf("[connections] WARNING: account has write privileges beyond SELECT — a read-only account is recommended")
	if _, err := db.Exec("DROP TEMPORARY TABLE mindset_probe"); err != nil {
		log.Printf("[connections] WARNING: could not drop probe table: %v", err)
	}
	return false, nil
}
