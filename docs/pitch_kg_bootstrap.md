# Pitch narrative update — the Knowledge Graph is auto-generated, not built

_Supersedes the framing implied in earlier materials that pipelines create the Knowledge Graph. Grounded in a live-tested feature (2026-07-20), not a roadmap promise. Background and full decision trail: `docs/analysis_log.md` Entries 87–98._

## The line to use

> **The Knowledge Graph is auto-generated at OPC-UA connect time, via ISA-95 mapping — with human validation before anything counts as confirmed. Pipelines enrich it; they don't build it. Context is there from day one.**

## Why this matters (the problem it answers)

The value of any data platform — automation, AI agents, dashboards, MCP-style integrations — depends entirely on the context existing first. If a customer has to spend weeks manually building pipelines just to get a structured picture of their own plant, the platform's ROI story is backwards: you're asking for the hard part before you deliver any of the value. We removed that step for the structural layer.

## What actually happens (accurate, not aspirational)

1. **Connect** to the plant's OPC-UA server.
2. **Discover** — one click browses the full tag tree. Structural data (which equipment exists, how it's organized) already lives there — it doesn't need to be typed in.
3. That structure is mapped to ISA-95 (Site → Area → Work Center → Equipment → Tag) automatically, using the plant's own tag-naming convention.
4. The graph is written immediately — flagged **pending**, never silently trusted.
5. A human **validates** — accept or reject, node by node — before anything is confirmed structure.
6. From there, **pipelines automate and enrich** an already-existing graph — cost calculation, cause detection, ERP linking. They are not the thing that makes the graph exist.

## What's proven, and what's still ahead (say both — precision builds trust)

**Proven, live-tested against a real OPC-UA server (not a simulation of the pitch):**
- End-to-end flow works: connect → discover → auto-generate → validate → confirmed graph, in under a minute, zero pipelines built.
- A real naming-convention edge case (Site.Line.Machine.Tag, not the originally-assumed shape) was found *and fixed* during testing — the kind of thing that only surfaces against real infrastructure, not in a slide.

**Explicitly not yet true — say this before a technical buyer finds it themselves:**
- This covers the **structural** layer (equipment, hierarchy) only. **Transactional** data — work orders, quality events, costs — still arrives progressively through pipelines, because it doesn't exist until it happens. That's not a gap, it's the nature of the data.
- IT-side (ERP) master data auto-generation is a plausible extension, not built yet.
- Validated at a small tag count (8 tags, 15 graph nodes). Scaling the human-validation UI to a real plant's tag volume (potentially hundreds) is untested.

## The one-sentence version, if you only get one sentence

**"Connect, and you have a structured, human-validated picture of the plant in under a minute — not after weeks of pipeline-building."**

## Visual summary

Published slide: https://claude.ai/code/artifact/23f40e6f-86a4-4968-a1cf-e359d1169f2c

Private by default — share from the artifact's own share menu when ready to send externally.
