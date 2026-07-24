import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import ForceGraph2D from 'react-force-graph-2d';

// Obsidian-style knowledge graph viewer for MindSet.
// 2026-07-02 — added to replace CytoscapeGraph in KnowledgeGraphPage.
//
// Design choices:
//   • Physics-alive force simulation (d3-force under the hood)
//   • Nodes as pure circles, size proportional to sqrt(degree)
//   • Curved links, subtle color, no arrows unless directional matters
//   • Hover: highlight neighbors + fade the rest — the Obsidian move
//   • Labels only on hover (or on zoom-in past a threshold)
//   • Colors derived from MindSet design tokens (see docs/frontend_redesign.md §3.1)

// Node colors keyed by type. Business categories get warm tones; platform gets
// cool tones. Falls back to a muted grey.
const NODE_COLORS = {
    // ─── Business (site fingerprint) ─────────────────────────
    // Site/Area/WorkCenter/Equipment/Tag form the structural-bootstrap hierarchy
    // (Entries 95-98/107) — colored as a deliberate gradient, broadest to most
    // granular, so the nesting reads visually. Tag is the most numerous type by
    // far (one per signal), so it's kept quiet/muted rather than eye-catching.
    Site:       '#93C5FD',  // light blue — broadest scope
    Area:       '#5EEAD4',  // teal — a zone within the site
    WorkCenter: '#FDBA74',  // light orange — a grouping level above Equipment (e.g. a line)
    Equipment:  '#F87171',  // status-stopped soft red — the physical machine
    Tag:        '#CBD5E1',  // muted slate — the leaf level, deliberately quiet (added Entry 111 — was previously missing, fell through to the grey fallback along with Site/Area)
    // SchemaMapping — the IT-side structural bootstrap (Entry 115, Track B):
    // an auto-suggested canonical mapping of a SQL table, pending validation
    // same as OT nodes above. Distinct hue (violet-blue) since it's neither
    // OT hierarchy nor a business event — added deliberately up front, not
    // left to fall through to FALLBACK_COLOR (the exact bug Entry 111 fixed).
    SchemaMapping: '#0EA5E9',
    Event:     '#E5A445',   // brand accent (amber)
    Cause:     '#A78BFA',   // purple
    Cost:      '#4ADE80',   // status-running green
    Operator:  '#FBBF24',   // status-warn amber
    OF:        '#FB923C',   // orange (production order)
    Product:   '#F472B6',   // pink
    Recipe:    '#FCD34D',   // yellow

    // ─── Platform (pipeline topology) ────────────────────────
    pipeline:   '#60A5FA',  // info blue
    topic:      '#22D3EE',  // cyan
    function:   '#94A3B8',  // slate
    connection: '#818CF8',  // indigo
    dashboard:  '#F0ABFC',  // pink-light
};
const FALLBACK_COLOR = '#6E6E7A';
// v0 structural bootstrap (Entry 95/96) — auto-generated nodes awaiting human
// validation render dimmer with a dashed amber ring, so they're visible in the
// graph but clearly not-yet-confirmed.
const PENDING_RING_COLOR = '#E5A445';
const PENDING_ALPHA = 0.55;

// Design tokens (mirrored here because canvas rendering needs concrete strings)
const CANVAS_BG      = '#0A0A0B';
const LINK_COLOR     = 'rgba(232, 232, 237, 0.15)';
const LINK_HIGHLIGHT = 'rgba(232, 232, 237, 0.55)';
const NODE_HIGHLIGHT_RING = '#E8E8ED';
const LABEL_COLOR    = 'rgba(232, 232, 237, 0.9)';

