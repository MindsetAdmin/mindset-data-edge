# SQL Connectors — V1 Implementation Plan

_Purpose: replace the current `sql_query` demo stub with a real, production-ready SQL connector that pulls IT context (OF, batch, product, schedule, quality) into MindSet Data pipelines. This is the first concrete step toward the OT/IT reconciliation moat (Fuzzy Join)._

Cross-refs: `docs/it_connectors.md` (strategy) · `docs/mindset.md` §8 (module 4) · `docs/impact_engine.md` (uses SQL data downstream)

---

## 0. Why now

The V0 platform can already ingest OPC-UA → contextualize to ISA-95 → detect stops → compute cost. But every enrichment is calculated from data the EDGE box already sees. The real leverage — the Impact Engine, the Fuzzy Join OF/batch, an AI-native context layer — needs data that lives in the customer's ERP / MES / LIMS.

SQL connectors are the fastest, most portable path to that data:
- ~90% of mid-sized manufacturers have at least one SQL-accessible IT system (SQL Server, Postgres, MySQL, Oracle)
- No vendor-specific SDK required
- No cloud round-trip — the edge box queries the customer's local DB
- Ships in weeks, not quarters (unlike native SAP/Werum connectors)

---

## 1. Scope — V1 in / out

### In V1

- **Drivers:** PostgreSQL, MySQL/MariaDB, Microsoft SQL Server, SQLite (dev)
- **Query modes:** parameterized SELECT (read-only)
- **Auth:** username + password, TLS optional
- **Result shape:** `[]map[string]interface{}` — one map per row
- **Config storage:** connections defined in `config/connections.yaml` OR via new `/api/connections` REST endpoints
- **Pipeline Studio UI:** SQL configuration panel with connection dropdown + query editor + parameter binding + preview
- **Simulation environment:** docker-compose + fake ERP schema + Go seed script

### Out of V1 (queued)

- Oracle (needs CGO — deferred to V1.5)
- SAP HANA (proprietary driver — V2 with SAP consultant)
- CDC / streaming (V2)
- Write queries (V1 = read-only for safety; no INSERT/UPDATE/DELETE)
- Stored procedures (V1.5)
- Multi-statement transactions (never — not the pattern we want on the edge)

---

## 2. Architecture — where the SQL connector fits

```
config/connections.yaml      ← DSNs (host, port, user, password refs)
      │
      ▼
internal/connections/        ← new package: registry, DSN parsing, pool management
      │
      ▼
internal/functions/connectors/sql_query.go   ← rewrites the demo stub
      │
      ▼
Pipeline engine executes sql_query as a normal node → outputs rows to next node
```

**Key design choices:**

1. **Connections are first-class**, not embedded in the pipeline YAML. Same DSN reused across many pipelines — one connection pool per DSN.
2. **Read-only DB user recommended** for every deployment. Enforced by documentation + connection health check that verifies user role.
3. **Query timeout mandatory.** No timeout = pipeline hangs = engine dead. Enforce 30s default, max 300s.
4. **Row limit mandatory.** Default 1000 rows per query. Prevents accidental table scans.
5. **The connector runs on the SAME binary that owns the pipeline** — for the default UI-controlled mode, that's `cmd/server`. The agent gets the same registration so agent-owned pipelines work too.

---

## 3. Function spec — `sql_query`

Registered in both `cmd/server/main.go#buildRegistry` and `cmd/agent/main.go`:

| Field | Value |
|---|---|
| **Name** | `sql_query` |
| **Type** | `connector` |
| **Description** | Execute a parameterized SELECT against a configured SQL connection. Returns rows. |
| **Inputs** | (optional) parameter values from upstream nodes |
| **Outputs** | `rows` — `[]map[string]interface{}` |
| **Handler** | `func(params map[string]interface{}) (interface{}, error)` |

### Config schema (in YAML pipeline)

```yaml
- id: fetch_current_of
  type: connector
  function: sql_query
  config:
    connection_id: "customer_erp_prod"        # ref to config/connections.yaml
    query: |
      SELECT of_number, product_code, planned_qty, actual_qty, status
      FROM work_orders
      WHERE work_center = :work_center
        AND status = 'RUNNING'
      LIMIT 1
    params:
      work_center: "{{ trigger.work_center }}"   # template — filled from trigger data
    timeout_seconds: 15                          # optional, default 30
    limit: 100                                   # optional, default 1000
```

### Output shape (what the next node receives)

```go
map[string]interface{}{
  "rows": []map[string]interface{}{
    {"of_number": "OF-2026-8891", "product_code": "PROD-A12", "planned_qty": 5000, "actual_qty": 3421, "status": "RUNNING"},
    // ...
  },
  "row_count": 1,
  "query_ms": 47,
}
```

