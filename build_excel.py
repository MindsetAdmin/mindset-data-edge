"""
Build MindSet vs Cognite comparison Excel workbook using only Python built-ins.
xlsx = ZIP of XML files — no third-party libraries needed.
"""
import zipfile, io
from xml.sax.saxutils import escape

OUT = r"C:\Users\khena\Desktop\MindSet_Architecture_Comparison.xlsx"

# ── colour palette ───────────────────────────────────────────────────────────
C_NAVY  = "FF1F3864"
C_WHITE = "FFFFFFFF"
C_LGREY = "FFF2F2F2"
C_GREEN = "FFE2EFDA"
C_AMBER = "FFFFF2CC"
C_RED   = "FFFFE0E0"
C_BLUE  = "FFD6E4F7"
C_GOLD  = "FFD4A017"
C_DKRED = "FFC00000"

# ── shared-string table ──────────────────────────────────────────────────────
_strings: list = []

def si(s: str) -> int:
    s = str(s)
    if s not in _strings:
        _strings.append(s)
    return _strings.index(s)

# ── style registry ───────────────────────────────────────────────────────────
# fills[0]=none, fills[1]=gray125 are mandatory by xlsx spec
_fills = ["none", "gray125"]
_fonts = []
_xfs   = []

def reg_fill(argb: str) -> int:
    if argb not in _fills:
        _fills.append(argb)
    return _fills.index(argb)

def reg_font(bold=False, color="FF000000", sz=11) -> int:
    f = (bold, color, sz)
    if f not in _fonts:
        _fonts.append(f)
    return _fonts.index(f)

def reg_xf(font_id=0, fill_id=0, wrap=True, halign=None, valign="top") -> int:
    x = (font_id, fill_id, wrap, halign, valign)
    if x not in _xfs:
        _xfs.append(x)
    return _xfs.index(x)

# pre-register fonts
FN = reg_font(False, "FF000000", 11)   # normal black
FB = reg_font(True,  "FF000000", 11)   # bold black
FH = reg_font(True,  C_WHITE,    11)   # white (header)
FR = reg_font(True,  C_DKRED,    11)   # dark red

# pre-register fills
FI_WHITE = reg_fill(C_WHITE)
FI_LGREY = reg_fill(C_LGREY)
FI_NAVY  = reg_fill(C_NAVY)
FI_GREEN = reg_fill(C_GREEN)
FI_AMBER = reg_fill(C_AMBER)
FI_RED   = reg_fill(C_RED)
FI_BLUE  = reg_fill(C_BLUE)
FI_GOLD  = reg_fill(C_GOLD)

# pre-register cell formats
XF_NORMAL  = reg_xf(FN, FI_WHITE)
XF_BOLD    = reg_xf(FB, FI_WHITE)
XF_HDR     = reg_xf(FH, FI_NAVY,  halign="center")
XF_LGREY   = reg_xf(FN, FI_LGREY)
XF_GREEN   = reg_xf(FN, FI_GREEN)
XF_AMBER   = reg_xf(FB, FI_AMBER)
XF_RED     = reg_xf(FR, FI_RED)
XF_BLUE    = reg_xf(FN, FI_BLUE)
XF_BOLD_G  = reg_xf(FB, FI_GREEN)
XF_BOLD_R  = reg_xf(FR, FI_RED)
XF_BOLD_LG = reg_xf(FB, FI_LGREY)

# ── cell address helpers ─────────────────────────────────────────────────────
def col_letter(n: int) -> str:
    s = ""
    while n:
        n, r = divmod(n - 1, 26)
        s = chr(65 + r) + s
    return s

def addr(r: int, c: int) -> str:
    return f"{col_letter(c)}{r}"

