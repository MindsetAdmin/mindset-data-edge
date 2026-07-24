# YC Application — Accomplishments section

> **Status: filled in live on the form, verified 2026-07-22.** The four answers below are what's actually saved in the `apply.ycombinator.com` bio-edit form right now (read directly from the page), not drafts — you wrote your own final versions, close to but improved on the earlier drafts. Kept for reference / in case you want to revise later.

> Working doc for the YC founder application "Accomplishments" section (`apply.ycombinator.com/bio/edit`). Drafted with Claude Code, grounded in `docs/linkedin_profile_recommendations.md` §12 (CLock) and the LinkedIn Experience entries logged there (ISPESAGE, Sonatrach, freelance embedded work). Not yet pasted into the live form. Items marked **[NEEDS YOUR INPUT]** are blocking — don't submit until resolved.

---

## Q1 — "Please tell us about a time you most successfully hacked some (non-computer) system to your advantage."

**[READY]**

**Your raw version:**
> During my time at university, I noticed that many students wanted to learn practical embedded systems skills such as FreeRTOS, STM32 development, and hardware debugging, but most learning opportunities were limited to formal coursework. Instead of waiting for a new course to be introduced, I "hacked" the system by creating informal learning groups and sharing resources, project templates, and development environments that allowed students to start building real projects immediately. I also connected students with more experienced peers and industry professionals whenever possible. As a result, several students who initially had little hands-on experience were able to complete embedded systems projects, participate in competitions, and gain practical skills beyond the standard curriculum. The lesson I learned is that many systems are constrained not by rules, but by assumptions. Sometimes the fastest way to improve outcomes is to rethink how existing resources are organized and shared.

**Tightened for the form (same story, trimmed):**
> The curriculum only taught embedded systems in formal courses, so most students never got hands-on with FreeRTOS, STM32, or real hardware debugging until years in. Instead of waiting for a course to be added, I built the shortcut myself: informal study groups with shared project templates and ready-to-go dev environments, plus direct connections to more experienced peers and industry contacts. Students with zero hands-on experience went from nothing to shipping embedded projects and entering competitions — faster than the official curriculum would have gotten them there. It taught me that most systems are limited by assumptions, not actual rules — the fastest fix is often just re-organizing what already exists.

Use the tightened version if there's a character limit; the raw version if there isn't.

---

## Q2 — "Please tell us in one or two sentences about the most impressive thing other than this startup that you have built or achieved."

**[NEEDS A CALL — two drafts below, pick one or merge.]** Also unconfirmed: is the traffic-congestion project (MDI prize) the same team/venture as CLock, or a separate one? Left them as separate below since nothing indicates otherwise.

**Option A — single strongest placement:**
> Won 3rd Prize at the MDI Business School Student Entrepreneur Competition for a smart system to reduce road traffic congestion.

**Option B — placement + AI/security label (shows range: hardware/AI + entrepreneurship):**
> Won 3rd Prize at the MDI Business School Student Entrepreneur Competition for a smart traffic-congestion system, and separately led an AI project — recognized with an official innovation label — that classified security incidents (theft, fire, violence) in real time for automated response.

Not included: Tstart Ooredoo (Finalist) — weaker than the other two for a "most impressive" framing; listed under Q4 instead. Swap it in here if you'd rather lead with it.

---

## Q3 — "Tell us about things you've built before. For example apps you've built, websites, open source contributions. Include URLs if possible."

No public app/website/open-source URLs on record — your build history is embedded systems + contract/freelance work, not public repos. Draft using what's documented on LinkedIn:

> - **CLock** (2020–2021, co-founder, École Nationale Polytechnique, Algiers) — AI security system: embedded hardware (cameras, microcontroller, siren, auto-dial) + computer vision (facial recognition, behavior detection) to catch intrusion/fire/violence in real time and auto-respond (alarm, emergency call, entry lockdown). Built the detection models and the embedded response logic.
> - **Sonatrach** (Feb–Aug 2021, internship) — real-time autonomous embedded system on STM32F103/PIC18F45K22: temperature/pressure/flow sensors driving closed-loop actuators, validated at <50ms response time on an industrial test bench.
> - **ISPESAGE** (May 2022–May 2023, contract) — embedded C/C++ firmware for STM32/PIC industrial weighing systems (UART/I2C/SPI, LCD UI).
> - **Freelance embedded/IoT work** (Dec 2024–May 2026) — ESP32 + FreeRTOS applications: MQTT, ESP-IDF, OTA updates, CI pipelines.
> - **MindSet Data** (current — excluded here per the question, covered elsewhere in the application) — Go edge agent, API server, pipeline engine, and auto-generated ISA-95 Knowledge Graph; live-verified against a real OPC-UA server and a real MySQL/ERP connector.

**[NEEDS YOUR INPUT if any of these have a public artifact]** — a demo video, a repo, a writeup — URLs strengthen this a lot; none are on record right now.

