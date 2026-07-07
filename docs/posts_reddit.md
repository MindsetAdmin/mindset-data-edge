# Reddit — Posts + How to Post

**Account:** Mohamed's personal Reddit account.
**Big rule:** Reddit is where credibility is won or lost. It's also where marketing gets shredded on sight. Different rules than Twitter or LinkedIn.

---

## Setup — do this once

1. **Use a real account with existing karma.** Do NOT create a new throwaway — Reddit's spam filter shadowbans low-karma accounts posting external links.
2. **Add `mindsetdata.io/beta` to your Reddit profile bio.** That's the ONLY place the beta link ever lives on Reddit.
3. **Optional:** create a Reddit account flair in each target subreddit if allowed — helps with credibility ("Founder / Building MindSet Data").

---

## THE 5 anti-shill rules — read TWICE before every post

1. **Never post the beta link in the body of a technical post.** Bio only. If someone directly asks in comments, reply with the link. Never earlier.
2. **Lead with the problem, not the product.** The word "we" appears at most 2× in the whole post. "MindSet" appears once at the end, in a parenthetical — or not at all.
3. **Post from a real account with karma** — spend 2 weeks commenting substantively in the target subreddits BEFORE your first long-form post there.
4. **Answer every comment for 48h.** Reddit rewards engagement, punishes drive-bys.
5. **If mods remove your post, do NOT repost.** DM the mods, ask what would make it acceptable, or move on. Reposting = permaban.

---

## How to post — step by step

### Composing a post