---

## 4. Connection configuration

### File-based (V1 default)

`config/connections.yaml`:

```yaml
connections:
  - id: customer_erp_prod
    name: "Customer ERP (production replica)"
    driver: postgres              # postgres | mysql | mssql | sqlite
    host: 10.42.1.15
    port: 5432
    database: erp_prod
    username: mindset_readonly
    password_env: MINDSET_ERP_PASSWORD   # never inline the password
    tls: true
    max_open_conns: 5
    max_idle_conns: 2
    conn_max_lifetime_seconds: 300
```

### Runtime (V1 optional)

`POST /api/connections` — same shape as YAML, persisted to `data/mindset.db` in a new `connections` table. Enables the Pipeline Studio to add connections without editing YAML.

### Credentials

- **Never in YAML.** Passwords always via `password_env` referencing an environment variable.
- Docker deployments: use `docker secrets` mounted to `/run/secrets/`.
- Rotate the DB user's password every 90 days minimum (customer's responsibility, documented).
- Add a health check on startup that verifies the user has `SELECT`-only permissions on the schemas listed.

---

## 5. Security

| Concern | V1 mitigation |
|---|---|
| **SQL injection** | Only parameterized queries. Named placeholders (`:name`) resolved via `database/sql` driver's `NamedArgs`. No string concatenation. |
| **Runaway queries** | Mandatory `timeout_seconds` (default 30, max 300) + mandatory `LIMIT` clause enforcement. Query without a LIMIT rejected at pipeline validation. |
| **Credential leakage** | Passwords via env vars only; never logged; redacted in error messages. |
| **Excessive load on customer DB** | Connection pool cap (`max_open_conns` default 5). Pipeline-level rate limiting deferred to V1.5. |
| **TLS** | Enabled by default when talking to remote DBs. Doc explicitly recommends against unencrypted connections. |
| **Read-only enforcement** | (a) documented DB-user recommendation, (b) startup health check verifies the user cannot INSERT/UPDATE/DELETE on the schemas listed, (c) sql_query handler rejects any query whose first token isn't `SELECT` (case-insensitive, whitespace-tolerant). |

---

## 6. Simulation environment — how to test without a real customer DB

The whole point of shipping SQL connectors early is dogfooding on OT/IT reconciliation. That needs an IT side that looks realistic. Bundle a fake ERP with the dev stack.

### 6.1. Files added

```
docker-compose.dev-erp.yml       # new — spins up Postgres with a seeded schema
sim/erp/schema.sql               # fake ERP tables + indexes
sim/erp/seed.sql                 # initial rows
cmd/erpsim/main.go               # Go binary: periodically inserts new OFs, advances batch states
```

### 6.2. Fake ERP schema

```sql
CREATE TABLE work_orders (
  of_number      TEXT PRIMARY KEY,
  product_code   TEXT NOT NULL,
  work_center    TEXT NOT NULL,
  planned_qty    INT  NOT NULL,
  actual_qty     INT  NOT NULL DEFAULT 0,
  status         TEXT NOT NULL,       -- PLANNED | RUNNING | DONE | ABORTED
  started_at     TIMESTAMPTZ,
  finished_at    TIMESTAMPTZ,
  operator_id    TEXT
);

CREATE TABLE batches (
  batch_id       TEXT PRIMARY KEY,
  of_number      TEXT REFERENCES work_orders,
  started_at     TIMESTAMPTZ NOT NULL,
  finished_at    TIMESTAMPTZ,
  quality_status TEXT                 -- PASS | FAIL | REWORK | NULL when in-flight
);

CREATE TABLE products (
  product_code   TEXT PRIMARY KEY,
  name           TEXT NOT NULL,
  target_rate    INT,                 -- units/hour target
  recipe_id      TEXT,
  hourly_margin  NUMERIC(10,2)        -- € margin per hour of runtime, if known
);

CREATE TABLE schedules (
  id             SERIAL PRIMARY KEY,
  work_center    TEXT NOT NULL,
  of_number      TEXT REFERENCES work_orders,
  planned_start  TIMESTAMPTZ NOT NULL,
  planned_end    TIMESTAMPTZ NOT NULL
);

CREATE TABLE quality_results (
  id             SERIAL PRIMARY KEY,
  batch_id       TEXT REFERENCES batches,
  measured_at    TIMESTAMPTZ NOT NULL,
  metric         TEXT NOT NULL,       -- e.g. "viscosity", "moisture", "temperature"
  value          NUMERIC(10,4) NOT NULL,
  spec_min       NUMERIC(10,4),
  spec_max       NUMERIC(10,4)
);

CREATE TABLE operators (
  operator_id    TEXT PRIMARY KEY,
  name           TEXT NOT NULL,
  shift          TEXT NOT NULL        -- MORNING | AFTERNOON | NIGHT
);

CREATE INDEX idx_wo_work_center_status ON work_orders(work_center, status);
CREATE INDEX idx_batches_of ON batches(of_number);
CREATE INDEX idx_quality_batch ON quality_results(batch_id, measured_at);
```

