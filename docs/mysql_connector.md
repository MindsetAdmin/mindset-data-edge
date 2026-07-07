# MySQL Connector — Detailed V1a Plan

_Executable slice of `docs/sql_connectors.md`. This doc is the ONLY thing you need open to implement the MySQL connector end-to-end. Postgres and MSSQL are deferred to V1b — they share ~80% of this code, so V1a lays the reusable foundation._

**Timebox:** 2 weeks. Ships to a customer smoke test at week 3.

**Related:**
- `docs/sql_connectors.md` — the broader strategy
- `docs/it_connectors.md` — the "why IT connectors" doc
- `docs/impact_engine.md` — downstream consumer of SQL-derived data

---

## 0. Why MySQL first

| Reason | Details |
|---|---|
| **Fastest to spin up locally** | Docker image `mysql:8` starts in <10s; MariaDB is API-identical |
| **Pure-Go driver** | `github.com/go-sql-driver/mysql` — no CGO, cross-compiles to the edge box on any OS |
| **Wide install base in target verticals** | Many mid-sized ERPs (Dolibarr, Odoo Community, custom Sage extensions), MES tools, and quality systems use MySQL/MariaDB — very common in agrifood and metallurgy plants |
| **Simple auth** | Username/password over TCP with optional TLS — no Kerberos, no SPN, no Windows Auth complexity |
| **Same wire protocol as MariaDB** | Two databases from one connector — free reach into the OSS market |
| **Fast learning loop** | We can dogfood the entire OT/IT flow in one afternoon against local Docker before touching a real customer |

---

## 1. Scope — V1a (this milestone)

### IN

- MySQL 8.x + MariaDB 10.x (same driver, same DSN, same code path)
- Read-only parameterized `SELECT` queries
- Username/password auth, optional TLS
- Connection pooling via `database/sql`
- Type mapping: all common MySQL types → JSON-friendly Go values
- YAML config for connections (`config/connections.yaml`)
- REST endpoints for connection CRUD + test + preview
- Pipeline Studio Connections page + SQL config panel
- Docker-compose fake MySQL "ERP" with seeded schema
- `cmd/erpsim` binary that simulates ERP activity against the fake DB
- Unit + integration tests

### OUT (deferred to V1b or later)

- Postgres, MSSQL (V1b — 1 more week after V1a; reuses this codebase)
- Oracle, SAP HANA (V2)
- Stored procedures (V1.5)
- CDC / streaming change events (V2)
- Non-SELECT queries (never — read-only is a safety property, not a limitation)

---

## 2. Files to add / modify

```
go.mod                                         MODIFY — add mysql driver
config/connections.yaml                        ADD    — connection definitions
docker-compose.dev-erp.yml                     ADD    — MySQL + init script

sim/erp/schema.mysql.sql                       ADD    — schema in MySQL 8 syntax
sim/erp/seed.mysql.sql                         ADD    — initial data

cmd/erpsim/main.go                             ADD    — ERP activity simulator (Go binary)

internal/connections/registry.go               ADD    — Registry, Open, Get, Close
internal/connections/config.go                 ADD    — Loader for connections.yaml
internal/connections/dsn.go                    ADD    — Build MySQL DSN from config
internal/connections/health.go                 ADD    — Startup verification: read-only?
internal/connections/registry_test.go          ADD    — unit tests

internal/functions/connectors/sql_query.go     REWRITE from stub
internal/functions/connectors/sql_query_test.go ADD   — unit tests

cmd/server/connections_handlers.go             ADD    — REST endpoints
cmd/server/main.go                             MODIFY — wire up the new handlers + register the sql_query function
cmd/agent/main.go                              MODIFY — register the sql_query function

frontend/pipeline-builder/src/pages/ConnectionsPage.jsx           ADD
frontend/pipeline-builder/src/components/SqlConfigPanel.jsx       ADD
frontend/pipeline-builder/src/api/client.js                       MODIFY — add connection endpoints
frontend/pipeline-builder/src/lib/functionDefaults.js             MODIFY — default sql_query config
frontend/pipeline-builder/src/lib/functionDocs.js                 MODIFY — sql_query docs

config/pipelines/examples/of_enrichment.yaml   ADD    — example pipeline
```

**~15 new files, ~5 modified files.** Nothing else touched.

---

## 3. Dependencies

Single new dependency:

```
go get github.com/go-sql-driver/mysql
```

Everything else (`database/sql`, `context`, `time`, `encoding/json`) is stdlib.

For integration tests, add:

```
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/mysql
```

That's it. Postgres and MSSQL will each add one driver later — no other structural changes.

---

## 4. MySQL DSN — the format that matters

The driver expects this format:

```
username:password@tcp(host:port)/dbname?param1=value1&param2=value2
```

**Parameters we always set:**

