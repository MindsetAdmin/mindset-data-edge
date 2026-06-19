import { useState, useEffect, useCallback } from 'react';
import ReactFlow, {
  Background,
  Controls,
  MiniMap,
  addEdge,
  useNodesState,
  useEdgesState,
  ReactFlowProvider,
  useReactFlow,
} from 'reactflow';
import 'reactflow/dist/style.css';

import Palette from '../components/Palette';
import NodeConfigPanel from '../components/NodeConfigPanel';
import PipelineNode from '../components/nodes/PipelineNode';
import TriggerNode from '../components/nodes/TriggerNode';
import { fetchFunctions, fetchPipelines, createPipeline, runPipeline } from '../api/client';
import { getCategory } from '../lib/functionMeta';
import { defaultConfigFor, triggerTypeFor } from '../lib/connectorTemplates';
import { flowToPipeline, pipelineToFlow, slugify, TRIGGER_ID } from '../lib/pipelineMapping';
import { useStudioStore } from '../store/studioStore';

const nodeTypes = { pipelineNode: PipelineNode, triggerNode: TriggerNode };

const makeInitialNodes = () => [
  {
    id: TRIGGER_ID,
    type: 'triggerNode',
    position: { x: 40, y: 200 },
    deletable: false,
    data: {
      triggerType: 'mqtt',
      function: 'mqtt_subscribe',
      config: { topic: 'mindset/events/status-change', qos: 1 },
    },
  },
];