### 6.3. Seed data

`sim/erp/seed.sql` — ~50 products, ~20 operators, ~200 historical OFs (last 30 days), 3 currently-RUNNING OFs on `machine1`/`machine2`/`machine3` to match the OPC-UA simulator's work centers.

### 6.4. ERP simulator binary — `cmd/erpsim`

A tiny Go binary that mimics ERP activity:

- Every **30 s** — pick one RUNNING OF, increment `actual_qty` by realistic amount (respecting product `target_rate`)
- Every **5 min** — 20% chance to finish a RUNNING OF (set `status=DONE`, `finished_at=now`) and start the next PLANNED one for that work_center
- Every **10 min** — insert a new `quality_results` row for each in-flight batch (with occasional out-of-spec value to trigger downstream detection)
- Every **1 h** — pre-plan a new OF for each work_center

This gives the OPC-UA simulator + pipeline engine live IT context to react to, so demos aren't static.

### 6.5. Bringing it up

```powershell
docker compose -f docker-compose.dev-erp.yml up -d       # Postgres + schema + seed
go run ./cmd/erpsim                                      # ERP activity simulator
.\run.ps1                                                # existing stack — server + agent + frontend
```

Add the fake ERP connection to `config/connections.yaml`:

```yaml
connections:
  - id: dev_erp
    name: "Dev ERP simulator"
    driver: postgres
    host: localhost
    port: 5432
    database: fake_erp
    username: mindset_readonly
    password_env: MINDSET_ERP_PASSWORD
    tls: false                   # local dev only
```

Then a demo pipeline in `config/pipelines/examples/of_enrichment.yaml` reads from it.

---

## 7. Implementation steps (ordered)

### Step 1 — pick drivers

Add to `go.mod`:
- `github.com/lib/pq` — Postgres
- `github.com/go-sql-driver/mysql` — MySQL/MariaDB
- `github.com/microsoft/go-mssqldb` — SQL Server
- SQLite: already using `modernc.org/sqlite`

All pure-Go — no CGO, cross-compiles cleanly to the edge box.

### Step 2 — new package: `internal/connections/`

```
internal/connections/
├── registry.go      # Registry holding open *sql.DB per connection ID
├── config.go        # Load config/connections.yaml
├── dsn.go           # Build driver-specific DSNs (each driver has different string format)
└── health.go        # Startup check: user has SELECT-only permissions
```

Exposes: `Get(id string) (*sql.DB, error)` — returns the pooled DB handle. Opens on first use, closed on shutdown.

### Step 3 — rewrite `internal/functions/connectors/sql_query.go`

Replace the demo stub:

```go
func handler(params map[string]interface{}) (interface{}, error) {
    connID, _ := params["connection_id"].(string)
    query, _ := params["query"].(string)
    queryParams, _ := params["params"].(map[string]interface{})
    timeout := durationOr(params["timeout_seconds"], 30*time.Second)
    limit := intOr(params["limit"], 1000)

    if err := validateReadOnly(query); err != nil { return nil, err }
    if err := enforceLimit(query, limit); err != nil { return nil, err }

    db, err := connections.Get(connID)
    if err != nil { return nil, err }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    started := time.Now()
    rows, err := db.QueryContext(ctx, resolvePlaceholders(query, queryParams), argsFrom(queryParams)...)
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

Helpers to implement:
- `validateReadOnly(sql string) error` — parses first non-whitespace token, must be `SELECT` (case-insensitive)
- `enforceLimit(sql string, max int) error` — rejects queries without a `LIMIT` clause; or auto-appends `LIMIT max`
- `resolvePlaceholders` — swaps `:name` → driver-specific placeholder syntax (`$1` for Postgres, `?` for MySQL/MSSQL/SQLite)
- `mapRows(*sql.Rows) ([]map[string]interface{}, error)` — column-type-aware, converts `sql.RawBytes` → `string`, `time.Time` → RFC3339, `[]byte` → `string` when text

### Step 4 — register in both binaries

Add to `buildRegistry()` in `cmd/server/main.go` AND `cmd/agent/main.go`:

```go
reg.Register(sqlquery.GetFunction())
```

### Step 5 — `/api/connections` REST endpoints

New handlers in `cmd/server/`:

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/api/connections` | List all defined connections (never returns passwords) |
| `POST` | `/api/connections` | Create/update — persisted to SQLite |
| `POST` | `/api/connections/{id}/test` | Ping the DB, return latency + row count on `SELECT 1` |
| `POST` | `/api/connections/{id}/preview` | Body: `{sql, params}` → returns 5 rows for the Pipeline Studio preview |
| `DELETE` | `/api/connections/{id}` | Remove |

