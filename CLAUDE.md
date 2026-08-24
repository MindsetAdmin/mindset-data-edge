# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

MindSet Data Edge is an industrial IoT edge platform with two Go binaries and a React frontend:

- **`cmd/server`** — HTTP API server (`:8080`). In the default mode it owns the OPC-UA session, manages live MQTT subscriptions, and exposes all REST + WebSocket endpoints for the UI.
- **`cmd/agent`** — Edge runtime: MQTT publishing, UNS contextualizer, rules engine, KG enrichment. Does NOT auto-connect to OPC-UA unless `opcua.auto_connect: true`.
- **`frontend/pipeline-builder`** — React/Vite app (`:5173`).

**Key coupling rule:** the frontend talks **only** to `cmd/server` via `/api`. The agent and server are loosely coupled — they never call each other directly; they share the MQTT broker and the SQLite database.

## Build and run

### External dependencies (not in this repo)

- **MQTT broker on `localhost:1883`** — required for almost everything (live tags, Run, auto-KG enrichment). Not bundled; quickest way to get one: `docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto:2 mosquitto -c /mosquitto-no-auth.conf`. `run.ps1` only warns if it's missing — it doesn't block startup.
- **OPC-UA source** — only needed for live tag demos. `config/agent.yaml`'s default endpoint (`opc.tcp://localhost:53530/OPCUA/Server1`) targets the free **Prosys OPC-UA Simulation Server**. Connect it from the UI (**Connecteurs → OPC-UA**, `/connect/opcua`), not from config — `opcua.auto_connect: false` means `cmd/server` owns the session via the UI. **`53530` can collide with a Windows/Hyper-V/WSL2-reserved dynamic port range** (`netsh interface ipv4 show excludedportrange protocol=tcp`) — if Prosys's UA TCP endpoint fails to start (its HTTPS endpoint working while TCP doesn't is the tell) but nothing shows as listening on the port, that's almost certainly why; move Prosys to a port outside all excluded ranges (e.g. `4840`, the IANA-standard OPC-UA port) and update `opcua.endpoint` to match (Entry 122). **`sim/opcua/mindset_simulation.xml`** — a NodeSet2 file (Objects menu → Import in Prosys) recreating the factory hierarchy this project's ISA-95 mapper expects (`Usine_Paris_Nord.Ligne1/2.Machine1/2/3.status`/`.temperature`, 4-level dot-notation — see `internal/uns/mapper.go`'s doc comment) — rebuild-from-scratch reference if Prosys's own simulation config is ever lost (Entry 123). After import, enable Prosys's own per-variable "Simulation" (not part of the standard NodeSet2 format) so values actually vary — `status` is the one Run/Stop detection keys off of.

### Full stack (Windows, recommended)
```powershell
.\run.ps1          # build both binaries, start server + agent + frontend, open browser
.\run.ps1 -NoAgent # skip the edge agent (when OPC-UA/MQTT aren't running)
.\run.ps1 -NoBuild # reuse existing bin/*.exe
```

### Manual build
```powershell
go build -o bin/server.exe ./cmd/server
go build -o bin/agent.exe  ./cmd/agent
```

### Run individually
```powershell
.\bin\server.exe   # API on :8080; flags: -config, -db, -pipelines, -addr
.\bin\agent.exe
```

### Frontend
```powershell
cd frontend/pipeline-builder
npm install        # first time only
npm run dev        # dev server on :5173
npm run build
npm run lint
```

### Tests
```powershell
go test ./...                                     # all packages (unit only)
go test ./internal/pipeline/...                   # single package
go test -run TestEngineName ./...                 # single test by name
go test -tags=integration ./internal/e2e/...       # sql_query vs a real MySQL testcontainer — needs Docker; skips cleanly if it's not running
```

### Fake ERP dev stack (`cmd/erpsim`)
Simulates a customer MySQL ERP so `sql_query` pipelines have realistic IT data to enrich against — see `config/pipelines/examples/of_enrichment.yaml`.
```powershell
docker compose -f docker-compose.dev.yml up -d       # mosquitto + MySQL (:3308, schema+seed from sim/erp/*.sql) + containerized erpsim
```
`docker-compose.dev.yml` runs all three dev dependencies (MQTT broker, fake-ERP MySQL, and the `erpsim` activity generator itself, built from `Dockerfile.erpsim`) in one network. To run `erpsim` locally instead of in the container (e.g. for faster iteration on `cmd/erpsim/main.go`):
```powershell
docker compose -f docker-compose.dev.yml up -d mosquitto mysql-erp   # skip the erpsim service
$env:MINDSET_ERP_PASSWORD = "readonly_dev"            # matches mindset_readonly's password in sim/erp/grant.mysql.sql — needed for the sql_query connector / config/connections.yaml's dev_erp entry
go run ./cmd/erpsim                                   # background loops: advance/rotate/quality/plan (see cmd/erpsim/main.go header)
```
Connects as `mindset_writer` (SELECT+INSERT+UPDATE, no DELETE) to `fake_erp`; `mindset_readonly` is the role intended for pipeline connectors. Both users are created by `sim/erp/grant.mysql.sql`. Configurable via `ERPSIM_DSN` and `ERPSIM_TICK_*` env vars. After the server is running, verify the connection from **Connectors → MySQL** (`/connectors/sql`) in the UI (**Test** on `dev_erp`) before running a pipeline that depends on it — `config/pipelines/examples/of_enrichment.yaml` is the shipped example (`work_center` param is a static value, edit it to `machine1`/`machine2`/`machine3` to match the seed data — no `{{ }}` templating exists yet, see Known limitations).

## Data flow