1. Go to the subreddit (e.g. `reddit.com/r/OPCUA`).
2. Click **"Create Post"** at the top.
3. Choose format: **Text** (for technical war stories) or **Image** (for screenshots — rare on the subs you'll target).
4. **Title** — max 300 chars but keep it under 100. Frame it as a real question or a specific observation, never a headline.
5. **Body** — 400–800 words. Use paragraph breaks. Use code blocks (three backticks) for code / dataset examples.
6. **Flair** — if the subreddit uses flairs, pick one. "Question" or "Discussion" usually fits.
7. Click **Post**.

### Commenting (the pre-post karma-building phase)

1. Browse `/r/PLC`, `/r/OPCUA`, `/r/golang` — read 2 weeks of top posts.
2. Comment substantively on 3–5 posts per week. Not "great post" — actually add technical value or ask a real question.
3. Aim for 100+ karma in each target subreddit BEFORE posting your own long-form.
4. Never link to `mindsetdata.io` in a comment unless someone explicitly asks. Even then, keep it minimal ("beta signup is in my profile bio if you want it").

### Replying to comments on your own post

1. Click **Reply** below any comment on your post.
2. Answer within 48h — Reddit's algorithm boosts posts with active author replies.
3. Do NOT plug the product in every reply. Answer the technical question first, then move on.

### Reddit awards / upvotes

- **Upvotes** feed the algorithm. Never ask for them. Never buy them.
- **Awards** cost coins. Ignore — they don't build reputation.

---

## Target subreddits — ranked by fit

| Subreddit | Members | Fit | Notes |
|---|---:|---|---|
| **r/PLC** | 300k+ | ★★★★★ | OT engineers. Best audience. Rules-strict but respects real content. |
| **r/OPCUA** | ~8k | ★★★★★ | Tiny but hyper-targeted. Even 200 upvotes here = 20 qualified leads. |
| **r/dataengineering** | 400k+ | ★★★★ | Pipeline engine + KG story fits well. |
| **r/selfhosted** | 400k+ | ★★★★ | On-prem angle. Expect "why not Home Assistant" comments — have an answer ready. |
| **r/golang** | 300k+ | ★★★ | Technical stories about the engine, SQLite choice, MQTT publisher. |
| **r/reactjs** | 500k+ | ★★★ | Frontend war stories (ReactFlow + Zustand + WebSocket dashboards). |
| **r/manufacturing** | 200k+ | ★★★ | Softer angle — micro-stop economics, not tech deep-dive. |
| **r/programming** | 6M | ★ | Skip unless the story is extraordinary. Marketing-hostile crowd. |

---

## Cadence

- **Max 1 long-form post per subreddit per month.** Any faster and mods flag you as promo.
- **1 comment per day** across target subreddits (rotating). Small, consistent, high-signal.
- **Space long-form posts across subreddits** — one Monday in r/OPCUA, next one 2 weeks later in r/golang.

---

## Anti-patterns — never do these

- ❌ Post the beta URL in the body of a technical post
- ❌ Create a new account to escape a mod ban — Reddit tracks IPs
- ❌ Ask for upvotes ("hey guys please upvote if you like this")
- ❌ Cross-post the same content in 3 subreddits the same day
- ❌ Comment "great post!" or "interesting!" — filler comments hurt karma
- ❌ Argue with hostile commenters — thank them for the pushback, respond technically, move on
- ❌ Delete a post that's not going well — leave it, learn, do better next time

---

## Metrics — weekly Friday review

| Metric | Month 1 target | Month 3 target |
|---|---:|---:|
| Karma net gain | +100 | +500 |
| Long-form posts published | 1 | 3 |
| Upvote ratio on long-form (≥) | 85% | 90% |
| Comments on long-form (avg) | 15 | 40 |
| Beta signups attributed (`?ref=reddit`) | 2 | 15 |

If a post hits <70% upvote ratio: it read as promo. Rewrite the next one with less "we" and more question-asking.

---

# The posts

Each block: **subreddit · publish window · title · body · notes**.

---

## Post 1 — r/OPCUA (Day 2, Thursday)

**Publish window:** 09:00–14:00 CET (US morning + EU afternoon)

**Title:**
```
How we auto-mapped 40k OPC-UA browse-path tags into ISA-95 hierarchy (open Q about heuristics)
```

**Body:**
```
Working on an edge platform for mid-sized factories. Ran into a problem I want to share and get feedback on — the OPC-UA server we tested with (Simatic S7 + Kepware gateway) exposes tags with browse names like:

  Usine_Paris_Nord.Ligne1.Machine1.Motor2.Speed
  Usine_Paris_Nord.Ligne1.Machine1.Status
  Usine_Paris_Nord.Ligne1.Machine2.Vibration.X

Naive assumption when I started: split on the first dot to get the machine name. That gave "Usine_Paris_Nord" for every tag. Site name, not machine.

Second attempt: split on the LAST dot to get the attribute, and treat everything before as the machine path. That works but you lose the ISA-95 layering (Site > Area > WorkCenter > WorkUnit > Attribute), which is what the operations team actually cares about.

What ended up working (for us, on this dataset — I'm not claiming universal):

  1. Split on '.'
  2. The LAST segment is always the attribute (Speed / Status / Vibration.X)
  3. The SECOND-TO-LAST segment is the WorkCenter or WorkUnit — the thing the plant manager calls "machine 1"
  4. Segments 0..N-2 map into Site / Area heuristically, unless the user overrides with an explicit ISA-95 mapping in the UI

The manual override matters more than the heuristic. In our workflow, when the OT engineer picks tags to subscribe to, they also pick a routing mode (raw / ISA-95 / both) and we let them explicitly assign site/area/work_center per tag. That mapping is what everything downstream (rules engine, dashboards, KG) uses. The heuristic is only a first guess to make the picker usable.

What I'd love feedback on:

- Anyone dealt with browse-name conventions that DON'T follow dot-notation? I've heard some Rockwell / FactoryTalk deployments use ':' or nested folders.
- Is there a standard OPC-UA companion spec that already publishes ISA-95 as structured metadata, so we could skip the heuristic entirely? I know Companion Spec exists but haven't seen it in the wild.
- For those running Kepware / Ignition / other gateways — do you flatten the hierarchy at the gateway or preserve it? And why?

Happy to answer questions.
```

**Notes:**
- Do NOT include the beta URL in the body — bio only
- 2 weeks before: comment substantively on 3–5 posts in r/OPCUA
- Answer every comment for 48h
- If someone asks "what are you building?" → reply "@mindsetdata.io — beta info in my profile bio, happy to answer questions here"

---

## Post 2 — r/golang (Day 15, Wednesday)

**Publish window:** 09:00–14:00 CET

**Title:**
```
Topological execution engine for a YAML pipeline DSL — ~300 lines of Go, would love a sanity check
```

**Body:**
```
Building an industrial data platform where users compose "pipelines" as YAML — a directed graph of typed nodes (connectors, transforms, calculators, conditions, outputs). Each node has a handler function that runs when its dependencies are satisfied.

Wrote our own execution engine in Go rather than pulling in Airflow / Prefect / Temporal — the whole binary needs to fit and run on an edge box next to a PLC, no external services.

Core loop, roughly:

- Parse YAML → in-memory Pipeline{Nodes []Node, Deps map[string][]string}
- Topological sort with cycle detection
- Execute nodes in order, storing each output in a map keyed by node ID
- When a node runs, merge previous-node outputs + trigger data + node's own config into a single map[string]interface{} passed to the handler
- Wrap every handler call in recover() so a panicking handler can't crash the engine

~300 lines total, no external deps beyond stdlib + a YAML parser.

Things I'm second-guessing:

- Should the "merge everything into one map" pattern be typed instead? It's convenient for handler authors but loses compile-time safety. I keep bouncing between the two.
- Right now dependencies that aren't node IDs (e.g. "trigger") are treated as already-satisfied. Works, but feels hacky.
- No parallelism yet — nodes at the same topological depth still run sequentially. Simple, deterministic, but leaves throughput on the table.

Anyone shipping something similar? Especially curious how you handle the typed-vs-map trade-off — every alternative I try feels worse for one reason or another.
```

**Notes:**
- Same rules: no beta URL in body, comment substantively in r/golang for 2 weeks before, answer every comment for 48h
- Golang crowd is technical AND opinionated. Welcome the pushback.

---

## Post 3 — r/selfhosted (Week 4)

**Publish window:** 09:00–14:00 CET

**Title:**
```
Self-hosted industrial edge platform (v0) — feedback on the docker-compose stack?
```

**Body outline** (draft closer to publish date):

- Set the scene: what industrial edge means (small industrial PC next to a PLC)
- The stack: Mosquitto MQTT broker + SQLite + Go binaries (server + agent) + React frontend, all in one docker-compose
- No external services, no cloud dependency — matches the r/selfhosted ethos
- Trade-offs vs. Home Assistant / Node-RED (have this ready — you WILL get asked)
- Genuine question at the end: what would a v0.5 need to be useful for people running homelabs that touch OT (someone with a PLC in their garage, e.g.)

**Notes:**
- Expect "why not Home Assistant / Node-RED?" as the top comment. Answer respectfully: HA is optimized for smart homes; Node-RED is a flow programming environment. MindSet is optimized for ISA-95 hierarchies + industrial protocol + edge-first data ownership. Not competing — different use cases.

---

## Post 4 — r/dataengineering (Week 5)

**Publish window:** 09:00–14:00 CET

**Title:**
```
Micro-batch rules engine on the edge: sub-second Run↔Stop detection
```

**Body outline**:

- Setup: raw PLC values arrive as MQTT messages. Need to detect transitions (Running ↔ Stopped) with sub-second latency.
- Trade-off: edge rules (fast, opinionated, low-latency) vs. streaming SQL (Flink / Materialize — flexible, higher latency, extra ops burden)
- Our choice: edge rules on the same box as the OPC-UA gateway. Deterministic. No network hop.
- The interesting bit: how we handle the state — one struct per work center, transitions written to an event topic that downstream (KG subscriber, dashboards) consume
- Open question: how do others handle the "small number of very hot state keys" pattern? Anyone tried Materialize for this and regretted it?

**Notes:**
- r/dataengineering respects real trade-off discussion. Avoid making it a MindSet promo — bio link only, technical war story in body.

---

## Coming up (queue for weeks 5+)

- **r/PLC** — "Building a pipeline builder for OT engineers who don't want to write code" (respectful, discussion of visual programming vs. ladder logic)
- **r/reactjs** — "ReactFlow + Zustand + WebSocket for a live industrial dashboard — what I learned"
- **r/manufacturing** (softer subreddit) — "The 40-second micro-stop your ERP never sees" (mirror of the LinkedIn Post 3 angle, reformatted for Reddit)

---

## Reference — accounts + rules

| Item | Value |
|---|---|
| Posting account | Mohamed's personal Reddit |
| Beta URL location | Profile bio ONLY, never in post body |
| UTM | `?ref=reddit` |
| Karma-build period | 2 weeks of substantive commenting before first long-form in each subreddit |
| Reply window on long-form | 48h |
| Cadence | Max 1 long-form / subreddit / month |

---

## First actions today (~10 min)

1. Add `Building MindSet Data — mindsetdata.io/beta` to your Reddit profile bio
2. Join (subscribe to) r/PLC, r/OPCUA, r/dataengineering, r/selfhosted, r/golang
3. Read the top 20 posts of each — get a feel for tone
4. Leave 2 substantive comments today across those subs — start the karma flywheel
