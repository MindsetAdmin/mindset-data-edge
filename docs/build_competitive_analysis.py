"""
Build the MindSet Competitive Analysis workbook (v2).

8 sheets:
  Output A (investor-facing) — sheets 1-4
    1. Positioning
    2. Competitive Matrix (15 dims x 4 competitors)
    3. The 5 Moats
    4. 3 Deployment Editions
  Output B (internal team) — sheets 5-8
    5. Locked Decisions
    6. Open Questions
    7. AI Agent Catalog (13 agents)
    8. Edge vs Cloud Map

Uses Python built-ins only (zipfile + xml). No openpyxl.
"""
import zipfile, io
from xml.sax.saxutils import escape

OUT = r"C:\Users\khena\Desktop\MINDSET\Project\mindset-data-edge\docs\MindSet_Competitive_Analysis_v2_3.xlsx"

# ── palette ──────────────────────────────────────────────────────────────────
C_NAVY    = "FF1F3864"
C_GOLD    = "FFD4A017"
C_WHITE   = "FFFFFFFF"
C_LGREY   = "FFF2F2F2"
C_DGREY   = "FFBFBFBF"
C_GREEN   = "FFE2EFDA"   # MindSet wins
C_DKGRN   = "FF385723"
C_AMBER   = "FFFFF2CC"   # neutral / TBD
C_RED     = "FFFFE0E0"   # competitor wins / problem
C_DKRED   = "FFC00000"
C_BLUE    = "FFD6E4F7"   # info / context
C_PURPLE  = "FFE4D6F7"   # MindSet column highlight

# ── shared strings ──────────────────────────────────────────────────────────
_strings = []
def si(s):
    s = str(s)
    if s not in _strings:
        _strings.append(s)
    return _strings.index(s)

# ── styles ──────────────────────────────────────────────────────────────────
_fills = ["none", "gray125"]
_fonts = []
_xfs   = []

def reg_fill(c):
    if c not in _fills: _fills.append(c)
    return _fills.index(c)

def reg_font(bold=False, color="FF000000", sz=11, italic=False):
    f = (bold, color, sz, italic)
    if f not in _fonts: _fonts.append(f)
    return _fonts.index(f)

def reg_xf(font_id=0, fill_id=0, wrap=True, halign=None, valign="top"):
    x = (font_id, fill_id, wrap, halign, valign)
    if x not in _xfs: _xfs.append(x)
    return _xfs.index(x)

# Fonts
FN     = reg_font(False, "FF000000", 11)
FB     = reg_font(True,  "FF000000", 11)
FH     = reg_font(True,  C_WHITE,    11)         # header white
FHB    = reg_font(True,  C_WHITE,    13)         # big header
FT     = reg_font(True,  C_NAVY,     16)         # title navy
FI     = reg_font(False, "FF595959", 10, True)   # italic small grey
FGR    = reg_font(True,  C_DKGRN,    11)         # bold green
FRD    = reg_font(True,  C_DKRED,    11)         # bold red

# Fills
FI_WHITE  = reg_fill(C_WHITE)
FI_LGREY  = reg_fill(C_LGREY)
FI_NAVY   = reg_fill(C_NAVY)
FI_GOLD   = reg_fill(C_GOLD)
FI_GREEN  = reg_fill(C_GREEN)
FI_AMBER  = reg_fill(C_AMBER)
FI_RED    = reg_fill(C_RED)
FI_BLUE   = reg_fill(C_BLUE)
FI_PURPLE = reg_fill(C_PURPLE)

# Cell formats
XF_NORMAL    = reg_xf(FN, FI_WHITE)
XF_BOLD      = reg_xf(FB, FI_WHITE)
XF_LGREY     = reg_xf(FN, FI_LGREY)
XF_BOLD_LG   = reg_xf(FB, FI_LGREY)
XF_HDR       = reg_xf(FH, FI_NAVY,  halign="center")
XF_HDRGLD    = reg_xf(FH, FI_GOLD,  halign="center")
XF_TITLE     = reg_xf(FT, FI_WHITE, halign="left")
XF_GREEN     = reg_xf(FGR,FI_GREEN)
XF_AMBER     = reg_xf(FB, FI_AMBER)
XF_RED       = reg_xf(FRD,FI_RED)
XF_BLUE      = reg_xf(FN, FI_BLUE)
XF_PURPLE    = reg_xf(FB, FI_PURPLE)
XF_ITALIC    = reg_xf(FI, FI_WHITE)
XF_NORMAL_C  = reg_xf(FN, FI_WHITE, halign="center")
XF_BOLD_C    = reg_xf(FB, FI_WHITE, halign="center")

# ── address helpers ─────────────────────────────────────────────────────────
def col_letter(n):
    s = ""
    while n:
        n, r = divmod(n - 1, 26)
        s = chr(65 + r) + s
    return s

def addr(r, c):
    return f"{col_letter(c)}{r}"

# ── XML builders ────────────────────────────────────────────────────────────
def styles_xml():
    fxml = ""
    for (bold, color, sz, italic) in _fonts:
        b = "<b/>" if bold else ""
        i = "<i/>" if italic else ""
        fxml += f'<font>{b}{i}<sz val="{sz}"/><color rgb="{color}"/><name val="Calibri"/><family val="2"/></font>'

    fillxml = ""
    for f in _fills:
        if f == "none":
            fillxml += '<fill><patternFill patternType="none"/></fill>'
        elif f == "gray125":
            fillxml += '<fill><patternFill patternType="gray125"/></fill>'
        else:
            fillxml += (f'<fill><patternFill patternType="solid">'
                        f'<fgColor rgb="{f}"/><bgColor indexed="64"/>'
                        f'</patternFill></fill>')

    borders_xml = '<border><left/><right/><top/><bottom/><diagonal/></border>'

    xfxml = ""
    for (font_id, fill_id, wrap, halign, valign) in _xfs:
        al = []
        if wrap:   al.append('wrapText="1"')
        if halign: al.append(f'horizontal="{halign}"')
        if valign: al.append(f'vertical="{valign}"')
        alxml = f'<alignment {" ".join(al)}/>' if al else ""
        xfxml += (f'<xf numFmtId="0" fontId="{font_id}" fillId="{fill_id}" '
                  f'borderId="0" xfId="0" applyFont="1" applyFill="1" applyAlignment="1">'
                  f'{alxml}</xf>')

    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">'
        f'<fonts count="{len(_fonts)}">{fxml}</fonts>'
        f'<fills count="{len(_fills)}">{fillxml}</fills>'
        f'<borders count="1">{borders_xml}</borders>'
        '<cellStyleXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellStyleXfs>'
        f'<cellXfs count="{len(_xfs)}">{xfxml}</cellXfs>'
        '<cellStyles count="1"><cellStyle name="Normal" xfId="0" builtinId="0"/></cellStyles>'
        '</styleSheet>'
    ).encode("utf-8")

def sst_xml():
    items = "".join(f'<si><t xml:space="preserve">{escape(s)}</t></si>' for s in _strings)
    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        f'<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" '
        f'count="{len(_strings)}" uniqueCount="{len(_strings)}">{items}</sst>'
    ).encode("utf-8")

# ── Sheet class ─────────────────────────────────────────────────────────────
class Sheet:
    def __init__(self, name):
        self.name = name
        self._cells = {}
        self._col_widths = {}
        self._row_heights = {}
        self._merges = []
        self._freeze_row = 0

    def write(self, r, c, value, xf_id=XF_NORMAL):
        self._cells[(r, c)] = (xf_id, str(value))

    def set_col_width(self, c, w):
        self._col_widths[c] = w

    def set_row_height(self, r, h):
        self._row_heights[r] = h

    def merge(self, r1, c1, r2, c2):
        self._merges.append(f"{addr(r1,c1)}:{addr(r2,c2)}")

    def freeze(self, row):
        self._freeze_row = row

    def to_xml(self):
        if self._freeze_row:
            top = addr(self._freeze_row + 1, 1)
            sv = ('<sheetViews><sheetView workbookViewId="0">'
                  f'<pane ySplit="{self._freeze_row}" topLeftCell="{top}" '
                  f'activePane="bottomLeft" state="frozen"/>'
                  f'<selection pane="bottomLeft" activeCell="{top}" sqref="{top}"/>'
                  '</sheetView></sheetViews>')
        else:
            sv = '<sheetViews><sheetView workbookViewId="0"/></sheetViews>'

        col_xml = ""
        for c, w in sorted(self._col_widths.items()):
            col_xml += f'<col min="{c}" max="{c}" width="{w}" customWidth="1"/>'
        if col_xml:
            col_xml = f"<cols>{col_xml}</cols>"

        by_row = {}
        for (r, c), (xf_id, val) in self._cells.items():
            by_row.setdefault(r, []).append((c, xf_id, val))

        rows_xml = ""
        for r in sorted(by_row):
            cells_xml = ""
            for c, xf_id, val in sorted(by_row[r]):
                idx = si(val)
                cells_xml += f'<c r="{addr(r,c)}" s="{xf_id}" t="s"><v>{idx}</v></c>'
            ht_attr = ""
            if r in self._row_heights:
                ht_attr = f' ht="{self._row_heights[r]}" customHeight="1"'
            rows_xml += f'<row r="{r}"{ht_attr}>{cells_xml}</row>'

        merge_xml = ""
        if self._merges:
            merge_xml = "<mergeCells>" + "".join(f'<mergeCell ref="{m}"/>' for m in self._merges) + "</mergeCells>"

        return (
            '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
            '<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">'
            f'{sv}{col_xml}<sheetData>{rows_xml}</sheetData>{merge_xml}'
            '</worksheet>'
        ).encode("utf-8")


