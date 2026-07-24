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
| `functions/{connectors,transforms,calculates,conditions,outputs}` | Concrete functions. `sql_query` is fully implemented (not a stub) — see §4.4. |
| `pipeline` | `types.go`, `loader.go` (YAML), `registry.go`, `engine.go` (topological exec), `builder.go`. |
| `kg` | Knowledge Graph: `graph.go` (unified CRUD, category‑tagged), `builder.go` (platform‑category rebuild), `subscriber.go` (auto‑enrich business category from micro‑stops), `types.go`. See §4.6 — this is now **one graph**, not two. |
| `storage` | SQLite store + schema. |
| `connections` | `Registry` pooling `*sql.DB` per SQL connection (MySQL only) for `sql_query`; lazy‑open + verify, runtime add/remove via `/api/connections`, `password_env` indirection (no inlined secrets). |
| `e2e` | Build‑tagged integration tests (`-tags=integration`) running `sql_query` against a real MySQL via testcontainers. |

### 4.4 The function catalog
| Function | Type | Purpose |
|---|---|---|
| `opcua_read` | connector | Read a value from an OPC‑UA node |
| `mqtt_subscribe` | connector | Subscribe to an MQTT topic (pipeline trigger) |
| `modbus_read` | connector | **demo stub** — errors if executed |
| `sql_query` | connector | **Fully implemented**, not a stub. Read‑only parameterized SELECT against a pooled MySQL connection (`internal/connections.Registry`); type‑coerces results and, with `field_map`/`value_map` configured, also returns a `canonical` copy whose object set aligns with ISA‑95's information model (`docs/mysql_connector.md` §6b; naming precision decided in `docs/analysis_log.md` Entry 92 — "aligns with," never "compliant," since the B2MML wire format isn't adopted). No `{{ trigger.field }}` templating — `params` are static values only. |
| `state_machine` | transform | Detect Run↔Stop transitions |
| `filter` | transform | Keep/drop by condition |
| `calculate_duration` | calculate | Duration between start/stop |
| `calculate_cost` | calculate | € cost from duration × hourly rate |
| `threshold` | condition | Is value within [min,max]? (micro‑stop window) |
| `add_to_dashboard` | output | Pin the data/event onto the Dashboard (`mindset/dashboard/<label>`) |
| `kg_save` | output | Save to KG — **not** in the server palette (KG is automatic) |

