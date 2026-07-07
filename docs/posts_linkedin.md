# LinkedIn — Posts + How to Post

**Account:** Mohamed (personal profile) — sole posting account.
**Default language:** English. French is used only for FR-specific audience pushes.

---

## Setup — do this once

1. Change your **current position** to `Building MindSet Data — mindsetdata.io/beta` (or add it as a second current position). This is what shows next to your name everywhere.
2. Update your **About** section — 3 sentences: what you're building, who it's for, one line CTA with `mindsetdata.io/beta`.
3. Add a **Featured** item pointing to `mindsetdata.io/beta` — appears at the top of your profile.
4. Turn on **Creator Mode** (Settings → Creator mode). Adds a "Follow" button, unlocks analytics, boosts reach.
5. Add "Manufacturing" and "Industrial Automation" to your **skills** — helps LinkedIn's algorithm send your posts to the right feed.

---

## How to post — the mechanics

| Element | Recommendation |
|---|---|
| **Length** | 250–500 words. LinkedIn's algorithm rewards depth over brevity, but ≤500 words keeps completion rate up. |
| **Hook** | First 2 lines are all that shows before "See more." Make them land — a contrarian claim, a surprising number, or a quote. |
| **Story arc** | Hook → problem → what we did → lesson → open question. Every anchor post follows this. |
| **Line breaks** | Every 1–3 sentences. Wall-of-text kills mobile completion rate. |
| **Emojis in body** | None. One in the hook line is acceptable. |
| **Beta link** | Near the end, one line, standalone (`→ mindsetdata.io/beta`). Never in line 1 or 2 — LinkedIn suppresses posts with early external links. |
| **Question at the end** | Every post ends with a genuine question. Comments = reach. |
| **Hashtags** | 3 max, at the end. Suggested: `#IndustrialIoT #Manufacturing #Industry40`. Skip if the post already reads focused. |
| **Images** | Optional. If included: a real screenshot (Pipeline Studio, Dashboard) beats a stock photo every time. Never AI-generated art. |
| **Timing** | **Tuesday or Wednesday, 08:00–09:00 CET.** Before French factories start. Avoid Monday morning (LinkedIn is noisy), avoid Friday (dead). |
| **Reply window** | Reply to every comment within 4h for the first 24h. LinkedIn boosts posts whose author responds. |
| **Cadence** | 1–2 posts/week. Consistency > volume. |

---

## Cadence — 1–2 posts / week

- **Tuesday** anchor topic (long-form, story arc)
- **Optional Thursday** — a shorter reaction post or a comment on someone else's post as if it were its own take (300 words max)

Miss two weeks in a row → the algorithm forgets you. If life gets in the way, post ONE line ("Ship day. Back Tuesday with a real post.") — better than silence.

---

## LinkedIn Groups (came back in 2024–25, worth 15 min/week)

Search + join 3–5:

- **Industry 4.0 / Industrial IoT**
- **OPC-UA / MQTT**
- **French manufacturing** — L'Usine Nouvelle group, Directeurs d'Usine
- **Vertical-specific** — pharma manufacturing, cosmetics operations, agrifood processing

**Rule:** answer 2 questions per week in these groups. Never drop the beta link in a body — only if someone directly asks. Groups feed reputation, not signups.

---

## Anti-patterns — never do these

- ❌ Post 3+ times per week — LinkedIn de-prioritizes overposters
- ❌ Copy the same post to Twitter — the algorithm penalizes cross-posts
- ❌ Ask for likes/comments explicitly ("please share!") — this reads amateur
- ❌ Post in Groups the same content you posted on your feed the same day
- ❌ Auto-comment on 20 posts a day — LinkedIn shadowbans engagement pods
- ❌ Buy followers or comments — obvious, kills credibility

---

## Metrics — weekly Friday review (15 min)

| Metric | Month 1 target | Month 3 target |
|---|---:|---:|
| Post reach (avg) | 3,000 | 15,000 |
| Comments per post (avg) | 5 | 25 |
| Follower growth (net/week) | +20 | +100 |
| DMs of interest | 1 | 5 |
| Beta signups attributed (`?ref=linkedin`) | 3 | 20 |

If 3 consecutive months = flat metrics → change the hook style, not the frequency.

---

# The posts

---