# ════════════════════════════════════════════════════════════════════════════
# DATA
# ════════════════════════════════════════════════════════════════════════════

COMPETITORS = ["MindSet Data", "UMH", "MaestroHub", "Cognite"]

# ── Sheet 1 — Positioning ───────────────────────────────────────────────────
s1 = Sheet("1. Positioning (INV)")
s1.set_col_width(1, 4)
s1.set_col_width(2, 110)
s1.set_col_width(3, 4)

# Title
s1.write(1, 2, "MindSet Data — Competitive Positioning", XF_TITLE)
s1.merge(1, 2, 1, 3)
s1.set_row_height(1, 36)

s1.write(2, 2, "AUDIENCE: investors + internal team   |   Date: 2026-06-27   |   Version: v2", XF_ITALIC)
s1.set_row_height(2, 18)

# One-paragraph positioning
s1.write(4, 2, "THE ONE-PARAGRAPH POSITIONING", XF_HDR)
s1.set_row_height(4, 22)
positioning = (
    "MindSet Data is the AI-native edge industrial PLATFORM for mid-market "
    "European manufacturers (ETI). One Docker command installs a single Go binary "
    "on a customer PC; in 48 hours the platform auto-discovers OT equipment, "
    "contextualizes data into an ISA-95 Unified Namespace, attributes every OT event "
    "to its active Fabrication Order via ERP-state reconciliation (robust to multi-hour "
    "clock skew, unlike time-window joins), and lets any MCP-compatible AI agent "
    "(Claude, Copilot, our native Ad-hoc Analyst) query the factory directly — without "
    "raw data ever leaving the customer network. V1 ships with 3 ready-to-use templates "
    "(micro-stop, energy waste, OEE/TRS); customers + their AI agents build more use "
    "cases on top. Three editions (On-Premise / Hybrid / Self-Hosted), all EU-sovereign, "
    "NEVER on US hyperscalers. Single-vendor, no per-tag fees, no Kepware-style "
    "middleware. Customer owns the cumulative site fingerprint."
)
s1.write(5, 2, positioning, XF_BLUE)
s1.set_row_height(5, 180)

# Three vs-statements
s1.write(7, 2, "POSITIONING VS THE 3 NAMED COMPETITORS", XF_HDR)
s1.set_row_height(7, 22)

vs_statements = [
    ("vs UMH (closest direct rival)",
     "UMH is the toolkit. MindSet is the turnkey product. UMH ships 6 open-source "
     "projects (Benthos + Node-RED + HiveMQ + Redpanda + TimescaleDB + Grafana) "
     "glued together on a Kubernetes cluster — you operate it, write Node-RED flows, "
     "build dashboards. MindSet ships one Go binary that already includes the rules "
     "engine, cost model in €, Fuzzy Join, AI agents, and dashboards — designed for "
     "an ETI Plant Manager to install in 48h on a single PC.  Plus: UMH enterprise = "
     "36k€/site/year — above MindSet's <30k€ Plant Manager signing threshold."),
    ("vs MaestroHub (different segment)",
     "MaestroHub targets the enterprise auto/chemicals/metals factory with an IT "
     "department and an AWS partnership. MindSet targets the 50-200 person ETI "
     "factory where the Plant Manager owns the buying decision and EU sovereignty "
     "is non-negotiable. MaestroHub partners with hyperscalers; MindSet explicitly "
     "rules them out."),
    ("vs Cognite (category benchmark, different segment)",
     "Cognite is the cloud platform for oil & gas majors. MindSet is the edge "
     "platform for European manufacturers. Both now have MCP — but Cognite's MCP "
     "runs inside their cloud (raw data leaves the factory). MindSet's MCP runs "
     "at the edge (data stays where it was generated). Cognite is enterprise "
     "contracts six-to-seven figures; MindSet is plant-manager-budget."),
]
row = 9
for title, body in vs_statements:
    s1.write(row, 2, title, XF_BOLD_LG)
    s1.set_row_height(row, 22)
    s1.write(row+1, 2, body, XF_NORMAL)
    s1.set_row_height(row+1, 100)
    row += 3

# V1 starter templates
s1.write(row, 2, "V1 — 3 STARTER USE-CASE TEMPLATES (CUSTOMERS BUILD MORE)", XF_HDR)
s1.set_row_height(row, 22)
row += 1
s1.write(row, 2,
    "1. MICRO-STOP DETECTION + COST IN €  —  Run→Stop→Run with 30s<duration<3min. Output: '47 micro-stops yesterday = 312€ lost on Line 2'.\n"
    "2. ENERGY WASTE DETECTION  —  Energy consumption when machine is stopped. Level 1 (alert only) works without ERP — fast week-1 ROI. Level 2 ('OF#456 wasted 18€ of steam') activates once ERP connector is configured.\n"
    "3. OEE / TRS DASHBOARD  —  Real availability vs declared. Output: 'Real TRS is 74%, not 88%. Eliminating jams recovers +5 TRS points = X€/week'.\n\n"
    "Plant Manager sees 3 working templates on day 1. The PLATFORM (rules engine + cost model + AI agents + KG + MCP) lets them build quality / changeover / predictive / any custom use case in days, not months. No vendor RFP, no consulting hours.",
    XF_PURPLE)
s1.set_row_height(row, 180)
row += 2

# OEE detection mechanism (Plant Manager + investor demo)
s1.write(row, 2, "HOW WE DETECT REAL OEE vs DECLARED OEE  (the strongest single demo)", XF_HDRGLD)
s1.set_row_height(row, 22)
row += 1
s1.write(row, 2,
    "DECLARED OEE = what the operator manually reports to MES/ERP. Typically optimistic — micro-stops not counted, downtime miscategorized. Usually 5-15 points HIGHER than reality.\n\n"
    "REAL OEE = what MindSet calculates from raw OT data:\n"
    "  • AVAILABILITY = (Planned_Time - Major_Stops - Micro-Stops) / Planned_Time  ← every state transition measured by the rules engine on OPC-UA Etat_Machine\n"
    "  • PERFORMANCE = Actual_Output / Theoretical_Output  ← Compteur_Pieces vs Cadence configured in cost wizard\n"
    "  • QUALITY = Good_Parts / Total_Parts  ← MES integration in V1.5; V1 uses customer-estimated defect rate\n\n"
    "THE PITCH: 'Your declared OEE is 88%. Your REAL OEE is 74%. The 14-point gap = 1h04 of hidden downtime/week = X€/week. Here is the Pareto of causes — top 3 fixes recover Y€.' The GAP itself is the value proposition.",
    XF_GREEN)
s1.set_row_height(row, 220)
row += 2

# Three-edition headline
s1.write(row, 2, "THE 3-EDITION DEPLOYMENT MODEL (NO HYPERSCALER EDITION — BY DESIGN)", XF_HDR)
s1.set_row_height(row, 22)
row += 1
s1.write(row, 2,
    "AIR-GAP  (defense, public sector, sensitive pharma)  —  ZERO cloud, per-site only.\n"
    "SOVEREIGN CLOUD  (default)  —  Scaleway FR / OVH FR for cross-site KG + remote dashboard + encrypted backup.\n"
    "BYOC  (large multi-site customers)  —  Customer deploys cloud tier on Hetzner / IONOS / T-Systems / 3DS Outscale / their own on-prem Kubernetes.\n\n"
    "AWS, Azure, GCP (including their EU regions) are explicitly EXCLUDED. The US CLOUD Act exposure invalidates the sovereignty pitch for public sector & defense.",
    XF_BLUE)
s1.set_row_height(row, 140)


# ── Sheet 2 — Competitive Matrix ───────────────────────────────────────────
s2 = Sheet("2. Comp Matrix (INV)")
s2.freeze(3)
s2.set_col_width(1, 28)
for c in range(2, 6): s2.set_col_width(c, 38)
s2.set_col_width(6, 22)

# Title row
s2.write(1, 1, "Competitive Matrix — 15 dimensions, sovereignty as the lens", XF_TITLE)
s2.merge(1, 1, 1, 6)
s2.set_row_height(1, 30)

s2.write(2, 1, "Color: GREEN = MindSet advantage, RED = competitor advantage, AMBER = neutral/parity. Read 'Advantage' column.", XF_ITALIC)
s2.merge(2, 1, 2, 6)

# Header row
s2.write(3, 1, "Dimension", XF_HDR)
s2.write(3, 2, "MindSet Data", XF_HDRGLD)   # gold = us
s2.write(3, 3, "UMH", XF_HDR)
s2.write(3, 4, "MaestroHub", XF_HDR)
s2.write(3, 5, "Cognite", XF_HDR)
s2.write(3, 6, "Advantage", XF_HDR)
s2.set_row_height(3, 28)