| Parameter | Value | Why |
|---|---|---|
| `parseTime=true` | (fixed) | Converts `DATETIME` / `TIMESTAMP` / `DATE` columns to `time.Time` instead of `[]byte`. Non-negotiable. |
| `loc=UTC` | (fixed) | All time values interpreted as UTC. Local-TZ conversion happens at the UI layer, never at the connector. |
| `charset=utf8mb4` | (fixed) | Real UTF-8 support (4-byte characters). Some MySQL setups still default to `utf8` (3-byte) which chokes on emojis + non-BMP chars. |
| `readTimeout=30s` | configurable | Per-query read timeout at the socket level. Belt-and-suspenders with `context.WithTimeout`. |
| `writeTimeout=10s` | configurable | Same, but for the query send. |
| `interpolateParams=false` | (fixed) | Use server-side prepared statements — safer against SQL injection than client-side interpolation. |
| `tls=true` or `tls=skip-verify` | configurable | Enable TLS. `skip-verify` accepts self-signed certs (dev only). |

**Example DSN we'll generate:**

```
mindset_readonly:PASSWORD@tcp(10.42.1.15:3306)/erp_prod?parseTime=true&loc=UTC&charset=utf8mb4&interpolateParams=false&readTimeout=30s&writeTimeout=10s&tls=true
```

The user never types this. `internal/connections/dsn.go` builds it from `config/connections.yaml`.

---

## 5. Connection config shape

`config/connections.yaml`:

```yaml
connections:
  - id: dev_erp
    name: "Dev ERP (MySQL simulator)"
    driver: mysql                     # V1a supports only "mysql" — switch on this later
    host: localhost
    port: 3306
    database: fake_erp
    username: mindset_readonly
    password_env: MINDSET_ERP_PASSWORD # NEVER inline the password
    tls: false                         # true | false | skip-verify
    read_timeout_seconds: 30
    write_timeout_seconds: 10
    max_open_conns: 5                  # driver default is unlimited — CAP THIS
    max_idle_conns: 2
    conn_max_lifetime_seconds: 300     # recycle every 5 min — avoids stale MySQL server-side timeouts
```

Password lives in an env var. The runtime resolves it at DSN build time. Empty or missing = startup error with a clear message.

---

## 6. Type mapping — MySQL → Go → JSON

The handler must return values that JSON-encode cleanly. MySQL types coerced as follows:

| MySQL type | Go type after scan | JSON output |
|---|---|---|
| `TINYINT(1)` | `int64` (careful — it's not `bool` by default) | number (0/1). Detect column type and coerce to bool if width == 1. |
| `TINYINT`, `SMALLINT`, `MEDIUMINT`, `INT`, `BIGINT` | `int64` | number |
| `FLOAT`, `DOUBLE` | `float64` | number |
| `DECIMAL`, `NUMERIC` | `string` (avoid float precision loss) | string. Document that clients requiring math should parse it themselves. |
| `VARCHAR`, `TEXT`, `CHAR` | `string` | string |
| `DATE`, `DATETIME`, `TIMESTAMP` | `time.Time` (via `parseTime=true`) | RFC3339 string |
| `TIME` | `[]byte` — driver quirk — convert to `string` (`HH:MM:SS`) | string |
| `BLOB`, `VARBINARY`, `BINARY` | `[]byte` | base64 string |
| `JSON` | `[]byte` — parse into `map[string]interface{}` or `[]interface{}` if valid, else keep as string | object / array / fallback string |
| `ENUM` | `string` (the label, not the ordinal) | string |
| `SET` | `string` (comma-separated) | string. Optionally split on comma if the user asks. |
| `NULL` | `nil` | JSON `null` |

**Code sketch for the coercion loop:**

```go
func mapRows(rows *sql.Rows) ([]map[string]interface{}, error) {
    cols, err := rows.Columns()
    if err != nil { return nil, err }
    colTypes, err := rows.ColumnTypes()
    if err != nil { return nil, err }

    out := []map[string]interface{}{}
    for rows.Next() {
        raw := make([]interface{}, len(cols))
        ptrs := make([]interface{}, len(cols))
        for i := range raw { ptrs[i] = &raw[i] }
        if err := rows.Scan(ptrs...); err != nil { return nil, err }

        row := make(map[string]interface{}, len(cols))
        for i, col := range cols {
            row[col] = coerce(raw[i], colTypes[i])
        }
        out = append(out, row)
    }
    return out, rows.Err()
}

func coerce(v interface{}, ct *sql.ColumnType) interface{} {
    if v == nil { return nil }
    switch x := v.(type) {
    case []byte:
        switch strings.ToUpper(ct.DatabaseTypeName()) {
        case "JSON":
            var j interface{}
            if err := json.Unmarshal(x, &j); err == nil { return j }
            return string(x)
        case "BLOB", "VARBINARY", "BINARY":
            return base64.StdEncoding.EncodeToString(x)
        default:
            return string(x)
        }
    case time.Time:
        return x.UTC().Format(time.RFC3339Nano)
    case int64:
        // TINYINT(1) → bool
        if strings.ToUpper(ct.DatabaseTypeName()) == "TINYINT" {
            if length, ok := ct.Length(); ok && length == 1 { return x != 0 }
        }
        return x
    default:
        return x
    }
}
```

Test every one of these paths in unit tests.

---

## 6b. Semantic mapping — the OTHER mapping layer

§6 above is the **type** layer: MySQL `INT` → Go `int64` → JSON `number`. That's identical across every customer.

The **semantic** layer is different. "The work order table is called `work_orders` with column `of_number`" is only true in ONE customer's DB. In production you'll meet all of these for the same conceptual object:

| System | Table for work orders | Column for OF number |
|---|---|---|
| SAP (ECC / S/4HANA) | `AUFK` | `AUFNR` |
| SAP MII | `PRODUCTION_ORDER` | `ORDER_NUMBER` |
| Odoo Community | `mrp_production` | `name` |
| Dolibarr | `llx_commande` | `ref` |
| Ignition MES | `production_order` | `wo_number` |
| Custom Access-derived thing | `Ordre_Fabrication` | `NumOF` |
| Excel export loaded into a MySQL table | `data_import_2024_wo` | `column_1` (yes, really) |

Same conceptual object. Six different schemas. Zero overlap.

If we let this heterogeneity leak into the rules engine, the Knowledge Graph, the Impact Engine, or the MCP, every downstream feature has to know every ERP's schema — which is exactly what we're competing AGAINST.

### The design — canonical model + `field_map`

We define a small, opinionated **canonical model** — the shape every downstream MindSet feature can assume:

```
WorkOrder { of_number, product_code, work_center, planned_qty, actual_qty, status,
            started_at, finished_at, operator_id }
Batch     { batch_id, of_number, started_at, finished_at, quality_status }
Product   { product_code, name, target_rate, recipe_id, hourly_margin }
Schedule  { work_center, of_number, planned_start, planned_end }
Quality   { batch_id, measured_at, metric, value, spec_min, spec_max }
Operator  { operator_id, name, shift }
```

These match the fake ERP schema in §10.2 exactly — which is by design: the sim IS the canonical shape.

The `sql_query` config gets a new optional field, `field_map`, that translates the customer's actual column names into canonical names:

```yaml
- id: fetch_current_of
  type: connector
  function: sql_query
  config:
    connection_id: "customer_erp_prod"
    query: |
      SELECT wo_no, mat_code, wc, plan_qty, act_qty, st, start_ts, end_ts, op_id
      FROM t_wo
      WHERE wc = :work_center AND st = 'R'
      LIMIT 1
    params:
      work_center: "{{ trigger.work_center }}"
    canonical: work_order              # tells downstream which canonical type this row is
    field_map:                         # NEW — customer columns → canonical field names
      of_number:    wo_no
      product_code: mat_code
      work_center:  wc
      planned_qty:  plan_qty
      actual_qty:   act_qty
      status:       st
      started_at:   start_ts
      finished_at:  end_ts
      operator_id:  op_id
```

### Runtime — what downstream nodes receive

The handler's output includes BOTH the raw row AND the canonical row for each result:

```json
{
  "rows": [
    {
      "wo_no": "WO-8891", "mat_code": "P-A12", "wc": "M1", "plan_qty": 5000,
      "act_qty": 3421, "st": "R", "start_ts": "2026-07-06T08:00:00Z", "end_ts": null, "op_id": "OP-004"
    }
  ],
  "canonical": [
    {
      "of_number": "WO-8891", "product_code": "P-A12", "work_center": "M1", "planned_qty": 5000,
      "actual_qty": 3421, "status": "R", "started_at": "2026-07-06T08:00:00Z",
      "finished_at": null, "operator_id": "OP-004"
    }
  ],
  "canonical_type": "work_order",
  "row_count": 1,
  "query_ms": 47
}
```

- Downstream nodes (rules, KG subscriber, Impact Engine, MCP) consume `canonical` — they don't care about the customer's column names.
- Raw `rows` stay available as an escape hatch — pipeline authors can access `wo_no` directly if they need something the canonical model doesn't expose.
- Missing `field_map` → `canonical` is empty and `canonical_type` is `null`. Handler doesn't error — degrades gracefully.

### Status normalization

Even after `field_map`, some values need value-level translation. SAP uses `CRTD`/`REL`/`TECO`; Odoo uses `draft`/`confirmed`/`progress`/`done`; Dolibarr uses `0`/`1`/`2`/`3`; a custom thing uses `EN COURS`/`FINI`. Downstream can't switch on all of these.

For V1a: add an optional `value_map` inside `field_map` for enum-like fields:

```yaml
field_map:
  status:
    from: st
    values:
      R:     RUNNING
      D:     DONE
      A:     ABORTED
      P:     PLANNED
```

Canonical `status` is always one of `PLANNED | RUNNING | DONE | ABORTED`. Downstream code can rely on that.

### System profiles — the follow-on

`field_map` handles ONE query. A customer with SAP has 20 canonical objects × N tables — no one types 200 field mappings by hand. Later we ship **system profiles** — pre-built libraries of `field_map`s + named queries for common systems:

| Profile | System | Ships in |
|---|---|---|
| `generic` | Any — user writes raw SQL + inline `field_map` (as shown above) | **V1a — this milestone** |
| `sap_mii_sql` | SAP MII SQL views | V1c |
| `odoo` | Odoo Community MySQL/PostgreSQL | V1c |
| `dolibarr` | Dolibarr open-source ERP (MySQL) | V1c |
| `ignition_mes` | Inductive Automation MES backend | V1c |
| `sap_ecc_direct` | Direct SAP ECC tables (rare — usually needs SAP consultant) | V2 |

A system profile is just a YAML file bundling named queries + field maps:

```yaml
# profiles/odoo.yaml
name: odoo
description: "Odoo Community 15+ ERP"
queries:
  current_work_order:
    sql: |
      SELECT id, name, product_id, workcenter_id, product_qty, qty_produced, state,
             date_start, date_finished
      FROM mrp_production
      WHERE workcenter_id = :work_center AND state = 'progress'
      LIMIT 1
    canonical: work_order
    field_map:
      of_number:    name
      product_code: product_id
      work_center:  workcenter_id
      planned_qty:  product_qty
      actual_qty:   qty_produced
      status:
        from: state
        values:
          draft:     PLANNED
          confirmed: PLANNED
          progress:  RUNNING
          done:      DONE
          cancel:    ABORTED
      started_at:   date_start
      finished_at:  date_finished
```

A pipeline using a profile just names the query:

```yaml
- id: fetch_current_of
  type: connector
  function: sql_query
  config:
    connection_id: "customer_odoo"      # has system_profile: odoo
    query_name: "current_work_order"    # resolved via the profile
    params:
      work_center: "{{ trigger.work_center }}"
```

The connection carries `system_profile: odoo`; the handler looks up `current_work_order` in the profile bundle; user never writes raw SQL for common cases. The **more profiles we ship, the bigger our library moat.**

### Why not push semantic mapping to a downstream `transform` node?

Three reasons:

1. **The canonical model is what everything downstream assumes.** Rules engine, KG, Impact Engine, MCP — they all consume `work_order.of_number`, not `wo_no` or `AUFNR`. Pushing mapping to a per-pipeline transform means every one of those has to be schema-aware. That's the mess we're avoiding.
2. **Library effect.** With `field_map` at the connector, we accumulate profiles as we onboard customers. Every SAP customer costs the same to onboard as the last one, not more. Pushing to per-pipeline transforms means every customer rebuilds it from scratch.
3. **Moat compounding.** System profiles are shippable IP. HighByte + Litmus both build this over time. Not having it = we lose to them on "1-day install" claims after the first pharma customer.

### V1a impact on the sprint

- **~1 extra day of work** to add `field_map` + `value_map` + canonical output shape to the handler
- Frontend adds a `field_map` editor to the SQL config panel — mostly a key/value grid + a small `value_map` sub-editor for enum fields
- Task list (§11) updated: extends Day 4 by 1 day (handler + tests) and Day 7 by half a day (UI panel)
- V1a still ships in 2 weeks

### V1c — system profiles

- 1–2 weeks after V1b (Postgres + MSSQL)
- Bundle profiles for Odoo, Dolibarr, Ignition MES, SAP MII
- Profile discovery: dropdown in the Connections page
- Custom profiles: a customer's SI writes their own YAML, drops it in `config/profiles/`, MindSet picks it up

---

## 7. Handler — full skeleton for `sql_query.go`

```go
package sqlquery

import (
    "context"
    "fmt"
    "strings"
    "time"

    "mindset-data-edge/internal/connections"
    "mindset-data-edge/internal/functions"
)

func GetFunction() *functions.Function {
    return &functions.Function{
        Name:        "sql_query",
        Type:        "connector",
        Description: "Execute a parameterized read-only SQL SELECT against a configured connection.",
        Inputs:      []functions.IO{ /* params from upstream */ },
        Outputs:     []functions.IO{{Name: "rows"}, {Name: "row_count"}, {Name: "query_ms"}},
        Handler:     handler,
    }
}

func handler(params map[string]interface{}) (interface{}, error) {
    connID, _ := params["connection_id"].(string)
    query, _  := params["query"].(string)
    if connID == "" { return nil, fmt.Errorf("sql_query: missing connection_id") }
    if query == ""  { return nil, fmt.Errorf("sql_query: missing query") }

    queryParams, _ := params["params"].(map[string]interface{})
    timeout        := durationOr(params["timeout_seconds"], 30*time.Second, 300*time.Second)
    limit          := intOr(params["limit"], 1000, 10000)

    if err := ensureSelectOnly(query); err != nil { return nil, err }
    query, args, err := bindPositional(query, queryParams)
    if err != nil { return nil, err }
    query = ensureLimit(query, limit)

    db, err := connections.Get(connID)
    if err != nil { return nil, err }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    started := time.Now()
    rows, err := db.QueryContext(ctx, query, args...)
    if err != nil { return nil, err }
    defer rows.Close()

    out, err := mapRows(rows)
    if err != nil { return nil, err }

    return map[string]interface{}{
        "rows":      out,
        "row_count": len(out),
        "query_ms":  time.Since(started).Milliseconds(),
    }, nil
}
```

### Guards to implement

**`ensureSelectOnly(query string) error`** — strips leading comments + whitespace, checks that the first keyword is `SELECT` (case-insensitive). Rejects `INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE`, `GRANT`, `REVOKE`, `CALL`, and multi-statement queries (any `;` before EOF that isn't inside a string literal).

**`bindPositional(query string, params map[string]interface{}) (string, []interface{}, error)`** — MySQL uses `?` positional placeholders. Users write `:name` for clarity; we convert `:name` → `?` and build the args slice in the order the placeholders appear. Multiple uses of the same `:name` = same value passed each time. Missing param = error.

**`ensureLimit(query string, max int) string`** — if the query already has a top-level `LIMIT` clause, respect it (but cap it at `max`). If not, append ` LIMIT {max}` before any trailing semicolon.

---

## 8. Connection registry — `internal/connections/`

```
internal/connections/
├── config.go        Load + validate config/connections.yaml
├── dsn.go           BuildMySQLDSN(cfg Config) string
├── registry.go      Registry{ mu, dbs map[string]*sql.DB }; Open, Get, Close
├── health.go        VerifyReadOnly(db *sql.DB) error — runs a canary INSERT and expects it to fail
└── registry_test.go
```

**Registry lifecycle:**

- `Load(path string)` reads YAML at startup
- `Get(id string) (*sql.DB, error)` — lazy open on first request; caches the `*sql.DB` (which is itself a pool)
- Applies `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime` from config
- `CloseAll()` on shutdown

**Health check (`health.go`):**

On first `Get`, run:

```sql
SELECT 1;                                         -- reachability + credentials
-- then, in a savepoint or transaction:
CREATE TEMPORARY TABLE mindset_probe (id INT);    -- fails if user can't write — good signal
```

If the CREATE succeeds (user CAN write), log a **warning** but don't refuse the connection. Enterprise IT sometimes provisions accounts with more permissions than they should. Log the warning so it appears in ops review.

Also cache a `read_only bool` on the connection metadata — exposed on `/api/connections` so the Pipeline Studio can badge non-read-only connections in red.

---

## 9. REST endpoints (`cmd/server/connections_handlers.go`)

| Method | Path | Body / Result |
|---|---|---|
| `GET` | `/api/connections` | Returns list — **never** includes passwords. Fields: id, name, driver, host, port, database, username, tls, read_only, status (`ok`/`error`), last_checked |
| `POST` | `/api/connections` | Create/update. Body: full config shape. Password field named `password` but stored as env var reference OR encrypted at rest (V1a: env var reference only). Returns 201 with the object. |
| `POST` | `/api/connections/{id}/test` | Runs the health check. Returns `{ok, latency_ms, read_only, error?}` |
| `POST` | `/api/connections/{id}/preview` | Body: `{query, params, limit}`. Returns first 5 rows (or `limit`, whichever is smaller) — used by the Pipeline Studio preview button. Enforces the same `ensureSelectOnly`/`ensureLimit` guards as `sql_query`. |
| `DELETE` | `/api/connections/{id}` | Removes from SQLite + closes the pool |

Persistence: new `connections` table in `data/mindset.db` (auto-created in `internal/storage/sqlite.go`).

Schema:
```sql
CREATE TABLE IF NOT EXISTS connections (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  driver TEXT NOT NULL,
  host TEXT NOT NULL,
  port INT NOT NULL,
  database TEXT NOT NULL,
  username TEXT NOT NULL,
  password_env TEXT NOT NULL,
  tls TEXT NOT NULL,          -- "true" | "false" | "skip-verify"
  read_timeout_seconds INT NOT NULL DEFAULT 30,
  write_timeout_seconds INT NOT NULL DEFAULT 10,
  max_open_conns INT NOT NULL DEFAULT 5,
  max_idle_conns INT NOT NULL DEFAULT 2,
  conn_max_lifetime_seconds INT NOT NULL DEFAULT 300,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## 10. Simulation environment — a working fake ERP in one command

### 10.1. `docker-compose.dev-erp.yml`

```yaml
services:
  mysql-erp:
    image: mysql:8
    container_name: mindset-erp
    restart: unless-stopped
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: rootdev
      MYSQL_DATABASE: fake_erp
      MYSQL_USER: mindset_readonly
      MYSQL_PASSWORD: readonly_dev
    volumes:
      - ./sim/erp/schema.mysql.sql:/docker-entrypoint-initdb.d/01-schema.sql:ro
      - ./sim/erp/seed.mysql.sql:/docker-entrypoint-initdb.d/02-seed.sql:ro
      - ./sim/erp/grant.mysql.sql:/docker-entrypoint-initdb.d/03-grant.sql:ro
      - erp-data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-uroot", "-prootdev"]
      interval: 5s
      timeout: 3s
      retries: 20
volumes:
  erp-data:
```

The three init scripts run in numeric order — schema → seed → grants. The grants script does:

```sql
-- sim/erp/grant.mysql.sql
GRANT SELECT ON fake_erp.* TO 'mindset_readonly'@'%';
FLUSH PRIVILEGES;
```

That gives `mindset_readonly` only `SELECT`, so our health check confirms read-only-ness on a real setup.

### 10.2. `sim/erp/schema.mysql.sql` (MySQL 8 syntax)

```sql
CREATE TABLE IF NOT EXISTS work_orders (
  of_number      VARCHAR(64) PRIMARY KEY,
  product_code   VARCHAR(64) NOT NULL,
  work_center    VARCHAR(64) NOT NULL,
  planned_qty    INT NOT NULL,
  actual_qty     INT NOT NULL DEFAULT 0,
  status         ENUM('PLANNED','RUNNING','DONE','ABORTED') NOT NULL DEFAULT 'PLANNED',
  started_at     DATETIME NULL,
  finished_at    DATETIME NULL,
  operator_id    VARCHAR(32) NULL,
  INDEX idx_wc_status (work_center, status)
);

CREATE TABLE IF NOT EXISTS batches (
  batch_id       VARCHAR(64) PRIMARY KEY,
  of_number      VARCHAR(64) NOT NULL,
  started_at     DATETIME NOT NULL,
  finished_at    DATETIME NULL,
  quality_status ENUM('PASS','FAIL','REWORK') NULL,
  FOREIGN KEY (of_number) REFERENCES work_orders(of_number),
  INDEX idx_batches_of (of_number)
);

CREATE TABLE IF NOT EXISTS products (
  product_code   VARCHAR(64) PRIMARY KEY,
  name           VARCHAR(255) NOT NULL,
  target_rate    INT NULL,
  recipe_id      VARCHAR(64) NULL,
  hourly_margin  DECIMAL(10,2) NULL
);

CREATE TABLE IF NOT EXISTS schedules (
  id             INT AUTO_INCREMENT PRIMARY KEY,
  work_center    VARCHAR(64) NOT NULL,
  of_number      VARCHAR(64) NOT NULL,
  planned_start  DATETIME NOT NULL,
  planned_end    DATETIME NOT NULL,
  FOREIGN KEY (of_number) REFERENCES work_orders(of_number)
);

CREATE TABLE IF NOT EXISTS quality_results (
  id             INT AUTO_INCREMENT PRIMARY KEY,
  batch_id       VARCHAR(64) NOT NULL,
  measured_at    DATETIME NOT NULL,
  metric         VARCHAR(64) NOT NULL,
  value          DECIMAL(10,4) NOT NULL,
  spec_min       DECIMAL(10,4) NULL,
  spec_max       DECIMAL(10,4) NULL,
  FOREIGN KEY (batch_id) REFERENCES batches(batch_id),
  INDEX idx_quality_batch (batch_id, measured_at)
);

CREATE TABLE IF NOT EXISTS operators (
  operator_id    VARCHAR(32) PRIMARY KEY,
  name           VARCHAR(255) NOT NULL,
  shift          ENUM('MORNING','AFTERNOON','NIGHT') NOT NULL
);
```

### 10.3. `sim/erp/seed.mysql.sql`

Initial data:

- ~50 products (`PROD-A01` through `PROD-A50`) with realistic `target_rate` and `hourly_margin`
- ~20 operators (`OP-001` through `OP-020`) rotating shifts
- ~200 historical work orders spanning the last 30 days (mostly `DONE`)
- **3 currently-RUNNING** work orders on `machine1`, `machine2`, `machine3` — matching the OPC-UA simulator's work_center names, so the demos wire up
- ~50 in-flight and closed batches
- ~500 quality_results
- Schedule entries for the next 24h across all 3 work centers

Realistic enough that a plant manager watching the demo doesn't laugh.

### 10.4. `cmd/erpsim/main.go` — the activity simulator

A small binary that runs alongside the dev stack and makes the ERP look alive:

```go
package main

import (
    "database/sql"
    "log"
    "math/rand"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "mindset_writer:writer_dev@tcp(localhost:3306)/fake_erp?parseTime=true&loc=UTC&charset=utf8mb4"
    db, err := sql.Open("mysql", dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()

    go loop(30*time.Second,  advanceRunningOFs, db)   // increment actual_qty on RUNNING OFs
    go loop(5*time.Minute,   rotateOFs, db)           // finish RUNNING → start next PLANNED
    go loop(10*time.Minute,  addQualityResult, db)    // one new quality reading per in-flight batch
    go loop(1*time.Hour,     planNewOF, db)           // create a new PLANNED OF for each work_center

    select {} // block forever
}

func loop(every time.Duration, fn func(*sql.DB) error, db *sql.DB) {
    for range time.Tick(every) {
        if err := fn(db); err != nil { log.Printf("erpsim: %v", err) }
    }
}
```

Note: `cmd/erpsim` uses a `mindset_writer` MySQL user (not `mindset_readonly`) — the sim needs write perms; the real connector must not have them. This is the boundary that the health check enforces.

Add a matching grant in `grant.mysql.sql`:

```sql
CREATE USER 'mindset_writer'@'%' IDENTIFIED BY 'writer_dev';
GRANT SELECT, INSERT, UPDATE ON fake_erp.* TO 'mindset_writer'@'%';
FLUSH PRIVILEGES;
```

### 10.5. Bring the stack up

```powershell
docker compose -f docker-compose.dev-erp.yml up -d
$env:MINDSET_ERP_PASSWORD = "readonly_dev"
go run ./cmd/erpsim                                    # run in one terminal
.\run.ps1                                              # existing stack, another terminal
```

Open the frontend at `http://localhost:5173` → Connections → add `dev_erp` (or read from `config/connections.yaml`) → run the example OF-enrichment pipeline.

---

## 11. Ordered task list — the 2-week execution plan

### Week 1 — infrastructure + backend

**Day 1 — Simulation stack**
- Create `docker-compose.dev-erp.yml`
- Write `sim/erp/schema.mysql.sql`, `seed.mysql.sql`, `grant.mysql.sql`
- Bring MySQL up, verify seed data via `mysql` CLI
- Verify `mindset_readonly` cannot INSERT (permissions correctly locked)

**Day 2 — `cmd/erpsim`**
- Scaffold the binary
- Implement `advanceRunningOFs`, `rotateOFs`, `addQualityResult`, `planNewOF`
- Watch the DB update in real time over lunch — sanity check the intervals

**Day 3 — `internal/connections/` package**
- Config loader
- DSN builder (with all the MySQL params from §4)
- Registry (Open, Get, CloseAll)
- Unit tests

**Day 4–5 — Rewrite `sql_query.go`** _(2 days — extended for the semantic mapping layer, §6b)_
- Handler skeleton (§7)
- `ensureSelectOnly`, `bindPositional`, `ensureLimit`
- `mapRows` + `coerce` with type-aware coercion
- **`applyFieldMap` — build canonical rows from raw rows + `field_map`**
- **`applyValueMap` — enum value-level translation for `status`-like fields**
- Unit tests with SQLite in-memory (fast, no docker) — cover both `rows` and `canonical` output paths

**Day 6 — REST endpoints** _(was Day 5 — shifted by 1 for the field_map work)_
- `cmd/server/connections_handlers.go` — 5 routes
- SQLite `connections` table
- Wire into `main.go`
- Register `sql_query` in both `cmd/server` and `cmd/agent`

### Week 2 — UI + tests + demo

**Day 6 — Connections page**
- `ConnectionsPage.jsx` — list + create form + Test button
- Add API client methods in `src/api/client.js`
- New route in `App.jsx`

**Day 7–8 — SQL config panel** _(now 1.5 days — includes the field_map editor)_
- `SqlConfigPanel.jsx` — connection dropdown + query editor + params grid + Preview button
- **`FieldMapEditor.jsx` sub-component** — canonical-type dropdown (`work_order` / `batch` / `product` / …) + key-value grid mapping raw column → canonical field + optional `value_map` sub-editor for enum fields
- Wire into `BuilderPage.jsx` so it appears when a `sql_query` node is selected
- Update `functionDefaults.js` + `functionDocs.js`

**Day 8 — Example pipeline**
- Write `config/pipelines/examples/of_enrichment.yaml`:
  - `mqtt_subscribe` on `mindset/events/status-change`
  - `sql_query` (current RUNNING OF for the work_center)
  - `mqtt_publish` on `mindset/events/status-change-enriched`
- Run against the live sim, verify enrichment in the dashboard

**Day 9 — Integration tests**
- `internal/e2e/sql_pipeline_test.go` — testcontainers MySQL + full pipeline execution
- Cover: happy path, timeout, injection attempt, non-SELECT rejection, read-only enforcement

**Day 10 — Polish + docs**
- Update `docs/COMPONENTS.md`, `docs/ARCHITECTURE.md`, `CLAUDE.md` function catalog (remove "demo stub" flag)
- Update `docs/mindset.md` §8 Module 4 to note SQL/MySQL V1a shipped
- Screen-record a 30-second demo of an enriched micro-stop for future beta pitches

### Week 3 — first customer smoke test

- Identify beta customer (pharma or agrifood — check whether their ERP is MySQL/MariaDB)
- If yes: provision a read-only user on their DB, get one small schema extract
- Wire one enrichment pipeline against their real data
- Show it in the dashboard within 48h of receiving credentials

---

## 12. Testing plan

### Unit (no docker, no network)

| Test | Assertion |
|---|---|
| `TestEnsureSelectOnly_accepts_SELECT` | Various leading whitespace + comment forms |
| `TestEnsureSelectOnly_rejects_INSERT_UPDATE_DELETE_DROP_ALTER_TRUNCATE_CALL` | All rejected |
| `TestEnsureSelectOnly_rejects_multi_statement` | `SELECT 1; DROP TABLE x;` → error |
| `TestBindPositional_named_to_question_mark` | `:foo AND :bar` + `{foo: 1, bar: 2}` → `? AND ?`, `[1, 2]` |
| `TestBindPositional_repeated_placeholder` | `:x + :x` → `? + ?`, `[v, v]` |
| `TestBindPositional_missing_param` | Error |
| `TestEnsureLimit_appends_when_missing` | `SELECT *` → `SELECT * LIMIT 1000` |
| `TestEnsureLimit_respects_smaller_user_limit` | `SELECT * LIMIT 10` → `SELECT * LIMIT 10` |
| `TestEnsureLimit_caps_larger_user_limit` | `SELECT * LIMIT 99999` → `SELECT * LIMIT 10000` |
| `TestCoerce_TINYINT_1_becomes_bool` | `TINYINT(1)` value `1` → `true` |
| `TestCoerce_JSON_parses` | JSON column with `{"a":1}` → `map[string]interface{}{"a":1}` |
| `TestCoerce_time_RFC3339` | `DATETIME` → RFC3339 string in UTC |
| `TestCoerce_DECIMAL_stays_string` | `DECIMAL(10,2)` `123.45` → `"123.45"` |
| `TestCoerce_null` | `nil` → `nil` (JSON `null`) |
| `TestRegistry_pool_reuse` | Two `Get(id)` calls return the same `*sql.DB` |

### Integration (testcontainers MySQL)

| Test | Assertion |
|---|---|
| `TestHappyPath` | Real MySQL + seed → SELECT returns rows with correct types |
| `TestTimeoutKicksIn` | Query `SELECT SLEEP(5)` with 1s timeout → context deadline error |
| `TestInjectionAttempt` | `:name` = `1; DROP TABLE x` → parameter, table not dropped |
| `TestReadOnlyEnforcement` | `sql_query` with `INSERT INTO ...` → rejected before connecting |
| `TestHealthCheck_readonly_user` | Verify `mindset_readonly` fails CREATE TEMPORARY TABLE — health returns `read_only: true` |

### End-to-end (docker-compose stack)

| Test | Assertion |
|---|---|
| Manual — full demo | `run.ps1` up, ERP up, sim running → dashboard shows enriched micro-stops with `of_number` populated within 5s of a stop event |

---

## 13. Known MySQL gotchas — how we handle each

| Gotcha | Impact | V1a handling |
|---|---|---|
| **TINYINT(1) is not `bool`** in the driver — returns `int64` by default | Booleans show as `0`/`1` in JSON | Detect column type + width in `coerce`, return `bool` for TINYINT(1) |
| **DECIMAL precision loss if scanned to `float64`** | Off-by-one-cent financial errors | Scan as `[]byte` → return as string; document that clients parse if they need math |
| **DATETIME vs TIMESTAMP timezone semantics differ** — DATETIME is naïve, TIMESTAMP is UTC | Wrong-hour displays | Force `parseTime=true&loc=UTC` in DSN; ONLY convert to plant TZ at the UI |
| **`utf8` charset ≠ real UTF-8** (3-byte) | Non-BMP chars break | Force `charset=utf8mb4` in DSN |
| **`max_allowed_packet` server-side limit** (default 4MB in some builds) | Very large result rejected | Enforce row + column-size limits at handler level; document the setting for customer DBs |
| **Client-side parameter interpolation** — driver has an option, but it's less safe | SQL injection surface | Force `interpolateParams=false` — server-side prepared statements only |
| **`sql.DB` connections lost after server-side timeout** (`wait_timeout` default 8h, some setups 60s) | Silent failures | Set `SetConnMaxLifetime(5m)` — recycles before typical timeouts |
| **Case-sensitivity of identifiers depends on OS/`lower_case_table_names`** | Query works dev, fails prod | Document — customer's DBA needs to confirm case sensitivity settings match dev |
| **BOOL / BOOLEAN is an alias for TINYINT(1)** | Same as above | Same handling |
| **JSON type only in MySQL 5.7+ and MariaDB 10.2+** | Query fails on older MariaDB | Document the minimum version; V1a claims MySQL 8 + MariaDB 10 |

---

## 14. Definition of Done — V1a MySQL

Ship when ALL of:

- [ ] `docker-compose.dev-erp.yml` starts MySQL + seeds + grants in <30s
- [ ] `cmd/erpsim` runs standalone and updates the DB every 30s minimum
- [ ] `internal/connections/` package: config loader + DSN builder + registry + health check + unit tests green
- [ ] `sql_query` handler: rewritten, all unit tests green (both raw + canonical output paths)
- [ ] `field_map` + `value_map` supported end-to-end (config → handler → downstream node receives `canonical` rows)
- [ ] Frontend `FieldMapEditor` component allows user to define + edit `field_map` visually
- [ ] `/api/connections` REST endpoints: 5 routes wired, tested via curl
- [ ] Connections page in Pipeline Studio: create, test, save, delete works
- [ ] SQL config panel in Pipeline Studio: connection dropdown + query + params + preview works
- [ ] `config/pipelines/examples/of_enrichment.yaml` runs end-to-end against the sim
- [ ] Integration tests (testcontainers MySQL) green
- [ ] Dashboard shows `of_number` on enriched micro-stops
- [ ] `CLAUDE.md` updated — `sql_query` no longer flagged "demo stub"
- [ ] 30-second demo screen recording saved to `docs/demos/`

---

## 15. What we get for free in V1b (Postgres + MSSQL)

The V1a architecture is deliberately driver-agnostic below the DSN layer:

- **Registry, Get, health check** — reused as-is
- **Handler shell** (`sql_query.go`) — reused; the guards (`ensureSelectOnly`, `bindPositional`, `ensureLimit`) work across drivers with minor placeholder-syntax swaps
- **REST endpoints** — reused; only the `driver` field on the connection payload changes
- **Frontend** — Connections page + SQL panel reused; only the "driver" dropdown expands to `mysql | postgres | mssql`

V1b adds:

- `internal/connections/dsn_postgres.go` — Postgres DSN builder + `$1`/`$2` placeholder syntax
- `internal/connections/dsn_mssql.go` — MSSQL DSN builder + `@p1`/`@p2` placeholder syntax
- Two more integration test files (testcontainers Postgres + MSSQL)
- Two more type-mapping test coverage extensions (Postgres arrays, MSSQL `datetime2`)

Estimated V1b timebox: **1 week**, from a clean V1a merge.

---

## 16. TL;DR for the 2-week sprint

**Week 1 — backend**
1. Docker MySQL + fake ERP schema + seed
2. `cmd/erpsim` — makes the ERP look alive
3. `internal/connections/` — DSN, registry, health
4. `sql_query.go` — real handler with all guards
5. `/api/connections` — 5 REST routes

**Week 2 — UI + demo**
6. Connections page
7. SQL config panel in Pipeline Studio
8. `of_enrichment.yaml` example pipeline
9. Integration tests
10. Docs + demo recording

**Week 3 — first customer smoke test.**

Ship the MySQL connector, then V1b (Postgres + MSSQL) is 1 more week on the same foundation.
