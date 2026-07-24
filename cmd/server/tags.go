package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"sort"
	"sync"
)

// Tag mirrors the raw OPC-UA payload published by the agent (mqtt/publisher.go
// TagMessage) on mindset/raw/<nodeID>.
type Tag struct {
	NodeID    string      `json:"node_id"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	DataType  string      `json:"data_type"`
	Timestamp int64       `json:"timestamp_ms"`
}

// TagRegistry keeps the live set of OPC-UA tags discovered by the agent. It is
// updated from mindset/raw/# and persisted to the shared SQLite DB so tags still
// show in the UI when the agent isn't running.
type TagRegistry struct {
	mu        sync.RWMutex
	tags      map[string]Tag
	db        *sql.DB
	persistCh chan Tag
}

// persistQueueSize is generous relative to real update rates (a handful of
// tags updating at most a few times/sec) so it only ever fills if the DB
// writer is genuinely stuck, not from ordinary bursts.
const persistQueueSize = 256

func NewTagRegistry(db *sql.DB) *TagRegistry {
	r := &TagRegistry{tags: make(map[string]Tag), db: db, persistCh: make(chan Tag, persistQueueSize)}
	if db != nil {
		r.ensureTable()
		r.loadFromDB()
		go r.persistLoop()
	}
	return r
}

// persistLoop is the only goroutine that ever writes to the tags table.
// upsert() hands it work over persistCh instead of writing inline (Entry
// 136): this SQLite file is also written concurrently by cmd/agent's KG
// subscriber and server.exe's own KG writes, and tag updates arrive multiple
// times per second per tag — a live incident showed the MQTT callback
// goroutine that drives the whole dashboard (tags, topic counters, machine
// state) going completely silent mid-session with no error logged anywhere,
// while the underlying OPC-UA subscription and MQTT publish kept running
// fine. The only blocking I/O on that callback's path was this DB write, and
// PRAGMA busy_timeout (Entry 131) bounds SQLite-level lock waits but not
// every possible driver/OS-level contention mode with a pure-Go SQLite
// driver under sustained concurrent writers. Moving persistence off the live
// path makes that class of freeze structurally impossible: a stuck DB write
// can now only stall this one background goroutine, never the live
// dashboard/tag/state pipeline, which is served entirely from the in-memory
// map upsert() updates synchronously before ever touching persistCh.
func (r *TagRegistry) persistLoop() {
	for t := range r.persistCh {
		valBytes, _ := json.Marshal(t.Value)
		_, err := r.db.Exec(`INSERT INTO tags (node_id, name, value, data_type, timestamp_ms, updated_at)
			VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(node_id) DO UPDATE SET
				name=excluded.name, value=excluded.value, data_type=excluded.data_type,
				timestamp_ms=excluded.timestamp_ms, updated_at=CURRENT_TIMESTAMP`,
			t.NodeID, t.Name, string(valBytes), t.DataType, t.Timestamp)
		if err != nil {
			log.Printf("[TAGS] persist: %v", err)
		}
	}
}

func (r *TagRegistry) ensureTable() {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS tags (
		node_id TEXT PRIMARY KEY,
		name TEXT,
		value TEXT,
		data_type TEXT,
		timestamp_ms INTEGER,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Printf("[TAGS] create table: %v", err)
	}
}

func (r *TagRegistry) loadFromDB() {
	rows, err := r.db.Query(`SELECT node_id, name, value, data_type, timestamp_ms FROM tags`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t Tag
		var valStr string
		if err := rows.Scan(&t.NodeID, &t.Name, &valStr, &t.DataType, &t.Timestamp); err == nil {
			_ = json.Unmarshal([]byte(valStr), &t.Value)
			r.tags[t.NodeID] = t
		}
	}
	log.Printf("[TAGS] Loaded %d persisted tags", len(r.tags))
}

func (r *TagRegistry) upsert(t Tag) {
	r.mu.Lock()
	r.tags[t.NodeID] = t
	r.mu.Unlock()

	if r.persistCh != nil {
		select {
		case r.persistCh <- t:
		default:
			// Queue full — persistLoop is falling behind or stuck. Drop this
			// sample rather than block the live MQTT callback: the in-memory
			// map above (what every live read actually serves) is already
			// correct, and the next tick will offer a fresher value anyway.
			log.Printf("[TAGS] persist queue full, dropping sample for %s", t.NodeID)
		}
	}
}

func (r *TagRegistry) list() []Tag {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Tag, 0, len(r.tags))
	for _, t := range r.tags {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].NodeID < out[j].NodeID })
	return out
}
