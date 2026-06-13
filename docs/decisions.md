# Mindset Data — Technical Decisions Log

> Running log of all architectural and implementation decisions.
> One entry per decision, newest first.
> Format: **Decision** → **Rationale** → **Alternatives rejected**

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
