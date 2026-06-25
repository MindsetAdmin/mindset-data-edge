# MindSet Data — Component Reference

A file‑by‑file description of what every script does, reflecting all current
features. For the high‑level design and data flow, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## A. Binaries (`cmd/`)

### `cmd/agent/` — the edge runtime (OPC‑UA → MQTT → UNS → rules → KG)
| File | What it does |
|---|---|
| `main.go` | Boots the **agent**: config → MQTT publisher → UNS contextualizer → rules engine → KG subscriber → functions/pipelines → OPC‑UA connect → browse tags → subscribe to value changes. **No HTTP**. |
| `init.go` | Helper constructors wiring each subsystem with logging. |

### `cmd/server/` — the HTTP/WebSocket API the web UI talks to
| File | What it does |
|---|---|
| `main.go` | Builds the function registry, opens the KG, connects MQTT, starts the KG auto‑enrich subscriber and the LiveHub, registers all `/api/*` routes + `/api/ws`, wraps in CORS, serves on `:8080`. Holds the REST handlers (functions, connectors, pipelines list/save/run + examples, tags, machines, topics, config, dashboard pins, kg technical/domain, stats, health). |
| `tags.go` | **`TagRegistry`** — latest value per OPC‑UA tag (from `mindset/raw/#`), persisted to the SQLite `tags` table. Backs `GET /api/tags`. |
| `live.go` | **`LiveHub`** — one `mindset/#` subscription. Feeds `TagRegistry`, `TopicRegistry` (msg/s rates + category), `StateTracker` (Running/Stopped + transition history), and keeps the latest **dashboard pin** per label. Broadcasts `tag`/`state`/`event`/`dashboard` over WebSocket. |
| `ws.go` | **`wsHub`** — upgrades `/api/ws`, tracks clients, `broadcast(type,data)` fans out JSON to all (concurrency‑safe, with a per‑client **write deadline** so one slow/dead client can't stall the feed). |

---

## B. Internal packages (`internal/`)

### `config/`
| File | What it does |
|---|---|
| `config.go` | Loads `config/agent.yaml`: `Site`, `OpcUA`, `Discovery`, `Mqtt`, `Cloud`, `Cost`. |

### `discovery/` — OPC‑UA
| File | What it does |
|---|---|
| `opcua.go` | Connect, `BrowseNodeTree` (discover tags + types), `Subscribe` (publish each value change to `mindset/raw/<nodeID>`), `WatchForChanges`. |

### `mqtt/`
| File | What it does |
|---|---|
| `publisher.go` | `PublishRaw` → `mindset/raw/<nodeID>`; `PublishEvent` → `mindset/events/<type>`. |

### `uns/` — Unified Namespace (ISA‑95)
| File | What it does |
|---|---|
| `mapper.go` | Tag → ISA‑95 node: normalize name, infer unit, build `mindset/site/...` topic. |
| `contextualizer.go` | Subscribe `mindset/raw/#` → map → republish to `mindset/site/...`. |

### `rules/`
| File | What it does |
|---|---|
| `engine.go` | Detect Run↔Stop transitions; publish `status-change` / micro‑stop events. |
| `state.go` | `StateStore` — thread‑safe current state + history per tag. |

### `functions/` — pipeline building blocks
| File | What it does |
|---|---|
| `registry.go` | Function registry; `Register/Get/List` + `ListFunctions`/`ByType` (UI info). |
| `types.go` | `FunctionType` constants + `FunctionInfo`/`Param`. |
| `connectors/opcu-read.go` | `opcua_read` — read an OPC‑UA node. |
| `connectors/mqt_subscribe.go` | `mqtt_subscribe` — subscribe to a topic (trigger). |
| `transforms/state_machine.go` | `state_machine` — Run/Stop transitions. |
| `transforms/uns_mapper.go` | `uns_mapper` — normalize to ISA‑95. |
| `transforms/filter.go` | `filter` — keep/drop by field/operator/value. |
| `calculates/duration.go` | `calculate_duration` — start→stop duration. |
| `calculates/cost.go` | `calculate_cost` — € cost; supports a per‑product **rate table** (`config.rates` keyed by the event's `product`). |
| `conditions/threshold.go` | `threshold` — value within `[min,max]`. |
| `outputs/mqtt_publish.go` | `mqtt_publish` — publish to a topic. |
| `outputs/dashboard.go` | `add_to_dashboard` — pin to the dashboard via `mindset/dashboard/<label>` (retained). |
| `outputs/kg_save.go` | `kg_save` — write event to KG (agent‑only; not in the server palette). |

### `pipeline/`
| File | What it does |
|---|---|
| `types.go` | `Pipeline`/`Node`/`Trigger` + execution types (timestamps `yaml:"-"`). |
| `loader.go` | Load `*.yaml` (non‑recursive), validate id/name, honor `enabled`. |
| `registry.go` | In‑memory registry + `GetHash`. |
| `engine.go` | Execute in dependency order (`"trigger"` dep = satisfied); panic‑guarded handler calls; emits events. |
| `builder.go` | Pipeline‑building helpers. |

### `kg/`
| File | What it does |
|---|---|
| `graph.go` | Domain KG CRUD over SQLite (`AddNode/AddEdge/AddMicroStop/AddCause/AddCost`, `GetFullGraph`). |
| `builder.go` | Technical KG — derived **only from your pipelines**; **groups** pipelines with the same function signature into one node listing all `tags`. External‑only model. |
| `subscriber.go` | Subscribe `mindset/events/micro-stop` → auto‑create Equipment/Event/Cause/Cost nodes (skips empty events). |
| `types.go` | `TechnicalNode/Edge/Graph` + type constants. |

### `storage/`
| File | What it does |
|---|---|
| `sqlite.go` | Open `data/mindset.db` (pure‑Go driver) + schema (`kg_nodes`, `kg_edges`, `events`); `tags` table created by the server. |

---

## C. Configuration (`config/`)
| Path | What it does |
|---|---|
| `agent.yaml` | Site, OPC‑UA, MQTT broker, cost rate. |
| `pipelines/*.yaml` | **Your** pipelines (engine + KG). Empty by default. |
| `pipelines/examples/*.yaml` | Shipped **templates** (load in the UI; not in engine/KG). |

---

## D. Frontend (`frontend/pipeline-builder/`)

### Entry / config
| File | What it does |
|---|---|
| `index.html` | HTML shell, title "MindSet Data". |
| `public/logo.png` | The MindSet Data logo, served at `/logo.png` (used in the NavBar). |
| `vite.config.js` | Dev `:5173`; proxy `/api` → `:8080` with `ws:true` (REST + WebSocket). |
| `src/main.jsx` | React entry; `BrowserRouter`. |
| `src/App.jsx` | Router shell (NavBar + 6 routes inside `ErrorBoundary`). |

### API & state
| File | What it does |
|---|---|
| `src/api/client.js` | All REST calls (functions, connectors, pipelines + examples + run, tags, machines, topics, config, **dashboard pins**, stats, KG). |
| `src/store/studioStore.js` | Zustand: cross‑page intents (Connect→Compose connector; Pipelines→Compose full pipeline object). |
| `src/lib/useLiveSocket.js` | WebSocket hook (`/api/ws`, auto‑reconnect). |

### Libraries (`src/lib/`)
| File | What it does |
|---|---|
| `pipelineMapping.js` | Canvas ⇄ backend Pipeline (zones, trigger, `depends_on`). |
| `functionMeta.js` | Icon/color/category per function type. |
| `functionDocs.js` | Description + per‑field label/help/example (guided panel). |
| `functionDefaults.js` | Default config seeded when a function node is added. |
| `connectorTemplates.js` | Default config + trigger type per connector. |
| `kgGraph.js` | KG JSON → Cytoscape elements/styles; drops dangling edges. |
| `dashboardData.js` | Domain graph → events (cause/cost), today/yesterday split. |

### Components (`src/components/`)
| File | What it does |
|---|---|
| `NavBar.jsx` | 6‑tab navigation + the **MindSet Data logo** (`public/logo.png`). |
| `Palette.jsx` | Draggable function blocks (connectors excluded). |
| `NodeConfigPanel.jsx` | Guided config: header (icon + category badge) + description, labelled fields with help/examples, pickers, **OPC‑UA machine/tag selector**, **cost source + CSV/Excel rate upload + live preview**, delete. |
| `PickerModal.jsx` | Generic searchable chooser. |
| `CytoscapeGraph.jsx` | React wrapper around Cytoscape. |
| `ErrorBoundary.jsx` | Catches render errors. |
| `LiveDataPanel.jsx` | Pick tag(s) → live multi‑line chart over WebSocket. |
| `DashboardWidgets.jsx` | **Interactive widgets** for `add_to_dashboard` data: add from available sources, pick chart type (line/bar/gauge/value/status) + time range (1m–24h), live stats (Last/Min/Max/Avg/Count), `✕`/`⚙️` controls, persisted in **localStorage**. Parses values (never raw JSON). |
| `nodes/PipelineNode.jsx` | Pipeline‑step node; **outputs are input‑only sinks** (no output port). |
| `nodes/TriggerNode.jsx` | The entry/trigger node. |
| `nodes/ZoneNode.jsx` | ENTRÉE/CŒUR/SORTIE background bands. |

### Pages (`src/pages/`)
| File | Route | What it does |
|---|---|---|
| `OverviewPage.jsx` | `/overview` | Landing: stats + quick links. |
| `ConnectPage.jsx` | `/connect` | Connector catalog; select → applies to trigger. |
| `BuilderPage.jsx` | `/compose` | Drag‑and‑drop builder; guided config; **smart validation + duplicate modal**; Save/Run/delete. |
| `PipelinesPage.jsx` | `/pipelines` | "Mes pipelines" (run/load) + "Modèles (exemples)" (load). |
| `DashboardPage.jsx` | `/dashboards` | Real‑time dashboard (WebSocket + 20s fallback): KPIs, pinned widgets, live tag chart, recent events, machine status, Gantt. |
| `KnowledgeGraphPage.jsx` | `/kg` | Cytoscape viewer, Technique/Domaine toggle + filters. |

---

## E. Root
| File | What it does |
|---|---|
| `run.ps1` | One‑command launcher (build + start server/agent/frontend; `-NoBuild`, `-NoAgent`). |
| `go.mod` / `go.sum` | Go deps (paho MQTT, gopcua, gorilla/websocket, modernc sqlite, yaml). |
| `.gitignore` | Ignores `bin/`, `*.exe`, `data/`, `node_modules/`, `dist/`, local settings. |
| `docs/` | `ARCHITECTURE.md` (design + data flow), this file, and design notes. |