# Each row: (dimension, mindset, umh, maestrohub, cognite, advantage)
matrix = [
    ("1. Sovereignty / data jurisdiction",
     "EU-sovereign by design. 3 editions, all excluding US hyperscalers. Raw data NEVER leaves customer network.",
     "Flexible deployment (on-prem / edge / cloud) but no opinionated sovereignty stance.",
     "AWS partner — hyperscaler-friendly. Sovereignty not a stated positioning.",
     "Cognite cloud only — proprietary SaaS. No sovereignty options.",
     "MindSet"),
    ("2. Open-source license + governance",
     "Proprietary (closed source) for first 2 years. Open-core / source-available options reconsidered in 2028.",
     "Apache 2.0 (relicensed from AGPL v3 in 2025). Free community edition + paid enterprise.",
     "Proprietary. Closed source.",
     "Proprietary. Closed source.",
     "UMH"),
    ("3. Deployment model",
     "Single Go binary in 1 Docker container. Edge-first. Cloud tier optional.",
     "Kubernetes Helm chart with 6+ components (Benthos, Node-RED, HiveMQ, Redpanda, TimescaleDB, Grafana). Heavy.",
     "Edge-to-cloud — strong cloud-side processing.",
     "Cloud-mandatory. Thin edge extractors.",
     "MindSet"),
    ("4. Target segment",
     "ETI mid-market manufacturing (50-200 person factories, <30k€/site signing).",
     "OSS-first — system integrators, manufacturing.",
     "Enterprise manufacturing (auto, appliances, chemicals, metals).",
     "Large enterprises (oil & gas, energy, utilities, six-to-seven-figure contracts).",
     "MindSet"),
    ("5. Geography",
     "France / EU first. EU sovereignty as core positioning.",
     "Germany / EU.",
     "EMEA expanding to NA via integrators.",
     "Norway / global. US presence.",
     "MindSet+UMH (EU rooted)"),
    ("6. Cost model / pricing",
     "TBD — likely free community Edge Agent + paid cloud + enterprise support. Edge-first means zero mandatory cloud spend.",
     "Free community edition + Enterprise: 36k€ / year / factory. ABOVE MindSet's <30k€ Plant Manager threshold.",
     "Not publicly listed — enterprise sales.",
     "Enterprise contracts (six-to-seven figures).",
     "MindSet (vs UMH+others)"),
    ("7. Connectors / protocol coverage",
     "V0: OPC-UA + Modbus. V1+: S7, SQL, REST, Files, FTP, MQTT, Sparkplug. Roadmap: 28+ protocols.",
     "Wide via Benthos + Node-RED ecosystem (Node-RED has 4000+ community nodes).",
     "40+ industrial protocols natively. Strong here.",
     "Many connectors via Cognite Extractors, but cloud-centric.",
     "MaestroHub (currently)"),
    ("8. UNS support",
     "Native — ISA-95 hierarchy, auto-generated from discovery.",
     "Native — ISA-95 UNS via MQTT (HiveMQ) + Kafka (Redpanda).",
     "Native — Unified Namespace as the core product.",
     "Industrial Knowledge Graph (proprietary data model), not strictly UNS.",
     "MindSet+UMH+MaestroHub (parity)"),
    ("9. OT/IT reconciliation (Fuzzy Join)",
     "NATIVE — OF-state-based attribution. Polls ERP for active OFs; tags every OT event with active OF. ROBUST TO MULTI-HOUR CLOCK SKEW (where time-window joins fail).",
     "Not built-in. User implements via Node-RED.",
     "Not visible as a dedicated feature.",
     "Cognite has entity contextualization (P&ID OCR, asset matching) — different problem, not OT/IT temporal attribution.",
     "MindSet"),
    ("10. Cost-in-€ at the edge",
     "NATIVE — Cost model V0 (3-field manual: hourly cost, cadence, margin). V1: auto-import from ERP.",
     "Not built-in. User builds via Grafana / external code.",
     "Not visible as a dedicated feature.",
     "Not at the edge — done in cloud reports.",
     "MindSet"),
    ("11. AI / LLM integration",
     "Local Phi-3 via Ollama by default. Optional remote LLM with explicit data-disclosure UI warning. Native MCP server AT THE EDGE.",
     "No native AI features. User integrates own AI.",
     "AI-ready data pipelines on public site. MCP support claimed by CEO in podcast — not on public docs (unverified, likely cloud-side if real).",
     "Cognite Atlas AI — proprietary AI agents in their cloud. MCP via Function Apps (cloud-only).",
     "MindSet (edge MCP is unique)"),
    ("12. Self-serve deployment time",
     "48h on-site (POC promise). Docker pull + run.",
     "Days-to-weeks (Kubernetes cluster + 6-component config + connector setup).",
     "60-day pilot — enterprise integration project.",
     "Months (enterprise integration project).",
     "MindSet"),
    ("13. Vendor lock-in",
     "LOW — Apache 2.0, open standards (MQTT, OPC-UA, ISA-95, SQLite). Customer owns historian.",
     "LOW — Apache 2.0, OSS stack. Customer owns everything.",
     "HIGH — proprietary platform.",
     "HIGH — proprietary data model, AI layer, storage all Cognite.",
     "MindSet+UMH (parity)"),
    ("14. Maturity",
     "Pre-POC / early. Implementation just starting. Strong vision, limited proof points.",
     "Funded, has customers, active GitHub community, mature OSS project.",
     "Commercially deployed across automotive/chemicals/metals — proven at enterprise.",
     "Mature SaaS, proven at scale (Aker BP, etc.).",
     "Cognite > UMH > MaestroHub > MindSet"),
    ("15. Compliance",
     "RGPD native, NIS2 architecture-ready. ISO 27001 / SOC 2 = roadmap.",
     "Enterprise edition: Audit Trail, RBAC, SSO. SOC 2 / ISO not stated.",
     "Standard enterprise compliance for a EU vendor.",
     "Enterprise compliance package (SOC 2, ISO 27001 likely).",
     "Cognite (currently), MindSet on track"),
    ("16. Resources needed (edge)",
     "1 Docker container on a single PC. MIN: 4 CPU / 8GB RAM / 50GB SSD. RECOMMENDED: 8 CPU / 16GB RAM / 100GB SSD (Phi-3 + buffers).",
     "Kubernetes cluster (1-3 nodes). 16GB RAM / 8 CPU PER NODE minimum. Significant DevOps overhead, K8s expertise required.",
     "Single host: 8 CPU / 16GB RAM / 200GB SSD (vendor spec — confirmed). Comparable to MindSet recommended, 2x storage.",
     "Thin extractor only at edge (~minimal). VAST cloud resources hosted at Cognite — customer pays the bill.",
     "MindSet (lightest at minimum spec)"),
]

row = 4
for (dim, mds, umh, mh, cog, adv) in matrix:
    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if row % 2 == 0 else XF_BOLD_LG
    s2.write(row, 1, dim, bg_b)
    s2.write(row, 2, mds, XF_PURPLE)   # us = always purple highlight
    s2.write(row, 3, umh, bg)
    s2.write(row, 4, mh, bg)
    s2.write(row, 5, cog, bg)
    if "MindSet" in adv and "+" not in adv:
        s2.write(row, 6, adv, XF_GREEN)
    elif "MindSet" in adv:
        s2.write(row, 6, adv, XF_AMBER)
    elif "Cognite" in adv or "UMH" in adv or "MaestroHub" in adv:
        s2.write(row, 6, adv, XF_RED)
    else:
        s2.write(row, 6, adv, XF_AMBER)
    s2.set_row_height(row, 80)
    row += 1


# ── Sheet 3 — The 5 Moats ──────────────────────────────────────────────────
s3 = Sheet("3. The 5 Moats (INV)")
s3.set_col_width(1, 4)
s3.set_col_width(2, 30)
s3.set_col_width(3, 70)
s3.set_col_width(4, 4)

s3.write(1, 2, "The 5 Defensive Moats", XF_TITLE)
s3.merge(1, 2, 1, 3)
s3.set_row_height(1, 30)

s3.write(2, 2, "Proposed revision of the 4 moats in docs/mindset.md section 15. Adds 'Edge sovereignty + MCP' as moat #5 based on new competitive findings.", XF_ITALIC)
s3.merge(2, 2, 2, 3)

s3.write(4, 2, "Moat", XF_HDR)
s3.write(4, 3, "What it is + why it's defensible", XF_HDR)
s3.set_row_height(4, 24)

