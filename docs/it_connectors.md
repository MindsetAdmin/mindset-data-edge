# MindSet — IT-Side Connectors

> **Strategy + spec.** The IT connector layer that powers MindSet's OT/IT reconciliation moat + the Impact Engine (`docs/impact_engine.md`).
> **In scope**: ERP · MES · CMMS · LIMS · QMS · Energy Management · L3 Historians.
> **Out of scope**: OT protocols (OPC-UA / Modbus / S7 / MQTT) — see `mindset.md` §5.
> **Last updated**: 2026-06-30

---

## Why IT connectors are NOT just "another connector"

OT connectors talk to PLCs and SCADA on the shop floor. IT connectors talk to the customer's enterprise systems. These look similar from a code perspective but they're radically different in practice.

| Dimension | OT connector | IT connector |
|---|---|---|
| **Customer gatekeeper** | Plant Manager / Maintenance | IT Director / CISO / IT Procurement |
| **Time to first connect** | Minutes (OPC-UA endpoint + IP) | Days to months (service account · firewall · security review · schema mapping) |
| **Read access** | Trivially granted by Plant Manager | Often requires formal security review |
| **Schema** | OPC-UA self-describes; Modbus has fingerprint DB | Heterogeneous per customer (no two SAP installations identical) |
| **Failure cost** | If we break it: production stops | If we break it: corrupted master data, lost audit trail, security incident |
| **Write capability** | Architecture prohibits all writes (Decision: read-only) | Same — never write. Hard rule. |
| **Refresh rate** | Real-time (ms-to-s) | Periodic polling (seconds-to-minutes) |
| **Reliability** | Edge buffer if OT drops | Cloud buffer if IT API drops |
| **Sovereignty risk** | OT data never leaves customer | IT data may include PII (operator names, customer names) — handle carefully |

**The implication**: IT connector strategy is as much organizational as technical. The customer's IT team can kill the deal if access is hard to grant. The architecture must make THEM say yes easily.

---

## The IT systems landscape (everything we may need to integrate with)

### Category A — ERP (orders · customers · margins · planning)

| System | Notes | FR mid-market prevalence |
|---|---|---|
| **SAP ECC** (on-prem) | Dominant in pharma + cosmetics + metallurgy enterprise | High |
| **SAP S/4HANA** (cloud/on-prem) | Modern SAP — REST + OData APIs | Growing |
| **Sage X3** | Common FR mid-market manufacturing | Very high |
| **Sage 100** | FR SMB | High |
| **Microsoft Dynamics 365** | Multi-vertical | Medium |
| **Odoo** | Modern FR SMB — open source | Growing |
| **Oracle Cloud / EBS** | Large enterprise, niche in FR mid-market | Low-medium |
| **Custom in-house ERPs** | Long tail | Variable |

### Category B — MES (recipes · WIP · quality · scheduling — Level 3)

| System | Notes | Vertical strength |
|---|---|---|
| **Werum PAS-X (Körber)** | Dominant pharma MES — mandatory in regulated pharma | Pharma · Cosmetics |
| **SAP MII** | Large enterprise — sits between SAP ERP and shop floor | Pharma · Cosmetics · Metallurgy |
| **AVEVA MES** (ex-Wonderware) | Multi-vertical | Metallurgy · Energy-intensive |
| **Siemens Opcenter** (ex-Camstar) | Modern Siemens MES | Pharma · Discrete manufacturing |
| **Plex (Rockwell)** | Cloud MES for mid-market | Mixed |
| **PSI Metals** | Specialised metallurgy MES | Metallurgy (dominant) |
| **Tulip** | Modern low-code MES | Agile mid-market |
| **iBASEt / Apriso (Dassault)** | Aerospace + complex discrete | Niche for us |
| **Aptean / CSB-System** | Agrifood-specialised | Agrifood |
| **Custom MES / no MES** | Long tail — many smaller agrifood factories have NO MES | Agrifood (especially FR independents) |

### Category C — CMMS (maintenance management)

IBM Maximo · infor EAM · dimo Maint (FR — significant) · openMAINT · custom.

### Category D — LIMS (lab quality — pharma + cosmetics critical)

LabWare LIMS · STARLIMS · Labvantage.

### Category E — QMS (quality management — pharma mandatory)

Sparta TrackWise · MasterControl · Veeva Vault QualityDocs.

### Category F — Energy Management Systems

Schneider EcoStruxure Energy · Siemens SIMATIC Energy Suite · custom.

### Category G — L3 Historians

