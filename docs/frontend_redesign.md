# MindSet Frontend Redesign — Design Brief

> **Purpose**: strip the "AI-generated" look from the current React UI without restructuring the app. Grounded in Grafana + n8n aesthetic. Mohamed executes solo — no intern.
> **Constraint**: same 7 pages, same tech stack, same logo. Skin change only.
> **Timeline**: 3-4 weeks of focused work OR 6-8 weeks interleaved with V1 build.
> **Last updated**: 2026-07-01

---

## 1. Constraints (locked with user 2026-07-01)

| Constraint | Value |
|---|---|
| Page architecture | **Unchanged.** 7 routes stay (Overview / Connect / OpcuaConnect / Compose / Pipelines / Dashboards / KG) |
| Tech stack | **Unchanged.** React 19 + Vite + Tailwind + Zustand + React Flow + Cytoscape + Recharts |
| Logo | **Keep existing** |
| Executor | **Mohamed solo.** No frontend intern (DevOps intern from Entry 42 pick unchanged) |
| Reference designs | **Grafana + n8n** |
| Scope | Design system (skin) — not restructure pages or workflow |

---

## 2. Design principles (10 rules — non-negotiable)

1. **Density over padding.** Grafana / TradingView density, not marketing-site whitespace.
2. **Semantic color over decorative color.** Green = running, red = stopped, amber = warn. Never used for "pretty".
3. **Typography hierarchy — 4 levels maximum.** Not 8 random text-xs/sm/md/lg variations.
4. **Monospace for all numbers.** Tabular figures make dashboards feel engineered.
5. **Icons, never emojis.** 📌 ● ⚠️ → replaced by Lucide icons.
6. **Precise borders (1px only).** No rounded-lg-with-thick-border padded cards everywhere.
7. **Purposeful whitespace.** Space between sections, not inside components.
8. **Consistent panel patterns.** Grafana-style panels on the Dashboard + n8n-style nodes in the Builder.
9. **Custom brand tokens, not Tailwind defaults.** Explicit design language, not "sensible defaults."
10. **Minimal motion.** No fun animations. Transitions <150ms, subtle only.

---

## 3. Design system tokens

### 3.1 Color palette — replace generic dark-slate

**Neutrals (background scale)** — 6 values, evenly spaced perceptually:

```
--bg-canvas     : #0A0A0B   (page background)
--bg-panel      : #131316   (panels / cards)
--bg-panel-alt  : #1B1B1F   (nested panels, table row hover)
--bg-elevated   : #232329   (dropdowns, popovers)
--border-subtle : #2A2A31   (1px dividers)
--border-strong : #3A3A44   (interactive borders)
```

**Text (opacity levels for hierarchy)**:

```
--text-primary   : #E8E8ED  (main content — 90%)
--text-secondary : #A8A8B2  (labels, secondary — 65%)
--text-tertiary  : #6E6E7A  (captions, disabled — 45%)
--text-muted     : #4A4A55  (metadata, timestamps — 30%)
```

**Brand accent — pick ONE color that isn't blue-500**:

Recommendation: **Amber-warm accent** (differentiates from blue-heavy competitors, hints at industrial/mechanical):

```
--accent-primary : #E5A445   (buttons, active states, brand emphasis)
--accent-muted   : #7A5620   (subtle accent — dividers, hover)
```

**Semantic states (from Grafana palette, muted to match dark canvas)**:

```
--status-running : #4ADE80   (Run — subtle green, not neon)
--status-stopped : #F87171   (Stop — muted red)
--status-warn    : #FBBF24   (Warn — amber)
--status-info    : #60A5FA   (Info — cool blue, used SPARINGLY, not the CTA color)
--status-idle    : #6E6E7A   (No signal / disconnected — text-tertiary)
```

**Anti-patterns to remove**:
- `bg-dark-900`, `text-dark-300` → replaced by the token names above
- `hover:bg-blue-500` → replaced by `hover:bg-accent-muted` or subtle border change
- Random gradients → eliminated

### 3.2 Typography

**Primary font**: **Inter** (self-hosted, open source, works at every weight). Load 3 weights only: 400 (regular), 500 (medium), 600 (semibold).

**Monospace font**: **JetBrains Mono** (self-hosted). Used for:
- All numbers in dashboards, KPIs, tables
- Tag values (`Etat_Machine`, `Compteur_Pieces`)
- Timestamps
- Code / config snippets
- OPC-UA node IDs

**Type scale (4 sizes only)**:

```
--text-xs  : 11px  (metadata, timestamps, secondary labels)
--text-sm  : 13px  (body text, form inputs — the default)
--text-md  : 15px  (page section headers)
--text-lg  : 20px  (page titles only)
```

No random 12px, 14px, 16px, 18px. Four sizes. Anywhere you're tempted to invent a new size, pick one of these.