moats = [
    ("1. Zero-manual auto-discovery",
     "Network scan + behavioral inference (10-15 min live pattern matching) + Phi-3 SLM classification. "
     "Auto-identifies Etat_Machine, Compteur_Pieces, Vitesse on opaque Modbus/S7 registers without tag names. "
     "48h deployment vs 3 months at a typical integrator. STRUCTURAL ADVANTAGE on client acquisition cost."),
    ("2. The Fuzzy Join (OT/IT attribution via OF state)",
     "MindSet attributes every OT event to its active Fabrication Order by reading OF STATE from the ERP — "
     "not by joining on timestamps. ROBUST TO MULTI-HOUR CLOCK SKEW typical of mid-market ERPs (where operators "
     "update records end-of-shift). UMH leaves OT/IT join to the user (Node-RED), MaestroHub doesn't address it as "
     "a dedicated feature, Cognite does entity contextualization (P&ID OCR, asset matching) — a different problem. "
     "Hard to build, invisible from outside — the most defensible component in the stack."),
    ("3. Cumulative site fingerprint (the KG)",
     "Site-specific non-replicable context: every micro-stop, every cause, every cost calibration accumulates. "
     "Replacing MindSet = losing all accumulated intelligence (causes per machine, energy patterns, "
     "calibrated cost models). Churn becomes structurally prohibitive after month 6."),
    ("4. Tribal knowledge structured (ships V1 — moat = the dataset, not the UX)",
     "V1 — 1-click cause dropdown + free-text on every stop event. The DATASET (sensor pattern → operator label) accumulates from day 1. "
     "After 6 months on-site, IMPOSSIBLE to reconstruct without real-time access — no competitor can copy post-hoc. "
     "V2 polish: Phi-3 conversational interview for richer capture. But the MOAT is the dataset, not the chatbot — both ship the same defensible asset. "
     "Compounds with the site fingerprint (Moat #3)."),
    ("5. Edge sovereignty + edge MCP (NEW)",
     "MindSet runs MCP server AT THE EDGE — AI agents query the factory floor directly without raw data "
     "leaving the customer network. Cognite has MCP but cloud-only (data ships to their cloud — sovereignty broken). "
     "MaestroHub CEO claimed MCP support in a podcast (unverified, not on public docs) — likely cloud-side if real. "
     "UMH has no MCP. Combined with 'no hyperscaler edition' through 2029, this is the strongest EU-regulatory moat. "
     "Defense, public sector, regulated pharma cannot use the others."),
]

row = 5
for (name, body) in moats:
    bg = XF_PURPLE if row % 2 == 1 else XF_BLUE
    s3.write(row, 2, name, XF_BOLD_LG)
    s3.write(row, 3, body, bg)
    s3.set_row_height(row, 90)
    row += 1


# ── Sheet 4 — Three Deployment Editions ─────────────────────────────────────
s4 = Sheet("4. 3 Editions (INV)")
s4.freeze(1)
s4.set_col_width(1, 24)
s4.set_col_width(2, 38)
s4.set_col_width(3, 38)
s4.set_col_width(4, 38)
s4.set_col_width(5, 36)

s4.write(1, 1, "Aspect", XF_HDR)
s4.write(1, 2, "ON-PREMISE", XF_HDR)
s4.write(1, 3, "HYBRID  (default)", XF_HDRGLD)
s4.write(1, 4, "SELF-HOSTED", XF_HDR)
s4.write(1, 5, "Hyperscaler  (NOT OFFERED — reconsider 2029)", XF_HDR)
s4.set_row_height(1, 28)

editions = [
    ("Target customer",
     "Defense, public sector, nuclear, sensitive pharma — air-gap mandatory",
     "Commercial ETI — default offer",
     "Large multi-site with existing EU cloud relationship or private K8s",
     "Customers asking for AWS / Azure / GCP — refused by design"),
    ("Cloud component",
     "NONE",
     "Scaleway FR / OVH FR (MindSet-managed)",
     "Hetzner / IONOS / T-Systems / 3DS Outscale / customer's on-prem K8s",
     "—"),
    ("Multi-site aggregation",
     "No (each site is independent)",
     "Yes",
     "Yes",
     "—"),
    ("Remote dashboard",
     "No (factory LAN only)",
     "Yes (app.mindsetdata.io)",
     "Yes (hosted on customer cloud)",
     "—"),
    ("KG snapshot backup",
     "Local NAS / customer responsibility",
     "Encrypted to MindSet cloud",
     "Encrypted to customer cloud",
     "—"),
    ("MCP server",
     "Edge only",
     "Edge + optional cloud relay (opt-in)",
     "Edge + optional cloud relay (opt-in)",
     "—"),
    ("AI agents",
     "Local Phi-3 only",
     "Local Phi-3 default + optional remote LLM (with disclosure)",
     "Local Phi-3 default + optional remote LLM (with disclosure)",
     "—"),
    ("RGPD / NIS2",
     "Perfect — zero external dependency",
     "Compliant — FR jurisdiction, EU clouds",
     "Compliant — customer choice of EU jurisdiction",
     "INCOMPATIBLE — US CLOUD Act exposure"),
    ("Pricing posture (TBD)",
     "Premium (specialized support, audit)",
     "Standard (the everyday offer)",
     "Standard + customer pays infra",
     "—"),
]

row = 2
for row_data in editions:
    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if row % 2 == 0 else XF_BOLD_LG
    s4.write(row, 1, row_data[0], bg_b)
    s4.write(row, 2, row_data[1], XF_BLUE)
    s4.write(row, 3, row_data[2], XF_GREEN)   # default = green
    s4.write(row, 4, row_data[3], XF_BLUE)
    s4.write(row, 5, row_data[4], XF_RED)     # not offered = red
    s4.set_row_height(row, 60)
    row += 1


# ── Sheet 5 — Locked Decisions (INTERNAL) ───────────────────────────────────
s5 = Sheet("5. Locked Decisions (INT)")
s5.freeze(1)
s5.set_col_width(1, 32)
s5.set_col_width(2, 50)
s5.set_col_width(3, 50)
s5.set_col_width(4, 18)

s5.write(1, 1, "Decision", XF_HDR)
s5.write(1, 2, "Rationale", XF_HDR)
s5.write(1, 3, "Alternatives rejected", XF_HDR)
s5.write(1, 4, "Date locked", XF_HDR)
s5.set_row_height(1, 24)