OSIsoft PI System (now AVEVA PI) · AVEVA Historian · InfluxDB · TimescaleDB · custom.

*(Often considered OT but accessed like IT — SQL or REST. Treated here.)*

---

## Connection methods (the technical layer)

### Method 1 — Direct SQL (V1)

**Covers**: SAP ECC (HANA/Oracle backend), Sage X3 (MSSQL), Sage 100 (MSSQL), Dynamics on-prem (MSSQL), Odoo (PostgreSQL), AVEVA Historian, TimescaleDB, custom DBs.

**Pros**: simple, fast, no API rate limits, customer's DBA can grant a read view in 30 min.

**Cons**: requires customer IT to grant DB read access (security review needed), schema knowledge required per customer, brittle to schema changes.

**V1 status**: ✅ **PostgreSQL + MSSQL + MySQL drivers** (locked decision). Oracle + HANA = V1.5.

### Method 2 — REST API (V1.5)

**Covers**: SAP S/4HANA · Dynamics 365 Cloud · Odoo · Tulip · most modern MES · Werum PAS-X (has REST endpoints) · Plex.

**Pros**: modern, well-documented, often token-auth, no direct DB access required (lower IT-team friction).

**Cons**: rate limits, may not expose everything via API, often requires customer IT to provision API key/service account.

**V1 status**: V1.5 generic REST connector.

### Method 3 — OData (V2)

**Covers**: SAP NetWeaver / Gateway · Dynamics 365 · SharePoint metadata.

**Pros**: standard, queryable, the "right" way to talk to SAP.

**Cons**: niche, extra learning curve, mostly only useful for SAP-heavy customers.

**V1 status**: V2 if SAP customer demand drives it.

### Method 4 — SOAP / XML (legacy)

**Covers**: legacy ERP installations · old SAP · some CMMS · IBM Maximo classic.

**Pros**: still works on legacy systems we can't avoid.

**Cons**: painful, XML-heavy, deprecated.

**V1 status**: **skip unless a paying customer specifically blocks on it.** Build only when forced.

### Method 5 — Files (CSV / Excel / JSON via FTP/SFTP/SMB)

**Covers**: ANY system that can export to flat files. Saves us when API access is impossible.

**Pros**: works for systems without API access · customer can manually email exports for the first pilot · lowest IT-team friction (no service account needed).

**Cons**: not real-time, manual setup, breaks if export format changes.

**V1 status**: V1.5 with FTP/SFTP polling.

### Method 6 — Proprietary SDKs

**Covers**: Werum PAS-X (Python SDK) · AVEVA (.NET integration) · some vendor-supplied connectors.

**Pros**: deeper functionality (full feature access) · vendor-supported · sometimes mandatory for certified integrations.

**Cons**: licensing fees · vendor lock-in · vendor pace controls our roadmap.

**V1 status**: V2+ — only build when a specific customer's deal depends on it. Werum is the most likely first one (pharma vertical).

---

## What we actually WANT from each IT system (data shopping list)

Driven by the OF-state Fuzzy Join + Impact Engine enrichments (`docs/impact_engine.md`).

### From ERP (V1 essentials)

| Field | Why |
|---|---|
| **Production orders** (OF) with status field | Fuzzy Join active-OF detection |
| OF → product mapping | Per-product margin lookup |
| OF → customer mapping + due date | Customer-commitment flag |
| Product master (`product_id`, name, margin per unit, theoretical cadence per line) | Per-product Impact Engine enrichment |
| Customer master (name, criticality) | Display name in dashboard |

**Note**: this is the SQL connector's V1 workload. Schema mapping per customer (configurable YAML — see "Configuration model" below).

### From MES (V1.5 essentials)

| Field | Why |
|---|---|
| Active batch / OF state (real-time) | Authoritative active-OF in pharma + cosmetics (more reliable than ERP) |
| Recipe per active OF | Per-product context, time-limit for spoilage (V2) |
| Real-time quality measurements (in-process) | Quality scrap on restart, quality-gap alerts |
| Historical scrap rate per stop type | Quality scrap V1.5 enrichment |
| Setup-time matrix (per product transition × machine) | Restart cost V1.5 enrichment (auto-replaces V1 manual config) |
| Operator-of-record per machine | Tribal knowledge capture context |
| Equipment availability status | Scheduling deviation detection |

### From CMMS (V1.5)

| Field | Why |
|---|---|
| Open maintenance work orders per machine | "This stop = maintenance planned, not a fault" filter |
| Asset criticality | Prioritization weight in Impact Engine |
| Historical MTBF / MTTR | Predictive scoring (V2) |