---

## Q4 — "List any competitions/awards you have won, or papers you've published."

No published papers on record. Awards/competitions (from your own list, this session):

> - **3rd Prize — MDI Business School Student Entrepreneur Competition** — smart system to reduce road traffic congestion.
> - **Innovation Label — AI in Security & Access Control** — for leading an AI project (CLock) that classified incidents (theft, fire, violence) for automated response.
> - **Finalist — Tstart Ooredoo Competition** — presented and defended an innovative concept before a jury.

**[CONFIRM]** — exact issuing body/year for the Label, and rough competition size for MDI/Tstart if you know it (adds credibility, but don't want to invent a number).

---

## Open items before this is submit-ready

1. Q1 story — needed from you, nothing to draft from.
2. Q2 — pick Option A or B, confirm CLock vs. traffic-project relationship.
3. Q3 — any public URLs (repo/demo/writeup) for CLock or the freelance work?
4. Q4 — confirm Label issuer/year; optional competition-size context.

---

# Main Application — gap analysis & plan

> Application: `apply.ycombinator.com/apps/5a5f0589-83a5-4d01-b857-3a246c40d7f3` — **MindSet Data, Fall 2026 batch.** Fetched (read-only) 2026-07-22. Founders: Cécilia Tran + Mohamed Khenafif, both profiles marked complete. This section is a plan only — nothing has been typed into the live form.

## Snapshot — what's already answered and solid

These read as strong, keep as-is unless you want polish:
- Founders section (how you met, who codes, not looking for a co-founder)
- Company description (50-char), company URL
- "What is your company going to make" — clear, specific, good use of the MCP/Knowledge Graph angle
- Location + location reasoning (Paris now, Houston TX later — good specific number: $260B TX manufacturing output)
- Progress ("how far along") — matches actual shipped state (OPC-UA/MQTT/SQL connectors, MCP server, Orchestrate builder) reasonably well
- "How long have you been working on this" — 4 months, Mohamed full-time
- "Why this idea / domain expertise" — the 40 discovery calls + France Industrie relationship is a strong, specific credibility signal
- Legal entity formed: yes / Fundraising: yes / Revenue: no / Users: no — all answered

## Gaps — currently blank ("Unanswered")

| Field | Why it matters | Can I draft it? |
|---|---|---|
| **Please provide a link to the product, if any** | YC will click this if present. If nothing's publicly demoable, better to leave blank than link something broken — your call. | No — needs a real URL or explicit "no public link yet" decision from you |
| **What tech stack are you using... Include AI models and AI coding tools** | YC explicitly asks about AI coding tools now — this is an easy, honest, strong answer given how this repo is actually built | **Yes** — see draft below |
| **If you have already participated in an incubator/accelerator...** | Docs mention Boost10x (`docs/advisors.md`) — need to confirm whether that counts as "participated in a program" or was informal advisor/network access, since those are different claims | Draft possible once you confirm Boost10x's actual relationship (formal program vs. informal advisor network) |
| **How do or will you make money? How much could you make?** | Blocking gap — YC weighs this heavily. `docs/mindset.md` §3.0 has real numbers (Motion #1 self-serve <30k€/site agrifood/metallurgy; Motion #2 enterprise 50-150k€/site pharma/cosmetics) | **Yes, draft below** — needs your confirmation these price points are still current, not just aspirational |
| **If you had other ideas you considered applying with** | Optional but YC says explicitly they sometimes fund what's listed here, not the main app | No — needs your actual answer, nothing to infer |
| **Please list all legal entities... state/country** | Straightforward factual gap | No — needs you (which entity, where incorporated) |
| **Equity breakdown among founders/employees, with titles** | Blocking, sensitive — YC wants exact % + titles | No — needs you |
| **Please provide relevant details about your current fundraise** | You said "yes" to fundraising but gave no terms | No — needs you (amount, stage, any commits so far — Entry 84's Polytechnique Ventures pre-seed lead may be relevant context) |
| **What convinced you to apply to YC? Did anyone encourage you? Been to any YC events?** | Low-stakes, easy to answer | Partially — needs your actual motivation/story, I can tighten prose once you give the raw version |
| **How did you hear about Y Combinator?** | Same as above | Needs you |
| **Founder Video / Demo Video** | Both explicitly "No video uploaded" | Not a doc task — recording/upload |

## Weak — answered, but worth strengthening

**"Who are your competitors? What do you understand about your business that they don't?"** — this is the biggest quality gap in the whole application. Current answer:
> *"We are entering the era of the Reasoning Enterprise, where companies no longer buy off-the-shelf software, but build custom, autonomous solutions trained on their own operational data."*

This never names a single competitor and doesn't answer the second half of the question at all. YC asks this specifically to see if you know your landscape — an unanswered "what do you understand that they don't" reads as not having done the homework, which isn't true here: `docs/analysis_log.md` has real competitor research (Entries 7, 66, 88 — MaestroHub, UMH, Cognite, LemonLime, Synapt.ai) and `docs/mindset.md` §15 has a "5 moats" framework built specifically to answer this. This is the single highest-leverage fix available before submitting.

**Draft, grounded in the actual research (not the "Reasoning Enterprise" abstraction):**
> Closest competitors are UMH (open-source UNS/Kafka stack — leaves reconciliation and prioritization to the customer), MaestroHub, and Cognite (does entity contextualization via P&ID/document OCR — a different, more manual problem than ours). What they don't do: reconcile OT events to the *active* production context (the work order/batch running right now) robustly across the multi-hour clock skew typical of mid-market ERPs, then price every event in € and rank it by actionable impact. Competitors stop at clean, contextualized data and hand it to the customer to act on; we cross every event with what makes it matter economically, in real time.

(Numbers/specific vendor framing should be sanity-checked against `docs/analysis_log.md` Entry 7 before use — that research is ~4 weeks old as of this application and competitor positioning can move fast.)

**"What tech stack are you using... AI models and AI coding tools":** verified against `go.mod` / `package.json` directly, not guessed. (Checked also: no CI/CD exists for this repo — no `.github/workflows` — not mentioned in the answer since nothing asked for a weakness inventory. `run.ps1`, the local PowerShell dev launcher, also left out — dev convenience, not product stack, low signal for YC. Docker mentioned generically — deliberately not naming the fake-ERP test simulator specifically, per your steer not to surface internal test-harness details in an external-facing answer.)

> Go for the entire backend — a single edge binary with an embedded database, no external DB dependency. Native, open-source protocol drivers for OT and IT connectivity (industrial protocols plus SQL/ERP integration) instead of paid middleware. React for the pipeline studio and live dashboard, with WebSocket for real-time updates. Docker for containerized dev/test and deployment. We build almost entirely with Claude Code — both founders ship production code through it, which is a real part of how 4 months (largely solo) got this far. Next: a local small language model for on-device tag classification, and an MCP server exposing our Knowledge Graph to AI agents — sovereignty-first, so customer data never has to leave the plant to be queried by AI.

**(Earlier, more implementation-detailed version, kept for reference — not for external use per your steer to stay at category level, not name specific libraries):**
> Go for the entire backend — one edge binary, no CGO (pure-Go SQLite via modernc.org/sqlite), native protocol drivers instead of middleware: gopcua for OPC-UA, Eclipse Paho for MQTT, go-sql-driver/mysql for ERP/MES reconciliation. React 19 + Vite + Tailwind + ReactFlow for the pipeline studio, WebSocket for the live dashboard. Docker for our containerized dev/test environment; production ships as a single Docker container. We build almost entirely with Claude Code — both founders ship production code through it, which is a real part of how 4 months (largely solo) got this far. Next: a local SLM (Phi-3 via Ollama) for on-device tag classification, and an MCP server exposing our Knowledge Graph as tools for Claude/ChatGPT/Copilot-class agents — sovereignty-first, so a customer's data never has to leave the plant to be queried by AI.

**Accuracy check done while drafting this** — worth flagging separately: the code has **no MCP server and no LLM/Ollama integration at all** (`grep -ri mcp **/*.go` → zero hits; no `internal/mcp` or `internal/llm` package exists). That's fine for *this* answer, since it correctly says "planning to use." But the application's own **Progress** section currently states *"Action Layer Ready: Native MCP Server + Orchestrate... ready to feed AI agents"* — "MCP Server" there is not shipped, only "Orchestrate" (the pipeline builder) is real. That's a claim in the live application that doesn't match the code, on a question YC reviewers may technically probe. Recommend revisiting that sentence in Progress separately — not fixed here since you asked me not to modify the application, just flagging it as found.

**"How do or will you make money? How much could you make?":** draft, needs your confirmation the pricing is current —
> Two motions matched to our two GTM segments: self-serve (agrifood, independent metallurgy) at <30k€/site/year, sold directly to the Plant Manager; enterprise IT-led (pharma, cosmetics, grouped metallurgy) at 50-150k€/site/year, 6-12 month RFP cycles. TAM is 15,000+ EU mid-sized factories across our four initial verticals.

## Suggested order of attack

1. **Competitors answer** — highest leverage, fully draftable now, just needs your sign-off on the vendor framing.
2. **Tech stack** — draftable now, low risk, quick win.
3. **Revenue model** — draftable, but confirm the € figures in `mindset.md` still reflect current thinking before using them in a YC application (these are strategy-doc numbers, not signed contracts).
4. **Legal entity / equity / fundraise details** — need you directly, sensitive, can't draft blind.
5. **Accelerator question** — needs a quick call on how to characterize Boost10x.
6. **"Other ideas" / "why YC" / "how did you hear"** — low effort once you give me the raw facts.
7. **Videos + product link** — outside doc scope, your action.
