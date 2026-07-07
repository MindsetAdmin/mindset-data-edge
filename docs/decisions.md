# Mindset Data — Technical Decisions Log

> Running log of all architectural and implementation decisions.
> One entry per decision, newest first.
> Format: **Decision** → **Rationale** → **Alternatives rejected**

---

## Corrections & Late Decisions (Sprint 3 cont. — June 2026)

---

### Knowledge Graph: merged into ONE unified graph with category tags (was 2 KGs)

**Decision (2026-07-02):** The previous split between "Domain KG" (persistent site fingerprint) and "Technical KG" (in-memory pipeline topology) is REMOVED. There is now **one Knowledge Graph** persisted in SQLite. Every node and edge carries a `category` tag: `business` (Equipment, Event, Cause, Cost, Operator, OF, Product, …) or `platform` (Connection, Topic, Function, Pipeline, Dashboard).

The unified API is `GET /api/kg?category=business|platform|all`. Legacy `/api/kg/domain` and `/api/kg/technical` are preserved as aliases mapping to `?category=business` and `?category=platform` respectively.

**Rationale:**
- Aligns with the "single trusted source for AI agents" principle (Prop #7 in `Prpopsitions1.md`) — one endpoint, one schema, one MCP tool surface
- Enables cross-category queries via a single graph traversal (e.g. *"which pipelines produce cost data for Line 2?"*)
- Removes the mental model confusion — "which KG?" — from user + investor + developer conversations
- Consolidates cache logic (the 5-min in-memory cache for the Technical KG is gone; platform sub-graph is rebuilt on pipeline register/deregister, then queried from SQLite like any other data)
- Simplifies the Impact Engine's data access (Entry 40)

**Trade-off accepted:** Platform sub-graph is now persisted in SQLite (was in-memory only). Adds a few kB per pipeline to the SQLite footprint — negligible.

**Backwards compatibility:** Legacy Go API (`GetFullGraph`, `GetTechnicalGraph`, `PurgeCache`) preserved as aliases. Legacy REST endpoints preserved. Old databases auto-migrate via `ALTER TABLE ADD COLUMN category` (SQLite idempotent-safe).

**Alternatives rejected:**
- Keep two KGs, rename Technical KG to "Platform Topology" (Option B in Entry 49) — clarifies the mental model but doesn't unlock cross-category queries
- Delete Technical KG entirely (Option C in Entry 49) — loses the pipeline-topology view which is useful for debugging and investor demos

---


---

### Target market: 15K+ EU mid-sized factories TAM, initial GTM focus on 4 high-value verticals

**Decision:** MindSet's total addressable market is **15,000+ European mid-sized factories**. Initial go-to-market focus is **4 high-value verticals**:
- 💊 **Pharma**
- 💄 **Cosmetics**
- 🌾 **Agrifood**
- ⚙️ **Metallurgy**

Geographic execution starts in France (founders' geography + Boost10x network), expanding to DACH + Italy + Spain + Nordics in V2-V3.

**Rationale:** All 4 verticals share traits that fit MindSet's product:
- High willingness to pay (high-value products, regulated)
- Sovereignty-sensitive (GMP / EU Cosmetic Regulation / HACCP / industry security requirements)
- High and measurable financial impact from downtime + waste
- Well-suited to the deterministic rules engine + OF-state Fuzzy Join + cost-in-€ model

TAM math: ~12,000-25,000 EU mid-sized factories across these 4 verticals. At conservative 30k€/site/year pricing → ~450M€ TAM. At pharma-supported 100k€/site/year → ~1.5B€ TAM. Both numbers are credible for a pre-seed investor pitch.

**Important downstream tension flagged (not locked yet, see analysis_log Entry 37):** the 4 verticals don't share the same sales motion. Agrifood + metallurgy fit the original ETI Plant-Manager self-serve <30k€ motion; pharma + cosmetics typically require enterprise IT-led 6-12 month cycles + ISO 27001 + GAMP 5. Likely needs two parallel sales motions — to be locked as a separate decision once explicitly confirmed.

**Alternatives rejected:**
- ETI manufacturing only (no vertical focus) — too generic, weakens positioning
- All verticals at once (utilities, energy, logistics included) — fragments positioning, weakens vs UMH / Cognite
- Pharma-only or single-vertical — narrows TAM unnecessarily, ties success to one regulatory environment
- Geographic execution multi-country from V1 — not realistic with 2 founders + 2 interns

---

### Local MQTT broker: bundled in multi-container docker-compose (NOT separate install)

**Decision:** Mosquitto runs as a sidecar container in MindSet's docker-compose. Customer install command remains a single `docker compose up`. Customer does NOT have to install or maintain Mosquitto separately.

**Rationale:** The "48h deployment" + "one Docker command" pitch claim is preserved. Customer's IT team doesn't have to evaluate / approve / patch Mosquitto separately — it's part of the MindSet image we ship and update. Mosquitto config is locked-down (localhost-only listener, no auth needed since intra-container) and bundled with the edge agent's deployment unit.

**Alternatives rejected:**
- Separate Mosquitto install (breaks "1 command" pitch, adds customer-side maintenance burden, adds CVE patching responsibility on the customer).
- Embed an MQTT broker library in the Go binary (e.g., embedded mochi-mqtt). Reduces operational complexity but adds binary size + tighter coupling; revisit at V2 if Mosquitto becomes a maintenance pain.

---

---

### Licensing model: PROPRIETARY (closed-source) for first 2 years — supersedes prior Apache 2.0 decision

**Decision:** MindSet Data is shipped as **closed-source proprietary software** for at least the first 2 years (V1, V2). Customers receive compiled Go binaries + React build, not source code. The decision is reconsidered in 2028 — open-core or source-available options remain on the table for that horizon.

**Rationale:** Early-stage company protection. Closed source preserves commercial control during PMF discovery and prevents free-rider competitors from forking before MindSet has captured customer relationships. Two years of closed runway also lets the team focus on customer acquisition rather than community management.

**Trade-off accepted:** Loses the "Apache 2.0 vs Cognite proprietary" line in the comp matrix. UMH now wins outright on OSS dimension. Compensated by stronger positioning on sovereignty + edge-native + OF-based Fuzzy Join + MCP + simplicity.

**Alternatives rejected:**
- Apache 2.0 from V1 (previous decision) — exposes IP to fast-follower competitors with no defensive moat in early years.
- Source-available BSL/PolyForm — adds legal complexity + customer confusion without the trust upside of true OSS.
- Open-core from V1 — splits engineering attention between OSS edge agent and proprietary cloud, premature for a 1-engineer team.

---

### Fuzzy Join algorithm: OF-state-based attribution — supersedes prior sliding-window decision

**Decision:** MindSet's OT/IT reconciliation works by **reading Fabrication Order (OF) state from the ERP** — polling for OFs currently in status "In Progress" / "Released" — and tagging every OT event happening during an active OF with that OF's metadata (product, OF ID, planned schedule). The algorithm joins on **OF state, not on timestamps**.

**Rationale:** The ±10 min sliding-window approach (in original docs) fails on real ERP data. Mid-market ERPs are updated by operators end-of-shift; ERP timestamps lag OT by hours, not minutes. OF-state-based attribution is robust to this multi-hour clock skew. The result: every micro-stop, every kWh, every defect is correctly tagged with product + OF without per-customer time-sync engineering.

**Trade-off accepted:** Requires ERP integration (no Fuzzy Join without ERP — which is fine since ERP connector is now V1). Cannot attribute events to OFs that aren't represented in the ERP at the time the event occurs (rare in practice).

**Alternatives rejected:** Sliding-window time-based join (breaks on real ERP latency); rely solely on operator-entered OF assignment (manual burden, error-prone).

---

### Edition naming: On-Premise / Hybrid / Self-Hosted — supersedes Air-Gap / Sovereign Cloud / BYOC

**Decision:** The three deployment editions are renamed to consumer-friendly terminology:
- **On-Premise** (formerly Air-Gap) — zero cloud, per-site only. Target: defense, public sector, sensitive pharma.
- **Hybrid** (formerly Sovereign Cloud) — Scaleway FR / OVH FR for multi-site KG + remote dashboard + backup. Default.
- **Self-Hosted** (formerly BYOC) — customer deploys cloud tier on their EU cloud or on-prem Kubernetes.

**Rationale:** Plain-language names land better with Plant Manager / CFO buyers who don't speak DevOps. "Air-Gap" is technical jargon; "BYOC" is industry acronym. "On-Premise / Hybrid / Self-Hosted" is the same thing in language any buyer understands.

**Alternatives rejected:** Keep technical names (less accessible); use "Local / Cloud / Custom" (loses the sovereignty implication).

---

### Hyperscaler edition: NOT offered in V1-V2-V3 — reconsider 2029 for international scaling

**Decision:** No AWS / Azure / GCP edition through V1, V2, V3 (through ~2029). The "no hyperscaler" stance holds for at least 3 years. **In 2029**, reconsider adding hyperscaler support for international (US / APAC) expansion — at that point the EU sovereignty moat is established and adding hyperscaler support targets a different market segment.

**Rationale:** Adding hyperscalers in V1-V3 collapses the sovereignty moat for the highest-value EU verticals (defense, public sector, regulated pharma). The TAM expansion isn't worth losing the differentiation. International expansion is a 2029+ concern — by then the EU footprint is real, the sovereignty pitch has been validated by reference customers, and a separate "Global Edition" with hyperscalers becomes a SEPARATE PRODUCT LINE, not a dilution of the core offering.

**Alternatives rejected:** Add hyperscalers in V1-V2 (kills sovereignty moat early); never add them (caps TAM permanently at EU manufacturing — leaves international upside on the table).

---

## V1 Scope & AI-Native Positioning (Sprint 3 — June 2026)

---

### Platform-first positioning with 3 starter use-case templates

**Decision:** MindSet is positioned as an AI-native edge **platform**, not a single-use-case product. V1 ships with 3 ready-to-use templates that customers can deploy on day 1: **micro-stop detection**, **energy waste detection**, and **OEE / TRS dashboard**. Customers and their AI agents can build additional use cases (quality, changeover, predictive, etc.) on top of the platform.

**Rationale:** "Don't impose micro-stops" — first customers will reveal which use cases to invest in. But "platform without a vertical" is the classic startup death — too generic to demo, too long a sales cycle. The 3 starter templates give Plant Managers something concrete to see in the demo while the platform claim defends broader TAM and customer-led roadmap.

**Alternatives rejected:**
- Single-use-case positioning ("micro-stop detection product") — narrows TAM and contradicts the AI-native + platform pitch.
- Pure platform ("you build everything") — no demo, no first-customer hook, indistinguishable from UMH.
- 5+ starter templates — over-scopes V1, dilutes the polish on each.

---

### AI-native from V1 (not V2)

**Decision:** AI capabilities ship in V1, not as a V2 add-on. V1 includes: Phi-3 / Ollama local runtime, edge MCP server exposing the KG, and one native AI agent (Ad-hoc Analyst). The product narrative becomes "AI-native edge industrial platform," not "industrial platform with AI added later."

**Rationale:** 2026 investor expectation is AI-native by default. Shipping AI in V2 makes the deck look behind the curve. AI integrated from the beginning also de-risks the architecture — no last-minute retrofit needed.

**Alternatives rejected:**
- AI in V2 (per original roadmap) — weak investor pitch in 2026 ("we'll add AI later"), risks late-stage architecture rework.
- Multiple agents at V1 — over-scopes the first ship, dilutes quality of each. One excellent agent beats five mediocre.

---

### ERP connectors in V1 (pulled forward from V1 mid-roadmap)

**Decision:** SQL connector (Fuzzy Join input) is part of V1, not V1.5. This makes Moat #2 (Fuzzy Join OT/IT) demoable from first customer install.

**Rationale:** Fuzzy Join is the technical moat that distinguishes MindSet from UMH (no built-in Fuzzy Join) and from MaestroHub (no clear OT/IT temporal alignment). Demoing it requires an ERP connector. Pulling the SQL connector forward makes the moat real at V1, not a future promise.

**Alternatives rejected:** Keep SQL connector at V1.5 — leaves Fuzzy Join undemoable in V1, undermines the strongest technical moat in the pitch.

---

### SQL connector V1 dialects: PostgreSQL + MSSQL + MySQL

**Decision:** V1 SQL connector supports PostgreSQL (via `pgx/v5`), MSSQL (via `microsoft/go-mssqldb`), and MySQL/MariaDB (via `go-sql-driver/mysql`). Oracle and SAP HANA dialects deferred to V1.5+ based on customer demand signal.

**Rationale:** PostgreSQL + MSSQL + MySQL covers ~80% of FR ETI ERP backends (Sage X3 = MSSQL, Dynamics 365 on-prem = MSSQL, Odoo = PostgreSQL, modern web ERPs = MySQL). Oracle is tied to large-account SAP deals which aren't first-customer targets. SAP HANA is enterprise S/4HANA territory — wrong segment for MindSet's ETI mid-market.

**Alternatives rejected:**
- PostgreSQL only (covers ~30% of FR ETIs — leaves Sage / Dynamics customers blocked).
- All 5 dialects in V1 (Oracle + HANA = high effort for wrong segment).

---

### MCP server: edge-only in V1, cloud relay deferred to V1.5+

**Decision:** V1 MCP server runs at the edge only, exposing the local KG to AI agents inside the customer's network (Claude Desktop, Copilot, MindSet's native agent). Cloud MCP relay for remote AI access is deferred to V1.5+ based on remote-access demand signal.

**Rationale:** Edge-only MCP simplifies V1 architecture (one binary, no cross-network auth), preserves the sovereignty default (data stays where it was generated), and is sufficient for the V1 customer profile (Plant Manager working inside the factory, occasional founder demo with Claude Desktop on a laptop on the factory LAN).

**Alternatives rejected:** Cloud-only MCP — breaks sovereignty default. Edge + cloud relay at V1 — doubles the V1 architecture surface for no first-customer benefit.

---

### V1 native AI agent: Ad-hoc Analyst (sole agent)

**Decision:** V1 ships exactly one native AI agent: **Ad-hoc Analyst** — chat UI embedded in the local dashboard, Phi-3 local default, grounded answers via MCP-tool access to the KG. Cites the KG nodes / events used. All other agents from the 13-agent catalog (Daily Briefing, Discovery Coach, Tribal Knowledge Chatbot, Causality Reasoner, etc.) are V1.5+ or V2.

**Rationale:** One excellent demoable agent beats five mediocre ones. Ad-hoc Analyst is the strongest demo (Plant Manager types a question, gets a real answer with sources — the chat UX every 2026 user expects). Built on the same MCP infrastructure, so it doubles as proof-of-concept for the MCP integration. Simple enough to ship within the V1 timeline.

**Alternatives rejected:**
- Discovery Coach first — useful but onboarding-only, weaker demo than Q&A.
- Daily Briefing first — needs accumulated data, doesn't demo well on day 1.
- 3-5 agents in V1 — over-scopes for a 1-engineer team.

---

### Tribal Knowledge moat ships in V1 via dropdown + free text (NOT V2 chatbot)

**Decision:** Moat #4 (Tribal Knowledge) ships in V1 as a 1-click cause dropdown + free-text field on every detected stop event, with the cause linked to the stop event in the KG. The V2 chatbot (Phi-3 conversational interview) is a stretch goal for richer capture UX, not the moat itself.

**Rationale:** Re-reading the moat definition in `docs/mindset.md` section 15: *"sensor pattern → operator label associations: impossible to reconstruct without on-site real-time access."* **The moat is the DATASET, not the UX that captures it.** A simple dropdown + free text accumulates the same site-specific pattern as a sophisticated chatbot would. The chatbot is polish, not the moat. This makes Tribal Knowledge a V1-realisable claim instead of a V2 promise.

**Alternatives rejected:**
- Defer Tribal Knowledge entirely to V2 — leaves one of the 5 moats undemonstrable at first customer install.
- Phi-3 conversational chatbot in V1 — Phi-3 conversational quality in French + operator jargon + interruption-handling is too risky for V1 ship.

---

## Strategic Positioning (Sprint 2 — June 2026)

---

### Three deployment editions: Air-Gap / Sovereign Cloud / BYOC — no hyperscaler edition

**Decision:** The product is offered in exactly three editions:
- **Air-Gap** — zero cloud component. Per-site only. No multi-site, no remote dashboard. Target: defense, public sector, sensitive pharma, nuclear.
- **Sovereign Cloud (default)** — Scaleway FR / OVH FR for cross-site KG aggregation, multi-site dashboard, encrypted backup, heartbeat monitor.
- **BYOC** — customer deploys the cloud tier on their own EU-jurisdiction cloud (Hetzner, IONOS, T-Systems, 3DS Outscale) or on-premise Kubernetes.

There is **no hyperscaler edition**. AWS, Azure, GCP (including their EU regions) are explicitly excluded.

**Rationale:** Lets customers self-select by sovereignty needs. Defense and regulated industries get pure air-gap. Commercial ETI gets the convenient default. Large multi-site customers with existing EU cloud relationships get BYOC. Excluding hyperscalers preserves the regulatory moat — US CLOUD Act exposure on AWS-EU and Azure-EU would invalidate the sovereignty pitch for public sector and defense buyers.

**Alternatives rejected:** Single one-size SaaS (excludes air-gap and regulated industries); supporting AWS/Azure to broaden TAM (breaks sovereignty story for the highest-value verticals).

---

### Cloud tier scope: aggregation + remote view + backup + heartbeat only

**Decision:** Cloud tier is limited to: cross-site KG aggregation, multi-site / remote dashboard serving, site management API (auth + keys + entitlements), encrypted KG snapshots, alerting heartbeat (liveness monitor). Nothing else runs in the cloud.

**Rationale:** A feature goes to the cloud only when it satisfies all three: (1) it needs to span multiple sites or be reached from outside the factory, (2) latency tolerates >1s round-trip, (3) only already-transformed data crosses the boundary. Discovery, contextualization, rules engine, cost model, Fuzzy Join, dashboards, alerting, AI agents, and MCP server all run at the edge.

**Alternatives rejected:** Cloud-side pipeline execution (latency + sovereignty issues); cloud-side rules engine (sub-second decisions impossible from cloud round-trip); cloud-side SLM (defeats the local-first sovereignty default).

---

### Alerting: edge-direct SMTP/Slack/Teams, cloud component is heartbeat monitor only

**Decision:** The edge agent sends emails, Slack, and Teams alerts directly via the customer's outbound HTTPS / SMTP. The cloud component related to alerting is a **liveness / heartbeat monitor** that alerts when an edge agent stops reporting — not a general-purpose SMTP relay.

**Rationale:** Honest framing of the actual value. In practice, almost all customers can send outbound SMTP from the factory. The real reason to have a cloud-side alerting component is to detect a dead edge agent. Reframing the component clarifies its actual purpose for both customers and internal team.

**Alternatives rejected:** Generic cloud SMTP relay (used by <10% of customers in practice, muddies the architecture story); no cloud alerting at all (loses ability to detect dead edge agents — important for operations).

---

### MCP server: essential feature, edge-default with optional cloud relay

**Decision:** MindSet exposes its Knowledge Graph and pipelines via a Model Context Protocol (MCP) server, running on the edge by default. External AI agents (Claude Desktop, Copilot, ChatGPT custom connectors, etc.) connect to the local MCP server. An optional cloud MCP relay is offered for customers needing remote AI access without setting up a VPN.

**Rationale:** MCP is becoming the de-facto standard for AI agent integrations. Native MCP support is a strong differentiator vs Cognite (closed proprietary AI Atlas SDK) and most mid-market rivals. Edge-default preserves the sovereignty pitch — the AI agent comes to the data, not the reverse.

**Alternatives rejected:** REST/GraphQL-only API for AI agents (won't plug into Claude Desktop / Copilot natively, weaker investor story); cloud-only MCP server (breaks the sovereignty default — defeats the purpose).

---

### AI provider strategy: local-default + optional remote with explicit disclosure (Option B)

**Decision:** Phi-3 via Ollama is the default local LLM. Customer can plug any LLM (OpenAI, Anthropic, Mistral, Aleph Alpha, Azure OpenAI) via configuration. When remote LLM is enabled, the UI explicitly warns the operator: *"Data will leave your network / EU."*

**Rationale:** Sovereignty pitch holds **by default**. Customer has flexibility for use cases where they want to use existing LLM contracts (e.g., an enterprise Azure OpenAI agreement). The explicit warning preserves the founder's honesty contract — informed consent rather than silent leak.

**Alternatives rejected:** Strict local-only (limits flexibility — locks out customers with existing LLM relationships); any-LLM-no-warnings (cannot claim sovereignty as default value prop — undermines the entire pitch); strict EU-LLM-only (Mistral/Aleph Alpha only — too narrow, prevents Azure OpenAI integration which is common in enterprise).

---

### BYOC scope: EU-jurisdiction cloud OR customer's on-prem Kubernetes only

**Decision:** "Bring Your Own Cloud" means EU-jurisdiction cloud (Scaleway, OVH, Hetzner, T-Systems, IONOS, 3DS Outscale) or the customer's own on-premise Kubernetes. AWS, Azure, GCP — including their EU regions — are explicitly excluded.

**Rationale:** AWS-EU and Azure-EU regions are subject to the US CLOUD Act, which lets US authorities compel data disclosure regardless of physical location. This invalidates the sovereignty pitch for FR public sector, defense, and regulated industries. The sovereignty moat must hold cleanly — accepting hyperscalers would let competitors (and customers) frame MindSet as "another SaaS that pretends to be sovereign."

**Alternatives rejected:** Any-cloud BYOC including AWS/Azure (broadens TAM but breaks the regulatory moat for the highest-value verticals); Scaleway/OVH only (too narrow — locks out customers with existing EU cloud relationships at Hetzner, T-Systems, etc.).

---

## OPC-UA Discovery (Sprint 1 — June 2026)

---

### NodeID as diff key in `tagsToMap()`
**Decision:** Use NodeID (not tag name / BrowseName) as the unique key for the `WatchForChanges()` diff.

**Rationale:** NodeID is the stable OPC-UA identifier assigned by the server. BrowseName / DisplayName can be renamed by a SCADA operator at any time without the NodeID changing. Diff keyed by NodeID is therefore resilient to tag renames — only genuine additions and removals trigger `onChange`.

**Alternatives rejected:** BrowseName as key — would generate false positives on every rename.

---

### Browse sleep: 50ms between requests
**Decision:** Insert a 50ms sleep between each `Browse` RPC call during node tree discovery.

**Rationale:** Two goals simultaneously:
1. Prevents flooding the OPC-UA server (some embedded SCADA servers have very limited concurrency).
2. Avoids session timeout on slow servers — a long synchronous browse without any activity risks a server-side session expiry.

**Alternatives rejected:** No sleep (fast but risky on fragile PLCs). 200ms (too slow on large node trees with hundreds of nodes).

---

### Continuation Points: always release via `BrowseNext(ReleaseContinuationPoints: true)`
**Decision:** After exhausting all continuation points, always call `BrowseNext` with `ReleaseContinuationPoints: true` to explicitly release server-side state.

**Rationale:** OPC-UA servers maintain server-side state for each active continuation point. Not releasing them leaks server resources and can exhaust the server's continuation point pool, causing browse failures for other clients or subsequent sessions.

**Alternatives rejected:** Fire-and-forget (skip final release) — would cause resource leaks on long-running deployments.

---

### `RequestedMaxReferencesPerNode`: 100
**Decision:** Set `RequestedMaxReferencesPerNode` to 100 in all Browse requests.

**Rationale:** Value of 0 (unlimited) is theoretically valid per the spec but many embedded OPC-UA servers ignore it or return errors. 100 is well within the safe range for all tested servers (Prosys, Siemens, Ignition). Produces continuation points naturally for nodes with many children, which is already handled.

**Alternatives rejected:** 0 (unlimited) — incompatible with some embedded servers. 1000 — risks memory spikes on servers with large flat namespaces.

---

### `ReferenceTypeID`: 0:33 (HierarchicalReferences), not 0:31
**Decision:** Use `ReferenceTypeID = 0:33` (`HierarchicalReferences`) for all Browse requests, not 0:31 (`References`).

**Rationale:** `HierarchicalReferences` (0:33) is the correct reference type for navigating the OPC-UA address space tree — it covers `Organizes`, `HasComponent`, `HasProperty`, and `FolderType` references, which are the relations used to build the node hierarchy. `References` (0:31) is the abstract base type and includes non-hierarchical references (e.g. `HasTypeDefinition`), which would pollute the discovery results with type nodes that are not data tags.

**Alternatives rejected:** 0:31 (References) — returns type nodes and non-navigable references. Causes false positives in tag discovery.

---

### Noise filter: skip Server, Types, Views, Aliases, StaticData, MyObjects
**Decision:** During node tree browse, skip any node whose BrowseName contains: `Server`, `Types`, `Views`, `Aliases`, `StaticData`, `MyObjects`.

**Rationale:** These are OPC-UA infrastructure nodes, not factory data. Including them would add dozens of non-actionable nodes (type definitions, server diagnostics, namespace metadata) to every discovery result. Filtering them at browse time is more efficient than filtering post-discovery.

**Alternatives rejected:** Post-discovery filter — wastes browse time and memory on useless nodes. No filter — unusable noise in tag list.

---

### DataType read: `AttributeIDDataType` + `AttributeIDValue` in single `ReadRequest`
**Decision:** Read both `AttributeIDDataType` and `AttributeIDValue` in a single `ReadRequest` per tag (two AttributeIDs in the same request).

**Rationale:** Halves the number of round-trips to the OPC-UA server during full tag discovery. On a 500-tag tree, this saves ~500 extra RPCs. The OPC-UA spec allows batching multiple AttributeIDs in one ReadRequest — the server returns one ReadValueID result per attribute.

**Alternatives rejected:** Two separate ReadRequests per tag — doubles round-trips, unnecessary overhead.

---

### Tag subscription: 500ms sampling interval via MonitoredItems
**Decision:** Subscribe to all discovered tags via MonitoredItems with a 500ms requested sampling interval.

**Rationale:** 500ms is the right tradeoff for the POC use cases:
- Fast enough to catch micro-stop transitions (shortest micro-stop threshold = 30s — 500ms gives 60 samples per event minimum).
- Slow enough to not flood SQLite storage or the rules engine with unnecessary samples.
- Below the OPC-UA default PublishingInterval in most SCADA configurations.

**Alternatives rejected:** 100ms — overkill for micro-stop detection, high write load on SQLite. 1000ms — might miss fast transient states on rapid jam sensors.

---

### `WatchForChanges()`: 20s polling interval for topology changes
**Decision:** `WatchForChanges()` re-browses the full node tree every 20 seconds to detect added or removed tags, and calls `onChange` only when the diff is non-empty.

**Rationale:** Factory OT topology changes are rare (new sensor added, SCADA tag renamed) but we want to detect them without manual restart. 20s is long enough to not impact Edge Agent CPU, short enough to catch changes within one scan cycle. Zero-overhead when nothing changes.

**Alternatives rejected:** On-demand only (no periodic scan) — would miss topology changes. 5s polling — unnecessary CPU load for rare events.

---

## Architecture Principles (Locked — all phases)

---

### Read-only on all source systems
**Decision:** The Edge Agent never writes to any PLC, SCADA, ERP, or MES.

**Rationale:** Core security promise to IT/OT security teams. Zero risk of corrupting production data or triggering unintended machine actions. Required for NIS2 compliance.

---

### Push-only to cloud
**Decision:** Only outbound HTTPS (TLS 1.3 + mTLS). Zero inbound open ports on client network.

**Rationale:** Direct sales argument to IT security and CISOs. Client network firewall requires zero rule changes — only standard HTTPS outbound. Eliminates the attack surface of an inbound-listening agent.

---

### Zero raw data in cloud
**Decision:** Only transformed events (micro-stop events, energy waste events, cost summaries) and aggregated KG snapshots leave the Edge. No raw OPC-UA tag values, no raw Modbus register values, ever sent to cloud.

**Rationale:** RGPD + industrial data sovereignty. Client retains full control of raw operational data. Cloud footprint stays minimal (Scaleway PLAY2-NANO ~3.99€/month). Enables BYOC option.

---

### No third-party middleware (no Kepware)
**Decision:** Native protocol drivers only — no Kepware, no OPC-Router, no third-party OT middleware.

**Rationale:** Direct sales argument to IT security teams (no additional attack surface). Removes vendor dependency and per-tag licensing costs (Kepware charges per tag count). Enables zero-dev deployment — client doesn't need to install, license, or maintain a separate middleware product.

**Alternatives rejected:** Kepware — licensing cost per tag, vendor lock-in, sales friction with security teams.

---

### ISA-95 as UNS ontology
**Decision:** The Unified Namespace follows the ISA-95 hierarchy: Enterprise → Site → Area → Work Center → Work Unit → Tag.

**Rationale:** ISA-95 is the manufacturing industry standard for hierarchical data modeling. Ensures maximum interoperability with third-party AI agents and ERP connectors. Aligns with EU Data Act positioning. Adopted by all major industrial automation vendors.

---

### Deterministic rules engine over ML for micro-stop detection
**Decision:** Micro-stop detection is implemented as a deterministic threshold-based rules engine, not a machine learning model.

**Rationale:** Deterministic rules are auditable, predictable, and require zero training data. Plant Managers can understand and adjust thresholds without data science expertise. ML would require labeled historical data that new clients don't have on Day 1. Rules are sufficient for the well-defined patterns (Run→Stop→Run, duration window, energy threshold).

**ML reserved for:** Tag classification (SLM Phi-3 — opaque tag names), predictive failure prediction (V2+).

---

*Internal document — Mindset Data — Confidential*
*Last updated: June 2026*
*License: Apache 2.0*
