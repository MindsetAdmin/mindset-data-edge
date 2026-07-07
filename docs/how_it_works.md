# MindSet — How It Works (Backend + Frontend Walkthrough)

> **Technical explainer** — how the currently-built system actually functions end-to-end.
> **Audience**: new developer joining · intern onboarding · advisor doing technical due diligence · Mohamed himself when re-reading later.
> **Grounded in code as of 2026-07-01.** For gaps (what's NOT built vs the V1 target), see `docs/analysis_log.md` Entry 42.
> **NOT covered**: vision (`mindset.md`), strategy (`analysis_log.md`), aspirational V1 features not yet coded.

---

## 1. 30,000-foot view

MindSet runs as **two Go binaries + one React frontend + one MQTT broker + one SQLite file**. All on the same edge machine (customer's PC).

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        CUSTOMER'S EDGE PC                                │
│                                                                           │
│  ┌───────────────┐  MQTT      ┌───────────────┐   HTTP    ┌────────────┐│
│  │  cmd/agent    │──────────▶│  Mosquitto    │◀────────▶│ cmd/server ││
│  │  (edge        │  publish  │  broker       │ subscribe │  (:8080    ││
│  │  runtime)     │           │  :1883        │           │   HTTP+WS) ││
│  └───────────────┘           └───────────────┘           └────────────┘│
│         │                            ▲                          ▲      │
│         │ OPC-UA                     │                          │      │
│         │ (if auto_connect=true)     │                          │      │
│         ▼                            │                          │      │
│  ┌──────────────┐                    │                          │      │
│  │  OPC-UA      │                    │       WebSocket + HTTP   │      │
│  │  server      │                    │       + REST /api/*      │      │
│  │  (Prosys /   │                    │                          │      │
│  │  SCADA)      │                    │                          │      │
│  └──────────────┘                    │                          │      │
│                                      │                          ▼      │
│                              ┌───────┴────────┐          ┌──────────────┐│
│                              │  data/mindset  │          │  React SPA   ││
│                              │  .db (SQLite)  │          │  :5173 dev   ││
│                              │  KG · tags ·   │          │  or built    ││
│                              │  events        │          │  static      ││
│                              └────────────────┘          └──────────────┘│
└─────────────────────────────────────────────────────────────────────────┘
```

**The frontend ONLY ever talks to `cmd/server` via `/api/*`.** It never touches the MQTT broker or `cmd/agent` directly.

---

## 2. The two binaries — what each does

### `cmd/server/` — the HTTP API + WebSocket + OPC-UA control plane

**Purpose**: everything the React UI needs. Owns the OPC-UA session in "UI-driven mode" (the default).

**Files**:
| File | Responsibility |
|---|---|
| `main.go` | Boots API, registers routes, wires everything together |
| `live.go` | `LiveHub` — one MQTT subscription (`mindset/#`) fanning out to Tag/Topic/State registries + WebSocket |
| `opcua.go` | `OPCUAManager` — dynamic connect / browse / subscribe driven by REST calls |
| `opcua_handlers.go` | HTTP handlers for `/api/opcua/*` |
| `tags.go` | `TagRegistry` (persisted to SQLite), `TopicRegistry`, `StateTracker` |
| `ws.go` | `wsHub` — WebSocket broadcaster |

**Listens on**: `:8080` (configurable via `-addr` flag).

### `cmd/agent/` — the edge runtime

**Purpose**: rules engine, KG enrichment, UNS contextualization, pipeline engine. Runs *without* HTTP. Consumes MQTT topics, produces MQTT topics.

**Files**:
| File | Responsibility |
|---|---|
| `main.go` | Boots agent in labeled steps (STEP 0 → 5) |
| `init.go` | Helper for one-time setup |

**Default mode**: `opcua.auto_connect = false` → agent runs the rules engine + KG subscriber + pipeline engine, but does NOT open its own OPC-UA session (server owns it).

**Alternative mode**: `opcua.auto_connect = true` → agent additionally opens the OPC-UA session, runs its own UNS contextualizer, discovery, and subscription. **This mode duplicates `mindset/site/#` publishing if `cmd/server` also has an OPC-UA session — so pick one.**

---

## 3. Startup sequence — `cmd/server`

Real sequence from `cmd/server/main.go`:

```
1. Parse flags: -config, -db, -pipelines, -addr
2. Load config/agent.yaml (fall back to defaults if absent)
3. Open Knowledge Graph SQLite (creates tables if new: kg_nodes, kg_edges, events, tags)
4. Connect to MQTT broker (best-effort — server keeps running if broker is down)
5. Start KG Subscriber if MQTT connected — subscribes to mindset/events/micro-stop
   to auto-enrich the domain KG when events fire
6. Create TagRegistry (SQLite-backed), TopicRegistry, StateTracker, WebSocket hub, LiveHub
7. LiveHub subscribes to mindset/# — fans out into the three registries + WebSocket
8. Create OPCUAManager (idle — waits for /api/opcua/connect REST call)
9. Build the function registry (11 functions across connectors/transforms/calculates/conditions/outputs)
10. Register HTTP routes on http.ServeMux (see route table below)
11. Listen on :8080 with permissive CORS
```

## 4. Startup sequence — `cmd/agent`

Real sequence from `cmd/agent/main.go`, labeled steps:

```
STEP 0    — MQTT Publisher (connect to tcp://localhost:1883)
STEP 0.5  — UNS Contextualizer (ONLY if opcua.auto_connect=true — otherwise the server does this)
STEP 0.6  — Rules Engine (subscribes to mindset/site/#, detects Run↔Stop transitions,
            publishes to mindset/events/status-change)
STEP 0.7  — Knowledge Graph + KG Subscriber (subscribes to mindset/events/micro-stop,
            writes Equipment/Event/Cause/Cost nodes to SQLite)
STEP 0.8  — Function registry: registers state_machine, uns_mapper, filter, calculate_duration,
            calculate_cost, threshold, kg_save (6 functions — a subset of what cmd/server registers)
STEP 0.9  — Pipeline engine: loads all YAML from config/pipelines/, registers them
STEP 0.10 — Technical KG: builds pipeline topology graph in memory (5-min cache)

Then:
IF opcua.auto_connect=false (default) → sit idle, keep rules/KG running, wait for Ctrl+C
IF opcua.auto_connect=true            → open OPC-UA, browse tags, subscribe, print value changes
```

Both binaries recover cleanly from missing dependencies (no MQTT broker → warnings, keep running).

---

## 5. Data flow — from OPC-UA tag change to dashboard render

Full trace of what happens when a machine's `Etat_Machine` transitions from Run to Stop:

```
[Real machine → OPC-UA server → cmd/server → MQTT → cmd/agent → MQTT → cmd/server → UI]

(1) Machine state changes in the PLC. OPC-UA server updates the node value.

(2) OPCUAManager (in cmd/server, running a MonitoredItem with 500ms sampling)
    receives the value change callback.

(3) OPCUAManager publishes to MQTT:
    Topic: mindset/raw/{nodeID}
    Payload: {"node_id":"...", "name":"Etat_Machine", "value":false,
              "data_type":"Boolean", "timestamp_ms":...}

(4a) LiveHub (in cmd/server) is subscribed to mindset/#. It:
     - Updates TagRegistry (persisted to SQLite)
     - Updates TopicRegistry (msg/s counter)
     - Updates StateTracker
     - Broadcasts via WebSocket to any connected frontend clients

(4b) UNS contextualizer (in cmd/server via the OPCUAManager's ISA-95 mapping mode,
     OR in cmd/agent if agent has an OPC-UA session) reads mindset/raw/{nodeID},
     applies the customer's ISA-95 mapping (site → area → work center → tag),
     publishes to:
     Topic: mindset/site/{site}/{area}/{work_center}/{tag_name}
     Payload: {"timestamp_ms":..., "value":false, "unit":"",
              "data_type":"Boolean", "metadata":{... ISA-95 fields ...}}

(5) Rules engine (in cmd/agent) is subscribed to mindset/site/#. It:
    - Updates internal state store
    - For tags with tag_name="status": checks for Run↔Stop transition
    - On Stop→Run: computes duration and publishes:
      Topic: mindset/events/status-change
      Payload: {"topic":"...", "work_center":"...", "duration_seconds":47,
                "previous_state":true, "current_state":false, ...}

(6) A pipeline like `microstop_detection` (defined in config/pipelines/*.yaml)
    is triggered by mindset/events/status-change (via a mqtt_subscribe node).
    Its nodes execute in topological order:
    - state_machine: classifies as Run/Stop pattern
    - calculate_duration: extracts duration
    - threshold: checks 30s < duration < 3min
    - calculate_cost: converts to €
    - mqtt_publish: publishes to mindset/events/micro-stop
    - kg_save: writes an Event node to the domain KG

(7) KG subscriber (in either binary) is subscribed to mindset/events/micro-stop.
    It writes Equipment, Event, Cause, Cost nodes and edges to the SQLite KG.

(8) LiveHub (cmd/server) sees ALL these MQTT publishes (subscribed to mindset/#).
    Any "dashboard-relevant" messages get broadcast via WebSocket.

(9) Frontend (React) receives WebSocket messages via useLiveSocket hook.
    Dashboard widgets, LiveDataPanel, machine status page all re-render in real time.
```

---

## 6. MQTT topic taxonomy (the "nervous system")

| Topic pattern | Payload shape | Producer | Consumer |
|---|---|---|---|
| `mindset/raw/{nodeID}` | `{node_id, name, value, data_type, timestamp_ms}` | OPCUAManager or discovery/opcua.go | UNS contextualizer, LiveHub |
| `mindset/site/{site}/{area}/{work_center}/{tag}` | ISA-95 enriched payload with metadata | UNS contextualizer | Rules engine, LiveHub, pipelines |
| `mindset/events/status-change` | `{topic, work_center, duration_seconds, previous_state, current_state, ...}` | Rules engine | Micro-stop detection pipeline |
| `mindset/events/micro-stop` | `{work_center, duration_seconds, cause, cost_eur, ...}` | Micro-stop pipeline (via `mqtt_publish` output) | KG subscriber, dashboards |
| `mindset/dashboard/{label}` | Any pinned value | `add_to_dashboard` output function | LiveHub → WebSocket → DashboardWidgets frontend |

**Key insight**: MQTT is the ONLY inter-process communication. `cmd/server` and `cmd/agent` never call each other directly — they exchange messages through the broker.

---

## 7. Pipeline engine — how execution works

Pipelines live in `config/pipelines/*.yaml`. Example:

```yaml
id: microstop_detection
name: Micro-Stop Detection
enabled: true
trigger:
  type: mqtt
  config:
    topic: mindset/events/status-change
nodes:
  - id: state_machine_step
    type: transform
    function: state_machine
    config: { machine_id: "line1" }
  - id: duration_step
    type: calculate
    function: calculate_duration
    depends_on: [state_machine_step]
  - id: threshold_step
    type: condition
    function: threshold
    config: { min: 30, max: 180 }
    depends_on: [duration_step]
  - id: cost_step
    type: calculate
    function: calculate_cost
    depends_on: [threshold_step]
  - id: publish_step
    type: output
    function: mqtt_publish
    config: { topic: "mindset/events/micro-stop" }
    depends_on: [cost_step]
```

### Loader
`internal/pipeline/loader.go` reads all `.yaml` under the configured directory (default `config/pipelines/`), parses via `gopkg.in/yaml.v3` into `Pipeline` structs, hands to the Engine.

### Engine execution
From `internal/pipeline/engine.go`:

1. Look up pipeline by ID
2. Build a `nodeMap` (id → Node)
3. Loop until all nodes executed:
   - For each unexecuted node: check if all `depends_on` are satisfied
   - A dep that isn't a node ID (e.g., `"trigger"`) is treated as **already satisfied** — that's the trigger placeholder
   - When deps are met: call `executeNode()`
4. `executeNode()`:
   - Merges `previous node outputs` + `trigger data` + node's own `config` map into one `params` map
   - Looks up the function in the registry
   - Calls the function's `Handler(params) (interface{}, error)`
   - A `recover()` wraps the call so a panicking handler doesn't crash the server
   - Result stored in `results[node.ID]`, emitted as a `PipelineEvent`

### Trigger types (declarative)
- `mqtt` — pipeline is triggered by a subscribed MQTT topic
- `http` — triggered by `POST /api/pipelines/{id}/run`
- `timer` — declared but wiring is minimal in the current code

### Failure mode
If any node returns an error, the whole pipeline is marked failed. Downstream nodes don't run. The `recover()` protection means bugs in handlers surface as errors, not crashes.

---

## 8. Function registry — the plug-in system

`internal/functions/registry.go` holds a map of `Name` → `*Function`.

Currently registered (11 in `cmd/server`, 6 in `cmd/agent`):

| Type | Function name | What it does |
|---|---|---|
| Connector | `opcua_read` | Read a value from an OPC-UA node |
| Connector | `mqtt_subscribe` | Subscribe to an MQTT topic (pipeline trigger) |
| Connector | `modbus_read` | Demo stub — errors if executed |
| Connector | `sql_query` | Demo stub — errors if executed |
| Transform | `state_machine` | Detect Run↔Stop transitions |
| Transform | `uns_mapper` | Normalize a tag to ISA-95 structure |
| Transform | `filter` | Keep/drop by condition |
| Calculate | `calculate_duration` | Duration between start/stop timestamps |
| Calculate | `calculate_cost` | € cost from duration × hourly rate (V0 — see `impact_engine.md` for what this SHOULD become) |
| Condition | `threshold` | Is value within `[min, max]` ? |
| Output | `mqtt_publish` | Publish to an MQTT topic |
| Output | `add_to_dashboard` | Pin data onto the dashboard |
| Output | `kg_save` | Save to KG (only registered when a KG instance is available) |

**Adding a new function**: create a file under `internal/functions/{category}/`, implement `GetFunction() *functions.Function`, register it in BOTH `cmd/server/main.go` and `cmd/agent/main.go`.

---

## 9. Knowledge Graph — one unified graph, two categories

**Refactored 2026-07-02** (see `docs/analysis_log.md` Entry 50). Previously the code had a Domain KG (persistent SQLite) and a Technical KG (in-memory, 5-min cached). They were merged into a single unified KG.

### Unified KG (persistent, SQLite)

`internal/kg/graph.go` + `internal/kg/subscriber.go` + `internal/kg/builder.go`.

**Storage**: SQLite tables `kg_nodes` (id, **category**, type, label, properties JSON) and `kg_edges` (id, **category**, from_id, to_id, relation, weight). Both tables gained a `category` column.

**Every node and edge is tagged with a category:**

| Category | What lives here | Populated by |
|---|---|---|
| `business` | Equipment · Event · Cause · Cost · Operator · OF · Product · Recipe … | `KGSubscriber` listening to `mindset/events/micro-stop` (auto) + `kg_save` output function (manual) |
| `platform` | Connection · Topic · Function · Pipeline · Dashboard | `KnowledgeGraph.RepopulatePlatform(registry)` — wipes+rebuilds when pipelines change |

**Cross-category edges are legal.** Example: a `Dashboard` platform node has an edge (`subscribes_to`) to a `Topic` platform node, which is `produced_by` a `Pipeline`, which processes `Event` business nodes from `Equipment`. AI agents can traverse across categories via MCP.

### Access

Single unified endpoint plus legacy aliases (both still work):

| Route | What it returns |
|---|---|
| `GET /api/kg?category=business` | Only site fingerprint (Equipment/Event/Cause/Cost/…) |
| `GET /api/kg?category=platform` | Only pipeline topology (Connection/Topic/Function/Pipeline/Dashboard) |
| `GET /api/kg?category=all` (or omit) | Both, in one graph |
| `GET /api/kg/domain` | Legacy alias → `category=business` |
| `GET /api/kg/technical` | Legacy alias → `category=platform` |

Returns JSON for Cytoscape rendering. Both nodes and edges carry a `category` field so the frontend can style them differently.

### Why unified

Aligns with:
- **Prop #7** — KG/UNS as the single trusted source for AI agents (one endpoint, one schema)
- **Prop #9** — pipeline updates automatically appear in "THE KG" (no more distinguishing "which KG?")
- **Impact Engine + Moat #3** — the cumulative site fingerprint has one home
- **MCP tools** — `kg_query`, `kg_list_events`, `kg_cost_summary` all target one graph

**Platform sub-graph rebuild** happens lazily on `GET /api/kg?category=platform|all` OR eagerly when pipelines are registered/deregistered. Replaces the old 5-min cache. No more "Technical KG is stale."

---

## 10. UNS contextualization

`internal/uns/mapper.go` + `internal/uns/contextualizer.go`.

**Purpose**: transform raw OPC-UA tags (e.g. `ns=2;s=Line1.Motor.Speed`) into a normalized ISA-95 topic structure.

**Contextualizer** (when active):
1. Subscribes to `mindset/raw/#`
2. For each message, calls the Mapper to normalize:
   - `Site` (from config, e.g. "usine-nord")
   - `Area` (parsed from tag path, e.g. "line1")
   - `WorkCenter` (parsed further)
   - `WorkUnit`, `TagName`
3. Republishes to `mindset/site/{site}/{area}/{work_center}/{tag_name}` with an enriched payload including the full ISA-95 metadata

**Who runs it**: `cmd/server` in UI-driven OPC-UA mode; `cmd/agent` if `opcua.auto_connect=true`.
**Never both** (would double-publish).

---

## 11. LiveHub + WebSocket

`cmd/server/live.go` + `cmd/server/ws.go`.

The `LiveHub` is the single MQTT subscriber on the server side. It subscribes ONCE to `mindset/#` and routes messages into:

- **TagRegistry** — persists tag latest values to SQLite
- **TopicRegistry** — tracks messages/second per topic
- **StateTracker** — Run/Stop state per work center
- **WebSocket hub** — broadcasts to frontend

The WebSocket sends message envelopes:

```json
{"type": "tag",       "data": {...}}   // OPC-UA tag update
{"type": "state",     "data": {...}}   // Machine state change
{"type": "event",     "data": {...}}   // Detected event (micro-stop, etc.)
{"type": "dashboard", "data": {...}}   // Pinned dashboard value
```

The frontend hook `useLiveSocket.js` handles auto-reconnect + a 20s fallback poll if the socket is lost.

---

## 12. Storage — SQLite via `modernc.org/sqlite`

Pure-Go SQLite (no CGO required). One file: `data/mindset.db`.

**Tables (auto-created)**:

| Table | Purpose |
|---|---|
| `kg_nodes` | Domain KG nodes (id, type, label, properties JSON) |
| `kg_edges` | Domain KG edges (id, from_id, to_id, relation, weight) |
| `events` | Reserved for event log |
| `tags` | Tag values persisted from LiveHub (node_id, name, value JSON, data_type, timestamp_ms) |

**Ring buffer / TTL cleanup**: not yet implemented (per Entry 42 gaps).

---

## 13. HTTP API surface (real routes registered in `cmd/server/main.go`)

CORS: `*` (permissive).

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/health` | `{status:"ok"}` liveness check |
| GET | `/api/functions[?type=]` | Function catalog (optionally filtered) |
| GET | `/api/connectors` | Connector-type functions |
| GET/POST | `/api/pipelines` | List / save YAML |
| GET | `/api/pipelines/examples` | Template pipelines |
| POST | `/api/pipelines/{id}/run` | Execute a pipeline synchronously |
| GET | `/api/tags` | Live OPC-UA tags + values |
| GET | `/api/machines` | Tags grouped by work center + Running/Stopped state |
| GET | `/api/topics` | Live MQTT topics + msg/s |
| GET | `/api/config` | Safe subset of config/agent.yaml |
| POST | `/api/opcua/connect` | Dynamically open OPC-UA session |
| GET | `/api/opcua/discover` | Browse the connected server's tree |
| POST | `/api/opcua/subscribe` | Monitor selected tags with per-tag mode (raw / isa95 / both) |
| POST | `/api/opcua/disconnect` | Close OPC-UA session |
| GET | `/api/opcua/status` | Connection status |
| GET | `/api/opcua/selections` | Per-tag routing + ISA-95 mapping |
| GET | `/api/dashboard/pins` | Snapshot of `add_to_dashboard` pins |
| GET | `/api/kg/technical` | Technical (architecture) graph JSON |
| GET | `/api/kg/domain` | Domain (data) graph JSON |
| GET | `/api/stats` | Counts + micro-stops summary + uptime |
| WS | `/api/ws` | WebSocket for live push |

---

## 14. Frontend architecture

**Stack**: React 19 · Vite · Tailwind CSS · Zustand (state) · React Flow (pipeline canvas) · Cytoscape (KG viewer) · Recharts (dashboard charts) · react-router-dom (routing).

**Root**: `frontend/pipeline-builder/src/App.jsx`

**Routes** (from App.jsx):

```
/            → redirect to /overview
/overview    → OverviewPage
/connect     → ConnectPage
/connect/opcua → OpcuaConnectPage
/compose     → BuilderPage (pipeline builder canvas)
/pipelines   → PipelinesPage
/dashboards  → DashboardPage
/kg          → KnowledgeGraphPage
```

**Layout**: `NavBar` on top + Routes below, all wrapped in `ErrorBoundary` that resets on route change.

**API base**: Vite proxies `/api/*` to `http://localhost:8080` (see `vite.config`). In production, the built static assets are served from the same origin as the API.

---

## 15. Frontend pages walkthrough

### OverviewPage (`/overview`)
Home: key stats (from `/api/stats`) + quick links.

### ConnectPage (`/connect`)
Pick a connector type. Selecting one applies it as the pipeline trigger (redirects to the Builder).

### OpcuaConnectPage (`/connect/opcua`)
Dynamic OPC-UA connect UI (built recently). Steps:
1. Enter endpoint + security settings → POST `/api/opcua/connect`
2. Browse tree → GET `/api/opcua/discover`
3. Select tags + per-tag mode (raw / isa95 / both) → POST `/api/opcua/subscribe`
4. Live values start flowing via WebSocket

### BuilderPage / `/compose` — the Pipeline Studio
The heart of the UI. React Flow canvas with:
- **Palette** (left): draggable functions organized by category (Connector · Transform · Calculate · Condition · Output)
- **Canvas** (center): drop functions here, connect them, delete them, run them
- **Config panel** (right): opens when a node is clicked — form fields per function type (from `lib/functionDocs.js` + `lib/functionDefaults.js`)
- Actions: **Save** (POST `/api/pipelines`, converts canvas → YAML), **Run** (POST `/api/pipelines/{id}/run`), Delete

### PipelinesPage (`/pipelines`)
Your saved pipelines list + example templates from `config/pipelines/examples/`. Can load a template into the Builder or run existing.

### DashboardPage (`/dashboards`)
Real-time ops dashboard. Shows:
- KPI cards (from `/api/stats` + WebSocket)
- **DashboardWidgets** — user pins values from pipeline outputs (`add_to_dashboard` function), configures chart type (line / bar / gauge / value / status) and time range
- Live machine state per work center
- Gantt-style timeline (basic)

### KnowledgeGraphPage (`/kg`)
Cytoscape viewer with a **Technique/Domaine** toggle:
- Technique: pipeline architecture (from `/api/kg/technical`)
- Domaine: cumulative site data (from `/api/kg/domain`)

Filters by node type (Connection · Function · Topic · Pipeline · Dashboard).

---

## 16. Key frontend components (what each does)

| Component | Purpose |
|---|---|
| `NavBar` | Top navigation |
| `ErrorBoundary` | Catches render errors, resets on route change |
| `Palette` | Draggable function catalog in the Builder |
| `nodes/PipelineNode` | ReactFlow node renderer for pipeline functions |
| `nodes/TriggerNode` | Special ReactFlow node for the pipeline trigger |
| `nodes/ZoneNode` | Visual zone container in the Builder |
| `NodeConfigPanel` | Right sidebar — dynamic form per function type |
| `DashboardWidgets` | Pin-based widget system (Recharts) |
| `LiveDataPanel` | Real-time tag + topic viewer |
| `OpcuaConnectionPanel` | OPC-UA connect form |
| `OpcuaTagSelector` | Tag tree browse + selection UI |
| `CytoscapeGraph` | KG visualization wrapper |
| `PickerModal` | Reusable modal for choosing items |

## 17. Frontend helper libraries (`src/lib/`)

| File | Purpose |
|---|---|
| `pipelineMapping.js` | Convert canvas graph ↔ backend `Pipeline` YAML |
| `functionMeta.js` | Icons + colors + categories per function type |
| `functionDocs.js` | Description + field labels/help per function |
| `functionDefaults.js` | Default config seeded when a node is added |
| `connectorTemplates.js` | Default config per connector + trigger type |
| `kgGraph.js` | Map KG JSON → Cytoscape elements + styles |
| `dashboardData.js` | Join domain graph → events (cause / cost / today / yesterday) |
| `pipelineLoading.js` | Load/save pipelines to the backend |
| `useLiveSocket.js` | WebSocket hook with auto-reconnect + 20s fallback poll |

---

## 18. State management — Zustand

`src/store/studioStore.js` holds:

- Current pipeline being built (nodes + edges + trigger)
- Selected node (for config panel)
- Function catalog cache
- OPC-UA session state
- User's tag selections
- ...

Components subscribe via `useStore(state => state.someField)`. Changes automatically re-render dependents.

## 19. `useLiveSocket` hook — real-time updates

Every page that needs live data calls:
```javascript
const connected = useLiveSocket((msg) => {
  // msg.type is one of: 'tag' | 'state' | 'event' | 'dashboard'
  // msg.data is the payload
});
```

The hook opens `ws://localhost:8080/api/ws` (proxied by Vite in dev), reconnects if the socket closes, and polls REST APIs every 20s as a fallback if disconnected.

---

## 20. Configuration

### `config/agent.yaml` (optional — defaults kick in if absent)

```yaml
site:
  name: "Local Test Site"
  id: "local-test"
opcua:
  name: "Prosys"
  endpoint: "opc.tcp://localhost:53530/OPCUA/SimulationServer"
  security_mode: "None"
  security_policy: "None"
  auto_connect: false     # UI drives OPC-UA in cmd/server
mqtt:
  broker: "tcp://localhost:1883"
cost:
  hourly_cost: 85.0
```

Loaded by `internal/config/config.go` — both binaries.

### `config/pipelines/*.yaml`

Every pipeline is one YAML file. Loaded at agent startup + reloadable via the Builder Save action.

---

## 21. Running locally (Windows)

### Automated
```powershell
.\run.ps1              # build both binaries, start server + agent + frontend, open browser
.\run.ps1 -NoAgent     # skip the edge agent
.\run.ps1 -NoBuild     # reuse existing bin/*.exe
```

### Manual
```powershell
# Build
go build -o bin/server.exe ./cmd/server
go build -o bin/agent.exe  ./cmd/agent

# Run backend
.\bin\server.exe                          # :8080 API + WebSocket
.\bin\agent.exe                           # edge runtime

# Run frontend (separate terminal)
cd frontend/pipeline-builder
npm install                               # first time only
npm run dev                               # Vite dev on :5173
```

Prerequisite: MQTT broker on `localhost:1883` (Mosquitto recommended). The system will run without it but MQTT-dependent features are limited.

For OPC-UA testing without a real factory: install Prosys OPC-UA Simulator (free) and point `opcua.endpoint` at it.

---

## 22. What's NOT yet built (see Entry 42 for detail)

At a high level:
- **Backend gaps**: Modbus connector · SQL connector (no drivers in go.mod) · MCP server · Phi-3/Ollama runtime · Ad-hoc Analyst · **OF-state Fuzzy Join engine (the moat)** · real Impact Engine (current cost.go is V0 stub) · OEE/TRS calculator · cloud push + heartbeat · alerting (SMTP/Slack/Teams) · license validator · SOPS secrets · network scanner · behavioral inference
- **Frontend gaps**: dedicated Gantt timeline · Pareto by € · OEE/TRS view (the killer demo) · ROI simulator · Tribal knowledge capture UI · Onboarding wizard · Ad-hoc Analyst chat panel
- **Infrastructure gaps**: NO `.github/workflows/` · NO Dockerfile · NO docker-compose · NO `_test.go` files (zero tests) · NO signed binaries · NO SBOM · NO Mosquitto bundle config

Full gap analysis: `docs/analysis_log.md` Entry 42.

---

## 23. Where to look in the code — quick index

| Want to understand… | Read this |
|---|---|
| Server boot + HTTP routes | `cmd/server/main.go` |
| Agent boot + startup steps | `cmd/agent/main.go` |
| OPC-UA reading | `internal/discovery/opcua.go` + `cmd/server/opcua.go` |
| Rules engine + state transitions | `internal/rules/engine.go` + `internal/rules/state.go` |
| Pipeline execution algorithm | `internal/pipeline/engine.go` |
| Function registry pattern | `internal/functions/registry.go` + any file under `internal/functions/*/` |
| Cost calculation (V0 stub) | `internal/functions/calculates/cost.go` |
| Domain KG | `internal/kg/graph.go` + `internal/kg/subscriber.go` |
| Technical KG | `internal/kg/builder.go` |
| UNS ISA-95 mapping | `internal/uns/mapper.go` + `internal/uns/contextualizer.go` |
| Storage | `internal/storage/sqlite.go` |
| Live data fan-out (server side) | `cmd/server/live.go` + `cmd/server/ws.go` |
| Frontend routing | `frontend/pipeline-builder/src/App.jsx` |
| Frontend state | `frontend/pipeline-builder/src/store/studioStore.js` |
| Real-time WebSocket hook | `frontend/pipeline-builder/src/lib/useLiveSocket.js` |
| Pipeline canvas → backend YAML | `frontend/pipeline-builder/src/lib/pipelineMapping.js` |
| Dashboard widgets | `frontend/pipeline-builder/src/components/DashboardWidgets.jsx` |

---

## 24. Summary

- **2 Go binaries** communicate via **1 MQTT broker**; both write to **1 SQLite file**.
- **`cmd/server`** owns HTTP + WebSocket + OPC-UA (in UI-driven mode). The frontend only talks to it.
- **`cmd/agent`** runs the rules engine, KG enrichment, pipeline engine, UNS contextualization.
- **Pipelines** are YAML files, executed by topological sort with `recover()`-protected handlers.
- **11 functions** in the registry cover connectors, transforms, calculates, conditions, outputs.
- **Two KGs**: a persistent Domain KG (site fingerprint) and an in-memory Technical KG (pipeline topology).
- **React SPA** with 7 routes, real-time via WebSocket + 20s fallback poll.
- **Everything is grounded in code** as of 2026-07-01. Gaps vs the V1 target catalog: see Entry 42.
