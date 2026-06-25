# MindSet Data — Architecture & How It Works

> Industrial IoT edge platform: it reads machine data from an **OPC‑UA** server,
> streams it over **MQTT**, contextualizes it into an **ISA‑95 Unified Namespace**,
> detects **micro‑stops**, computes **costs**, builds a **Knowledge Graph**, and
> exposes everything through a **web studio** (pipeline builder + live dashboard).

---

## 1. The big picture

The solution is made of **three runtimes** plus an MQTT broker and an OPC‑UA server:

```
 ┌─────────────┐   OPC-UA    ┌──────────────────────┐   MQTT    ┌───────────────────────┐
 │  Prosys     │◄───────────►│  Edge Agent          │──────────►│      MQTT Broker      │
 │  OPC-UA Sim │  :53530     │  (cmd/agent)         │  publish  │      :1883            │
 └─────────────┘             │                      │           │                       │
                             │ discover → publish   │           │  mindset/raw/#        │
                             │ contextualize (UNS)  │◄──────────│  mindset/site/#       │
                             │ rules (micro-stops)  │ subscribe │  mindset/events/#     │
                             └──────────────────────┘           └───────────┬───────────┘
                                                                            │ subscribe
                                       reads config + writes KG             │
                             ┌──────────────────────┐                       │
                             │  API Server          │◄──────────────────────┘
                             │  (cmd/server) :8080  │
                             │  REST /api/*         │   reads/writes  ┌──────────────────┐
                             │  LiveHub (mindset/#) │◄───────────────►│ data/mindset.db  │
                             └──────────┬───────────┘                 │ (SQLite: KG+tags)│
                                        │ HTTP /api (Vite proxy)      └──────────────────┘
                             ┌──────────▼───────────┐
                             │  Web Studio (React)  │
                             │  frontend :5173      │
                             │  Overview/Connect/   │
                             │  Compose/Pipelines/  │
                             │  Dashboards/KG       │
                             └──────────────────────┘
```

| Runtime | Path | Role | Needs |
|---|---|---|---|
| **Edge Agent** | `cmd/agent` | The on‑site runtime. Connects to OPC‑UA, publishes raw tags, contextualizes to ISA‑95, detects micro‑stops, enriches the KG. | OPC‑UA server + MQTT broker |
| **API Server** | `cmd/server` | Backend the UI talks to. Serves functions, pipelines, KG, live tags/topics/machines, config, pipeline execution. Auto‑enriches the KG and tracks live state. | Nothing to start; MQTT only for live data/Run |
| **Web Studio** | `frontend/pipeline-builder` | React UI: build pipelines, pick connectors, view the KG, watch the live dashboard. | The API server |

**Key coupling rule:** the frontend talks **only** to the API server (`/api`). The
agent and server are **loosely coupled** — they never call each other directly;
they share the **MQTT broker** and the **SQLite database**.

---

## 2. End‑to‑end data flow

```
OPC-UA tag value changes
        │  discovery/opcua.go  (Subscribe → monitored items)
        ▼
mindset/raw/<nodeID>              ← raw tag JSON {node_id,name,value,data_type,timestamp_ms}
        │
        ├──► uns/contextualizer.go  subscribes raw, maps to ISA-95
        │         ▼
        │   mindset/site/<site>/<area>/<machine>/<tag>   ← contextualized value
        │
        └──► rules/engine.go  watches status tags, detects Run↔Stop
                  ▼
            mindset/events/status-change   (transition + duration)
                  │  (a micro-stop pipeline / rules)
                  ▼
            mindset/events/micro-stop      {work_center,duration_seconds,cause,cost_eur,...}
                  │
                  ├──► kg/subscriber.go   → writes Equipment/Event/Cause/Cost nodes to SQLite
                  └──► (cost pipeline)    → mindset/events/micro-stop-cost
```

On the **server side**, `cmd/server/live.go`’s **LiveHub** subscribes to
`mindset/#` and feeds three in‑memory views used by the UI:

- **TagRegistry** (`tags.go`) — latest value per OPC‑UA tag (persisted to SQLite `tags`).
- **TopicRegistry** — every topic seen + a 5‑second **msg/s rate** + category (raw/site/events).
- **StateTracker** — machine **Running/Stopped** + transition history (from `*.status` tags).

