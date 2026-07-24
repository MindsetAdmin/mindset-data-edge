# MindSet Data Edge

Industrial IoT edge platform that turns OT (OPC-UA) and IT (SQL/ERP) data into a live, structured Knowledge Graph — with AI agents and automation built on top, not bolted on.

**What's new:** the Knowledge Graph now **auto-generates its structural skeleton at OPC-UA connect time**, via ISA-95 mapping, with mandatory human validation before anything counts as confirmed. Context exists from the moment you connect — not after weeks of manually building pipelines. Live-tested end-to-end against a real Prosys OPC-UA server. See [§6](#6-the-structural-bootstrap-how-it-works) for the full flow.

---

## 1. Installation

**Prerequisites**
- Go (see `go.mod`) and Node.js (`npm`)
- Docker Desktop — only if you want the SQL connector / fake-ERP demo
- An MQTT broker on `:1883` — not bundled. Quickest option:
  ```powershell
  docker run -d --name mosquitto -p 1883:1883 eclipse-mosquitto:2 mosquitto -c /mosquitto-no-auth.conf
  ```
- An OPC-UA source — a real PLC/SCADA system, or the free [Prosys OPC-UA Simulation Server](https://www.prosysopc.com/products/opc-ua-simulation-server/) for testing.

**Build**
```powershell
go build -o bin/server.exe ./cmd/server
go build -o bin/agent.exe  ./cmd/agent
cd frontend/pipeline-builder && npm install
```

**Run everything**
```powershell
.\run.ps1            # builds + starts server, agent, frontend; opens the browser
.\run.ps1 -NoAgent    # skip the edge agent (fine for UI-only work)
```

Opens the UI at `http://localhost:5173`, API at `http://localhost:8080`.

## 2. Configuration

`config/agent.yaml` — the file that matters most for the structural bootstrap:

```yaml
site:
  id: "local-test"                                    # becomes the graph's Site node

opcua:
  endpoint: "opc.tcp://<host>:53530/OPCUA/Server1"     # your OPC-UA source
  security_mode: "None"                                # Sign/SignAndEncrypt not yet wired
  auto_connect: false                                  # OPC-UA is driven from the UI, not the agent

mqtt:
  broker: "tcp://localhost:1883"
```

`config/connections.yaml` — SQL connections for the (separate, IT-side) `sql_query` connector. Not required for the OT structural bootstrap.

## 3. Usage — the structural bootstrap, step by step

1. Start the stack (`.\run.ps1`), with an MQTT broker and an OPC-UA source reachable.
2. In the UI, go to **Connect → OPC-UA**. Enter your endpoint, click **Connect**.
3. Click **Discover tags**. This browses the server's full node tree — and, as a side effect, **auto-generates the site's structural skeleton** (Site → Area → Work Center → Equipment → Tag) into the Knowledge Graph, using the existing ISA-95 tag-naming mapper. Every generated node is flagged pending, not live.
4. Go to **Knowledge Graph**. The sidebar shows a **"Pending validation (N)"** list — every auto-generated node, with **Accept**/**Reject** buttons. Unvalidated nodes render in the graph with a dashed amber ring, distinct from confirmed structure.
5. Review and accept (or reject) each node. Accepted nodes become ordinary graph nodes; rejected ones are deleted.
6. From here, build pipelines in **Compose** to automate and further enrich — the graph already exists, pipelines don't create it from scratch.

## 4. API endpoints (structural bootstrap)

| Method & path | Purpose |
|---|---|
| `POST /api/opcua/connect` | Connect to an OPC-UA endpoint |
| `GET /api/opcua/discover` | Browse the server's tags — **triggers the KG seed as a side effect** |
| `GET /api/kg/pending` | List business-category nodes awaiting validation |
| `POST /api/kg/pending/{id}/validate` | Confirm an auto-generated node |
| `POST /api/kg/pending/{id}/reject` | Discard an auto-generated node (deletes it + its direct edges) |
| `GET /api/kg?category=business\|platform\|all` | Read the unified Knowledge Graph |

Full API surface (pipelines, connections, dashboards, etc.) is documented in `CLAUDE.md`.

## 5. UI

| Page | Route | Role in this flow |
|---|---|---|
| Connect → OPC-UA | `/connect/opcua` | Connect + browse an OPC-UA source; triggers the seed |
| Knowledge Graph | `/kg` | Graph viewer with the **Pending validation** list; auto-generated nodes render with a dashed amber ring until accepted |
| Compose | `/compose` | Build pipelines that automate and enrich the (already-existing) graph |

## 6. The structural bootstrap — how it works

```
Connect (OPC-UA)
    │
    ▼
Discover  ──►  BrowseNodeTree()  ──►  full tag/node tree, in one call
    │
    ▼
ISA-95 mapping (existing tag-naming mapper)
    │
    ▼
SeedFromDiscovery()  ──►  writes Site → Area → Work Center → Equipment → Tag
                          into the Knowledge Graph, every node flagged `pending: true`
    │
    ▼
Human validation (UI: Pending validation list)
    │
    ├── Accept  ──► node becomes confirmed structure
    └── Reject  ──► node deleted
```

**Design principle:** the structural skeleton (what equipment exists, how it's organized) is auto-generated because that information already exists at connect time — no reason to make a human type it in. What can *never* be auto-generated is operational/transactional data (work orders, quality events, costs) — that arrives progressively, through pipelines, because it doesn't exist until it happens. Pipelines automate and enrich an already-existing graph; they don't bootstrap it from zero.

## 7. Tech stack

- **Backend:** Go, `modernc.org/sqlite` (pure-Go, no CGO), `gopcua` (OPC-UA client), Eclipse Paho (MQTT)
- **Frontend:** React 19, Vite, ReactFlow (pipeline canvas), a force-directed graph viewer (Knowledge Graph), Zustand, Tailwind
- **Storage:** SQLite — one unified graph table set (`kg_nodes`/`kg_edges`), category-tagged (`business` vs `platform`)
- **Protocols:** OPC-UA (OT), MQTT (internal bus), MySQL (IT/ERP, separate connector)

## 8. Known limitations

- **Flat-list validation doesn't scale-test past a small demo server.** Verified usable at 8 tags / 15 pending nodes. A real industrial OPC-UA server can expose hundreds of tags — whether a flat accept/reject list stays usable at that volume is untested and likely needs a grouped/tree view.
- **IT-side (ERP) master data isn't auto-generated.** The structural bootstrap is OT-only today. Master data (assets, materials, products) *could* in principle be bootstrapped similarly, but that needs a schema-discovery capability the SQL connector doesn't have yet — a separate, larger effort.
- **Rejecting a node doesn't cascade.** Rejecting a Site or Area leaves its children as orphaned-but-still-pending entries in the flat list (individually rejectable, just not automatic). Not an issue at flat-list scale; would need addressing if this becomes a tree view.
- **Auto-generated and reactive nodes can silently coexist unmerged.** If the bootstrap runs before a real operational event (a micro-stop) references the same equipment, the node stays `pending` even after real data starts flowing, until explicitly validated.
- Also inherited from the base platform: secure OPC-UA modes (Sign/SignAndEncrypt) aren't wired up; the OPC-UA session holds one subscription at a time.

Full technical detail, architecture, and the complete decision history behind this feature: `CLAUDE.md`, `docs/ARCHITECTURE.md`, `docs/new_member_guide.md`, and `docs/analysis_log.md` (Entries 87–98).