**Removed (`docs/analysis_log.md` Entry 119):** `uns_mapper` (duplicated what `OPCUAManager.route()` already does automatically) and `mqtt_publish` (a pipeline's terminal node now auto‑publishes to MQTT without an explicit node — see §4.5).

Outputs are **sinks**: in the builder they have an input port only (no output
port). `calculate_cost` supports a rate **source** (Manual / from `agent.yaml` /
from a tag) and a per‑product **rate table** uploaded from CSV/Excel
(`config.rates`, keyed by the event's `product`).

### 4.5 Pipelines
A pipeline is YAML in `config/pipelines/` with a **trigger** + **nodes** (each has
`function`, `config`, `depends_on`). The engine (`pipeline/engine.go`) runs nodes
in dependency order (a dependency that isn’t a node — e.g. `"trigger"` — counts as
satisfied). After a successful run, the pipeline's declared `output` node's result
auto‑publishes to MQTT (`cmd/server/pipeline_output.go`) — topic from the optional
`output_topic` YAML field if set, else auto‑derived as
`mindset/pipelines/<id>/output`. No output node required for this. Built‑ins:
- `microstop_detection` — `mqtt_subscribe → state_machine → calculate_duration → threshold` (auto‑publishes to `mindset/events/micro-stop`)
- `cost_calculation` — `mqtt_subscribe → calculate_cost` (auto‑publishes to `mindset/events/micro-stop-cost`)

### 4.6 Knowledge Graph: one graph, two categories (updated 2026‑07‑20 — see `docs/analysis_log.md` Entry 50 for the merge, Entries 87/89/90 for the bootstrapping gap below)

**This section previously described two separate graphs (a SQLite "domain" graph + a 5‑minute‑cached in‑memory "technical" graph). That's no longer accurate.** As of the 2026‑07‑02 merge, it's **one SQLite‑backed graph** (`kg/graph.go`), every node/edge tagged with a `category` column:

- **`business`** — the *data*: `Equipment`, `Event` (micro‑stops), `Cause`, `Cost`, `Operator`, `Product`, `OF`. Auto‑enriched by `KGSubscriber`.
- **`platform`** — the *architecture*: Pipeline/Function/Connection/Topic/Dashboard, derived only from pipelines you build (empty until you save one — shipped examples don't count). Rebuilt by `RepopulatePlatform`, no‑op'd by a registry‑hash check (the old 5‑minute cache is gone).

Read through one function, `GetGraph(category)` → `GET /api/kg?category=business|platform|all`. `/api/kg/technical` and `/api/kg/domain` still work as thin legacy aliases.

**Known gap — the `business` category has no structure‑discovery bootstrapping path.** `Equipment` nodes are created in exactly one place (`kg/subscriber.go`, triggered only by `mindset/events/micro-stop`) — nothing connects OPC‑UA discovery (`internal/discovery.BrowseNodeTree`) to the graph, and the ISA‑95 mapper (`internal/uns/mapper.go`) that could seed it only shapes MQTT topic names today. Same gap on the SQL side: `sql_query`'s `canonical` output (§4.4) isn't consumed by anything either, despite `docs/mysql_connector.md` describing it as if it were. **Proposed fix, agreed in principle, not yet built:** auto‑generate the Equipment/Area/Site skeleton at OPC‑UA connect time via the existing ISA‑95 mapper, gate it behind mandatory human validation before it's live, and wire `sql_query`'s canonical output (object set expanded per Entry 92 — `WorkOrder`/`Batch`/`Product`/`Schedule`/`Quality`/`Operator`/`Material`/`Asset`/`ProcessSegment`) into the same bridge. `Asset` specifically is meant to entity-resolve against this same `Equipment` node type once built. Full analysis in `docs/analysis_log.md` Entries 87, 89, 90, 92.

### 4.7 SQLite schema (`data/mindset.db`)
| Table | Columns | Written by |
|---|---|---|
| `kg_nodes` | id, category, type, label, properties(JSON), created_at | KG subscriber / graph |
| `kg_edges` | id, category, from_id, to_id, relation, weight, created_at | KG subscriber / graph |
| `tags` | node_id, name, value(JSON), data_type, timestamp_ms, updated_at | server TagRegistry |
| `events` | id, type, work_center, duration_seconds, cause, cost_eur, payload, timestamp | (reserved) |
| `connections` | SQL connection definitions created via `/api/connections` | Connections API |

Legacy DBs get `category` backfilled to `'business'` by a migration in `internal/storage/sqlite.go`. Note also (Entry 89): `confidence` is not a general primitive — it only exists on `Cause` nodes/edges (`AddCause`, reusing the generic edge `weight` column); every other edge type hardcodes `weight = 1.0`.

---

## 5. HTTP API (server `:8080`)

| Method & path | Returns |
|---|---|
| `GET /api/health` | `{status:"ok"}` |
| `GET /api/config` | safe subset of `agent.yaml` (opcua endpoint/security, broker, hourly rate, site/area) |
| `GET /api/functions[?type=]` | function catalog (optionally filtered by type) |
| `GET /api/connectors` | connector functions only (thin `type=connector` alias over `/api/functions`) |
| `GET /api/pipelines` | your pipelines loaded from `config/pipelines/` |
| `POST /api/pipelines` | save a pipeline as YAML (validated, id sanitized) |
| `GET /api/pipelines/examples` | shipped template pipelines (`config/pipelines/examples/`) |
| `POST /api/pipelines/{id}/run` | execute a pipeline; returns per‑node status + timing |
| `GET /api/tags` | live OPC‑UA tags + values (persisted) |
| `GET /api/machines` | tags grouped by work center + live Running/Stopped state |
| `GET /api/topics` | live topics + msg/s + category + broker_connected |
| `GET /api/config` | safe subset of `agent.yaml` |
| `POST /api/opcua/connect` | connect to a user‑specified OPC‑UA server (endpoint/security/auth/timeout) |
| `GET /api/opcua/discover` | browse the connected server's node tree |
| `POST /api/opcua/subscribe` | monitor selected tags with per‑tag mode (`raw`\|`isa95`\|`both`) |
| `POST /api/opcua/disconnect` | close the OPC‑UA session |
| `GET /api/opcua/status` | connection status (status/endpoint/tag_count/error) |
| `GET /api/opcua/selections` | current per‑tag routing + ISA‑95 mapping (builder governance) |
| `GET /api/dashboard/pins` | current dashboard pins (from `add_to_dashboard`) |
| `GET/POST /api/connections` | list SQL connections (never returns passwords) / create‑or‑replace one |
| `POST /api/connections/{id}/test` | force a fresh read‑only health check |
| `POST /api/connections/{id}/preview` | run a query through the same guards as `sql_query`, capped at 5 rows |
| `DELETE /api/connections/{id}` | remove a connection |
| `GET /api/kg?category=business\|platform\|all` | **unified** KG read (default `all`) — see §4.6 |
| `GET /api/kg/technical` | legacy alias → `category=platform` |
| `GET /api/kg/domain` | legacy alias → `category=business` |
| `GET /api/stats` | counts + micro‑stops/downtime/cost + uptime + broker status |
| `WS  /api/ws` | **WebSocket** live push: `{type:"tag"\|"state"\|"event"\|"dashboard"}` |

CORS is enabled; in dev the Vite proxy forwards `/api` (HTTP **and** the
WebSocket upgrade, `ws:true`) → `:8080`.

---

## 6. Frontend (React + Vite)

Stack: **React 19**, **Vite**, **React Router**, **Zustand** (cross‑page state),
**ReactFlow** (canvas), **ForceGraph** (KG — see note below), **Recharts** (charts), **xlsx**
(CSV/Excel parsing), **WebSocket** (live push), **Tailwind**.

> **Correction:** the KG viewer now renders via `ForceGraph.jsx`, consuming the unified `/api/kg?category=` endpoint directly. `CytoscapeGraph.jsx` still exists in the repo but is dead code — not imported anywhere.

### 6.1 Pages (`src/pages/`)
| Page | Route | What it does |
|---|---|---|
| `OverviewPage` | `/overview` | Landing: key stats + quick links |
| `ConnectPage` | `/connect` | Pick a connector → applies to the pipeline trigger |
| `OpcuaConnectPage` | `/connect/opcua` | Dedicated OPC‑UA flow: connect → discover → per‑tag raw/ISA‑95/both routing → apply |
| `BuilderPage` | `/compose` | **Drag‑and‑drop builder** (ReactFlow): ENTRÉE/CŒUR/SORTIE bands, guided config panel, function/field pickers, **smart validation + duplicate prevention**, Save (→YAML), Run, delete node |
| `PipelinesPage` | `/pipelines` | Two sections: **your** pipelines (run/load) + **templates** (load) |
| `ConnectionsPage` | `/connections` | SQL connections: list + create + Test — backs the `sql_query` connector |
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
│   ├── kg/               # Knowledge Graph — one graph, category-tagged (see §4.6)
│   ├── storage/          # SQLite store + schema
│   ├── connections/      # SQL connection pool/registry (MySQL only) for sql_query
│   └── e2e/              # build-tagged integration tests (real MySQL via testcontainers)
├── config/
│   ├── agent.yaml        # site, opcua, mqtt, cost
│   ├── connections.yaml  # SQL connection definitions (dev_erp, etc.)
│   └── pipelines/*.yaml  # predefined pipelines
├── frontend/pipeline-builder/   # React + Vite web studio
│   ├── src/{pages,components,lib,store,api}/
│   ├── vite.config.js    # proxy /api -> :8080
│   └── package.json
├── sim/erp/               # fake ERP schema/seed/grants (MySQL)
├── cmd/erpsim/            # fake ERP data generator (advance/rotate/quality/plan loops)
├── docker-compose.dev-erp.yml   # fake ERP MySQL container
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

**Prerequisites for live data:** an MQTT broker on `:1883` (not bundled — e.g. `docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto:2 mosquitto -c /mosquitto-no-auth.conf`) and (for real tags) the
Prosys OPC‑UA simulator at the endpoint in `config/agent.yaml`. Without them the UI
still runs; live values/rates/state are limited and persisted tags are shown.

**For the SQL connector / fake ERP demo:**
```powershell
docker compose -f docker-compose.dev-erp.yml up -d   # MySQL on host :3308
$env:MINDSET_ERP_PASSWORD = "readonly_dev"             # matches mindset_readonly in sim/erp/grant.mysql.sql
go run ./cmd/erpsim                                     # background data generator
```
Then verify via the **Connections** page (**Test** on `dev_erp`) before running `config/pipelines/examples/of_enrichment.yaml`.

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
18. **Dynamic, user‑controlled OPC‑UA** — the API server owns the OPC‑UA session
    (`cmd/server/opcua.go`, `OPCUAManager`); the UI connects, browses and selects
    tags with per‑tag **Raw / ISA‑95 / Both** routing via `/api/opcua/*`. Raw is
    published to `mindset/raw/#` (storage); ISA‑95/Both are also mapped and
    published to `mindset/site/#` (functions). The agent no longer auto‑connects
    (`opcua.auto_connect=false`); it keeps running rules + KG and consumes the
    server‑published `site/#`. Builder function pickers are restricted to
    ISA‑95/Both work centers (`/api/opcua/selections`).
19. **MySQL connector V1a** (2026‑07‑06 to 07‑18) — `sql_query` fully implemented: `internal/connections.Registry`, `field_map`/`value_map` → `canonical` output aligned to ISA‑95's information model, `/api/connections` CRUD + `/test` + `/preview`, `ConnectionsPage`/`SqlConfigPanel`/`FieldMapEditor` frontend, the fake‑ERP dev stack (`cmd/erpsim`, `docker-compose.dev-erp.yml`), and `-tags=integration` tests against a real MySQL testcontainer. See `docs/mysql_connector.md` and `docs/analysis_log.md` Entries 58‑82.
20. **KG unified into one graph** (2026‑07‑02) — `business`/`platform` categories replace the old two‑graph (SQLite domain + in‑memory technical) design; see §4.6.

## 11. Notes & limitations
- **TRS** on the dashboard is an *availability proxy* (downtime vs an 8h shift);
  true OEE needs production‑count + quality data not yet captured.
- **"vs hier"** deltas need events spanning two days to show a baseline.
- **Modbus** is a demo stub (metadata only, errors if run). **SQL (`sql_query`) is fully implemented**, not a stub — see item 19 above.
- **The Knowledge Graph's `business` category has no structure‑discovery bootstrapping path** — `Equipment` nodes only appear reactively, from a micro‑stop event, not from OPC‑UA discovery or the SQL `canonical` output. See §4.6 and `docs/analysis_log.md` Entries 87/89/90/92 for the full analysis and the proposed (not yet built) fix.
- Live values/rates/state require the agent to be actively publishing; otherwise
  tags persist (SQLite) but rates read 0 and state is last‑known.
- The dashboard is **WebSocket‑driven** with a **20s polling fallback**.
- Pipeline **Run** executes on demand via the server; the agent's live data flow
  (raw→UNS→events) runs independently in `cmd/agent`.
- **Dynamic OPC‑UA** owns one session at a time; re‑applying selections stacks a new
  subscription, so changing the selection currently means reconnecting. Secure modes
  (`Sign`/`SignAndEncrypt`) need a client certificate (not yet wired) — use `None`.
  The session idle‑timeout defaults to **300s** so the browse→select gap can't drop it.







