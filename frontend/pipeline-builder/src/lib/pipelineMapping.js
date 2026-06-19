// Conversion between the ReactFlow graph (nodes + edges) and the backend
// Pipeline shape (internal/pipeline/types.go). The trigger node has the fixed
// id "trigger" so `depends_on: ["trigger"]` matches the existing YAML convention.

import { getCategory } from './functionMeta';

export const TRIGGER_ID = 'trigger';

export function slugify(s) {
  return (s || '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '_')
    .replace(/^_+|_+$/g, '');
}

// flowToPipeline: canvas -> backend JSON ready for POST /api/pipelines
export function flowToPipeline({ id, name, description, version, nodes, edges }) {
  const triggerNode = nodes.find((n) => n.type === 'triggerNode');
  const stepNodes = nodes.filter((n) => n.type !== 'triggerNode');

  // depends_on derived from incoming edges
  const deps = {};
  edges.forEach((e) => {
    if (!deps[e.target]) deps[e.target] = [];
    deps[e.target].push(e.source);
  });

  const pipelineNodes = stepNodes.map((n) => ({
    id: n.id,
    name: n.data.name || n.id,
    type: n.data.type,
    function: n.data.function,
    config: n.data.config || {},
    depends_on: deps[n.id] || [],
  }));

  // output = a step node with no outgoing edge (terminal)
  const sources = new Set(edges.map((e) => e.source));
  const terminal = stepNodes.find((n) => !sources.has(n.id));

  const trigger = triggerNode
    ? {
        type: triggerNode.data.triggerType || 'mqtt',
        function: triggerNode.data.function || 'mqtt_subscribe',
        config: triggerNode.data.config || {},
      }
    : { type: 'mqtt', function: 'mqtt_subscribe', config: {} };

  return {
    id,
    name,
    description: description || '',
    version: version || '1.0',
    trigger,
    nodes: pipelineNodes,
    output: terminal ? terminal.id : '',
  };
}

// pipelineToFlow: backend pipeline -> canvas nodes/edges, with a simple
// dependency-depth layout (columns left-to-right).
export function pipelineToFlow(p) {
  const nodes = [];
  const edges = [];

  nodes.push({
    id: TRIGGER_ID,
    type: 'triggerNode',
    position: { x: 40, y: 200 },
    deletable: false,
    data: {
      triggerType: p.trigger?.type || 'mqtt',
      function: p.trigger?.function || 'mqtt_subscribe',
      config: p.trigger?.config || {},
    },
  });

  const steps = p.nodes || [];
  const byId = {};
  steps.forEach((n) => (byId[n.id] = n));

  // depth = longest dependency chain from the trigger
  const depthCache = {};
  const depthOf = (id, seen = new Set()) => {
    if (id === TRIGGER_ID) return 0;
    if (depthCache[id] != null) return depthCache[id];
    if (seen.has(id)) return 1; // cycle guard
    seen.add(id);
    const node = byId[id];
    const parents = (node?.depends_on || []).filter((d) => d !== TRIGGER_ID);
    const d = parents.length ? 1 + Math.max(...parents.map((pp) => depthOf(pp, seen))) : 1;
    depthCache[id] = d;
    return d;
  };

  const rowByCol = {};
  steps.forEach((n) => {
    const col = depthOf(n.id);
    const row = rowByCol[col] || 0;
    rowByCol[col] = row + 1;
    nodes.push({
      id: n.id,
      type: 'pipelineNode',
      position: { x: 40 + col * 240, y: 60 + row * 120 },
      data: {
        name: n.name || n.id,
        type: n.type,
        function: n.function,
        config: n.config || {},
        category: getCategory(n.type),
      },
    });
    (n.depends_on || []).forEach((dep) => {
      edges.push({ id: `e-${dep}-${n.id}`, source: dep, target: n.id, animated: true });
    });
  });

  return { nodes, edges };
}
