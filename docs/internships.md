# MindSet — Intern Recruitment Catalog (2-month free stages)

> **Context.** Pre-funding, no payroll possible. French *stage convention* allows free internships up to 2 months (>60 days triggers mandatory gratification ~600€/mo). Window: 8 weeks max, July-September 2026.
> **Capacity constraint** (Entry 34 of analysis_log): Mohamed (solo engineer) realistically supervises 1-2 technical interns at ~30% of his time. Cécilia can supervise 1-2 non-technical interns. **Sweet spot: 2 interns total.** 3+ erodes Mohamed's V1 throughput too much.
> **Realistic productive window per intern**: ~5-6 weeks (weeks 1-2 = onboarding, week 8 = handover). Plan scope accordingly.
> **Last updated**: 2026-06-30

---

## Quick reference — full intern type catalog

| # | Type | Tier | Supervisor | Mentorship load | V1 impact |
|---|---|---|---|---|---|
| 1 | Frontend Developer (React) | **1** | Mohamed | Low (~10%) | HIGH — unblocks dashboard tabs |
| 2 | Marketing Site / Next.js Dev | **1** | Mohamed | Low (~10%) | HIGH — mindsetdata.io live |
| 3 | DevOps / SRE / Platform Eng | **1** | Mohamed | Low (~10%) | HIGH — V1 security additions (CI/CD, signed binaries, SBOM) |
| 4 | Sales Development / BDR | **1** | Cécilia | Medium (~25%) | HIGH — first pilot pipeline |
| 5 | Customer Discovery / Research | **1** | Cécilia | Medium (~25%) | HIGH — V1.5 use-case prioritization data |
| 6 | Go Backend Developer (junior) | **2** | Mohamed | Medium (~25%) | MEDIUM — 1 specific connector ships faster |
| 7 | Test Engineer / QA | **2** | Mohamed | Medium (~20%) | MEDIUM — V1 quality baseline |
| 8 | AI / ML / NLP intern | **2** | Mohamed | High (~35%) | MEDIUM-HIGH — Phi-3 tag classifier OR MCP server core |
| 9 | UX / UI Designer | **2** | Both | Low (~15%) | MEDIUM — pitch deck + dashboard polish |
| 10 | Marketing / Content | **2** | Cécilia | Low (~10%) | MEDIUM — X + LinkedIn cadence |
| 11 | Industrial Software Eng (OT-savvy) | **2** | Mohamed | High (~35%) | HIGH if found — but rare profile |
| 12 | Security / Pentesting | **3** | Mohamed | Medium (~25%) | LOW for V1 ship, HIGH for pharma readiness V1.5 |
| 13 | Product Designer (UX research) | **3** | Cécilia | Medium (~20%) | MEDIUM — V1.5 design |
| 14 | Business Analyst / Strategy | **3** | Cécilia | Low (~15%) | MEDIUM — competitive deep-dive + pricing data |
| 15 | Legal / Compliance | **3** | Cécilia | Low (~10%) | LOW for V1, HIGH if pharma deal lands |
| 16 | Technical Writer / Documentation | **3** | Both | Low (~10%) | LOW for V1, HIGH for V1.5 customer onboarding |

---

# TIER 1 — most impactful for V1 ship (recruit one or two of these first)

---

## 1. Frontend Developer (React)

**Why now**: V1 dashboard components (Gantt, Pareto, OEE/TRS, ROI simulator) are bottlenecks if Mohamed builds them solo. Self-contained, well-specified work.

**8-week scope**:
- Build 2-3 dashboard components from the V1 inventory (U5 Gantt timeline + U6 Pareto + U7 OEE view + U9 Tribal knowledge dropdown)
- Polish the dashboard skeleton (U4) for first-pilot demoability
- Integration with the existing React + Vite + Tailwind + Zustand stack

**Concrete deliverable**: 2-3 production-quality dashboard components shipped + integrated into `frontend/` codebase + their unit tests + screenshots in the pitch deck.

**Profile sought**:
- Strong React experience (3+ real projects in portfolio)
- Tailwind CSS comfort
- Has shipped charts / data-viz (Chart.js, Recharts, D3, ApexCharts)
- French OR English communication (Mohamed bilingual)

**Source schools (FR)**:
- **EPITECH** (Paris/Lyon) — strong dev culture, project-based
- **EPITA** (Paris) — strong on web/frontend
- **École 42** (Paris) — peer-learning, portfolio-driven students
- **HETIC** (Paris) — web + media specialty