function BuilderInner() {
  const [nodes, setNodes, onNodesChange] = useNodesState(makeInitialNodes());
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [functions, setFunctions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pipelines, setPipelines] = useState([]);
  const [selectedId, setSelectedId] = useState(null);
  const [meta, setMeta] = useState({ id: '', name: '', description: '' });
  const [status, setStatus] = useState(null);

  const { screenToFlowPosition } = useReactFlow();

  const pendingConnector = useStudioStore((s) => s.pendingConnector);
  const pipelineToLoad = useStudioStore((s) => s.pipelineToLoad);
  const clearPending = useStudioStore((s) => s.clearPending);

  useEffect(() => {
    loadFunctions();
    refreshPipelines();
  }, []);

  async function loadFunctions() {
    try {
      setLoading(true);
      const data = await fetchFunctions();
      setFunctions(data.functions || []);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function refreshPipelines() {
    try {
      const data = await fetchPipelines();
      setPipelines(data.pipelines || []);
    } catch {
      /* ignore — listing is best-effort */
    }
  }

  const onConnect = useCallback(
    (params) => setEdges((eds) => addEdge({ ...params, animated: true }, eds)),
    [setEdges]
  );

  const onDragOver = useCallback((event) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event) => {
      event.preventDefault();
      const raw = event.dataTransfer.getData('application/reactflow');
      if (!raw) return;
      const fn = JSON.parse(raw);
      const position = screenToFlowPosition({ x: event.clientX, y: event.clientY });
      const id = `${slugify(fn.name)}_${Math.random().toString(36).slice(2, 7)}`;
      // Connector nodes get their default config template pre-filled.
      const config = fn.type === 'connector' ? defaultConfigFor(fn.name) : {};
      setNodes((nds) =>
        nds.concat({
          id,
          type: 'pipelineNode',
          position,
          data: { name: fn.name, type: fn.type, function: fn.name, config, category: getCategory(fn.type) },
        })
      );
    },
    [screenToFlowPosition, setNodes]
  );

  const updateNodeData = useCallback(
    (id, updater) => setNodes((nds) => nds.map((n) => (n.id === id ? { ...n, data: updater(n.data) } : n))),
    [setNodes]
  );

  const selectedNode = nodes.find((n) => n.id === selectedId) || null;
  const connectors = functions.filter((f) => f.type === 'connector');

  function handleNameChange(value) {
    setMeta((m) => ({ ...m, name: value, id: m.id || slugify(value) }));
  }

  async function handleSave() {
    if (!meta.id || !meta.name) {
      setStatus({ type: 'error', msg: 'ID et nom sont requis.' });
      return;
    }
    try {
      const pipeline = flowToPipeline({ ...meta, version: '1.0', nodes, edges });
      await createPipeline(pipeline);
      setStatus({ type: 'ok', msg: `Pipeline « ${meta.id} » sauvegardé.` });
      refreshPipelines();
    } catch (err) {
      setStatus({ type: 'error', msg: err.message });
    }
  }

  function handleLoad(id) {
    if (!id) return;
    const p = pipelines.find((x) => x.id === id);
    if (!p) return;
    const flow = pipelineToFlow(p);
    setNodes(flow.nodes);
    setEdges(flow.edges);
    setMeta({ id: p.id, name: p.name, description: p.description || '' });
    setSelectedId(null);
    setStatus({ type: 'ok', msg: `Pipeline « ${p.id} » chargé.` });
  }

  function handleNew() {
    setNodes(makeInitialNodes());
    setEdges([]);
    setMeta({ id: '', name: '', description: '' });
    setSelectedId(null);
    setStatus(null);
  }

  async function handleRun() {
    if (!meta.id) {
      setStatus({ type: 'error', msg: 'Sauvegardez le pipeline avant de l’exécuter.' });
      return;
    }
    setStatus({ type: 'pending', msg: 'Exécution…' });
    try {
      const res = await runPipeline(meta.id);
      const byNode = {};
      (res.nodes || []).forEach((n) => {
        byNode[n.node_id] = n.status;
      });
      setNodes((nds) => nds.map((n) => ({ ...n, data: { ...n.data, runStatus: byNode[n.id] } })));
      const okCount = (res.nodes || []).filter((n) => n.status === 'success').length;
      setStatus({
        type: res.status === 'success' ? 'ok' : 'error',
        msg: `Exécution : ${res.status} — ${okCount}/${(res.nodes || []).length} nœuds OK`,
      });
    } catch (e) {
      setStatus({ type: 'error', msg: e.message });
    }
  }

  // Consume cross-page intents from Connect / Pipelines.
  useEffect(() => {
    if (!pendingConnector) return;
    updateNodeData(TRIGGER_ID, (d) => ({
      ...d,
      function: pendingConnector.name,
      triggerType: triggerTypeFor(pendingConnector.name),
      config: defaultConfigFor(pendingConnector.name),
    }));
    setStatus({ type: 'ok', msg: `Connecteur « ${pendingConnector.name} » appliqué au trigger.` });
    clearPending();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pendingConnector]);

  useEffect(() => {
    if (pipelineToLoad && pipelines.length) {
      handleLoad(pipelineToLoad);
      clearPending();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pipelineToLoad, pipelines]);

  return (
    <div className="flex h-full">
      <Palette functions={functions} loading={loading} error={error} />

      <div className="flex-1 flex flex-col min-w-0">
        {/* Toolbar */}
        <div className="bg-dark-900 border-b border-dark-700 px-4 py-2 flex items-center gap-2 flex-wrap">
          <input
            value={meta.name}
            onChange={(e) => handleNameChange(e.target.value)}
            placeholder="Nom du pipeline"
            className="bg-dark-950 border border-dark-700 rounded-md px-3 py-1.5 text-sm w-48 focus:outline-none focus:border-blue-500"
          />
          <input
            value={meta.id}
            onChange={(e) => setMeta((m) => ({ ...m, id: slugify(e.target.value) }))}
            placeholder="id"
            className="bg-dark-950 border border-dark-700 rounded-md px-3 py-1.5 text-sm w-44 font-mono text-dark-300 focus:outline-none focus:border-blue-500"
          />
          <button
            onClick={handleSave}
            className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-4 py-1.5 rounded-md transition"
          >
            💾 Sauvegarder
          </button>
          <button
            onClick={handleRun}
            className="bg-green-600 hover:bg-green-500 text-white text-sm px-4 py-1.5 rounded-md transition"
          >
            ▶️ Exécuter
          </button>
          <button
            onClick={handleNew}
            className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition"
          >
            ✨ Nouveau
          </button>

          <select
            onChange={(e) => handleLoad(e.target.value)}
            value=""
            className="bg-dark-950 border border-dark-700 rounded-md px-3 py-1.5 text-sm text-dark-300 focus:outline-none focus:border-blue-500"
          >
            <option value="">📂 Charger…</option>
            {pipelines.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>

          {status && (
            <span
              className={`text-sm ml-2 ${
                status.type === 'ok' ? 'text-green-400' : status.type === 'pending' ? 'text-dark-400' : 'text-red-400'
              }`}
            >
              {status.type === 'ok' ? '✅' : status.type === 'pending' ? '⏳' : '❌'} {status.msg}
            </span>
          )}
        </div>

        {/* Canvas */}
        <div className="flex-1 bg-dark-950" onDrop={onDrop} onDragOver={onDragOver}>
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            nodeTypes={nodeTypes}
            onNodeClick={(_, n) => setSelectedId(n.id)}
            onPaneClick={() => setSelectedId(null)}
            fitView
            proOptions={{ hideAttribution: true }}
          >
            <Background color="#1e293b" gap={16} />
            <Controls className="!bg-dark-800 !border-dark-700" />
            <MiniMap pannable zoomable nodeColor={() => '#334155'} maskColor="rgba(2,6,23,0.6)" />
          </ReactFlow>
        </div>
      </div>

      <NodeConfigPanel
        node={selectedNode}
        connectors={connectors}
        onChange={updateNodeData}
        onClose={() => setSelectedId(null)}
      />
    </div>
  );
}

export default function BuilderPage() {
  return (
    <ReactFlowProvider>
      <BuilderInner />
    </ReactFlowProvider>
  );
}