And `kg.NewKGSubscriber` (started in `cmd/server`) listens to
`mindset/events/micro-stop` and **auto‑enriches the Knowledge Graph** — the user
never adds a "save to KG" step.

---

## 3. MQTT topic taxonomy

| Topic | Payload | Produced by | Consumed by |
|---|---|---|---|
| `mindset/raw/<nodeID>` | `{node_id,name,value,data_type,timestamp_ms}` | agent `mqtt/publisher.go` | UNS contextualizer, server LiveHub |
| `mindset/site/<site>/<area>/<machine>/<tag>` | contextualized value (ISA‑95) | agent `uns/contextualizer.go` | server LiveHub (topic list) |
| `mindset/events/status-change` | transition + duration | agent `rules/engine.go` | micro‑stop logic |
| `mindset/events/micro-stop` | `{work_center,duration_seconds,cause,confidence,cost_eur,...}` | rules / pipelines | server KG subscriber |
| `mindset/events/micro-stop-cost` | cost event | cost pipeline | dashboards / KG |

---

## 4. Backend (Go)

Module: `github.com/MindsetAdmin/mindset-data-edge` (Go 1.26). SQLite driver is
pure‑Go (`modernc.org/sqlite`) so binaries are self‑contained (no CGO).

### 4.1 `cmd/agent` — the edge runtime
- `main.go` — boots everything: config → MQTT publisher → UNS contextualizer →
  rules engine → KG subscriber → functions/pipelines → OPC‑UA discovery →
  live subscription. **No HTTP** (that lives in `cmd/server`).
- `init.go` — helper constructors for each subsystem.

### 4.2 `cmd/server` — the API
- `main.go` — wires dependencies and routes; owns config, MQTT client, registries.
- `tags.go` — **TagRegistry**: live OPC‑UA tags + SQLite `tags` persistence.
- `live.go` — **LiveHub** (`mindset/#`), **TopicRegistry** (rates), **StateTracker**.

### 4.3 `internal/` packages
| Package | Responsibility |
|---|---|
| `config` | Loads `config/agent.yaml` (site, opcua, mqtt, cost). |
| `discovery` | OPC‑UA connect, `BrowseNodeTree`, subscribe to value changes, publish raw. |
| `mqtt` | MQTT publisher (`PublishRaw`, `PublishEvent`). |
| `uns` | `mapper.go` (tag → ISA‑95 normalization) + `contextualizer.go` (raw→site). |
| `rules` | `engine.go` (Run/Stop detection, micro‑stop events) + `state.go` (StateStore). |
| `functions` | Function **registry** + types; the catalog of pipeline building blocks. |
| `functions/{connectors,transforms,calculates,conditions,outputs}` | Concrete functions. |
| `pipeline` | `types.go`, `loader.go` (YAML), `registry.go`, `engine.go` (topological exec), `builder.go`. |
| `kg` | Knowledge Graph: `graph.go` (domain CRUD), `builder.go` (technical graph), `subscriber.go` (auto‑enrich), `types.go`. |
| `storage` | SQLite store + schema. |

### 4.4 The function catalog
| Function | Type | Purpose |
|---|---|---|
| `opcua_read` | connector | Read a value from an OPC‑UA node |
| `mqtt_subscribe` | connector | Subscribe to an MQTT topic (pipeline trigger) |
| `modbus_read`, `sql_query` | connector | **demo stubs** (registered for the picker) |
| `state_machine` | transform | Detect Run↔Stop transitions |
| `uns_mapper` | transform | Normalize a tag to ISA‑95 |
| `filter` | transform | Keep/drop by condition |
| `calculate_duration` | calculate | Duration between start/stop |
| `calculate_cost` | calculate | € cost from duration × hourly rate |
| `threshold` | condition | Is value within [min,max]? (micro‑stop window) |
| `mqtt_publish` | output | Publish to an MQTT topic |
| `add_to_dashboard` | output | Pin the data/event onto the Dashboard (`mindset/dashboard/<label>`) |
| `kg_save` | output | Save to KG — **not** in the server palette (KG is automatic) |

