# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

MindSet Data Edge is an industrial IoT edge platform with two Go binaries and a React frontend:

- **`cmd/server`** — HTTP API server (`:8080`). In the default mode it owns the OPC-UA session, manages live MQTT subscriptions, and exposes all REST + WebSocket endpoints for the UI.
- **`cmd/agent`** — Edge runtime: MQTT publishing, UNS contextualizer, rules engine, KG enrichment. Does NOT auto-connect to OPC-UA unless `opcua.auto_connect: true`.
- **`frontend/pipeline-builder`** — React/Vite app (`:5173`).

**Key coupling rule:** the frontend talks **only** to `cmd/server` via `/api`. The agent and server are loosely coupled — they never call each other directly; they share the MQTT broker and the SQLite database.

## Build and run

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
go test ./...                           # all packages
go test ./internal/pipeline/...         # single package
go test -run TestEngineName ./...       # single test by name
```

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
| `internal/kg` | SQLite-backed domain KG (micro-stops) + in-memory technical KG (pipeline topology) |
| `internal/rules` | Detects machine stop/start transitions from `mindset/site/#`, publishes to `mindset/events/` |
| `internal/uns` | `mapper.go` (tag→ISA-95 normalization) + `contextualizer.go` (raw→site republish) |
| `internal/mqtt` | `Publisher` wrapping paho MQTT client (`PublishRaw`, `PublishEvent`) |
| `internal/storage` | `SQLiteStore` using `modernc.org/sqlite` (pure-Go, no CGO) — auto-creates schema |
| `internal/discovery` | `OPCUADiscovery` — browse + subscribe + publish raw |
| `cmd/server` | `OPCUAManager`, `LiveHub`, `TagRegistry` (persisted to SQLite), `TopicRegistry`, `StateTracker`, WebSocket hub |

## Function catalog

All functions must be registered in **both** `cmd/server/main.go#buildRegistry` and `cmd/agent/main.go`.

| Function | Type | Purpose |
|---|---|---|
| `opcua_read` | connector | Read a value from an OPC-UA node |
| `mqtt_subscribe` | connector | Subscribe to an MQTT topic (pipeline trigger) |
| `modbus_read`, `sql_query` | connector | **Demo stubs** — error if executed |
| `state_machine` | transform | Detect Run↔Stop transitions |
| `uns_mapper` | transform | Normalize a tag to ISA-95 structure |
| `filter` | transform | Keep/drop by condition |
| `calculate_duration` | calculate | Duration between start/stop timestamps |
| `calculate_cost` | calculate | € cost from duration × hourly rate; supports Manual/Config/Tag rate source + CSV/Excel per-product table |
| `threshold` | condition | Is value within [min, max]? |
| `mqtt_publish` | output | Publish to an MQTT topic |
| `add_to_dashboard` | output | Pin data onto the dashboard (`mindset/dashboard/<label>`) |
| `kg_save` | output | Save to KG — **not** in the server palette (KG enriches itself automatically via KGSubscriber) |

Outputs are **sinks**: in the canvas they have an input port only, no output port.

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

Built-in example pipelines (in `config/pipelines/`): `microstop_detection`, `opcua_to_uns`, `cost_calculation`.

The engine merges previous node outputs + trigger data + the node's own `config` map into one `params` map passed to the handler. A `recover()` in `callFunction` prevents a panicking handler from crashing the server.

## Knowledge Graph

Two distinct graphs:

- **Domain KG** (`internal/kg/graph.go`, `/api/kg/domain`): persisted in SQLite, stores `Equipment`, `Event` (micro-stops), `Cost`, and `Cause` nodes. Auto-enriched by `KGSubscriber` from `mindset/events/micro-stop`.
- **Technical KG** (`internal/kg/builder.go`, `/api/kg/technical`): computed in-memory from the pipeline registry — the architecture view (Connections, Topics, Pipelines, Dashboards). Cached 5 minutes, cache-busted by a registry hash. Empty until you create a pipeline (shipped examples don't appear here).

SQLite schema auto-created on `NewSQLiteStore`: `kg_nodes`, `kg_edges`, `events`, `tags`.

| Table | Key columns | Written by |
|---|---|---|
| `kg_nodes` | id, type, label, properties(JSON) | KG subscriber / graph |
| `kg_edges` | id, from_id, to_id, relation, weight | KG subscriber / graph |
| `tags` | node_id, name, value(JSON), data_type, timestamp_ms | server TagRegistry |
| `events` | id, type, work_center, duration_seconds, cause, cost_eur | (reserved) |

## Configuration

`config/agent.yaml` — the server falls back to defaults if the file is absent:

- `opcua.auto_connect: false` — default; OPC-UA driven from UI via `cmd/server`. Set `true` to restore legacy agent-owned auto-discovery.
- `opcua.endpoint` — e.g. `opc.tcp://localhost:53530/OPCUA/SimulationServer`
- `mqtt.broker` — default `tcp://localhost:1883`
- `cost.hourly_cost` — default `85.0` EUR/h

## API surface (`cmd/server`)

All routes under `/api/`. CORS is open (`*`).

| Method & path | Returns |
|---|---|
| `GET /api/health` | `{status:"ok"}` |
| `GET /api/functions[?type=]` | function catalog, optionally filtered by type |
| `GET/POST /api/pipelines` | list pipelines / save as YAML |
| `GET /api/pipelines/examples` | template pipelines from `config/pipelines/examples/` |
| `POST /api/pipelines/{id}/run` | execute a pipeline; returns per-node `ExecutionResult` |
| `GET /api/tags` | live OPC-UA tags + values (SQLite-persisted) |
| `GET /api/machines` | tags grouped by work center + live Running/Stopped state |
| `GET /api/topics` | live topics + msg/s + category + broker_connected |
| `GET /api/config` | safe subset of `agent.yaml` (opcua, mqtt, cost, site) |
| `POST /api/opcua/connect` | dynamically connect to a user-specified OPC-UA endpoint |
| `GET /api/opcua/discover` | browse the connected server's node tree |
| `POST /api/opcua/subscribe` | monitor selected tags with per-tag mode (`raw`\|`isa95`\|`both`) |
| `POST /api/opcua/disconnect` | close the OPC-UA session |
| `GET /api/opcua/status` | connection status |
| `GET /api/opcua/selections` | per-tag routing + ISA-95 mapping |
| `GET /api/dashboard/pins` | current `add_to_dashboard` pins |
| `GET /api/kg/technical` | technical (architecture) graph for Cytoscape |
| `GET /api/kg/domain` | domain (data) graph for Cytoscape |
| `GET /api/stats` | counts + micro-stops/downtime/cost + uptime + broker status |
| `WS /api/ws` | WebSocket live push: `{type:"tag"\|"state"\|"event"\|"dashboard"}` |

## Frontend

React 19 + Vite + Tailwind. State managed with Zustand (`src/store/studioStore.js`). Vite proxies all `/api` requests (including WebSocket upgrades) to `http://localhost:8080`.

### Pages (`src/pages/`)

| Page | Route | Purpose |
|---|---|---|
| `OverviewPage` | `/overview` | Key stats + quick links |
| `ConnectPage` | `/connect` | Pick a connector → applies to the pipeline trigger |
| `BuilderPage` | `/compose` | Drag-and-drop pipeline builder (ReactFlow); ENTRÉE/CŒUR/SORTIE bands, guided config panel, Save→YAML, Run, delete |
| `PipelinesPage` | `/pipelines` | Your pipelines (run/load) + templates (load) |
| `DashboardPage` | `/dashboards` | Real-time ops dashboard, WebSocket-driven (20s fallback); KPIs, pinned widgets, live chart, machine status, Gantt |
| `KnowledgeGraphPage` | `/kg` | Cytoscape viewer, Technique/Domaine toggle + type filters |
| `OpcuaConnectPage` | `/opcua` | Dynamic OPC-UA connect/browse/subscribe UI |

### `src/lib/` helpers

| File | Purpose |
|---|---|
| `pipelineMapping.js` | Convert canvas ↔ backend `Pipeline` (zones, trigger, depends_on) |
| `functionMeta.js` | Icons/colors/categories per function type |
| `functionDocs.js` | Per-function description + field labels/help/examples |
| `functionDefaults.js` | Default config seeded when a function node is added |
| `connectorTemplates.js` | Default config per connector + trigger type |
| `kgGraph.js` | Map KG JSON → Cytoscape elements + styles |
| `dashboardData.js` | Join domain graph → events (cause/cost), today/yesterday |
| `useLiveSocket.js` | WebSocket hook (auto-reconnect) for live push |

`src/api/client.js` — all `fetch` calls to `/api/*`.

## Adding a new function

1. Create a file under the appropriate sub-package (e.g. `internal/functions/transforms/my_transform.go`)
2. Implement a struct with `GetFunction() *functions.Function` returning a `*functions.Function` with `Name`, `Type`, `Description`, `Inputs`, `Outputs`, and `Handler func(map[string]interface{}) (interface{}, error)`
3. Register it in `buildRegistry()` in `cmd/server/main.go` **and** in `cmd/agent/main.go`

## Known limitations

- Secure OPC-UA modes (`Sign`/`SignAndEncrypt`) need a client certificate not yet wired — use `None`.
- The OPC-UA session holds one subscription at a time; changing tag selections requires a disconnect + reconnect cycle.
- `modbus_read` and `sql_query` are metadata-only stubs that error if executed.
- The technical KG is empty until you save at least one pipeline; shipped example pipelines are templates only.
- Dashboard "vs hier" deltas need events spanning two days; TRS is an availability proxy, not true OEE.
