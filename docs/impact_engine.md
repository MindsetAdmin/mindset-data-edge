# MindSet — Cost Function / Impact Engine

> **Product spec.** What the cost function should become — going beyond V0's `duration × hourly_cost` toward a financially-grade prioritization engine that leverages OT/IT reconciliation (ERP + MES + others).
> **Status**: V0 shipped in `internal/functions/calculates/cost.go` · V1 enrichments scoped (this doc) · V1.5+ phased.
> **Related**: `docs/analysis_log.md` Entries 40 + 41 for the strategic reasoning · Entries 69-71 for the reconciliation with `docs/Cost_function.md` folded into this doc below · `docs/Cost_function.md` for the original provider-pattern brainstorm (superseded by §"Providers" below, kept for history) · `docs/mindset.md` §8 Module 5 (placeholder for the upgraded module) · `docs/decisions.md` for the locked architecture decision.
> **Last updated**: 2026-07-14 — reconciled with `docs/Cost_function.md` (Entries 69-71): resolved flag-vs-priced-number split, dropped the standalone Fuzzy Join engine in favor of provider nodes the pipeline engine already supports, added the Stock-shortage dimension.

---

## Why "cost function" undersells what this is

The current `internal/functions/calculates/cost.go` computes `TotalImpact = duration_minutes × (HourlyRate / 60)`, with `HourlyRate` swappable via a per-product CSV/Excel table (corrected 2026-07-14 — see the V0 baseline section below for what actually shipped vs. what this doc originally assumed). That's a stub.

A real Impact Engine does three things, not one:
1. **PRICE** every event in € — using real product/customer/recipe context, not constants
2. **PRIORITIZE** by actionable impact (not just absolute cost) — so the Plant Manager fixes the right thing first
3. **PROJECT** the consequence of action — *"if you fix cause #1, you recover Y €/week, conditional on cause #2 not becoming next bottleneck"*

PRICE is what V0 attempts (badly). PRIORITIZE is the decision-time-reduction the user asked for. PROJECT is the CFO-grade differentiator we ship in V2.

**Naming convention going forward**: "Impact Engine" in pitch + docs. "Cost function" still acceptable as informal shorthand internally.

---

## The 14 dimensions a real Impact Engine should cover

Categorized by what they require and when they ship.

### The pricing rule (locked 2026-07-14, Entry 71)

Every dimension below is either **priced with a visible formula** or **flagged** — never a bare asserted number with no traceable derivation. This resolves the tension between this doc's original caution (customer-commitment as flag-only) and `docs/Cost_function.md`'s brainstorm (which priced delivery risk at a flat, unauditable "+1200€"):

- **Contractual / opaque** (penalty clauses, SLA credits) → flag only. Most ERPs don't expose penalty-clause terms in a queryable field — pricing it would produce a number the CFO can't audit. Stays V3 per the original plan (dimension 2, and the new Customer Penalty in V3 below).
- **Computable from ERP/MES fields customers DO expose** (stock levels, reorder points, lead times, product margins) → price it, but derive it from a stated formula, never hardcode it. This is what unlocks the new Stock-shortage dimension (14) below.
- **Judgment-call proxies** (a 1-3 "criticality" rating with no underlying margin data) → allowed only as a LOW-confidence fallback when the real figure (per-product margin) isn't available — never the primary mechanism.

### Category A — Reconciliation-powered (the moat-leveraging ones)

These ONLY work because MindSet reconciles OT events to active production context from ERP / MES / others. This is what "core value is reconciliation" means in practice.

