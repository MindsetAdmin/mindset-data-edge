# Feedback on `docs/Prpopsitions1.md`

> **My analysis + recommendations on each of your propositions.**
> Grounded in the actual code + the 47 locked decisions in `analysis_log.md`.
> Where a proposition conflicts with something locked, I flag it — you decide whether to override.
> Last updated: 2026-07-02

**Filename typo flag**: your file is `Prpopsitions1.md`. If you want I can rename to `Propositions.md` for grep-ability. Non-blocking.

---

## Q1 — "How does the Knowledge Graph pick the data? From DB or from broker?"

### Direct answer: BOTH — one for writing, one for reading. Two different KGs.

There are actually **2 KGs** in the code (see `docs/how_it_works.md` §9):

**Domain KG** (site fingerprint — micro-stops, equipment, causes, costs):
- **Writes come from the MQTT broker.** `internal/kg/subscriber.go` (KGSubscriber) subscribes to `mindset/events/micro-stop` and inserts Equipment / Event / Cause / Cost nodes into SQLite.
- **Reads come from the DB.** `GET /api/kg/domain` queries SQLite tables (`kg_nodes`, `kg_edges`) and returns JSON to the frontend Cytoscape.

**Technical KG** (pipeline topology — Connectors / Functions / Topics / Pipelines / Dashboards):
- **NOT from the broker.** Built entirely in-memory from `internal/kg/builder.go` reading the pipeline registry.
- Cached 5 min, rebuilt when the registry hash changes.
- Reads via `GET /api/kg/technical`.

**In short**:
- Real-time enrichment happens over MQTT (Domain KG).
- Persistence + queries happen against SQLite (Domain KG).
- Pipeline architecture is computed in-memory (Technical KG).

This matches your intuition — nothing broken here.

---

## Q2 + Q3 — "Every event must publish to MQTT with ISA-95 topic naming"

### My analysis: partially already true, but needs tightening.

**Current state (from `docs/how_it_works.md` §6)**:

| Topic pattern | ISA-95? | Notes |
|---|---|---|
| `mindset/raw/{nodeID}` | ❌ NO — uses OPC-UA nodeID | The rawest layer |
| `mindset/site/{site}/{area}/{work_center}/{tag}` | ✅ YES | This IS ISA-95 |
| `mindset/events/status-change` | ❌ NO — flat, no site context | Global event bus |
| `mindset/events/micro-stop` | ❌ NO — flat, no site context | Global event bus |
| `mindset/dashboard/{label}` | ❌ NO — user-defined labels | Pin-based |

So we already have ISA-95 topics for the contextualized layer, but events and raw data don't follow ISA-95.

### Your proposal, refined

I think you're proposing: **every event carries its ISA-95 location** so any subscriber (KG, dashboard, AI agent) knows where in the site hierarchy it happened, without needing to look up metadata separately.

There are 2 ways to enforce this:

**Option A — Nest events under site paths**:
```
mindset/site/{site}/{area}/{work_center}/events/status-change
mindset/site/{site}/{area}/{work_center}/events/micro-stop
```
Pros: fully ISA-95, subscribers can filter by site/area easily
Cons: no more single-subscription "give me all micro-stops"; harder for KG to aggregate

**Option B — Keep flat event topics, ENFORCE ISA-95 fields in payload**:
```
Topic: mindset/events/micro-stop
Payload: {
  "site": "...", "area": "...", "work_center": "...", "work_unit": "...",
  "tag_name": "...", "uns_topic": "mindset/site/.../etat_machine",
  "timestamp_ms": ..., "duration_seconds": 47, "cost_eur": 18, ...
}
```
Pros: keeps single-topic subscription simple; still full ISA-95 in the payload; matches how the code already partially works (see `internal/rules/engine.go` where `EnrichedMessage` has all the ISA-95 metadata fields)
Cons: less MQTT-purist, but more practical

### My recommendation: **Option B (or hybrid A+B)**

Enforce that every event payload includes the complete ISA-95 metadata block, following the existing `EnrichedMessage.Metadata` struct pattern. Publish to the flat topic for global subscribers (KG, dashboard) AND republish to the nested ISA-95 topic for site-specific subscribers (a hybrid).

**Code impact**: minor — augment publishers in `internal/rules/engine.go` and pipelines that publish to `mindset/events/*` to include the full ISA-95 block. `internal/uns/mapper.go` already does the mapping — reuse it.