# ── XML builders ─────────────────────────────────────────────────────────────
def styles_xml() -> bytes:
    # fonts
    fonts_xml = ""
    for (bold, color, sz) in _fonts:
        b = "<b/>" if bold else ""
        fonts_xml += f'<font>{b}<sz val="{sz}"/><color rgb="{color}"/><name val="Calibri"/><family val="2"/></font>'

    # fills — solid fills need fgColor + bgColor indexed="64"
    fills_xml = ""
    for f in _fills:
        if f == "none":
            fills_xml += '<fill><patternFill patternType="none"/></fill>'
        elif f == "gray125":
            fills_xml += '<fill><patternFill patternType="gray125"/></fill>'
        else:
            fills_xml += (
                f'<fill><patternFill patternType="solid">'
                f'<fgColor rgb="{f}"/><bgColor indexed="64"/>'
                f'</patternFill></fill>'
            )

    # border (one default)
    borders_xml = '<border><left/><right/><top/><bottom/><diagonal/></border>'

    # cell xfs
    xfs_xml = ""
    for (font_id, fill_id, wrap, halign, valign) in _xfs:
        al_attrs = []
        if wrap:   al_attrs.append('wrapText="1"')
        if halign: al_attrs.append(f'horizontal="{halign}"')
        if valign: al_attrs.append(f'vertical="{valign}"')
        al = f'<alignment {" ".join(al_attrs)}/>' if al_attrs else ""
        xfs_xml += (
            f'<xf numFmtId="0" fontId="{font_id}" fillId="{fill_id}" '
            f'borderId="0" xfId="0" applyFont="1" applyFill="1" applyAlignment="1">'
            f'{al}</xf>'
        )

    xml = (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">'
        f'<fonts count="{len(_fonts)}">{fonts_xml}</fonts>'
        f'<fills count="{len(_fills)}">{fills_xml}</fills>'
        f'<borders count="1">{borders_xml}</borders>'
        '<cellStyleXfs count="1">'
        '<xf numFmtId="0" fontId="0" fillId="0" borderId="0"/>'
        '</cellStyleXfs>'
        f'<cellXfs count="{len(_xfs)}">{xfs_xml}</cellXfs>'
        '<cellStyles count="1">'
        '<cellStyle name="Normal" xfId="0" builtinId="0"/>'
        '</cellStyles>'
        '</styleSheet>'
    )
    return xml.encode("utf-8")

def sst_xml() -> bytes:
    items = "".join(
        f'<si><t xml:space="preserve">{escape(s)}</t></si>'
        for s in _strings
    )
    xml = (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        f'<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" '
        f'count="{len(_strings)}" uniqueCount="{len(_strings)}">'
        f'{items}</sst>'
    )
    return xml.encode("utf-8")

# ── Sheet class ──────────────────────────────────────────────────────────────
class Sheet:
    def __init__(self, name: str):
        self.name = name
        self._cells: dict = {}          # (r,c) -> (xf_id, str_value)
        self._col_widths: dict = {}     # col -> float
        self._freeze_row: int = 0

    def write(self, r: int, c: int, value, xf_id: int = XF_NORMAL):
        self._cells[(r, c)] = (xf_id, str(value))

    def set_col_width(self, c: int, w: float):
        self._col_widths[c] = w

    def freeze(self, row: int):
        self._freeze_row = row

    def to_xml(self) -> bytes:
        # sheetViews — always emit a default view; add freeze pane if set
        if self._freeze_row:
            top = addr(self._freeze_row + 1, 1)
            sv = (
                '<sheetViews><sheetView tabSelected="0" workbookViewId="0">'
                f'<pane ySplit="{self._freeze_row}" topLeftCell="{top}" '
                f'activePane="bottomLeft" state="frozen"/>'
                f'<selection pane="bottomLeft" activeCell="{top}" sqref="{top}"/>'
                '</sheetView></sheetViews>'
            )
        else:
            sv = '<sheetViews><sheetView workbookViewId="0"/></sheetViews>'

        # column widths
        col_xml = ""
        for c, w in sorted(self._col_widths.items()):
            col_xml += f'<col min="{c}" max="{c}" width="{w}" customWidth="1"/>'
        if col_xml:
            col_xml = f"<cols>{col_xml}</cols>"

        # rows
        by_row: dict = {}
        for (r, c), (xf_id, val) in self._cells.items():
            by_row.setdefault(r, []).append((c, xf_id, val))

        rows_xml = ""
        for r in sorted(by_row):
            cells_xml = ""
            for c, xf_id, val in sorted(by_row[r]):
                idx = si(val)
                cells_xml += f'<c r="{addr(r,c)}" s="{xf_id}" t="s"><v>{idx}</v></c>'
            rows_xml += f'<row r="{r}">{cells_xml}</row>'

        xml = (
            '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
            '<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">'
            f'{sv}{col_xml}'
            f'<sheetData>{rows_xml}</sheetData>'
            '</worksheet>'
        )
        return xml.encode("utf-8")


# ════════════════════════════════════════════════════════════════════════════
# DATA
# ════════════════════════════════════════════════════════════════════════════

# ── Sheet 1: Current Architecture ───────────────────────────────────────────
s1 = Sheet("1. Current Architecture")
s1.freeze(1)
s1.set_col_width(1, 42); s1.set_col_width(2, 26); s1.set_col_width(3, 34); s1.set_col_width(4, 48)

for ci, h in enumerate(["Component","Where it runs TODAY","Storage / location","Notes"], 1):
    s1.write(1, ci, h, XF_HDR)