| # | Dimension | Primary IT source | Ships |
|---|---|---|---|
| 1 | Per-product margin (not generic line margin) | ERP product master OR MES recipe | **V1** |
| 2 | Customer-commitment flag (delivery risk indicator, FLAG only — see pricing rule) | ERP — OF→customer→due-date | **V1** |
| 3 | Downstream idle propagation (Theory of Constraints) | Line-layout config (V1) → MES auto (V1.5) | **V1** (manual) / V1.5 (auto) |
| 4 | Setup/restart cost (re-warming, re-priming, cleaning) | Config V1 → MES setup-time matrix V1.5 | **V1** (config) / V1.5 (auto) |
| 5 | Quality scrap on restart | MES historical scrap rate per stop type | V1.5 |
| 6 | Energy peak penalty (restart triggers tariff window) | Energy Management System + grid contract | V1.5 |
| 7 | Material spoilage (cold chain, fermentation, alloy stability) | MES batch state + recipe time-limits | V2 |
| 8 | Regulatory cost (pharma: compliance reporting + revalidation) | MES + QMS regulatory flags | V2 |
| 9 | Schedule criticality (buffer vs critical-path order) | ERP MRP / planning data | V2 |
| 14 | **Stock-shortage risk (PRICED, computed — not the contractual penalty)** | ERP stock level + reorder point + lead time | **V1** |

### Category B — Aggregation / decision-time reducers (the "prioritize impacts" ones)

These compute ON TOP of per-event impacts to drive action.

| # | Dimension | Primary input | Ships |
|---|---|---|---|
| 10 | Bottleneck identification | Cumulative KG observation (30-day window) | V1.5 |
| 11 | Actionability score = impact × ease-of-fix | Impact (computed) + operator-rated ease (V1 tribal knowledge dropdown) | V1.5 |
| 12 | Forward simulation ("if you fix cause #1, recover Y €/week conditional on...") | Computed from causes + bottleneck model | V2 |
| 13 | Time-weighted urgency (this week vs cumulative) | Computed from event log | V1 (simple) / V1.5 (richer) |

---

## V0 baseline — what's already coded

**Correction (2026-07-14, Entry 69/71): this section previously cited a nonexistent `internal/cost/model.go` with a 3-component formula. The real V0 code is `internal/functions/calculates/cost.go`, and it's simpler than that — one component, not three:**

```
TotalImpact_V0 = duration_minutes × (HourlyRate / 60)
```

`HourlyRate` resolves in this order: default config → per-call `config.hourly_rate` override → a **per-product lookup** in a CSV/Excel-loaded rate table, if the event carries a `product` field that matches a row.

That per-product lookup is worth noting explicitly: **it's already a crude precursor of Enrichment #1 (per-product margin)** — the codebase isn't starting from zero on that dimension. The gap, fixed as part of V1 below: it substitutes into `HourlyRate` (a TimeLoss proxy) instead of being its own `ProductionLoss` term, which conflates two economically different things. V1 splits them apart.

**Issue with V0**: no ProductionLoss or EnergyLoss component at all yet, and the one dimension that exists (product-specific rate) is wired into the wrong term. Plant Manager + CFO will sanity-check and find the math too simple to trust.

---

## V1 Impact Engine — concrete additions (~3-4 weeks of work on top of current Track 1)