**Update to lock**: this becomes a new architectural decision in `decisions.md` — *"All events on MQTT MUST carry the full ISA-95 metadata block."*

---

## Q4 — "I want the KG to be LIVE (real-time) — reflect what's in the pipeline"

### Direct answer: doable, with 2 sub-questions.

**Sub-question 4a — Live Technical KG (pipeline topology)**

Currently: cached 5 min, rebuilt on registry hash change.
Your proposal: real-time updates.

**How to do it**:
1. Invalidate the cache immediately every time a pipeline is registered / updated / deleted (already happens by hash, but polling)
2. Add a WebSocket message type `{type: "kg-technical-update"}` to `internal/kg/builder.go` — push whenever the graph rebuilds
3. Frontend `KnowledgeGraphPage` subscribes via `useLiveSocket` and re-renders

**Effort**: 2-3 days.

**Sub-question 4b — Live Domain KG (site fingerprint)**

Currently: enriched via MQTT events, frontend re-fetches on user action.
Your proposal: push updates to frontend as they arrive.

**How to do it**:
1. `LiveHub` in `cmd/server/live.go` already sees all `mindset/#` traffic. When it sees a `mindset/events/*` message, push `{type: "kg-domain-update", data: {node: ..., edges: [...]}}` to the WebSocket
2. Frontend Cytoscape adds/updates the node without full reload

**Effort**: 2-3 days.

### My recommendation: **do both**

Both are consistent with the "AI-native + real-time" positioning. Total effort ~1 week. Great demo item for investors — "watch the KG grow live as the factory runs."

---

## Q5 — ⚠️ MAJOR STRATEGIC FLAG — "Clean KG + pipelines on program restart from frontend"

### This CONFLICTS with a locked moat. Full analysis below.

Your proposal: when the user clicks "Stop" from the frontend, wipe the KG + pipelines. Fresh session on next start.

**Current state**: KG persists in `data/mindset.db` — survives restarts. Pipelines are YAML files in `config/pipelines/` — also persist.

### The conflict

**Moat #3 (site fingerprint) — locked in `decisions.md`**:
> *"Cumulative site fingerprint — every micro-stop, every cause, every cost calibration accumulates over months. Replacing MindSet = losing all accumulated intelligence. Churn becomes structurally prohibitive after month 6."*

If the KG wipes on every restart, **the moat evaporates**. Six months of accumulated stop causes, calibrated costs, tribal knowledge — all gone. Customer replacement cost becomes low.

### What I think you actually mean (challenge)

There are 3 possible interpretations:

| Interpretation | My read | Recommendation |
|---|---|---|
| **A — DEV workflow**: during development / demo prep, you want to wipe state to start clean | Likely what you want in practice | Add a **dev-mode toggle** in the UI (or a CLI flag `--reset-state`) that clears KG + config/pipelines. Not exposed to end customers. |
| **B — Session-scoped for MULTIPLE CUSTOMERS on the same install**: each customer runs a session, nothing persists between | Unlikely — each customer install is separate | Reject — doesn't match the ETI deployment model (one install per site) |
| **C — Production behavior**: end customers get clean state every restart | Destroys the moat | **STRONGLY REJECT** — this contradicts Moat #3 |

### My recommendation

Add a **dev-mode "Reset session" button** in the UI (probably in Overview or a hidden Settings page). When clicked:
- Wipes `data/mindset.db` (KG + tags + events)
- Deletes user pipelines in `config/pipelines/` (KEEPS the `examples/` templates)
- Restarts state trackers

This satisfies your workflow need WITHOUT destroying the moat. Production customers never touch this — it's a founder / dev / demo tool.

**Decision needed from you**: A, B, or C? If C, be explicit — we'll need to update `decisions.md` and the moat framing across all docs (mindset.md §15, competitive Excel Sheet 3, memo_cecilia_FR.md).

---

## Q6 — "Dashboard shows only machines the user selected (ISA-95 normalized) in Connect"

### Direct answer: sensible UX, small implementation.

**Current state**: Dashboard (`/dashboards`) shows all pinned widgets (from `add_to_dashboard` function outputs) plus all discovered machines from `/api/machines`.

**Your proposal**: filter dashboard to only show machines the user explicitly selected + mapped to ISA-95 (site/area/work_center/work_unit) in the OPC-UA Connect page.

### Why this is good

- Reduces noise on the dashboard — Plant Manager sees THEIR configured machines, not all raw tags
- Forces the ISA-95 mapping discipline — no shortcut around configuring the hierarchy
- Aligns with the "guided setup / onboarding wizard" (U10 in the V1 inventory)

