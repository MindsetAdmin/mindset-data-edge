# Self-creating pipelines — 18 concrete examples + feasibility

_The MindSet analog to LemonLime's "self-creating automations" — patterns the Pipeline Studio could detect in the KG and turn into ready-to-accept pipelines._
_See Entry 67 for context._

## How it would work

1. Detection module runs periodically (e.g., nightly) over the KG's accumulated events
2. Each detector matches a specific pattern template
3. When a match crosses confidence threshold, system generates a **pipeline template** (YAML) + a plain-language explanation
4. Surfaces in the DataOps Studio as "Suggestions" — user reviews, edits if needed, clicks Accept → pipeline live
5. User never has to know what to build; the system proposes based on what it sees

**Precision > recall.** In industrial contexts, a bad suggestion firing 3 alarms a night gets the tool killed. Better: 5 well-scoped suggestions/month that each save real money.

---

## The 18 examples, grouped by pattern family

### A. Temporal / recurring patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 1 | `machine1` has ≥5 micro-stops in the same hour, on 3+ consecutive shifts | Alert pipeline: notify shift lead + tag Cause="recurring" | **EASY** — count-aggregation + threshold, Go stdlib |
| 2 | Line output drops >15% within 30 min of every shift change | Dashboard widget + operator handover checklist | **EASY** — window aggregation over shift boundaries |
| 3 | Batch duration for `PROD-A03` trending upward: +12% over 4 weeks | Alert for process review; add to weekly ops report | **MEDIUM** — trend regression on batch durations |

### B. Correlation / cascade patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 4 | `machine2` slowdown always precedes `machine3` stop within 15 min | Upstream alert: pre-warn `machine3` operator when `machine2` slows | **MEDIUM** — lag-correlation between event streams |
| 5 | `PROD-B01` batches fail 3× more often when viscosity metric > 890 | Threshold-check pipeline: flag future PROD-B01 batches near limit | **MEDIUM** — grouped correlation over quality outcomes |

### C. Root cause patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 6 | 40% of `machine1` micro-stops tagged with same Cause node | Maintenance ticket workflow + vendor escalation if repeated | **EASY** — cause-aggregation via KG traversal |
| 7 | Cost concentration: `PROD-C01` accounts for 45% of downtime cost this month | Focused-optimization pipeline: deep-dive dashboard for that SKU | **MEDIUM** — needs Impact Engine cost data joined to events |

### D. Data quality patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 8 | `franco_de_port` on supplier X unchanged for 3 years while peers update quarterly | Data-quality alert flagging STALE; ask acheteur to verify | **EASY** — timestamp diff + peer comparison |
| 9 | `operator_id` null on 30% of last week's work orders | Data-entry gap alert to production planner | **EASY** — null-count aggregation |
| 10 | `target_rate` for `PROD-A05` is 3× the average of similar-category products | Outlier flag: probable data-entry error | **EASY** — z-score against peer group |

### E. Throughput / rate patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 11 | Machine X producing 20% below its 30-day moving average this week | Performance-investigation alert with root cause hint | **EASY** — moving average + deviation |
| 12 | Sensor drift: reactor temperature climbing 0.3°C/week over 6 weeks | Predictive-maintenance alert before hard alarm fires | **MEDIUM** — linear regression on time-series |

### F. Sequence / state-machine patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 13 | Line frequently transitions `Running → Fault → Idle` instead of expected `Running → Idle → Clean → Running` | State-machine anomaly alert; probable process shortcut | **MEDIUM** — state-graph walk analysis |
| 14 | Batches usually get 3 quality readings but `PROD-B01` batches get only 1 | QA-completeness alert to lab | **EASY** — count grouping by product |

### G. Supplier / procurement patterns (Deroche-relevant)

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 15 | Supplier X late on 3 of last 8 deliveries | Supplier-scoring pipeline + monthly review alert | **EASY** — count aggregation on delivery events |
| 16 | Single supplier > 60% of category volume | Diversification alert (matches Deroche cahier des charges) | **EASY** — share % calculation |

### H. DLC / expiry patterns (Deroche-relevant)

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 17 | `PROD-A07` has DLC losses on ≥30% of end-of-week days | Cadence adjustment: recommend reducing order size or shifting delivery day | **EASY** — day-of-week grouping on waste events |