**Supervisor**: Mohamed (technical review weekly).
**Mentorship load**: ~10% of Mohamed's time. Lowest of all engineering options because frontend is well-separated from his Go backend work.

**Risk if poorly executed**: dashboard tabs ship late or buggy → first-pilot demo lands poorly.

---

## 2. Marketing Site / Next.js Developer

**Why now**: `mindsetdata.io` doesn't exist yet but it's mentioned in the V1 roadmap + the X strategy (LinkedIn URL in bios should point somewhere). No founder time available to build it.

**8-week scope**:
- Build complete Next.js site on Vercel: landing + /product + /use-cases + /security + /contact (+ /download for V1.5 when product is ready)
- SEO baseline (meta tags, sitemap, structured data)
- Connect to a form backend (Typeform / Formspree / Plausible analytics)
- Brand consistency with logo / colors

**Concrete deliverable**: `mindsetdata.io` LIVE on Vercel by week 6. Cécilia + Mohamed can share the URL in every conversation.

**Profile sought**:
- Next.js 14+ experience
- Tailwind + design sensibility (or works well with a designer)
- Has shipped a marketing/landing site (portfolio required)
- SEO basics

**Source schools (FR)**:
- **HETIC** (Paris) — web + design specialty
- **EPITECH** (Paris/Lyon)
- **École 42** (Paris)
- **Sup de Web / SUPINFO** alternates

**Supervisor**: Cécilia for content; Mohamed for technical deploy.
**Mentorship load**: ~10% of founder time total.

**Risk if poorly executed**: no public URL to point investors / customers / advisors / X followers to → credibility gap.

---

## 3. DevOps / SRE / Platform Engineer

**Why now**: V1 security additions (SEC2 signed binaries + SEC3 SBOM + CVE scanning) are blockers for the security pitch. Pharma/cosmetics need this V1.5 latest. Setting up the CI/CD pipeline now means every release auto-signs + SBOMs from day 1.

**8-week scope**:
- GitHub Actions pipeline: build → test → sign → SBOM → CVE scan → push to private registry
- Cosign + Sigstore integration for binary signatures
- CycloneDX SBOM auto-generation
- Trivy or Snyk for CVE scanning (gate merges by severity)
- Docker multi-arch builds (amd64 + arm64)
- Setup the customer-side private registry (`registry.mindsetdata.io`)
- Documentation for "how to update the Edge Agent" customer-facing

**Concrete deliverable**: Every `git push main` triggers signed-binary + SBOM + CVE-scanned Docker image published to private registry. CI green badge in README.

**Profile sought**:
- Experience with GitHub Actions OR GitLab CI
- Familiar with Docker, containerd
- Has worked with cosign / sigstore / SBOMs (rare in interns — even basic familiarity is enough)
- Linux comfort, shell scripting