Outputs are **sinks**: in the builder they have an input port only (no output
port). `calculate_cost` supports a rate **source** (Manual / from `agent.yaml` /
from a tag) and a per‑product **rate table** uploaded from CSV/Excel
(`config.rates`, keyed by the event's `product`).

### 4.5 Pipelines
A pipeline is YAML in `config/pipelines/` with a **trigger** + **nodes** (each has
`function`, `config`, `depends_on`). The engine (`pipeline/engine.go`) runs nodes
in dependency order (a dependency that isn’t a node — e.g. `"trigger"` — counts as
satisfied). Built‑ins:
- `microstop_detection` — `mqtt_subscribe → state_machine → calculate_duration → threshold → mqtt_publish`
- `opcua_to_uns` — `mqtt_subscribe → uns_mapper → mqtt_publish`
- `cost_calculation` — `mqtt_subscribe → calculate_cost → mqtt_publish`

### 4.6 Knowledge Graph: two graphs
- **Technical graph** (`kg/builder.go`, `/api/kg/technical`) — the *architecture*:
  Connections, Topics, Pipelines, Dashboards. Derived **only from the pipelines you
  build** (empty until you create one; the shipped examples don't appear).
  **External‑only model**: a Pipeline node links to connections/topics/dashboards
  and carries a `functions` count — it does **not** expose its internal functions.
  Pipelines with the **same function signature** are **grouped** into one node that
  lists all their tags (property `tags`).
- **Domain graph** (`kg/graph.go`, `/api/kg/domain`) — the *data*: `Equipment`,
  `Event` (micro‑stops), `Cause`, `Cost` nodes with `occurred_at` / `caused_by` /
  `costs` edges.

### 4.7 SQLite schema (`data/mindset.db`)
| Table | Columns | Written by |
|---|---|---|
| `kg_nodes` | id, type, label, properties(JSON), created_at | KG subscriber / graph |
| `kg_edges` | id, from_id, to_id, relation, weight, created_at | KG subscriber / graph |
| `tags` | node_id, name, value(JSON), data_type, timestamp_ms, updated_at | server TagRegistry |
| `events` | id, type, work_center, duration_seconds, cause, cost_eur, payload, timestamp | (reserved) |

---

## 5. HTTP API (server `:8080`)

| Method & path | Returns |
|---|---|
| `GET /api/health` | `{status:"ok"}` |
| `GET /api/config` | safe subset of `agent.yaml` (opcua endpoint/security, broker, hourly rate, site/area) |
| `GET /api/functions[?type=]` | function catalog (optionally filtered by type) |
| `GET /api/connectors` | connector functions only |
| `GET /api/pipelines` | your pipelines loaded from `config/pipelines/` |
| `POST /api/pipelines` | save a pipeline as YAML (validated, id sanitized) |
| `GET /api/pipelines/examples` | shipped template pipelines (`config/pipelines/examples/`) |
| `POST /api/pipelines/{id}/run` | execute a pipeline; returns per‑node status + timing |
| `GET /api/tags` | live OPC‑UA tags + values (persisted) |
| `GET /api/machines` | tags grouped by work center + live Running/Stopped state |
| `GET /api/topics` | live topics + msg/s + category + broker_connected |
| `GET /api/config` | safe subset of `agent.yaml` |
| `GET /api/dashboard/pins` | current dashboard pins (from `add_to_dashboard`) |
| `GET /api/kg/technical` | technical (architecture) graph |
| `GET /api/kg/domain` | domain (data) graph |
| `GET /api/stats` | counts + micro‑stops/downtime/cost + uptime + broker status |
| `WS  /api/ws` | **WebSocket** live push: `{type:"tag"\|"state"\|"event"\|"dashboard"}` |

CORS is enabled; in dev the Vite proxy forwards `/api` (HTTP **and** the
WebSocket upgrade, `ws:true`) → `:8080`.

---

## 6. Frontend (React + Vite)

Stack: **React 19**, **Vite**, **React Router**, **Zustand** (cross‑page state),
**ReactFlow** (canvas), **Cytoscape** (KG), **Recharts** (charts), **xlsx**
(CSV/Excel parsing), **WebSocket** (live push), **Tailwind**.

### 6.1 Pages (`src/pages/`)
| Page | Route | What it does |
|---|---|---|
| `OverviewPage` | `/overview` | Landing: key stats + quick links |
| `ConnectPage` | `/connect` | Pick a connector → applies to the pipeline trigger |
| `BuilderPage` | `/compose` | **Drag‑and‑drop builder** (ReactFlow): ENTRÉE/CŒUR/SORTIE bands, guided config panel, function/field pickers, **smart validation + duplicate prevention**, Save (→YAML), Run, delete node |
| `PipelinesPage` | `/pipelines` | Two sections: **your** pipelines (run/load) + **templates** (load) |
| `DashboardPage` | `/dashboards` | **Real‑time ops dashboard**, **WebSocket‑driven** (20s heartbeat fallback): KPIs, pinned widgets, live tag chart, recent events, machine status, Gantt |
| `KnowledgeGraphPage` | `/kg` | Cytoscape viewer, Technique/Domaine toggle + type filters |

### 6.2 Components (`src/components/`)
`NavBar` (shows the **MindSet Data logo**, `public/logo.png`), `Palette`,
`NodeConfigPanel` (guided config: header/badge/description, labelled fields +
help/examples, pickers, OPC‑UA machine+tag selector, cost source + CSV/Excel rate
upload + live preview, delete), `PickerModal`, `CytoscapeGraph`, `ErrorBoundary`,
**`LiveDataPanel`** (pick tags → live chart), **`DashboardWidgets`** (interactive
`add_to_dashboard` widgets: line/bar/gauge/value/status charts, time ranges, live
stats, `✕`/`⚙️`, persisted in localStorage), and ReactFlow nodes
`nodes/{PipelineNode,TriggerNode,ZoneNode}` (outputs are input‑only sinks).

### 6.3 Libraries (`src/lib/`)
| File | Purpose |
|---|---|
| `pipelineMapping.js` | Convert canvas ⇄ backend Pipeline (zones, trigger, depends_on) |
| `functionMeta.js` | Icons/colors/categories per function type |
| `functionDocs.js` | Per‑function description + field labels/help/examples (guided panel) |
| `functionDefaults.js` | Default config seeded when a function node is added |
| `connectorTemplates.js` | Default config per connector + trigger type |
| `kgGraph.js` | Map KG JSON → Cytoscape elements + styles |
| `dashboardData.js` | Join domain graph → events (cause/cost), today/yesterday |
| `useLiveSocket.js` | WebSocket hook (auto‑reconnect) for live push |

### 6.4 State & API
- `src/store/studioStore.js` (Zustand) carries cross‑page intents: *Connect → apply
  connector to trigger*, *Pipelines → load pipeline into Compose*.
- `src/api/client.js` — all `fetch` calls to `/api/*`.

---

## 7. Frontend ↔ Backend connection

The browser **only ever talks to the API server** (`cmd/server`, port `:8080`).
There are exactly **two channels**:

1. **REST over HTTP** — request/response for everything (load + actions).
2. **WebSocket** — one persistent connection for live server→browser push.

### 7.1 Request lifecycle (dev)
```
Browser (React, :5173)
   │  fetch('/api/...')              │  new WebSocket('ws://host/api/ws')
   ▼                                 ▼
Vite dev server (:5173)  ── proxy '/api' (http + ws:true) ──►  Go API server (:8080)
                                                                 ├─ REST handlers (mux)
                                                                 └─ wsHub  (/api/ws)
```
In dev, Vite proxies every `/api/*` call (including the WebSocket upgrade) to
`:8080`, so the browser uses same‑origin URLs and there are no CORS issues. In a
production build you serve the static frontend and point it at the API host
(CORS is already enabled server‑side via `withCORS`).

### 7.2 Which scripts are responsible

**Frontend (`frontend/pipeline-builder/`)**
| File | Responsibility |
|---|---|
| `vite.config.js` | Dev proxy: forwards `/api` → `http://localhost:8080`, with `ws: true` so the `/api/ws` **WebSocket** upgrade is proxied too. |
| `src/api/client.js` | **All REST calls.** Wraps `fetch('/api/...')` — `fetchFunctions`, `fetchConnectors`, `fetchPipelines`, `createPipeline`, `runPipeline`, `fetchTags`, `fetchMachines`, `fetchTopics`, `fetchConfig`, `fetchStats`, `fetchKnowledgeGraph`. |
| `src/lib/useLiveSocket.js` | **WebSocket client.** React hook that connects to `/api/ws`, auto‑reconnects, and calls back on each `{type,data}` message. |
| pages/components | Call `client.js` for data and (Dashboard) `useLiveSocket` for live push. |

**Backend (`cmd/server/`)**
| File | Responsibility |
|---|---|
| `main.go` | Creates the `http.ServeMux`, registers every `/api/*` route, wraps it in `withCORS`, and serves on `:8080` (`http.ListenAndServe`). Contains the REST handlers. |
| `ws.go` | `wsHub` — upgrades `/api/ws` to a WebSocket and broadcasts `{type,data}` to all connected clients (per‑client write deadline so a slow/dead client can't stall the feed). |
| `live.go` | `LiveHub` — subscribes to MQTT `mindset/#` and calls `wsHub.broadcast(...)` so tag/state/event changes are pushed to the browser in real time. |
| `tags.go` | Tag registry behind `GET /api/tags`. |

### 7.3 End‑to‑end example (a micro‑stop appears)
```
Agent → MQTT mindset/events/micro-stop
   → cmd/server live.go (LiveHub) receives it
   → ws.go wsHub.broadcast({type:"event", ...})
   → browser useLiveSocket onmessage
   → DashboardPage triggers a debounced refresh()
   → client.js re-fetches /api/stats + /api/kg/domain
   → KPIs/charts update (~sub-second)
```

> Note: the OPC‑UA connection, tag discovery, raw publishing and UNS
> contextualization currently run inside **`cmd/agent`** and are **not yet
> controllable from the frontend** — the UI observes their results via the API.
> Adding a control plane (`POST /api/opcua/connect`, `/browse`, `/subscribe`,
> `/uns/start`) to `cmd/server` would let the browser drive them directly.

---

## 8. Project structure

```
mindset-data-edge/
├── cmd/
│   ├── agent/            # Edge runtime (OPC-UA → MQTT → UNS → rules → KG)
│   │   ├── main.go
│   │   └── init.go
│   └── server/           # HTTP API the UI talks to
│       ├── main.go       # routes, deps, handlers
│       ├── tags.go       # TagRegistry (+ SQLite persistence)
│       └── live.go       # LiveHub: topics rates + machine state
├── internal/
│   ├── config/           # agent.yaml loader
│   ├── discovery/        # OPC-UA browse + subscribe + publish raw
│   ├── mqtt/             # MQTT publisher
│   ├── uns/              # ISA-95 mapper + contextualizer
│   ├── rules/            # Run/Stop detection + StateStore
│   ├── functions/        # function registry + connectors/transforms/...
│   ├── pipeline/         # YAML loader, registry, execution engine
│   ├── kg/               # Knowledge Graph (domain + technical + subscriber)
│   └── storage/          # SQLite store + schema
├── config/
│   ├── agent.yaml        # site, opcua, mqtt, cost
│   └── pipelines/*.yaml  # predefined pipelines
├── frontend/pipeline-builder/   # React + Vite web studio
│   ├── src/{pages,components,lib,store,api}/
│   ├── vite.config.js    # proxy /api -> :8080
│   └── package.json
├── data/mindset.db       # SQLite (gitignored)
├── docs/                 # this file + design notes
├── run.ps1               # one-command launcher (build + start all)
├── go.mod / go.sum
└── .gitignore
```

---

## 9. Running it

### One command (Windows)
```powershell
.\run.ps1            # builds bin/*.exe, starts server + agent + frontend, opens the UI
.\run.ps1 -NoAgent   # UI + API only (no OPC-UA/MQTT needed)
.\run.ps1 -NoBuild   # reuse existing binaries
```

### Manually
```powershell
go run ./cmd/server                              # API on :8080
cd frontend/pipeline-builder; npm run dev        # UI on :5173
go run ./cmd/agent                               # edge runtime (needs OPC-UA + MQTT)
```

### Build .exe
```powershell
go build -o bin/server.exe ./cmd/server
go build -o bin/agent.exe  ./cmd/agent
```

**Prerequisites for live data:** an MQTT broker on `:1883` and (for real tags) the
Prosys OPC‑UA simulator at the endpoint in `config/agent.yaml`. Without them the UI
still runs; live values/rates/state are limited and persisted tags are shown.

---

## 10. Feature log (what was built, in order)

1. **Standalone API server** (`cmd/server`) — runs the whole UI backend with no
   OPC‑UA/MQTT required to start.
2. **Drag‑and‑drop builder** (Compose) — ReactFlow canvas, palette, save → YAML.
3. **Connector selection** — pick a connector, applied to the trigger; templates.
4. **In‑app Knowledge Graph** — Cytoscape, Technique/Domaine toggle + filters.
5. **Real‑time dashboard** — KPIs, recent events, machine status, Gantt timeline.
6. **Agent HTTP removed** — `cmd/agent` is purely the edge runtime now.
7. **Live OPC‑UA tag bridge** — server subscribes `mindset/raw/#`, `/api/tags`,
   tag picker shows real tags + values (persisted to SQLite).
8. **Real‑data bridge** — `/api/machines` (status), `/api/topics` (rates/category),
   `/api/config` (from `agent.yaml`); pickers use real data.
9. **WebSocket live push** — `/api/ws`; the dashboard updates on push.
10. **Live tag chart** (`LiveDataPanel`) — pick tags, watch values stream.
11. **KG = only your pipelines**; predefined ones become **loadable templates**
    (`/api/pipelines/examples`).
12. **`add_to_dashboard`** output + **pinned widgets** panel (snapshot + live).
13. **Guided builder** — descriptions, labelled fields with help/examples, OPC‑UA
    machine+tag selection, smart validation, **duplicate prevention**.
14. **Cost config** — rate source (Manual/Config/Tag), currency, **CSV/Excel**
    per‑product rates, live preview.
15. **Output sinks** (input‑only ports), **KG tag grouping**, **Pareto removed**.
16. **Interactive dashboard widgets** (`DashboardWidgets`) — add from data sources,
    chart type (line/bar/gauge/value/status), time range, live stats, localStorage.
17. **MindSet Data logo** in the NavBar; **WebSocket write‑deadline** hardening.

## 11. Notes & limitations
- **TRS** on the dashboard is an *availability proxy* (downtime vs an 8h shift);
  true OEE needs production‑count + quality data not yet captured.
- **"vs hier"** deltas need events spanning two days to show a baseline.
- **Modbus/SQL** connectors are demo stubs (metadata only).
- Live values/rates/state require the agent to be actively publishing; otherwise
  tags persist (SQLite) but rates read 0 and state is last‑known.
- The dashboard is **WebSocket‑driven** with a **20s polling fallback**.
- Pipeline **Run** executes on demand via the server; the agent's live data flow
  (raw→UNS→events) runs independently in `cmd/agent`.











HTTP polling, not WebSocket/SSE.

How "real-time" actually works here
There are two different links, and only the backend↔broker one is true push:

1. Frontend ↔ API Server = REST over HTTP (polling)

All UI data goes through fetch() calls in src/api/client.js, proxied by Vite (/api → :8080).
The Dashboard creates the "live" feel by polling every 5 seconds:
// DashboardPage.jsx
const REFRESH_MS = 5000;
timer.current = setInterval(refresh, REFRESH_MS);
The builder pickers refresh on demand (e.g. tags/machines reload each time you select a node).
No WebSocket, no SSE between browser and server.
2. API Server ↔ MQTT Broker = real push (MQTT subscriptions)

This is where the genuine real-time is. The server subscribes to the broker and keeps in-memory state fresh:
LiveHub (cmd/server/live.go) subscribes to mindset/# → updates TagRegistry, TopicRegistry (rates), StateTracker.
KGSubscriber subscribes to mindset/events/micro-stop → auto-enriches the KG.
So the pattern is:
Agent ──MQTT push──► Broker ──MQTT push──► Server (in-memory live state)
                                              ▲
                              Browser ──HTTP poll every 5s──┘
The server is updated in real time by MQTT; the browser samples that state every 5s via REST. So end-to-end latency is "broker push + up to 5s polling delay."
