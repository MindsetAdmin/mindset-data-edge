// cmd/erpsim — ERP activity simulator against the fake_erp MySQL database.
//
// Runs 4 background loops that mimic a live mid-sized manufacturing ERP:
//
//	advance  — every  30s: nudge actual_qty on each RUNNING work order
//	rotate   — every   5m: 20% chance to finish a RUNNING OF + start the next PLANNED
//	quality  — every  10m: insert a new quality reading on each in-flight batch (10% out-of-spec)
//	plan     — every   1h: create a new PLANNED work order per work_center
//
// Uses the mindset_writer user (SELECT+INSERT+UPDATE, no DELETE) — the same
// boundary enforced in production for erpsim-shaped workloads.
//
// Configure via env vars (all optional):
//
//	ERPSIM_DSN     — full driver DSN (default: mindset_writer:writer_dev@tcp(localhost:3308)/fake_erp?parseTime=true&loc=UTC&charset=utf8mb4)
//	ERPSIM_TICK_ADVANCE  — advance loop interval  (default 30s)
//	ERPSIM_TICK_ROTATE   — rotate  loop interval  (default 5m)
//	ERPSIM_TICK_QUALITY  — quality loop interval  (default 10m)
//	ERPSIM_TICK_PLAN     — plan    loop interval  (default 1h)

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	products    = []string{"PROD-A01", "PROD-A02", "PROD-A03", "PROD-A04", "PROD-A05", "PROD-A06", "PROD-A07", "PROD-A08", "PROD-B01", "PROD-B02", "PROD-C01", "PROD-C02"}
	workCenters = []string{"machine1", "machine2", "machine3"}
	operators   = []string{"OP-001", "OP-002", "OP-003", "OP-004", "OP-005", "OP-006", "OP-007", "OP-008"}
	metrics     = []metricSpec{
		{"viscosity", 800, 900},
		{"temperature", 3.5, 4.5},
		{"ph", 4.2, 4.6},
		{"moisture", 20, 25},
	}

	nextOFCounter    int64 = 9100
	nextBatchCounter int64 = 6100
	counterMu        sync.Mutex
)

type metricSpec struct {
	name    string
	specMin float64
	specMax float64
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("[erpsim] starting")

	dsn := envOr("ERPSIM_DSN", "mindset_writer:writer_dev@tcp(localhost:3308)/fake_erp?parseTime=true&loc=UTC&charset=utf8mb4")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[erpsim] open: %v", err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[erpsim] ping: %v", err)
	}
	log.Println("[erpsim] connected to fake_erp")

	if err := initCounters(db); err != nil {
		log.Fatalf("[erpsim] init counters: %v", err)
	}
	log.Printf("[erpsim] counters initialized: next OF=%d, next batch=%d", nextOFCounter+1, nextBatchCounter+1)

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("[erpsim] shutdown requested")
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(4)
	go loop(ctx, &wg, "advance", envDurOr("ERPSIM_TICK_ADVANCE", 30*time.Second), db, advanceRunningOFs)
	go loop(ctx, &wg, "rotate", envDurOr("ERPSIM_TICK_ROTATE", 5*time.Minute), db, rotateOFs)
	go loop(ctx, &wg, "quality", envDurOr("ERPSIM_TICK_QUALITY", 10*time.Minute), db, addQualityResult)
	go loop(ctx, &wg, "plan", envDurOr("ERPSIM_TICK_PLAN", 1*time.Hour), db, planNewOF)

	wg.Wait()
	log.Println("[erpsim] stopped")
}

// ---------------------------------------------------------------------------
// Loop driver
// ---------------------------------------------------------------------------