// Normalize either MindSet KG shape into ForceGraph nodes/links.
// The unified /api/kg endpoint returns {nodes, edges} where nodes have
// id/category/type/label and edges have id/category/from_id/to_id/relation.
// We also accept the legacy shapes (technical: node.name, edge.from/to/type)
// so this component works pre- and post-Entry 50 refactor.
export function normalizeGraph(rawGraph) {
    if (!rawGraph) return { nodes: [], links: [] };

    const rawNodes = rawGraph.nodes || [];
    const rawEdges = rawGraph.edges || [];

    const nodes = rawNodes.map((n) => ({
        id: n.id,
        label: n.label || n.name || n.id,
        type: n.type,
        category: n.category || (isPlatformType(n.type) ? 'platform' : 'business'),
        properties: n.properties || {},
        raw: n,
    }));

    const nodeIds = new Set(nodes.map((n) => n.id));

    const links = rawEdges
        .map((e) => ({
            id: e.id,
            source: e.from_id || e.from,
            target: e.to_id || e.to,
            relation: e.relation || e.type,
            category: e.category,
        }))
        // Drop dangling edges — ForceGraph2D throws otherwise.
        .filter((l) => nodeIds.has(l.source) && nodeIds.has(l.target));

    return { nodes, links };
}

function isPlatformType(type) {
    return (
        type === 'pipeline' ||
        type === 'topic' ||
        type === 'function' ||
        type === 'connection' ||
        type === 'dashboard'
    );
}

/**
 * @param {object} props
 * @param {object} props.graph      raw MindSet KG (unified shape)
 * @param {Function} props.onNodeSelect callback(nodeOrNull)
 */
