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
import PickerModal from '../components/PickerModal';
import PipelineNode from '../components/nodes/PipelineNode';
import TriggerNode from '../components/nodes/TriggerNode';
import ZoneNode from '../components/nodes/ZoneNode';
import {
  fetchFunctions,
  fetchPipelines,
  createPipeline,
  runPipeline,
  fetchTags,
  fetchMachines,
  fetchTopics,
  fetchConfig,
} from '../api/client';
import { getCategory } from '../lib/functionMeta';
import { defaultConfigFor, triggerTypeFor } from '../lib/connectorTemplates';
import { defaultFunctionConfig } from '../lib/functionDefaults';
import {
  flowToPipeline,
  pipelineToFlow,
  slugify,
  TRIGGER_ID,
  makeZoneNodes,
  makeTriggerNode,
} from '../lib/pipelineMapping';
import { useStudioStore } from '../store/studioStore';

const nodeTypes = { pipelineNode: PipelineNode, triggerNode: TriggerNode, zoneNode: ZoneNode };

const makeInitialNodes = () => [...makeZoneNodes(), makeTriggerNode()];

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
  const [showFnPicker, setShowFnPicker] = useState(false);
  const [fieldPickers, setFieldPickers] = useState({ machine_id: [], topic: [], broker: [], node_id: [] });
  const [configDefaults, setConfigDefaults] = useState(null);
  const [machines, setMachines] = useState([]);
  const [dupModal, setDupModal] = useState(null); // { existing } when a duplicate is found

  const { screenToFlowPosition } = useReactFlow();

  const pendingConnector = useStudioStore((s) => s.pendingConnector);
  const pipelineToLoad = useStudioStore((s) => s.pipelineToLoad);
  const clearPending = useStudioStore((s) => s.clearPending);

  useEffect(() => {
    loadFunctions();
    refreshPipelines();
    loadPickerOptions();
  }, []);

  async function loadPickerOptions() {
    let machine_id = [];
    let topic = [];
    let node_id = [];
    let broker = [];
    // Machines (with live status) + their tags grouped by machine.
    try {
      const m = await fetchMachines();
      setMachines((m.machines || []).filter((x) => x.work_center !== '(autres)'));
      machine_id = (m.machines || [])
        .filter((x) => x.work_center !== '(autres)')
        .map((x) => ({
          value: x.work_center,
          label: x.work_center,
          sub: x.state ? (x.state.running ? 'Running ✅' : 'Stopped ❌') : `${x.tags.length} tags`,
        }));
      node_id = (m.machines || []).flatMap((x) =>
        (x.tags || []).map((t) => ({
          value: t.node_id,
          label: t.name || t.node_id,
          sub: `valeur: ${t.value} · ${t.data_type}`,
          badge: t.node_id,
          group: x.work_center,
        }))
      );
    } catch { /* best-effort */ }
    // Live MQTT topics with msg/s + category.
    try {
      const tp = await fetchTopics();
      topic = (tp.topics || []).map((t) => ({
        value: t.topic,
        label: t.topic,
        sub: `${t.rate_per_sec.toFixed(1)} msg/s`,
        group: t.category,
      }));
    } catch { /* best-effort */ }
    // Broker from agent.yaml.
    try {
      const c = await fetchConfig();
      setConfigDefaults(c);
      if (c.mqtt?.broker) broker = [{ value: c.mqtt.broker, label: c.mqtt.broker, badge: 'config' }];
    } catch { /* best-effort */ }
    if (broker.length === 0) broker = [{ value: 'tcp://localhost:1883', label: 'tcp://localhost:1883', badge: 'défaut' }];
    setFieldPickers({ machine_id, topic, broker, node_id });
  }

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
      // Seed sensible default config so config fields + pickers appear.
      const config = fn.type === 'connector' ? defaultConfigFor(fn.name) : defaultFunctionConfig(fn.name);
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

  const handleDeleteNode = useCallback(
    (id) => {
      setNodes((nds) => nds.filter((n) => n.id !== id));
      setEdges((eds) => eds.filter((e) => e.source !== id && e.target !== id));
      setSelectedId(null);
    },
    [setNodes, setEdges]
  );

  // Refresh picker options (live tags / machines / topics) each time a node is
  // selected, so the opcua_read tag picker always shows current discovered tags.
  useEffect(() => {
    if (selectedId) loadPickerOptions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedId]);

  const selectedNode = nodes.find((n) => n.id === selectedId) || null;
  const connectors = functions.filter((f) => f.type === 'connector');
  const functionPickerOptions = functions
    .filter((f) => f.type !== 'connector')
    .map((f) => ({ value: f.name, label: f.name, sub: f.description, group: getCategory(f.type), fnType: f.type }));

  function addFunctionNode(o) {
    const id = `${slugify(o.value)}_${Math.random().toString(36).slice(2, 7)}`;
    setNodes((nds) =>
      nds.concat({
        id,
        type: 'pipelineNode',
        position: { x: 300 + Math.random() * 160, y: 90 + Math.random() * 240 },
        data: { name: o.value, type: o.fnType, function: o.value, config: defaultFunctionConfig(o.value), category: getCategory(o.fnType) },
      })
    );
    setShowFnPicker(false);
  }

  function handleNameChange(value) {
    setMeta((m) => ({ ...m, name: value, id: m.id || slugify(value) }));
  }

  // Smart, actionable validation. Returns an error string or null.
  function validate() {
    if (!meta.name) return '❌ Veuillez donner un titre à votre pipeline (ex: "Micro-stop Detection").';
    if (!meta.id) return '❌ Veuillez donner un nom à votre pipeline. Le nom sera utilisé comme identifiant unique.';

    const trigger = nodes.find((n) => n.type === 'triggerNode');
    if (!trigger || !trigger.data.function) return '❌ Aucun connecteur (trigger) trouvé. Veuillez ajouter un connecteur dans la zone ENTRÉE.';
    if (trigger.data.function === 'mqtt_subscribe' && !trigger.data.config?.topic)
      return '❌ Veuillez sélectionner un topic pour "mqtt_subscribe".';
    if (trigger.data.function === 'opcua_read' && !(trigger.data.config?.tags?.length || trigger.data.config?.node_id))
      return '❌ Veuillez sélectionner au moins un tag pour "opcua_read".';

    const steps = nodes.filter((n) => n.type === 'pipelineNode');
    if (!steps.some((n) => n.data.type === 'output'))
      return '❌ Aucune sortie trouvée. Veuillez ajouter "mqtt_publish" ou "add_to_dashboard" dans la zone SORTIE.';

    for (const n of steps) {
      if (n.data.function === 'state_machine' && !n.data.config?.machine_id)
        return '❌ Veuillez sélectionner une machine pour "state_machine".';
      if (n.data.function === 'opcua_read' && !(n.data.config?.tags?.length || n.data.config?.node_id))
        return '❌ Veuillez sélectionner au moins un tag pour "opcua_read".';
    }
    return null;
  }

  // A signature for duplicate detection: trigger topic/tags + the set of functions.
  function signature(p) {
    const fns = (p.nodes || []).map((n) => n.function).sort().join(',');
    const trig = `${p.trigger?.function || ''}:${p.trigger?.config?.topic || ''}:${(p.trigger?.config?.tags || []).slice().sort().join('|')}`;
    return `${trig}#${fns}`;
  }

  async function doSave() {
    try {
      const pipeline = flowToPipeline({ ...meta, version: '1.0', nodes, edges });
      await createPipeline(pipeline);
      setStatus({ type: 'ok', msg: `✅ Pipeline « ${meta.name} » sauvegardé — visible dans le Knowledge Graph.` });
      refreshPipelines();
    } catch (err) {
      setStatus({ type: 'error', msg: err.message });
    }
  }

  async function handleSave() {
    const err = validate();
    if (err) {
      setStatus({ type: 'error', msg: err });
      return;
    }
    // Duplicate check (same trigger + functions as an existing pipeline, different id).
    const sig = signature(flowToPipeline({ ...meta, version: '1.0', nodes, edges }));
    const existing = pipelines.find((p) => p.id !== meta.id && signature(p) === sig);
    if (existing) {
      setDupModal({ existing });
      return;
    }
    doSave();
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
      config: {
        ...defaultConfigFor(pendingConnector.name),
        ...(pendingConnector.name === 'opcua_read' && configDefaults?.opcua?.endpoint
          ? { endpoint: configDefaults.opcua.endpoint }
          : {}),
      },
    }));
    setStatus({ type: 'ok', msg: `Connecteur « ${pendingConnector.name} » appliqué au trigger.` });
    clearPending();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pendingConnector]);

  useEffect(() => {
    if (!pipelineToLoad) return;
    const p = pipelineToLoad;
    const flow = pipelineToFlow(p);
    setNodes(flow.nodes);
    setEdges(flow.edges);
    setMeta({ id: p.id, name: p.name, description: p.description || '' });
    setSelectedId(null);
    setStatus({ type: 'ok', msg: `Pipeline « ${p.name || p.id} » chargé.` });
    clearPending();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pipelineToLoad]);

  return (
    <div className="flex h-full">
      <Palette functions={functions.filter((f) => f.type !== 'connector')} loading={loading} error={error} />

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
            onClick={() => setShowFnPicker(true)}
            className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition"
          >
            ➕ Fonction
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
            onNodeClick={(_, n) => setSelectedId(n.type === 'zoneNode' ? null : n.id)}
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
        fieldPickers={fieldPickers}
        configDefaults={configDefaults}
        machines={machines}
        onChange={updateNodeData}
        onDelete={handleDeleteNode}
        onClose={() => setSelectedId(null)}
      />

      {showFnPicker && (
        <PickerModal
          title="⚙️ Ajouter une fonction"
          options={functionPickerOptions}
          onSelect={addFunctionNode}
          onClose={() => setShowFnPicker(false)}
        />
      )}

      {dupModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" onClick={() => setDupModal(null)}>
          <div className="bg-dark-900 border border-dark-700 rounded-xl w-full max-w-md p-5 shadow-2xl" onClick={(e) => e.stopPropagation()}>
            <h3 className="text-white font-semibold mb-2">⚠️ Pipeline en double</h3>
            <p className="text-sm text-dark-300 mb-4">
              Une pipeline avec cette configuration existe déjà : « {dupModal.existing.name} ».
            </p>
            <div className="flex justify-end gap-2">
              <button
                onClick={() => { setMeta((m) => ({ ...m, id: dupModal.existing.id, name: dupModal.existing.name })); setDupModal(null); doSave(); }}
                className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md"
              >
                Modifier l'existante
              </button>
              <button
                onClick={() => { const s = `${meta.id}_v2`; setMeta((m) => ({ ...m, id: s, name: `${m.name} v2` })); setDupModal(null); setStatus({ type: 'pending', msg: 'Renommée — cliquez Sauvegarder.' }); }}
                className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-1.5 rounded-md"
              >
                Nouvelle version
              </button>
              <button onClick={() => setDupModal(null)} className="text-dark-400 hover:text-white text-sm px-3 py-1.5">Annuler</button>
            </div>
          </div>
        </div>
      )}
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
