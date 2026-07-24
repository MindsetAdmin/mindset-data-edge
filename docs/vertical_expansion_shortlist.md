# Beyond Manufacturing — Vertical Expansion Shortlist

_For Cécilia — companion to the outreach pivot discussion (2026-07-16/17). Full internal reasoning and technical build plan logged in `docs/analysis_log.md`, Entries 71–72._

---

## 1. Why look beyond manufacturing now

Manufacturing traction takes time to build, and we don't want to get stuck in long industrial sales cycles before we've proven MindSet Data can scale. The investor feedback has been consistent: "-X% energy optimization for factories" pitches fail because they're a point solution. What lands is the broader claim — **we connect fragmented systems, put a € figure on every event, and tell you what to act on first.**

That claim isn't manufacturing-specific. Manufacturing is just the vertical where the silos happen to include OT (machine-level signals). The underlying pattern — fragmented systems, an event happens, nobody prices it in real time, the decision comes too late — shows up anywhere operations run across disconnected systems with real money on the line.

This doc proposes 4 verticals to test that pattern against, so outreach can find out where it resonates fastest without waiting on manufacturing alone.

## 2. The core thesis, unchanged

> Every platform stops at clean, contextualized data and hands it off for someone else to act on. We cross every event in real time with what makes it matter — and rank it by economic stake, not at random.

This is the pitch already in `Blurb_Invest.md`. Nothing about it is manufacturing-locked. What changes per vertical is only: which systems are the silos, and what the €-priced event looks like.

## 3. How this shortlist was screened

Three filters, in order of importance:

1. **Real pattern fit** — a real-time operational signal, a separate business/records system, and an event that has a genuine € cost nobody prices as it happens.
2. **Technical distance from what's already built** — does the "real-time signal" side look like what we already read (industrial protocols, sensor telemetry), or is it a different world entirely?
3. **Sales-motion fit** — can we still sell direct to an on-site operations lead, fast, under our target deal size — or does the vertical drag us back into procurement-heavy, IT-gatekept cycles (exactly what we're trying to avoid)?

## 4. The shortlist

### 🥇 Warehousing & 3PL operations — strongest fit

- **The silos:** conveyor/sorter control systems, dock-door sensors, and forklift/AGV telemetry on one side; warehouse management (WMS) and transport management (TMS) systems on the other.
- **The unpriced event:** a sorter jam or dock congestion event is invisible to the WMS/TMS until it cascades into a missed truck departure — SLA penalty, overtime, or (cold chain) spoilage. This is the exact "sub-minute stoppage the ERP never sees" story, word for word.
- **Buyer:** warehouse or 3PL site operations manager — same profile as a plant manager, same fast, direct sales motion.
- **Feasibility note:** the real-time signal side speaks a protocol family we already support. This is the closest vertical to what's shipped today — a real pilot could be stood up quickly if a lead converts.
- **Watch-out:** WMS/TMS vendors (Manhattan, Blue Yonder, etc.) are more consolidated than industrial ERP — expect more "why not just use our vendor's module" pushback than in manufacturing.

### 🥈 Commercial real estate / multi-site facilities management

- **The silos:** building management systems (HVAC, energy, access control) on one side; maintenance ticketing and lease/SLA contract terms on the other, usually in entirely separate tools.
- **The unpriced event:** an HVAC or access-control fault isn't priced until someone files a ticket, often hours later — but it's directly tied to SLA penalty clauses or tenant churn risk. Same gap as our "customer-commitment" concept in the Impact Engine — a technical fault with a contractual € consequence nobody surfaces in time.
- **Buyer:** operations director at a facilities-management outsourcer (ISS/Sodexo-type) or a REIT's portfolio ops lead — multi-site, budget authority below IT, fits a per-site pricing model almost unchanged from manufacturing.
- **Feasibility note:** building-system protocols are a genuinely new integration surface for us, but a thinner first version is realistic without full protocol coverage.
- **Watch-out:** € stakes per individual event tend to be smaller than a manufacturing line stoppage — this is a better "many small losses add up" story than a single dramatic number.

### 🥉 Data center / critical infrastructure operations — validate the message, don't build yet

- **The silos:** power and cooling telemetry on one side; uptime SLA contracts on the other.
- **The unpriced event:** this is actually the easiest € story of all four — data-center contracts already price downtime in €/minute. We wouldn't be inventing the number, just catching the anomaly and connecting it to the number that already exists.
- **Buyer risk:** this vertical skews toward larger enterprise and colocation operators with heavy security review — which risks recreating exactly the long sales cycle this whole exercise is meant to escape.
- **Recommendation:** great material for investor conversations and content (the € story sells itself) — not a first outreach cohort.

### Hospital / clinic operations (non-clinical: OR scheduling, bed turnover, equipment) — validate the message, don't build yet

- **The silos:** biomedical equipment and building telemetry on one side; patient records, bed management, and staffing schedules on the other.
- **The unpriced event:** an idle operating room or a down piece of equipment is extremely expensive per hour, and reconciling why is famously manual — the highest pain of all four verticals.
- **Buyer risk:** procurement and compliance review in healthcare will almost certainly reintroduce the long sales cycle we're trying to avoid — the weakest fit of the four despite the strongest pain point.
- **Recommendation:** same as data centers — use it to sharpen the pitch, not as a near-term pilot target.

## 5. Recommendation

**Lead outreach with warehousing/3PL.** It's the closest match to what we've already built and already sell, and the buyer profile preserves our fast, direct sales motion. Run **CRE/facilities management** as a parallel secondary test — nearly as fast to pitch, slightly more integration work if a lead converts.

Treat **data centers** and **hospital operations** as messaging and credibility material — strong stories for investor conversations and content (like the LinkedIn post that sparked this) — rather than active outreach targets for now, since both risk pulling us back into the long sales cycles this pivot exists to avoid.

## 6. Open questions for you

- Do we test messaging on all 4 in outreach copy, or focus first contact on warehousing/3PL + CRE only, and hold the other two for content/investor conversations as recommended above?
- Any existing warm intros or network in warehousing/3PL or facilities management we should prioritize before cold outreach?
- Should we draft vertical-specific hooks (one message variant per vertical) for you to test against the LinkedIn-commenter outreach you're already running?