```
OPC-UA server
    ↓  (OPCUAManager in cmd/server, or discovery/ in agent when auto_connect=true)
MQTT  mindset/raw/#          ← {node_id, name, value, data_type, timestamp_ms}
    ↓  (UNS Contextualizer — agent only when auto_connect=true)
MQTT  mindset/site/<site>/<area>/<workCenter>/<tag>   ← ISA-95 enriched payload
    ├──► Rules Engine → mindset/events/status-change → mindset/events/micro-stop
    ├──► KG Subscriber → writes Equipment/Event/Cause/Cost nodes to data/mindset.db
    └──► LiveHub (cmd/server) → WebSocket /api/ws → React UI
```

In the default UI-controlled mode (`opcua.auto_connect: false`), `cmd/server` owns the OPC-UA session. The agent runs only the Rules Engine and KG Subscriber, consuming `mindset/site/#` published by the server.

## MQTT topic taxonomy

| Topic | Payload | Produced by | Consumed by |
|---|---|---|---|
| `mindset/raw/<nodeID>` | `{node_id,name,value,data_type,timestamp_ms}` | `mqtt/publisher.go` | UNS contextualizer, LiveHub |
| `mindset/site/<site>/<area>/<machine>/<tag>` | ISA-95 contextualized value | `uns/contextualizer.go` | Rules engine, LiveHub |
| `mindset/events/status-change` | transition + duration | `rules/engine.go` | micro-stop logic |
| `mindset/events/micro-stop` | `{work_center,duration_seconds,cause,cost_eur,...}` | rules / pipelines | KG subscriber, dashboards |

## Key packages

| Package | Role |
|---|---|
| `internal/config` | Loads `config/agent.yaml` into `Config` struct |
| `internal/pipeline` | `Pipeline` YAML model, `Engine` (topological execution), `Loader`, `Registry` |
| `internal/functions` | `Function` + `Registry`; sub-packages: `connectors/`, `transforms/`, `calculates/`, `conditions/`, `outputs/` |
| `internal/kg` | Unified SQLite-backed KG (`kg_nodes`/`kg_edges` tagged with a `category` column: `business` = site data, `platform` = pipeline topology) |
| `internal/rules` | Detects machine stop/start transitions from `mindset/site/#`, publishes to `mindset/events/` |
| `internal/uns` | `mapper.go` (tag→ISA-95 normalization) + `contextualizer.go` (raw→site republish) |
| `internal/mqtt` | `Publisher` wrapping paho MQTT client (`PublishRaw`, `PublishEvent`) |
| `internal/storage` | `SQLiteStore` using `modernc.org/sqlite` (pure-Go, no CGO) — auto-creates schema |
| `internal/connections` | `Registry` — pools `*sql.DB` per known connection (MySQL only in V1a); `Get` lazy-opens + verifies once, `Test` always re-verifies, `Add`/`Remove`/`List` manage definitions at runtime, `CloseAll` on shutdown |
| `internal/discovery` | `OPCUADiscovery` — browse + subscribe + publish raw |
| `cmd/server` | `OPCUAManager`, `LiveHub`, `TagRegistry` (persisted to SQLite), `TopicRegistry`, `StateTracker`, WebSocket hub, MCP server (`/mcp`, Track A — see "MCP server" below) |

## Function catalog

All functions must be registered in **both** `cmd/server/main.go#buildRegistry` and `cmd/agent/main.go`.

| Function | Type | Purpose |
|---|---|---|
| `opcua_read` | connector | Read a value from an OPC-UA node |
| `mqtt_subscribe` | connector | Subscribe to an MQTT topic (pipeline trigger) |
| `modbus_read` | connector | **Demo stub** — errors if executed |
| `sql_query` | connector | Read-only, parameterized SELECT (`:name` placeholders) against a connection from `internal/connections.Registry`; type-coerces results and, when `field_map`/`value_map` are configured, also returns a `canonical` copy per `docs/mysql_connector.md` §6b. Backed by `config/connections.yaml` + `POST /api/connections` |
| `state_machine` | transform | Detect Run↔Stop transitions |
| `filter` | transform | Keep/drop by condition |
| `calculate_duration` | calculate | Duration between start/stop timestamps |
| `calculate_cost` | calculate | € cost from duration × hourly rate; supports Manual/Config/Tag rate source + CSV/Excel per-product table |
| `threshold` | condition | Is value within [min, max]? |
| `add_to_dashboard` | output | Pin data onto the dashboard (`mindset/dashboard/<label>`) |
| `kg_save` | output | Save to KG — **not** in the server palette (KG enriches itself automatically via KGSubscriber) |

**Removed (Entry 119):** `uns_mapper` — deleted, not just deregistered. It was a thin wrapper around the exact ISA-95 mapping `OPCUAManager.route()` already does automatically for `isa95`/`both`-routed tags; keeping it meant a pipeline could hand-build, node by node, something the platform already does natively. `mqtt_publish` — also deleted. A pipeline no longer needs an explicit output node to publish to MQTT; see "Automatic output publishing" below.

Outputs are **sinks**: in the canvas they have an input port only, no output port. `add_to_dashboard` remains — it's a distinct, still-manual action (pinning to the UI dashboard), not superseded by the automatic MQTT publish below.

### Automatic output publishing (replaces the old `mqtt_publish` node, Entry 119)

A pipeline's declared `output` node (the terminal node — the one nothing else depends on; `pipelineMapping.js` computes this from the canvas automatically) has its result **auto-published to MQTT** after a successful run — `cmd/server/pipeline_output.go`'s `publishPipelineOutput`, called right after `Engine.Execute` in `handleRunPipeline`. No node, no manual wiring.