### I. Multi-machine / cross-site patterns

| # | Detected pattern | Suggested pipeline | Feasibility |
|---|---|---|---|
| 18 | `machine1` has >3× the stops of peer machines with same product | Comparative dashboard + investigation ticket | **EASY** — grouped aggregation + rank |

---

## Feasibility summary

| Feasibility | Count | Technical requirements | Time-to-build |
|---|---:|---|---|
| **EASY** | 12 / 18 | Go stdlib + basic aggregation on KG (SQLite queries + go stdlib `math` / optional `gonum/stat`) | 1–2 weeks per detector |
| **MEDIUM** | 6 / 18 | Add trend regression + lag correlation + state-graph walks. Introduce `gonum/stat` as a dependency. | 2–4 weeks per detector |
| **HARD** (not shown) | 0 in this list | Multi-variate anomaly detection with ML autoencoders, deep sequence models, NLP over free-text — V2+ | Months, needs ML expertise |

**All 18 examples are buildable with our current stack** (Go, KG in SQLite, event history). None require ML frameworks or GPUs. This is the important point: we don't need a data-science team to ship a first version of the suggestion engine.

---

## Architecture — how to build it

```
KG event history (SQLite)
       │
       ▼
Detector registry (internal/suggestions/)
  ├── detectors/temporal/     (patterns 1, 2, 3)
  ├── detectors/correlation/  (patterns 4, 5)
  ├── detectors/quality/      (patterns 8, 9, 10)
  ├── detectors/throughput/   (patterns 11, 12)
  ├── detectors/sequence/     (patterns 13, 14)
  ├── detectors/supplier/     (patterns 15, 16)
  ├── detectors/dlc/          (patterns 17)
  └── detectors/multi/        (pattern 18)
       │
       ▼
Each detector returns:
  {
    detected: bool,
    confidence: 0-1,
    evidence: [event IDs],
    explanation_fr: "Ligne 3 a 5 micro-arrêts dans la même heure sur 3 shifts consécutifs",
    proposed_pipeline_yaml: "..."
  }
       │
       ▼
Suggestion queue (SQLite table)
       │
       ▼
DataOps Studio UI: "Suggestions" panel
  Each suggestion: [Accept] [Edit] [Dismiss] [Never suggest this again]
```

### Precision guardrails (non-negotiable)

- Every detector has a minimum confidence threshold (default 0.8) — no suggestion below it
- User feedback loop: dismissed suggestions decay that detector's future weight
- Rate limit: max 3 suggestions/user/week to avoid alert fatigue
- Explainability: every suggestion carries a plain-French explanation of what triggered it + the specific evidence (event IDs)

---

## Sequencing recommendation

**V1.5 M0-M1 — Foundation:**
- Build the `internal/suggestions/` registry framework
- Ship 3 EASY detectors as proof of concept (recommend: patterns 1, 6, 15)
- Build the DataOps Studio Suggestions panel

**V1.5 M1-M2 — Deroche-facing:**
- Add patterns 15, 16, 17 (supplier + DLC — directly matches Deroche cahier des charges)
- These are our proof to Deroche that we don't just move data — we actively surface what to improve

**V1.5 M2-M4 — Breadth:**
- Roll out remaining EASY detectors
- Start MEDIUM detectors (patterns 3, 4, 12, 13) once V1a+V1b MSSQL are in production

**V2 — Advanced:**
- ML-based anomaly detection over multi-variate signals
- Cross-site pattern learning across MindSet customers (with opt-in — federated by default)

---

## Strategic angle

Building this positions MindSet ahead of both HighByte/Litmus (no suggestions, purely data-plumbing) AND LemonLime (SaaS suggestions, not OT). The industrial domain knowledge required for these detectors — knowing that a `Running → Fault → Idle` cycle matters, knowing that supplier concentration > 60% is a red flag — is the moat LemonLime can't cheaply replicate.

**Talking point for Cécilia:**
> *"On ne se contente pas de livrer la donnée. Notre plateforme regarde en permanence ce qui se passe dans l'usine et suggère à l'équipe ce qu'il vaut la peine d'optimiser en premier — classé par impact économique. Pas d'IA magique, pas de boîte noire : des patterns déterministes que l'acheteur ou l'ingénieur peut vérifier."*
