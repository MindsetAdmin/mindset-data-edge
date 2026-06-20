// Conversion between the ReactFlow graph (nodes + edges) and the backend
// Pipeline shape (internal/pipeline/types.go). The trigger node has the fixed
// id "trigger" so `depends_on: ["trigger"]` matches the existing YAML convention.

import { getCategory } from './functionMeta';

export const TRIGGER_ID = 'trigger';

// Visual bands on the Compose canvas. The pipeline "core" sits in the middle;
// entry (trigger) and exit (outputs) are shown as separate IN/OUT bands.
export const ZONES = {
  in: { x: 0, width: 190, label: 'ENTRÉE', className: 'border-cyan-500/40 text-cyan-300' },
  core: { x: 220, width: 580, label: 'CŒUR', className: 'border-dark-600 text-dark-400' },
  out: { x: 820, width: 200, label: 'SORTIE', className: 'border-red-500/40 text-red-300' },
};
const ZONE_HEIGHT = 520;

// makeZoneNodes returns the three background band nodes.
export function makeZoneNodes() {
  return Object.entries(ZONES).map(([key, z]) => ({
    id: `zone_${key}`,
    type: 'zoneNode',
    position: { x: z.x, y: 0 },
    data: { label: z.label, className: z.className },
    style: { width: z.width, height: ZONE_HEIGHT },
    selectable: false,
    draggable: false,
    deletable: false,
    connectable: false,
    zIndex: -1,
  }));
}

// makeTriggerNode returns the entry node placed in the IN band.
export function makeTriggerNode() {
  return {
    id: TRIGGER_ID,
    type: 'triggerNode',
    position: { x: ZONES.in.x + 20, y: 90 },
    deletable: false,
    data: {
      triggerType: 'mqtt',
      function: 'mqtt_subscribe',
      config: { topic: 'mindset/events/status-change', qos: 1 },
    },
  };
}

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
  // Only pipelineNode entries are real steps (zone bands are visual-only).
  const stepNodes = nodes.filter((n) => n.type === 'pipelineNode');

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

// pipelineToFlow: backend pipeline -> canvas nodes/edges. Entry goes in the IN
// band, output nodes in the OUT band, the rest spread across the CŒUR band.
export function pipelineToFlow(p) {
  const nodes = [...makeZoneNodes()];
  const edges = [];

  nodes.push({
    ...makeTriggerNode(),
    data: {
      triggerType: p.trigger?.type || 'mqtt',
      function: p.trigger?.function || 'mqtt_subscribe',
      config: p.trigger?.config || {},
    },
  });

  const steps = p.nodes || [];
  let coreRow = 0;
  let outRow = 0;

  steps.forEach((n) => {
    const isOutput = n.type === 'output';
    let position;
    if (isOutput) {
      position = { x: ZONES.out.x + 20, y: 90 + outRow * 110 };
      outRow += 1;
    } else {
      position = { x: ZONES.core.x + 30 + (coreRow % 2) * 250, y: 70 + coreRow * 95 };
      coreRow += 1;
    }
    nodes.push({
      id: n.id,
      type: 'pipelineNode',
      position,
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
