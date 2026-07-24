# LinkedIn Profile — Recommended Updates

_Reviewed live on linkedin.com/in/mohamed-khenafif-52844335b, 2026-07-18. Scope: headline, About, Experience, all 5 posts in Activity, and the "Open to" setting — not a full audit (Skills, Recommendations, Certifications, and comment/reaction history weren't reviewed). Recommendations are grounded in the positioning already settled in `docs/Blurb_Invest.md` and the finance/shopfloor wedge decision in `docs/analysis_log.md` Entry 74, so this reinforces the outreach work rather than reinventing the message._

## Decisions made so far (2026-07-18)

- **Item 1 (Experience description) — declined.** MindSet Data hasn't officially launched; user doesn't want to write an official company description yet. Leaving the entry blank as-is.
- **Visibility level for headline/About — "light framing, no specifics."** Not the full €/silos pitch (that's what got declined in item 1, for the same reason) — a vague directional line only. Exact wording below in items 3–4, drafted but **not yet applied to the live profile** — pending final confirmation.
- **Featured section (pinning the MindSet post) — skipped**, deliberately. It carries the full pitch language, which pulls against "light framing, no specifics." Revisit once the company is officially announced.
- **Posts/Activity audit — done**, see new section below. Only one candidate for removal, and it's low-stakes.
- **Stealth Startup entry, URL cleanup — still pending execution**, not yet done on the live profile.
- **Headline and About — applied to the live profile.** User applied the agreed light-framing headline manually. While editing About, the original technical-bio paragraph was accidentally deleted, leaving only the new opening line. Restored verbatim from this doc's Appendix and re-saved (2026-07-18) — About now reads: the new opening line, then the full original bio paragraph, unchanged, exactly as before.
- **About didn't reflect current MindSet Data work — fixed.** User flagged that the bio (STM32/ESP32, UART/SPI) reads as if MindSet Data itself is a hardware/embedded shop, which it isn't. Fixed by (a) adding "Before that" + switching the old bio to past tense, so it reads as background, not current work, and (b) adding a new middle paragraph naming actual current technical scope (Go, OPC-UA/MQTT, SQL connectors, data pipelines — drawn from the real work in this repo, not invented). Both applied and live — see item 4 and the Appendix.
- **Top Skills — added.** Go (Programming Language), SQL, REST APIs, Data Pipelines added to the profile's Skills section, matching the new About paragraph. Applied and live.
- **Public profile URL — resolved on its own.** Checked after the other edits: it's now `linkedin.com/in/mohamed-khenafif` (clean, no numeric suffix) without anyone explicitly changing it — LinkedIn appears to have auto-claimed the vanity slug at some point during today's edits. No action needed on item 5 anymore.

---

## 1. The MindSet Data experience entry has no description — declined for now

Today it reads:

> **Co-founder** — MindSet Data — Jun 2026 – Present · 2 mos · France
> *(no description)*

This is the section anyone doing real due diligence — an investor, a warm lead from Cécilia's LinkedIn outreach, a prospective co-design partner — reads right after the headline. It's currently blank while your freelance and contract roles below it both have descriptions. That's backwards for a co-founder profile.

**Original suggestion** (condensed from `Blurb_Invest.md`, kept here for later — not being used now):

> MindSet Data is the transversal data infrastructure that unifies industrial operations — connecting OT and IT silos, translating every event into its economic stake, and telling teams what to act on first. Not another point tool: the foundation manufacturers build their AI strategy on.
>
> Building this with Cécilia Tran. If you're sitting on shopfloor or ERP data that isn't turning into decisions, let's talk.

**Decision:** not writing this yet — MindSet Data isn't officially launched, and the user doesn't want to describe the company publicly before that. Entry stays blank. Revisit at official launch.

## 2. Two overlapping entries for the same company — merge or remove "Stealth Startup"

Right below the MindSet Data entry:

> **Co-Founder** — Stealth Startup · Permanent — Feb 2026 – Present · 6 mos · France
> "Powering the infrastructure behind manufacturing. Stay tuned."

Same person, same country, overlapping dates, obviously the same venture before it had a public name. Now that MindSet Data is announced, having both listed side by side reads as sloppy bookkeeping, not as a stealth-to-launch story — a careful reader (which is exactly who you want reading a co-founder profile) will notice.

**Fix:** delete the "Stealth Startup" entry. If you want to preserve the real start date, edit MindSet Data's start date back to Feb 2026 instead of keeping both.

## 3. Headline says nothing about what you're building

Current:

> Embedded Software Engineer | IoT & Low-Power Systems | FreeRTOS, STM32, ESP32 | Firmware Development

Accurate, but it's pure job-title framing — no mission, no company. The headline is the single most-seen piece of text on LinkedIn (search results, comments, connection requests, the small text under your name everywhere) — right now it reads "engineer looking for embedded work," not "co-founder building X."

**Original suggestion** (full pitch — too specific given the "not officially launched" constraint, kept for later):

> Co-Founder & CTO @ MindSet Data — turning invisible shopfloor losses into € decisions | ex-embedded (FreeRTOS / STM32 / ESP32)

**Agreed direction: "light framing, no specifics."** No €/silos pitch — just a directional line naming the company, nothing about what it does yet:

> Co-Founder @ MindSet Data — building the future of industrial data | Embedded Software Engineer | FreeRTOS, STM32, ESP32

Status: **applied to the live profile.**

## 4. About section is 100% technical bio — no problem framing

The current About is a well-written firmware-engineer resume paragraph (FreeRTOS, embedded Linux, MQTT/TLS, low-power design) with zero mention of the problem MindSet solves. For a co-founder, About is usually the second thing read after the headline+experience — right now it undersells the reason the company exists.

**Original suggestion** (full problem framing pulled from `Blurb_Invest.md`'s "The Problem" paragraph — too specific for now, kept for later): open with 2–3 sentences on silos/invisible losses/decision latency, then pivot into the existing technical bio as evidence.

**Agreed direction: same light framing as the headline.** One new opening line, existing bio paragraph unchanged below it. Status: **applied to the live profile.** (Note: the original bio paragraph was briefly lost mid-edit — user deleted it by mistake while applying this — and restored verbatim from this doc's Appendix.)

**Follow-up: the bio didn't reflect current MindSet Data work.** User caught that the technical bio (STM32/ESP32, UART/SPI) reads as if MindSet Data builds embedded hardware, which it doesn't — that's the pre-MindSet background. Fixed with two changes, both live:
1. Reframed the old bio as clearly past ("Before that, I spent years as..." + past tense throughout) instead of present-tense current work.
2. Added a new paragraph naming the real current technical scope at MindSet Data.

**Final About text, live on the profile:**

> Co-founder of MindSet Data, where we're building the future of industrial data infrastructure — more on that soon.
>
> At MindSet Data, I work on Go-based backend systems, industrial protocol integration (OPC-UA, MQTT), SQL data connectors, and real-time data pipeline architecture — bridging OT and modern software systems.
>
> Before that, I spent years as an Embedded Software Engineer specialized in real-time and low-power systems. I designed robust and scalable IoT solutions using FreeRTOS, Embedded Linux, and bare-metal C/C++ on platforms like STM32, ESP32, and PIC.
>
> My work included:
> Firmware development with secure communication (MQTT, TLS)
> Deep sleep & energy optimization for battery-powered devices
> Multi-protocol integration (UART, SPI, I2C, CAN)
> Cloud integration with HiveMQ, ThingSpeak, and custom dashboards
> Linux kernel customization and device driver development
>
> Passionate about clean architecture, efficient resource use, and delivering reliable embedded products — from concept to deployment.

**Skills section — added, live:** Go (Programming Language), SQL, REST APIs, Data Pipelines. Chosen to match the new middle paragraph — real technical scope from the actual work in this repo (`internal/connections`, `sql_query`, the pipeline engine, OPC-UA/MQTT bridging), not invented.

## 5. Smaller polish items

- ~~**Public profile URL**~~ — **resolved on its own.** Was `linkedin.com/in/mohamed-khenafif-52844335b`; checked after today's other edits and it's now the clean `linkedin.com/in/mohamed-khenafif`, without anyone explicitly setting it. LinkedIn appears to auto-claim the vanity slug once available — no action needed.
- **"Open to"** is currently unset. Not a problem by default, but worth a deliberate choice ("Providing services" is the closest fit if you want something set) rather than leaving it blank, since prospects and investors sometimes check it.
- **Post mix:** your three highest-reach posts (RTOS+MPI — 4.1K impressions, RTOS scheduling — 4.3K, ESP32 deep sleep — 7.5K) are all deep embedded tutorials for an embedded-engineer audience. Your MindSet Data mission post is the most recent but lowest-reach of the five (461 impressions) — expected, different audience, not a problem on its own. **Decision: skip Featured-pinning it** — it carries the full €/silos pitch, which contradicts "light framing, no specifics." Leaving the technical posts as the most visible content is actually consistent with staying quiet on MindSet specifics right now. Revisit at official launch.

---

## 6. Posts / Activity audit — full history, 5 items total

Checked the complete activity list (`/recent-activity/all/`, clicked through "Show more results" — confirmed 5 is everything, not a partial view).

| # | Post | Age | Reach | Verdict |
|---|---|---|---|---|
| 1 | MindSet Data / Cécilia Tran post — OT data-to-decision thesis | 1mo | 461 impressions, **9 comments** (best engagement of the 5) | **Keep.** Already public, best comment engagement of anything on the profile — deleting it would waste real traction. See note below on the tension with "light framing." |
| 2 | Repost — Talel BELHAJ SALEM's TechLeef Yocto/embedded-Linux academy | 6mo | — | **Candidate for removal**, low priority. Not harmful, but pure community-support content, unrelated to the founder narrative — dilutes profile focus without adding anything. |
| 3 | "RTOS + MPI" technical post (Embedded C Programming group) | 1yr | 4.1K impressions | **Keep** — real technical credibility, no conflict with staying quiet on MindSet specifics. |
| 4 | "RTOS Task Scheduling" technical post (same group) | 1yr | 4.3K impressions | **Keep**, same reasoning. |
| 5 | "ESP32 Deep Sleep" technical post (same group) | 1yr | 7.5K impressions | **Keep** — highest-reach post on the profile. |

**Tension flagged, not yet resolved:** post #1 already does exactly what "light framing, no specifics" is meant to avoid for the headline/About — it names Cécilia Tran, uses `#MindSetData`, states the full silos→€ pitch, and invites "let's talk." It predates today's visibility decision. Not recommending deletion (real engagement, already public, retracting it would look worse than leaving it), but flagging that "light framing" going forward is really about *not adding more* specific public language right now, not about the post history being fully scrubbed of it. Worth an explicit call if it matters.

Net: one low-stakes optional removal (item 2), nothing urgent.

---

## 7. Content strategy — groups + "indirect" post ideas (2026-07-18)

User's goal: build audience/followers who'd be interested in MindSet Data's solution or skills, without posting about the solution directly — the same recipe that already works for the embedded posts (#3–5 above: 4.1K–7.5K impressions, zero company pitch), just aimed at the new stack (Go, SQL connectors, OPC-UA/MQTT, data pipelines) instead of embedded hardware.

### Candidate groups — found via live LinkedIn group search, not guessed

**Broad reach (volume):**

| Group | Members | Fit |
|---|---|---|
| Industry 4.0 & the Industrial Internet | 31K | IoT/IIoT leaders, broad OT audience |
| Fourth Industrial Revolution (Industry 4.0) | 15K | Smart factory / digital transformation crowd |
| Mechatronics and Industrial IoT with ML/DL/AI | 19K | Engineers + data scientists — bridges both audiences |
| IIoT (Industrial Internet of Things) and Industrial Ethernet | 16K | Automation/IIoT executives |
| Analytics & AI in Supply Chain and Manufacturing | 8K | Public, data-analytics-in-manufacturing angle |

**Small but exactly on-theme (where the actual future customers/partners are):**

| Group | Members | Fit |
|---|---|---|
| MES & Smart Manufacturing Community | 117 | Explicitly lists "OPC UA," "SQL & Manufacturing Data Analytics" as discussion topics |
| IT-OT Integration & Smart Factory Automation | 49 | The OT/IT convergence problem, as a named group topic |
| vNode Automation DACH — IIoT Edge \| OPC UA \| MQTT \| REST | 175 | Same protocol stack actually used in the product |
| Industrial Automation & Industry 4.0 Network | 2K | Explicit "AI & data-driven manufacturing" focus, moderated (max 1 post/week rule — signals real activity, not a dead group) |

Status: **list only, not joined yet** — needs a go-ahead before joining (account-state change), and the user hasn't confirmed.

### Post drafts, full text, sequenced — groups joined by user 2026-07-18

Sequencing logic: safest/most-technical topics first (pure protocol/engineering credibility, zero pitch risk), the closest-to-the-wedge observational post held for last (Post 7) so it reads as insight from an established voice, not a stranger's opener. Different post per group per week — never the same content blasted to all 9 groups at once. "Industrial Automation & Industry 4.0 Network" has a stated max-1-post/week rule, respected by only posting there twice, weeks 3 and 5.

**Week 1 — Post 1: "OPC-UA or MQTT? Wrong question."**
Target: vNode Automation DACH — IIoT Edge | OPC UA | MQTT | REST (175, exact protocol match) + IIoT and Industrial Ethernet (16K, broad reach)

> 🔌 OPC-UA or MQTT? Wrong question.
>
> I keep seeing this framed as an either/or choice. It isn't.
>
> OPC-UA and MQTT solve different problems, and most real industrial data architectures need both.
>
> 📌 OPC-UA — the machine-to-software layer
> - Rich, self-describing data model (types, units, relationships baked in)
> - Built for structured request/response and subscriptions from PLCs/SCADA to your software
> - Strong native security model (certificates, signing, encryption)
> - Where it struggles: not built for lightweight, high-fanout pub/sub across a WAN
>
> 📌 MQTT — the broadcast layer
> - Minimal overhead, built for pub/sub at scale
> - Perfect for distributing already-contextualized data to many consumers (dashboards, cloud, other services)
> - Where it struggles: no native data model — you get bytes, not meaning
>
> 🔧 The pattern that actually works: read from OPC-UA close to the machine, normalize/contextualize the data, then republish over MQTT for everything downstream. You get OPC-UA's structure at the source and MQTT's simplicity everywhere else.
>
> I've built this exact bridge more than once now, and the failure mode is always the same: teams pick one protocol for the whole stack and then fight it for years.
>
> Where does your architecture draw this line — one protocol, or a bridge like this?
>
> #OPCUA #MQTT #IndustrialIoT #Industry40 #IIoT #DataArchitecture

**Week 1 — Post 2: "ISA-95 sounds academic until you actually try to normalize plant data."**
Target: MES & Smart Manufacturing Community (117, exact topic match — group explicitly lists OPC UA/tag normalization) + Mechatronics and Industrial IoT with ML/DL/AI (19K)

> 🏭 ISA-95 sounds academic until you actually try to normalize plant data.
>
> Every plant I've looked at has tags named like this:
> `machine2.ligne1.presion`
> `M1_TEMP_01`
> `Zone3-Press-A`
>
> Three different naming conventions, same underlying concept (a pressure reading on a piece of equipment), zero consistency.
>
> 📌 ISA-95 gives you the hierarchy: Site → Area → Work Center → Work Unit → Tag. The theory is simple. The practice is where it gets interesting.
>
> 🔧 What actually works to normalize raw tags at scale:
> - A dot-count or delimiter heuristic to infer hierarchy depth from naming convention
> - A tag-name dictionary that catches abbreviations *and* language variants — a plant's tags are rarely English-only (presion/pression/pressure all mean the same thing)
> - Treating "can't confidently map this" as a valid outcome, not a bug — better to flag an unmapped tag than silently guess wrong
> - Unit inference as a separate pass from name normalization — a "temperature" tag doesn't imply Celsius
>
> 📌 The payoff isn't cosmetic. Once tags are contextualized, everything downstream — dashboards, rules engines, historians — stops needing to know the raw tag naming convention at all. New machine, new naming scheme, zero downstream changes.
>
> Curious how others have handled the last mile here — do you normalize at ingestion, or push the mess downstream and deal with it in the dashboard layer?
>
> #ISA95 #IndustrialIoT #MES #SmartManufacturing #DataEngineering #Industry40

**Week 2 — Post 3: "Giving software read access to a customer's ERP is scarier than it sounds."**
Target: IT-OT Integration & Smart Factory Automation (49, exact OT/IT + data-access match) + Analytics & AI in Supply Chain and Manufacturing (8K)

> 🔒 Giving software read access to a customer's ERP is scarier than it sounds.
>
> A few guardrails I've learned are non-negotiable when building a SQL connector that reads from someone else's production database:
>
> 📌 Read-only isn't a suggestion, it's enforced twice
> - At the account level (a DB user with SELECT only — no INSERT/UPDATE/DELETE grants)
> - At the query level (reject anything that isn't a SELECT before you even open a connection)
>
> 🔧 The guardrails that actually matter:
> ✅ Parameterized queries only — named placeholders resolved server-side, never string concatenation
> ✅ Mandatory query timeout — a slow query on someone else's production DB isn't your problem until it's everyone's problem
> ✅ Mandatory row limit — cap it, don't trust the caller to add LIMIT themselves
> ✅ Reject multi-statement queries — one SELECT in, one result set out, nothing else gets to ride along
> ✅ A startup health check that verifies the account really is read-only — don't just document it, prove it
>
> ⚠️ The subtle one: what happens when the account CAN write, because someone on the customer side over-provisioned it? Don't refuse to connect — log a loud warning and keep working. Being pedantic about someone else's IT mistake helps nobody.
>
> None of this is exotic. It's just the boring, unglamorous 20% of the connector that determines whether it's safe to point at a real customer's database.
>
> What's the guardrail you've seen skipped that came back to bite someone?
>
> #DataEngineering #SQL #SoftwareArchitecture #IndustrialIoT #Security

**Week 3 — Post 4: "If your SQL-injection test doesn't touch a real database, it isn't testing anything."**
Target: Industrial Automation & Industry 4.0 Network (2K — 1st post here, respects the 1/week rule) + Analytics & AI in Supply Chain and Manufacturing (8K, 2nd appearance, different angle than Post 3)

> 🧪 If your SQL-injection test doesn't touch a real database, it isn't testing anything.
>
> Mocking the database for a connector's unit tests is fine — until you get to the tests that actually matter: does the read-only enforcement work, does the injection guard hold, does the timeout actually fire.
>
> 📌 Those need a real database, or they're theater.
>
> 🔧 What that looks like in practice:
> - Spin up a disposable container (MySQL, Postgres, whatever your target) for the test run, seeded with a tiny known schema
> - Create a second, genuinely read-only DB user in the same container, alongside the writer — you need both to prove the health check actually discriminates between them, not just returns "true" unconditionally
> - Run the real injection attempt: bind a string like `1; DROP TABLE x` as a *parameter*, then actually re-query the table afterward to confirm it's still there — don't just assert "no error," prove nothing broke
> - Test the timeout with a query that's guaranteed slow (`SELECT SLEEP(5)` against a 1-second timeout) and assert it actually cuts off in time, not just that it eventually errors
>
> 📌 The one that catches people out: your "read-only" test needs a matching writer-account test in the same run. If your read-only assertion would also pass for a writable account, it isn't testing what you think it's testing.
>
> Slower to run than a mocked suite. Also the only version of that test suite I actually trust.
>
> #Testing #SQL #SoftwareEngineering #DataEngineering #CI

**Week 4 — Post 5: "Detecting when a machine stops sounds trivial. It isn't."**
Target: Industry 4.0 & the Industrial Internet (31K, broad) + Fourth Industrial Revolution / Industry 4.0 (15K, broad)

> ⚙️ Detecting when a machine stops sounds trivial. It isn't.
>
> On paper: machine running, value drops, machine stopped. Done.
>
> In practice, that naive version fires false positives constantly.
>
> 📌 What actually happens on a real signal:
> - Sensor noise causes momentary drops that aren't real stops
> - A machine can be "running" at zero output during a legitimate changeover
> - The signal that indicates state isn't always a clean boolean — sometimes it's a speed, a current draw, a status word with meaning buried in specific bits
>
> 🔧 What holds up better:
> - A debounce window — a transition only counts once it's held for N seconds, not on the first sample
> - Explicit state, not inferred state — track Running/Stopped/Unknown as its own value, don't just react to raw signal deltas
> - Duration only calculated on a *confirmed* transition, so a flicker doesn't generate a phantom "12-second stop" in your history
> - Treating the ambiguous case (signal missing, sensor offline) as its own state, not silently defaulting to either Running or Stopped
>
> 📌 The payoff: once state detection is solid, everything built on top of it — micro-stop detection, cost calculations, dashboards — inherits that reliability for free. Get this layer wrong and every report above it is quietly wrong too.
>
> What's the noisiest signal you've had to build state detection around?
>
> #IndustrialIoT #Automation #StateMachines #Manufacturing #Industry40

**Week 4 — Post 6: "Same concept, six different database schemas. Every single time."**
Target: MES & Smart Manufacturing Community (117, 2nd appearance, deepens niche presence) + vNode Automation DACH (175, 2nd appearance)

> 🗂️ Same concept, six different database schemas. Every single time.
>
> Ask five different ERPs for "the current work order on this line" and you'll get five different table names, five different column names, and at least two different status-code conventions.
>
> 📌 SAP calls it AUFNR. Odoo calls it name. A custom Access-derived system calls it NumOF. Same concept, zero shared vocabulary.
>
> If you let that heterogeneity leak into every downstream system — dashboards, rules engines, cost calculations — every one of them has to become schema-aware. That doesn't scale past your second customer.
>
> 🔧 What scales instead:
> - A small, opinionated canonical model — a fixed set of fields every downstream feature can assume exist (of_number, product_code, status, and so on), independent of any one customer's schema
> - A translation layer at the connector, not downstream — map raw columns to canonical fields once, at the edge, so nothing past that point needs to know the source system
> - Enum translation as its own concern — a status field isn't fully normalized just because the column is renamed; `CRTD/REL/TECO` and `draft/confirmed/done` need mapping to the same canonical values too
> - Treating an incomplete mapping as a soft failure — a field with no mapping should degrade gracefully, not break the pipeline
>
> 📌 The real payoff shows up on the second customer. If onboarding them costs the same as the first, the abstraction is doing its job.
>
> Anyone else building against a zoo of ERPs with zero shared schema? How do you keep the translation layer from becoming its own maintenance burden?
>
> #DataEngineering #ERP #SystemIntegration #Manufacturing #SoftwareArchitecture

**Week 5 — Post 7: "The most expensive minute in a factory is usually the one nobody measured."**
Target: Industrial Automation & Industry 4.0 Network (2K, 2nd appearance, respects the 1/week rule — 2 weeks after Post 4) + Fourth Industrial Revolution (15K, 2nd appearance)
Deliberately closest to the wedge, held for last — no company name, no CTA, unlike the existing MindSet post.

> 📉 The most expensive minute in a factory is usually the one nobody measured.
>
> A pattern I keep running into, across very different plants and industries:
>
> A line stops for 45 seconds. Nobody notices in real time — it's below the threshold anyone bothers to log manually. It happens six more times that shift. By the time anyone looks at a weekly report, the actual cost is invisible, folded into a vague "downtime" number nobody can act on.
>
> 📌 Not because the data doesn't exist. The PLC saw every one of those stops. It's because the path from "a machine changed state" to "someone with budget authority understands what that cost" has never been built to survive contact with a real factory's mess of legacy equipment, inconsistent naming, and systems that don't talk to each other.
>
> 🔧 The gap isn't collection — most plants already have more raw signal than anyone's using. It's context. A stop means nothing on its own; a stop tied to the work order that was running, the product's margin, and what was scheduled next means something a plant manager can act on in the moment, not three weeks later in a spreadsheet.
>
> I don't think this is a data problem anymore. I think it's an infrastructure problem — one layer that was never really designed, sitting between raw signals and the decisions people actually need to make.
>
> Curious whether others see this as an integration gap, a tooling gap, or something more structural.
>
> #Manufacturing #OT #Industry40 #DataOps #Downtime

Status: **7 posts drafted, full text ready. None posted yet.**

---

## 8. Messages inbox review + co-design outreach candidates (2026-07-19)

Reviewed all 17 threads in the LinkedIn inbox to find existing contacts worth a co-design ask.

**Investor threads (already covered in the earlier status check):** Maxime Lhoustau (Motier Ventures), Antoine Loiseau (ex-VC/angel), Pierre Ben Kiran — all routed to `cecilia@mindset-data.com`, awaiting her follow-up. Not co-design targets — these are fundraising conversations.

**Everything else in the inbox — checked individually, none are co-design targets:** Benoit Camus (vendor pitch, SaaS/CRM building), Tijn van Daelen (vendor pitch, AI dev tooling), Arslan Akram (vendor pitch, QA outsourcing, already closed — "this will be my last message"), Sérène Dupré (recruiter), Théo Louro (SaaS content creator, unrelated), Zeshan Ali (vendor pitch, Zero-Axis Technology dev agency), Julien Lijeour (unclear service provider), G.Rodney Rendambo (recruiter, offered Mohamed a job, declined), Cécilia (co-founder, informal), Ramzi Belkhelfa (personal, unrelated to business).

**Finding: zero existing manufacturing/industrial-operations prospects in the inbox.** The whole network skews toward people pitching *Mohamed* (dev agencies, recruiters), investors, and personal contacts — not plant managers or ops leads. A co-design ask can't come from mining the inbox; it needs fresh targets.

### Candidates identified — via LinkedIn people search among 1st-degree connections

| Name | Profile | Relationship | Why |
|---|---|---|---|
| **Rami BOUMEKHITA** | Automation & Control engineer (École Polytechnique), Master's in Mobile Automatic Systems | Existing personal contact (casual rapport, Arabic) | User's pick — an engineer who could open a door to his manager |
| **Boudjemaa Abdelhadi TELLI** | Ingénieur Automatisme & Informatique Industrielle \| Industrie 4.0 \| Systèmes Intelligents & MES/SCADA | 1st-degree, 7 mutual connections, no prior conversation | User's pick — near-exact profile match to MindSet's domain |
| **Randy LENDOYE** | Industrial IT \| PLC \| SCADA \| MES \| ERP \| EMS \| BAS \| Industrial Innovation \| Industry 4.0 | 1st-degree, found via search, no prior conversation | Proposed — literally spans the OT/IT bridge MindSet builds |
| **Doria Belahbib** | Industrial engineer and information system project director | 1st-degree, found via search, no prior conversation | Proposed — director-level, closer to being the decision-maker herself than someone who needs to escalate |

### Draft outreach messages — not yet sent

**Rami** (matches his existing casual tone):
> Salam Rami,
> Content d'avoir de tes nouvelles ! Au fait, je suis en train de monter MindSet Data avec ma co-fondatrice — on bosse sur l'infrastructure data pour l'industrie, pile dans ton domaine (SCADA/automatisme). Ça te dirait qu'on en discute vite fait un jour ? Et si jamais ça résonne avec ce que vous faites chez vous, je serais curieux de savoir si ça vaudrait le coup d'en toucher un mot à ton responsable. Zéro pression, juste curieux d'avoir ton avis d'ingénieur terrain.
> À bientôt,
> Mohamed

**Boudjemaa** (first contact, references mutual network):
> Bonjour Boudjemaa,
> On a plusieurs connexions en commun donc je me permets de te contacter directement. Je co-fonde MindSet Data — on travaille sur l'infrastructure data industrielle (bridging OT/IT, SCADA/MES), exactement ton domaine vu ton profil. Je serais curieux d'échanger 15-20 min pour avoir ton regard terrain sur les problématiques que vous rencontrez côté MES/SCADA. Et si ça résonne avec ce que vit ton équipe, ce serait avec plaisir d'être mis en contact avec la bonne personne côté décision pour explorer un co-design.
> Bien à toi,
> Mohamed

**Randy and Doria:** not yet drafted — offered, not requested yet.

Status: **candidates identified, 2 messages drafted, none sent.**

---

## 9. The real outreach tracker — `Outbound` Google Sheet (shared by the user, 2026-07-19)

User shared the actual shared outreach spreadsheet: https://docs.google.com/spreadsheets/d/1RhLRMJgKAAAVOOMkLB9bkf7Q2juuWbJ2GYkVH3DeJU4/ — this supersedes the ad-hoc candidate-hunting in §8; it's the real, existing source of truth for who's already being tracked.

**5 tabs, reviewed in full:**

| Tab | Contents |
|---|---|
| **Outreach List** | ~45 named contacts, mostly "Responsable Automatisme/Informatique Industrielle/MES/IT-OT/Amélioration Continue" and "Directeur Technique" roles at real mid-sized industrial companies matching `Blurb_Invest.md`'s target verticals: pharma (Curium Pharma, Ethypharm, Delpharm ×4 sites, Lesieur Cristal), food/dairy (Candia, LAITA, Fleury Michon, Groupe Aoste, Darégal, Savencia), chemicals (Alsachimie, WeylChem). Has a "Lemlist" sub-section (cold-email tool) with full contact details (email, LinkedIn URL, company, role). No status/reply-tracking column on this tab. |
| **Prospecting factories MVP** | A smaller, separate list — smaller/family-owned food companies (Charcuterie Vendéenne, Bonilait, Ingredia, Brasserie Caulier, Guyader Gastronomie) — this one **does** have a live status column ("Pending," "2ème relance," "2ème relance envoyé"), plus a "NETWORK" section of personal warm contacts with intro counts. Looks like the active pilot/MVP-customer track. |
| **SI** | Nearly empty — one note: "Check Siemens/Schneider partners" (future systems-integrator channel idea, ties to the OT Integrators GTM pillar in `Blurb_Invest.md`). |
| **Digital solutions/startups** | Barely started — 2 rows (FocusProd, Doqcheck). |
| **Online Forum / Groups** | Completely empty — this is where §7's 9-groups + 7-posts content plan belongs; nobody has populated it yet. |

**Correction to §8:** confirmed Rami BOUMEKHITA (Verkor) and Randy LENDOYE (Curium Pharma) are already on the "Outreach List" tab — my LinkedIn search yesterday found real, correctly-targeted people, they just weren't new discoveries. Boudjemaa Abdelhadi TELLI and Doria Belahbib are genuinely not on the sheet — still net-new.

Status: **reviewed, not yet acted on** — offered to populate "Online Forum / Groups," add Boudjemaa/Doria as new rows, or find more candidates. User chose the third.

## 10. Second candidate batch — new LinkedIn search, cross-checked against the sheet (2026-07-19)

User asked specifically for candidates *not* already on the "Outreach List." Ran 2 more targeted searches ("Responsable Informatique Industrielle," "responsable automatisme cosmétique" — the latter to hit the cosmetics vertical, which `Blurb_Invest.md` names but the sheet barely covers) and cross-checked every company name against the ~45 already tracked. All 11 below are confirmed not on the sheet.

| Name | Title | Company | Vertical | Notes |
|---|---|---|---|---|
| Yannick Martineau | Directeur Pôle Automation / Directeur SI / **CODIR** | — | Industrial | Executive-committee level — the decision-maker himself, no escalation needed |
| Guillaume Merlier | Responsable informatique industrielle | UGITECH | Metallurgy | Named target vertical, underrepresented on the sheet; mutual connection (Sarah EL AYNAOUI) |
| Etienne FANTONE | Directeur — Responsable bureau d'étude Automatisme/Informatique Industrielle | — | Industrial | Director level; mutual connection (Sandrine Rondeau) |
| Guillaume Labadie | Responsable Informatique Industrielle (Process Control) | Chevron Oronite | Petrochemicals | Le Havre |
| Romain Doaré | Responsable du pôle informatique industrielle (DSI) | GRDF | Utilities | Mutual connection (Maxime De Oliveira) |
| Gwenaël LOYZANCE | Responsable projet Automatisme et Informatique Industrielle | OET | Industrial | Rennes |
| Samuel Benavente | Automation and data processing manager | ex-Boccard (agro/beverage/cosmetics integrator) | Cosmetics-adjacent | Mutual connection (Prince Noukounwoui) |
| Clément BRUNEL | Responsable Technique Automatisme et Informatique Industrielle | — | Industrial | Lyon |
| Emmanuel LEBRETON | Responsable Automatisme et Informatique Industrielle | — | Industrial | Pays de la Loire |
| Louis Plassais | Responsable projets Informatique Industrielle | Ose Group | Industrial | Angers |
| Léandre DE ROECK | Responsable Projets | PACOSPHARM | Pharma | Mutual connection (MAZELIN Hervé) |

**Suggested priority:** Yannick Martineau + Etienne FANTONE first (director/exec level, straight to the decision-maker), then Guillaume Merlier (fills the metallurgy gap), then the rest.

Status: **11 candidates identified, none added to the sheet.** Romain Doaré accepted Mohamed's connection request (2026-07-19) — first outreach message drafted below, not yet sent.

### Draft outreach message — Romain Doaré (2026-07-19, not yet sent)

First contact, post-connection-acceptance. He's DSI-level (head of the industrial-IT department at GRDF, a gas utility) — the decision-maker himself, so the message asks directly rather than requesting an escalation (unlike the Boudjemaa template in §8). Vouvoiement used given seniority/corporate register — flag if `tu` reads better for this network.

> Bonjour Romain,
>
> Merci d'avoir accepté ma demande de connexion ! On a Maxime De Oliveira en commun, donc je me permets de vous contacter directement.
>
> Je co-fonde MindSet Data — on construit l'infrastructure data pour les environnements industriels, avec un focus sur le pont entre OT et IT (SCADA, informatique industrielle). Vu votre rôle chez GRDF, je serais curieux d'échanger 15-20 minutes sur les problématiques que vous rencontrez côté remontée et contextualisation de la donnée terrain.
>
> Si ça résonne avec vos enjeux actuels, ce serait avec plaisir d'explorer ensemble un potentiel co-design.
>
> Bien à vous,
> Mohamed

---

## 11. Investor lead — Polytechnique Ventures / Denis Lucquin Catalyst Initiative (2026-07-19)

Not a profile/content item — a live investor thread the user surfaced via 2 WhatsApp screenshots (`docs/WhatsApp Image 2026-07-19 at 1.19.57 AM.jpeg` and `...(1).jpeg`) plus a LinkedIn post. Kept in this doc rather than a new file since it came out of the same messaging-review thread — split it out later if it grows into its own workstream.

**The screenshots:** Cécilia messaged Quentin Sanchez (met at Vivatech, discussed the X-UP program — École Polytechnique's incubator) asking for an intro to Polytechnique Ventures after seeing their pre-seed announcement. Quentin's reply: *"Étrange! Si ton associé est polytechnicien, il peut contacter directement Gaspard Devissaguet: gaspard.devissaguet@polytechnique-ventures.fr. Il se fera une joie de vous accorder un peu de temps ;)"* — recommending Mohamed reach out directly rather than go through an intro, specifically because he's the Polytechnicien in the founding pair. Cécilia: "Super! Merci beaucoup je vais relayer le message :)".

**The LinkedIn post** (`linkedin.com/feed/update/urn:li:activity:7480293384188837888/`): Polytechnique Ventures announced the **Denis Lucquin Catalyst Initiative** — 5% of their fund now goes to pre-seed startups, average ticket **€150k**, 1–2 calls/year for the most ambitious young ventures. Eligibility: deeptech, and at least one founder graduated from, is a researcher at, was incubated by, or emerges from École Polytechnique's labs. Gaspard Devissaguet is tagged directly in the post.

**Eligibility caveat, flagged to the user:** the LinkedIn profile lists "Institut Polytechnique de Paris — Master 2 ROSP," the broader graduate consortium (includes École Polytechnique, Télécom Paris, ENSTA, etc.), not necessarily the same as an "X" (École Polytechnique) engineering degree the post specifically names. Recommended being precise about this in the email rather than letting Gaspard assume and find out later.

**Draft v1** (superseded — kept for reference): led with a generic OT/IT-silos description instead of a concrete hook, and included a self-doubting line about the Institut Polytechnique de Paris vs. École Polytechnique distinction that would have undercut the ask before Gaspard even engaged.

> Objet : Polytechnicien × MindSet Data — pré-amorçage / Denis Lucquin Catalyst Initiative
>
> Bonjour Gaspard,
>
> Quentin Sanchez m'a suggéré de vous contacter directement. Je suis diplômé de l'Institut Polytechnique de Paris (Master 2 ROSP) et co-fondateur de MindSet Data, avec Cécilia Tran.
>
> J'ai vu votre post sur la Denis Lucquin Catalyst Initiative et le nouveau focus pré-amorçage — ça correspond bien à où on en est. MindSet Data est une infrastructure data industrielle : on connecte les silos OT/IT en usine (PLC/SCADA, ERP/MES), on traduit chaque événement en son enjeu économique, et on pousse la bonne action à la bonne personne. On est en phase d'exploration avec des investisseurs européens pour un tour de table solide.
>
> Auriez-vous 20-30 minutes pour en discuter ?
>
> Bien à vous,
> Mohamed Khenafif

**Draft v2 — "excellent," rewritten** — superseded per user feedback below, kept for reference: leads with the finance/shopfloor wedge (the concrete 45-second-stop pattern, same wedge decided in Entry 74) instead of an abstract infra description.

> Objet : Mohamed Khenafif (Institut Polytechnique de Paris) — MindSet Data, via Quentin Sanchez
>
> Bonjour Gaspard,
>
> Quentin Sanchez m'a recommandé de vous écrire directement, après avoir vu votre post sur la Denis Lucquin Catalyst Initiative.
>
> Je suis diplômé de l'Institut Polytechnique de Paris (Master 2 ROSP) et co-fondateur de MindSet Data, avec Cécilia Tran (EDHEC).
>
> En résumé : dans une usine, un arrêt de 45 secondes ne coûte rien sur le papier — la donnée existe, le PLC l'a vue. Le problème, c'est que personne ne relie cet événement à ce qu'il coûte vraiment (l'OF en cours, la marge produit, le planning) avant qu'il soit trop tard pour agir. MindSet Data connecte les silos OT/IT en usine et traduit chaque événement en son enjeu économique, en temps réel — pas un énième tableau de bord, la couche qui manque entre le signal et la décision.
>
> On explore actuellement un tour de pré-amorçage avec des investisseurs européens, en parallèle d'un pipeline de partenariats co-design avec des groupes industriels de taille intermédiaire (pharma, agroalimentaire, chimie, métallurgie).
>
> Auriez-vous 20–30 minutes dans les prochaines semaines ? Je peux aussi vous envoyer un executive summary d'une page avant, si c'est plus simple pour vous.
>
> Bien à vous,
> Mohamed Khenafif
> mohamed@mindset-data.com

**Draft v3 — user feedback: no operational use-case (micro-stop) in an investor email; lead with finance/AI/silos/decision instead.** Swapped the middle paragraph for the macro thesis — pulls from `Blurb_Invest.md`'s "Why now" (AI agents hitting a wall on OT context, margin pressure) and "Why existing solutions fall short" (silos, decision latency, economic prioritization) rather than a shopfloor anecdote. Matches VC-audience instinct better than customer-audience instinct: investors want the market-timing thesis, not a single operational vignette. Not yet sent, not yet confirmed.

> Objet : Mohamed Khenafif (Institut Polytechnique de Paris) — MindSet Data, via Quentin Sanchez
>
> Bonjour Gaspard,
>
> Quentin Sanchez m'a recommandé de vous écrire directement, après avoir vu votre post sur la Denis Lucquin Catalyst Initiative.
>
> Je suis diplômé de l'Institut Polytechnique de Paris (Master 2 ROSP) et co-fondateur de MindSet Data, avec Cécilia Tran (EDHEC).
>
> En résumé : les industriels s'appuient sur des systèmes cloisonnés — OT, ERP, MES — qui ne communiquent pas entre eux, et perdent un temps considérable à reconstituer manuellement une vision d'ensemble. Résultat : des décisions prises sur des données d'hier, au moment précis où toutes les stratégies IA butent sur ce même manque de contexte temps réel côté OT. MindSet Data connecte ces silos et traduit chaque événement en son enjeu économique, pour que les équipes — et demain les agents IA — sachent quoi prioriser, en temps réel.
>
> On explore actuellement un tour de pré-amorçage avec des investisseurs européens, en parallèle d'un pipeline de partenariats co-design avec des groupes industriels de taille intermédiaire (pharma, agroalimentaire, chimie, métallurgie).
>
> Auriez-vous 20–30 minutes dans les prochaines semaines ? Je peux aussi vous envoyer un executive summary d'une page avant, si c'est plus simple pour vous.
>
> Bien à vous,
> Mohamed Khenafif
> mohamed@mindset-data.com

Status: **v3 is current, not sent, not yet confirmed by user.**

---

## 12. Add a 2021 startup project — "CLock" — to Experience (2026-07-19)

User provided `docs/CLock.pptx` (10-slide pitch deck, image-only slides — read via unzip + image extraction) and asked to add this as a LinkedIn Experience entry. Not yet applied — draft below, pending confirmation.

**What the deck says:**

- **Slide 1 (title):** "CLock" — tagline "Protecting you is our mission." Two named authors: Hechaichi Abdelbasset and Mohamed Khenafif.
- **Slide 2 (The Problem):** Cites 2015 Algeria stats — 17,318 night robberies (houses/apartments/factories), 5,352 car thefts. Framed need: "Properties and assets Protection. Surveillance: Roads, indoor applications and etc."
- **Slide 3 (Conventional security systems):** Existing options (ultrasonic sensors, dome cameras, alarm keypads, sirens) called out as: expensive (multiple sensors needed), not reliable (false alarms), not efficient (cameras only record, don't act).
- **Slide 4 (Our Solution):** "An intelligent security system" that detects suspicious actions (burglaries, arson, accidents) and reacts automatically — triggers alarm, calls police and owner, closes entries where applicable. Positioning line: "Higher reliability for a reasonable price."
- **Slide 5 (How it works):** Two-part system — **Hardware** (cameras, microcontroller, siren alarm, auto-dial module) + **Algorithms** (facial recognition, behavior detection).
- **Slide 6 (Our Team):** "The National Polytechnic School Algiers, 5th Year" — three named members:
  - Hechaichi Abdelbasset — 9th place Algeria Innovation Prize, INJAZ EL DJAZAIR participant, Sponsor Coordinator Hult Prize ENSV, B2C AIESEC — Background: Power Engineering. *(Note: this achievement line is listed under his personal bio, not explicitly stated as CLock's own competition result — don't present "9th place" as the project's placement unless you can confirm it was for CLock specifically.)*
  - Mohand Boudraa Amokrane — INJAZ EL DJAZAIR participant, Drillbotics competition finalist (Germany), Skills: AI — Background: Power Engineering.
  - **Mohamed Khenafif** (you) — Skills: Electronics, AI, Python, C++ — Background: Electronics Engineering.
- **Slide 7 (Competition):** 2×2 "Expensive/Affordable × Efficient/Not Efficient" positioning map. Competitors placed as expensive (Eyepix Solutions, APE Systems, SES – Sécurité Electronic Système) or affordable-but-not-efficient (SOSETEL). CLock placed alone in the "Efficient + Affordable" quadrant.
- **Slide 8 (Business Model):** Manufacturing → Product, sold via delivery-company partnerships to **Real Estate Owners** and **Shop Owners**, and via dedicated personal assistance to **Government** and **Institutions**.
- **Slide 9 (Financials):** 5-year projection table (2022–2026), revenue growing 48.5M → 382.6M (currency not labeled, presumably DZD), net profit swinging from -3.3M in 2022 to +13.3M by 2026 — standard illustrative business-plan projections, not actuals.
- **Slide 10 (closing):** "Thank You For Your Attention" — contact `abdelbasset.hechaichi@g.enp.edu.dz`.

**Reading of context:** this is a student entrepreneurship/innovation project from École Nationale Polytechnique (ENP Algiers), built with two classmates during your final ("5th") year — consistent with the ENP dates already on your profile (2016–2021) and just before/overlapping the Sonatrach internship (Feb–Aug 2021). No slide states exact project dates or a confirmed competition result for CLock as a team — treat "2021" and any competition placement as needing your confirmation, not asserted fact.

**Draft Experience entry (not yet applied):**

> **Title:** Co-Founder
> **Company:** CLock
> **Employment type:** Self-employed (or leave blank — LinkedIn doesn't have a great fit for "student venture")
> **Dates:** Oct 2020 – Dec 2021
> **Location:** Algiers, Algeria
>
> **Description:**
> Co-founded CLock, an intelligent security system, at École Nationale Polytechnique (Algiers). Conventional security systems (cameras, alarm keypads, sensors) are expensive and reactive — they record or alert, but don't act. CLock combined embedded hardware (cameras, microcontroller, siren, auto-dial module) with AI (facial recognition, behavior detection) to detect violence, fire, theft, and intrusion **before they occur**, by analyzing suspicious behavior in real time — then respond automatically: trigger the alarm, call the owner and police, and close entries where applicable.
>
> My role: designed and built the AI detection layer — training and integrating the facial-recognition and behavior-analysis models that flag suspicious activity in real time, and writing the Python/C++ logic that turns a detection into an automated response (alarm, emergency calls, entry lockdown).
>
> Built for everywhere protection matters — homes, shops, and government/institutional buildings — reached through delivery-company partnerships and dedicated assistance for larger accounts.

**Status: applied and live** (2026-07-19). Added via Add section → Experience on the live profile: Title "Co-Founder", Company "CLock" (free text — no matching LinkedIn company page existed), Employment type "Self-employed", Oct 2020 – Dec 2021, Location "Algiers, Algeria", description as drafted above. Verified via `/details/experience/` — now the 6th entry on the profile, below ISPESAGE. Confirmed changes made during drafting: (a) no naming/tagging of co-founders; (b) title is plain "Co-Founder", no "Electronics"; (c) detected situations = "violence, fire, theft, and intrusion"; (d) no "student venture" framing; (e) predictive framing ("before they occur", real-time behavior analysis); (f) dates: Oct 2020 – Dec 2021; (g) role line expanded (training/integrating detection models + response logic); (h) target line broadened from a narrow B2B buyer-segment list to "everywhere protection matters" — houses/shops/government, not framed as a sales-buyer taxonomy. Note: right after saving, the main profile page (`/in/mohamed-khenafif/`) briefly didn't render the Experience section at all (stale client-side cache) — the `/details/experience/` view confirmed the save was correct; a manual refresh of the main profile should show it too.

**Follow-up (2026-07-19): skills tagged on the entry.** User flagged the entry was saved without skills (every other Experience entry has 4-5 tagged skills; this one had none). Added, consistent with "forget the electronic" — AI/software skills only, nothing hardware-framed: **Artificial Intelligence (AI), Python (Programming Language), C++ (Programming Language), Computer Vision, Machine Learning.**

## Priority order / status

1. ~~Write the MindSet Data experience description~~ — **declined**, revisit at official launch
2. Remove/merge the "Stealth Startup" duplicate — **agreed, not yet executed** — the one item still outstanding
3. ~~Update the headline~~ — **done, live**
4. ~~Rework About~~ — **done, live**, including the follow-up fix so it reflects current MindSet Data work, not just pre-MindSet background (recovered after an accidental mid-edit deletion along the way — see Appendix)
5. ~~URL cleanup~~ — **resolved on its own**, now `linkedin.com/in/mohamed-khenafif`
6. ~~Add Skills (Go, SQL, REST APIs, Data Pipelines)~~ — **done, live**
7. Featured section — **skipped**, revisit at official launch
8. Talel repost removal — **optional**, low priority, user hasn't weighed in yet
9. ~~Join the 9 candidate groups (§7)~~ — **done, user joined 2026-07-18**
10. Draft/post the first "indirect" content post (§7) — **ideas only, not started**
11. Send co-design outreach to Rami + Boudjemaa (§8) — **drafted, not yet sent**
12. Draft + send outreach to Randy LENDOYE + Doria Belahbib (§8) — **candidates identified, messages not drafted**
13. Add Boudjemaa + Doria to the real "Outreach List" sheet (§9) — **not done, user chose to keep prospecting instead**
14. Populate the empty "Online Forum / Groups" tab with §7's plan (§9) — **not done, offered, not requested**
15. Add/draft outreach for the 11 new candidates (§10) — **identified, nothing else done yet**
16. Send the Gaspard Devissaguet / Polytechnique Ventures email (§11) — **drafted, not yet sent**
17. Add the CLock Experience entry (§12) — **drafted, not yet applied** — pending confirmation on dates + co-founder tagging

---

## Appendix — full profile snapshot, current as of 2026-07-18 (post-edits, raw, for reference)

Captured so this doc is self-contained and doesn't require re-pulling the live profile to reconstruct context. This is the **live, current** state (after today's edits) — the pre-edit originals are preserved above in items 3–5 for reference (what changed and why).

**Header:** Mohamed KHENAFIF · Paris, Île-de-France, France · 185 connections · 192 followers
**Public profile URL:** `linkedin.com/in/mohamed-khenafif` (auto-resolved from the numeric-suffix slug — see item 5)
**Profile language:** English
**"Open to":** unset — nothing currently selected
**Analytics (private to the user):** 223 profile views · 13 search appearances

**Headline (live):**
> Co-Founder @ MindSet Data | building the future of industrial data | Embedded Software Engineer | FreeRTOS, STM32, ESP32

**About (live, full text — see item 4 for the edit history):**
> Co-founder of MindSet Data, where we're building the future of industrial data infrastructure — more on that soon.
>
> At MindSet Data, I work on Go-based backend systems, industrial protocol integration (OPC-UA, MQTT), SQL data connectors, and real-time data pipeline architecture — bridging OT and modern software systems.
>
> Before that, I spent years as an Embedded Software Engineer specialized in real-time and low-power systems. I designed robust and scalable IoT solutions using FreeRTOS, Embedded Linux, and bare-metal C/C++ on platforms like STM32, ESP32, and PIC.
>
> My work included:
> Firmware development with secure communication (MQTT, TLS)
> Deep sleep & energy optimization for battery-powered devices
> Multi-protocol integration (UART, SPI, I2C, CAN)
> Cloud integration with HiveMQ, ThingSpeak, and custom dashboards
> Linux kernel customization and device driver development
>
> Passionate about clean architecture, efficient resource use, and delivering reliable embedded products — from concept to deployment.

**Top Skills (live):** Go (Programming Language) · SQL · REST APIs · Data Pipelines
**Full Skills section:** 30 skills total (not itemized here — mostly embedded/protocol skills predating today's edits, e.g. Communication Protocols, Low-level Drivers).

**Experience — all 5 entries** (a 5th, the Sonatrach internship, was missed in the first pass of this doc — found when re-checking the full live page):

1. **Co-founder** — MindSet Data — Jun 2026 – Present · 2 mos · France — *(no description — item 1, declined)*
2. **Co-Founder** — Stealth Startup · Permanent — Feb 2026 – Present · 6 mos · France — "Powering the infrastructure behind manufacturing. Stay tuned." *(item 2 — still pending removal)*
3. **Embedded Software Developer** — Freelance (Self employed) · Freelance — Dec 2024 – May 2026 · 1 yr 6 mos — "Developing embedded applications and IoT systems using C/C++ and ESP32," FreeRTOS/MQTT/ESP-IDF/OTA/CI details · Skills tagged: Embedded Devices, Embedded Software, +4 more
4. **Embedded Software Engineer** — ISPESAGE · Contract — May 2022 – May 2023 · 1 yr 1 mo — embedded C/C++ firmware for STM32/PIC industrial weighing systems, UART/I2C/SPI, LCD UI · Skills tagged: Internet Protocol Suite (TCP/IP), C (Programming Language), +3 more
5. **Embedded Software Developer – Intern** — Sonatrach · Internship — Feb 2021 – Aug 2021 · 7 mos — real-time autonomous embedded system (temperature/pressure/flow sensors + closed-loop actuators), STM32F103/PIC18F45K22, <50ms response time validated on an industrial test bench · Skills tagged: Real-Time Operating Systems (RTOS), C (Programming Language), +4 more

**Education (not previously captured in this doc):**
- Institut Polytechnique de Paris — Master 2 ROSP — 2023–2024
- Ecole Nationale Polytechnique — Electronics engineer, Electrical and Electronics Engineering — 2016–2021 · Skills tagged: LTSpice, I2C, +3 more

**Connected apps:** Gamma, IntelliJ IDEA, HubSpot, Replit
**Companies followed (Interests):** Stripe, Polytechnique Ventures, and others (not fully itemized)

**Post #1 — the MindSet Data / Cécilia Tran post, full text** (1mo old, 480 impressions, 9 comments, 1 repost, visible to anyone on or off LinkedIn):

> I've spent the last years deep in embedded systems, across industries, across protocols, across constraints. 🛠️
>
> One thing never changes.
>
> The closer you get to the machine, the more data there is. And the less of it actually makes it anywhere useful. 📉
>
> Not because it's hard to collect. But because the path from data to decision has never been built to be secure and contextual at the same time.
>
> Industrial environments are some of the most complex I've tackled. Legacy equipment. Heterogeneous systems. Strict security constraints. No two factories are the same. 🏭
>
> Here's the real problem: for decades, factories have been adding software layer after software layer. Custom integrations to make them talk, building technical debt that compounds every year.
>
> Today, every industrial direction wants to deploy AI. Tomorrow, it'll be something else.
> Same question every time: how many new systems, how many custom integrations before the data finally serves its purpose?
>
> The infrastructure was never designed for this. And patching it won't work anymore. 🛑
> That complexity is exactly why the problem is still unsolved, and exactly why it's worth solving.
>
> In the US, this shift has been underway for years. Europe is catching up. And with our industrial heritage, we have every reason to lead.
>
> That's why, together with Cécilia Tran, we decided to build something different. Not another tool. The foundation.
>
> If you're building or working in this space, or just tired of watching good data not transforming into ROI, let's talk! 🤝
>
> #OT #Manufacturing #DataOps #AI #Industry40 #Polytechnique #MindSetData

**Posts #2–5** — summarized in §6's table above (full technical content not reproduced here; they're long-form RTOS/ESP32 tutorials, publicly viewable on the profile under the "Embedded C Programming" group and via the repost).