Keep V0 as the floor. Add 4 enrichments (plus Enrichment #14, folded in below) that leverage `work_order`/`product`/`stock`-tagged provider nodes + ERP/MES reads.

### Enrichment #1 — Per-product margin (replaces generic)

**Math:**
```
ProductMargin = lookup(active_OF.product_id) from ERP product master OR MES recipe table
fallback = LineProductMargin (V0 generic)
```

**Data source**: ERP `products` table OR MES `recipes` table — first one available wins.

**Code**: extend `internal/functions/calculates/cost.go` to read the `active_of`/`product` context from a `canonical_type`-tagged provider node in `params` (not a Fuzzy Join engine — see Architecture below); add `internal/functions/calculates/product_margin.go` for the lookup.

**Configuration** (`config/cost.yaml`):
```yaml
product_margin_source: erp        # or "mes"
product_margin_table: products    # or path in mes schema
product_margin_field: margin_eur_per_unit
fallback_margin: 0.08             # used if lookup fails
```

**Output change**: `ProductionLoss` now uses per-product margin. Number changes dynamically per active OF.

**Confidence tiering (2026-07-14, Entry 71)**: this ERP/MES lookup is the HIGH-confidence path. When it fails or isn't configured, fall back to V0's existing CSV/Excel per-product rate table (already shipped — see the V0 baseline correction above) as a LOW-confidence proxy, rather than a flat generic constant. Never invent a bare 1-3 "criticality" multiplier with no underlying margin data as the primary mechanism — that's the judgment-call proxy the pricing rule above demotes to fallback-only.

---

### Enrichment #2 — Customer-commitment flag

**Math:**
```
IsCustomerCommitted = active_OF.customer_id != null
                      AND active_OF.due_date - now < 7 days
                      AND active_OF.status IN ("In Progress", "Released")
DisplayFlag = "⚠️ HIGH PRIORITY — customer-committed delivery" if true
```

**Important V1 limit**: we DO NOT add a penalty cost number to the impact. Penalty clauses live in contracts that most ERPs don't expose. We FLAG the event so it surfaces in dashboard priorities; quantification waits for V3.

**Data source**: ERP `production_orders` (customer_id + due_date fields). MES also has this in many setups.

**Code**: `internal/functions/calculates/customer_flag.go` reads active OF context, applies the rule, attaches flag to event.

**Configuration** (`config/cost.yaml`):
```yaml
customer_commitment:
  enabled: true
  due_date_window_days: 7
  source: erp                      # or "mes"
  customer_id_field: customer_id
  due_date_field: requested_delivery_date
```

**Output change**: events emit a `customer_committed: bool` field. Dashboard Pareto shows the flag visually + prioritizes flagged events in the "Top 3 actions this week" widget.

---

### Enrichment #14 — Stock-shortage risk (NEW 2026-07-14, Entry 71 — reconciled from `docs/Cost_function.md`)

**Origin**: `docs/Cost_function.md`'s brainstorm wanted a priced "delivery risk" (its example: stock=0 → flat +1200€). That number had no visible derivation, which fails this doc's trust principles. This enrichment keeps the instinct — stock risk IS priceable, unlike the contractual penalty in Enrichment #2 — but derives it from a stated formula instead of asserting it.

**Math:**
```
StockDeficitUnits    = max(0, DemandDuringLeadTime - CurrentStock)
DemandDuringLeadTime  = AvgDailyDemand × ReplenishmentLeadTimeDays
ExpectedShortageCost  = StockDeficitUnits × ProductMarginPerUnit
```

If `ReplenishmentLeadTimeDays` or `AvgDailyDemand` aren't available from the ERP, **do not fall back to a guessed number** — omit the field entirely (`stock_risk_eur: null`, not `0`) and mark that event's confidence no higher than MEDIUM. A missing value must be visibly missing, not silently zero.

**Data source**: ERP `stock_levels` (current stock, reorder point) + `products` (avg daily demand, lead time) — same tables Enrichment #1's product-margin lookup already reads, so this rides on the same provider connection.

**Code**: `internal/functions/calculates/stock_risk.go`.

**Configuration** (`config/cost.yaml`):
```yaml
stock_risk:
  enabled: true
  source: erp
  stock_table: stock_levels
  lead_time_field: replenishment_lead_time_days
  demand_field: avg_daily_demand
```

**Output change**: events emit `stock_risk_eur: float | null`. Never rendered with false precision — the dashboard breakdown must show the formula inputs (deficit units × margin), matching trust principle #1.

---

### Enrichment #3 — Simple downstream idle propagation

**Math** (simplified Theory of Constraints):
```
DownstreamIdleCost = if (line_layout configured for affected machine):
    nb_downstream_machines × LineHourlyCost / nb_machines_on_line × duration
else:
    0   # silent fallback, no error
```

**Why simplified**: full TOC needs bottleneck model (V1.5). V1 just says "if the stopped machine has downstream machines, they're also idle — cost them".

**Data source**: customer-supplied `config/line_layout.yaml`:
```yaml
lines:
  Line 1:
    machines:
      - { id: M1, upstream: [], downstream: [M2] }
      - { id: M2, upstream: [M1], downstream: [M3] }
      - { id: M3, upstream: [M2], downstream: [] }
    hourly_cost_total: 85           # whole-line hourly cost
```

**Code**: `internal/functions/calculates/line_layout.go` loads the YAML at startup, exposes lookup; `internal/functions/calculates/cost.go` adds DownstreamIdleCost component.

**Output change**: events emit a `downstream_idle_cost: float` component. Total impact now includes propagation when layout is configured.

---

### Enrichment #4 — Setup/restart cost (manual config V1)

**Math:**
```
RestartCost = lookup(stop_type, product_id, machine_id) from config/setup_costs.yaml
fallback = 0
```

**Configuration** (`config/setup_costs.yaml`):
```yaml
restart_costs:
  - { machine: M1, product_type: "Chocolate", stop_type: "Jam", cost_eur: 25 }
  - { machine: M2, product_type: "*",         stop_type: "Air Pressure", cost_eur: 5 }
  - { machine: "*", product_type: "*",        stop_type: "Series Change", cost_eur: 80 }
```

**V1.5 upgrade**: replace manual config with auto-computed from MES historical setup-time matrix × labor cost.

**Code**: `internal/functions/calculates/restart_cost.go`.

**Output change**: events emit `restart_cost: float`. Visible in dashboard breakdown.

---

### V1 Total Impact formula

```
TotalImpact_V1 = TimeLoss                          # from V0
                + ProductionLoss × per_product_margin   # Enrichment #1
                + EnergyLoss                       # from V0
                + DownstreamIdleCost               # Enrichment #3 (if line_layout configured)
                + RestartCost                      # Enrichment #4
                + (StockRiskEur if resolved else 0) # Enrichment #14 — additive when present, never guessed when absent

Event metadata:
  customer_committed: bool                          # Enrichment #2 (flag, not cost)
  stock_risk_eur: float | null                       # Enrichment #14 — null when lead-time/demand data unavailable
  active_of: string                                 # from the work_order provider node
  product_id: string                                # from the work_order provider node
  confidence: HIGH | MEDIUM | LOW                    # see confidence tiering, Trust principle #4
```

**Dashboard impact**: new "Top 3 actions this week" widget driven by:
1. Sort events by `TotalImpact_V1` desc
2. Boost events with `customer_committed: true`
3. Group by `cause` → display top 3 cause-groups with totals

---

## V1.5 Impact Engine (Q3 2027, with MES integration in production)

| # | Enrichment | New formula |
|---|---|---|
| 10 | **Bottleneck identification** | `BottleneckScore[machine] = time_as_constraint_last_30d / total_production_time_last_30d`. Display: *"Line 1 bottleneck: Machine M2 73% of the time"* |
| 11 | **Actionability score** | `ActionabilityScore = AvoidableImpact × EaseOfFix` where AvoidableImpact = projected next-30-day cost if pattern continues, EaseOfFix is operator-rated 1-5 (from tribal knowledge dropdown) |
| 5 | **Quality scrap on restart** | `ScrapCost = QuantityScrappedAfterRestart × MaterialUnitCost`. QuantityScrapped from MES historical defect rate per stop type. |
| 6 | **Energy peak penalty** | `PeakPenalty = if (restart_within_peak_tariff_window): PeakPowerDraw × PeakTariff × restart_duration` |

V1.5 also: **replace V1 manual restart-cost config with MES auto-computed setup-time matrix × labor cost.**

---

## V2 Impact Engine (Q4 2027 — CFO-grade forward simulation)

### Enrichment #12 — Forward simulation (the killer demo)

**Math**:
```
For each top-N cause:
  CurrentImpact[cause]   = sum(impact_last_30d where root_cause = cause)
  AssumedFixEffectiveness = 0.7 (configurable per cause type)
  ProjectedRecovery       = CurrentImpact × AssumedFixEffectiveness
  ConditionalCheck        = whether fixing this cause moves bottleneck to next cause
                            (re-run bottleneck identification with this cause's events removed)
  
Display: "If you fix [Cause #1: Jam on Line 1], you recover ~ X €/week.
         CAVEAT: doing so will make [Cause #2: Air Pressure] the new bottleneck (~ Y €/week recoverable next)."
```

This is the demo that closes CFO and Ops Director. It shows MindSet doesn't just price events — it projects strategic action.

### Other V2 enrichments

| # | Enrichment | Notes |
|---|---|---|
| 7 | **Material spoilage** | If `active_OF.recipe.time_limit_minutes` exists AND `duration > time_limit` → entire batch lost. Especially relevant in agrifood (cold chain) + pharma (fermentation). |
| 8 | **Regulatory cost** | If stop triggers compliance report (MES/QMS rule) → add reporting + revalidation cost. Pharma-specific. |
| 9 | **Schedule criticality** | Weight events on critical-path OFs higher than buffer-time OFs. Needs ERP MRP integration. |

---

## V3 Impact Engine

| # | Enrichment | Notes |
|---|---|---|
| — | **Customer penalty model** | Real penalty cost (not just flag) — requires ERP penalty-clause data extraction, often custom per customer |
| — | **Predictive cost** | "This stop pattern will cause X € in the next 7 days if not addressed" — needs accumulated KG + ML on historical patterns |

---

## Architecture

### Pipeline of an event through the Impact Engine

**Revised 2026-07-14 (Entry 71): there is no separate "Fuzzy Join engine" subsystem.** `internal/pipeline/engine.go`'s `executeNode` already pre-loads every node's `params` with every previous node's raw output keyed by node ID, plus trigger data, plus config (see `internal/pipeline/engine.go:196-211`) — that's the exact substrate active-OF/product/stock/quality context needs. Building `internal/fuzzy/of_state.go` as a bespoke coupling would duplicate what the pipeline engine's dependency-ordered execution already does for free. Instead: **provider nodes** — ordinary `sql_query` connector nodes (once the SQL connector lands, see `docs/mysql_connector.md`) that `calculate_cost` depends on via the pipeline YAML's `depends_on`. `calculate_cost` discovers them by scanning `params` for values carrying a `canonical_type` tag (`work_order`, `product`, `stock`, `quality` — the semantic-mapping layer already spec'd in `docs/mysql_connector.md` §6b) rather than hardcoding provider node IDs, so a pipeline author can name provider nodes however they like.

```
OPC-UA event (state transition)
    ↓ rules engine (detect → micro-stop / energy waste / OEE event)
    ↓ provider nodes (sql_query, canonical_type-tagged: work_order / product / stock / quality)
    ↓ Impact Engine (calculate_cost — discovers providers in `params` by canonical_type)
        ├─ price (V0 floor + V1+ enrichments)
        ├─ flag (customer-commitment, regulatory)
        ├─ score (V1.5+ actionability)
        └─ project (V2+ forward simulation, aggregate level only)
    ↓ KG (persist enriched event)
    ↓ MQTT (publish event with full impact breakdown)
    ↓ Dashboard (Pareto by impact, Top 3 actions widget)
    ↓ AI agent (Ad-hoc Analyst can answer "why this number?")
```

### Inputs per enrichment (table)

| Enrichment | OT input | IT input | Config input |
|---|---|---|---|
| TimeLoss (V0) | event duration | — | LineHourlyCost |
| ProductionLoss (V0) | nb stops | — | Cadence, generic margin |
| EnergyLoss (V0) | OffProdConsumption | — | EnergyPrice |
| Per-product margin (V1) | active_OF.product_id | ERP/MES products table | source config |
| Customer-commitment flag (V1) | active_OF | ERP customer_id + due_date | window_days |
| Downstream idle (V1) | event.machine_id | — | line_layout.yaml |
| Restart cost (V1) | event.cause, machine, product | — | setup_costs.yaml |
| Quality scrap (V1.5) | event.cause | MES defect history | — |
| Energy peak penalty (V1.5) | restart duration | Energy Mgmt System tariff windows | peak threshold |
| Bottleneck ID (V1.5) | machine state history (30d) | optional MES line-state | — |
| Actionability (V1.5) | impact | operator-rated ease (tribal knowledge) | — |
| Forward simulation (V2) | accumulated KG | — | fix-effectiveness defaults |
| Material spoilage (V2) | event duration | MES recipe time-limit + batch value | — |
| Regulatory cost (V2) | event metadata | MES/QMS rules | reporting cost |

### Configuration files this introduces

V1:
- `config/cost.yaml` — sources for product_margin + customer_commitment
- `config/line_layout.yaml` — for downstream propagation
- `config/setup_costs.yaml` — for restart costs (V1 manual)

V1.5+:
- `config/peak_tariffs.yaml` — energy peak windows + tariffs
- `config/scrap_rates.yaml` — optional override of MES-computed values

---

## Trust + transparency principles (non-negotiable)

The Impact Engine MUST satisfy these — otherwise the Plant Manager + CFO sanity-check kills adoption:

1. **Every number is traceable**. Dashboard hovers / clicks reveal: *"This 312€ = Σ over 47 events: TimeLoss 142€ + ProductionLoss 130€ + DownstreamIdle 35€ + RestartCost 5€. Weighted by per-product margin pulled from ERP product master at 2026-07-04 09:32."*
2. **No black boxes**. Plant Manager can audit every component cost.
3. **Customer-tunable weights**. If a customer thinks `EnergyLoss` is over-counted, they can adjust the multiplier in `config/cost.yaml` (with sane defaults).
4. **Confidence levels per enrichment**. If `per_product_margin` lookup failed and we fell back to generic, flag the event with `confidence: medium`. Aggregate dashboard shows % of events at each confidence level.
5. **Versioned cost model**. If we change the formula, old events keep their old impact number (with the version that computed them). Don't silently rewrite history.
6. **Reproducibility**. Given the same inputs + same `config/cost.yaml`, the same number must always come out. No non-determinism.
7. **Phased build = phased trust**. V1 ships 5 enrichments (4 original + Stock-shortage risk) PROVEN against pilot customer data before V1.5 enrichments are added. Don't drown the V1 customer in 14 dimensions on day 1.

---

## Example outputs the Plant Manager actually sees

### Per-event detail

```
Event: Line 1 — Machine M2 stop
Date: 2026-07-04 10:15:32
Duration: 47 seconds
Cause: Jam (Capteur_Bourrage triggered)
Active OF: OF#456 (Product A — Chocolate biscuits)
Customer: GrandeDistribution-X — DUE 2026-07-05 (⚠️ tomorrow)

IMPACT BREAKDOWN:
  Time loss              : 1.10 €  (47s × 85€/h)
  Production loss        : 2.40 €  (1 stop × 3600 u/h × 0.08€/u × 47/3600)
  Downstream idle        : 1.10 €  (M3 idle, same line cost)
  Restart cost           : 5.00 €  (Jam on M2 with Product A — config)
  Energy loss            : 0.00 €
  TOTAL                  : 9.60 €
  Customer-committed     : ⚠️ YES (due tomorrow)
  Confidence             : HIGH (all enrichments resolved)
```

### Weekly Pareto (Top 3 Actions widget)

```
WEEK OF 2026-06-29 — TOP 3 ACTIONS

#1. JAM on Line 1            — 312 € lost / 47 events / 12 customer-committed OFs
    Fix: sensor calibration  — Operator ease 2/5 (easy) → ACTIONABILITY 624
    
#2. SERIES CHANGE excess time — 280 € lost / 8 events / 3 customer-committed OFs
    Fix: SMED workshop       — Operator ease 4/5 (medium) → ACTIONABILITY 70
    
#3. AIR PRESSURE drops       — 145 € lost / 23 events / 0 customer-committed
    Fix: compressor service   — Operator ease 3/5 → ACTIONABILITY 48
```

(Actionability score = AvoidableImpact ÷ EaseOfFix — lower ease number = easier = higher actionability. V1.5 feature.)

---

## Implementation notes

### Code structure changes

**Revised 2026-07-14 (Entry 71): the package is `internal/functions/calculates/`, not a standalone `internal/cost/` — that's where `cost.go` already lives (90 lines, V0 shipped). New V1 files extend the existing package rather than founding a new one.**

```
internal/functions/calculates/
  cost.go                      # SHIPPED V0 — extend with V1 orchestration + provider discovery
  product_margin.go            # NEW V1 — ERP/MES product lookup (HIGH confidence) + CSV/Excel fallback (LOW confidence, reuses V0's existing table)
  customer_flag.go             # NEW V1 — commitment flag rule (flag only, no € — pricing rule)
  stock_risk.go                # NEW V1 — Enrichment #14, formula-derived, nullable when data absent
  downstream.go                # NEW V1 — propagation if line_layout configured
  restart_cost.go              # NEW V1 — lookup from config; V1.5 from MES
  line_layout.go               # NEW V1 — YAML loader
  confidence.go                # NEW V1 — per-enrichment confidence flag
  versioning.go                # NEW V1 — formula version stamped on each event
  # V1.5
  bottleneck.go                # NEW V1.5
  actionability.go              # NEW V1.5
  quality_scrap.go             # NEW V1.5 (needs MES connector)
  energy_peak.go               # NEW V1.5
  # V2
  forward_simulation.go        # NEW V2 — CFO killer
  spoilage.go                  # NEW V2
  regulatory.go                # NEW V2 (needs QMS/LIMS connectors)
```

**Sequencing (Entry 71)** — the SQL connector this all depends on for live ERP/MES reads is stuck at Day 2 of a 10-day plan (Entry 69: `internal/connections/` and the real `sql_query` handler don't exist yet; `sql_query` is still `stubConnector(...)` in `cmd/server/main.go:179`). Two pieces need **zero** connector work and can ship now: `downstream.go` + `line_layout.go` (pure YAML), `restart_cost.go` (pure YAML), and `confidence.go` + `versioning.go` (scaffolding only). `product_margin.go`'s ERP/MES path, `customer_flag.go`, and `stock_risk.go` all need a live provider query — don't start wiring them until `internal/connections/` + real `sql_query` land.

### Integration points

- **Pipeline engine** (`internal/pipeline/engine.go`): Impact Engine becomes a function in the function registry. Existing `calculate_cost` function name preserved for compatibility (not renamed to `calculate_business_impact` — see naming decision below); behavior upgraded.
- **Provider nodes, not a Fuzzy Join engine** (revised 2026-07-14, Entry 71): there is no `internal/fuzzy/of_state.go`. `calculate_cost` receives active-OF/product/stock/quality context because the pipeline engine already merges every upstream node's output into `params` keyed by node ID (`internal/pipeline/engine.go:196-211`) — provider context arrives the same way any other node's output does. `calculate_cost` discovers it by scanning `params` for `canonical_type`-tagged values, so provider nodes need no special engine support.
- **KG persistence** (`internal/kg/`): every event persisted with the full impact breakdown (not just the total). Enables drill-down + audit.
- **Dashboard** (`frontend/`): new "Top 3 Actions" widget consumes the actionability-sorted Pareto endpoint.
- **MCP server** (`internal/mcp/`): exposes new tool `kg_cost_breakdown(event_id)` for AI agent drill-down.

### Naming (locked 2026-07-14, Entry 71)

Keep `calculate_cost` as the function identifier. `docs/Cost_function.md` proposed renaming it to `calculate_business_impact` — treat that as UI-copy only ("Impact Engine" branding in the frontend label), not a YAML/config contract change. Renaming the function would break `config/pipelines/examples/cost_calculation.yaml` and `config/pipelines/pipeline_cost_calculation.yaml` for no functional gain.

### Configuration discoverability

A `config/cost.yaml.example` ships with the Edge Agent containing every option commented. Customer copies and edits. Onboarding wizard (U10 in the V1 inventory) walks Plant Manager through the 3 mandatory fields + offers to map a line layout.

### Dependency on provider nodes (revised 2026-07-14, Entry 71 — was "Dependency on Fuzzy Join")

**The Impact Engine is meaningless without active production context.** Without a `work_order`-tagged provider node in the pipeline, all V1 enrichments fall back to V0 defaults (generic margin, no customer flag, no downstream idle, no restart cost, no stock risk). That context now comes from an ordinary `sql_query` provider node the SQL connector work ships (`docs/mysql_connector.md`), wired into the pipeline via `depends_on` — not a bespoke join engine. Practically this still means: the SQL connector (Entry 69) gates the ERP/MES-dependent enrichments, but no longer gates the whole Impact Engine — the config-only enrichments (downstream idle, restart cost) ship independently of it.

---

## Validation plan (how we know V1 Impact Engine is credible)

Before declaring V1 ready:

1. **Reproduce a 1-week real factory dataset** with the engine. Compare totals to operator-reported daily losses. Target: within ±15% of operator estimate.
2. **Sanity-check with first pilot customer's Plant Manager**: show them 10 events with full breakdown. Ask: "Does this match your intuition?" Iterate until at least 8/10 are validated.
3. **CFO drill-down test**: pick a high-impact event, expand all the way down (time loss formula → product margin lookup → ERP record). Every step must be readable in the dashboard. No "magic numbers".
4. **Confidence-level distribution**: across all events in the pilot week, what % are HIGH / MEDIUM / LOW confidence? Target: >70% HIGH after first-week setup.
5. **Determinism**: re-run the engine on the same week's data → numbers match exactly.

---

## Open questions

| Question | Decision needed by | Notes |
|---|---|---|
| Do we ship V1 with manual `setup_costs.yaml` or wait for MES auto-compute? | V1 sprint planning | Manual is simpler. MES auto requires customer with mature MES. Recommend manual V1 + MES V1.5. |
| Confidence-level threshold for "trustworthy event" | V1 design | Recommend ≥4 of 5 enrichments resolved = HIGH. ≥3 = MEDIUM. <3 = LOW (still computed, flagged). Updated 2026-07-14 for the added Stock-shortage enrichment. |
| Per-vertical configuration defaults | V1.5 (after first pilot in each vertical) | Pharma probably needs different defaults than agrifood. Capture as we onboard each vertical. |
| Forward simulation algorithm (V2) | V2 design — months away | Options: deterministic counterfactual (simple) vs Monte Carlo (more accurate, slower) |
| Customer-managed override of formula | V2 | Some sophisticated customers will want to plug in their own cost model. Architect for it but don't ship V1. |
| Open-source the Impact Engine spec? | 2028 license reconsideration | The math should be defensibly open even if the code stays closed-source — builds trust. |

---

## Pitch language (how this shows up externally)

For investor deck / customer pitch / X content:

**One-liner:**
> *"MindSet's Impact Engine doesn't just count micro-stops — it prices every event in € using your ERP + MES + recipe + customer-commitment context. Then it ranks by what's actionable, not just what's expensive. The Plant Manager sees the top 3 things to fix this week — not 200 anonymous events."*

**Three-line version (CFO-flavored):**
> *"Detecting events is table stakes. Pricing them in € — using your ERP product margins, your MES recipes, your customer-commitment risk, your stock-shortage exposure — is the moat. Ranking by actionable financial impact is the decision-time reducer. Our 5 V1 enrichments turn operational data into CFO-grade prioritization. Built on the OT/IT reconciliation that competitors don't have."*

**X content thread material** (Mohamed):
> Thread: "Why I rewrote our cost function this month. From `duration × hourly_cost` to a reconciliation-powered Impact Engine. Here's the 5 enrichments that change everything for the Plant Manager."

---

## Today's takeaway

| | |
|---|---|
| **Current state** | V0 shipped: 1 component (TimeLoss-style, `duration × hourly_rate/60`) + a per-product CSV/Excel rate override — thinner than this doc originally assumed. Doesn't leverage reconciliation yet. Plant Manager + CFO sanity-check will kill trust. |
| **V1 target (3-4 weeks on top of Track 1)** | 5 enrichments: per-product margin · customer-commitment flag (FLAG only) · downstream idle · restart cost · stock-shortage risk (PRICED, formula-derived). Plus "Top 3 Actions" dashboard widget. Config-only pieces (downstream idle, restart cost, confidence/versioning scaffolding) ship independent of the SQL connector; the rest wait on it (Entry 69). |
| **V1.5 target** | Bottleneck ID · actionability score · quality scrap · energy peak. |
| **V2 target (the CFO killer)** | Forward simulation: *"fix cause #1 → recover Y €/week, conditional on cause #2 becoming next bottleneck"*. |
| **Critical dependency** | Provider nodes (`sql_query`, canonical_type-tagged) must exist for the ERP/MES-dependent enrichments — not a bespoke Fuzzy Join engine (revised 2026-07-14, Entry 71). Config-only enrichments don't wait on this. |
| **Trust principle** | Every number traceable. No black boxes. Customer-tunable. Versioned. Priced-or-flagged, never asserted (the pricing rule, Entry 71). |