decisions = [
    # --- Corrections to prior decisions (newest, top) ---
    ("Licensing: PROPRIETARY closed-source for first 2 years (REVERSES Apache 2.0)",
     "Early-stage commercial protection. Closed source prevents fast-follower forks during PMF. Open-core / source-available reconsidered in 2028.",
     "Apache 2.0 from V1 (REVERSED — exposes IP too early); source-available BSL (legal complexity without OSS trust upside); open-core from V1 (premature for 1-engineer team).",
     "2026-06-28"),
    ("Fuzzy Join: OF-state-based attribution (REVERSES sliding window)",
     "Real-world ERP timestamps lag OT by hours (operators update end-of-shift). Time-window joins break. OF-state-based attribution is robust to clock skew.",
     "Sliding-window ±10 min (REVERSED — fails on real ERP data); operator-entered OF assignment (manual burden).",
     "2026-06-28"),
    ("Edition naming: On-Premise / Hybrid / Self-Hosted (REPLACES Air-Gap / Sovereign Cloud / BYOC)",
     "Plain-language names for non-technical Plant Manager / CFO buyers. 'Air-Gap' and 'BYOC' are jargon.",
     "Keep technical names (less accessible); Local/Cloud/Custom (loses sovereignty implication).",
     "2026-06-28"),
    ("No hyperscaler edition through 2029. Reconsider for international scaling (US/APAC) in year 4.",
     "Adding AWS/Azure/GCP in V1-V3 collapses the EU sovereignty moat for defense, public sector, regulated pharma. International expansion = 2029+ topic, then a SEPARATE product line.",
     "Add hyperscalers in V1-V2 (kills sovereignty moat early); never add them (caps TAM at EU permanently).",
     "2026-06-28"),
    # --- V1 scope decisions ---
    ("Platform positioning + 3 starter use-case templates",
     "Platform pitch with 3 ready templates (micro-stop, energy waste, OEE/TRS). Customers + AI agents build more use cases on top.",
     "Single-use-case product (narrows TAM); pure platform without templates (no demo, classic startup death).",
     "2026-06-27"),
    ("AI-native from V1 (not V2)",
     "Phi-3 runtime + edge MCP + Ad-hoc Analyst agent ship in V1. Narrative becomes 'AI-native edge platform', not 'industrial platform with AI later'.",
     "AI in V2 (weak 2026 investor pitch); multiple agents in V1 (over-scope, 5 mediocre vs 1 excellent).",
     "2026-06-27"),
    ("ERP connectors pulled forward to V1 (was V1.5)",
     "Fuzzy Join (Moat #2) needs ERP input. Pulling forward makes the moat demoable from first customer install.",
     "Keep ERP at V1.5 (leaves Fuzzy Join undemoable in V1).",
     "2026-06-27"),
    ("SQL connector V1 dialects: PostgreSQL + MSSQL + MySQL",
     "Covers ~80% of FR ETI ERPs (Sage X3 = MSSQL, Dynamics = MSSQL, Odoo = PostgreSQL, web ERPs = MySQL). Oracle + HANA = V1.5+ based on demand.",
     "PostgreSQL only (too narrow — leaves Sage/Dynamics blocked); all 5 dialects (Oracle + HANA = wrong segment, high effort).",
     "2026-06-27"),
    ("MCP server edge-only in V1, cloud relay = V1.5+",
     "Edge-only simplifies V1 architecture, preserves sovereignty default. Sufficient for V1 profile (Plant Manager inside factory + founder demos with Claude Desktop on factory LAN).",
     "Cloud-only MCP (breaks sovereignty); edge + cloud relay at V1 (doubles V1 surface).",
     "2026-06-27"),
    ("V1 native AI agent: Ad-hoc Analyst (sole agent)",
     "One excellent demoable agent > 5 mediocre. Chat UI in dashboard, Phi-3 default, grounded via MCP. Demos the MCP integration. All other 12 agents = V1.5+/V2.",
     "Discovery Coach first (onboarding-only, weaker demo); Daily Briefing first (needs accumulated data, weak day-1 demo); 3-5 agents in V1 (over-scope).",
     "2026-06-27"),
    ("Tribal Knowledge moat ships in V1 via dropdown + free text",
     "The MOAT is the DATASET (sensor pattern → operator label), not the chatbot UX. Dropdown + free text accumulates the dataset from V1. Phi-3 chatbot = V2 polish.",
     "Defer Tribal Knowledge entirely to V2 (loses Moat #4 at V1); Phi-3 chatbot in V1 (FR conversational quality + jargon too risky).",
     "2026-06-27"),
    # --- Sprint 2 strategic decisions ---
    ("3 deployment editions: Air-Gap / Sovereign Cloud / BYOC",
     "Lets customers self-select by sovereignty needs. Excludes hyperscalers — preserves regulatory moat for public sector + defense.",
     "Single one-size SaaS (excludes air-gap); supporting AWS/Azure (breaks sovereignty story).",
     "2026-06-27"),
    ("Cloud tier scope: aggregation + remote view + backup + heartbeat only",
     "A feature goes to cloud only if it spans multi-site, latency tolerates >1s, and data crossing is already-transformed.",
     "Cloud-side pipelines (latency); cloud-side rules engine (sub-sec impossible); cloud-side SLM (defeats sovereignty).",
     "2026-06-27"),
    ("Alerting: edge-direct + cloud heartbeat monitor (not SMTP relay)",
     "Honest framing. Generic SMTP relay was a fallback for <10% of customers. Heartbeat detects dead edge agent — real value.",
     "Generic cloud SMTP relay (muddies architecture); no cloud alerting (loses dead-agent detection).",
     "2026-06-27"),
    ("MCP server: essential, edge-default, optional cloud relay",
     "MCP becoming de-facto standard. Edge-default preserves sovereignty. Cloud relay opt-in for remote AI access.",
     "REST/GraphQL only (won't plug into Claude/Copilot natively); cloud-only MCP (breaks sovereignty default).",
     "2026-06-27"),
    ("AI provider: local-default + optional remote with disclosure (Option B)",
     "Sovereignty by default. Customer flexibility for existing LLM contracts. UI warns explicitly when remote is enabled.",
     "Strict local-only (limits flexibility); any-LLM-no-warnings (cannot claim sovereignty); EU-LLM-only (too narrow).",
     "2026-06-27"),
    ("BYOC scope: EU-jurisdiction cloud or on-prem K8s only",
     "AWS-EU / Azure-EU are subject to US CLOUD Act. Excluding hyperscalers preserves regulatory moat for highest-value verticals.",
     "Any-cloud BYOC (breaks moat); Scaleway/OVH only (too narrow for customers with existing EU cloud).",
     "2026-06-27"),
    ("Hardware play: software-only",
     "Zero-hardware sales promise. Avoids margin distraction and inventory risk.",
     "Partnered mini-PC reseller (Beelink/Lenovo) — distracts from software focus.",
     "2026-06-27"),
    ("Competitor frame: both (mid-market direct + giants as reference)",
     "Investors need 'where we fit' AND 'who we beat in deals' stories.",
     "Cognite-only (wrong segment); mid-market-only (no investor reference point).",
     "2026-06-27"),
    ("ISA-95 as UNS ontology",
     "Industry standard for hierarchical manufacturing data modeling. Maximum interoperability.",
     "Custom ontology (no interop); flat namespace (loses semantic value).",
     "2026 (existing)"),
    ("Apache 2.0 license  ⚠ REVERSED 2026-06-28 — see Proprietary decision at top of sheet",
     "Was: open, permissive, aligns with OSS-first strategy. Reversed because early-stage commercial protection needed.",
     "AGPL (scares enterprise); MIT (no patent grant); proprietary — now the chosen path.",
     "2026 → REVERSED"),
    ("Deterministic rules engine (not ML) for micro-stop detection",
     "Auditable, predictable, zero training data needed. Plant Manager can adjust thresholds without DS expertise.",
     "ML model (requires labeled historical data customers don't have on day 1).",
     "2026 (existing)"),
    ("Read-only on all source systems",
     "Core security promise to IT/OT teams. Zero risk of corrupting production data. NIS2 compliant.",
     "Write-back to PLC (security risk, regulatory blocker).",
     "2026 (existing)"),
    ("Push-only outbound HTTPS + mTLS",
     "Zero inbound open ports on customer network. Removes attack surface, eliminates firewall friction.",
     "Inbound API (firewall changes required, attack surface).",
     "2026 (existing)"),
    ("No third-party middleware (no Kepware)",
     "Direct sales argument to security teams. Removes per-tag licensing cost. Zero-dev deployment.",
     "Kepware (per-tag cost, vendor lock-in, sales friction).",
     "2026 (existing)"),
]

row = 2
for (dec, rat, alt, date) in decisions:
    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if row % 2 == 0 else XF_BOLD_LG
    s5.write(row, 1, dec, bg_b)
    s5.write(row, 2, rat, bg)
    s5.write(row, 3, alt, bg)
    s5.write(row, 4, date, XF_AMBER if "2026-06" in date else XF_BLUE)
    s5.set_row_height(row, 70)
    row += 1


# ── Sheet 6 — Open Questions (INTERNAL) ─────────────────────────────────────
s6 = Sheet("6. Open Questions (INT)")
s6.freeze(1)
s6.set_col_width(1, 28)
s6.set_col_width(2, 36)
s6.set_col_width(3, 56)
s6.set_col_width(4, 22)

s6.write(1, 1, "Topic", XF_HDR)
s6.write(1, 2, "Question still open", XF_HDR)
s6.write(1, 3, "Options on the table", XF_HDR)
s6.write(1, 4, "Decision owner / when", XF_HDR)
s6.set_row_height(1, 24)

open_qs = [
    # Currently most important — these block first customers / first investor pitch
    ("Pricing model",
     "How do we monetise?",
     "A) Open-core (free Edge Agent + paid cloud + support)  B) Per-site SaaS  C) Freemium + paid support  D) Tiered (community / pro / enterprise like UMH's 36k€/site/year)",
     "Cécilia (CEO) — before investor deck"),
    ("4th and 5th starter templates (post-V1)",
     "V1 ships 3 templates (micro-stop / energy / OEE). First customers will reveal what to build next.",
     "A) Quality (defect detection)  B) Changeover (setup-time analysis)  C) Predictive maintenance  D) Schedule deviation  E) Wait for customer signal, don't pre-decide",
     "Cécilia + first 5 customers — Q4 2026"),
    ("International expansion: add hyperscaler edition in 2029?",
     "Locked through 2029: no hyperscaler edition. Reconsider for US/APAC expansion when EU footprint is established.",
     "A) Add 'Global Edition' on AWS/Azure/GCP as separate product line  B) Stay EU-only (cap TAM)  C) Partner with EU-equivalent in other regions (sovereignty federation)",
     "Cécilia + Mohamed — 2029 strategic review"),
    ("Second hire engineering profile",
     "When seed funding closes, what's the first engineer profile to hire?",
     "A) Full-stack Go + React (clone Mohamed) — speed  B) ML / data eng (boost AI agents) — depth  C) DevOps / cloud platform (BYOC + multi-tenant) — scale  D) OT/industrial-protocols specialist — coverage",
     "Cécilia + Mohamed — at seed close"),
    ("Multi-tenant SaaS",
     "Do we offer a hosted SaaS option for SMBs that can't run Docker?",
     "A) Never (single-tenant always)  B) Hosted SaaS for SMB tier  C) Reseller model via integrators",
     "Long-term — after first 20 customers"),
    ("Historian integration depth",
     "How deep do we go on PI / Wonderware / InfluxDB push?",
     "A) Push enriched events only (current plan)  B) Full bi-directional with historian queries  C) Replace historian for small customers",
     "Mohamed — V1.5 design"),
    ("V2 native AI agents — which 3 next after Ad-hoc Analyst",
     "V1 ships Ad-hoc Analyst only. Which V1.5/V2 agents next?",
     "Top candidates: Daily Briefing (Plant Manager value), Tribal Knowledge Chatbot (moat polish), Causality Reasoner (technical depth)",
     "Mohamed + first customers — Q1 2027"),
    ("Sales motion: self-serve vs integrator-led",
     "Roadmap says self-serve Docker pull. UMH + MaestroHub both go via integrators.",
     "A) Self-serve only  B) Self-serve + curated integrator partners  C) Pivot to integrator-led if self-serve doesn't convert",
     "Cécilia — after first 10 customers"),
    ("Geographic expansion order",
     "FR-first locked. Then what?",
     "A) DACH next — manufacturing density  B) Italy + Spain — agrifood density  C) Nordics — pharma + sovereignty culture",
     "Cécilia — Year 2 plan"),
    ("Industry vertical focus",
     "Roadmap names manufacturing broadly. Niche down first?",
     "A) Agrifood (cost-€ + energy resonate)  B) Pharma (compliance pitch)  C) General manufacturing (broader TAM, harder to differentiate)",
     "Cécilia — Year 1 GTM"),
]

row = 2
for (topic, q, opts, owner) in open_qs:
    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if row % 2 == 0 else XF_BOLD_LG
    s6.write(row, 1, topic, bg_b)
    s6.write(row, 2, q, XF_AMBER)
    s6.write(row, 3, opts, bg)
    s6.write(row, 4, owner, XF_BLUE)
    s6.set_row_height(row, 75)
    row += 1