### How to implement

1. When user selects OPC-UA tags in `/connect/opcua`, store the ISA-95 mapping per selected tag (already exists — `/api/opcua/selections` returns per-tag routing)
2. Add `machine_id` field (derived from site/area/work_center) to the selection metadata
3. `DashboardPage` calls `/api/machines?configured=true` (new filter) → only returns machines with completed ISA-95 mapping
4. Unmapped tags: hidden by default; toggle "Show raw / unmapped" if user wants to see them

**Effort**: 3-4 days (mostly UX work — the routing infra is already there).

### My recommendation: **do it, ship in the frontend redesign Phase 3**

Fits naturally with the Grafana-style Dashboard redesign we started. Include as a Priority 1 / DashboardPage task.

---

## Q7 — "Use KG/UNS as ONE trusted source — AI agents get everything from there"

### Direct answer: this is already the architectural principle. Let's make it explicit.

**What the current code does**:
- UNS lives in MQTT (`mindset/site/#` topics) — real-time
- Domain KG persists in SQLite — historical + aggregate
- AI agents (V1: Ad-hoc Analyst) will access both via **MCP tools** — see `impact_engine.md` and `advisors.md`

**Your proposal**: enforce KG/UNS as the SINGLE source. AI agents can't read raw MQTT / raw SQL / raw OPC-UA — they only see the KG.

### Why this is right

- Architecturally clean — one interface, one contract, one auditable path
- Aligns perfectly with the "Impact Engine + MCP" plan (Entries 40-41)
- Prevents AI agents from bypassing the reconciliation layer (which is the moat)
- Matches Cognite's approach with their Industrial Knowledge Graph — but at the edge

### What to do

**Lock this as a new architectural decision** in `decisions.md`:

> *"AI agents access data only via the KG/UNS layer exposed by the MCP server. AI agents MUST NOT read raw MQTT / raw SQLite / raw connector data directly. This ensures every AI query benefits from reconciliation (Fuzzy Join), context enrichment (Impact Engine), and audit logging."*

**MCP tool naming enforces this**:
- `kg_query(node_type, filter)` — read KG
- `kg_list_events(since, machine)` — read events
- `kg_cost_summary(range, group_by)` — read Impact Engine outputs
- `uns_get_current(topic_pattern)` — read UNS current state

NO tool named `sqlite_query` or `mqtt_subscribe`. Direct raw access is architecturally forbidden.

### My recommendation: **lock this decision, add to `decisions.md`, update `mindset.md` §8 Module 9 (MCP server)**

This is a strong principle. Investors will love this — clean, defensible architecture. Adds nothing to build time (MCP server was already scoped this way).

---

## Q8 — "Can `calculate_duration` be used with other functions, or only next to `state_machine`?"

### Direct answer: yes, it's already generic. Not tied to `state_machine`.

I read `internal/functions/calculates/duration.go`:

- The function takes any `event_id` string + optional `end_time`
- On FIRST call for an event_id, stores the start time
- On SECOND call for the same event_id, computes the duration and cleans up

**No dependency on `state_machine`** — the handler is agnostic. Any pipeline node that produces `event_id` and `end_time` params can trigger it.

### Practical implication

You can use `calculate_duration` for:
- Time between two Run→Stop transitions (current use)
- Time between an alert triggering and resolution
- Time between an OF starting and finishing (in the Fuzzy Join context)
- Cycle time between two counter increments
- Anything with a defined start + end event

### Caveats (worth knowing)

1. It stores start times **in memory** — if the process restarts mid-cycle, the start time is lost
2. It's keyed by `event_id` — if two pipelines use the same event_id concurrently, they'll conflict
3. **Not thread-safe** for concurrent updates to the same event_id (Go map)

For production robustness (V1.5+): persist start times to SQLite and add mutex. Not urgent.

### My recommendation: **document the generic usage in `functionDocs.js`**

Currently the frontend might imply this function is only for state_machine outputs. Update the docstring to say "usable anywhere you have paired start/end events sharing an event_id."

---

## Q9 — "Delete `kg_save`, every pipeline added will be shown directly in the KG"

### Direct answer: nuanced — depends on WHICH KG you mean.

There are 2 KGs (per Q1 answer). `kg_save` affects the Domain KG (site fingerprint).

**What already auto-populates without kg_save:**