Persistence: new `connections` table in `data/mindset.db` (auto-created on first use).

### Step 6 — Pipeline Studio SQL config panel

Frontend changes in `frontend/pipeline-builder/`:

- New page `src/pages/ConnectionsPage.jsx` — list + create/edit form + Test button
- In `BuilderPage.jsx`, when the user drags an `sql_query` node → show a config panel with:
  - **Connection** — dropdown of available connections
  - **Query** — Monaco editor (SQL syntax highlighting), placeholders as `:name`
  - **Params** — key/value grid, values can be static or `{{ trigger.field }}` template
  - **Timeout / limit** — number inputs with sensible defaults
  - **Preview** button — calls `/api/connections/{id}/preview`, shows first 5 rows in a mini-table
- Update `src/lib/functionDocs.js` — proper description + field help for `sql_query`
- Update `src/lib/functionDefaults.js` — default config seeded when the node is added

### Step 7 — simulation stack

- New `docker-compose.dev-erp.yml`
- `sim/erp/schema.sql` + `sim/erp/seed.sql`
- New binary `cmd/erpsim/main.go`
- New pipeline example `config/pipelines/examples/of_enrichment.yaml`

### Step 8 — integration tests

- `internal/connections/registry_test.go` — pool management, DSN parsing, config loading
- `internal/functions/connectors/sql_query_test.go` — using SQLite in-memory: happy path, timeout, limit, injection attempt rejected, non-SELECT rejected
- `internal/e2e/sql_pipeline_test.go` — spin up testcontainers Postgres, seed, run a real pipeline, assert enriched output

### Step 9 — docs update

- `docs/mindset.md` §8 — flag Module 4 as V1 in progress
- `docs/how_it_works.md` — new section on the SQL connector data path
- `docs/COMPONENTS.md` — add `internal/connections/` package
- `docs/ARCHITECTURE.md` — data-flow diagram: add ERP → connector arrow

### Step 10 — first customer smoke test

Once the simulator works end-to-end, get the beta cohort's first pharma or agrifood plant to:
1. Provision a read-only SQL user on a schema they consider safe
2. Give us a small `schedules` or `work_orders` extract
3. Wire ONE enrichment pipeline
4. Show it in the dashboard within 48h

That's the proof of concept for the paid engagement.

---

## 8. Common V1 use cases (example pipelines to ship)

### 8.1. Enrich a micro-stop event with current OF

**Trigger:** `mindset/events/status-change` (Run → Stopped)
**Pipeline:** `mqtt_subscribe → sql_query (current OF for that work_center) → merge → mqtt_publish (mindset/events/micro-stop-enriched)`
**Result:** every micro-stop event carries `of_number`, `product_code`, `operator`, planned throughput.

### 8.2. Load the shift schedule at start of shift

**Trigger:** cron / manual
**Pipeline:** `sql_query (SELECT FROM schedules WHERE work_center = ? AND planned_start BETWEEN ? AND ?) → add_to_dashboard`
**Result:** shift-view Gantt on the dashboard shows the planned work.

### 8.3. Fetch product master data on OF change

**Trigger:** `mindset/events/of-change` (detected upstream)
**Pipeline:** `mqtt_subscribe → sql_query (product master) → add_to_dashboard`
**Result:** operator sees target rate + recipe + hourly margin for the current product.

### 8.4. Cost-model input — hourly margin per product

**Used by:** `calculate_cost` function (see `docs/impact_engine.md`)
**Pipeline:** every event enrichment pulls `products.hourly_margin` so cost isn't a flat plant rate but reflects actual product economics.

---

## 9. UI changes (summary)

| Change | Where | Notes |
|---|---|---|
| New Connections page | `/connections` route | List, add, edit, delete, test. Matches the OPC-UA Connect page pattern. |
| SQL node config panel | Pipeline Studio | Connection dropdown + SQL editor + params + preview |
| Function meta update | `src/lib/functionMeta.js` | Icon: database. Color: teal (already partly set — verify). |
| Function docs update | `src/lib/functionDocs.js` | Description, per-field help text, working examples |
| Function defaults | `src/lib/functionDefaults.js` | Sensible starting query template so new nodes aren't empty |