**Line height**: 1.4 for body, 1.2 for headings.
**Letter-spacing**: -0.01em for headings; 0 for body; 0 for monospace.

### 3.3 Iconography

**Icon library**: **Lucide React** (`lucide-react`). Already React-native, clean, matches n8n's aesthetic.

**Size scale**: 14px, 16px, 20px (three sizes — pick from these).
**Stroke width**: 1.5 (consistent everywhere — not 2 for some, 1 for others).

**Emoji removal checklist** (grep the codebase):
- 📌 → `<Pin />` from Lucide
- ● (live indicator) → `<Circle fill="var(--status-running)" size={8} />`
- ⚠️ → `<AlertTriangle />`
- ✅ → `<CheckCircle2 />`
- ❌ → `<XCircle />`
- 📡 → `<Radio />`
- 🛑 → `<Square />` or `<CircleStop />`
- ▶️ → `<Play />`
- 📤 → `<Upload />`

### 3.4 Spacing scale — 5 values

```
--space-1 :  4px
--space-2 :  8px
--space-3 : 16px
--space-4 : 24px
--space-5 : 32px
```

No 6px, 10px, 12px, 20px. When you're tempted, pick from the 5.

### 3.5 Component primitives — adopt shadcn/ui patterns

Install `shadcn/ui` for the primitives (Radix-based composable components). Apply MindSet's design tokens on top. This eliminates the "handmade dropdown that looks slightly off" problem.

Priority primitives:
- `Button` (3 variants: primary, secondary, ghost — nothing more)
- `Input` / `Textarea` / `Select`
- `Dialog` (replacing PickerModal)
- `Tooltip`
- `Popover` (replacing custom hover panels)
- `Tabs`
- `ScrollArea`
- `Command` (for search + palette pattern)

---

## 4. Custom components — MindSet-specific

Beyond shadcn primitives, build these:

### 4.1 `<Panel>` — the Grafana-style panel container

Every dashboard section becomes a `<Panel>`. Standard structure:

```
┌────────────────────────────────────────┐
│ HEADER: title  |  toolbar | actions ┋ │
├────────────────────────────────────────┤
│                                        │
│           panel body                   │
│                                        │
└────────────────────────────────────────┘
```

Props: `title`, `subtitle?`, `toolbar?`, `actions?` (menu), `loading?`, `error?`, `noPadding?`.
Height: fills grid cell.
Border: 1px `--border-subtle`, hover: `--border-strong`.
Background: `--bg-panel`.
Corner radius: 4px (not 12px — tighter, more precise).

### 4.2 `<StatCard>` — dense KPI card (Grafana stat panel)

```
┌────────────────────┐
│ Micro-stops (24h)  │  ← --text-xs --text-secondary
│                    │
│  47                │  ← Monospace, --text-lg, --text-primary
│  ↓ 8 vs yesterday  │  ← --text-xs, semantic color for delta
└────────────────────┘
```

Props: `label`, `value`, `unit?`, `delta?`, `deltaLabel?`, `sparkline?` (optional inline chart).

### 4.3 `<PipelineNode>` — n8n-style node in the Builder

Replace the current ReactFlow node rendering:

```
┌──────────────────────────────┐
│ ● Category    ⚙              │  ← left dot = semantic category color; right = config
│                              │
│  <Icon /> Function Name      │  ← lucide icon + label
│                              │
│  ▸ short description         │  ← --text-xs --text-tertiary
├──────────────────────────────┤
│ ○───          ───○           │  ← input port (left) → output port (right)
└──────────────────────────────┘
```

Category color coding (semantic + brand):
- **Connector**: `#60A5FA` (info blue)
- **Transform**: `#A78BFA` (purple)
- **Calculate**: `#4ADE80` (running green)
- **Condition**: `#FBBF24` (warn amber)
- **Output**: `#E5A445` (accent — outputs are the "goal")

### 4.4 `<StatusDot>` — replace ● emoji + hex circles

```
<StatusDot state="running" />   → filled green circle, 8px
<StatusDot state="stopped" />   → filled red circle, 8px
<StatusDot state="warn" pulse /> → amber circle with subtle pulse animation
```

### 4.5 `<TimeRangeSelector>` — Grafana-style time picker

Standard values: `5m`, `15m`, `1h`, `6h`, `24h`, `7d`, `Custom`.
Position: top-right of DashboardPage.
Visual: segmented control, not dropdown.

### 4.6 `<DenseTable>` — replace card-per-item lists

For OPC-UA tag lists, pipeline lists, event history:
- 32px row height (dense)
- Zebra striping via `--bg-panel-alt`
- Monospace for numeric + ID columns
- 1px `--border-subtle` between rows
- Hover: `--bg-panel-alt`
- Right-aligned numbers
- Left-aligned text
- No card wrapping