### From LIMS (V2 — pharma only)

| Field | Why |
|---|---|
| Sample results per batch | Quality scrap precision |
| In-process measurements (deviation alerts) | Trigger Impact Engine quality cost |

### From QMS (V2 — pharma)

| Field | Why |
|---|---|
| Compliance-flag rules per stop type | "This stop triggers a deviation report" detection |
| Validation status per equipment | Block events on non-validated equipment |

### From Energy Management Systems (V1.5)

| Field | Why |
|---|---|
| Real-time consumption per machine | Energy waste V1 enrichment + Level 2 (with OF context) |
| Tariff windows (peak hours) | Energy peak penalty V1.5 |
| Grid contract pricing | Peak cost calculation |

### From L3 Historians (V1.5)

Bidirectional: read historical OT context, write enriched events back.

---

## Per-vertical priority matrix (which IT systems matter most)

| Vertical | Priority 1 (V1) | Priority 2 (V1.5) | Priority 3 (V2) |
|---|---|---|---|
| 💊 **Pharma** | ERP (SAP MII often) + **Werum PAS-X** | LIMS + QMS + CMMS | Proprietary Werum SDK |
| 💄 **Cosmetics** | ERP (SAP/Dynamics) + MES (Werum/SAP MII) | LIMS + Energy Mgmt | Proprietary integrations |
| 🌾 **Agrifood** | ERP (Sage X3, Aptean) | MES if present (often absent in independents) + Energy Mgmt | Files/FTP for non-API customers |
| ⚙️ **Metallurgy** | ERP (SAP) + **PSI Metals MES** | Energy Mgmt + CMMS | LIMS for alloy quality |

**GTM implication**:
- **Agrifood first pilot is easiest** — ERP only, often Sage X3 (well-understood MSSQL backend), can ship with V1 SQL connector alone
- **Pharma + cosmetics demand MES integration before they'll consider us** — Werum PAS-X integration is a gating item for pharma deals (high effort, high payoff)
- **Metallurgy needs PSI Metals understanding** — niche but high-value

---

## Authentication + security model (the IT-team-says-yes model)

To make a CISO say yes:

1. **Read-only by architecture**. No write capability anywhere in the IT connector layer. `SELECT` queries only via SQL; `GET` only via REST. Architecturally enforced — not a config flag the customer trusts us to flip.
2. **Dedicated service account at customer's IT**. Customer creates a service user with minimum privileges (specific read views/tables only). MindSet never gets DBA access.
3. **mTLS for all REST APIs**. Customer-issued client certs preferred over API keys.
4. **Audit log of every query**. Every SQL query / REST call MindSet issues is logged with timestamp + source + result hash. Exportable to customer's SIEM (per security framework SEC4 — Entry 20).
5. **Secrets management (SOPS)**. Credentials encrypted at rest with customer-controlled key. No plaintext DB passwords on disk.
6. **Customer can revoke access at any time**. Single config change, no MindSet involvement needed.
7. **PII handling discipline**. IT systems contain operator names, customer names, supplier names — PII. **We tag PII at ingestion + apply pseudonymization for cloud upload + never include in MQTT/event payloads.**
8. **No data egress to MindSet cloud without explicit per-customer agreement.** On-Premise edition: zero egress. Hybrid: only transformed events go up (no raw IT records).
9. **Network architecture**: customer's IT firewall sees only outbound HTTPS from MindSet edge — no inbound. Same as the OT-side principle (Entry 5).

---

## Configuration model

**Per-customer YAML mapping**. Templates per vertical.

`config/it_connectors.yaml`:
```yaml
# Customer: AcmePharma — pharma vertical, SAP MII + Werum PAS-X
erp:
  type: sap_mii
  connection:
    method: rest                       # rest | sql | files | soap
    base_url: https://sap-mii.acme.internal
    auth: mtls
    cert_path: /etc/mindset/secrets/sap-client.pem
  schema:
    production_orders_endpoint: /odata/v4/ProductionOrders
    fields:
      of_id: OrderID
      product_id: MaterialID
      customer_id: SoldToParty
      due_date: RequestedDeliveryDate
      status: OrderStatus
    status_active_values: ["REL", "INPR"]
  poll_interval_seconds: 30

mes:
  type: werum_pas_x
  connection:
    method: rest                       # later: proprietary_sdk
    base_url: https://werum.acme.internal/api
    auth: api_key
    secret_ref: WERUM_API_KEY          # resolved via SOPS
  schema:
    active_batches_endpoint: /batches?state=in_process
    recipe_endpoint: /recipes/{batch_id}
  poll_interval_seconds: 15

cmms:
  type: maximo
  enabled: false                       # V1.5 — not yet

lims:
  enabled: false                       # V2 — pharma later

energy:
  type: schneider_ecostruxure
  enabled: false                       # V1.5
```