s1_data = [
    ("cmd/server  (API, WebSocket, OPC-UA manager)",
     "On-premises (edge server)", "—",
     "HTTP :8080  -  owns OPC-UA session in UI-driven mode"),
    ("cmd/agent  (MQTT, rules engine, UNS contextualizer, KG subscriber)",
     "On-premises (edge server)", "—",
     "Publishes to local MQTT broker  -  runs rules & KG enrichment"),
    ("Knowledge Graph - domain KG  (micro-stops, events, costs, causes)",
     "On-premises", "data/mindset.db  (SQLite, local disk)",
     "Auto-enriched by KGSubscriber from mindset/events/micro-stop"),
    ("Knowledge Graph - technical KG  (pipeline topology / architecture view)",
     "On-premises", "In-memory  (cached 5 min, busted by registry hash)",
     "Rebuilt from pipeline YAML registry  -  empty until you save a pipeline"),
    ("Pipeline engine",
     "On-premises", "config/pipelines/*.yaml  (local files)",
     "Topological execution, YAML-defined, recover() protects server from panics"),
    ("MQTT broker",
     "On-premises", "localhost:1883",
     "Mosquitto or compatible  -  zero cloud dependency"),
    ("OPC-UA session",
     "On-premises", "—",
     "Connects to local / plant OPC-UA server  -  UI-driven by default"),
    ("Dashboards / React UI",
     "On-premises", "Served from :8080 / :5173 (dev)",
     "WebSocket live push from LiveHub  -  Zustand state management"),
    ("MCP server",
     "NOT YET IMPLEMENTED", "TBD - see Sheet 3",
     "Planned key differentiator vs Cognite  -  options: edge or cloud"),
]

for ri, row in enumerate(s1_data, 2):
    bg = XF_NORMAL if ri % 2 == 0 else XF_LGREY
    for ci, val in enumerate(row, 1):
        xfid = XF_RED if "NOT YET" in val else bg
        s1.write(ri, ci, val, xfid)


# ── Sheet 2: MindSet vs Cognite ─────────────────────────────────────────────
s2 = Sheet("2. MindSet vs Cognite")
s2.freeze(1)
s2.set_col_width(1, 30); s2.set_col_width(2, 50); s2.set_col_width(3, 50); s2.set_col_width(4, 16)

for ci, h in enumerate(["Dimension","Cognite CDF","MindSet Data Edge","Advantage"], 1):
    s2.write(1, ci, h, XF_HDR)

s2_data = [
    ("On-premises footprint",
     "Thin extractors only (OPC-UA, historians). Everything else is cloud.",
     "Full stack runs 100% on-premises - air-gap capable today.",
     "MindSet"),
    ("Cloud dependency",
     "MANDATORY - KG, pipelines, AI, dashboards all live in Cognite cloud.",
     "OPTIONAL - works without any cloud connection.",
     "MindSet"),
    ("Knowledge Graph location",
     "Cognite cloud (proprietary Industrial Data Model / CDF data model).",
     "Local SQLite today; can replicate to any cloud DB.",
     "MindSet"),
    ("Contextualization / UNS",
     "Cloud-side transformation pipelines inside Cognite tenant.",
     "Edge-side UNS contextualizer (cmd/agent) - ISA-95 enrichment on the wire.",
     "MindSet"),
    ("Pipeline / transformation engine",
     "Cognite Functions (cloud FaaS, Python). Cognite Transformations for raw->clean.",
     "Local Go engine, YAML-defined, edge-native, no cloud round-trip.",
     "MindSet"),
    ("AI agents",
     "AI Atlas - proprietary, locked to Cognite cloud ecosystem.",
     "AI-agnostic: OpenAI, Azure AI, Bedrock, local Ollama - customer's choice.",
     "MindSet"),
    ("MCP support",
     "Not native. Open REST/GraphQL API + AI Atlas SDK (closed ecosystem).",
     "Native MCP planned - any MCP-compatible AI client can query the KG.",
     "MindSet"),
    ("Vendor lock-in",
     "HIGH - data model, AI layer, and storage are all Cognite-proprietary.",
     "LOW - open standards: MQTT, OPC-UA, ISA-95, SQLite, YAML.",
     "MindSet"),
    ("Dashboards",
     "Cognite cloud (Charts, InField, Industrial Canvas).",
     "Local React UI - can be deployed to any cloud or served from edge.",
     "MindSet"),
    ("Target market",
     "Large enterprises - oil & gas, energy, utilities. High ACV.",
     "SMEs + mid-market, any industry. Lower barrier to entry.",
     "Neutral"),
    ("Ecosystem maturity",
     "Mature SaaS product, proven at scale (Aker BP, Cognite Hub).",
     "Early-stage - higher agility, lower maturity.",
     "Cognite"),
    ("Sovereign / air-gap deployment",
     "Not possible - Cognite cloud is mandatory.",
     "Fully supported today (Scenario A).",
     "MindSet"),
    ("Pricing model",
     "Enterprise contracts, per-asset / per-user licensing. High TCO.",
     "TBD - edge-first model has no mandatory cloud spend.",
     "MindSet"),
    ("Open ecosystem",
     "Partially open (REST API) but AI layer is proprietary.",
     "Fully open - bring your own cloud, AI model, broker.",
     "MindSet"),
]