---

## 5. Page-by-page redesign priorities

User instruction: **keep the priority pages** — meaning don't reshuffle which pages matter. Mohamed picks execution order. My suggested order (by demo/customer impact):

### Priority 1 — DashboardPage (`/dashboards`) — 3-4 days

The investor demo + customer daily driver.

**Changes**:
- Top row: 4-6 `<StatCard>`s (KPIs — running machines, micro-stops today, downtime €, OEE)
- Middle: grid of `<Panel>`s with dedicated chart types (line for time-series, gauge for OEE, bar for Pareto)
- Right side: `<TimeRangeSelector>` sticky at top
- Bottom: `<DenseTable>` for recent events
- Remove pin-per-widget UI as PRIMARY — keep as advanced feature
- Ship with 3-4 default panels that make sense (not empty state waiting for user to pin)
- Real-time indicator: `<StatusDot state="running" pulse />` in header, not `● live`

**Recharts config**: strip the defaults. Use `--text-tertiary` for grid lines, `--accent-primary` for main series, hide the legend if only one series, monospace ticks.

### Priority 2 — BuilderPage / Compose (`/compose`) — 3-4 days

The pipeline builder — MindSet's differentiator.

**Changes**:
- Replace current ReactFlow nodes with `<PipelineNode>` — clean I/O ports, category color left, config icon right
- Palette (left): grouped by category with color chips, `<Command>`-style search at top
- Config panel (right): tabbed (`Config` / `Inputs` / `Outputs` / `Docs`) instead of scrolling wall
- Canvas: subtle 24px grid, minimum viewport indicators (`Zoom 100%` / `Fit`) bottom-right
- Toolbar: `Save`, `Run`, `Undo/Redo`, `Delete` — Lucide icons + text
- No emojis anywhere

### Priority 3 — OverviewPage (`/overview`) — 1-2 days

**Changes**:
- Grafana-style stat grid (6 KPIs)
- Recent-activity feed (`<DenseTable>` with last 10 events)
- "Getting started" hidden after first pipeline saved (not permanent)
- No emojis

### Priority 4 — KnowledgeGraphPage (`/kg`) — 2 days

**Changes**:
- Cytoscape styles rewritten to match tokens (nodes = `--bg-panel`, edges = `--border-strong`)
- Domain vs Technical toggle: segmented control top-right (like TimeRangeSelector)
- Filter panel: simplified — checkboxes with category color chips
- Node hover: `<Popover>` with node details (not tooltip)

### Priority 5 — OpcuaConnectPage (`/connect/opcua`) — 2 days

**Changes**:
- Connection form: pfSense / Ubuntu network config feel — dense form, monospace inputs
- Discovered tag tree: `<ScrollArea>` with dense rows (like VS Code file tree)
- Selected tags: right panel, `<DenseTable>` with remove buttons
- Status: precise ("Connected · 3 monitored items · 500ms sampling") not vague ("● Connected")

### Priority 6 — ConnectPage (`/connect`) — 1 day

**Changes**:
- Grid of connector cards → `<DenseTable>` of connector rows with icon + name + description + "Select" button
- Removes visual noise

### Priority 7 — PipelinesPage (`/pipelines`) — 1 day

**Changes**:
- List view: `<DenseTable>` with columns (Name / Trigger / Nodes / Last Run / Status / Actions menu)
- Actions consolidated behind a `⋮` menu (not 4 buttons per row)
- Example templates: separate `<Panel>` below the user's pipelines

**Total estimated pages**: ~13-16 days of focused work.

---

## 6. Anti-patterns to actively remove (grep-and-destroy list)

Before starting, spend a day auditing the codebase:

- [ ] Every emoji in JSX (`grep -r "📌\|📡\|▶️\|🛑\|⚠️\|✅\|❌"` — replace with Lucide)
- [ ] Every `bg-dark-*`, `text-dark-*` (replace with CSS custom properties from §3.1)
- [ ] Every `rounded-lg` on containers (replace with `rounded-[4px]` — 4px is our radius)
- [ ] Every `border-dark-700` (replace with `border-subtle` token)
- [ ] Every `hover:bg-blue-500` on buttons (replace with `hover:bg-accent-muted`)
- [ ] Every `text-sm`, `text-xs`, `text-md`, `text-lg` (audit — is it one of our 4 sizes?)
- [ ] Recharts default colors (override via `<CartesianGrid stroke="var(--border-subtle)"/>` and series colors)
- [ ] Non-monospace numbers (audit dashboards for `font-family` on numeric spans)
- [ ] `hover:scale-*` (remove — no fun transforms)
- [ ] Gradient backgrounds (grep `gradient-` — remove all)

---