## Post 1 — ANCHOR (Day 0, Tuesday, English)

**Publish time:** Tuesday 08:00 CET
**Length:** ~340 words

```
I keep meeting industrial software vendors who confuse "sophisticated" with "useful."

We spent 6 weeks talking to plant managers across pharma, cosmetics, and agrifood factories in France. Every single one said the same thing: they don't need a smarter model, a fancier UI, or a new protocol. They need someone to explain to them, in plain terms, WHY the line stopped for 47 seconds at 14:12, and what that just cost them.

That's it.

So at MindSet Data we made a hard choice: no hyperscaler lock-in, no 6-month integration, no forced cloud. Three deployment modes — Self-Hosted (on your server), On-Premise (on our edge box), or Hybrid (edge + EU-BYOC) — because a pharma lab, a cosmetics logistics hub, and a metallurgy site have different sovereignty and connectivity constraints. OPC-UA in, ISA-95 out, dashboards in the browser. Ships in a day.

The tech under the hood matters — a topological pipeline engine in Go, a knowledge graph built directly on top of the tag registry, and an AI layer that has access to the plant's operational context from day one. But none of that is the point. The point is: the plant manager gets an answer to "why did we lose 23 minutes yesterday" in one click, on the shop floor, without an IT ticket.

Speed of execution, not tech for tech.

If you're running a mid-sized factory in Europe and this resonates, we're opening a small beta cohort in 4 weeks. Priority for pharma, cosmetics, agrifood, metallurgy.

→ mindsetdata.io/beta

What are the "sophisticated" tools your factory quietly stopped using?
```

---

## Post 2 — "Why 3 deployment modes" (Day 7, Tuesday, English)

**Publish time:** Tuesday 08:00 CET
**Length:** ~330 words

```
"You should be cloud-first."

I heard this three times last month from smart investors. Every time from someone who has never set foot in a mid-sized European factory.

Here is what actually happens on the shop floor:

1. A pharma plant in Alsace runs on ISA-95 with GxP audit trails. Sending unit-level production data to a US hyperscaler creates a compliance question they cannot answer today. Cloud-first means: never a customer.

2. An agrifood site in Bretagne has intermittent 4G. When the connection drops, the line does not stop — but their SCADA data does. Cloud-first means: dashboards go blank at the worst moment.

3. A cosmetics logistics hub outside Paris shares its historian with 2 sister plants for OEE benchmarking. They will not send raw production data to a shared cloud tenant. Cloud-first means: dead in the demo.

At MindSet Data we ship three deployment modes, and we let the plant choose:

- Self-Hosted — our binaries on your own server. Zero infrastructure lock-in.
- On-Premise — our sealed edge box, dropped in the electrical cabinet, provisioned by us.
- Hybrid — edge stays local for real-time, EU BYOC handles fleet aggregation and long-term storage.

Same platform, same features, same UI. Different data-residency footprint.

This is not a technical choice. It is a customer-empathy choice. European factories do not want a one-size-fits-all cloud. They want optionality.

Beta cohort opens in 3 weeks. Priority for pharma, cosmetics, agrifood, metallurgy.

→ mindsetdata.io/beta

What sovereignty constraints have blocked your team's tool choices?
```

---

## Post 3 — "The 40-second micro-stop no one measures" (Week 3, Tuesday, English)

**Publish time:** Tuesday 08:00 CET
**Angle:** agrifood story, rules-engine hook

```
An agrifood plant I visited last month runs a bottling line at 42,000 units/hour.

Their ERP said the line hit 91% availability last week. The plant manager said the number was "about right."

We connected our edge platform for 4 days. Rules engine on the raw PLC signals, sub-second resolution.

Real availability: 79%.

The 12-point gap? Sub-minute stops that never made it into the shift report. 23–47-second stalls from label misfeed, cap-torque errors, and a bearing on Line 3 that clicks every 90 minutes. None of it visible in the ERP because the operators cleared them before the 60-second threshold that triggered the stop code.

Multiply 12 points of availability across a 3-shift operation at that throughput. Six figures per week. Invisible.

The reason nobody measures this is not technical. It is architectural. The ERP is the wrong layer to catch a 30-second event, and pushing raw PLC data to the cloud for analysis takes longer than the event itself. Both approaches lose the data or lose the deadline.

Our answer: a rules engine at the edge, on the same box as the OPC-UA gateway, running against ISA-95-normalized tags. It sees the stall, times it, tags it with product / batch / operator context from the plant's own systems, and pushes the enriched event to a local KG in < 200 ms. The dashboard shows it before the operator has closed the alarm window.

That is the difference between reporting the past and understanding the present.

Beta cohort opens in 2 weeks. Priority for pharma, cosmetics, agrifood, metallurgy.

→ mindsetdata.io/beta

What is the sub-minute event your plant is not catching today?
```