func loop(ctx context.Context, wg *sync.WaitGroup, name string, every time.Duration, db *sql.DB, fn func(*sql.DB) error) {
	defer wg.Done()
	log.Printf("[erpsim/%s] tick=%s", name, every)
	// Run once immediately so the user sees action within seconds.
	if err := fn(db); err != nil {
		log.Printf("[erpsim/%s] %v", name, err)
	}
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if err := fn(db); err != nil {
				log.Printf("[erpsim/%s] %v", name, err)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Counters — pick up above existing seed
// ---------------------------------------------------------------------------

func initCounters(db *sql.DB) error {
	// OF counter — max existing WO-2026-* numeric suffix.
	rows, err := db.Query(`SELECT of_number FROM work_orders WHERE of_number LIKE 'WO-2026-%'`)
	if err != nil {
		return err
	}
	max := nextOFCounter
	for rows.Next() {
		var of string
		if err := rows.Scan(&of); err != nil {
			rows.Close()
			return err
		}
		if n, err := strconv.ParseInt(of[len("WO-2026-"):], 10, 64); err == nil && n > max {
			max = n
		}
	}
	rows.Close()
	nextOFCounter = max

	// Batch counter — max existing B-2026-* numeric suffix.
	rows2, err := db.Query(`SELECT batch_id FROM batches WHERE batch_id LIKE 'B-2026-%'`)
	if err != nil {
		return err
	}
	maxB := nextBatchCounter
	for rows2.Next() {
		var b string
		if err := rows2.Scan(&b); err != nil {
			rows2.Close()
			return err
		}
		if n, err := strconv.ParseInt(b[len("B-2026-"):], 10, 64); err == nil && n > maxB {
			maxB = n
		}
	}
	rows2.Close()
	nextBatchCounter = maxB
	return nil
}

func nextOF() string {
	counterMu.Lock()
	nextOFCounter++
	n := nextOFCounter
	counterMu.Unlock()
	return fmt.Sprintf("WO-2026-%d", n)
}

func nextBatch() string {
	counterMu.Lock()
	nextBatchCounter++
	n := nextBatchCounter
	counterMu.Unlock()
	return fmt.Sprintf("B-2026-%d", n)
}

// ---------------------------------------------------------------------------
// Move 1 — advanceRunningOFs
// ---------------------------------------------------------------------------

func advanceRunningOFs(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT w.of_number, w.actual_qty, w.planned_qty, COALESCE(p.target_rate, 1000)
		FROM work_orders w LEFT JOIN products p ON p.product_code = w.product_code
		WHERE w.status = 'RUNNING'`)
	if err != nil {
		return fmt.Errorf("query running: %w", err)
	}
	type row struct {
		of         string
		actual     int
		planned    int
		targetRate int
	}
	var running []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.of, &r.actual, &r.planned, &r.targetRate); err != nil {
			rows.Close()
			return err
		}
		running = append(running, r)
	}
	rows.Close()

	for _, r := range running {
		if r.actual >= r.planned {
			continue
		}
		// target_rate is units/hour; we tick every 30s → target_rate/120 per tick, ±20% jitter.
		base := float64(r.targetRate) / 120.0
		delta := int(base * (0.8 + 0.4*rand.Float64()))
		if delta < 1 {
			delta = 1
		}
		newQty := r.actual + delta
		if newQty > r.planned {
			newQty = r.planned
		}
		if _, err := db.Exec(`UPDATE work_orders SET actual_qty=? WHERE of_number=?`, newQty, r.of); err != nil {
			return fmt.Errorf("update %s: %w", r.of, err)
		}
		log.Printf("[erpsim/advance] %s: %d → %d (+%d, target %d)", r.of, r.actual, newQty, delta, r.planned)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Move 2 — rotateOFs (finish + start next planned)
// ---------------------------------------------------------------------------

func rotateOFs(db *sql.DB) error {
	rows, err := db.Query(`SELECT of_number, work_center FROM work_orders WHERE status='RUNNING'`)
	if err != nil {
		return fmt.Errorf("query running: %w", err)
	}
	type openOF struct{ of, wc string }
	var open []openOF
	for rows.Next() {
		var r openOF
		if err := rows.Scan(&r.of, &r.wc); err != nil {
			rows.Close()
			return err
		}
		open = append(open, r)
	}
	rows.Close()

	for _, r := range open {
		if rand.Float64() > 0.2 {
			continue // 80% skip
		}
		if _, err := db.Exec(`UPDATE work_orders SET status='DONE', finished_at=NOW() WHERE of_number=?`, r.of); err != nil {
			log.Printf("[erpsim/rotate] finish %s: %v", r.of, err)
			continue
		}
		batchStatus := "PASS"
		if rand.Float64() < 0.10 {
			batchStatus = "REWORK"
		}
		if _, err := db.Exec(
			`UPDATE batches SET finished_at=NOW(), quality_status=? WHERE of_number=? AND finished_at IS NULL`,
			batchStatus, r.of,
		); err != nil {
			log.Printf("[erpsim/rotate] finish batch for %s: %v", r.of, err)
		}
		log.Printf("[erpsim/rotate] %s → DONE (batch %s)", r.of, batchStatus)

		// Start next PLANNED OF for this work_center.
		var nextOFNum string
		err := db.QueryRow(
			`SELECT of_number FROM work_orders WHERE work_center=? AND status='PLANNED' ORDER BY of_number LIMIT 1`,
			r.wc,
		).Scan(&nextOFNum)
		if err == sql.ErrNoRows {
			log.Printf("[erpsim/rotate] no PLANNED OF on %s — nothing to start", r.wc)
			continue
		}
		if err != nil {
			log.Printf("[erpsim/rotate] find planned on %s: %v", r.wc, err)
			continue
		}
		op := operators[rand.Intn(len(operators))]
		if _, err := db.Exec(
			`UPDATE work_orders SET status='RUNNING', started_at=NOW(), operator_id=? WHERE of_number=?`,
			op, nextOFNum,
		); err != nil {
			log.Printf("[erpsim/rotate] start %s: %v", nextOFNum, err)
			continue
		}
		batchID := nextBatch()
		if _, err := db.Exec(
			`INSERT INTO batches (batch_id, of_number, started_at) VALUES (?, ?, NOW())`,
			batchID, nextOFNum,
		); err != nil {
			log.Printf("[erpsim/rotate] new batch %s for %s: %v", batchID, nextOFNum, err)
		}
		log.Printf("[erpsim/rotate] started %s on %s (op=%s, batch=%s)", nextOFNum, r.wc, op, batchID)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Move 3 — addQualityResult (10% out-of-spec)
// ---------------------------------------------------------------------------

func addQualityResult(db *sql.DB) error {
	rows, err := db.Query(`SELECT batch_id FROM batches WHERE finished_at IS NULL`)
	if err != nil {
		return fmt.Errorf("query inflight: %w", err)
	}
	var batches []string
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			rows.Close()
			return err
		}
		batches = append(batches, b)
	}
	rows.Close()

	for _, batchID := range batches {
		m := metrics[rand.Intn(len(metrics))]
		var value float64
		if rand.Float64() < 0.10 {
			// Out-of-spec — nudge ±5% of the range past the boundary.
			delta := (m.specMax - m.specMin) * 0.05
			if rand.Intn(2) == 0 {
				value = m.specMin - delta
			} else {
				value = m.specMax + delta
			}
		} else {
			// In-spec — uniform inside range.
			value = m.specMin + (m.specMax-m.specMin)*rand.Float64()
		}
		if _, err := db.Exec(
			`INSERT INTO quality_results (batch_id, measured_at, metric, value, spec_min, spec_max) VALUES (?, NOW(), ?, ?, ?, ?)`,
			batchID, m.name, value, m.specMin, m.specMax,
		); err != nil {
			log.Printf("[erpsim/quality] %s %s: %v", batchID, m.name, err)
			continue
		}
		flag := ""
		if value < m.specMin || value > m.specMax {
			flag = " OUT-OF-SPEC"
		}
		log.Printf("[erpsim/quality] %s %s=%.2f (spec %.2f–%.2f)%s", batchID, m.name, value, m.specMin, m.specMax, flag)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Move 4 — planNewOF (one per work_center per tick)
// ---------------------------------------------------------------------------

func planNewOF(db *sql.DB) error {
	for _, wc := range workCenters {
		ofNum := nextOF()
		product := products[rand.Intn(len(products))]
		plannedQty := 1000 + rand.Intn(5000)
		if _, err := db.Exec(
			`INSERT INTO work_orders (of_number, product_code, work_center, planned_qty, status) VALUES (?, ?, ?, ?, 'PLANNED')`,
			ofNum, product, wc, plannedQty,
		); err != nil {
			log.Printf("[erpsim/plan] %s: %v", ofNum, err)
			continue
		}
		log.Printf("[erpsim/plan] %s: %s × %d on %s", ofNum, product, plannedQty, wc)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Env helpers
// ---------------------------------------------------------------------------

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envDurOr(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