- **Technical KG** (pipeline topology): fully automatic. `internal/kg/builder.go` reads the pipeline registry — as soon as a pipeline is saved, it appears in the Technical KG at next cache refresh. **You don't need `kg_save` for this**.
- **Domain KG** (events): automatic via `KGSubscriber`. Any `mindset/events/micro-stop` message auto-inserts Equipment / Event / Cause / Cost nodes. **You don't need `kg_save` for the standard event flow either**.

**What `kg_save` currently uniquely does:**

Looking at `internal/functions/outputs/kg_save.go` — it's an EXPLICIT output function for CUSTOM enrichment. When a pipeline author wants to write something to the Domain KG that isn't a standard event (e.g., a machine metadata update, an operator note, a tribal knowledge entry that doesn't come through the auto-enrichment path).

### My analysis

**If you delete `kg_save`**:
- ✅ Technical KG still auto-updates — no impact
- ✅ Standard event enrichment still works — KGSubscriber handles it
- ❌ You lose the escape hatch for custom pipelines that want to write non-event data to the Domain KG

**The escape hatch matters for**: tribal knowledge capture (Moat #4 — dropdown + free text on stops); machine metadata mapping from the Connect page; custom V1.5+ agents wanting to attach findings to the KG.

### My recommendation: **keep it, but make its purpose clearer**

Two options:

**Option 1 (minimum change)**: rename `kg_save` → `kg_write_advanced`. Mark in the function catalog as "Advanced — usually not needed. The Technical KG auto-populates from pipeline definitions; event enrichment auto-populates from `mindset/events/*` topics. Use this only for custom writes."

**Option 2 (bigger cleanup)**: split into 3 named functions with clear semantics:
- `kg_write_event` — for custom event nodes (what most people mean by kg_save today)
- `kg_write_metadata` — for machine/tag/operator metadata
- `kg_write_note` — for tribal knowledge notes

**Effort**: 1 day for option 1; 2-3 days for option 2.

### If you INSIST on deleting it

You can — everything critical still works. Just document clearly that:
- Technical KG auto-populates on pipeline save (already true)
- Domain KG auto-populates from events on the standard MQTT topics
- Custom KG writes are no longer possible without extending the codebase

**My recommendation**: Option 1 (rename + document better). Removes user confusion without losing capability.

---

## Summary — what I'd change based on your propositions

| # | Your proposition | My recommendation | Priority | Effort |
|---|---|---|---|---|
| 1 | How does KG pick data | Q — answered | — | — |
| 2+3 | Everything MQTT with ISA-95 topics | **DO** — enforce ISA-95 payload on all events (Option B). Lock as decision. | HIGH | ~2 days |
| 4 | Live KG | **DO both** (Technical + Domain live updates via WebSocket) | HIGH | ~1 week |
| 5 | Clean KG + pipelines on restart | ⚠️ **STRATEGIC** — this destroys Moat #3. My proposal: add **dev-mode reset button** instead. Waiting on your call: A / B / C. | BLOCKING DECISION | ~1 day for reset button |
| 6 | Dashboard = only ISA-95 configured machines | **DO** — as part of frontend redesign Phase 3 | MEDIUM | 3-4 days |
| 7 | KG/UNS as single source for AI agents | **LOCK as decision** — matches existing MCP plan | HIGH (docs) | ~2h to write decision |
| 8 | Is calculate_duration reusable | Q — yes, already generic. Update docs only | LOW | 30 min |
| 9 | Delete kg_save | **KEEP but rename** to `kg_write_advanced` with clearer docs | LOW | ~1 day |

**Total impact**: 1 blocking decision (Q5) + ~2 weeks of engineering work for the "do all" path. Non-trivial but delivers a much cleaner architecture + fixes real UX issues.

---

## Questions back to you

1. **Q5 interpretation**: A (dev-mode reset button) / B (multi-customer session model) / C (production behavior — destroys moat)?
2. **Q2 nested vs flat event topics**: pure ISA-95 nesting (Option A) or ISA-95 payload block on flat topics (Option B — my pick)?
3. **Q9 kg_save**: rename to `kg_write_advanced` (my pick) / split into 3 / delete entirely?
4. **Priority order**: given ~2 weeks of engineering time, which do you want first — Live KG (Q4)? Event ISA-95 tightening (Q2)? Dashboard filter (Q6)? Frontend redesign continuation?

Once you answer these, I can turn any of them into concrete code changes + logged decisions.