---

## Post 4 — "We don't sell AI. We sell context." (Week 4, Tuesday, English)

**Publish time:** Tuesday 08:00 CET
**Angle:** positioning piece — AI without OT context = chatbot

```
Every industrial AI demo I have watched in the last 6 months does the same trick.

You type: "Why did Line 3 slow down yesterday?"
It answers: "Line 3 slowed down yesterday due to reduced throughput."

That is not intelligence. That is a wrapper around a database query.

Real answers require context the AI does not have out of the box. What product was running. What batch. Whether the operator on shift was new. Whether the ambient temperature crossed a threshold that affects viscosity. Whether Line 2 was starving Line 3 upstream. Whether procurement changed suppliers of the input material 3 weeks ago and nobody told the plant.

Any AI without that context is a chatbot with a factory theme.

MindSet is not a smarter model. It is the layer underneath, that gives the model access to the plant's operational context the moment you ask a question — the ISA-95 hierarchy, the schedule, the recipe, the quality data, the maintenance history. Then the model actually has something to reason with.

Speed of execution, not tech for tech. AI without context is a slide, not a product.

Beta cohort opens next week. Priority for pharma, cosmetics, agrifood, metallurgy.

→ mindsetdata.io/beta

What did the last industrial AI demo you watched actually explain?
```

---

## Post 5 — Anchor in French (Day 14, optional)

**Only publish if there's a specific FR audience push planned** — otherwise skip and cycle back to English.

```
"Vous devriez être cloud-first."

Je l'ai entendu trois fois le mois dernier, d'investisseurs intelligents. À chaque fois de quelqu'un qui n'a jamais mis les pieds dans une usine mid-sized européenne.

Voici ce qui se passe vraiment sur le terrain :

1. Une usine pharma en Alsace fonctionne en ISA-95 avec des audit trails GxP. Envoyer des données de production niveau unité vers un hyperscaler américain crée une question de conformité qu'ils ne savent pas répondre aujourd'hui. Cloud-first = jamais client.

2. Un site agroalimentaire en Bretagne a du 4G intermittent. Quand la connexion tombe, la ligne ne s'arrête pas — mais leurs données SCADA, si. Cloud-first = dashboards à zéro au pire moment.

3. Un hub logistique cosmétique près de Paris partage son historian avec 2 sites frères pour benchmarker l'OEE. Ils n'enverront pas leurs données brutes vers un tenant cloud partagé. Cloud-first = mort en démo.

Chez MindSet Data on livre trois modes de déploiement, et on laisse l'usine choisir :

- Self-Hosted — nos binaires sur votre serveur. Zéro dépendance infrastructure.
- On-Premise — notre boîte edge scellée, posée dans l'armoire électrique, provisionnée par nous.
- Hybrid — l'edge reste local pour le temps réel, BYOC UE gère l'agrégation multi-site et le stockage long terme.

Même plateforme, mêmes fonctions, même UI. Empreinte de résidence de données différente.

Ce n'est pas un choix technique. C'est un choix d'empathie client. Les usines européennes ne veulent pas d'un cloud one-size-fits-all. Elles veulent le choix.

Cohorte beta dans 3 semaines. Priorité pharma, cosmétique, agroalimentaire, métallurgie.

→ mindsetdata.io/beta

Quelles contraintes de souveraineté ont bloqué les choix d'outils de votre équipe ?
```

---

## First actions today (~30 min)

1. Add `Building MindSet Data — mindsetdata.io/beta` to your current position
2. Update About section — 3 sentences
3. Turn on Creator Mode
4. Join 3 industrial LinkedIn Groups (Industry 4.0, OPC-UA/MQTT, French manufacturing)
5. Schedule Post 1 for next Tuesday 08:00 CET (or publish live)