---

## 10. Testing checklist

**Unit** (fast, no docker):
- [ ] `validateReadOnly` accepts SELECT, rejects INSERT/UPDATE/DELETE/DROP/... including whitespace + comment variants
- [ ] `enforceLimit` appends `LIMIT` when absent, rejects when explicit limit > max
- [ ] `resolvePlaceholders` handles all 4 driver syntaxes
- [ ] `mapRows` type conversions round-trip cleanly (int, float, bool, string, time, null)
- [ ] Connection registry pooling — same DSN returns same handle

**Integration** (testcontainers):
- [ ] Postgres: happy path returns rows
- [ ] Postgres: timeout kicks in and returns error
- [ ] Postgres: injection attempt via a parameter fails safely
- [ ] MySQL + MSSQL smoke tests

**End-to-end**:
- [ ] docker-compose ERP + `cmd/erpsim` + `run.ps1` → dashboard shows enriched micro-stops with `of_number` populated
- [ ] Preview button in Pipeline Studio returns rows within 2s
- [ ] Connections page: create → test → save → use in a pipeline round-trip works

---

## 11. Rollout

| Milestone | Target date | Blocker |
|---|---|---|
| Docker-compose ERP + erpsim + schema seed working | end of this week | none — pure infra |
| `sql_query` handler + registry + basic unit tests | +1 week | pick drivers, refactor connector package |
| `/api/connections` + Connections page | +1 week | UI polish |
| Integration tests green | +2 weeks | testcontainers-go setup |
| First beta customer wires 1 SQL enrichment | +3 weeks | customer scheduling |

Total: ~3 weeks from today to first customer smoke test.

---

## 12. Risks + known unknowns

| Risk | Mitigation |
|---|---|
| **SAP tables locked / obscured** — customer refuses direct SQL access | Deferred to MES-native V2 (Werum PAS-X / SAP MII). For V1, target ERPs that expose SQL (agrifood, metallurgy). |
| **Customer DB user has too many perms** — accidental writes possible | Startup health check verifies read-only. If not read-only, log warning + refuse to start `sql_query` handler. |
| **Timezone drift** — customer DB in UTC vs. plant local time | All timestamps stored as `TIMESTAMPTZ` (Postgres) / `datetimeoffset` (MSSQL). Convert to plant TZ only at UI render. |
| **Long queries block engine** | Mandatory `timeout_seconds`. Engine runs each node with `context.WithTimeout`. |
| **Connection storm** — pipeline fires many times/second | Per-connection pool cap. Pipeline-level rate limiting deferred to V1.5. |
| **Oracle demand** | Deferred to V1.5 with a CGO-tagged build variant OR the ODPI-C driver — decision when a customer requires it. |
| **Schema changes at customer site** | Query breakage. V1 mitigation: preview shows current shape. V1.5: connection health check runs a canary query nightly and flags on schema drift. |

---

## 13. What comes after SQL (queue)

1. **REST connectors** (V1.5) — modern MES/ERP APIs (SAP OData, Sage REST, Odoo)
2. **Pipeline suggestion engine** (V1.5) — DataOps Studio detects recurring patterns in the KG (repeat stops, quality clusters, supplier anomalies) and auto-proposes pipeline templates the user accepts with one click. Feature-borrow from LemonLime's *"self-creating automations"* — see Entry 67. Requires V1a+V1b populated first so there's data to reason over. Design constraint: precision > recall (hallucinated suggestions in industrial contexts cost money).
3. **Werum PAS-X connector** (V2) — pharma-mandatory; needs a Werum consultant + certification
4. **SAP MII / RFC connector** (V2) — for the largest SAP customers who won't grant SQL access
5. **AVEVA System Platform / OSIsoft PI connector** (V2) — historian access, complementary to real-time OPC-UA
6. **CDC / change-data-capture** (V2) — Debezium-style, when customers want "notify me when this table changes" instead of polling

---

## 14. TL;DR for the next 3 weeks

1. **Week 1** — Docker ERP simulator + schema + `cmd/erpsim`. Get the fake IT side alive.
2. **Week 2** — Replace `sql_query` stub with the real handler. Wire it through both binaries. Unit tests green.
3. **Week 3** — Connections page in the Pipeline Studio + one working demo pipeline (OF enrichment) + the first customer smoke test scheduled.

Ships in a month. Gets us to the first real OT/IT reconciliation demo — which is the moat.