# ── Sheet 7 — AI Agent Catalog (INTERNAL) ──────────────────────────────────
s7 = Sheet("7. AI Agent Catalog (INT)")
s7.freeze(1)
s7.set_col_width(1, 6)
s7.set_col_width(2, 22)
s7.set_col_width(3, 26)
s7.set_col_width(4, 46)
s7.set_col_width(5, 22)
s7.set_col_width(6, 12)
s7.set_col_width(7, 10)

s7.write(1, 1, "ID", XF_HDR)
s7.write(1, 2, "Pillar", XF_HDR)
s7.write(1, 3, "Agent name", XF_HDR)
s7.write(1, 4, "What it does", XF_HDR)
s7.write(1, 5, "Primary persona", XF_HDR)
s7.write(1, 6, "Priority", XF_HDR)
s7.write(1, 7, "Version", XF_HDR)
s7.set_row_height(1, 24)

agents = [
    # === V1 — ONLY ONE AGENT SHIPS HERE ===
    ("V2", "Visualise", "Ad-hoc Analyst  ★ V1 SOLE AGENT",
     "Free-text Q&A in dashboard chat. 'How did Line 2 perform this week?' 'Which product had the most micro-stops?' Grounded via MCP, cites KG sources.",
     "Plant Manager / CFO / Ops Director", "P0", "V1"),
    # === V1.5 — POST-FIRST-CUSTOMER ===
    ("V1", "Visualise", "Daily Briefing Agent",
     "Shift start: 'Last 24h: 47 events, 312€ impact. Top cause: jam on Line 1 (62%). Action: check sensor calibration.'",
     "Plant Manager", "P0", "V1.5"),
    ("A1", "Act", "Alert Triage Agent",
     "When €-threshold breached: pings Plant Manager with cause + recommended action + 1-click acknowledge.",
     "Plant Manager", "P0", "V1.5"),
    ("C1", "Connect", "Discovery Coach",
     "'I scanned your network. Found 3 OPC-UA servers, 1 Modbus device. Want me to walk through what each tag likely is?'",
     "IT/OT Manager", "P0", "V1.5"),
    ("C2", "Connect", "Tag Classifier",
     "Explains Phi-3 classification confidence. Asks user to confirm low-confidence tags. (Note: SLM classification itself ships V1 — this is the agentic explanation layer.)",
     "IT/OT Manager", "P0", "V1.5"),
    # === V2 — MOAT DEPTH ===
    ("X1", "Contextualise", "Tribal Knowledge Chatbot  (moat polish — V1 ships dropdown)",
     "Phi-3 conversational interview with operator after each stop. Extracts richer cause + resolution. NOTE: dropdown + free text capture ships V1 — chatbot is V2 UX polish, not the moat itself.",
     "Operator + Plant Manager", "P0", "V2"),
    ("X2", "Contextualise", "Causality Reasoner",
     "When event has no obvious cause: queries related tags + recent events. 'Pressure dropped 12s before — could be a leak.'",
     "Plant Manager", "P1", "V2"),
    ("V4", "Visualise", "Trend Spotter",
     "Proactively surfaces patterns: 'Events on Line 3 doubled in 3 days. Same cause as last month's incident.'",
     "Plant Manager", "P1", "V2"),
    ("V3", "Visualise", "Multi-site Benchmarker",
     "Site A vs Site B on TRS, stop frequency, cause distribution.",
     "Ops Director / CEO", "P1", "V2 (multi-site)"),
    ("X3", "Contextualise", "Cost Coach",
     "'Your hourly cost is 85€/h but for product X with margin 0.08€/unit, your true cost-per-stop-minute is 4.30€. Want to refine?'",
     "CFO + Plant Manager", "P1", "V2"),
    # === V3+ — STRETCH ===
    ("C3", "Connect", "Connector Recommender",
     "'You mentioned SAP. I see your ERP exposes REST. Want me to configure the connector?'",
     "IT/OT Manager", "P2", "V3"),
    ("A2", "Act", "Maintenance Scheduler",
     "'Sensor S3 has 17 false alarms this week. Recommend recalibration. Want to draft a maintenance ticket?'",
     "Maintenance team", "P2", "V3"),
    ("A3", "Act", "Compliance Reporter",
     "Generates NIS2 / RGPD audit reports on demand.",
     "IT/OT Manager + CISO", "P2", "V3"),
]

row = 2
for (aid, pillar, name, desc, persona, prio, ver) in agents:
    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY
    pillar_color = {"Connect": XF_BLUE, "Contextualise": XF_PURPLE, "Visualise": XF_GREEN, "Act": XF_AMBER}[pillar]
    s7.write(row, 1, aid, XF_BOLD_C)
    s7.write(row, 2, pillar, pillar_color)
    s7.write(row, 3, name, XF_BOLD_LG if row % 2 else XF_BOLD)
    s7.write(row, 4, desc, bg)
    s7.write(row, 5, persona, bg)
    prio_xf = XF_GREEN if prio == "P0" else (XF_AMBER if prio == "P1" else XF_LGREY)
    s7.write(row, 6, prio, prio_xf)
    s7.write(row, 7, ver, XF_NORMAL_C)
    s7.set_row_height(row, 60)
    row += 1


# ── Sheet 8 — Edge vs Cloud Map (INTERNAL) — FULL 63-COMPONENT INVENTORY ───
s8 = Sheet("8. Edge vs Cloud (INT)")
s8.freeze(2)
s8.set_col_width(1, 22)   # Category
s8.set_col_width(2, 8)    # ID
s8.set_col_width(3, 38)   # Component
s8.set_col_width(4, 12)   # Runs on
s8.set_col_width(5, 14)   # Status
s8.set_col_width(6, 50)   # Notes

s8.write(1, 1, "Full Edge + Cloud Inventory — 63 components total (~51 ship in V1)", XF_TITLE)
s8.merge(1, 1, 1, 6)
s8.set_row_height(1, 30)

s8.write(2, 1, "Category", XF_HDR)
s8.write(2, 2, "ID", XF_HDR)
s8.write(2, 3, "Component", XF_HDR)
s8.write(2, 4, "Runs on", XF_HDR)
s8.write(2, 5, "Status", XF_HDR)
s8.write(2, 6, "Notes / why this location", XF_HDR)
s8.set_row_height(2, 24)

