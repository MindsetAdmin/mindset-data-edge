# MINDSET DATA — Project Documentation

> **Vision:** Connect every machine, system, and data flow in a factory to a single reliable and exploitable source of truth — the Unified Namespace (UNS) — to transform every signal into a business decision.

---

## Table of Contents

1. [The Problem](#1-the-problem)
2. [The Solution](#2-the-solution)
3. [Client Segmentation](#3-client-segmentation)
4. [Technical Architecture](#4-technical-architecture)
5. [Protocols & Connectors](#5-protocols--connectors)
6. [Auto-Discovery & Detection](#6-auto-discovery--detection)
7. [The Unified Namespace (UNS)](#7-the-unified-namespace-uns)
8. [Product Modules](#8-product-modules)
9. [Use Cases](#9-use-cases)
10. [Full Roadmap](#10-full-roadmap)
11. [GitHub Repository Structure](#11-github-repository-structure)
12. [Infrastructure & Distribution](#12-infrastructure--distribution)
13. [Tech Stack](#13-tech-stack)
14. [Security & Compliance](#14-security--compliance)
15. [Tech Moat](#15-tech-moat)
16. [Development Workflow](#16-development-workflow)
17. [Hardware Requirements](#17-hardware-requirements)

---

## 1. The Problem

### 1.1 Four structural failures of the mid-market factory

**Data silos**
Machines, ERP, MES, SCADA generate siloed data that never talks to each other. Hours are lost reconciling numbers before anyone can make a decision.

**Operational noise**
Industrial signals are dense and multiple. Managers decode the past manually while operational losses accumulate in silence.

**Data rich, decision poor.**
Some factories already stream data in real time — and get nothing out of it. Every new source requires custom development. Every dashboard visualizes without explaining. Teams know something is wrong. They don't know how much it costs, or what to fix first.

**Tribal knowledge not captured**
The operational knowledge of experienced technicians — the real causes of stops, the settings that work — is never structured. It disappears when they retire.

### 1.2 The consequence

> *"Slow decision loops are the primary reason manufacturing productivity stays below 1% per year."* — McKinsey

Unreconciled data slows decisions and silently destroys margin.

---

## 2. The Solution

Mindset Data is the **unified, reliable, real-time data infrastructure layer** that connects all factory sources (machines, ERP, MES, energy meters) to a UNS that contextualizes every signal into financial impact — without replacing existing systems, without additional hardware, without a heavy IT project.

### The four-step flow

```
01 CONNECT        02 CONTEXTUALISE   03 VISUALISE      04 ACT
──────────────    ────────────────   ────────────────  ────────────────
Auto-discovery    UNS auto-generated Top 3 priorities  Push alerting
OT network        ISA-95 ontology    coded in €        Email / Slack
OPC-UA, Modbus,   Knowledge Graph    3-tab dashboard   Cause capture
S7, SQL, Files    Fuzzy Join OT/IT   Gantt / Pareto    KG enriched
Zero dev client   Enriched ongoing   ROI simulated     Loop closed
```

### Product segmentation

| Step | Why this step? | POC — Value proven from Day 1 | Post-POC |
|------|---------------|-------------------------------|----------|
| **Connection** | Connect to equipment via native protocols. Hardware agnostic, deployed in under half a day, without touching existing infrastructure. **Promise: Zero infra impact, deployment in half a day.** | Go Edge Agent, OPC-UA auto-discovery, Modbus, manual entry of business parameters. Base graph with UNS Topics. → *"We connect to your machines in half a day without touching your infrastructure."* | Extension to legacy protocols (S7, MQTT…). Native ERP/MES connector. |
| **Contextualisation** | The Knowledge Graph gives business meaning to raw signals. Pre-built Functions transform OT tags into operational and financial events — without a developer. | Pre-built micro-stop + energy pipeline instantiated automatically on validated tags. → *"We automatically build the context of your main line from day one."* | Graph enriched continuously. New no-code Functions per use case (quality, changeover…). |
| **Visualisation & Recommendation** | Eliminate time spent searching and interpreting. Make the urgent issue visible — and its cost. | Top 3 daily priorities coded in €. → *"Sensor X stopped 180 times this week. Aggregated, these micro-stops total 2h of invisible downtime and cost 2,160€ on the chocolate biscuit production batch."* | Multi-site cockpit. Real-time P&L sliding financial dashboard. |
| **Action** | Turn the recommendation into an executed decision. Close the loop detection → field resolution with traceability. | Push alert (Email/Slack) when a € threshold is crossed. Cause entered in 1 click. → *"Your team receives the right information, at the right time, with the right priority."* | Automated continuous improvement loop. Every resolved action feeds the KG. |

---

## 3. Client Segmentation

### IT/OT Manager
- **Pain:** Integration complexity and maintenance of factory data.
- **Motivation:** Simplify deployment, reduce maintenance, standardize data flows.
- **Solution:** Native connection to industrial protocols and automatic creation of the contextual data foundation — without custom development.
- **Benefits:** Unified operational context / Integration without custom code / Real-time data access everywhere
- **Budget:** Technical prescriber. Maintenance or IT tooling budget — validates infra before any decision.
- **Blurb:** *"Connecting the shop floor to the rest of the business usually means heavy custom development and maintenance hell. Our solution integrates transparently into your current architecture: it automatically maps your factory data to create a unified model, without a single line of custom code."*

### Operations Director
- **Pain:** Too much noise and delay between machines, lines and sites. Difficulty prioritizing next actions.
- **Motivation:** Produce more and sell more. Eliminate noise at group scale to maximize volumes and global performance.
- **Solution:** Unified contextual layer that transforms fragmented signals into field truth, with integrated governance and automatic prioritization by operational impact.
- **Benefits:** Faster root cause analysis / Real OEE vs declared OEE / Continuous improvement based on field reality
- **Budget:** Decision maker / Co-funder. Operational excellence budget or multi-site CAPEX.
- **Verbatim:** *"To be productive in the daily meeting on what to prioritize, you need the right data and to know how to interpret it. I have data engineers as intermediaries — I want it to be doable by whoever understands the shop floor."*
- **Blurb:** *"Standardizing the performance of your sites shouldn't require overhauling your IT infrastructure. Our solution filters industrial background noise in real time to surface only critical drifts. You move from a culture of reacting to emergencies to global steering based on clear, standardized priorities."*

### Production Director / Plant Manager
- **Pain:** Too much noise and useless alerts on the lines daily. Difficulty understanding cause-effect links to prioritize.
- **Motivation:** Produce more and sell more. Maximize availability rate, make production reliable.
- **Solution:** Automatic micro-stop detection and cause-effect correlation at source, sub-second latency. Every drift translated into financial impact and priority action.
- **Benefits:** Immediate downtime reduction / Focus on the real bottleneck / Instant understanding cause → cost → action
- **Budget:** Main field buyer. Line maintenance or continuous improvement budget. Autonomous signing threshold (<30k€/site).
- **Verbatim:** *"If we don't collect information with a minimum of precision, it's impossible to interpret KPIs and make decisions. I know my line has problems, but what I need is to understand their financial impact on my revenue and know what to prioritize without spending hours searching."*
- **Blurb:** *"To produce more and sell more, your teams need clarity, not more data. We eliminate noise on your lines by analyzing causes and effects at source, with sub-second latency. Your managers no longer have to guess: the system tells them exactly which priority action to take to unblock the bottleneck of the day."*

### CFO / CEO
- **Pain:** Lost margins visible too late — end of month in the P&L. Decisions based on filtered and delayed data.
- **Motivation:** Consolidated multi-site visibility to anticipate and reduce costs immediately.
- **Solution:** Automatic real-time translation of operational field drifts into direct financial impact — visible before losses are consumed.
- **Benefits:** Real-time financial visibility / Decisions based on avoided costs / Immediate measurable P&L impact
- **Budget:** Strategic buyer. Global strategic investment budget — arbitrates multi-site expansion.
- **Blurb:** *"Today, your factory drifts only show up in your P&L at end of month, when it's too late to correct. We connect and contextualize your field data in real time to translate every operational loss into immediate financial impact. You don't wait for the report: you steer your margins and see secured EBITDA before losses are even consumed."*

---

## 4. Technical Architecture

### 4.1 Global overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         CLIENT NETWORK (OT + IT)                    │
│                                                                     │
│  OT SOURCES                      IT SOURCES                        │
│  ├── SCADA / PLC (OPC-UA)        ├── ERP (SQL / REST API)          │
│  ├── Siemens PLC (S7)            ├── MES (SQL / REST API)          │
│  ├── Legacy PLC (Modbus TCP)     ├── Production Orders / Recipes   │
│  └── Energy meters (Modbus TCP)  └── Hourly costs / Margins        │
│           │                               │                        │
│           └───────────────┬───────────────┘                        │
│                           ▼                                        │
│              ┌────────────────────────┐                            │
│              │     EDGE AGENT (Go)    │  ← Docker, client's PC     │
│              │                        │                            │
│              │  Discovery Layer       │                            │
│              │  ├── Network scanner   │  ← Auto-detect endpoints   │
│              │  ├── OPC-UA browse     │  ← Node tree extraction    │
│              │  ├── Modbus scan       │  ← Register fingerprint    │
│              │  └── S7 scan           │  ← DB scan                 │
│              │                        │                            │
│              │  Intelligence Layer    │                            │
│              │  ├── Behavioral infer  │  ← Live pattern matching   │
│              │  ├── SLM Phi-3 local   │  ← Tag classification      │
│              │  ├── Local UNS ISA-95  │  ← In-memory routing       │
│              │  └── Fuzzy Join        │  ← OT/IT reconciliation    │
│              │                        │                            │
│              │  Processing Layer      │                            │
│              │  ├── Rules engine      │  ← Micro-stops detection   │
│              │  ├── Energy rules      │  ← Waste detection         │
│              │  └── Cost model (€)    │  ← Real-time € calculation │
│              │                        │                            │
│              │  Storage Layer         │                            │
│              │  ├── SQLite 7-15 days  │  ← Local ring buffer       │
│              │  └── Push → Historian  │  ← Client's own system     │
│              │                        │                            │
│              │  Serving Layer         │                            │
│              │  ├── Local dashboard   │  ← localhost:3000          │
│              │  ├── Local alerting    │  ← Direct SMTP/Slack       │
│              │  └── Ollama / Phi-3    │  ← SLM fully local         │
│              └───────────┬────────────┘                            │
│                          │ HTTPS Push-only                         │
│                          │ Aggregated snapshots only               │
│                          │ (KG delta + TRS summary + alerts)       │
└──────────────────────────┼──────────────────────────────────────────┘
                           ▼
              ┌────────────────────────┐
              │    CLOUD (Scaleway FR) │  ← Minimal footprint
              │                        │
              │  ├── KG aggregation    │  ← Cross-site (V1+)
              │  ├── Remote dashboard  │  ← CEO/CFO access
              │  ├── Site management   │  ← API keys, auth
              │  ├── Backup KG         │  ← Snapshots
              │  └── Alerting relay    │  ← If local down
              └────────────────────────┘
```

### 4.2 Core principles

| Principle | Implementation |
|-----------|---------------|
| Zero additional hardware | Edge Agent runs on existing client infrastructure (VM, industrial PC) |
| Zero raw data in cloud | Only transformed events and aggregated snapshots go up |
| Push-only | Outbound HTTPS only, zero inbound open ports |
| Read-only on source systems | Zero writes to PLC, SCADA, ERP |
| Zero manual connection work | Network scan + behavioral inference + SLM classification |
| Maximum local processing | Rules engine, cost model, Fuzzy Join, dashboard — all at the Edge |
| Offline resilience | Full operation if cloud unreachable, queues sync on reconnect |
| Cloud sovereignty | FR certified cloud (Scaleway/OVH) or BYOC |

### 4.3 Storage strategy

```
EDGE (local — rolling window)     CLIENT HISTORIAN          YOUR CLOUD
─────────────────────────────     ─────────────────         ──────────
SQLite ring buffer                PI System /               Aggregated
7 to 15 days configurable         Wonderware /              KG snapshots
                                  InfluxDB /                only
Raw events + causes + costs €     MSSQL / TimescaleDB
                                  (already exists)
Auto-purge after TTL              Client owns their         You own the
                                  long-term history         context layer
Push enriched events              via your SQL/REST         No raw history
to historian (V1)                 connector                 ever stored
```

### 4.4 Local UNS vs Cloud UNS

```
LOCAL UNS (Edge — volatile)              CLOUD UNS (Persistent)
────────────────────────────             ───────────────────────────
Lives in memory, rebuilt on restart      Lives in PostgreSQL, grows forever
Current tag → UNS topic mapping          Full site fingerprint since day 1
Used for: real-time routing,             Used for: KG API, remote dashboard,
rules engine, cost calculation           AI agents (V2+), cross-site views
Scope: present moment                    Scope: full history + context

"ns=2;s=Line1.Motor.Speed"              "site/usine-nord/line1/motor/speed
→ site/line1/motor/speed"                → Vitesse_Convoyeur
→ apply micro-stop rule                  → 847 micro-stops since install
→ calculate cost €                       → primary cause: bourrage (67%)
→ push event to cloud"                   → avg cost: 18€/event"
```

---

## 5. Protocols & Connectors

**Strategy: No Kepware, no third-party middleware. Native drivers only.**
This removes vendor dependency and is a direct sales argument to IT security teams.

### V0 — Build now

| # | Protocol | ETI Frequency | Deploy Difficulty | Why |
|---|----------|--------------|-------------------|-----|
| 1 | **OPC-UA** | ⭐⭐⭐⭐⭐ | 🟢 Low | Self-describing, SCADA cheat code, native auto-discovery, EU Data Act standard. Absolute foundation. |
| 2 | **Modbus TCP** | ⭐⭐⭐⭐⭐ | 🟡 Medium | Ubiquitous on legacy equipment and energy meters. Device fingerprint DB needed. |

### V1 — Months 2-4

| # | Protocol | ETI Frequency | Deploy Difficulty | Why |
|---|----------|--------------|-------------------|-----|
| 3 | **Siemens S7** | ⭐⭐⭐⭐ | 🟡 Medium | 30-40% of European industrial park. Go lib `gos7`. Direct PLC access without SCADA. |
| 4 | **SQL** | ⭐⭐⭐⭐⭐ | 🟢 Low | ERP/MES/historians (SAP, Sage, Dynamics, PI System). Fuzzy Join foundation. |
| 5 | **REST API** | ⭐⭐⭐⭐ | 🟢 Low | Modern ERP (SAP S/4HANA, Oracle Cloud, D365). YAML config, configurable polling. |
| 6 | **Files** (CSV/Excel/JSON) | ⭐⭐⭐⭐ | 🟢 Low | ETI that export to flat files. Unlocks clients without API. |
| 7 | **FTP / SFTP** | ⭐⭐⭐ | 🟢 Low | Automatic export from legacy MES/ERP. Complement to Files. |

### V2 — Months 4-8

| # | Protocol | ETI Frequency | Deploy Difficulty | Why |
|---|----------|--------------|-------------------|-----|
| 8 | **MQTT** | ⭐⭐⭐ | 🟢 Low | Recent IIoT gateways, post-2018 equipment. Broker often already present. |
| 9 | **Sparkplug B** | ⭐⭐ | 🟢 Low | MQTT + native ISA-95 structured payload. If client has Sparkplug, UNS is near-free. |
| 10 | **EtherNet/IP** | ⭐⭐⭐ | 🔴 High | Rockwell/Allen-Bradley. Connect via PLC/SCADA that aggregates it — never directly. |
| 11 | **Ignition** | ⭐⭐⭐ | 🟢 Low | Agrifood/pharma SCADA. Exposes OPC-UA natively — transparent connection. |
| 12 | **InfluxDB** | ⭐⭐⭐ | 🟢 Low | Time-series historian common in digitalized ETI. Simple API. |
| 13 | **OPC-DA** | ⭐⭐⭐ | 🔴 High | Legacy Windows COM/DCOM. Connect via OPC-DA→UA wrapper — never implement natively. |
| 14 | **MTConnect** | ⭐⭐ | 🟡 Medium | CNC/machining/metallurgy standard. XML over HTTP — simple Go implementation. |
| 15 | **BACnet/IP** | ⭐⭐ | 🟡 Medium | Building/HVAC. Energy use case on energy-intensive sites. |
| 16 | **Omron FINS** | ⭐⭐ | 🟡 Medium | Omron PLC — agrifood and pharma. Available Go lib. Recurring niche. |
| 17 | **MongoDB** | ⭐⭐ | 🟢 Low | Some modern MES. Native Go driver. |
| 18 | **RabbitMQ** | ⭐⭐ | 🟢 Low | IT-side message broker. Clients with microservices architecture. |
| 19 | **SMB / CIFS** | ⭐⭐ | 🟢 Low | Windows network shares — Excel/CSV on server. Files complement. |

### V3+ — Month 9+

| # | Protocol | Deploy Difficulty | Why |
|---|----------|-------------------|-----|
| 20 | **Kafka** | 🟡 Medium | High-volume streaming multi-site. Overkill in V0/V1. |
| 21 | **Redis** | 🟢 Low | Real-time cache. Intermediate layer, rarely direct source. |
| 22 | **AWS S3 / Azure Blob** | 🟢 Low | Cloud-first clients dumping OT data to object storage. |
| 23 | **AWS IoT SiteWise** | 🟡 Medium | AWS industrial clients. Relevant for marketplace partnership. |
| 24 | **Azure IoT Hub** | 🟡 Medium | Advanced Microsoft industrial clients. |
| 25 | **Google Pub/Sub** | 🟢 Low | Rare in European manufacturing ETI. |
| 26 | **LoRaWAN** | 🔴 High | Long-range low-power IoT. Isolated field sensors. Niche. |
| 27 | **NATS** | 🟢 Low | Modern message broker. Rising in cloud-native IIoT. |
| 28 | **Elasticsearch** | 🟢 Low | IT log management. Rare as direct OT source. |

### Outbound alerting (transversal — from V0)

| Protocol | Usage | Version |
|----------|-------|---------|
| **SMTP** | Email alerting | V0 |
| **Slack webhook** | Team alerting | V0 |
| **Microsoft Teams webhook** | Microsoft-centric ETI | V1 |

---

## 6. Auto-Discovery & Detection

### 6.1 Network scanning — 100% automatic

```
Subnet scan (client provides network access):
  → Port 4840 (OPC-UA)  → Handshake    → Browse node tree
  → Port 502  (Modbus)  → Probe FC03   → Register scan + fingerprint
  → Port 102  (S7)      → ISO TSAP     → DB scan
  → Port 1883 (MQTT)    → CONNECT      → Subscribe wildcard #
```

Zero manual input. Agent receives network access, finds all endpoints in minutes.

### 6.2 OPC-UA auto-discovery pipeline

```
1. Browse node tree  → extract all tags
2. Noise filter      → remove constants, duplicates, signals > 10Hz
3. SLM (Phi-3)       → semantic classification of readable tags
4. Behavioral        → behavioral inference on opaque tags
5. Confidence score  → High / Medium / Low per tag
6. UNS mapping       → Site > Area > Work Center > Work Unit > Tag
```

### 6.3 Behavioral inference — the real differentiator

Observe live behavior for 10-15 minutes without knowing tag names:

```
Boolean TRUE >80% of time, FALSE transitions in 30s-3min bursts
→ High confidence: Etat_Machine (Run/Stop)

Monotonically incrementing integer, occasional reset
→ High confidence: Compteur_Pieces or Compteur_Cycles

Float drops to 0 exactly when Boolean goes FALSE
→ High confidence: Vitesse_Instantanee

Boolean TRUE <5s at exact moment of machine stop
→ High confidence: Capteur_Defaut / Capteur_Bourrage

Float non-zero even when machine is stopped
→ High confidence: Compteur_Energie (energy waste signal)
```

**Works on opaque Modbus registers and S7 DBs too** — behavior reveals role without name.

### 6.4 Modbus device fingerprint database

- Internal DB of 20-30 most common devices (Schneider M340/M580, Siemens S7-1200, Danfoss VFD, SEW Eurodrive…)
- IP/MAC fingerprint → register map loaded automatically
- Unknown device → behavioral inference on populated registers
- Every new documented device enriches the DB for all future deployments

### 6.5 Confidence levels and UX

| Confidence | Condition | UX |
|------------|-----------|-----|
| High (>90%) | Readable name + consistent type | Pre-checked, 1-click validation |
| Medium (60-80%) | Abbreviated name or behavioral inference | Score shown, confirmation required |
| Low (<50%) | Opaque address, unknown | Flagged red, optional manual entry |

**V0 rule:** require a SCADA with OPC-UA enabled. The SCADA already aggregates all fieldbus protocols and exposes named tags. Connection = 1 OPC-UA endpoint, tags already readable.

---

## 7. The Unified Namespace (UNS)

### 7.1 ISA-95 hierarchy

```
Enterprise
  └── Site (Usine Paris-Nord)
        └── Area (Ligne 1, Ligne 2)
              └── Work Center (Machine A, Machine B)
                    └── Work Unit (Pressure Sensor, Motor X)
                          └── Tags (Etat_Machine, Vitesse, Conso_kW)
```

### 7.2 Auto-generation of the Knowledge Graph

1. Network scan → automatic detection of all endpoints
2. Auto-discovery → OPC-UA node tree extraction
3. SLM + behavioral inference → tag classification
4. Automatic mapping to ISA-95 hierarchy in **under 48h**
5. Knowledge Graph enriched continuously (tribal knowledge, validated causes, learned patterns)

### 7.3 What the UNS enables

- Every event **tagged** with its context (product, production order, team, shift)
- Every signal **attributed** to its financial impact in real time
- Native or third-party AI agents can **query reliable context** via semantic API
- Every deployment generates a **unique cumulative site fingerprint**

---

## 8. Product Modules

### Module 1 — Edge Agent & Auto-Discovery
**Differentiation: Strong** — Zero manual work vs 2-3 weeks of mapping at integrators.

| Component | Detail |
|-----------|--------|
| Language | Go (lightweight, static binary, low footprint) |
| Deployment | Docker — one command, any OS (Windows/Linux) |
| Network scanner | Auto-detect OPC-UA / Modbus / S7 / MQTT endpoints |
| OPC-UA auto-discovery | Full node tree browse, all tags extracted |
| Behavioral inference | Live behavioral pattern matching |
| SLM classification | Phi-3 local (Ollama) — sovereign, zero latency |
| Noise filtering | Constants, duplicates, signals >10Hz auto-excluded |
| Modbus device DB | Fingerprint of 20-30 most common devices |

### Module 2 — Pipeline & UNS
**Differentiation: Moderate** — Auto-generation and no-code extensibility are differentiating.

| Component | Detail |
|-----------|--------|
| Transformation | Redpanda Connect — declarative YAML pipelines |
| Structure | ISA-95 UNS auto-generated |
| No-code Functions | Detection, cause correlation, cost calculation, aggregation — pre-built per use case |
| Extensibility | New data source = new YAML config, zero dev ticket |

### Module 3 — Deterministic Rules Engine
**Differentiation: Strong** — Pre-configured rules per use case, zero client code.

| Use Case | Rule |
|----------|------|
| Micro-stop | `IF Etat_Machine=Stop AND 30s < duration < 3min → Micro-Stop event` |
| Energy waste | `IF Etat_Machine=Stop AND Energy > Threshold → Waste alert` |
| Schedule deviation | `IF Actual_Duration_OT > ERP_Schedule_Time + 15% → Margin Eroded` |
| Configurable threshold | UI wizard <15 min, zero dev |

### Module 4 — Fuzzy Join / Temporal Engine
**Differentiation: Strong** — Universal unsolved problem in mid-market. Makes data 100% reliable.

| Component | Detail |
|-----------|--------|
| Algorithm | Sliding window ±10 min around physical start detected at Edge |
| Detection | Physical start peak (Run transition on Etat_Machine) → matched to ERP production order |
| Result | Every micro-stop and every kWh tagged with correct production order and product |
| Precision | Sub-second at the Edge |

### Module 5 — Cost Model
**Differentiation: Strong** — Nobody does this immediate calculation in the mid-market.

**V0 — Manual entry (3-field onboarding wizard):**
```
Line_Hourly_Cost (€/h)        → Plant Manager
Theoretical_Rate (units/h)    → Plant Manager
Product_Margin (€/unit)       → Plant Manager
```

**V1 — Automatic ERP import via SQL/API connector**

**Edge calculation (V0 and V1):**
```
Time loss      = Micro-Stop duration (min) × Line_Hourly_Cost
Production loss = Nb_Micro-Stops × Rate × Product_Margin
Energy loss    = Off-Prod consumption (kWh) × Energy_Price
──────────────────────────────────────────────────────────
Total impact € = Time loss + Production loss + Energy loss
```

### Module 6 — Dashboard (Local + Cloud)
**Served locally from the Edge Agent. No internet required.**

**V0 — 3 tabs:**

| Tab | Visual | Purpose |
|-----|--------|---------|
| **Real Time** | Gantt timeline (Green=Run, Red=Stop, Orange=Setup) | Validate the tool captures reality to the second |
| **Micro-Stops** | "Xh lost this week" counter + Pareto causes by € | The visual shock — see money lost, know where to act |
| **Impact & ROI** | Real TRS vs declared + gain potential if cause #1 resolved | Justify action and estimate Mindset Data ROI |

**V1 — Additional tabs:**

| Tab | Visual | Purpose |
|-----|--------|---------|
| **Schedule Gaps** | Production orders: ERP Planned vs Edge Actual vs Gap € | Orders that destroyed margin |
| **Energy** | Consumption by production order + off-prod waste | Energy-heavy products and inter-batch leaks |
| **Financial Impact** | "Margin lost this month: X€" + financial Pareto | CFO language |

### Module 7 — Alerting
**From V0 — outbound only, never polling.**

| Trigger | Channel | Example |
|---------|---------|---------|
| Micro-stop exceeds € threshold | Email / Slack | "Line 2 — 3 jams in 1h = 240€ lost" |
| Energy waste detected | Email / Slack | "Steam active during stop OF#456 — 18€/h" |
| TRS below threshold | Email | "TRS Line 1 below 70% for 2h" |
| Schedule gap >15% | Email | "OF#789 exceeds schedule time — margin eroded" |

### Module 8 — Tribal Knowledge (V2)
**Differentiation: Strong** — Field knowledge structured before it disappears.

| Version | Mechanism |
|---------|-----------|
| V0 | Pre-filled dropdown in dashboard (Jam, Air Pressure, Series Change, Material Wait, Adjustment, Other) — 1 click |
| V1 | Dropdown + free text field saved to local KG |
| V2 | Local SLM chatbot (Ollama/Phi-3) — natural language, structured response, indexed in Knowledge Graph |

---

## 9. Use Cases

### Use Case 1 — Micro-Stops + OT/IT Reconciliation (POC, Priority #1)

**Required tags (OPC-UA via SCADA preferred):**
```
Etat_Machine          (Run / Stop / Setup / Alarm)
Compteur_Pieces       (cumulative integer)
Vitesse_Instantanee   (analog)
Capteur_Bourrage      (boolean)
Pression_Air          (analog, bar)
```

**Analysis pipeline:**
```
STEP 1 — Detection
Run → Stop → Run with 30s < duration < 3min → "Micro-Stop" cumulated

STEP 2 — Causality
IF Capteur_Bourrage = 1      → "Jam"
IF Pression_Air < 4 bars     → "Air Pressure"
ELSE                         → dropdown V0 / text V1 / chatbot V2

STEP 3 — OT/IT Reconciliation (Fuzzy Join)
Physical start detected on Edge → match to ERP production order ±10 min window
Every micro-stop tagged: production order, product, shift, planned schedule

STEP 4 — Real TRS
Real_Availability = (Planned_Time - Major_Stops - Micro-Stops) / Planned_Time
→ "Your real TRS is 74%, not 88%"

STEP 5 — Gain potential
→ "Eliminating 80% of jams recovers 1h04/week = +5 TRS points = +X€/week"
```

**POC success criterion:**
> In under 48h after installation, the Plant Manager knows exactly what micro-stops are costing this week, to the euro, on which production orders, with the causes — automatically.

---

### Use Case 2 — Energy Waste (POC, Priority #2)

**Additional tags:**
```
Puissance_Electrique   (kW, Modbus TCP if not in SCADA)
Debit_Air_Comprime     (m³/h)
Debit_Vapeur           (kg/h)
```

**Rules:**
```
Off-prod waste : IF Etat_Machine=Stop AND Energy > Threshold → ALERT
Cost per batch : Consumption during production order × Energy_Price → batch cost
```

**Fast ROI argument:** Energy waste is visible within the first week, requires no ERP connection, and typically represents 10-15% reducible cost immediately.

---

### Use Case 3 — Schedule Gaps & Margin (V1)

```
IF Actual_Duration_OT > ERP_Schedule_Time + 15% → Alert "Margin Eroded"
Loss = Gap_minutes × Line_Hourly_Cost
```

**Output:** *"Production order #123 (Product A) lost 450€ — 12 micro-stops + 50€ steam waste. Action: optimize changeover for this product."*

---

## 10. Full Roadmap

### Phase 0 — Foundations (Weeks 1-2)

**Objective:** Development infrastructure in place. No features yet, but everything ready to build fast.

| Task | Priority |
|------|----------|
| Create GitHub repos: `mindset-data-edge`, `mindset-data-platform`, `mindset-data-website` (all private) | P0 |
| Go project structure: `cmd/`, `internal/`, `config/`, `Dockerfile` | P0 |
| `go mod init` + core dependencies (gopcua, goburrow/modbus, gos7, paho) | P0 |
| Install Prosys OPC-UA Simulator (local test without factory) | P0 |
| Install Claude Code: `npm install -g @anthropic-ai/claude-code` | P0 |
| Docker Compose local dev stack (Edge Agent + OPC-UA simulator + API + PostgreSQL + Redis) | P0 |
| Apache 2.0 LICENSE file in all repos | P0 |
| Create `docs/mindset.md`, `docs/decisions.md`, `docs/context_starter.md` | P0 |

**Exit criterion:** `go run cmd/agent/main.go` connects to Prosys simulator and prints node tree in terminal.

---

### POC — Micro-Stops + Energy + OT/IT Reconciliation (Weeks 3-10)

**Objective:** Plant Manager sees in under 48h what their micro-stops cost this week, to the euro, on which production orders, with causes — automatically.

#### Sprint 1 — Connection & Discovery (Weeks 3-4)

| Module | Task | Target file |
|--------|------|-------------|
| Network scanner | Subnet scan, OPC-UA/Modbus/S7/MQTT port detection | `internal/discovery/network.go` |
| OPC-UA connector | Connection, node tree browse, tag extraction | `internal/discovery/opcua.go` |
| Modbus connector | TCP connection, register scan, device fingerprint | `internal/discovery/modbus.go` |
| S7 connector | Siemens S7 direct connection (gos7) | `internal/discovery/s7.go` |
| Tag filter | Remove constants, duplicates, signals >10Hz | `internal/discovery/filter.go` |
| SLM classifier | Phi-3 local via Ollama — semantic classification | `internal/classifier/slm.go` |
| Behavioral inference | Live pattern matching 10-15 min | `internal/classifier/behavioral.go` |
| UNS mapper | Tags → ISA-95 hierarchy | `internal/uns/mapper.go` |
| Config loader | YAML config (site params, endpoints, thresholds) | `internal/config/config.go` |

**Exit criterion:** Agent auto-identifies `Etat_Machine`, `Compteur_Pieces`, `Vitesse` on simulator without manual input.

#### Sprint 2 — Rules Engine & Cost Model (Weeks 5-6)

| Module | Task | Target file |
|--------|------|-------------|
| Rules engine | Deterministic engine core | `internal/rules/engine.go` |
| Micro-stop logic | Run→Stop→Run pattern, configurable thresholds | `internal/rules/microstop.go` |
| Energy rules | Off-prod waste detection | `internal/rules/energy.go` |
| Causality | Tag correlation at stop timestamp | `internal/rules/causality.go` |
| Cost model V0 | Manual 3-field entry, € calculation at Edge | `internal/cost/model.go` |
| TRS calculator | Real_Availability, OEE, gain potential | `internal/cost/trs.go` |
| SQLite storage | Local ring buffer 7-15 days | `internal/storage/sqlite.go` |

**Exit criterion:** `"Micro-stop Line1 — 47s — Cause: Jam — Cost: 18€"` output in terminal.

#### Sprint 3 — Fuzzy Join OT/IT (Weeks 7-8)

| Module | Task | Target file |
|--------|------|-------------|
| SQL connector | PostgreSQL, MSSQL, Oracle via YAML config | `internal/connectors/sql.go` |
| REST connector | Modern ERP polling, configurable | `internal/connectors/rest.go` |
| Fuzzy Join engine | Sliding window ±10 min OT/IT alignment | `internal/fuzzy/join.go` |
| Production order matching | Edge events → ERP production order | `internal/fuzzy/matcher.go` |
| Schedule gap detection | Actual_OT vs ERP_Schedule > 15% | `internal/rules/schedule.go` |
| Automated cost model | Import costs from ERP | `internal/cost/erp.go` |

**Exit criterion:** Every micro-stop tagged with correct production order, product, and planned schedule automatically.

#### Sprint 4 — Push, Security & Local Dashboard (Weeks 9-10)

| Module | Task | Target file |
|--------|------|-------------|
| HTTPS Push | Send events to cloud API | `internal/push/cloud.go` |
| TLS 1.3 + mTLS | Encryption + mutual auth | `internal/push/tls.go` |
| API Key auth | Per-site unique key | `internal/push/auth.go` |
| Offline queue | Queue locally if cloud unreachable, auto-sync | `internal/push/queue.go` |
| Historian push | Push enriched events to client historian | `internal/push/historian.go` |
| Cloud receiver | API endpoint that receives events (Go) | `api/handlers/events.go` |
| Local dashboard | React app served from Edge Agent | `frontend/src/` |
| Gantt timeline | Run/Stop/Setup timeline | `frontend/src/components/GanttTimeline.jsx` |
| Pareto chart | Causes by € impact | `frontend/src/components/ParetoChart.jsx` |
| ROI simulator | Real TRS vs declared + gain potential | `frontend/src/components/ROISimulator.jsx` |
| Cause dropdown | V0 tribal knowledge (1-click) | `frontend/src/components/CauseDropdown.jsx` |
| Onboarding wizard | 3-field cost entry | `frontend/src/pages/Onboarding.jsx` |
| Local SMTP/Slack | Direct alerting from Edge | `internal/alerting/smtp.go` |
| Docker image | Build + push `mindsetdata/edge-agent:v0` | `Dockerfile` |

**Full POC exit criterion:**
> Docker install in 1 command → automatic OPC-UA connection → micro-stop detection → OT/IT reconciliation → energy waste → dashboard with € losses by production order — in under 48h on site.

---

### Website & Distribution (Weeks 5-8, parallel)

| Task | Priority |
|------|----------|
| Next.js setup + deploy to Vercel (mindsetdata.io) | P0 |
| Landing page (pitch, value prop, CTA) | P0 |
| `/product` page (Connect / Contextualise / Decide / Act) | P0 |
| `/use-cases` page (micro-stops + energy with € numbers) | P0 |
| `/security` page (push-only, mTLS, RGPD/NIS2) | P0 |
| `/download` page (form → API key → Docker command) | P0 |
| `/contact` page (demo request) | P1 |
| Docker Hub `mindsetdata/edge-agent` public image | P0 |
| Email onboarding sequence (API key + 5-min setup guide) | P1 |

**Client install flow:**
```powershell
# Windows
docker pull mindsetdata/edge-agent:latest
docker run -d `
  -e API_KEY=their-unique-key `
  -e SITE_NAME="Usine Paris Nord" `
  --network host `
  mindsetdata/edge-agent:latest
# → Agent starts, scans network, finds OPC-UA equipment
# → Local dashboard available at http://localhost:3000 immediately
# → Cloud dashboard available at app.mindsetdata.io within 1h
```

---

### V1 — Multi-site + Deep ERP + Historian (Weeks 11-18)

| Module | Description | Priority |
|--------|-------------|----------|
| Multi-site aggregation | Cross-site KG, inter-site benchmarking | P0 |
| Files + FTP/SFTP connector | CSV/Excel/JSON import from network shares | P0 |
| Historian connector (PI/Wonderware/InfluxDB) | Push enriched events to client's historian | P0 |
| MQTT connector | Modern IIoT gateways | P1 |
| Sparkplug B | Native ISA-95 MQTT payload | P1 |
| InfluxDB connector | Digitalized ETI historians | P1 |
| Multi-site dashboard V1 | Cross-site Pareto, site comparison | P0 |
| Microsoft Teams alerting | Microsoft-centric ETI | P1 |
| BYOC deployment | Docker-compose for on-premise cloud | P1 |
| Auto-update Edge Agent | Automatic update from Docker Hub | P1 |

**V1 success criterion:**
> Every completed production order automatically shows: actual time vs ERP schedule, exact energy cost, realized margin vs theoretical margin. Across all connected sites.

---

### V2 — Open UNS + AI Agents + Tribal Knowledge (Weeks 19-32)

| Module | Description | Priority |
|--------|-------------|----------|
| Semantic UNS API | REST/GraphQL on Knowledge Graph for AI agents | P0 |
| Tribal Knowledge chatbot | Ollama/Phi-3 local — natural language cause capture | P0 |
| Native AI agents | Automatic diagnosis, action suggestion | P1 |
| MTConnect | CNC/machining/metallurgy | P1 |
| BACnet/IP | Building energy management | P2 |
| Omron FINS | Agrifood/pharma niche | P2 |
| MongoDB connector | Modern MES | P2 |
| Partner SDK | Integrators and ISV can connect to UNS | P2 |
| Functions marketplace | No-code function library by industry | P2 |
| Predictive model | ML on historical patterns → failure prediction | P2 |

**V2 success criterion:**
> A third-party AI agent queries factory context and produces reliable recommendations without hallucination, using the UNS as source of truth.

---

### V3+ — European OS Infrastructure (Month 9+)

| Axis | Description |
|------|-------------|
| UNS as a Service | Open standard for ETI manufacturing data infrastructure |
| Ecosystem | Integrators, ISV, AI agents as partners on the UNS |
| Cloud connectors | Kafka, AWS S3, Azure Blob, Google Pub/Sub, NATS |
| Data Act positioning | EU Data Act compliance infrastructure |
| Cross-industry KG | Cross-site benchmarks → sectoral reference data |
| Acquisition moat | Each site = irreplaceable cumulative context → structurally impossible churn |

---

## 11. GitHub Repository Structure

### Repositories

```
github.com/mindset-data/
├── mindset-data-edge          ← Edge Agent Go (private)
├── mindset-data-platform      ← Cloud API + React dashboard (private)
└── mindset-data-website       ← Next.js marketing site (private)
```

### mindset-data-edge (complete structure)

```
mindset-data-edge/
├── cmd/
│   └── agent/
│       └── main.go                    ← Entry point
├── internal/
│   ├── discovery/
│   │   ├── network.go                 ← Subnet scan, port detection
│   │   ├── opcua.go                   ← OPC-UA connection + node tree
│   │   ├── modbus.go                  ← Modbus TCP + device fingerprint
│   │   ├── s7.go                      ← Siemens S7 direct connection
│   │   └── filter.go                  ← Noise/constant removal
│   ├── classifier/
│   │   ├── slm.go                     ← Phi-3 local via Ollama
│   │   └── behavioral.go              ← Live behavioral pattern matching
│   ├── uns/
│   │   ├── graph.go                   ← Local in-memory Knowledge Graph
│   │   ├── isa95.go                   ← ISA-95 ontology definitions
│   │   └── mapper.go                  ← Tag → UNS topic mapping
│   ├── rules/
│   │   ├── engine.go                  ← Deterministic rules engine core
│   │   ├── microstop.go               ← Micro-stop detection logic
│   │   ├── energy.go                  ← Energy waste detection
│   │   ├── causality.go               ← Tag correlation at stop timestamp
│   │   └── schedule.go                ← ERP schedule gap detection
│   ├── fuzzy/
│   │   ├── join.go                    ← Fuzzy Join temporal engine
│   │   └── matcher.go                 ← OT event → ERP production order
│   ├── connectors/
│   │   ├── sql.go                     ← PostgreSQL/MSSQL/Oracle
│   │   ├── rest.go                    ← REST API polling
│   │   ├── files.go                   ← CSV/Excel/JSON import
│   │   └── ftp.go                     ← FTP/SFTP
│   ├── cost/
│   │   ├── model.go                   ← Cost calculation (€) at Edge
│   │   ├── trs.go                     ← Real TRS, OEE, gain potential
│   │   └── erp.go                     ← Automated cost from ERP
│   ├── storage/
│   │   ├── sqlite.go                  ← Local ring buffer 7-15 days
│   │   └── purge.go                   ← TTL auto-purge
│   ├── push/
│   │   ├── cloud.go                   ← HTTPS event push to cloud
│   │   ├── tls.go                     ← TLS 1.3 + mTLS
│   │   ├── auth.go                    ← API key per site
│   │   ├── queue.go                   ← Offline queue + auto-sync
│   │   └── historian.go               ← Push to client historian
│   ├── alerting/
│   │   ├── smtp.go                    ← Direct email from Edge
│   │   └── slack.go                   ← Direct Slack webhook from Edge
│   └── config/
│       └── config.go                  ← YAML config loader
├── frontend/                          ← React local dashboard
│   ├── src/
│   │   ├── components/
│   │   │   ├── GanttTimeline.jsx      ← Run/Stop/Setup timeline
│   │   │   ├── ParetoChart.jsx        ← Causes by € impact
│   │   │   ├── ROISimulator.jsx       ← TRS real vs declared
│   │   │   └── CauseDropdown.jsx      ← V0 tribal knowledge
│   │   ├── pages/
│   │   │   ├── Realtime.jsx           ← Tab 1: Real time
│   │   │   ├── MicroStops.jsx         ← Tab 2: Micro-stops analysis
│   │   │   ├── Impact.jsx             ← Tab 3: ROI & impact
│   │   │   └── Onboarding.jsx         ← 3-field cost wizard
│   │   └── App.jsx
│   └── package.json
├── config/
│   └── agent.yaml                     ← Default config template
├── docs/
│   ├── mindset.md                     ← This file
│   ├── decisions.md                   ← Technical decisions log
│   └── context_starter.md             ← Claude session brief
├── .github/
│   └── workflows/
│       └── ci.yml                     ← GitHub Actions CI/CD
├── Dockerfile
├── docker-compose.yml                 ← Local dev stack
├── go.mod
├── go.sum
├── LICENSE                            ← Apache 2.0
└── README.md
```

### mindset-data-platform (complete structure)

```
mindset-data-platform/
├── api/
│   ├── handlers/
│   │   ├── events.go                  ← Receive events from Edge
│   │   ├── sites.go                   ← Site management
│   │   ├── dashboard.go               ← Dashboard data API
│   │   └── auth.go                    ← API key validation
│   ├── middleware/
│   │   ├── auth.go                    ← Auth middleware
│   │   └── cors.go                    ← CORS
│   └── main.go
├── internal/
│   ├── kg/
│   │   ├── graph.go                   ← Persistent Knowledge Graph
│   │   └── aggregator.go              ← Cross-site KG aggregation
│   ├── storage/
│   │   ├── postgres.go                ← PostgreSQL + TimescaleDB
│   │   └── redis.go                   ← Real-time cache
│   └── alerting/
│       ├── smtp.go                    ← Cloud email relay
│       ├── slack.go                   ← Cloud Slack relay
│       └── teams.go                   ← Microsoft Teams (V1)
├── frontend/                          ← Cloud dashboard (remote access)
│   ├── src/
│   │   ├── components/                ← Same components as Edge dashboard
│   │   └── pages/
│   │       ├── MultiSite.jsx          ← V1: cross-site view
│   │       └── Financial.jsx          ← CFO/CEO dashboard
│   └── package.json
├── pipelines/
│   ├── microstop.yaml                 ← Redpanda Connect pipeline
│   └── energy.yaml
├── migrations/                        ← PostgreSQL schema migrations
│   └── 001_init.sql
├── docker-compose.yml
├── .github/
│   └── workflows/
│       └── ci.yml
├── LICENSE
└── README.md
```

### mindset-data-website (complete structure)

```
mindset-data-website/
├── app/
│   ├── page.jsx                       ← Landing page
│   ├── product/page.jsx               ← How it works
│   ├── use-cases/page.jsx             ← Micro-stops + energy with €
│   ├── security/page.jsx              ← Push-only, mTLS, RGPD/NIS2
│   ├── download/page.jsx              ← Form → API key → Docker command
│   └── contact/page.jsx               ← Demo request
├── components/
│   ├── Hero.jsx
│   ├── HowItWorks.jsx
│   ├── UseCases.jsx
│   └── DownloadForm.jsx
├── public/
├── .github/
│   └── workflows/
│       └── deploy.yml                 ← Auto-deploy to Vercel on push
├── LICENSE
└── README.md
```

### GitHub Actions CI/CD (all repos)

```yaml
# .github/workflows/ci.yml (edge-agent)
name: CI/CD

on:
  push:
    branches: [main, dev]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go test ./...

  build-push:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          push: true
          tags: mindsetdata/edge-agent:latest
```

---

## 12. Infrastructure & Distribution

### Production endpoints

```
mindsetdata.io              ← Marketing (Next.js / Vercel — free)
app.mindsetdata.io          ← Cloud dashboard (React / Scaleway)
api.mindsetdata.io          ← Cloud API receiver (Go / Scaleway)
hub.docker.com/mindsetdata  ← Edge Agent Docker image (free public)
```

### Scaleway FR — V0 minimal cost

```
1× PLAY2-NANO VPS (2 vCPU, 2GB RAM)   → 3.99€/month
1× Managed PostgreSQL DEV-1500         → 9.99€/month
1× Object Storage (backups)            → pay per use ~1€/month
──────────────────────────────────────────────────────
Total V0: ~15€/month
```

### Local dev docker-compose

```yaml
version: '3.8'
services:
  edge-agent:
    build: .
    environment:
      - API_KEY=dev-key-local
      - SITE_NAME=Local Test Site
      - CLOUD_ENDPOINT=http://api:8080
      - OPC_UA_ENDPOINT=opc.tcp://host.docker.internal:4840
    network_mode: host

  api:
    build: ./mindset-data-platform/api
    ports:
      - "8080:8080"
    environment:
      - DB_URL=postgres://mindset:mindset@postgres:5432/mindset
      - REDIS_URL=redis:6379
    depends_on:
      - postgres
      - redis

  postgres:
    image: timescale/timescaledb:latest-pg15
    environment:
      - POSTGRES_DB=mindset
      - POSTGRES_USER=mindset
      - POSTGRES_PASSWORD=mindset
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  redpanda-connect:
    image: redpandadata/connect:latest
    volumes:
      - ./mindset-data-platform/pipelines:/pipelines
    command: run /pipelines/microstop.yaml

volumes:
  postgres_data:
```

---

## 13. Tech Stack

| Layer | Technology | Version | Justification |
|-------|------------|---------|---------------|
| Edge Agent | Go + Docker | Go 1.22+ | Lightweight, static binary, cross-platform, low footprint |
| OPC-UA | `gopcua` | latest | Mature, browse + subscriptions + security |
| Modbus | `goburrow/modbus` | latest | TCP + RTU |
| Siemens S7 | `gos7` | latest | Direct PLC Siemens access |
| MQTT | `paho.mqtt.golang` | latest | Eclipse standard |
| Sparkplug B | `sparkplug-go` | latest | Structured industrial payload |
| SQL | `database/sql` + drivers | — | pgx, go-mssqldb, go-ora |
| Local SLM | Phi-3 via Ollama | latest | Sovereign, zero latency, zero cloud |
| Local DB | SQLite | — | Zero deps, ring buffer |
| Pipeline | Redpanda Connect | latest | ex-Benthos, declarative YAML |
| Ontology | ISA-95 | standard | Maximum interoperability |
| Encryption | mTLS + TLS 1.3 | — | NIS2 compliant |
| Cloud DB | PostgreSQL + TimescaleDB | pg15 | Time-series + events |
| Cache | Redis | 7 | Real-time dashboard |
| Cloud | Scaleway / OVH | — | FR sovereignty, RGPD |
| Frontend | React | 18+ | Local + cloud dashboard |
| Website | Next.js | 14+ | SEO, Vercel deploy |
| CI/CD | GitHub Actions | — | Free, build + test + push Docker |
| License | Apache 2.0 | — | Open, permissive |

---

## 14. Security & Compliance

| Principle | Implementation |
|-----------|---------------|
| Read-only | Zero writes to PLC, SCADA, ERP |
| Push-only | Outbound HTTPS only, zero inbound open ports |
| Data minimisation | Only transformed events leave the Edge, never raw data |
| Transit encryption | mTLS (mutual authentication) + TLS 1.3 |
| Edge processing | Raw data never leaves client network |
| Offline resilience | Full operation if cloud unreachable |
| Sovereignty | FR certified cloud or BYOC (zero data outside territory) |
| RGPD | Native compliance — industrial data, minimisation, FR location |
| NIS2 | OT/IT compliant architecture, audit trail |
| Audit trail | Every event timestamped and signed |

---

## 15. Tech Moat

### The 4 defensive elements

**1. Zero-manual auto-discovery**
Network scan + behavioral inference + SLM = automatic connection without mapping. 48h vs 3 months for an integrator. Structural advantage on client acquisition cost.

**2. The Fuzzy Join**
OT/IT temporal alignment — the universal unsolved problem in mid-market. The sliding window algorithm that re-attaches the ERP clock and the machine clock is the most defensive component in the stack: hard to build, invisible from the outside.

**3. The cumulative Knowledge Graph (site fingerprint)**
Site-specific non-replicable context. Replacing Mindset Data = losing all accumulated intelligence. Churn structurally becomes prohibitive.

**4. Tribal Knowledge structured over time (V2)**
`sensor pattern → operator label` associations: impossible to reconstruct without access to the same site in real time. No competitor can copy this dataset.

### Value curve over time

```
Day 0       → Base UNS, generic rules, first micro-stops detected
Week 4      → Labeled causes, site-specific patterns, real TRS calculated
Month 3     → Product/machine correlations, OF history, Fuzzy Join calibrated
Month 6     → Site-specific predictive model, structured tribal knowledge
Month 12+   → Cross-site benchmarks, reliable AI agents
            → Churn = irreversible destruction of accumulated intelligence
```

---

## 16. Development Workflow

### Claude session protocol

```
START OF SESSION
1. Open claude.ai (architecture/decisions) or Claude Code (active coding)
2. Paste context_starter.md
3. Describe session objective

DURING SESSION
- Architecture / decisions   → claude.ai (this chat)
- Active code in repo        → Claude Code (terminal in repo folder)
- Debug                      → paste code + error + context

END OF SESSION
- Ask: "Summarize technical decisions made today"
- Update decisions.md
- Update context_starter.md (section "Current stage")
- Commit: git commit -m "session: [topic]"
```

### Windows setup (one-time)

```powershell
# 1. Git → git-scm.com/download/win
# 2. Go  → go.dev/dl (go1.22.x.windows-amd64.msi)
# 3. Docker Desktop → docs.docker.com/desktop/install/windows-install
# 4. Node.js LTS → nodejs.org
# 5. VS Code → code.visualstudio.com

# Then in PowerShell:
npm install -g @anthropic-ai/claude-code

git clone https://github.com/mindset-data/mindset-data-edge
cd mindset-data-edge
claude   # ← Claude Code takes over from here
```

### context_starter.md (update every session)

```markdown
# Mindset Data — Context Brief

## What we're building
Industrial data infrastructure for manufacturing ETI.
Edge Agent (Go) → UNS (ISA-95) → Decision dashboard.
Maximum local processing. Zero raw data to cloud.
Zero hardware, zero dev client, deployment in 48h.

## Team
- Mohamed (CTO) — Polytechnique, IoT & embedded systems (Windows PC)
- Cécilia (CEO) — EDHEC, ex-VC AgriFoodTech

## Stack
- Edge Agent: Go + Docker
- V0 protocols: OPC-UA (gopcua) + Modbus (goburrow/modbus)
- V1 protocols: Siemens S7 (gos7), SQL, REST, Files/FTP
- Local SLM: Phi-3 / Ollama
- Local DB: SQLite (7-15 day ring buffer)
- Pipeline: Redpanda Connect
- Cloud: Scaleway FR (minimal — KG backup + remote dashboard)
- Frontend: React (local first) + Next.js (website)
- License: Apache 2.0

## Current stage — [UPDATE EACH SESSION]
Working on: _______________
Last decision: _______________
Current blocker: _______________

## POC scope (active)
Micro-stops + OT/IT reconciliation + Energy.
Tags: Etat_Machine, Compteur_Pieces, Vitesse, Capteur_Bourrage, Pression_Air.
Cost model: manual 3-field entry V0, auto ERP V1.
Fuzzy Join: ±10 min sliding window.
Dashboard: local React served from Edge (localhost:3000).
Push: HTTPS TLS 1.3 + mTLS, compressed snapshots only.
Storage: SQLite 7-15 days local + push to client historian.
```

### Parallel development sessions

| Session | Track | Objective | Expected output |
|---------|-------|-----------|-----------------|
| **1** | Edge | Go skeleton + OPC-UA connect | `go run` → node tree in terminal |
| **2** | Platform | Go API + PostgreSQL schema + Docker Compose | Full local stack running |
| **3** | Website | Next.js + Vercel deploy | mindsetdata.io live |
| **4** | Edge | Behavioral inference + SLM classification | Agent auto-identifies Etat_Machine |
| **5** | Dashboard | React skeleton + Gantt component | Local dashboard at localhost:3000 |
| **6** | Edge | Rules engine — micro-stops + energy | `"Micro-stop 47s — Jam — 18€"` |
| **7** | Edge | Fuzzy Join + automated cost model | Every stop tagged with production order |
| **8** | Dashboard | Pareto €, ROI simulator, energy tab | Plant Manager reads week's losses |
| **9** | Edge | SQLite ring buffer + historian push | Local storage + client historian fed |
| **10** | All | GitHub Actions CI/CD | Auto-deploy on every push to main |

---

## 17. Hardware Requirements

### Client machine for Edge Agent

| Spec | Minimum (V0) | Recommended (V0) | Ideal (V1+) |
|------|-------------|-----------------|-------------|
| CPU | 4 cores 2.0 GHz | 8 cores 2.5 GHz | 8 cores 3.0 GHz |
| RAM | 8 GB | 16 GB | 32 GB |
| Storage | 50 GB SSD | 100 GB SSD | 256 GB NVMe |
| Network | 1 NIC + VLAN | 2 NIC (OT + IT) | 2 NIC mandatory |
| GPU | ✗ | ✗ | Optional RTX 3060 |
| OS | Windows 10 / Ubuntu 20.04 | Ubuntu 22.04 LTS | Ubuntu 22.04 LTS |
| Docker | Required | Required | Required |
| Internet | Outbound HTTPS only | Outbound HTTPS only | Outbound HTTPS only |

### Phi-3 local SLM constraints

```
Model size on disk   : 2.2 GB (Q4 quantized)
RAM needed to run    : 2.5 GB
Inference speed:
  CPU 4 cores        → 8-15 tokens/sec (acceptable for tag classification)
  CPU 8 cores        → 15-25 tokens/sec (good)
  GPU RTX 3060       → 60-80 tokens/sec (fast)

Tag classification at install:
  500 tags to classify → ~15-25 min on 4-core CPU
  Done once at setup, not in real time → perfectly acceptable
```

### Network requirement — 2 NICs preferred

```
OT Network (isolated)          IT Network (connected)
─────────────────────          ──────────────────────
PLC / SCADA / Sensors          ERP / MES / Internet
Read-only access               Cloud sync (HTTPS out)
No internet                    SQL queries (V1)
NIC 1 (e.g. 192.168.10.x)    NIC 2 (e.g. 192.168.1.x)
         │                              │
         └──────── EDGE AGENT ──────────┘
                  (bridges both,
                  read-only OT,
                  push-only IT)
```

### What to tell the client (sales conversation)

> *"To install our agent, you need an existing PC or server in your factory network — it can be a desktop PC you already have, a spare server, or a small industrial box if you prefer something dedicated. You need minimum 8GB RAM, 50GB free storage, and ideally two network ports — one on the machine side, one on the IT side. We install everything with a single Docker command. No modifications to your infrastructure."*

### Hardware decision tree for onboarding

```
Client has existing PC/server?
        │
        ├── YES → Check: RAM ≥ 8GB + Storage ≥ 50GB?
        │              │
        │              ├── YES → Use it. Install Docker. Done.
        │              │
        │              └── NO  → Upgrade RAM (50-100€) or
        │                        use recommended mini-PC below
        │
        └── NO  → Recommend:
                  Budget    : Beelink Mini PC (i5, 16GB, 500GB) — 350€
                  Standard  : Lenovo ThinkCentre M90q — 600€
                  Industrial: Advantech UNO-2484G — 1,200€
```

---

*Internal document — Mindset Data — Confidential*
*Last updated: May 2026*
*License: Apache 2.0*