**Source schools (FR)**:
- **Polytechnique** (Mohamed's network — DIRECT INTRO POSSIBLE)
- **INSA** (Lyon, Toulouse, Rennes)
- **Centrale** (Paris, Lyon)
- **Mines** (Paris, Saint-Étienne)
- **Telecom ParisTech** (cybersecurity adjacent)

**Supervisor**: Mohamed (architecture decisions).
**Mentorship load**: ~10% of Mohamed's time. Highly bounded work.

**Risk if poorly executed**: V1 ships without supply-chain security → can't credibly pitch pharma/cosmetics → V1.5 deals delayed.

---

## 4. Sales Development / Business Development Representative (BDR)

**Why now**: First pilot customers must be identified BEFORE V1 ships, not after. Cécilia is the lead salesperson but has 100 other things. A BDR builds the target list + does cold outreach + qualifies leads.

**8-week scope**:
- Build a targeted list of 100 ETI customers across the 4 verticals (pharma, cosmetics, agrifood, metallurgy) — France-based
- Per-target research: company size, plant locations, key contact names + emails, recent news (M&A, expansion, regulation pressure)
- Cold outreach via LinkedIn + email (Cécilia approves messaging)
- Qualify responses → book demos for Cécilia
- Track in a CRM (HubSpot free tier or Airtable)

**Concrete deliverable**: 100-account list + 50 first-contact touches sent + 5-10 qualified demo bookings (target: 5%-10% reply rate, 10%-20% of replies book).

**Profile sought**:
- Business school undergrad or recent grad
- Energetic, comfortable on phone + LinkedIn
- French native + English working
- Ideally industrial / B2B internship background already

**Source schools (FR)**:
- **EM Lyon** (entrepreneurship track)
- **ESCP** (B2B sales)
- **EDHEC** (Cécilia's network — DIRECT INTRO POSSIBLE)
- **HEC** (Bachelor / Master Sales)
- **Audencia, Skema, Kedge** (regional B-schools)

**Supervisor**: Cécilia.
**Mentorship load**: ~25% of Cécilia's time (1-on-1 daily check-ins first week, weekly afterward, all outbound messaging reviewed before send).

**Risk if poorly executed**: pipeline empty when V1 ships → 3-month sales delay → cash crunch pre-Series A.

---

## 5. Customer Discovery / Research Intern

**Why now**: V1.5 use-case prioritization (4th + 5th starter templates beyond micro-stop / energy / OEE) depends on REAL customer signal — not founder guesses. This intern collects the signal.

**8-week scope**:
- Conduct 10-15 structured 1-hour interviews with Plant Managers + IT/OT Managers in the 4 target verticals
- Mix: 4-5 in agrifood (easiest access), 3-4 in metallurgy, 2-3 in pharma (harder), 1-2 in cosmetics
- Structured interview guide (Cécilia helps craft) covering: current OEE measurement, biggest operational pains, current IT/OT integration state, willingness to pay
- Synthesize findings: pain Pareto, vertical-specific quirks, "hair on fire" use cases beyond the 3 starter templates
- Build a vertical-specific pitch one-pager per vertical based on findings

**Concrete deliverable**: structured research report + 10-15 verbatim transcripts + ranked Pareto of pains + 4 vertical-specific pitch one-pagers.

**Profile sought**:
- Curious, asks good follow-up questions
- French native (for the interviews)
- Business school OR sociology / STS background
- Comfortable cold-calling factories

**Source schools (FR)**:
- **EDHEC** (Cécilia's network)
- **HEC** (business research track)
- **Sciences Po** (Master Sociology, applied research)
- **ESCP** (B2B research)
- **Mines Paris** (industrial engineering with field research)

**Supervisor**: Cécilia.
**Mentorship load**: ~25% of Cécilia's time (interview prep, debriefs, synthesis review).

**Risk if poorly executed**: V1.5 builds in the dark → wrong template choices → wasted engineering quarter post-funding.

---

# TIER 2 — high value, defer if bandwidth limited

---

## 6. Go Backend Developer (junior)

**Why now**: V1 connector work (Modbus TCP, Files/FTP) is well-specified. A junior Go dev can own a specific connector end-to-end while Mohamed builds critical-path (OF-state Fuzzy Join, Impact Engine).

**8-week scope**: own ONE of:
- **Option A**: Modbus TCP connector + device fingerprint DB (C2 in V1 inventory) — 8 weeks tight but doable
- **Option B**: Files / FTP / SFTP connector (C6) — easier, 5-6 weeks
- **Option C**: Help with the SQL connector multi-dialect (C3) under Mohamed's lead architecture — supportive role

**Concrete deliverable**: 1 production-ready connector with tests + integrated into the function registry.

**Profile sought**: 1-2 years Go (or Rust/C++ transferable), comfortable with `database/sql` and binary protocols.

**Source schools (FR)**: Polytechnique, INSA, Centrale, ENSIMAG (Grenoble), ENSEEIHT (Toulouse), Mines.

**Supervisor**: Mohamed.
**Mentorship load**: ~25% — Mohamed needs to review architecture + code quality carefully on a Go intern.

**Risk if poorly executed**: bad Go code that ships into V1 → maintenance debt + sovereignty incident risk.

---

## 7. Test Engineer / QA Engineer

**Why now**: V1 ships without an integration test harness today. Pre-pilot, this becomes urgent. After a paying customer, even more.

**8-week scope**:
- Build the integration test harness: dockerized OPC-UA simulator (Prosys), mock Modbus device, mock ERP (PostgreSQL with sample schema), mock MES
- Write 30-50 integration tests covering V1 connector contracts
- Set up test data fixtures per vertical
- CI integration (handoff to DevOps intern if both are hired)
- Document a manual QA checklist for the first-pilot install

**Concrete deliverable**: green CI badge + 30+ integration tests + manual QA checklist for pilot install.

**Profile sought**: comfortable with Docker, basic Go OR Python for tests, attention to detail.

**Source schools (FR)**: EPITA (engineering track), INSA, Centrale, university IT programs.

**Supervisor**: Mohamed.
**Mentorship load**: ~20% — moderate review on testing patterns.

**Risk if poorly executed**: first-pilot install reveals bugs that should have been caught in CI → customer trust loss.

---

## 8. AI / ML / NLP Intern

**Why now**: Phi-3 tag classifier (D5 + A1 in inventory) needs prompt engineering + eval. The MCP server schema design (A2) is open architecture work. Either is high-value but high-mentorship.

**8-week scope**: ONE of:
- **Option A** — Phi-3 tag classifier: prompt engineering + eval harness, get 70%+ accuracy on tag semantic classification + behavioral inference
- **Option B** — MCP server schema design: implement V1 MCP tools (`kg_query`, `kg_list_events`, `kg_cost_summary`), test against Claude Desktop, document schema
- **Option C** — Ad-hoc Analyst agent: design prompt + tool-use pattern for the Q&A agent, get demoable answers grounded in KG via MCP

**Concrete deliverable**: one of the V1 AI components shipped + eval results + demo video.

**Profile sought**:
- ML/AI master student OR strong AI side projects
- Familiar with Ollama, prompt engineering, eval harnesses (BLEU, BERTScore, custom)
- Python OR Go comfort
- For Option B: comfortable with JSON-RPC, MCP spec

**Source schools (FR)**:
- **Polytechnique IA track**
- **Mines AI Paris**
- **ENS Paris-Saclay** (ML master)
- **PSL / Université Paris-Dauphine** AI master
- **Sorbonne** AI master (M2 IASD)
- **INSA Lyon** AI track

**Supervisor**: Mohamed.
**Mentorship load**: **~35%** — highest of all options. AI work has many wrong paths, intern needs heavy guidance + eval discipline.

**Risk if poorly executed**: AI demo fails on day 1 of pilot → "AI-native" pitch collapses for that customer.

---


## 9. Marketing / Content Intern

**Why now**: X + LinkedIn cadence (`x_strategy.md`) requires consistent posting. Mohamed + Cécilia can't do 3-5 posts/week each indefinitely without help. An intern drafts → founder reviews → publishes.

**8-week scope**:
- Build content calendar for X + LinkedIn (8 weeks pre-scheduled)
- Draft 16-24 LinkedIn posts (Cécilia voice) + 8-12 X threads (Mohamed voice) for founders to edit
- Write 2-3 blog posts for `mindsetdata.io` (SEO-optimized — on EU sovereignty, OEE, edge MCP)
- SEO keyword research for the marketing site
- Manage X engagement (replies, monitoring, follows) under Mohamed's supervision
- LinkedIn outreach support for the BDR

**Concrete deliverable**: 8 weeks of pre-drafted content + 2-3 blog posts + SEO keyword list + 50+ X follows.

**Profile sought**:
- Strong writer (French + English), tech-curious
- Has run a brand's social media (portfolio: examples)
- Comfortable in B2B context

**Source schools (FR)**:
- **CELSA** (communication, Sorbonne)
- **EFAP** (communication + brand)
- **ESCP, EDHEC** (Cécilia's network — marketing tracks)
- **HEC** (marketing)
- **École W** (newer communication school)

**Supervisor**: Cécilia (content); Mohamed (technical X threads).
**Mentorship load**: ~10% — drafts come back for edit, founders publish.

**Risk if poorly executed**: bland content that doesn't capture the founder voices → low engagement, low signal.

---

## 10. Industrial Software Engineer (OT-savvy) — HARD to find

**Why now**: The OF-state Fuzzy Join + ERP schema templates per vertical NEED domain knowledge. A generic Go dev struggles with OF/recipe/batch concepts. This profile is RARE in interns but if found = high value.

**8-week scope**:
- Work alongside Mohamed on the OF-state Fuzzy Join engine
- Build initial ERP schema templates (Sage X3 first — most accessible) with deep field mapping
- Document the "active production context" abstraction for vertical-specific extensions
- Help design the Impact Engine integration

**Concrete deliverable**: OF-state Fuzzy Join working with Sage X3 + documented schema mapping + one vertical config template.

**Profile sought** (RARE):
- Operations research / industrial engineering background
- Knows what an OF / batch / recipe / MES is
- Programming experience (Go OR Python)
- Likely has done a previous stage in a factory or industrial software vendor

**Source schools (FR)**:
- **Polytechnique** (industrial engineering specialization)
- **Mines** (Paris, Saint-Étienne) — strong industrial program
- **Arts et Métiers (ENSAM)** — manufacturing-specialised
- **Centrale Lille** (industrial engineering)
- **ENSGSI** (Nancy — industrial engineering)
- **IMT Atlantique** (Nantes — manufacturing)

**Supervisor**: Mohamed.
**Mentorship load**: ~35% — Mohamed must architect, intern must learn industrial terminology fast.

**Risk if poorly executed**: muddled OF-state design that needs full rewrite → V1 critical path slips.

**Note on availability**: this profile is RARE. Realistic to find ONE candidate in the next 2 weeks via Mohamed's Polytechnique network + Mines + Arts et Métiers career offices. If not found, fall back to a generic Go intern (#6) under Mohamed's domain-architect lead.

---

# TIER 3 — useful but lower urgency

---

## 11. Security / Pentesting Intern

**Why now**: V1 isn't deeply security-mature yet (security framework decision Entry 20 still pending). For pharma/cosmetics V1.5+, baseline security work needs to start now. Less urgent for V1 ship.

**8-week scope**:
- Threat model for the V1 architecture (STRIDE-based)
- Baseline pentest on V0 stack (OPC-UA endpoint exposure, MQTT broker, REST API, dashboard)
- Write `SECURITY.md` + vulnerability disclosure policy
- Initial audit log architecture design (SEC4)
- ISO 27001 readiness gap analysis

**Concrete deliverable**: threat model doc + pentest report + SECURITY.md + audit log spec + ISO 27001 gap report.

**Profile sought**:
- Cybersecurity school student
- Familiar with OWASP, STRIDE, basic offensive security
- Understanding of supply-chain security (cosign, SBOM)

**Source schools (FR)**:
- **EPITA SRS** (Sécurité Réseaux Systèmes)
- **ESIEA Cybersécurité**
- **Telecom SudParis** cyber track
- **INSA Cybersécurité**
- **ENSIBS** (Vannes — cyber specialisation)

**Supervisor**: Mohamed (occasionally Edmond Tahar / Boost10x specialist for legal touchpoints).
**Mentorship load**: ~25%.

**Risk if poorly executed**: V1 ships with security gaps surfaced later by a real pentest → reputation damage with a regulated customer.

---

## 12. Business Analyst / Strategy Intern

**Why now**: Pricing model is open (the biggest investor-deck blocker). Competitive intelligence on Cognite / UMH / MaestroHub / Braincube can be deeper. TAM/SAM/SOM modeling can be refined with data.

**8-week scope**:
- Competitive product deep-dive: install / trial / map every feature of UMH community edition, MaestroHub free tier, document Cognite Atlas AI capabilities
- Pricing benchmark research: willingness-to-pay interviews with 5-10 prospects across the 4 verticals
- TAM/SAM/SOM modeling per vertical with real FR + EU industrial data
- Pricing model recommendation backed by data

**Concrete deliverable**: competitive analysis report + pricing benchmark with data + TAM/SAM/SOM model + final pricing recommendation memo.

**Profile sought**: business school strategy major OR consulting-track student.

**Source schools**: **HEC**, **ESCP**, **EDHEC**, **EM Lyon**, **Sciences Po**.

**Supervisor**: Cécilia (with input from Boost10x's Djamil on pricing).
**Mentorship load**: ~15%.

**Risk if poorly executed**: pricing decision still gut-feel post-internship → investor pitch weakness persists.

---

# Recommended selection — top 2 (best fit for "right now")

**My pick for 2 interns**:

| Pick | Role | Why this combo |
|---|---|---|
| **A** | **Frontend Developer (React)** — type #1 | Unblocks dashboard tabs (Gantt/Pareto/OEE/Tribal Knowledge UI). Mohamed's lowest supervision cost on the engineering side. Highest near-term V1 demo impact. |
| **B** | **Sales Development / BDR** — type #4 | Builds pilot pipeline so customers are queued when V1 ships. Cécilia-led, complementary to A. No engineering supervision cost. |

**Combined mentorship load**: ~10% Mohamed + ~25% Cécilia = both founders preserve ~75%+ of own building time.

**Alternative pick if you want both engineering**:
- A: Frontend Developer (type #1)
- B: DevOps / SRE (type #3) — also low Mohamed-supervision, ships V1 security additions

**Alternative pick if Cécilia wants 2 non-engineering**:
- A: BDR (type #4) — pipeline
- B: Customer Discovery (type #5) — V1.5 use-case signal

**My strong recommendation: don't go 3+**. The math from Entry 34: 3 interns = ~50% mgmt time for Mohamed if 2 are technical → he stops building V1 critical path → V1 slips.

---

# Recruitment timeline (starts today, week of 2026-06-30)

| Week | Action |
|---|---|
| **This week** (June 30 - July 4) | Mohamed pings Polytechnique career office + Cécilia pings EDHEC career office. Post stage offers on EPITA / EPITECH / 42 / HETIC career portals. WhatsApp Jalil/Djamil for Boost10x network referrals. |
| **Week of July 7** | Screen CVs (target: 5-10 per role). Cécilia owns scheduling first-round 30-min screens. |
| **Week of July 14** | 30-min first screens (Cécilia + Mohamed depending on role). Filter to 2-3 finalists per role. |
| **Week of July 21** | Final-round (60 min): Mohamed for technical, Cécilia for non-technical. Boost10x specialist (Maureen Rousseau) reviews intern convention templates. |
| **Week of July 28** | Offers + convention de stage signed. |
| **Week of August 4** | **Interns START.** 8-week stage runs Aug 4 → Sep 26. Productive output: Aug 18 → Sep 19 (~5 weeks). |
| **Week of September 22-26** | Handover docs + wrap-up. |

---

# Selection guidance — how to pick

Ask yourself in order:

1. **What's the #1 V1 risk this month?**
   - Dashboard not demo-ready → Frontend dev (#1)
   - No marketing site → Next.js dev (#2)
   - No CI/CD or signed binaries → DevOps (#3)
   - No customer pipeline → BDR (#4)
   - Wrong use-cases → Customer Discovery (#5)

2. **Whose bandwidth is the bottleneck — Mohamed or Cécilia?**
   - Mohamed bottlenecked → pick a low-supervision-cost engineering intern (#1 or #3 — both ~10% mentorship)
   - Cécilia bottlenecked → pick BDR (#4) or Customer Discovery (#5)
   - Both → mix one engineering + one non-engineering

3. **What's the path-of-least-resistance for finding the candidate?**
   - Mohamed's Polytechnique network → DevOps (#3) or Industrial Software (#11)
   - Cécilia's EDHEC network → BDR (#4), Customer Discovery (#5), Marketing (#10), Strategy (#14)
   - Public portals (EPITA / EPITECH / 42) → Frontend (#1), Marketing Site (#2)

4. **Are any of these BLOCKING the pitch deck / fundraise?**
   - Pricing decision → Strategy intern (#14) collects data → unblocks
   - No customer signal → Customer Discovery (#5)
   - Weak deck design → UX/UI Designer (#9)

5. **What can wait until V1.5?**
   - Security pentest (#12) — useful but not V1-critical
   - Product designer / UX research (#13) — V1.5 refinement
   - Legal (#15) — only fires when first customer signs
   - Technical writer (#16) — V1.5 customer onboarding wave

---

## Bottom line

| | |
|---|---|
| **Total roles in catalog** | 16 |
| **Tier 1 (most impactful for V1)** | 5 |
| **Recommended pick (Mohamed's lens)** | Frontend Dev (#1) + BDR (#4) — balanced engineering + GTM, lowest combined supervision cost |
| **Recommended pick (engineering-only)** | Frontend Dev (#1) + DevOps (#3) — V1 demo + V1 security infra |
| **Recommended pick (Cécilia-only)** | BDR (#4) + Customer Discovery (#5) — pipeline + signal |
| **Hard rule** | Max 2 interns. 3+ erodes Mohamed's V1 throughput. |
| **Best source for FR engineering interns** | Polytechnique (Mohamed's network) → INSA → Centrale → EPITA → EPITECH |
| **Best source for FR business interns** | EDHEC (Cécilia's network) → ESCP → HEC → EM Lyon → Sciences Po |
| **Stage start date target** | August 4, 2026 |
| **Productive output window per intern** | ~5 weeks |