# (category, id, name, location, status, notes)
# location: EDGE / CLOUD / OPT-IN / MARKETING
# status: Built / V1 / V1.5 / V2 / V3 / Pending
inventory = [
    # === STORAGE ===
    ("1. Storage", "S1", "SQLite ring buffer (raw events 7-15d, auto-purge)", "EDGE", "V1", "Pure-Go modernc/sqlite, zero deps. Raw history stays where generated."),
    ("1. Storage", "S2", "Domain KG (cumulative site fingerprint)", "EDGE", "V1 (partial)", "Equipment / Event / Cause / Cost nodes. Grows forever. THE MOAT #3 DATASET."),
    ("1. Storage", "S3", "Technical KG (pipeline topology, in-memory)", "EDGE", "Built", "5-min cache, busted by registry hash. Rebuilt from YAML pipelines."),
    ("1. Storage", "S4", "Tag registry (persisted)", "EDGE", "Built", "Discovered OPC-UA tags + values + types. Survives restart."),
    ("1. Storage", "S5", "Topic registry (persisted)", "EDGE", "Built", "Live MQTT topics + msg/s + category."),
    ("1. Storage", "S6", "State tracker (in-memory)", "EDGE", "Built", "Current Run/Stop/Setup states per work center."),
    # === MESSAGE BUS ===
    ("2. Message bus", "M1", "Local MQTT broker (Mosquitto)", "EDGE", "V1", "SEPARATE process. tcp://localhost:1883. Internal nervous system. Bundled in docker-compose."),
    # === DISCOVERY ===
    ("3. Discovery", "D1", "Network scanner (subnet scan for OPC-UA / Modbus / S7 / MQTT)", "EDGE", "V1", "Scans customer OT subnet. Cannot be done from outside firewall."),
    ("3. Discovery", "D2", "OPC-UA browse engine (node tree + continuation points)", "EDGE", "Built", "Direct PLC/SCADA access. Read-only on source."),
    ("3. Discovery", "D3", "Modbus device fingerprint DB (20-30 common devices)", "EDGE", "V1", "Schneider, Siemens, Danfoss, SEW. Auto-loads register map per IP/MAC."),
    ("3. Discovery", "D4", "Behavioral inference engine (10-15 min pattern matching)", "EDGE", "V1", "Auto-classifies opaque Modbus/S7 registers. Needs raw tag stream — cloud forbidden."),
    ("3. Discovery", "D5", "UNS ISA-95 mapper (tag → Site/Area/WC/Tag)", "EDGE", "V1 (partial)", "Transforms raw tags into UNS topics. Output goes up, never raw input."),
    # === CONNECTORS ===
    ("4. Connectors", "C1", "OPC-UA", "EDGE", "Built", "gopcua. Mature, secure modes need hardening."),
    ("4. Connectors", "C2", "Modbus TCP", "EDGE", "V1", "goburrow/modbus."),
    ("4. Connectors", "C3", "SQL multi-dialect (PostgreSQL + MSSQL + MySQL)", "EDGE", "V1 (NEW)", "Covers ~80% of FR ETI ERPs. Oracle + HANA = V1.5+."),
    ("4. Connectors", "C4", "Siemens S7", "EDGE", "V1.5", "gos7. 30-40% of European industrial park."),
    ("4. Connectors", "C5", "REST (modern ERPs)", "EDGE", "V1.5", "SAP S/4HANA, D365, Sage X3."),
    ("4. Connectors", "C6", "Files / FTP / SFTP (CSV/Excel/JSON)", "EDGE", "V1.5", "Unblocks ETIs without API."),
    ("4. Connectors", "C7", "MQTT generic", "EDGE", "V2", "Modern IIoT gateways."),
    ("4. Connectors", "C8", "Sparkplug B", "EDGE", "V2", "MQTT with ISA-95 structured payload."),
    ("4. Connectors", "C9", "MTConnect", "EDGE", "V2", "CNC / machining / metallurgy."),
    ("4. Connectors", "C10", "BACnet/IP", "EDGE", "V2", "Building / HVAC."),
    ("4. Connectors", "C11", "Omron FINS / MongoDB / InfluxDB", "EDGE", "V2/V3", "Niche connectors, demand-driven."),
    # === PROCESSING ===
    ("5. Processing", "P1", "Pipeline engine (topological YAML execution)", "EDGE", "Built", "recover()-protected. Core platform mechanism."),
    ("5. Processing", "P2", "Function registry (connectors / transforms / outputs)", "EDGE", "Built", "Pluggable function catalog."),
    ("5. Processing", "P3", "Rules engine (deterministic threshold-based)", "EDGE", "V1 (partial)", "Micro-stop, energy, OEE/TRS templates."),
    ("5. Processing", "P4", "OF-state Fuzzy Join engine (the moat)", "EDGE", "V1", "NOT sliding window. Robust to multi-hour ERP clock skew."),
    ("5. Processing", "P5", "Cost model in € (3-field wizard + ERP auto V1.5)", "EDGE", "V1", "Multiplies local events × rates. Trivial compute."),
    ("5. Processing", "P6", "OEE / TRS calculator (real availability vs declared)", "EDGE", "V1", "The killer demo: declared-vs-real OEE gap in € per week."),
    # === KG INTEGRATION ===
    ("6. KG integration", "K1", "KG subscriber (enriches Domain KG from MQTT events)", "EDGE", "Built", "Listens to mindset/events/* topics."),
    ("6. KG integration", "K2", "KG builder (computes Technical KG from pipeline registry)", "EDGE", "Built", "5-min cache."),
    ("6. KG integration", "K3", "KG REST API (GET /api/kg/domain + /api/kg/technical)", "EDGE", "Built", "Feeds dashboard + MCP server."),
    # === LOCAL UI ===
    ("7. Local UI", "U1", "React app skeleton (Vite + Tailwind + Zustand)", "EDGE", "Built", "localhost:8080. SPA."),
    ("7. Local UI", "U2", "Pipeline Studio (React Flow drag-drop canvas)", "EDGE", "Built", "Visual pipeline builder."),
    ("7. Local UI", "U3", "KG viewer (Cytoscape)", "EDGE", "Built", "Technical + domain graphs."),
    ("7. Local UI", "U4", "Dashboard skeleton + WebSocket live hub", "EDGE", "Built (partial)", "Real-time push to React."),
    ("7. Local UI", "U5", "Live Gantt timeline (Run/Stop/Setup)", "EDGE", "V1", "Per-machine timeline view."),
    ("7. Local UI", "U6", "Pareto chart (causes by €)", "EDGE", "V1", "Visual shock — money lost where."),
    ("7. Local UI", "U7", "OEE / TRS view (declared vs real gap)", "EDGE", "V1", "THE KILLER DEMO."),
    ("7. Local UI", "U8", "ROI simulator", "EDGE", "V1", "Gain potential if cause #1 resolved."),
    ("7. Local UI", "U9", "Tribal knowledge capture UI (dropdown + free text)", "EDGE", "V1", "MOAT #4 CAPTURE. The dataset ships V1."),
    ("7. Local UI", "U10", "Onboarding wizard (3-field cost + endpoints)", "EDGE", "V1", "Plant Manager-facing setup flow."),
    # === AI LAYER ===
    ("8. AI layer", "A1", "Phi-3 runtime (via Ollama, local process)", "EDGE", "V1", "Default LLM. Sovereign. ~2.5GB RAM."),
    ("8. AI layer", "A2", "MCP server (edge-default)", "EDGE", "V1", "Exposes KG tools to AI clients. localhost:5000."),
    ("8. AI layer", "A3", "Ad-hoc Analyst agent (chat UI in dashboard)", "EDGE", "V1", "V1 SOLE NATIVE AGENT. Grounded via MCP."),
    ("8. AI layer", "A4", "Remote LLM proxy (OpenAI/Claude/Mistral config + warning)", "EDGE", "V1", "Customer opt-in. UI warns when data leaves network."),
    # === COMMUNICATION ===
    ("9. Communication", "O1", "WebSocket live hub (push to local dashboard)", "EDGE", "Built", "Real-time UI updates."),
    ("9. Communication", "O2", "HTTPS pusher to cloud (mTLS + offline queue)", "EDGE", "V1", "Hybrid/Self-Hosted only. Push-only outbound."),
    ("9. Communication", "O3", "SMTP / Slack / Teams alerting (direct from edge)", "EDGE", "V1", "No cloud relay needed for the common case."),
    ("9. Communication", "O4", "Heartbeat sender (to cloud)", "EDGE", "V1", "Hybrid/Self-Hosted only. Cloud detects dead edge."),
    ("9. Communication", "O5", "Historian push (PI / Wonderware / InfluxDB)", "EDGE", "V1.5", "Push enriched events to customer's existing historian."),
    # === INFRASTRUCTURE ===
    ("10. Infrastructure", "I1", "Config loader (YAML)", "EDGE", "Built", "Loads agent.yaml + pipelines/*.yaml."),
    ("10. Infrastructure", "I2", "Structured logger", "EDGE", "Built", "JSON logs to stdout."),
    ("10. Infrastructure", "I3", "Local secrets management (SOPS)", "EDGE", "V1", "ERP credentials, LLM API keys, cloud auth keys."),
    ("10. Infrastructure", "I4", "License key validator (proprietary gate)", "EDGE", "V1", "Validates license against cloud, gracefully degrades if offline."),
    ("10. Infrastructure", "I5", "Health check endpoints (/api/health)", "EDGE", "Built", "Used by heartbeat + customer monitoring."),
    ("10. Infrastructure", "I6", "Auto-update mechanism (signed pulls)", "EDGE", "V1.5", "Opt-in update channel."),
    # === SECURITY (pending Entry 20 decision) ===
    ("11. Security", "SEC1", "SQLite encryption at rest (SQLCipher)", "EDGE", "Pending", "Pending security framework decision."),
    ("11. Security", "SEC2", "Signed binaries (cosign + Sigstore)", "EDGE", "Pending", "Customer can verify binary authenticity."),
    ("11. Security", "SEC3", "SBOM (CycloneDX) shipped with releases", "EDGE", "Pending", "Customer CISO can audit dependencies."),
    ("11. Security", "SEC4", "Audit log (immutable + SIEM-exportable)", "EDGE", "Pending", "NIS2 requirement."),
    ("11. Security", "SEC5", "RBAC engine (Admin / Operator / Read-only / Auditor)", "EDGE", "Pending V1.5", "Multi-user per site."),
    ("11. Security", "SEC6", "SSO (SAML / OIDC)", "EDGE", "Pending V1.5", "Customer identity provider integration."),
    # === CLOUD COMPONENTS ===
    ("12. Cloud (Hybrid + Self-Hosted only)", "CL1", "Cross-site KG aggregator", "CLOUD", "V1.5", "Multi-site only. Sees only transformed events, never raw."),
    ("12. Cloud (Hybrid + Self-Hosted only)", "CL2", "Multi-site dashboard", "CLOUD", "V1.5", "Cross-site Pareto + benchmark views for CEO/Ops."),
    ("12. Cloud (Hybrid + Self-Hosted only)", "CL3", "Remote single-site dashboard", "CLOUD", "V1", "Same React UI served from cloud for off-site users."),
    ("12. Cloud (Hybrid + Self-Hosted only)", "CL4", "Site management API (auth + keys + entitlements)", "CLOUD", "V1", "Central registry of sites + customer admin panel."),
    ("12. Cloud (Hybrid + Self-Hosted only)", "CL5", "KG snapshot backup (encrypted)", "CLOUD", "V1", "DR. Snapshots encrypted at edge — cloud holds opaque blobs."),
    ("12. Cloud (Hybrid + Self-Hosted only)", "CL6", "Heartbeat / liveness monitor", "CLOUD", "V1", "Alerts if edge agent stops reporting >5min."),
    # === OPT-IN CLOUD ===
    ("13. Opt-in cloud", "OP1", "Remote LLM proxy (cloud-side routing)", "OPT-IN", "V1", "Routes to OpenAI/Claude/Mistral. UI warns explicitly."),
    ("13. Opt-in cloud", "OP2", "Cloud MCP relay (remote AI without VPN)", "OPT-IN", "V1.5", "For CEO/CFO querying from laptop on home network."),
    # === MARKETING ===
    ("14. Marketing", "W1", "Public marketing site (mindsetdata.io)", "MARKETING", "V1", "Next.js / Vercel. No customer data."),
    ("14. Marketing", "W2", "Edge Agent image distribution (private registry)", "MARKETING", "V1", "Closed-source license → private registry, not public Docker Hub."),
]

