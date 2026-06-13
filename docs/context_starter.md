# Mindset Data — Context Brief

## What we're building
Industrial data infrastructure for manufacturing ETI.
Edge Agent (Go + Docker) → UNS (ISA-95) → Decision dashboard.
Maximum local processing. Zero raw data to cloud.
Zero hardware, zero dev client, deployment in 48h.

## Team
- Mohamed (CTO) — Polytechnique, IoT & embedded systems (Windows PC)
- Cécilia (CEO) — EDHEC, ex-VC AgriFoodTech

## Stack
- Edge Agent: Go + Docker (Go 1.22+, static binary, cross-platform)
- V0 protocols: OPC-UA (`gopcua`) + Modbus TCP (`goburrow/modbus`)
- V1 protocols: Siemens S7 (`gos7`), SQL, REST, Files/FTP
- Local SLM: Phi-3 / Ollama (tag classification, fully local)
- Local DB: SQLite (7-15 day ring buffer, TTL auto-purge)
- Pipeline: Redpanda Connect (declarative YAML, no-code)
- Cloud: Scaleway FR (minimal — KG aggregation + remote dashboard)
- Frontend: React (local first, localhost:3000) + Next.js (website)
- License: Apache 2.0

## Dev environment
- OS: Windows PC + PowerShell
- Docker Desktop + Prosys OPC-UA Simulator (local dev without factory)
- Go 1.22+, VS Code, Claude Code in repo terminal

---

## Current stage — [UPDATE EACH SESSION]

Working on: _______________
Last decision: _______________
Current blocker: _______________

---

## What is BUILT and working (June 2026)

### `internal/discovery/opcua.go` ✅
- OPC-UA connection to any server (tested on Prosys simulator)
- Recursive node tree browse with Continuation Point handling
- Tag discovery: reads NodeID, Name, DataType (`AttributeIDDataType`),
  Value (`AttributeIDValue`) — single ReadRequest per tag
- Noise filter: skips Server, Types, Views, Aliases, StaticData, MyObjects
- `mapDataType()` + `opcuaTypeNodeIDToString()` → clean type names
  (Boolean / Float / Double / Int32 / Int64 / String / DateTime…)
- Live subscription via MonitoredItems, 500ms interval, value change callback
- `WatchForChanges()`: polls node tree every 20s, diffs added/removed tags,
  calls `onChange` only when something actually changed (zero overhead otherwise)
- `browseNodeSilent()`: silent browse for periodic watcher (no logs)
- `tagsToMap()`: O(n) diff keyed by NodeID (stable, survives SCADA renames)

### Key technical decisions locked in
- `ReferenceTypeID`: 0:33 (HierarchicalReferences) — not 0:31
- `RequestedMaxReferencesPerNode`: 100 — safe for all OPC-UA servers
- Continuation Points: always released via `BrowseNext(ReleaseContinuationPoints: true)`
- Browse sleep: 50ms between requests — prevents server flooding, avoids session timeout
- NodeID as diff key — stable identifier, survives tag renames

---

## What is NOT built yet (next sessions)

| Module | File target |
|--------|-------------|
| Behavioral inference (10-15 min live pattern matching) | `internal/classifier/behavioral.go` |
| SLM Phi-3 via Ollama (tag classification) | `internal/classifier/slm.go` |
| UNS ISA-95 mapper (tag → Site/Area/WorkCenter/Tag) | `internal/uns/mapper.go` |
| Rules engine core | `internal/rules/engine.go` |
| Micro-stop detection (Run→Stop→Run, 30s < t < 3min) | `internal/rules/microstop.go` |
| Energy waste detection (Modbus TCP) | `internal/rules/energy.go` |
| Causality (tag correlation at stop timestamp) | `internal/rules/causality.go` |
| Fuzzy Join OT/IT (±10 min sliding window) | `internal/fuzzy/join.go` |
| Cost model V0 (manual 3-field: Coût_Horaire, Cadence, Marge_Produit) | `internal/cost/model.go` |
| TRS calculator (real availability, OEE, gain potential) | `internal/cost/trs.go` |
| SQLite ring buffer (7-15 days, TTL auto-purge) | `internal/storage/sqlite.go` |
| HTTPS push to cloud (TLS 1.3 + mTLS, offline queue) | `internal/push/cloud.go` |
| SMTP / Slack alerting from Edge | `internal/alerting/smtp.go` |
| Modbus TCP connector | `internal/discovery/modbus.go` |
| Siemens S7 connector | `internal/discovery/s7.go` |
| React local dashboard (Gantt, Pareto, ROI tabs) | `frontend/src/` |
| Cloud API receiver (Go, events handler) | `api/handlers/events.go` |

---

## POC scope (active target)

Three use cases in one POC:
1. **Micro-stop detection** — OPC-UA, 5 key tags
2. **OT/IT reconciliation** — Fuzzy Join ±10 min, ERP production orders
3. **Energy waste detection** — Modbus TCP energy meters

**Required tags:**
```
Etat_Machine        (Boolean / Int — Run/Stop/Setup/Alarm)
Compteur_Pieces     (cumulative Int)
Vitesse_Instantanee (Float, analog)
Capteur_Bourrage    (Boolean — jam sensor)
Pression_Air        (Float, bar)
Puissance_Electrique (Float, kW — Modbus TCP if not in SCADA)
```

**Cost model V0:** manual 3-field wizard (Coût_Horaire €/h, Cadence unités/h, Marge_Produit €/unité)
**Cost model V1:** automatic import from ERP via SQL/REST connector

**POC success criterion:**
> Docker install in 1 command → automatic OPC-UA connection → micro-stop detection
> → OT/IT reconciliation → energy waste → dashboard with € losses by production order
> — in under 48h on site.

---

## Claude session protocol

```
START OF SESSION
1. Open claude.ai (architecture/decisions) or Claude Code (active coding)
2. Paste context_starter.md
3. Describe session objective

DURING SESSION
- Architecture / decisions  → claude.ai (chat)
- Active code in repo       → Claude Code (terminal in repo folder)
- Debug                     → paste code + error + context

END OF SESSION
- Ask: "Summarize technical decisions made today"
- Update decisions.md
- Update context_starter.md (section "Current stage")
- Commit: git commit -m "session: [topic]"
```