## 7. Implementation phases (Mohamed solo)

### Phase 1 — Foundation (Week 1)

- Day 1: Audit + destroy list (§6). No new features — just clean up.
- Day 2-3: Install shadcn/ui, wire tokens (§3), set up Inter + JetBrains Mono, replace Tailwind palette in config
- Day 4-5: Build `<Panel>`, `<StatCard>`, `<StatusDot>`, `<DenseTable>`, `<TimeRangeSelector>` primitives with Storybook (or a `/design-system` route to browse)

### Phase 2 — Priority pages (Week 2-3)

- Day 6-9: DashboardPage redesign (Priority 1)
- Day 10-13: BuilderPage redesign (Priority 2)

### Phase 3 — Remaining pages (Week 4)

- Day 14-16: OverviewPage, KnowledgeGraphPage
- Day 17-19: OpcuaConnectPage, ConnectPage, PipelinesPage

Optional Phase 4 — Polish (rolling, throughout V1 development):
- Icon consistency audit
- Empty-state design (Grafana-style: "No data yet — do X to see something")
- Loading states (subtle skeletons, not spinners)
- Error boundaries with actionable messages

---

## 8. Timeline reality check

**Focused solo work**: ~3-4 weeks (15-19 days).

**Interleaved with V1 critical path**: Realistically ~6-8 weeks of calendar time because Mohamed needs to keep working on OF-state Fuzzy Join + MCP server + connectors in parallel.

**V1 ship impact**: adds ~3 weeks to V1 timeline if not carefully sequenced with the AI + connector work.

**Recommendation**: dedicate 2 consecutive focused weeks for Phase 1 + Phase 2 (Dashboard + Builder — the demoable pages). Phase 3 can happen in evenings/weekends over the following 4 weeks without blocking V1 critical path.

---

## 9. Validation

Before declaring "done":

1. **Screenshot before/after every page**. Side-by-side comparison. If the "after" doesn't feel meaningfully better, don't ship it.
2. **Show to Cécilia**. Business-eye test — does it look like a real product?
3. **Show to Djamil via Boost10x WhatsApp**. PMM eye — does the design match the sovereignty + industrial positioning?
4. **Optional**: show to a Plant Manager during customer discovery — does the density feel like their existing SCADA / Grafana, or too spa-like?

---

## 10. What's NOT in scope (explicitly)

- Restructuring page hierarchy (keep 7 routes as-is)
- Changing the logo (kept as-is)
- New pages / features (this is skin only)
- Adding motion / animations beyond subtle transitions
- Marketing site (`mindsetdata.io` is separate — Next.js, addressed by V1 marketing)
- Customer-facing installer / onboarding wizard (V1 build, not redesign scope)
- Ad-hoc Analyst chat UI (V1 build — will be designed when the AI agent lands)

---

## 11. Files that will change

| Category | Files |
|---|---|
| Global | `tailwind.config.js` (extend colors + typography) · `src/index.css` (CSS custom properties) · `src/App.css` (remove dark-* utilities) |
| New components | `src/components/ui/*` (shadcn primitives) · `src/components/panels/*` (Panel, StatCard) · `src/components/StatusDot.jsx` · `src/components/TimeRangeSelector.jsx` · `src/components/DenseTable.jsx` |
| Modified pages | All 7 page files under `src/pages/` |
| Modified components | `DashboardWidgets.jsx` (remove emoji, adopt Panel wrapper) · `LiveDataPanel.jsx` · `NodeConfigPanel.jsx` · nodes (`PipelineNode.jsx`, `TriggerNode.jsx`, `ZoneNode.jsx`) · `NavBar.jsx` · `Palette.jsx` |
| New assets | Inter + JetBrains Mono self-hosted fonts under `src/assets/fonts/` |
| Config | `package.json` (add `lucide-react`, `@radix-ui/*` via shadcn/ui install) |

---

## 12. Recommended immediate next 3 actions

1. **Grep the codebase for emojis** and Tailwind default color classes — get a concrete count of the audit surface (should take 30 min)
2. **Install Inter + JetBrains Mono + Lucide + shadcn/ui base** — foundation ready to build against (~2 hours)
3. **Write the tokens** as CSS custom properties (§3) — non-negotiable "no more `dark-900`" rule (~1 hour)

Then start Phase 1 Day 1.

---

## 13. Intern recommendation impact (from Entry 42 + Entry 44)

**Frontend intern (type #1 with design chops)** — **CANCELLED** per user 2026-07-01. Mohamed does redesign himself.

**DevOps / SRE intern (type #3)** — **STILL RECOMMENDED**. CI/CD + signed binaries + SBOM + Docker + tests are non-optional. Redesign work does not replace infrastructure work.

**Result**: 1 intern in the July-Sep 2026 window — DevOps only.