# Render with category section headers
last_cat = None
row = 3
for (cat, id_, name, where, status, notes) in inventory:
    # Category section break
    if cat != last_cat:
        s8.write(row, 1, cat, XF_HDRGLD)
        s8.merge(row, 1, row, 6)
        s8.set_row_height(row, 22)
        row += 1
        last_cat = cat

    # Location style
    if where == "EDGE":
        loc_xf = XF_GREEN
    elif where == "CLOUD":
        loc_xf = XF_BLUE
    elif where == "OPT-IN":
        loc_xf = XF_AMBER
    else:  # MARKETING
        loc_xf = XF_LGREY

    # Status style
    if status == "Built":
        status_xf = XF_GREEN
    elif status.startswith("V1"):
        status_xf = XF_PURPLE
    elif status.startswith("V1.5"):
        status_xf = XF_BLUE
    elif status.startswith("V2"):
        status_xf = XF_AMBER
    elif status.startswith("V3") or status.startswith("V2/V3"):
        status_xf = XF_LGREY
    elif status.startswith("Pending"):
        status_xf = XF_RED
    else:
        status_xf = XF_NORMAL

    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY

    s8.write(row, 1, "", bg)  # Category column blank in detail rows
    s8.write(row, 2, id_, XF_BOLD_C)
    s8.write(row, 3, name, XF_BOLD_LG if row % 2 else XF_BOLD)
    s8.write(row, 4, where, loc_xf)
    s8.write(row, 5, status, status_xf)
    s8.write(row, 6, notes, bg)
    s8.set_row_height(row, 38)
    row += 1


# ── Sheet 9 — Technical Differentiation (NEW — for engineering DD) ─────────
s9 = Sheet("9. Technical Diff (DD)")
s9.freeze(2)
s9.set_col_width(1, 28)
s9.set_col_width(2, 36)
s9.set_col_width(3, 32)
s9.set_col_width(4, 32)
s9.set_col_width(5, 32)

s9.write(1, 1, "Technical Differentiation — engineer-to-engineer comparison (for technical due diligence)", XF_TITLE)
s9.merge(1, 1, 1, 5)
s9.set_row_height(1, 30)

s9.write(2, 1, "Technical dimension", XF_HDR)
s9.write(2, 2, "MindSet Data", XF_HDRGLD)
s9.write(2, 3, "UMH", XF_HDR)
s9.write(2, 4, "MaestroHub", XF_HDR)
s9.write(2, 5, "Cognite", XF_HDR)
s9.set_row_height(2, 28)

tech_diff = [
    ("Edge runtime footprint",
     "1 Go binary, ~200MB RAM idle, +2.5GB when Phi-3 loaded. Single process.",
     "Kubernetes cluster — 6+ containers (Benthos, Node-RED, HiveMQ, Redpanda, TimescaleDB, Grafana). 4-8GB RAM min per node.",
     "Single host: 8 CPU / 16GB RAM / 200GB SSD per vendor spec.",
     "Thin extractor (small) — heavy compute lives in Cognite cloud."),
    ("Deployment unit",
     "docker run  (one command)",
     "helm install on existing K8s cluster (or create one)",
     "Vendor-installed",
     "Vendor-installed extractor + Cognite cloud provisioning"),
    ("Time to first event (cold start)",
     "<60 seconds (Docker pull + auto-discovery + first OPC-UA tag value)",
     "5-15 minutes (K8s init + multi-container startup + connector config)",
     "Days (vendor scheduling)",
     "Days-to-weeks (cloud tenant provisioning)"),
    ("Storage at edge",
     "SQLite ring buffer 7-15 days (configurable, TTL auto-purge). Pure-Go modernc/sqlite — zero deps.",
     "TimescaleDB + Kafka log retention (Redpanda). Persistent volumes required.",
     "Not publicly documented.",
     "Minimal — events pushed to cloud immediately."),
    ("Real-time processing model",
     "Go in-process; rules engine runs at <500ms latency on detected state transitions.",
     "Multi-service via Kafka topics; pub-sub latency adds ms-to-s.",
     "Edge-to-cloud — much processing happens cloud-side.",
     "Cloud-side (extractor → cloud transformation pipelines)."),
    ("OT/IT join algorithm",
     "OF-state-based attribution: poll ERP for active OFs, tag every OT event with current OF. Robust to multi-hour clock skew.",
     "User-built — typically via Node-RED flows.",
     "Not documented as a dedicated feature.",
     "Entity contextualization (P&ID OCR, asset matching) — different problem."),
    ("LLM runtime location",
     "Local Phi-3 via Ollama (default). Optional remote LLM (any) with explicit data-disclosure warning.",
     "None native — user integrates own AI separately.",
     "None native — 'AI-ready' framing only.",
     "Cognite Atlas AI runs in Cognite cloud."),
    ("MCP server location",
     "EDGE (in-process with cmd/server). Customer's AI agents connect from inside the factory network.",
     "None.",
     "CEO podcast claim — unverified, likely cloud-side if real.",
     "CLOUD (via Function Apps endpoint). Data ships to Cognite cloud."),
    ("Failure mode if cloud unreachable",
     "Full local operation continues. Rules, dashboard, alerts, AI agent all work. Events queue for later sync.",
     "Same — fully OSS local stack continues to function.",
     "Likely degraded (cloud-side processing path).",
     "Significant degradation — AI, dashboards, contextualization live in cloud."),
    ("Customer audit surface (code review)",
     "1 Go binary delivered as compiled artifact (proprietary closed-source for 2 years). Customer cannot read source.",
     "Multiple OSS projects on GitHub — fully auditable. Apache 2.0.",
     "Not available — proprietary.",
     "Not available — proprietary."),
    ("Per-tag licensing fees",
     "NONE.",
     "NONE (OSS community edition; enterprise = per-factory).",
     "Not publicly disclosed.",
     "Enterprise contract — typically includes tag count tiers."),
    ("Third-party middleware (Kepware-style)",
     "NO — native protocol drivers.",
     "NO — native via Benthos / Node-RED.",
     "Unknown — vendor stack.",
     "NO — native Cognite extractors."),
]

row = 3
for (dim, mds, umh, mh, cog) in tech_diff:
    bg = XF_NORMAL if row % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if row % 2 == 0 else XF_BOLD_LG
    s9.write(row, 1, dim, bg_b)
    s9.write(row, 2, mds, XF_PURPLE)
    s9.write(row, 3, umh, bg)
    s9.write(row, 4, mh, bg)
    s9.write(row, 5, cog, bg)
    s9.set_row_height(row, 75)
    row += 1


# ════════════════════════════════════════════════════════════════════════════
# ASSEMBLE
# ════════════════════════════════════════════════════════════════════════════
sheets = [s1, s2, s3, s4, s5, s6, s7, s8, s9]

workbook_xml = (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
    '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" '
    'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">'
    '<bookViews><workbookView xWindow="0" yWindow="0" windowWidth="16000" windowHeight="9000"/></bookViews>'
    '<sheets>'
    + "".join(f'<sheet name="{escape(sh.name)}" sheetId="{i}" r:id="rId{i}"/>'
              for i, sh in enumerate(sheets, 1))
    + '</sheets></workbook>'
).encode("utf-8")

wb_rels = (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
    '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">'
    + "".join(
        f'<Relationship Id="rId{i}" '
        f'Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" '
        f'Target="worksheets/sheet{i}.xml"/>'
        for i in range(1, len(sheets)+1)
    )
    + f'<Relationship Id="rId{len(sheets)+1}" '
    f'Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings" '
    f'Target="sharedStrings.xml"/>'
    f'<Relationship Id="rId{len(sheets)+2}" '
    f'Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" '
    f'Target="styles.xml"/>'
    '</Relationships>'
).encode("utf-8")

pkg_rels = (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
    '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">'
    '<Relationship Id="rId1" '
    'Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" '
    'Target="xl/workbook.xml"/>'
    '</Relationships>'
).encode("utf-8")

content_types = (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
    '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">'
    '<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>'
    '<Default Extension="xml" ContentType="application/xml"/>'
    '<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>'
    '<Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>'
    '<Override PartName="/xl/sharedStrings.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"/>'
    + "".join(
        f'<Override PartName="/xl/worksheets/sheet{i}.xml" '
        f'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>'
        for i in range(1, len(sheets)+1)
    )
    + '</Types>'
).encode("utf-8")

# Pre-render sheet XML so all si() calls populate _strings first
sheet_xmls = [sh.to_xml() for sh in sheets]

buf = io.BytesIO()
with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
    zf.writestr("[Content_Types].xml", content_types)
    zf.writestr("_rels/.rels", pkg_rels)
    zf.writestr("xl/workbook.xml", workbook_xml)
    zf.writestr("xl/_rels/workbook.xml.rels", wb_rels)
    zf.writestr("xl/styles.xml", styles_xml())
    zf.writestr("xl/sharedStrings.xml", sst_xml())
    for i, xml in enumerate(sheet_xmls, 1):
        zf.writestr(f"xl/worksheets/sheet{i}.xml", xml)

with open(OUT, "wb") as f:
    f.write(buf.getvalue())

kb = len(buf.getvalue()) // 1024
print(f"OK -- {OUT}  ({kb} KB, {len(_strings)} strings, {len(_xfs)} cell formats, {len(sheets)} sheets)")