export default function ForceGraph({ graph, onNodeSelect }) {
    const containerRef = useRef(null);
    const fgRef = useRef(null);

    // Track container size (ForceGraph2D needs explicit width/height for canvas)
    const [size, setSize] = useState({ w: 800, h: 600 });
    useEffect(() => {
        const el = containerRef.current;
        if (!el) return undefined;
        const ro = new ResizeObserver(([entry]) => {
            const cr = entry.contentRect;
            setSize({ w: Math.max(100, Math.floor(cr.width)), h: Math.max(100, Math.floor(cr.height)) });
        });
        ro.observe(el);
        return () => ro.disconnect();
    }, []);

    // Normalize + compute degree per node (used for sizing).
    const data = useMemo(() => {
        const norm = normalizeGraph(graph);
        const degree = new Map();
        norm.links.forEach((l) => {
            degree.set(l.source, (degree.get(l.source) || 0) + 1);
            degree.set(l.target, (degree.get(l.target) || 0) + 1);
        });
        norm.nodes.forEach((n) => {
            n.degree = degree.get(n.id) || 0;
        });
        return norm;
    }, [graph]);

    // Neighbor sets for hover highlight.
    const neighbors = useMemo(() => {
        const map = new Map();
        data.nodes.forEach((n) => map.set(n.id, new Set()));
        data.links.forEach((l) => {
            map.get(l.source)?.add(l.target);
            map.get(l.target)?.add(l.source);
        });
        return map;
    }, [data]);

    const [hoverId, setHoverId] = useState(null);
    const [highlightNodes, setHighlightNodes] = useState(new Set());
    const [highlightLinks, setHighlightLinks] = useState(new Set());

    const onNodeHover = useCallback(
        (node) => {
            if (!node) {
                setHoverId(null);
                setHighlightNodes(new Set());
                setHighlightLinks(new Set());
                return;
            }
            const hn = new Set([node.id]);
            neighbors.get(node.id)?.forEach((id) => hn.add(id));
            const hl = new Set(
                data.links
                    .filter((l) => l.source.id === node.id || l.target.id === node.id ||
                                   l.source === node.id || l.target === node.id)
                    .map((l) => l.id)
            );
            setHoverId(node.id);
            setHighlightNodes(hn);
            setHighlightLinks(hl);
        },
        [neighbors, data.links]
    );

    // Autofit once after layout stabilizes.
    const [fitDone, setFitDone] = useState(false);
    useEffect(() => setFitDone(false), [graph]);
    const onEngineStop = useCallback(() => {
        if (fitDone || !fgRef.current) return;
        fgRef.current.zoomToFit(400, 60);
        setFitDone(true);
    }, [fitDone]);

    // Custom node rendering — precise circles + hover ring + hover label.
    const nodeCanvasObject = useCallback(
        (node, ctx, globalScale) => {
            const isHovered = hoverId === node.id;
            const isFaded = hoverId != null && !highlightNodes.has(node.id);
            const baseSize = 3 + Math.sqrt(node.degree || 0) * 1.6; // 3–8ish
            const size = isHovered ? baseSize * 1.35 : baseSize;

            const color = NODE_COLORS[node.type] || FALLBACK_COLOR;
            const isPending = Boolean(node.properties?.pending);
            ctx.globalAlpha = isFaded ? 0.15 : isPending ? PENDING_ALPHA : 1;

            // Filled circle
            ctx.beginPath();
            ctx.arc(node.x, node.y, size, 0, 2 * Math.PI, false);
            ctx.fillStyle = color;
            ctx.fill();

            // Pending ring — dashed amber, so an unvalidated auto-generated
            // node is visibly different even before hovering it.
            if (isPending && !isFaded) {
                ctx.beginPath();
                ctx.setLineDash([2 / globalScale, 2 / globalScale]);
                ctx.arc(node.x, node.y, size + 2, 0, 2 * Math.PI, false);
                ctx.strokeStyle = PENDING_RING_COLOR;
                ctx.lineWidth = 1 / globalScale;
                ctx.stroke();
                ctx.setLineDash([]);
            }

            // Hover ring
            if (isHovered) {
                ctx.beginPath();
                ctx.arc(node.x, node.y, size + 2.5, 0, 2 * Math.PI, false);
                ctx.strokeStyle = NODE_HIGHLIGHT_RING;
                ctx.lineWidth = 1.2 / globalScale;
                ctx.stroke();
            }

            // Label — show only on hover, or when zoomed in enough
            const showLabel = isHovered || globalScale > 2.2;
            if (showLabel) {
                const label = node.label || node.id;
                const fontSize = Math.max(10 / globalScale, 3);
                ctx.font = `500 ${fontSize}px Inter, system-ui, sans-serif`;
                ctx.textAlign = 'center';
                ctx.textBaseline = 'top';
                ctx.fillStyle = LABEL_COLOR;
                ctx.fillText(label, node.x, node.y + size + 3);
            }

            ctx.globalAlpha = 1;
        },
        [hoverId, highlightNodes]
    );

    // Pointer area for click hit-detection = circle around node
    const nodePointerAreaPaint = useCallback((node, color, ctx) => {
        const size = 3 + Math.sqrt(node.degree || 0) * 1.6;
        ctx.beginPath();
        ctx.arc(node.x, node.y, size + 4, 0, 2 * Math.PI, false);
        ctx.fillStyle = color;
        ctx.fill();
    }, []);

    const linkColor = useCallback(
        (link) => {
            if (hoverId != null && highlightLinks.has(link.id)) return LINK_HIGHLIGHT;
            if (hoverId != null) return 'rgba(232,232,237,0.05)';
            return LINK_COLOR;
        },
        [hoverId, highlightLinks]
    );

    return (
        <div ref={containerRef} className="w-full h-full">
            <ForceGraph2D
                ref={fgRef}
                width={size.w}
                height={size.h}
                graphData={data}
                backgroundColor={CANVAS_BG}
                nodeId="id"
                nodeRelSize={4}
                nodeLabel={(n) => n.label}
                nodeCanvasObject={nodeCanvasObject}
                nodePointerAreaPaint={nodePointerAreaPaint}
                linkColor={linkColor}
                linkWidth={1}
                linkCurvature={0.18}
                linkDirectionalArrowLength={0}
                cooldownTicks={200}
                d3AlphaDecay={0.025}
                d3VelocityDecay={0.35}
                onNodeHover={onNodeHover}
                onNodeClick={(node) => onNodeSelect && onNodeSelect(node?.raw || null)}
                onBackgroundClick={() => onNodeSelect && onNodeSelect(null)}
                onEngineStop={onEngineStop}
                enableNodeDrag={true}
                enableZoomInteraction={true}
                enablePanInteraction={true}
            />
        </div>
    );
}

// Distinct node types present, for the legend (color chip mapping).
export function typesPresent(graph) {
    const set = new Set((graph?.nodes || []).map((n) => n.type));
    return [...set];
}

export { NODE_COLORS, FALLBACK_COLOR };