Topic resolution: the pipeline YAML's optional `output_topic` field, if set — required whenever the topic name is load-bearing (chained into another pipeline's trigger topic, or consumed by `internal/kg/subscriber.go`'s hardcoded `mindset/events/micro-stop` subscription) — otherwise auto-derived as `mindset/pipelines/<pipeline_id>/output`. All 3 shipped example pipelines (`microstop_detection`, `cost_calculation`, `of_enrichment`) set it explicitly, preserving their exact pre-existing topic names so the KG auto-enrichment chain and inter-pipeline chaining keep working unchanged.

## Pipeline system

Pipelines are YAML files in `config/pipelines/`. Each pipeline has typed `nodes` executed in dependency order (topological sort). A dependency that isn't a node ID (e.g. `"trigger"`) is treated as already satisfied.

```yaml
id: my_pipeline
name: My Pipeline
enabled: true
trigger:
  type: mqtt
nodes:
  - id: step1
    type: connector
    function: mqtt_subscribe
    config: { topic: "mindset/raw/#" }
  - id: step2
    type: transform
    function: state_machine
    depends_on: [step1]
```

Template pipelines live in `config/pipelines/examples/` (`microstop_detection`, `cost_calculation`, `of_enrichment` — served read-only via `GET /api/pipelines/examples`) and are copied into `config/pipelines/` (via `PipelinesPage`'s "load") once a user saves them as their own. (`opcua_to_uns` was removed in Entry 119 — it existed only to demonstrate the now-deleted `uns_mapper`, something `OPCUAManager.route()` already does automatically.)

The engine merges previous node outputs + trigger data + the node's own `config` map into one `params` map passed to the handler. A `recover()` in `callFunction` prevents a panicking handler from crashing the server.

## Knowledge Graph

One SQLite-backed graph (`internal/kg/graph.go`), with every node/edge tagged by a `category` column (2026-07-02 merge — `docs/analysis_log.md` Entry 50; see `internal/kg/types.go`):

- **`business`** — site/operational data: `Equipment`, `Event` (micro-stops), `Cause`, `Cost`, `Operator`, `Product`, `OF`. Auto-enriched by `KGSubscriber` from `mindset/events/micro-stop`.
- **`platform`** — pipeline topology: `Pipeline`, `Function`, `Connection`, `Topic`, `Dashboard`. Rebuilt by `KnowledgeGraph.RepopulatePlatform` (`internal/kg/builder.go` provides the in-memory build algorithm; the result is persisted, replacing the old 5-min in-memory cache). `RepopulatePlatform` no-ops if the pipeline registry's hash hasn't changed since the last rebuild. Empty until you save at least one pipeline (shipped examples in `config/pipelines/examples/` don't populate it).

`GetGraph(category)` (`"business"`, `"platform"`, or `"all"`) is the single read path, backing `GET /api/kg?category=`. `/api/kg/technical` and `/api/kg/domain` remain as legacy aliases (`handleTechnicalGraph`/`handleDomainGraph` in `cmd/server/main.go`) that call the same store under the old response shapes.

SQLite schema auto-created on `NewSQLiteStore`: `kg_nodes`, `kg_edges`, `events`, `tags`. Legacy DBs get `category` backfilled to `'business'` by a migration (`internal/storage/sqlite.go`).

### Structural bootstrap (OT auto-generation, `internal/kg/bootstrap.go`)

`GET /api/opcua/discover` no longer just lists tags — as a side effect (`OPCUAManager.seedKG` in `cmd/server/opcua.go`), it ISA-95-maps every discovered tag (`internal/uns.Mapper`) and calls `KnowledgeGraph.SeedFromDiscovery` to write the Site → Area → [WorkCenter] → Equipment → Tag skeleton straight into the `business` category — no waiting for a micro-stop event. `Equipment` node IDs (`machine_<name>`) match `AddMicroStop`'s scheme so both paths converge on the same node once real operational data arrives.

**Preview + correct before committing (Entry 124)**: the mapping + confidence, previously computed only for the KG-seed side effect, is now returned in the `/api/opcua/discover` response too (`discoveredTag.Site/Area/WorkCenter/WorkUnit/TagName/Confidence/Pending`) — `OpcuaTagSelector.jsx` shows it per tag with editable Area/WorkCenter/WorkUnit/TagName fields, pre-filled with the auto-guess. A correction submitted with `POST /api/opcua/subscribe` (optional `area`/`work_center`/`work_unit`/`tag_name` per selection) does two things, not just one: (1) `OPCUAManager.route()` — which recomputes the ISA-95 mapping on **every single live value**, not just once — checks a per-node `overrides` map first, so a correction actually changes what topic the tag's live data publishes to, not just a one-time KG write; (2) the corrected mapping is written to the KG with `confidence: 1.0` (human-confirmed), via the same `SeedFromDiscovery` path. **Known limitation, same class as the existing "reject doesn't cascade" one**: this is additive — the original auto-guessed node isn't retracted or deleted, it just sits alongside the corrected one in the KG (worse than the pending-node case if the original guess scored ≥ `AutoAcceptThreshold`, since it was never in the Pending list to reject from in the first place; no UI currently exists to clean up a stale-but-already-accepted node).

Each entry gets a 0.0–1.0 confidence score from two cheap heuristics computed across the whole discovered batch (not ML): does the tag's dot-depth match the batch's modal depth, and does its mapped name collide with another tag already claimed by the same equipment. Nodes scoring **≥ `kg.AutoAcceptThreshold` (0.7)** are written already confirmed (`pending: false`); nodes below it are flagged `pending: true` for human review — so a human only reviews the ones the heuristic itself is unsure about, not every generated node. `SeedFromDiscovery` is idempotent (`INSERT OR IGNORE`), so it's safe to call after every Discover.

**`OpcuaTagSelector.jsx` routing choice simplified to ISA-95 only (Entry 125)**: the Type column and the Brut/Les deux routing options were removed from display — a checkbox now selects ISA-95 publication per tag, replacing the 3-way radio group. Nothing changed about what raw storage actually does: any selected/monitored tag is still always raw-published by `discovery.Subscribe` regardless of this choice (that was always true — "raw" and "both" modes differed from "isa95" only in whether `route()` additionally published to `mindset/site/#`, and "isa95"/"both" were already functionally identical, so this simplification drops no real capability).

The UI surfaces this on `KnowledgeGraphPage` (`/kg`): a "Pending validation (N)" list with Accept/Reject per node, and pending nodes render with a dashed amber ring in the graph. `ValidateNode`/`RejectNode` back `POST /api/kg/pending/{id}/validate` and `.../reject`; `ListPending` backs `GET /api/kg/pending`. See `docs/analysis_log.md` Entries 87–98 and 107 for the design history (why WorkCenter/WorkUnit swap roles at depth ≥ 4, the live-Prosys-test bug that motivated the equipment-identity fix, and the confidence gate).

**The WorkCenter/WorkUnit depth-4 rule is now centralized, not re-derived per caller (Entry 127)**: `uns.UNSNode.EquipmentIdentity()` — WorkUnit if depth ≥ 4 and non-empty, else WorkCenter — is the single source of truth for "what's the machine's real identity," used by `internal/kg.SeedFromDiscovery` (via `kg.HierarchyEntry`, kept intentionally decoupled from `uns.UNSNode` — the KG package stays protocol-agnostic), `OPCUAManager.route()`, `OPCUAManager.computeMappings()`, and `OPCUAManager.SelectionsDetailed()`. Before this entry, only the KG-write path had the fix — `route()` (the live MQTT-publish path) and `SelectionsDetailed()` (backing `/api/opcua/selections` and `/api/machines`'s grouping) both used the raw `WorkCenter` field directly, so at 4-level tag names every machine on the same line silently shared one `StateTracker` entry (confirmed live: `Machine1` and `Machine2` under `Ligne1` merged into one Run/Stop state, and `kg_current_state("Machine1")` found nothing because the tracked key was actually `"Ligne1"`). `SelectionsDetailed()` also now checks a node's `overrides` entry (Entry 124) before falling back to the raw mapper — previously it always recomputed from scratch, so a manual correction wouldn't show up there even though it was already correctly affecting the live topic.

### IT-side structural bootstrap (`internal/connections/schema.go` + `canonical_suggest.go`, `internal/kg/it_bootstrap.go`)

The Track B counterpart to the OT bootstrap above (`docs/analysis_log.md` Entry 115, built in Entry 116) — same pattern, applied to SQL connections instead of OPC-UA. `GET /api/connections/{id}/discover` browses the connection's schema (`information_schema.columns` — MySQL only, matching `sql_query`'s V1a scope), then scores every table against 2 canonical types (`work_order`, `product` — deliberately not the full 9-object set from Entry 92 yet) by column-name synonym matching: core fields (id/status/product/work_center references) are worth 80% of the confidence score, bonus fields (customer/due-date/margin) 20% — mid-market ERPs often lack the bonus fields, so a table shouldn't be penalized out of usefulness for missing them. Tables scoring below `suggestionFloor` (0.5) against every type aren't suggested at all.

Results are written as a **new KG node type**, `SchemaMapping` (business category), through `KnowledgeGraph.SeedSchemaMappings` — gated by the exact same `kg.AutoAcceptThreshold` and surfaced through the exact same `ListPending`/`ValidateNode`/`RejectNode`/`KnowledgeGraphPage` pending-list UI as OT nodes, unchanged. Live-verified against the real `dev_erp` fake-ERP connection (`sim/erp/schema.mysql.sql`): `products` → confidence 1.0 (auto-accepted), `work_orders` → confidence 0.8 (auto-accepted), and `operators`/`batches`/`schedules`/`quality_results` correctly excluded (all below the floor).

**Phase 4 built**: `cmd/server/active_production.go`'s `ActiveProduction` queries every validated `work_order` mapping for orders whose status matches a hardcoded in-progress token set, exposed both as the MCP tool `kg_active_production` and as `GET /api/production/active` (Entry 120 added the REST route). Live-verified against real erpsim-seeded data: 3 active orders returned across `machine1`/`machine2`/`machine3`. **Still not built**: retroactive event-tagging, so "how long did product B take yesterday" remains unanswerable — `kg_active_production`/`/api/production/active` only ever reflect the current moment.

**Second priority axis — deadline urgency (Entry 133)**: alongside `kg_cost_summary`'s cost ranking, `ActiveProductionFact` optionally carries `due_date`/`customer_id`/`days_until_due` — sourced from `due_date`/`customer_id` as *bonus* fields in the `work_order` canonical mapping (already scored by `internal/connections/canonical_suggest.go` since Track B's original design; the fake ERP's schema just didn't have matching columns until this entry). Absent when a mapping doesn't resolve them — never guessed. Cost priority and deadline priority are genuinely different axes and can disagree (a cheap stop on a today-due order can matter more than an expensive one with no deadline pressure) — the demo's seed data deliberately makes them disagree, so both tools stay meaningfully distinct rather than one being a subset of the other.

**Entity resolution (Entry 120)**: `cmd/server/entity_resolution.go`'s `ResolveWorkCenters` — runs automatically as part of every `/discover`, right after mapping seeding. For every validated `work_order` mapping, queries the real distinct `work_center` values present in the ERP table and matches each (case-insensitive, exact — not fuzzy) against known OT `Equipment` nodes' `work_center` property. Where they match, writes a persisted `same_as` edge (`Equipment → SchemaMapping`, business category). This closes the gap Entry 109 flagged as entirely missing ("nothing computes the `same_as` OT↔IT match"). Live-verified: 3 `same_as` edges created against the real `dev_erp` data (`machine1`/`machine2`/`machine3` all resolved to their `machine_Machine<N>` Equipment nodes). Each `ActiveProductionFact` also carries the resolved `equipment_id` directly (empty string if no OT node matched — never guessed), computed via the same normalized matching, independent of whether `ResolveWorkCenters` has run — so the live answer and the persisted edge can't drift out of sync with each other.

**Surfaced on the Dashboard (Entry 120)**: `DashboardPage.jsx` now fetches `/api/production/active` and renders a "🔗 Production active (ERP)" panel (machine / product / OF / status / OT-link badge) — only when there's data, so it's invisible until a work-order mapping exists. Previously this data was reachable only via MCP, with zero UI surface.

**Browse everything on connect (Entry 126)**: `internal/connections.ListDatabasesAndTables` — every database + table + column visible to a connection's user, in one `information_schema` pass, deliberately separate from `DiscoverSchema`/`/discover`'s single-database-scoped canonical-mapping flow (that one stays tied to the connection's own configured `database`, since `SchemaMapping` nodes assume one connection = one database). Scoped automatically by the account's real MySQL grants — `information_schema.schemata` only lists databases the user has privileges on, so this can't leak visibility the account doesn't already have. `SqlConnectionsPage.jsx`'s "Connecter" button (already the existing Test action) now also populates this tree automatically on success — click-to-expand per table shows its columns with the primary key flagged.

| Table | Key columns | Written by |
|---|---|---|
| `kg_nodes` | id, category, type, label, properties(JSON), created_at | KG subscriber / graph |
| `kg_edges` | id, category, from_id, to_id, relation, weight, created_at | KG subscriber / graph |
| `tags` | node_id, name, value(JSON), data_type, timestamp_ms | server TagRegistry |
| `events` | id, type, work_center, duration_seconds, cause, cost_eur | (reserved) |

## Configuration

`config/agent.yaml` — the server falls back to defaults if the file is absent:

- `opcua.auto_connect: false` — default; OPC-UA driven from UI via `cmd/server`. Set `true` to restore legacy agent-owned auto-discovery.
- `opcua.endpoint` — e.g. `opc.tcp://localhost:53530/OPCUA/SimulationServer`
- `mqtt.broker` — default `tcp://localhost:1883`
- `cost.hourly_cost` — default `85.0` EUR/h

`config/connections.yaml` — SQL connection definitions for `internal/connections`, loaded by both `cmd/server` and `cmd/agent` on startup (missing file → empty connection set, not fatal). Passwords are never inlined; each entry's `password_env` names an env var resolved at connection-open time (e.g. the shipped `dev_erp` entry expects `MINDSET_ERP_PASSWORD`, matching the `mindset_readonly` user created by `sim/erp/grant.mysql.sql`). Connections created via `POST /api/connections` persist to the `connections` table in `data/mindset.db` and are re-loaded into the registry on the next startup, taking precedence over a YAML entry with the same id.

## API surface (`cmd/server`)

All routes under `/api/`. CORS is open (`*`).

| Method & path | Returns |
|---|---|
| `GET /api/health` | `{status:"ok"}` |
| `GET /api/functions[?type=]` | function catalog, optionally filtered by type |
| `GET /api/connectors` | function catalog filtered to `type=connector` (thin alias) |
| `GET/POST /api/pipelines` | list pipelines / save as YAML |
| `GET /api/pipelines/examples` | template pipelines from `config/pipelines/examples/` |
| `POST /api/pipelines/{id}/run` | execute a pipeline; returns per-node `ExecutionResult`. Optional body `{"trigger_data": {...}}` (Entry 131) supplies what a live trigger message would have carried — needed because no pipeline auto-fires on real MQTT today (see Known limitations); an empty/absent body preserves the old no-trigger-data behavior |
| `DELETE /api/pipelines/{id}` | remove a saved pipeline's YAML file |
| `GET /api/tags` | live OPC-UA tags + values (SQLite-persisted) |
| `GET /api/machines` | tags grouped by work center + live Running/Stopped state |
| `GET /api/topics` | live topics + msg/s + category + broker_connected |
| `GET /api/config` | safe subset of `agent.yaml` (opcua, mqtt, cost, site) |
| `POST /api/opcua/connect` | dynamically connect to a user-specified OPC-UA endpoint |
| `GET /api/opcua/discover` | browse the connected server's node tree; each tag now also carries its auto-computed ISA-95 mapping + confidence (Entry 124), previously computed but never returned |
| `POST /api/opcua/subscribe` | monitor selected tags with per-tag mode (`raw`\|`isa95`\|`both`); each selection may also carry an optional `area`/`work_center`/`work_unit`/`tag_name` correction (Entry 124) — see "Structural bootstrap" below |
| `POST /api/opcua/disconnect` | close the OPC-UA session |
| `GET /api/opcua/status` | connection status |
| `GET /api/opcua/selections` | per-tag routing + ISA-95 mapping |
| `GET /api/dashboard/pins` | current `add_to_dashboard` pins |
| `GET/POST /api/connections` | list SQL connections (never returns passwords) / create-or-replace one |
| `POST /api/connections/{id}/test` | force a fresh read-only health check; `{ok, latency_ms, read_only, error?}` |
| `POST /api/connections/{id}/preview` | run a query through the same guards as `sql_query`, capped at 5 rows |
| `GET /api/connections/{id}/discover` | browse the connection's schema (`information_schema`) and auto-suggest canonical mappings (`work_order`/`product`) into the KG as confidence-gated `SchemaMapping` nodes, then resolve any validated `work_order` mapping's work centers against OT `Equipment` nodes — the IT-side analog of `/api/opcua/discover`; see "IT-side structural bootstrap" below |
| `GET /api/connections/{id}/databases` | browse every database + table (with columns) visible to this connection's user in one call (Entry 126) — pure read-only visibility, no KG side effect, scoped by the account's real MySQL grants |
| `GET /api/production/active[?work_center=]` | live active production order + product per machine, from a validated ERP work-order mapping; each entry carries `equipment_id` when resolved against an OT node |
| `DELETE /api/connections/{id}` | remove a connection (registry + `data/mindset.db`) |
| `GET /api/kg?category=business\|platform\|all` | unified KG read (default `all`); the frontend's `client.js` calls this |
| `GET /api/kg/technical` | legacy alias → `category=platform`, old `TechnicalGraph` response shape |
| `GET /api/kg/domain` | legacy alias → `category=business` |
| `GET /api/kg/pending` | business-category nodes still awaiting human validation (structural bootstrap) |
| `POST /api/kg/pending/{id}/validate` | confirm an auto-generated node |
| `POST /api/kg/pending/{id}/reject` | discard an auto-generated node (deletes it + its direct edges) |
| `GET /api/stats` | counts + micro-stops/downtime/cost + uptime + broker status |
| `WS /api/ws` | WebSocket live push: `{type:"tag"\|"state"\|"event"\|"dashboard"}` |
| `POST /mcp` | MCP server (Streamable HTTP, stateless) — see "MCP server" below |

## MCP server (agent tool access)

`cmd/server/mcp_server.go` mounts a read-only MCP server at `/mcp` (`github.com/modelcontextprotocol/go-sdk`, Streamable HTTP transport, stateless — every tool call is self-contained). Started as **Track A** (`docs/analysis_log.md` Entry 113/114): 4 tools over data that already worked end to end from day one. Entry 117 added a 5th tool, `kg_active_production` (Track B Phase 4), once the IT-side structural bootstrap made a live product/work-order signal available. There is still **no tool for historical product-scoped questions** ("how long did product B run yesterday") — that needs retroactive event-tagging, which doesn't exist, so no tool claims to answer it.

**Two transports (Entry 121)**: the default is HTTP at `/mcp` (`mountMCP`), for remote/URL-based clients — but those generally require `https://`, which a local dev server doesn't have. For local clients that launch the server as a subprocess instead (Claude Desktop's `mcpServers` config), run the same binary with **`-mcp-stdio`**: no HTTP listener is started at all (no port bound, so it runs alongside the normal `:8080` instance with zero conflict — a distinct MQTT client ID, `mindset-mcp-stdio` vs `mindset-api-server`, prevents the two processes from kicking each other off the broker), and the identical tool set is served over stdin/stdout (`runMCPStdio`, `mcp.StdioTransport`). All logging already goes to `log.Printf` (stderr by default in Go) — verified live that stdout carries nothing but protocol JSON in stdio mode, since any stray print there would corrupt the stream.

| Tool | Wraps |
|---|---|
| `kg_query_events` | `internal/kg.QueryEvents` — Event nodes joined with Cause/Cost, filtered by work_center/cause/time window |
| `kg_cost_summary` | `internal/kg.CostSummary` — same data, aggregated by cause or work_center. When grouped by work_center (default), merged server-side (Entry 134) with `kg_active_production`'s deadline data: a machine due within 7 days ranks ahead of costlier-but-not-urgent ones, since a missed deadline has real cost even though it's flagged rather than priced (no penalty-clause data to price it from) |
| `kg_current_state` | `StateTracker` (`cmd/server/live.go`) — Running/Stopped per machine |
| `kg_describe_node` | `KnowledgeGraph.GetGraph("all")` — any node + its direct edges, either direction |
| `kg_active_production` | `cmd/server/active_production.go`'s `ActiveProduction` — live query against a human-validated `work_order` `SchemaMapping`; answers "what's running now," not historical duration; each fact includes `equipment_id` when resolved against an OT node (Entry 120) |

The query logic (`internal/kg/query.go`, `active_production.go`) is transport-agnostic on purpose — reusable outside MCP (e.g. a future dashboard widget) and independently testable.

**Connecting Claude Desktop (stdio):** Claude Desktop's custom-connector UI validates for `https://`, which a local dev server doesn't have — use the stdio transport instead (Settings → Developer → "Local MCP servers" on the packaged Windows build; **not** the web-only "Connectors" list, which is OAuth/HTTPS-only and a different feature entirely). Prerequisites: `bin/server.exe` built, MQTT broker + `mindset-erp` MySQL containers up (`docker compose -f docker-compose.dev.yml up -d`).

**All 4 path flags must be absolute** (Entry 122) — `-config`, `-db`, `-pipelines`, and `-connections` (new flag; was a hardcoded `"config/connections.yaml"` string until this entry, unfixable via config alone). Their defaults are relative to the current working directory, which is fine when a shell has already `cd`'d into the project root (`run.ps1`, `claude mcp add` run from the repo) — but Claude Desktop launches the subprocess from its own working directory (confirmed: `C:\Windows\System32` in one live test), so relative defaults silently fail (`open config/agent.yaml: The system cannot find the path specified.` → `log.Fatalf` → the process exits immediately → Claude Desktop reports **"Server disconnected"**, which gives no hint the actual cause is a path problem):
```json
{
  "mcpServers": {
    "mindset-data": {
      "command": "C:\\path\\to\\mindset-data-edge\\bin\\server.exe",
      "args": [
        "-mcp-stdio",
        "-config", "C:\\path\\to\\mindset-data-edge\\config\\agent.yaml",
        "-db", "C:\\path\\to\\mindset-data-edge\\data\\mindset.db",
        "-pipelines", "C:\\path\\to\\mindset-data-edge\\config\\pipelines",
        "-connections", "C:\\path\\to\\mindset-data-edge\\config\\connections.yaml"
      ],
      "env": { "MINDSET_ERP_PASSWORD": "readonly_dev" }
    }
  }
}
```
Fully quit Claude Desktop (not just close the window — on the Windows Store build, background windows/tray icon can survive a plain close) and reopen it; config is only read at startup. Verify via Settings → Developer → Local MCP servers that `mindset-data` shows connected (not "failed"), then ask a real question in a new conversation (e.g. *"which product is running on machine1 right now?"*); approve the tool-use prompt when it appears. This spawns a separate process from any already-running `:8080` HTTP server — distinct MQTT client ID, no port bound, no conflict (see "Two transports" above).

## Frontend

React 19 + Vite + Tailwind. State managed with Zustand (`src/store/studioStore.js`). Vite proxies all `/api` requests (including WebSocket upgrades) to `http://localhost:8080`.

**i18n (Entry 129):** `react-i18next`, initialized in `src/i18n.js` from `src/locales/{en,fr}.json` (default `fr`, persisted per-user in `localStorage['mindset_lang']`). A FR/EN toggle lives in `NavBar.jsx`. Components use `useTranslation()` + `t('namespace.key')`; the couple of plain (non-component) modules that need strings — `functionDocs.js`, `functionMeta.js`, `dashboardData.js` — import the `i18n` singleton and call `.t()` directly instead. Two things are deliberately **not** wired to the toggle: the three canvas zone labels (`ENTRÉE`/`CŒUR`/`SORTIE`, `pipelineMapping.js`'s `ZONES` constant) are fixed product jargon baked in at canvas-node-creation time in a non-reactive module; and function *descriptions* shown in the Palette/NodeConfigPanel come from the Go backend's function catalog (`GET /api/functions`), not the frontend, so translating them is a backend task.

### Pages (`src/pages/`)

| Page | Route | Purpose |
|---|---|---|
| `OverviewPage` | `/overview` | Key stats + quick links |
| `ConnectorsPage` | `/connectors` | Connector gallery — tiles for every catalog connector, `implemented: true` ones link out, the rest shown honestly as roadmap (`docs/mindset.md` §5) |
| `OpcuaConnectPage` | `/connect/opcua` | Dynamic OPC-UA connect/browse/subscribe UI; Discover also triggers the KG structural bootstrap |
| `MqttConnectPage` | `/connectors/mqtt` | MQTT connector config |
| `SqlConnectionsPage` | `/connectors/sql` | SQL connections: list + create + Test, backs the `sql_query` connector |
| `BuilderPage` | `/compose` | Drag-and-drop pipeline builder (ReactFlow); ENTRÉE/CŒUR/SORTIE bands, guided config panel, Save→YAML, Run, delete |
| `PipelinesPage` | `/pipelines` | Your pipelines (run/load) + templates (load) |
| `DashboardPage` | `/dashboards` | Real-time ops dashboard, WebSocket-driven (20s fallback); KPIs, pinned widgets, live chart, machine status, Gantt |
| `KnowledgeGraphPage` | `/kg` | `ForceGraph.jsx` viewer (not Cytoscape — `CytoscapeGraph.jsx` exists in the repo but is dead code, unused; corrected 2026-07-28, see `analysis_log.md` Entry 140), Technique/Domaine toggle + type filters; **Pending validation** list (Accept/Reject) for structural-bootstrap nodes, rendered with a dashed amber ring until confirmed |

### `src/lib/` helpers

| File | Purpose |
|---|---|
| `pipelineMapping.js` | Convert canvas ↔ backend `Pipeline` (zones, trigger, depends_on) |
| `pipelineLoading.js` | Load a saved pipeline as a chain-only flow — keeps node types/edges, resets each node's config to its default so the user reconfigures before saving |
| `functionMeta.js` | Icons/colors/categories per function type |
| `functionDocs.js` | Per-function description + field labels/help/examples |
| `functionDefaults.js` | Default config seeded when a function node is added |
| `connectorTemplates.js` | Default config per connector + trigger type |
| `kgGraph.js` | **Dead code** — was the Cytoscape mapping helper; unused since the KG viewer moved to `ForceGraph.jsx` (which handles its own KG JSON → graph-elements transform internally via `typesPresent`/`NODE_COLORS`/`FALLBACK_COLOR`). Not imported anywhere. Corrected 2026-07-28, `analysis_log.md` Entry 140. |
| `dashboardData.js` | Join domain graph → events (cause/cost), today/yesterday |
| `useLiveSocket.js` | WebSocket hook (auto-reconnect) for live push |

`src/api/client.js` — all `fetch` calls to `/api/*`.

## Adding a new function

1. Create a file under the appropriate sub-package (e.g. `internal/functions/transforms/my_transform.go`)
2. Implement a struct with `GetFunction() *functions.Function` returning a `*functions.Function` with `Name`, `Type`, `Description`, `Inputs`, `Outputs`, and `Handler func(map[string]interface{}) (interface{}, error)`
3. Register it in `buildRegistry()` in `cmd/server/main.go` **and** in `cmd/agent/main.go`

## Known limitations

- **No pipeline is ever auto-triggered by a live MQTT message (Entry 131).** A pipeline YAML's `trigger: mqtt_subscribe` block is declarative only (drives the Compose UI's ENTRÉE zone display) — `pipeline.Engine.Execute` is called from exactly one place in the codebase, `handleRunPipeline` (the manual Run button/API), always with the trigger data the caller supplies (empty by default; see the `/api/pipelines/{id}/run` row above for the optional `trigger_data` body). The rules engine (`internal/rules/engine.go`) publishing `mindset/events/status-change` does NOT cause any pipeline to run — nothing currently subscribes that topic and calls Execute. `config/pipelines/pipeline_microstop_detection.yaml` was rebuilt (Entry 131) to be a correct target for that wiring once it exists (trigger → threshold → calculate_cost, all fields it needs already present in a status-change payload) but still requires either an automatic trigger dispatcher (not built) or a manual/scripted Run call with the right `trigger_data` in the meantime.
- **KG subscribers need distinct MQTT client IDs, or they silently evict each other (Entry 131, fixed).** `cmd/server` and `cmd/agent` each run their own `KGSubscriber`; `internal/kg.NewKGSubscriber` now requires a `clientID` argument for exactly this reason — don't reintroduce a shared literal.
- **Two OS processes writing the same SQLite file need a busy timeout (Entry 131, fixed).** `internal/storage.NewSQLiteStore` sets `PRAGMA busy_timeout = 5000` right after opening — without it, `cmd/server` and `cmd/agent` racing to write the KG at the same moment silently drop writes via `SQLITE_BUSY` (subscriber.go logs the error but never retries).
- Secure OPC-UA modes (`Sign`/`SignAndEncrypt`) need a client certificate not yet wired — use `None`.
- The OPC-UA session holds one subscription at a time; changing tag selections requires a disconnect + reconnect cycle.
- `modbus_read` is a metadata-only stub that errors if executed.
- `sql_query` only supports `driver: mysql` (V1a); no `system_profile`/`query_name` resolution yet (system profiles are V1c — see `docs/mysql_connector.md` §6b). No `{{ trigger.field }}`-style templating for its `params` values either — despite the doc's own examples showing that syntax, `internal/pipeline` has no `{{ }}` handling; params are static values only (see `docs/analysis_log.md` Entry 79). Verified end-to-end against a real MySQL container (Entry 82): `go test -tags=integration ./internal/e2e/...`, `/api/connections/{id}/test`+`/preview`, and a full `of_enrichment.yaml` run all pass. The `SqlConnectionsPage`/`SqlConfigPanel`/`FieldMapEditor` frontend is still only lint/build-verified, not click-tested — no browser available in this environment.
- Local dev machines sometimes have something else already bound to whatever port `docker-compose.dev.yml` maps MySQL to (3306 conflicts with a Windows `mysqld.exe` service; on at least one dev machine 3307 collided with a XAMPP install too — hence 3308 as of Entry 82). If `/api/connections/*/test` fails with an auth error despite correct credentials, check `docker port mindset-erp` actually shows the mapping you expect, and that nothing else owns that host port, before assuming it's a code bug.
- The KG's `platform` category is empty until you save at least one pipeline; shipped example pipelines are templates only and don't populate it.
- **The structural bootstrap (see above) is v0 and OT-only.** Verified against a real Prosys server only at small scale (8 tags / 15 pending nodes) — the flat accept/reject list is untested past that and likely needs a grouped/tree view for a real industrial server's hundreds of tags. Rejecting a Site or Area doesn't cascade to its children (they become orphaned-but-still-pending, individually rejectable). If the bootstrap seeds an Equipment node before a real micro-stop references the same machine, it stays `pending` even after real data starts flowing, until someone explicitly validates it (`AddNodeCat` is `INSERT OR IGNORE`, so whichever path writes first wins the row). **IT-side (ERP) master data now has an equivalent auto-generation (Entry 115/116, see "IT-side structural bootstrap" above)** — schema discovery + confidence-gated `SchemaMapping` suggestion, live-verified against the real `dev_erp` connection, and a validated mapping is now actually consumed (`kg_active_production` MCP tool, Phase 4, also live-verified). Still narrower than the OT side: only 2 canonical types (`work_order`, `product`), and no retroactive event-tagging — "what's running now" works, "how long did product X run yesterday" still doesn't.
- **The `sql_query` "canonical model" (`docs/mysql_connector.md` §6b) isn't consumed by anything yet, despite the doc saying it is.** The doc states *"Downstream nodes (rules, KG subscriber, Impact Engine, MCP) consume canonical"* — that's a design intent, not current behavior. `internal/kg/subscriber.go` subscribes to exactly one MQTT topic (`mindset/events/micro-stop`); nothing reads `canonical` output. Same missing bridge as the OPC-UA/KG gap above — see `docs/analysis_log.md` Entry 90. **Object set (decided Entry 92):** `WorkOrder`/`Batch`/`Product`/`Schedule`/`Quality`/`Operator`/`Material`/`Asset`/`ProcessSegment` — deliberately aligned to ISA-95 Part 2's information model. `Asset` specifically is meant to eventually entity-resolve against the KG's OT-derived `Equipment` node type (see the bootstrapping bullet above), not just a 9th flat object. **Locked external language:** "aligns with ISA-95's information model" — never "ISA-95-compliant" (that implies B2MML/wire-format conformance, explicitly out of scope).
- Dashboard "vs hier" deltas need events spanning two days; TRS is an availability proxy, not true OEE.
