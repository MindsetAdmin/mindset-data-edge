// Mapping + styling for the two knowledge graphs served by the API:
//  - technical (/api/kg/technical): node.name, edge.from/to/type
//  - domain    (/api/kg/domain):    node.label, edge.from_id/to_id/relation

export const NODE_COLORS = {
  // technical graph
  connection: '#0ea5e9',
  topic: '#34d399',
  function: '#fbbf24',
  pipeline: '#f87171',
  dashboard: '#a78bfa',
  // domain graph
  Equipment: '#3b82f6',
  Event: '#f59e0b',
  Cause: '#8b5cf6',
  Cost: '#10b981',
};

export const FALLBACK_COLOR = '#64748b';

// Normalize either graph shape into Cytoscape elements.
export function toElements(kind, graph) {
  if (!graph) return [];
  const nodes = (graph.nodes || []).map((n) => ({
    data: {
      id: n.id,
      label: kind === 'technical' ? n.name : n.label,
      type: n.type,
      properties: n.properties || {},
    },
  }));
  // Defensive: Cytoscape throws if an edge references a missing node, so drop
  // any edge whose endpoints aren't present.
  const ids = new Set(nodes.map((n) => n.data.id));
  const edges = (graph.edges || [])
    .map((e) => ({
      data: {
        id: e.id,
        source: kind === 'technical' ? e.from : e.from_id,
        target: kind === 'technical' ? e.to : e.to_id,
        label: kind === 'technical' ? e.type : e.relation,
      },
    }))
    .filter((e) => ids.has(e.data.source) && ids.has(e.data.target));
  return [...nodes, ...edges];
}

// Distinct node types present, for the legend.
export function typesPresent(graph) {
  const set = new Set((graph?.nodes || []).map((n) => n.type));
  return [...set];
}

export const STYLESHEET = [
  {
    selector: 'node',
    style: {
      'background-color': (ele) => NODE_COLORS[ele.data('type')] || FALLBACK_COLOR,
      label: 'data(label)',
      color: '#e2e8f0',
      'font-size': '9px',
      'text-valign': 'top',
      'text-halign': 'center',
      'text-margin-y': '-4px',
      'text-wrap': 'wrap',
      'text-max-width': '90px',
      width: 30,
      height: 30,
      'border-width': 2,
      'border-color': '#0f172a',
    },
  },
  {
    selector: 'edge',
    style: {
      width: 1.4,
      'line-color': '#475569',
      'target-arrow-color': '#475569',
      'target-arrow-shape': 'triangle',
      'curve-style': 'bezier',
      label: 'data(label)',
      'font-size': '7px',
      color: '#94a3b8',
      'text-rotation': 'autorotate',
      'text-background-color': '#0f172a',
      'text-background-opacity': 0.85,
      'text-background-padding': '2px',
    },
  },
  {
    selector: 'node:selected',
    style: { 'border-color': '#60a5fa', 'border-width': 3 },
  },
];

export const LAYOUTS = {
  technical: { name: 'breadthfirst', directed: true, padding: 30, spacingFactor: 1.4 },
  domain: { name: 'concentric', padding: 30, minNodeSpacing: 40 },
};
