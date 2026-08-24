# Mindset Data — Context Brief

> **For a Claude Code session in this repo**: don't paste this file — `CLAUDE.md` is auto-loaded automatically and is the accurate, current technical source of truth (data flow, packages, API surface, known limitations). This file is for **claude.ai / web-chat strategy sessions** where `CLAUDE.md` isn't available, or as a fast human-readable orientation doc.
>
> **Rewritten 2026-07-28** — the previous version of this file described a stack that never matched what got built (Redpanda Connect, Phi-3/Ollama, Apache 2.0 license, a "not built yet" table where nearly everything listed was actually already built). See `docs/analysis_log.md` Entries 137-138 for how that drift was found. This version reflects what's actually in the repo today; the full always-current version of most of this lives in `CLAUDE.md`.

## What we're building

Industrial data infrastructure for manufacturing ETI/mid-market factories.
Two Go binaries (`cmd/server`, `cmd/agent`) + a React pipeline-builder frontend → ISA-95-aligned Knowledge Graph → real-time dashboard + MCP access for AI agents.
Currently a **single repo, single-machine deployment** (server + agent + frontend on one box, talking over local MQTT) — not yet the multi-edition Cloud/Hybrid/On-Premise distribution model described in `docs/mindset.md`'s vision sections; that's roadmap, not current state.

## Team

- Mohamed (CTO) — Polytechnique, IoT & embedded systems (Windows PC)
- Cécilia (CEO) — EDHEC, ex-VC AgriFoodTech

## Stack (actually in `go.mod` / `package.json` — verify there if in doubt)

- Backend: Go, two binaries (`cmd/server` — API + WebSocket + MCP; `cmd/agent` — rules engine + KG subscriber)
- OPC-UA: `github.com/gopcua/opcua`
- MQTT: `github.com/eclipse/paho.mqtt.golang` (broker required, e.g. Mosquitto — not bundled by default)
- SQL connector (V1a, MySQL only): `github.com/go-sql-driver/mysql`
- WebSocket: `github.com/gorilla/websocket`
- MCP server: `github.com/modelcontextprotocol/go-sdk` — HTTP (`/mcp`) and stdio (`-mcp-stdio`) transports, 5 read-only tools
- Local DB: SQLite, pure-Go driver `modernc.org/sqlite` (no CGO) — one file, `data/mindset.db`, shared by both binaries with `PRAGMA busy_timeout`
- Pipeline engine: **custom Go implementation** (`internal/pipeline`) — YAML-defined nodes, topological execution. **Not Redpanda Connect** — that was vision-doc content that was never built.
- Local SLM (Phi-3/Ollama) / behavioral inference: **not built.** No local LLM runtime exists in the repo yet.
- Frontend: React 19 + Vite + Tailwind + Zustand, ReactFlow (pipeline canvas), Cytoscape (KG viewer), i18n via `react-i18next` (FR default, EN toggle)
- Dev-only fake ERP: `cmd/erpsim` + `sim/erp/*.sql` — simulates a customer MySQL ERP for testing the SQL connector; explicitly a dev fixture, not a real integration
- Dev-only OPC-UA source: Prosys OPC-UA Simulation Server (free, local) — not a real factory
- **License**: **Proprietary (closed-source), 2-year minimum** — per `docs/decisions.md` (restored 2026-07-28 from git history): "supersedes prior Apache 2.0 decision." No real `LICENSE` file exists in the repo yet — that's a genuine open gap, independent of which model was chosen.

## Dev environment

- OS: Windows PC + PowerShell (primary shell; Bash also available)
- `run.ps1` — builds both binaries, starts server + agent + frontend, opens browser (`-NoAgent` / `-NoBuild` flags available)
- Docker Desktop + Prosys OPC-UA Simulator (local dev without a factory)
- Go 1.26+, Node/npm for the frontend, Claude Code in repo terminal
- `docker-compose.dev.yml` — mosquitto + fake-ERP MySQL (+ optional erpsim container) for local dev

---

## Current stage — 2026-07-28

Working on: outbound prospecting (lemlist) for US/UK/France/Benelux industrial IT contacts, drafting personalized outreach emails, and a docs-integrity pass (found + partially fixed stale/contradictory content in `mindset.md`/`context_starter.md`, found unsourced claims in the pitch copy, deleted hand-tuned demo KG data).
Last decision: KG's demo Event/Cost data (40+40 nodes, hand-tuned per `analysis_log.md` Entries 133/135) was deleted — the KG currently holds only real structural data (Equipment/Site/Area/WorkCenter/Tag/SchemaMapping from the OPC-UA/ERP simulators), no micro-stop/cost numbers, until reseeded or fed real data.
Current blocker: none outstanding from this pass — `docs/decisions.md` (restored 2026-07-28) is untracked in git; commit it when you're ready to make the restore permanent.

---

## What is actually BUILT and working (per `CLAUDE.md`, verify there for full detail)

- OPC-UA connect/browse/subscribe, live tag ingestion, ISA-95 auto-mapping with confidence scoring (`internal/uns`)
- MQTT-based decoupling between `cmd/server` and `cmd/agent` (raw → contextualized → event topics)
- Rules engine: Run/Stop state tracking, micro-stop detection
- Cost calculation (manual/config/tag rate sources, CSV/Excel per-product tables)
- SQLite-backed unified Knowledge Graph (`business` + `platform` categories), with structural auto-bootstrap from both OPC-UA (OT side) and SQL schema discovery (IT side, `SchemaMapping` nodes), confidence-gated pending/validate/reject workflow
- Entity resolution: OT `Equipment` ↔ IT `work_center` matching (`same_as` edges)
- Active-production tracking from a validated ERP `work_order` mapping, merged with cost-priority ranking (urgency vs. cost as two distinct, sometimes-disagreeing axes)
- MCP server, 5 tools, both HTTP (`/mcp`) and stdio transports (Claude Desktop-compatible)
- Pipeline Studio (drag-and-drop, ReactFlow), KG viewer (Cytoscape) with pending-validation UI, real-time dashboard (WebSocket-driven)
- i18n (FR default / EN toggle)
- `cmd/erpsim` fake-ERP dev stack for exercising the SQL connector end-to-end

## What is NOT built yet

- Local LLM / SLM (Phi-3, Ollama) — no local inference runtime in the repo
- Behavioral inference / tag classification beyond the ISA-95 depth/collision heuristics already in `internal/uns`
- Redpanda Connect or any pipeline engine other than the custom Go one
- Modbus and Siemens S7 connectors — `modbus_read` is a metadata-only stub that errors if executed; S7 isn't started
- Any cloud tier (Hybrid/Self-Hosted editions, cross-site KG aggregation, remote dashboard, license-key distribution) — everything today runs on one machine
- Automatic pipeline triggering from live MQTT — pipelines only run via manual/API-triggered `Run`, see `CLAUDE.md`'s "Known limitations"
- Retroactive/historical product-scoped queries ("how long did product X run yesterday") — only live/current-moment queries are answered today

---

## Session protocol

```
START OF SESSION
- Claude Code in this repo: CLAUDE.md loads automatically, no need to paste this file.
- claude.ai / web chat (architecture, strategy, no repo access): paste this file.

DURING SESSION
- Architecture / decisions   → claude.ai (chat) or ask Claude Code to read docs/analysis_log.md
- Active code in repo        → Claude Code (terminal in repo folder)
- Debug                      → paste code + error + context

END OF SESSION (if the session produced a real decision or state change)
- Log it in docs/analysis_log.md (project convention — see MEMORY.md "Analysis Log Convention")
- Update this file's "Current stage" section
- Commit if asked to — never commit automatically
```