**Templates per vertical** ship in `config/templates/`:
- `pharma-werum.yaml` (Werum PAS-X + SAP MII)
- `agrifood-sage.yaml` (Sage X3 only)
- `metallurgy-psi.yaml` (PSI Metals + SAP)
- `cosmetics-werum.yaml` (Werum + Dynamics 365)
- `minimum-erp-only.yaml` (SQL ERP only — works for any customer)

**Auto-detection (V1.5+)**: scan customer's network for known IT-system signatures (SAP MII HTTP fingerprint, Werum login page, etc.) → suggest the template.

---

## Phased rollout

### V1 (with the SQL connector — currently locked) — ~3 weeks

| Item | Status |
|---|---|
| SQL connector — PostgreSQL · MSSQL · MySQL drivers | Decision locked, in scope |
| Per-customer schema YAML loader | NEW for V1 |
| `minimum-erp-only.yaml` template | NEW for V1 |
| Read-only enforcement (architectural) | NEW for V1 |
| SOPS secrets management for DB credentials | NEW for V1 |
| Audit log of queries issued | NEW for V1 (links to security framework) |

**V1 capability**: read OF + product master + customer master from ERPs that expose a SQL backend. Powers Fuzzy Join + 2 of the 4 V1 Impact Engine enrichments (per-product margin + customer-commitment flag).

### V1.5 — Q3 2027 — ~6-8 weeks

| Item | Notes |
|---|---|
| Generic REST connector | Covers modern ERP (S/4HANA, D365, Odoo cloud) + most modern MES |
| Generic MES read layer | Plug-in for Werum PAS-X REST endpoints, SAP MII OData, AVEVA MES |
| Files / FTP / SFTP / SMB connector | For customers with no API access |
| CMMS connector (Maximo + infor EAM) | Read open work orders + asset criticality |
| Energy Management connector (Schneider EcoStruxure) | Real-time consumption + tariff windows |
| Auto-detection of IT systems (network scan) | Suggest the right template |
| Oracle SQL driver | Adds Oracle ERP coverage |

### V2 — Q4 2027 — ~10-12 weeks

| Item | Notes |
|---|---|
| Werum PAS-X proprietary SDK integration | Pharma-vertical credibility gate — only when first pharma deal forces it |
| SAP HANA driver | Modern SAP-heavy customers |
| OData connector | SAP NetWeaver / Dynamics 365 deep integration |
| LIMS connectors (LabWare + STARLIMS) | Pharma quality data |
| QMS connectors (TrackWise + MasterControl) | Pharma compliance flags |
| Vertical templates expanded | Per-customer customization patterns extracted to reusable templates |

### V3+ (2028-2029)

- PSI Metals SDK (metallurgy)
- Aptean / CSB-System (agrifood mid-tier)
- Legacy SOAP connector (only on customer demand)
- Cloud connectors for IT-in-cloud customers (S3, Azure Blob — already deferred per sovereignty stance)

---

## Strategic question — generic vs vendor-specific connectors

Two approaches with real trade-offs:

| Approach | Pros | Cons | Recommendation |
|---|---|---|---|
| **Generic only** (SQL + REST + Files) | Faster to ship, broadest customer fit, no vendor lock-in for us, lower maintenance | Requires per-customer schema config, breaks on weird MES, no "turn-key" pharma demo | **V1 + V1.5 path** |
| **Vendor-specific connectors** (`werum.go`, `sap_mii.go`, `psi_metals.go`) | Turn-key for that vendor's customers, deeper integration possible, better pharma sales | High build cost per connector (weeks each), narrow applicability, vendor dependency, maintenance burden | **V2+ on customer demand only** |
| **Hybrid (recommended)** | Generic V1 ships broad coverage; vendor-specific V2 adds turn-key for high-priority verticals (pharma → Werum first) | Need to manage both code paths | ✅ **This is the plan** |

**Decision principle**: every vendor-specific connector must be paid for by a paying customer. No speculative builds.

---

## Open strategic questions

