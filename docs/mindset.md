# MINDSET DATA — Project Documentation

> **Vision:** Connect every machine, system, and data flow in a factory to a single reliable and exploitable source of truth — the Unified Namespace (UNS) — to transform every signal into a business decision.

> **Document status:** Vision narrative — updated 2026-06-28. **Reality-check note added 2026-07-28 — read it before trusting anything below as current state.**
> **Canonical decision log:** `docs/decisions.md` — restored 2026-07-28 from git history (was deleted in commit `c310c6f`, recovered from `c25337d`). Back and readable; still untracked in git until committed.
> **Competitive analysis:** `docs/MindSet_Competitive_Analysis_v2_3.xlsx` (the `v2_2` file this used to point to is also deleted).
> **Audit trail:** `docs/analysis_log.md` — now well past Entries 1-16, currently up to **Entry 138**. Search it for the real, current status of anything below rather than trusting a specific section number here.

---

## ⚠️ Reality check (2026-07-28) — read this before the rest of the document

**This document is a vision/roadmap narrative, not a description of the current build.** A large fraction of what follows — separate `mindset-data-platform` and `mindset-data-website` repos, Redpanda Connect pipelines, a local Phi-3/Ollama SLM, a PostgreSQL+TimescaleDB+Redis cloud stack, the 3-edition Cloud/Hybrid/On-Premise distribution model, license-key-gated Docker Hub distribution — **does not exist in the actual codebase today**.

**License, resolved**: §13 (proprietary/closed-source) is correct; the footer (Apache 2.0) was stale. `docs/decisions.md` is restored (see above) — its locked decision: **"Licensing model: PROPRIETARY (closed-source) for first 2 years — supersedes prior Apache 2.0 decision."** Apache 2.0 was the *original* choice, later explicitly overridden. Even `decisions.md`'s own footer had this same stale-Apache-2.0 bug before the restore — a pre-existing doc bug, now fixed there too. No real `LICENSE` file exists in the repo yet either way — that's still a genuine gap.

**For an accurate, actively-maintained picture of what's actually built, use `CLAUDE.md`** (repo root) — it's auto-loaded into every Claude Code session and documents the real architecture: a single repo, two Go binaries (`cmd/server`, `cmd/agent`), a custom pipeline engine (not Redpanda), a SQLite-backed Knowledge Graph with confidence-gated OT+IT structural auto-bootstrap, an MCP server (HTTP + stdio, 5 tools), and a React frontend — all running on one machine today, not the distributed multi-edition system described in §4/§11/§12 below. `docs/context_starter.md` (rewritten 2026-07-28) has a compact built-vs-not-built summary if you want it without opening `CLAUDE.md`.

**What's genuinely still true and load-bearing from this doc**: the problem framing (§1), the ISA-95/UNS positioning (§2, §7), the persona segmentation (§3) — modulo the unsourced McKinsey quote and two persona "Verbatim" quotes flagged in `analysis_log.md` Entry 137, which need a real source or a rewrite before external use. Treat everything else — especially §4 (Architecture), §11 (repo structure), §12 (infrastructure/distribution), §13 (tech stack) — as a **future-state proposal**, not a status report.

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

**Target market**: 15,000+ European mid-sized factories. **Initial GTM focus on 4 high-value verticals** — pharma 💊 · cosmetics 💄 · agrifood 🌾 · metallurgy ⚙️ — where downtime cost is highest, regulation is strictest, and the EU-sovereignty pitch lands hardest. See §3 for vertical detail.

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

### 3.0 Target market + initial verticals (V1 GTM focus)

**TAM**: 15,000+ European mid-sized factories.

**Initial GTM focus — 4 high-value verticals**:

| Vertical | Why it fits MindSet | Sales motion | Indicative deal size |
|---|---|---|---|
| 💊 **Pharma** | GMP / FDA / EMA regulated → sovereignty pitch lands hard · ~50k€/h downtime cost · mature OEE culture · willingness to pay | Enterprise IT-led, 6-12 month cycle, RFP, ISO 27001 + GAMP 5 required | 50-150k€ / site / year |
| 💄 **Cosmetics** | EU Cosmetic Regulation = procurement-sensitive · high-margin products · brand-reputation-sensitive (no public IT incidents) · often part of large groups (LVMH, L'Oréal, Estée Lauder) | Enterprise IT-led, similar to pharma | 50-100k€ / site / year |
| 🌾 **Agrifood** | Largest FR industrial vertical · high cost of waste + energy-intensive · strict EU/FR regulation (HACCP, traceability) · many independent mid-sized + family-owned | **Self-serve Plant-Manager <30k€** (the original ETI motion) | 15-30k€ / site / year |
| ⚙️ **Metallurgy** | Capital-intensive · high downtime cost · complex OF/scheduling · energy-heavy (Level 2 energy waste demos perfectly) · mix of independent + grouped | Mixed — self-serve for independents, enterprise for grouped | 20-50k€ / site / year |

**Geographic execution**: starts in **France** (founders' geography + Boost10x network), expands to **DACH + Italy + Spain + Nordics** in V2-V3.

**Two parallel sales motions** (downstream of the vertical mix):
- **Motion #1 — Self-serve Plant Manager** (agrifood + independent metallurgy): Docker pull, 48h deployment, <30k€/site, Plant Manager signs autonomously
- **Motion #2 — Enterprise IT-led** (pharma + cosmetics + grouped metallurgy): RFP-driven, 6-12 month cycles, 50-150k€/site, requires ISO 27001 + GAMP 5 + RBAC + audit log

The buyer personas below (IT/OT Manager, Operations Director, Plant Manager, CFO/CEO) apply across both motions — but the **decision authority** shifts: Plant Manager-led in Motion #1; CISO + IT Procurement + Plant Manager + Ops Director in Motion #2.

---

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

### 4.1 Three deployment editions

The product ships in **exactly three editions**. No hyperscaler edition through 2029 (reconsidered for international scaling — see `decisions.md`).

| Edition | Cloud component | Multi-site | Remote dashboard | Target customer |
|---|---|---|---|---|
| **On-Premise** | NONE — zero cloud | No (per-site only) | No (factory LAN only) | Defense, public sector, sensitive pharma |
| **Hybrid** *(default)* | Scaleway FR / OVH FR (MindSet-managed) | Yes | Yes (`app.mindsetdata.io`) | Commercial ETI — the everyday offer |
| **Self-Hosted** | Customer's EU-jurisdiction cloud (Hetzner, IONOS, T-Systems, 3DS Outscale, Bleu) OR customer's on-prem Kubernetes | Yes | Yes (customer-hosted) | Large multi-site with existing EU cloud relationship |

**Explicitly NOT offered:** AWS, Azure, GCP — including their EU regions. The US CLOUD Act exposure invalidates the sovereignty pitch for public sector and defense customers.

### 4.2 Global overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         CLIENT NETWORK (OT + IT)                    │
│                                                                     │
│  OT SOURCES                      IT SOURCES                        │
│  ├── SCADA / PLC (OPC-UA)        ├── ERP (SQL: PostgreSQL/         │
│  ├── Siemens PLC (S7, V1.5)      │    MSSQL/MySQL — V1)            │
│  ├── Legacy PLC (Modbus TCP)     ├── ERP REST (V1.5)               │
│  └── Energy meters (Modbus TCP)  └── Files / FTP (V1.5)            │
│           │                               │                        │
│           └───────────────┬───────────────┘                        │
│                           ▼                                        │
│              ┌────────────────────────┐                            │
│              │     EDGE AGENT (Go)    │  ← Docker, client's PC     │
│              │     1 binary           │  ← 8GB RAM min, 16GB rec   │
│              │                        │                            │
│              │  Discovery Layer       │                            │
│              │  ├── Network scanner   │  ← Auto-detect endpoints   │
│              │  ├── OPC-UA browse     │  ← Node tree extraction    │
│              │  ├── Modbus scan       │  ← Register fingerprint    │
│              │  └── S7 scan (V1.5)    │                            │
│              │                        │                            │
│              │  Intelligence Layer    │                            │
│              │  ├── Behavioral infer  │  ← Live pattern matching   │
│              │  ├── SLM Phi-3 local   │  ← Tag classif. (Ollama)   │
│              │  ├── UNS ISA-95 mapper │  ← Tag → Site/Area/WC/Tag  │
│              │  └── OF-state Fuzzy J. │  ← OT events → active OF   │
│              │                        │    (NOT sliding window)    │
│              │                        │                            │
│              │  Processing Layer      │                            │
│              │  ├── Rules engine      │  ← Templates: micro-stop,  │
│              │  │                     │    energy waste, OEE/TRS   │
│              │  └── Cost model (€)    │  ← Real-time € calculation │
│              │                        │                            │
│              │  AI Layer (V1)         │                            │
│              │  ├── MCP server        │  ← localhost:5000          │
│              │  └── Ad-hoc Analyst    │  ← Chat agent in dashboard │
│              │                        │                            │
│              │  Storage Layer         │                            │
│              │  ├── SQLite 7-15 days  │  ← Local ring buffer       │
│              │  ├── KG (domain+tech)  │  ← Site fingerprint        │
│              │  └── Push → Historian  │  ← Client's own system     │
│              │                        │                            │
│              │  Serving Layer         │                            │
│              │  ├── Local dashboard   │  ← localhost:8080 (React)  │
│              │  └── Local alerting    │  ← Direct SMTP/Slack/Teams │
│              └───────────┬────────────┘                            │
│                          │ HTTPS Push-only (mTLS + TLS 1.3)        │
│                          │ Transformed events + KG snapshots ONLY  │
│                          │ Raw OT data NEVER leaves                │
└──────────────────────────┼──────────────────────────────────────────┘
                           ▼
              ┌────────────────────────────────┐
              │ CLOUD TIER (Hybrid + Self-Hosted only — On-Premise has no cloud) │
              │                                                                  │
              │  ├── Cross-site KG aggregation    ← Multi-site only             │
              │  ├── Remote dashboard             ← CEO/CFO/Ops outside factory │
              │  ├── Site management API          ← Auth, API keys, licenses    │
              │  ├── KG snapshot backup           ← Encrypted at edge first     │
              │  └── Heartbeat monitor            ← Detect dead edge agent      │
              │                                                                  │
              │  Hybrid: Scaleway FR / OVH FR (MindSet-hosted, ~15€/mo)         │
              │  Self-Hosted: customer's EU cloud or on-prem K8s                │
              └─────────────────────────────────────────────────────────────────┘
```

### 4.3 Core principles

| Principle | Implementation |
|---|---|
| Zero additional hardware | Edge Agent runs on existing client infrastructure (PC, VM, industrial box) |
| Zero raw data in cloud | Only transformed events + aggregated snapshots go up |
| Push-only | Outbound HTTPS (mTLS + TLS 1.3) only — zero inbound open ports |
| Read-only on source systems | Zero writes to PLC, SCADA, ERP |
| Zero manual connection work | Network scan + behavioral inference + Phi-3 SLM classification |
| Maximum local processing | Rules engine, cost model, OF-state Fuzzy Join, dashboard, MCP server, AI agent — all at the Edge |
| Offline resilience | Full operation if cloud unreachable; queues sync on reconnect |
| EU sovereignty | EU jurisdiction by design — no hyperscaler edition through 2029 |
| Single-vendor, no middleware | No Kepware, no per-tag fees, native protocol drivers in Go |

### 4.4 What runs WHERE — strict rule

**A feature goes to the cloud tier ONLY IF all three hold:**
1. It needs to span multiple sites OR be reachable from outside the factory.
2. Latency tolerates >1s round-trip.
3. Only already-transformed data crosses the boundary (raw OT values forbidden).

Everything else runs at the edge.

**The edge agent encapsulates ~63 components across 11 categories.** Full inventory in `docs/analysis_log.md` Entry 23 + Sheet 8 of `MindSet_Competitive_Analysis_v2_3.xlsx`. Compact summary:

| Category | # | Examples | V1 ship |
|---|---|---|---|
| Storage | 6 | SQLite ring buffer · Domain KG (the moat dataset) · Technical KG · Tag/Topic registries · State tracker | All 6 |
| Message bus | 1 | Local Mosquitto MQTT broker (bundled in docker-compose) | 1 |
| Discovery + Classification | 5 | Network scanner · OPC-UA browse · Modbus fingerprint DB · Behavioral inference · UNS ISA-95 mapper | All 5 |
| Connectors | 11 | OPC-UA, Modbus TCP, SQL multi-dialect (V1) · S7, REST, Files (V1.5) · MQTT, Sparkplug B, MTConnect, BACnet (V2+) | 3 |
| Processing engines | 6 | Pipeline engine · Function registry · Rules engine · OF-state Fuzzy Join · Cost model · OEE/TRS calculator | All 6 |
| KG integration | 3 | KG subscriber · KG builder · KG REST API | All 3 |
| Local UI | 10 | React skeleton · Pipeline Studio · KG viewer · Dashboard + WebSocket hub · Gantt · Pareto · OEE view · ROI sim · Tribal knowledge capture · Onboarding wizard | All 10 |
| AI layer | 4 | Phi-3 runtime · MCP server · Ad-hoc Analyst agent · Remote LLM proxy | All 4 |
| Communication outbound | 5 | WebSocket hub · HTTPS cloud pusher · SMTP/Slack/Teams alerting · Heartbeat sender · Historian push | 4 |
| Infrastructure | 6 | Config loader · Logger · Secrets management · License validator · Health endpoints · Auto-update | 5 |
| Security (pending decision) | 6 | Encryption at rest · Signed binaries · SBOM · Audit log · RBAC · SSO | 4 V1 + 2 V1.5 |
| **TOTAL EDGE** | **63** | | **~51** |

**Components running IN THE CLOUD (Hybrid + Self-Hosted only — On-Premise edition has none):**
Cross-site KG aggregator · multi-site dashboard · single-site remote dashboard · site management API (auth + keys + entitlements) · encrypted KG snapshot backup · heartbeat / liveness monitor.

**Opt-in cloud (customer decision, with explicit disclosure):**
Remote LLM proxy (OpenAI/Claude/Mistral) instead of local Phi-3 · cloud MCP relay for remote AI agents (V1.5+).

**Why this matters for the pitch:** "1 Go binary + 1 React UI" is honest at 30,000ft view — but the binary actually encapsulates ~51 production-grade components at V1, all in one Docker container that installs in 48h. UMH ships ~10 OSS components requiring Kubernetes expertise. That's the deployment-simplicity moat.

### 4.5 Storage strategy

```
EDGE (always)                    CLIENT HISTORIAN          MindSet CLOUD (Hybrid only)
─────────────────────            ─────────────────         ──────────────────────────
SQLite ring buffer 7-15 days     PI System /               Aggregated KG snapshots only
+ persistent KG (site fingerprint)  Wonderware /           Transformed events for cross-site
Raw events + causes + costs €    InfluxDB / MSSQL          NEVER raw OT data
Auto-purge by TTL                Push enriched events      Encrypted before upload
                                 (V1.5 connector)
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

### Module 4 — OF-State-Based Attribution Engine (formerly "Fuzzy Join")
**Differentiation: Strong** — Universal unsolved problem in mid-market. Makes data 100% reliable.

**Important correction (2026-06-28):** the original "sliding window ±10 min" approach fails on real-world ERP data — mid-market ERPs are updated by operators end-of-shift, so ERP timestamps lag OT by hours, not minutes. We use **OF-state-based attribution** instead:

| Component | Detail |
|---|---|
| Algorithm | Poll ERP for OFs in status "In Progress" / "Released". Tag every OT event with the currently active OF. |
| Robustness | **Survives multi-hour clock skew** between OT and IT — joins on OF state, not on timestamps. |
| Detection | Active OF window: from OF "start" status → "complete" status (regardless of when ERP records were updated). |
| Result | Every micro-stop, every kWh, every defect tagged with correct OF, product, planned schedule. |
| Precision | Sub-second at the Edge once OF state is known; ERP poll interval is the bottleneck (1-5 min typical). |
| Failure mode | If no OF is active in the ERP (rare), event is queued and back-attributed when the next OF appears. |

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

### Module 8 — Tribal Knowledge
**Differentiation: Strong** — Field knowledge structured before it disappears.

**Important framing (2026-06-28):** the MOAT is the DATASET (sensor pattern → operator label associations), not the UX that captures it. V1 dropdown + free text accumulates the same site-specific dataset as a sophisticated chatbot. **The chatbot is polish, not the moat.** This means the moat ships in V1, not V2.

| Version | Mechanism | Status |
|---|---|---|
| V0 | Pre-filled dropdown in dashboard (Jam, Air Pressure, Series Change, Material Wait, Adjustment, Other) — 1 click | **V1** |
| V1 | Dropdown + free text field saved to local KG; cause linked to stop event | **V1** |
| V2 | Local SLM chatbot (Ollama/Phi-3) — conversational interview with operator, structured extraction | V2 polish |

The site-specific dataset (operator label paired with sensor pattern at the moment of stop) compounds over months and becomes **impossible to reconstruct without real-time on-site access** — no competitor can copy it post-hoc.

---

### Module 9 — MCP Server (NEW — V1)
**Differentiation: Strong** — Edge MCP is the AI-native edge sovereignty story.

| Component | Detail |
|---|---|
| What it is | Model Context Protocol server embedded in `cmd/server`. Exposes the local Knowledge Graph + pipelines as MCP tools (`kg_query`, `kg_describe_node`, `kg_list_events`, `kg_cost_summary`, etc.). |
| Location | **Edge** by default (V1). Listens on `localhost:5000`. Customer's AI agents connect from inside the factory network. |
| Compatible clients | Claude Desktop, Microsoft Copilot, ChatGPT custom connectors, any MCP-compatible AI agent (MCP is the de-facto standard since 2026). |
| Sovereignty | Data never leaves the customer network. AI agent comes to the data, not the reverse. |
| Cloud relay | Optional opt-in in V1.5+ for remote AI access scenarios (e.g., CEO at home asking via Claude Desktop on laptop). |
| Differentiator | Cognite has MCP but **cloud-only** (data ships to Cognite). MaestroHub CEO claimed MCP in a podcast (unverified, likely cloud-side). UMH has none. **MindSet is the only edge MCP**. |

---

### Module 10 — Ad-hoc Analyst AI Agent (NEW — V1)
**Differentiation: Strong** — Native AI in the product from day 1, not a V2 add-on.

| Component | Detail |
|---|---|
| What it is | Chat panel embedded in the local React dashboard. Plant Manager types natural-language questions, gets grounded answers with KG sources cited. |
| Example prompts | "How did Line 2 perform yesterday?" / "Top 5 micro-stop causes this month with their €cost." / "Which product had the most jams last week?" |
| LLM runtime | **Phi-3 via Ollama (local) by default**. Optional remote LLM (any: OpenAI, Claude, Mistral) with explicit UI disclosure: *"Data will leave your network / EU."* |
| Tool access | All MCP tools exposed by Module 9 — `kg_query`, `kg_list_events`, `kg_cost_summary`, etc. |
| Grounding | Every answer cites the KG nodes / events that informed it. No free-text speculation. |
| Persona | Plant Manager (primary), CFO / Ops Director (secondary). |
| Out of scope V1 | Multi-turn complex reasoning · action-taking (recommendations) · tribal-knowledge interview · multi-site comparison. All deferred to V1.5/V2 (other agents in the 13-agent catalog — see `docs/MindSet_Competitive_Analysis_v2_2.xlsx` Sheet 7). |

**V1 ships exactly this one agent.** The other 12 agents in the catalog (Daily Briefing, Discovery Coach, Tribal Knowledge Chatbot, Causality Reasoner, Trend Spotter, Multi-site Benchmarker, Cost Coach, Alert Triage, Maintenance Scheduler, Compliance Reporter, Connector Recommender, Tag Classifier) are V1.5+, V2, or V3+ based on customer demand signals.

---

## 9. Use Cases

> **Framing (2026-06-28):** the use cases below are the **3 starter templates** that ship in V1. They are NOT the product. The product is the PLATFORM (rules engine + cost model + UNS + KG + MCP + AI agents). Customers + their AI agents build additional use cases (quality, changeover, predictive maintenance, schedule deviation, etc.) on top of the platform. The 3 starter templates exist so Plant Managers see concrete value in the day-1 demo while the platform pitch defends a broader TAM and customer-led roadmap. **Don't impose micro-stops as the only thing MindSet does** — pitch it as one of 3 ready templates, with more buildable in days.

### Use Case 1 — Micro-Stops + OT/IT Reconciliation (V1 starter template #1)

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

### Use Case 2 — Energy Waste (V1 starter template #2)

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

**3 levels of detection** (depends on ERP availability):

| Level | What it detects | Needs ERP? |
|---|---|---|
| **Level 1 — Basic** | "Energy > X kW while machine state = Stop" → ALERT | NO — fast week-1 ROI argument |
| **Level 2 — Cost-attributed** | "OF#456 wasted 18€ of steam during stop" | YES — needs OF context |
| **Level 3 — Comparative** | "Product A uses 12% more energy than Product B for same output" | YES — needs OF + product context |

**Sales angle:** Level 1 ships and demonstrates value before the ERP integration is done. Level 2/3 activate once the V1 ERP connector is configured. Energy typically represents **10-15% reducible cost immediately**.

---

### Use Case 3 — Real OEE vs Declared OEE (V1 starter template #3)

The **single strongest demo** for the investor pitch and for the Plant Manager close. The whole story in one screen.

**The two values:**

**DECLARED OEE** — what operator/supervisor manually reports to MES/ERP. Typically optimistic: micro-stops not counted, downtime miscategorized, optimistic rounding. Usually **5-15 percentage points HIGHER** than reality.

**REAL OEE** — what MindSet calculates from raw OT data:
```
Availability = (Planned_Time - Major_Stops - Micro-Stops) / Planned_Time
  ← from OPC-UA Etat_Machine state transitions via rules engine
Performance  = Actual_Output / Theoretical_Output
  ← from Compteur_Pieces vs Cadence (cost-model config)
Quality      = Good_Parts / Total_Parts
  ← MES integration in V1.5; V1 uses customer-estimated defect rate
OEE = Availability × Performance × Quality
```

**Schedule-gap detection (sub-use-case, V1.5):**
```
IF Actual_Duration_OT > ERP_Schedule_Time + 15% → Alert "Margin Eroded"
Loss = Gap_minutes × Line_Hourly_Cost
```

**The pitch:**
> *"Your declared OEE is 88%. We measured every micro-stop on Line 1 last week — your REAL OEE is 74%. The 14-point gap = 1h04 of hidden downtime per week = X€/week. Here's the Pareto of causes — top 3 fixes recover Y€."*

**The GAP itself is the value proposition** — it directly equals € the Plant Manager didn't know they were losing. Output example: *"OF#123 (Product A) lost 450€ — 12 micro-stops + 50€ steam waste. Action: optimize changeover for this product."*

---

### Use cases that customers + AI agents build NEXT (no commitment from MindSet)

Quality / defect detection · Changeover / setup-time analysis · Predictive maintenance · Schedule deviation alerting · Operator productivity benchmarking · Compliance audit trail generation · Multi-site OEE benchmarking · Anything customer's AI agents can compose from the MCP toolset.

---

## 10. Full Roadmap

> **Rewritten 2026-06-28.** The previous roadmap (32-week plan with parallel "Sessions 1-10") was over-optimistic for a 1-engineer team and didn't reflect the AI-native pull-forward + ERP-in-V1 decisions (see `docs/decisions.md`). New plan: AI-native V1, ERP connectors in V1, realistic timelines for solo dev + Claude Code, explicit hiring milestones, deferred items moved to V1.5 / V2 / V3+.

### Phase 0 — Foundations (DONE — June 2026)

Already shipped or in-progress at June 2026:
- ✅ OPC-UA discovery (`internal/discovery/opcua.go`)
- ✅ HTTP API server (`cmd/server`) on port 8080
- ✅ Edge agent (`cmd/agent`) with MQTT + UNS contextualizer + KG subscriber
- ✅ Knowledge Graph in SQLite (`internal/kg/`)
- ✅ Pipeline engine — YAML pipelines, topological execution, recover()-protected
- ✅ Pipeline Studio (React Flow) — drag-and-drop pipeline builder UI
- ✅ KG viewer (Cytoscape) — technical + domain graphs
- ✅ Dashboard skeleton with live WebSocket push
- ✅ Tag registry + topic registry persisted to SQLite
- ✅ State tracker

**V1 work continues from this base, not from scratch.**

---

### V1 — AI-Native POC (target: end Q1 2027 — ~6-9 months from today)

**Thesis:** three concurrent tracks. Realistic for 1 engineer (Mohamed) with Claude Code acceleration. Hiring a 2nd engineer in this window compresses by ~30%.

**Vertical sequencing in V1**: first pilot customer = **agrifood OR independent metallurgy** (self-serve Plant Manager motion, 48h deployment, <30k€/site — fastest sales cycle). **Pharma + cosmetics deferred to V1.5+** because they require enterprise IT-led sales (6-12 month cycle) + ISO 27001 + GAMP 5 + RBAC — those security additions aren't shipped until V1.5 per the security framework discussion (see analysis_log Entry 20).

**Exit criterion:**
> A first pilot customer in agrifood or metallurgy installs the Edge Agent in 48h, sees the 3 starter templates (micro-stop, energy waste, OEE/TRS) working with real data, queries the KG via the Ad-hoc Analyst chat, and a founder can demo Claude Desktop connecting to the edge MCP server during the customer meeting.

#### Track 1 — Core data pipeline + ERP

| Module | Detail | Target file |
|---|---|---|
| OPC-UA polish | Already shipped. Hardening: secure modes (Sign/SignAndEncrypt cert chain), session resilience, multiple-subscription support | `internal/discovery/opcua.go` |
| Modbus TCP connector | TCP connection, register scan, device fingerprint DB (20-30 common devices) | `internal/discovery/modbus.go` |
| **SQL connector — multi-dialect** | PostgreSQL (`pgx/v5`), MSSQL (`microsoft/go-mssqldb`), MySQL (`go-sql-driver/mysql`). Per-customer dialect via YAML config. | `internal/connectors/sql.go` |
| **OF-state Fuzzy Join engine** | Poll ERP for active OFs; tag every OT event with current OF. **NOT sliding-window** — robust to multi-hour clock skew. (See Module 4 in §8.) | `internal/fuzzy/of_state.go` |
| Cost model | 3-field manual entry wizard (hourly cost, cadence, margin); €/event + €/OF calculation | `internal/cost/model.go` |
| OEE / TRS calculator | Real availability, performance, quality. Computes the declared-vs-real gap. | `internal/cost/oee.go` |
| SQLite ring buffer | 7-15 day rolling window + TTL auto-purge | `internal/storage/ringbuffer.go` |
| Push to cloud (Hybrid edition only) | mTLS + TLS 1.3, offline queue with auto-sync | `internal/push/` |
| Heartbeat monitor | Reports liveness to cloud every 60s; cloud alerts if missing | `internal/push/heartbeat.go` |
| **Local MQTT broker bundle** | Mosquitto as sidecar in docker-compose. Localhost-only listener, no auth needed (intra-container). | `deploy/docker-compose.yml` + `mosquitto.conf` |
| **License key validator** | Validates license against cloud at startup; gracefully degrades to cached entitlements if offline | `internal/license/validator.go` |
| **Local secrets management** | SOPS-encrypted config files for ERP credentials + LLM API keys + cloud auth keys | `internal/secrets/sops.go` + `config/secrets.enc.yaml` |

#### Track 2 — AI core (NEW priority — pulled forward from V2)

| Module | Detail | Target file |
|---|---|---|
| Phi-3 + Ollama integration | Local LLM runtime, model loading, prompt execution, health check | `internal/llm/ollama.go` |
| **MCP server (edge)** | Wraps KG + cost + events API as MCP tools (`kg_query`, `kg_list_events`, `kg_cost_summary`, `kg_describe_node`, etc.). Listens on `localhost:5000`. | `internal/mcp/server.go` |
| **Ad-hoc Analyst agent** | Chat UI in dashboard. Phi-3 default + optional remote LLM with disclosure warning. Grounded answers cite KG sources. | `frontend/src/components/AdHocChat.jsx` + `internal/agents/adhoc.go` |
| Remote LLM config | UI toggle to plug OpenAI / Claude / Mistral. Explicit warning displayed when enabled. | `internal/llm/remote.go` |

#### Track 3 — 3 starter templates + UX polish

| Template / module | Detail | Persona |
|---|---|---|
| **Template 1 — Micro-stop + cost €** | Run→Stop→Run (30s<dur<3min), cause dropdown, Pareto by €. Already partial. | Plant Manager |
| **Template 2 — Energy waste** | Level 1 (no ERP): "Energy > X kW while machine stopped" → alert. Level 2 (with ERP): "OF#456 wasted 18€ of steam." | Plant Manager + CFO |
| **Template 3 — OEE / TRS dashboard** | Real availability vs declared. Shows the gap in €/week. **The killer demo.** | Plant Manager + CFO |
| Onboarding wizard | 3-field cost entry + OPC-UA endpoint config + ERP credentials | All |
| **Tribal Knowledge V1** | 1-click cause dropdown + free-text field on every stop event; saved to KG with link to the event. **The moat ships in V1.** | Operator |
| Dashboard polish | Gantt, Pareto €, ROI simulator, real-vs-declared OEE | Plant Manager |
| Local alerting | SMTP + Slack + Teams direct from edge | All |
| Docker image | `mindsetdata/edge-agent:v1` on PRIVATE registry (proprietary license — distribution-controlled) | — |

#### Parallel track — Website + first customers

| Task | Notes |
|---|---|
| `mindsetdata.io` marketing site (Next.js / Vercel) | Landing + product + use-cases + security + demo-request flow |
| Distribution decision | Private registry (proprietary license affects this) — NOT public Docker Hub like the original roadmap assumed |
| Identify 1-3 pilot customers | Target FR ETI manufacturing — agrifood, pharma, light manufacturing |
| Sales-deck | Already done — `MindSet_Competitive_Analysis_v2_2.xlsx`. Refresh per customer-context. |

#### Hiring milestone in V1

- **Engineer #2 (full-stack Go)** — hire within 4 months of seed close. Compresses V1 timeline by ~30%.
  - Profile: senior Go developer, comfortable with both edge connectors AND React/dashboard work.

---

### V1.5 — Multi-site + AI agent expansion (target: Q2-Q3 2027)

**Triggered when:** 5+ pilot customers signed AND first multi-site customer requests aggregation.

| Module | Description | Priority |
|---|---|---|
| Cross-site KG aggregation | Hybrid edition: cloud-side KG receives transformed events from multiple edges | P0 |
| Multi-site dashboard | Site-vs-site Pareto, benchmark views (CEO / Ops Director persona) | P0 |
| **AI agent expansion (4-5 new agents)** | Daily Briefing, Alert Triage, Discovery Coach, Tag Classifier (agentic) | P0 |
| S7 connector (Siemens) | `gos7` — covers 30-40% of European industrial park | P1 |
| REST connector | Modern ERPs (SAP S/4HANA, D365, Sage X3) | P1 |
| Files / FTP / SFTP connector | CSV / Excel / JSON import from network shares | P1 |
| Cloud MCP relay (opt-in) | For customers wanting remote AI access without VPN | P1 |
| Historian connector (PI / Wonderware / InfluxDB) | Push enriched events to customer's existing historian | P1 |
| Microsoft Teams alerting | Microsoft-centric ETIs | P2 |
| 4th + 5th starter templates | Based on first-customer signals (candidates: Quality, Changeover, Schedule deviation) | P0 |
| Schedule-gap detection (sub-feature of OEE template) | `IF Actual_Duration_OT > ERP_Schedule_Time + 15% → Margin Eroded` alert | P1 |

**V1.5 exit criterion:**
> Multi-site customer sees consolidated cross-site OEE + cost + Pareto in cloud dashboard. AI Daily Briefing arrives in inbox at 6am every weekday.

**Hiring milestone:** Engineer #3 (DevOps / cloud platform) for multi-tenant cloud + K8s automation for Self-Hosted edition.

---

### V2 — Deep AI + ecosystem (target: Q4 2027 - Q1 2028)

**Triggered when:** 15+ paying customers, V1.5 stable, 6+ months of accumulated tribal-knowledge dataset.

| Module | Description | Priority |
|---|---|---|
| **Tribal Knowledge Chatbot** | Phi-3 conversational interview with operator after each stop — replaces V1 dropdown for richer capture. (Note: the MOAT — the dataset — already ships in V1. This is UX polish, not the moat.) | P0 |
| **Causality Reasoner agent** | "Pressure dropped 12s before this stop — could be a leak." Multi-step LLM reasoning. | P0 |
| **Trend Spotter agent** | Proactive surfacing: "Stops on Line 3 doubled in 3 days." | P1 |
| **Multi-site Benchmarker agent** | "Site A vs Site B on TRS / causes / cost." | P1 |
| **Cost Coach agent** | Explains cost model to CFO; suggests refinements | P1 |
| Sparkplug B | MQTT with ISA-95 structured payload | P1 |
| MQTT generic | Modern IIoT gateways | P1 |
| MTConnect | CNC / machining / metallurgy | P2 |
| BACnet/IP | Building / HVAC for energy-intensive sites | P2 |
| BYOC deployment automation | Helm charts + docker-compose for customer on-prem K8s | P0 |
| **Public MCP tool catalog** | Documented MCP schema so partners can build agents on top | P1 |
| Auto-update Edge Agent | Roll out new versions to consenting customers | P1 |

**V2 exit criterion:**
> A third-party AI agent (e.g., Claude Desktop on a customer's CISO laptop) queries the factory KG via MCP and produces a grounded compliance report — no hallucinations, all sources cited.

**Hiring milestone:** Engineer #4 (ML / data eng) — owns AI agent quality, eval harness, and the tribal-knowledge dataset operationalization.

---

### V3+ — Scaling + ecosystem (2028-2029)

**Triggered when:** 50+ paying customers, V2 stable, considering international expansion.

| Axis | Description |
|---|---|
| **License reconsideration (2028)** | Per locked decision: evaluate open-core or source-available models for the edge agent. Cloud + enterprise stays proprietary. Decision driven by whether OSS pressure from UMH/others is causing lost deals. |
| **Hyperscaler edition reconsideration (2029)** | Per locked decision: separate product line for US/APAC scaling (AWS / Azure / GCP). EU pipeline stays no-hyperscaler. Decision driven by international demand signals and whether sovereignty moat is established enough to survive a separate "global" SKU. |
| **Predictive ML** | On accumulated KG dataset — failure prediction per machine, product, season. Per-site model trained on local data. |
| **Cross-industry KG benchmarks** | Sectoral anonymized data: "you are at the X-th percentile of agrifood ETIs for stop frequency." Opt-in. |
| **Partner SDK** | Integrators + ISVs build custom MCP tools + connectors + agents on the platform |
| **Functions marketplace** | Curated community library of YAML pipelines per industry. Build only if community demand materializes (current: no, defer). |
| **Additional protocols (V3 catalog)** | Kafka, AWS S3, Azure Blob, NATS, MongoDB, Omron FINS, OPC-DA, LoRaWAN — driven by customer demand, not pre-decided. |
| **Acquisition moat at scale** | Each site = irreversible cumulative site fingerprint → structurally impossible churn. Cross-industry KG → sectoral reference data only MindSet has. |

---

### Timeline & headcount summary

| Phase | Target | Headcount required | Realistic? |
|---|---|---|---|
| V1 (POC complete, 1-3 pilot customers) | Q1 2027 | 1-2 engineers | YES with 2nd engineer hired month 1-2 post-seed. TIGHT solo. |
| V1.5 (multi-site, 5-10 customers) | Q3 2027 | 3 engineers | Feasible |
| V2 (deep AI, 15-25 customers) | Q1 2028 | 4 engineers | Feasible |
| V3+ (50+ customers, international option) | 2028-2029 | 6-8 engineers | Requires Series A funding |

**Key insight on team size:** the previous roadmap assumed parallel "Sessions 1-10" execution which requires multiple engineers running in parallel. With **1 engineer (current state)**, realistic V1 ship is **6-9 months solo, OR 4-5 months with a 2nd engineer hired in month 1-2 of post-seed funding.** The investor pitch should explicitly request funding to hire 2-3 engineers within 6 months — not pretend the 2-founder team can ship the V2 vision alone.

---

### Distribution model (revised — proprietary license affects this)

Original roadmap assumed public Docker Hub distribution. The locked closed-source decision changes this:

**V1 distribution:**
```powershell
# Customer receives a license key + private registry credentials after signing.
docker login registry.mindsetdata.io -u customer-name -p <license-key>
docker pull registry.mindsetdata.io/edge-agent:v1
docker run -d `
  -e LICENSE_KEY=<key> `
  -e SITE_NAME="Usine Paris Nord" `
  -e EDITION=hybrid `
  --network host `
  registry.mindsetdata.io/edge-agent:v1
# → Agent starts, scans network, finds OPC-UA equipment
# → Local dashboard available at http://localhost:8080 immediately
# → Hybrid edition: cloud dashboard available at app.mindsetdata.io within 1h
```

Distribution gates: license key required, telemetry of license validity, no free public-Docker-Hub pull. **Reconsidered in 2028 if license model shifts to open-core** (then Edge Agent goes back to public Docker Hub; only cloud/enterprise features stay license-gated).

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
app.mindsetdata.io          ← Hybrid edition cloud dashboard (React / Scaleway FR)
api.mindsetdata.io          ← Hybrid edition cloud API receiver (Go / Scaleway FR)
hub.docker.com/mindsetdata  ← Edge Agent Docker image (proprietary binary distribution)
```

### Cost per edition (per site)

| Edition | Cloud cost / month | What drives it |
|---|---|---|
| **On-Premise** | **0 €** — zero cloud component | Customer pays own edge hardware (~600€ amortized PC) |
| **Hybrid (default)** | **~15 €** at V0 (Scaleway FR PLAY2-NANO + Managed PostgreSQL + Object Storage). Scales with multi-site customers (~30-50€/mo for 5+ sites). | KG aggregation + remote dashboard + encrypted backup + heartbeat monitor |
| **Self-Hosted** | **Variable — customer pays** their EU cloud account. Indicative: Hetzner CX21 ~5€, Scaleway VPS 7-15€, OVH 8-12€, Bleu = TBD, on-prem K8s = sunk cost. | Same workload as Hybrid, on customer infra |

**Bleu note:** Bleu is the FR sovereign cloud joint venture (Orange + Capgemini using Microsoft Azure tech under FR jurisdiction). Valid Self-Hosted target for customers wanting Azure-style services with FR sovereignty guarantees.

### Scaleway FR — Hybrid edition V0 detail

```
1× PLAY2-NANO VPS (2 vCPU, 2GB RAM)   → 3.99€/month
1× Managed PostgreSQL DEV-1500         → 9.99€/month
1× Object Storage (KG snapshots)       → pay per use ~1€/month
──────────────────────────────────────────────────────
Total V0 per customer:                 → ~15€/month
```

### NOT offered — Hyperscaler edition

AWS / Azure / GCP (including their EU regions) are **explicitly excluded** from the catalog through 2029. Reason: US CLOUD Act exposure invalidates the sovereignty pitch for the highest-value verticals (defense, public sector, regulated pharma). Reconsidered for international scaling in 2029 — at that point it becomes a separate product line, not a dilution of the core offering.

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
|---|---|---|---|
| Edge Agent | Go + Docker | Go 1.22+ | Lightweight, static binary, cross-platform, low footprint |
| OPC-UA | `gopcua` | latest | Mature, browse + subscriptions + security |
| Modbus | `goburrow/modbus` | latest | TCP + RTU |
| Siemens S7 (V1.5+) | `gos7` | latest | Direct PLC Siemens access |
| MQTT | `paho.mqtt.golang` | latest | Eclipse standard |
| Sparkplug B (V2) | `sparkplug-go` | latest | Structured industrial payload |
| SQL — PostgreSQL | `pgx/v5` | latest | Modern ERPs, Odoo, custom |
| SQL — MSSQL | `microsoft/go-mssqldb` | latest | Sage X3, Dynamics 365 on-prem |
| SQL — MySQL | `go-sql-driver/mysql` | latest | Web-based ERPs |
| SQL — Oracle (V1.5+) | `sijms/go-ora` | latest | SAP large accounts |
| Local SLM | Phi-3 via Ollama | latest | Sovereign, zero latency, zero cloud |
| **MCP server** | Anthropic MCP spec (de-facto standard since 2026) | latest | Edge-default; exposes KG to AI agents |
| Local DB | SQLite (`modernc.org/sqlite`) | — | Zero deps, ring buffer, pure-Go |
| Pipeline | Redpanda Connect | latest | ex-Benthos, declarative YAML |
| Ontology | ISA-95 | standard | Maximum interoperability |
| Encryption | mTLS + TLS 1.3 | — | NIS2 compliant |
| Cloud DB (Hybrid) | PostgreSQL + TimescaleDB | pg15 | Time-series + events |
| Cache (Hybrid) | Redis | 7 | Real-time dashboard |
| Cloud (Hybrid) | Scaleway FR / OVH FR | — | FR sovereignty, RGPD. **No hyperscaler edition through 2029.** |
| Cloud (Self-Hosted) | Hetzner / IONOS / T-Systems / 3DS Outscale / Bleu / customer K8s | — | Customer's EU jurisdiction choice |
| Frontend | React | 19+ | Local + cloud dashboard |
| Website | Next.js | 14+ | SEO, Vercel deploy |
| CI/CD | GitHub Actions | — | Build + test + push Docker |
| **License** | **Proprietary (closed-source, 2-year minimum)** | — | Early-stage commercial protection. Open-core / source-available reconsidered in 2028. |

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

### The 5 defensive elements (revised 2026-06-28)

**1. Zero-manual auto-discovery**
Network scan + behavioral inference + Phi-3 SLM classification = automatic connection without mapping. 48h vs 3 months for an integrator. Structural advantage on client acquisition cost.

**2. OF-state-based attribution (formerly "Fuzzy Join")**
OT/IT reconciliation by reading **Fabrication Order state from the ERP** — not by joining on timestamps. Robust to multi-hour clock skew typical of mid-market ERPs (where operators update records end-of-shift). Every micro-stop, every kWh, every defect correctly tagged with product + OF, without per-customer time-sync engineering. UMH leaves this to the user (Node-RED); MaestroHub doesn't address it as a dedicated feature; Cognite does entity contextualization (P&ID OCR) — a different problem. Hard to build, invisible from outside — the most defensible component in the stack.

**3. The cumulative Knowledge Graph (site fingerprint)**
Site-specific non-replicable context: every micro-stop, every cause, every calibrated cost model accumulates over months. Replacing MindSet = losing all accumulated intelligence. After 6 months on-site, **churn becomes structurally prohibitive**.

**4. Tribal Knowledge structured over time — ships in V1**
`sensor pattern → operator label` associations: impossible to reconstruct without real-time on-site access. **The MOAT is the dataset, NOT the chatbot UX** — V1 captures the dataset via 1-click dropdown + free text. V2 chatbot is polish, not the moat. Compounds with moat #3 (site fingerprint).

**5. Edge sovereignty + edge MCP (NEW)**
MindSet runs MCP server AT THE EDGE — AI agents (Claude, Copilot, our Ad-hoc Analyst) query the factory floor directly without raw data leaving the customer network. Cognite has MCP but cloud-only (data ships to Cognite cloud). MaestroHub CEO claims MCP in a podcast (unverified, likely cloud-side if real). UMH has none. Combined with **no-hyperscaler-edition through 2029**, this is the strongest EU-regulatory moat — defense, public sector, regulated pharma cannot use the others.

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
*Last updated: May 2026 (body) / 2026-07-28 (reality-check header only — see top of doc)*
*License: Proprietary (closed-source), 2-year minimum — this footer previously said Apache 2.0, which was the superseded original decision, not the current one. See the reality-check block at the top of this document.*