for ri, (dim, cog, mds, adv) in enumerate(s2_data, 2):
    bg = XF_NORMAL if ri % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if ri % 2 == 0 else XF_BOLD_LG
    s2.write(ri, 1, dim, bg_b)
    s2.write(ri, 2, cog, bg)
    s2.write(ri, 3, mds, bg)
    if adv == "MindSet":
        s2.write(ri, 4, "MindSet", XF_BOLD_G)
    elif adv == "Cognite":
        s2.write(ri, 4, "Cognite", XF_BOLD_R)
    else:
        s2.write(ri, 4, "Neutral", XF_AMBER)


# ── Sheet 3: Deployment Scenarios ───────────────────────────────────────────
s3 = Sheet("3. Deployment Scenarios")
s3.freeze(1)
s3.set_col_width(1, 28); s3.set_col_width(2, 44)
for c in range(3, 9): s3.set_col_width(c, 26)
s3.set_col_width(9, 30)

hdrs3 = ["Scenario","Description","KG location","MCP server",
         "AI agents","Pipelines","Dashboards","Cloud cost","Status"]
for ci, h in enumerate(hdrs3, 1):
    s3.write(1, ci, h, XF_HDR)

s3_data = [
    ("A - Fully on-premises  (CURRENT)",
     "Everything runs on customer hardware. Zero cloud dependency.",
     "Local SQLite",
     "Edge (planned)",
     "Local model or none",
     "Edge execution, local YAML",
     "Local React UI",
     "0 EUR cloud",
     "Implemented today"),
    ("B - Hybrid edge + cloud sync",
     "Edge processes OT data; KG/events replicated to cloud DB.",
     "Edge primary + cloud replica (Postgres / CosmosDB / DynamoDB)",
     "Edge OR cloud",
     "Cloud (OpenAI / Azure AI / Bedrock / Ollama)",
     "Edge execution, cloud copy for audit",
     "Cloud-hosted (Vercel / Azure Static / S3)",
     "Low  (storage + AI API calls)",
     "To be built"),
    ("C - Edge agent + customer cloud",
     "Customer deploys MindSet stack on their own Azure / AWS / GCP tenant.",
     "Customer VPC  (any DB)",
     "Customer cloud",
     "Customer's AI infra",
     "Customer cloud",
     "Customer cloud",
     "Customer pays infra",
     "Needs containerisation (Docker)"),
    ("D - MindSet SaaS  (multi-tenant)",
     "MindSet hosts cloud tier; customers get isolated tenant + edge agent.",
     "MindSet cloud  (per-tenant DB)",
     "MindSet cloud",
     "MindSet-managed AI",
     "MindSet cloud",
     "MindSet cloud",
     "Highest  (multi-tenant infra + AI)",
     "Long-term roadmap"),
]

scenario_xfs = [XF_GREEN, XF_BLUE, XF_LGREY, XF_AMBER]
for ri, row in enumerate(s3_data, 2):
    xfid = scenario_xfs[ri - 2]
    for ci, val in enumerate(row, 1):
        s3.write(ri, ci, val, xfid)


# ── Sheet 4: Certain vs TBD ─────────────────────────────────────────────────
s4 = Sheet("4. Certain vs TBD")
s4.freeze(1)
s4.set_col_width(1, 26); s4.set_col_width(2, 46); s4.set_col_width(3, 38); s4.set_col_width(4, 56)

for ci, h in enumerate(["Topic","What we know for certain","What is TBD","Options / next decision"], 1):
    s4.write(1, ci, h, XF_HDR)