| Question | Decision needed by | Notes |
|---|---|---|
| Werum partnership? | Q1 2027 (first pharma deal) | Becoming a Werum-certified integration partner unlocks pharma sales but takes 6-12 months + €€. Decide after first pharma signal. |
| SAP partnership? | Q2 2027 | Similar — SAP partner program is real distribution but heavy. |
| Open the schema specs publicly? | 2028 | If MindSet shifts to open-core, publishing the schema-mapping YAMLs as community templates accelerates adoption. Hold for now (closed-source decision). |
| Per-vertical config templates: ship publicly or per-customer only? | V1.5 | Recommend ship `minimum-erp-only.yaml` + 2 vertical templates publicly with V1.5. Specific customer configs stay private. |
| Cloud-side IT ingestion (for multi-site reconciliation)? | V2 design | Currently all IT reads happen at the edge. If a customer's IT is centralized, we may need cloud-side ingestion that fans out to per-site edges. Architectural shift — flag for V2. |

---

## Critical implementation notes

### Where in the code

```
internal/
  connectors/
    sql/
      driver.go              # multi-dialect dispatcher (V1)
      postgres.go            # V1
      mssql.go               # V1
      mysql.go               # V1
      oracle.go              # V1.5
      hana.go                # V2
    rest/
      generic.go             # V1.5 — schema-driven REST
      odata.go               # V2 — SAP/Dynamics specific
    files/
      csv.go                 # V1.5
      excel.go               # V1.5
      ftp_sftp.go            # V1.5
    vendor_specific/
      werum_pas_x.go         # V2 — proprietary SDK wrapper
      psi_metals.go          # V3 — when metallurgy customer signs
      sap_mii_native.go      # V2 — beyond generic SQL/REST
  it_schema/
    loader.go                # YAML config loader
    fuzzy_join_input.go      # presents IT data to Fuzzy Join in a uniform shape
    audit_log.go             # logs every query issued (security framework SEC4)
  templates/
    minimum-erp-only.yaml
    pharma-werum.yaml
    agrifood-sage.yaml
    metallurgy-psi.yaml
    cosmetics-werum.yaml
```

### Integration with the rest of the platform

- **Fuzzy Join engine** (`internal/fuzzy/of_state.go`) — consumes a uniform "active production context" object that the IT connector layer assembles, regardless of source system
- **Impact Engine** (`internal/cost/*`) — consumes the same context for per-product margin + customer flag + setup cost lookups
- **UNS contextualization** (`internal/uns/*`) — can be enriched by MES recipe data when available
- **Pipeline engine** (`internal/pipeline/engine.go`) — IT connectors become functions in the function registry; pipelines can wire them into custom flows

### Testing strategy

- **V1 SQL connector**: integration tests against PostgreSQL + MSSQL + MySQL via Docker containers in CI
- **V1.5 REST connector**: contract tests against mock REST servers per vendor
- **V2 vendor-specific**: requires vendor sandbox access (Werum has one; SAP has one; PSI Metals has one) — engagement begins when first vertical customer signs

---

## Today's takeaway

| | |
|---|---|
| **V1 (next 3 weeks)** | SQL connector — PG/MSSQL/MySQL — reads ERP OF + product + customer + due-date. Powers Fuzzy Join + 2 V1 Impact Engine enrichments. |
| **V1.5 (Q3 2027)** | REST · MES · Files/FTP · CMMS · Energy Management · auto-detection. Powers full V1.5 Impact Engine. |
| **V2 (Q4 2027)** | Vendor-specific Werum + PSI Metals — only when paying customers demand it. LIMS + QMS for pharma compliance. |
| **V3+ (2028+)** | Cloud-side IT ingestion · legacy SOAP · partnership-driven SDK integrations. |
| **Hard rule** | **Read-only by architecture.** No writes. Ever. |
| **Hard rule** | **PII tagged at ingestion.** No raw IT records leave the customer network. Hybrid edition only sends pseudonymized aggregates. |
| **Strategic** | Hybrid generic+specific approach. Each vendor-specific connector built only on paying-customer demand. |
| **GTM unblocker** | The agrifood first pilot needs only the V1 SQL connector to a Sage X3 backend. We can start there. |

---

## Open questions for the user

1. **Pricing-relevant**: Do we want to charge differently for customers requiring vendor-specific connectors (Werum, PSI Metals)? The build cost is real.
2. **GTM-relevant**: Does the Werum partnership question (Q1 2027) become a Cécilia investigation now, given pharma is one of the 4 initial verticals?
3. **Resource-relevant**: SAP MII integration is heavy. Do we wait until a SAP customer signs, or pre-build during V1.5 to win pharma deals?
4. **Architecture-relevant**: Cloud-side IT ingestion vs edge-only — defer to V2 design or revisit now?
