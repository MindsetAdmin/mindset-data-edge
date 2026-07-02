# AI Concepts for MindSet — Beginner Guide

**Audience:** Mohamed (CTO) — strong Go / IoT / embedded engineer, new to AI.
**Goal:** Understand AI Agents, MCP, Native MCP, and agent types — well enough to design + ship V1 + speak to investors with precision.
**Date:** 2026-06-28

---

## 1. Foundation — What is an LLM?

LLM = Large Language Model. Examples: GPT-4, Claude, Phi-3 (MindSet's choice), Mistral, Llama.

Mental model: **a very expensive autocomplete**. Given some input text, it predicts the most probable next chunk of text, one token at a time. That's literally all it does at the base level.

Three properties to internalize:

| Property | What it means | Why it matters for MindSet |
|---|---|---|
| **Stateless** | Every API call is independent. The model has no memory of yesterday's conversation. | You must re-send context every call — or build it back from your KG. |
| **No actions on its own** | The model produces text. It cannot click, query a DB, or call an API by itself. | An LLM alone is not useful in production. You need to wrap it. |
| **Hallucinates** | If asked something it doesn't know, it makes up a plausible-sounding answer. | NEVER trust raw LLM output for factual claims. You MUST ground it in real data. |

Phi-3 specifically: 3.8B parameters quantized to ~2.2GB. Runs on a CPU. Slower than cloud models (8-25 tokens/sec). Weaker reasoning than GPT-4 / Claude. **Good enough** for grounded Q&A, simple instruction-following, and classification. **Not good enough** for complex multi-step reasoning, long conversations, or nuanced French operator slang.

---

## 2. What is an AI Agent?

**An AI agent = LLM + tools + a loop.**

Without tools, an LLM can only talk. With tools, it can act.

The pattern:

```
1. User gives the agent a goal      ("how many micro-stops yesterday?")
2. Agent calls the LLM              ("Plan: I need to query the KG")
3. LLM decides which tool to call   (kg_query(date=yesterday, type=microstop))
4. Tool executes (real code runs)   (SQLite returns 47 rows)
5. Result feeds back to the LLM     ("47 micro-stops, here are the details")
6. LLM produces final answer        ("Yesterday: 47 micro-stops, 312€ lost.")
7. Loop continues if more steps     (or stops here)
```

The "loop" is what makes it an *agent* rather than just an LLM call. Some answers need only step 6. Some need 10 loops (plan → call tool → see result → re-plan → call another tool → ...).

**Analogy:**
- An LLM is like a **chess engine that only suggests moves** but never moves the pieces.
- An Agent is like a **player who actually moves pieces and reads the board after each move.**

### What MindSet's first agent (Ad-hoc Analyst) does

```
Plant Manager types: "How did Line 2 perform yesterday?"
    ↓
Ad-hoc Analyst (running on Phi-3 in your edge agent)
    ↓
Phi-3 plans: "I need stop events for Line 2 in the last 24h"
    ↓
Calls tool: kg_list_events(machine="Line 2", since="-24h")
    ↓
MCP server returns: [{stop1...}, {stop2...}, ..., 47 events]
    ↓
Phi-3 summarizes: "47 micro-stops, top cause = jam (62%), 312€ lost"
    ↓
Answer shows in the dashboard chat, with the KG events cited as sources
```

The agent didn't "know" any of this. It used the tool to find out. **That's the whole trick.**

---

## 3. The problem MCP solves

Before MCP, every AI client had its own way to define and call tools.

| AI client | Tool system |
|---|---|
| Anthropic Claude | Anthropic Tool Use spec |
| OpenAI GPT | OpenAI Function Calling |
| LangChain | LangChain Tools |
| Cursor / Copilot | Their own custom protocol |

So if you wanted your product to be queryable by all of them, you wrote N × M integrations. Same as USB cables before USB-C: every device had its own connector.

**MCP fixes this.** One protocol. Any AI client. Any tool.

---

## 4. What is MCP?

**MCP = Model Context Protocol.**

Open standard published by Anthropic in November 2024. By March 2026, supported by every major AI vendor (Anthropic, OpenAI, Google, Microsoft, AWS, etc.). 97M+ monthly SDK downloads. **De-facto standard for AI ↔ tools.**

Nickname: *"USB-C for the AI world."*

### MCP architecture (3 minutes to grasp)

```
┌──────────────────┐         MCP Protocol           ┌──────────────────┐
│   AI Client      │ ◄──────────────────────────►   │   MCP Server     │
│   (Claude,       │   (JSON-RPC over stdio        │   (your code)    │
│    Copilot,      │    or HTTP/SSE or WebSocket)   │                  │
│    Cursor, ...)  │                                │   Exposes:       │
│                  │                                │   - Tools        │
│                  │                                │   - Resources    │
│                  │                                │   - Prompts      │
└──────────────────┘                                └──────────────────┘
```

### The 3 things an MCP server exposes

1. **Tools** — actions the AI can call (functions with parameters)
   - Example: `kg_query(node_type, filter)`, `kg_list_events(date_range)`
2. **Resources** — read-only data the AI can fetch
   - Example: `kg://nodes/equipment/line-2`, `kg://schema`
3. **Prompts** — reusable prompt templates
   - Example: "explain_micro_stop_cause" template

For V1, MindSet only needs **Tools** to start. Resources + Prompts come later.

### Real-world example for MindSet

You'd expose tools like:

```typescript
// Pseudocode for clarity
mcp_server.tool("kg_query", {
  inputs: { node_type: string, filter?: object },
  output: NodeList,
  handler: (args) => sqlite.query(...)
})

mcp_server.tool("kg_list_events", {
  inputs: { since: string, machine?: string, type?: string },
  output: EventList,
  handler: (args) => sqlite.query(...)
})

mcp_server.tool("kg_cost_summary", {
  inputs: { date_range: string, group_by?: string },
  output: CostBreakdown,
  handler: (args) => costEngine.summarize(...)
})
```

Once these tools are registered, **any** MCP client (Claude Desktop on your laptop, Cursor in your IDE, the Ad-hoc Analyst agent in your dashboard) can call them. **One implementation, infinite clients.** That's the leverage.

---

## 5. What is "Native MCP"?

The phrase has two meanings depending on context. Both important.

### Meaning A — "Native MCP server" (what MindSet has)

A server that **speaks MCP directly**, not via a separate adapter or wrapper.

| Type | Architecture | Example |
|---|---|---|
| **Native MCP** | Your application IS the MCP server. The MCP code is part of `cmd/server`. | MindSet's planned implementation |
| **MCP wrapper / shim** | You have an existing API. You write a separate adapter that translates MCP calls into your existing API calls. | A third-party PostgreSQL MCP server that wraps the PostgreSQL wire protocol |
| **No MCP** | Your app has its own REST/GraphQL API. To use it with AI, the user writes a custom integration. | Most pre-2024 SaaS products |

**Why "native" matters strategically**: a native MCP server is first-class — versioned with the product, no impedance mismatch, no extra deployment. A wrapper is a separate piece that can lag or break. **Investors / technical buyers respect native.**

### Meaning B — "Native MCP client" (what MindSet uses)

An AI tool that **natively understands MCP** — no plugin or extension needed.

| AI client | MCP support |
|---|---|
| Claude Desktop | NATIVE (built-in) |
| Microsoft Copilot | NATIVE (since 2026) |
| ChatGPT (custom connectors) | NATIVE |
| Cursor (IDE) | NATIVE |
| LangChain / older frameworks | Adapter required |

For MindSet's positioning: when you say "any MCP-compatible AI agent (Claude, Copilot, our native Ad-hoc Analyst) can query the factory directly," you're leveraging the native-client side. Customers don't need to write integrations — their existing AI tools just work.

---

## 6. Types of AI Agents — the taxonomy

Agents differ along multiple dimensions. Here are the main types, with MindSet examples.

### Type 1 — Analytic / Q&A Agents (read-only, reactive)

**What they do:** Answer questions by querying data. Read-only.

**Examples in MindSet's 13-agent catalog:**
- **Ad-hoc Analyst (V1 sole agent)** — "How did Line 2 perform yesterday?"
- **Cost Coach (V2)** — explains the cost model to CFO
- **Multi-site Benchmarker (V2)** — "Site A vs Site B on TRS"

**Properties:**
- ✅ Lowest risk — they read, they don't write
- ✅ Easiest to demo
- ✅ Easiest to evaluate (right answer or wrong)
- ⚠️ Phi-3 can handle simple analytic queries. Complex ones (multi-step joins, statistical reasoning) might need a bigger model.

**Why MindSet's V1 ships ONLY this type:** maximum value, minimum risk, easiest to ship in one engineer-quarter.

---

### Type 2 — Monitoring / Watch Agents (proactive, scheduled)

**What they do:** Run on a schedule or react to events. Push insights to the user without being asked.

**Examples:**
- **Daily Briefing Agent (V1.5)** — runs at 6am, sends "last 24h summary" to Plant Manager via email/Slack
- **Alert Triage Agent (V1.5)** — when €-threshold breached, decides priority + draft a recommended action
- **Trend Spotter (V2)** — proactively surfaces "stops on Line 3 doubled in 3 days"

**Properties:**
- ✅ High user value — they tell you things you wouldn't think to ask
- ⚠️ Need stateful design — what was the baseline yesterday? what's "normal"?
- ⚠️ False positives kill trust — better silent than spammy

---

### Type 3 — Conversational Agents (multi-turn, context-aware)

**What they do:** Hold a dialogue. Multi-turn. Need memory of what was just said.

**Examples:**
- **Tribal Knowledge Chatbot (V2)** — interviews operators after a stop: "What did you see? What did you do? Was it the same as last week?"

**Properties:**
- ✅ Best UX when done well
- ❌ HARDEST to ship. Phi-3 local struggles with French operator jargon, interruptions, multi-turn coherence.
- ❌ Need eval harness — easy to drift into nonsense
- For MindSet: V1 ships the **simpler 1-click dropdown + free text** (captures the same dataset — the moat). The chatbot is V2 polish, not the moat.

---

### Type 4 — Action-taking / Workflow Agents (write, execute)

**What they do:** Don't just suggest — actually DO something. Send an email, create a ticket, change a setting.

**Examples:**
- **Maintenance Scheduler (V2)** — "Sensor S3 has 17 false alarms. Want me to draft a maintenance ticket?" (then writes the ticket if confirmed)
- **Connector Recommender (V3)** — "I see your ERP is SAP. Want me to configure the connector?" (then writes the config)

**Properties:**
- ⚠️ HIGH RISK — they can break things. Always need approval flows / dry-run / undo.
- ⚠️ Permission model becomes a major design question. Who's allowed to approve?
- For MindSet: defer to V2+. V1 stays read-only.

---

### Type 5 — Reasoning / Planning Agents (multi-step problem-solving)

**What they do:** Decompose a goal into sub-tasks. Plan a sequence of tool calls. Adapt based on intermediate results.

**Examples:**
- **Causality Reasoner (V2)** — "Stop happened. Query: pressure tag history (-30s). Query: motor temp history. Compare to last 5 stops on this machine. Hypothesis: leak."

**Properties:**
- ✅ Most "intelligent-looking" — when they work, they impress investors
- ❌ Computationally expensive (many LLM calls)
- ❌ Hardest to evaluate (multiple right answers, many wrong ones)
- ❌ Phi-3 will struggle. Likely needs Mistral-Large or Claude — which breaks the local-default sovereignty unless customer opts in.

---

### Type 6 — Orchestration / Multi-Agent Agents

**What they do:** Manage other agents. One agent dispatches work to specialized sub-agents.

**Examples:** Not in MindSet's V1-V2 catalog. Would only emerge if you had 10+ agents and needed coordination.

**Properties:**
- Advanced pattern. Don't try this in V1.
- Risk: cascading failures (one bad agent breaks the whole pipeline)

---

## 7. The dimensions every agent has

When designing or evaluating an agent, think about these axes:

| Dimension | Options | What it affects |
|---|---|---|
| **Read vs Write** | Read-only / Read-write | Security, blast radius, approval flows |
| **Sync vs Async** | Blocks on response / Runs in background | UX, infrastructure complexity |
| **Stateless vs Stateful** | Forgets after each query / Remembers context | Memory architecture, conversation quality |
| **Local vs Remote LLM** | Phi-3 on edge / OpenAI in cloud | Cost, latency, sovereignty, quality |
| **Single-turn vs Multi-turn** | One question = one answer / Dialogue | Complexity, evaluation difficulty |
| **Tool-using vs Pure-text** | Calls external functions / Just talks | Capability, error surface |
| **Grounded vs Free-form** | Cites sources / Speculative | Trust, hallucination risk |

For MindSet V1 (Ad-hoc Analyst), the defaults should be: **Read-only, Sync, Stateless (or short-context), Local LLM (Phi-3), Single-turn (or short multi-turn), Tool-using (MCP), Grounded (cites KG sources)**. That's the safest, easiest-to-ship configuration.

---

## 8. Honest limitations & gotchas (the things vendors don't tell you)

1. **LLMs hallucinate. Always.** Even GPT-4. Even with tools. **Grounding (citing sources from your KG) is the only defense.** Never let the agent free-text answer factual questions without source citations.

2. **Local small models are MUCH weaker than cloud big models.** Phi-3 (3.8B params) is great for classification, simple Q&A, instruction following. It's NOT great for: complex reasoning, long conversations, nuanced language understanding, code generation. Don't promise what Phi-3 can't deliver.

3. **Agents fail in unpredictable ways.** They take wrong tools. They loop infinitely. They give confidently wrong answers. **Build an eval harness early** — a set of fixed questions with known correct answers, run after every change.

4. **"AI-powered" ≠ AI Agent.** Most products that claim "AI" just call an LLM once. A real agent has a tool-calling loop. Don't be fooled by marketing copy in competitor analysis.

5. **MCP is a protocol, not a product.** Implementing MCP doesn't automatically make your product AI-native. The MCP server has to expose USEFUL tools with USEFUL schemas. Schema design is the hard part.

6. **Latency matters.** A local Phi-3 takes 5-30 seconds per query on a 4-core CPU. Plant Managers won't wait 30 seconds for a chart. Set expectations or upgrade hardware.

7. **Cost compounds with usage.** Local Phi-3 = free (electricity). Remote LLM = $0.001 to $0.10 per query depending on model + tokens. If a customer enables remote LLM and runs 10,000 queries/day, that's $10-$1000/day they're paying — make sure the UI shows them this.

8. **Multi-turn agents need a persistent context store.** "Remember our earlier conversation" requires the agent to look up the past. SQLite + an embeddings index is fine for V1; vector DBs (Qdrant, Weaviate) for V2+.

---

## 9. How this maps to MindSet's V1

What you're actually shipping in V1:

| Component | What it does | Status |
|---|---|---|
| **Phi-3 via Ollama** | Local LLM runtime. Health check, model loading, basic prompt execution. | NEW in V1 |
| **MCP server** at `cmd/server` | Exposes KG tools (`kg_query`, `kg_list_events`, `kg_cost_summary`, etc.) to any MCP client | NEW in V1 |
| **Ad-hoc Analyst agent** | Chat UI in dashboard. Type a question, get grounded answer with KG sources cited. Phi-3 by default. | NEW in V1 |
| Optional remote LLM | Customer can configure to call OpenAI/Claude/Mistral instead. UI warns "data leaves your network." | Configurable |

**That's it for V1 AI.** No multi-agent orchestration. No conversational chatbot. No action-taking. No causality reasoning. Just one well-built Q&A agent + MCP server.

**Why minimal V1 is right:**
- 1 engineer can ship it in 6-9 months
- Ad-hoc Analyst alone is enough to claim "AI-native" in investor pitch
- MCP server is enough for the "external AI agents work too" pitch (Claude Desktop demo)
- All other 12 agents in the catalog are V1.5+ / V2 / V3+ — built based on first-customer signals, not pre-decided

---

## 10. Mini glossary

| Term | Meaning |
|---|---|
| **LLM** | Large Language Model. The text-prediction engine. Examples: Phi-3, Claude, GPT-4. |
| **Token** | Smallest unit an LLM reads/writes. ~4 characters in English, more in code. |
| **Context window** | How much text an LLM can "see" at once. Phi-3 = 128K tokens. |
| **Hallucination** | LLM confidently makes up false information. The fundamental risk. |
| **Grounding** | Forcing an LLM to base answers on retrieved data + cite sources. Anti-hallucination technique. |
| **Tool / Function calling** | LLM asks the runtime to call an external function. The basis of all agents. |
| **Agent** | LLM + tools + loop. Can plan, act, observe, re-plan. |
| **MCP (Model Context Protocol)** | Open standard (Anthropic, 2024) for AI ↔ tools. De-facto standard by 2026. |
| **MCP server** | Code that exposes tools/resources via MCP. (Yours, in `cmd/server`.) |
| **MCP client** | AI app that consumes MCP servers. (Claude Desktop, Copilot, your Ad-hoc Analyst.) |
| **Native MCP** | First-class MCP support built into the product, not via a wrapper. |
| **RAG (Retrieval-Augmented Generation)** | Pattern: retrieve relevant data, then ask LLM to answer using ONLY that data. Form of grounding. |
| **Eval harness** | Test suite for agents. Fixed questions with known answers. Run after every change. |
| **Quantization** | Compressing an LLM (e.g., from 16-bit to 4-bit) to run on smaller hardware. Trades quality for speed. |
| **Ollama** | Tool that runs LLMs locally. Wraps Phi-3, Mistral, Llama, etc. MindSet's chosen runtime. |

---

## 11. Recommended reading order (when you have 1-2 hours)

1. [Anthropic's "Building Effective Agents" guide](https://www.anthropic.com/research/building-effective-agents) — 20 min, foundational
2. [Model Context Protocol intro](https://modelcontextprotocol.io/docs/getting-started/intro) — 15 min, official spec overview
3. Pick one MCP server example on GitHub (PostgreSQL or Filesystem) — 20 min reading source, gives you the implementation feel
4. Watch an MCP demo video (e.g., Claude Desktop with a custom MCP server) — 10 min, gives you the "user-facing" feel

You don't need to read papers on transformers / attention / training. You're a CONSUMER of LLMs, not building them. Focus on agent design + MCP integration + grounding/eval. That's 80% of practical AI engineering for products like MindSet.

---

## 12. What I'd suggest as your FIRST hands-on experiment (1 weekend)

To go from concept to muscle memory:

1. Install Ollama. Pull `phi3:mini`. Run `ollama run phi3` and chat with it locally — 10 min. **Feel the speed + the quality.**
2. Write a tiny Go program that POSTs to `http://localhost:11434/api/generate` with a prompt, gets the response — 30 min. **You now know how to call a local LLM from your edge agent.**
3. Read the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) README — 20 min. **You see the shape of an MCP server in Go.**
4. Build a TOY MCP server with ONE tool (`get_current_time`) — 1 hour. **You've shipped an MCP server.**
5. Connect Claude Desktop to your toy MCP server — 30 min. **You've experienced the end-to-end loop.**
6. Replace `get_current_time` with `kg_query` against your real SQLite KG — 2 hours. **You now have a working MindSet edge MCP server.**

After that weekend, you're not a beginner anymore. The rest is iteration + product polish.

---

*Questions? Ask anytime — concepts can be revisited. The goal isn't to memorize this doc; it's to be able to use the terms precisely in your investor deck + engineering decisions.*