s4_data = [
    ("KG storage",
     "SQLite at data/mindset.db on the edge server. Works today, air-gap safe.",
     "Cloud sync strategy - if/when we push KG data to the cloud.",
     "A) Keep edge-only  B) Replicate to Postgres/CosmosDB  C) Publish KG events to cloud bus"),
    ("MCP server placement",
     "Not yet implemented - zero MCP code exists in the codebase.",
     "Where does the MCP server run?",
     "A) Edge (same process as cmd/server) - offline capable  B) Cloud - accessible to remote AI agents  C) Both (edge + cloud relay)"),
    ("Cloud tier",
     "None exists today - 100% on-premises.",
     "Do we build a cloud tier? Which cloud? Who hosts it?",
     "A) None (pure on-prem product)  B) Customer cloud (Scenario C)  C) MindSet SaaS (Scenario D)"),
    ("AI integration",
     "No AI calls today. Architecture is AI-agnostic by design.",
     "Which AI providers do we target first? How do agents auth to MCP?",
     "A) Local Ollama (private, no API cost)  B) OpenAI / Anthropic API at edge  C) Cloud agents calling edge MCP over HTTPS"),
    ("Multi-site aggregation",
     "Single edge instance per site - no cross-site aggregation today.",
     "How do we aggregate data from multiple factory sites?",
     "A) Each site independent  B) Cloud aggregation layer  C) Site-to-site MQTT federation"),
    ("Containerisation",
     "Runs as raw Go binaries (server.exe, agent.exe) today.",
     "Docker / Kubernetes packaging needed for cloud scenarios.",
     "Needed for Scenarios C and D  -  straightforward given pure-Go + SQLite stack"),
    ("Pipeline engine location",
     "100% on-premises, YAML-defined, Go execution.",
     "For SaaS / cloud: cloud-hosted pipeline execution?",
     "Edge execution preferred (lower latency, works offline). Cloud copy for audit is optional."),
    ("Dashboard hosting",
     "React UI served from edge server (:8080) or Vite dev server (:5173).",
     "For remote access / SaaS: cloud-hosted UI build.",
     "Easy to deploy as static build (npm run build) to any CDN or cloud storage."),
    ("Pricing / licensing",
     "No pricing model defined yet.",
     "How do we monetise? Per site? Per asset? SaaS subscription?",
     "Edge-first = zero mandatory cloud spend -> strong SME pitch. SaaS adds recurring revenue."),
]

for ri, row in enumerate(s4_data, 2):
    bg = XF_NORMAL if ri % 2 == 0 else XF_LGREY
    bg_b = XF_BOLD if ri % 2 == 0 else XF_BOLD_LG
    s4.write(ri, 1, row[0], bg_b)
    s4.write(ri, 2, row[1], XF_GREEN)
    s4.write(ri, 3, row[2], XF_AMBER)
    s4.write(ri, 4, row[3], bg)


# ════════════════════════════════════════════════════════════════════════════
# ASSEMBLE XLSX
# ════════════════════════════════════════════════════════════════════════════
sheets = [s1, s2, s3, s4]

# Call sst_xml LAST so all si() calls have populated _strings first
# (styles_xml reads _fonts/_fills/_xfs which are fully built already)

workbook_xml = (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
    '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" '
    'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">'
    '<bookViews><workbookView xWindow="0" yWindow="0" windowWidth="16000" windowHeight="9000"/></bookViews>'
    '<sheets>'
    + "".join(
        f'<sheet name="{escape(sh.name)}" sheetId="{i}" r:id="rId{i}"/>'
        for i, sh in enumerate(sheets, 1)
    )
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
    f'</Relationships>'
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
    '<Override PartName="/xl/workbook.xml" '
    'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>'
    '<Override PartName="/xl/styles.xml" '
    'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>'
    '<Override PartName="/xl/sharedStrings.xml" '
    'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"/>'
    + "".join(
        f'<Override PartName="/xl/worksheets/sheet{i}.xml" '
        f'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>'
        for i in range(1, len(sheets)+1)
    )
    + '</Types>'
).encode("utf-8")

# Build sheet XMLs first so all shared strings are registered
sheet_xmls = [sh.to_xml() for sh in sheets]

buf = io.BytesIO()
with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
    zf.writestr("[Content_Types].xml", content_types)
    zf.writestr("_rels/.rels", pkg_rels)
    zf.writestr("xl/workbook.xml", workbook_xml)
    zf.writestr("xl/_rels/workbook.xml.rels", wb_rels)
    zf.writestr("xl/styles.xml", styles_xml())
    zf.writestr("xl/sharedStrings.xml", sst_xml())
    for i, (sh, xml) in enumerate(zip(sheets, sheet_xmls), 1):
        zf.writestr(f"xl/worksheets/sheet{i}.xml", xml)

with open(OUT, "wb") as f:
    f.write(buf.getvalue())

kb = len(buf.getvalue()) // 1024
print(f"OK -- written to {OUT}  ({kb} KB, {len(_strings)} shared strings, {len(_xfs)} cell formats)")
