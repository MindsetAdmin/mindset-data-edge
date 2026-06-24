// Right-hand panel to configure the selected node: a guided panel with the
// function's description, labelled fields with help/examples, pickers, a cost
// preview, and an OPC-UA machine/tag selector.

import { useState } from 'react';
import PickerModal from './PickerModal';
import { defaultConfigFor, triggerTypeFor } from '../lib/connectorTemplates';
import { typeStyle, getCategory } from '../lib/functionMeta';
import { functionDoc, fieldDoc } from '../lib/functionDocs';

function coerce(v) {
  if (typeof v !== 'string') return v; // arrays / numbers pass through untouched
  if (v === '') return '';
  if (v === 'true') return true;
  if (v === 'false') return false;
  const n = Number(v);
  return v.trim() !== '' && !Number.isNaN(n) ? n : v;
}

export default function NodeConfigPanel({ node, connectors = [], fieldPickers = {}, configDefaults = null, machines = [], onChange, onDelete, onClose }) {
  const [pickerKey, setPickerKey] = useState(null);

  if (!node) {
    return (
      <aside className="w-72 bg-dark-900 border-l border-dark-700 p-4 shrink-0">
        <p className="text-sm text-dark-500">Sélectionnez un nœud pour le configurer.</p>
      </aside>
    );
  }

  const isTrigger = node.type === 'triggerNode';
  const data = node.data;
  const fn = data.function;
  const config = data.config || {};
  const doc = functionDoc(fn);
  const style = typeStyle(isTrigger ? 'connector' : data.type);
  const category = getCategory(isTrigger ? 'connector' : data.type);

  const setField = (key, value) => onChange(node.id, (d) => ({ ...d, [key]: value }));
  const setConfig = (key, value) =>
    onChange(node.id, (d) => ({ ...d, config: { ...d.config, [key]: coerce(value) } }));
  const setConfigRaw = (key, value) =>
    onChange(node.id, (d) => ({ ...d, config: { ...d.config, [key]: value } }));
  const renameKey = (oldKey, newKey) =>
    onChange(node.id, (d) => {
      const c = { ...d.config };
      const val = c[oldKey];
      delete c[oldKey];
      c[newKey] = val;
      return { ...d, config: c };
    });
  const removeKey = (key) =>
    onChange(node.id, (d) => {
      const c = { ...d.config };
      delete c[key];
      return { ...d, config: c };
    });
  const addKey = () => onChange(node.id, (d) => ({ ...d, config: { ...d.config, ['nouveau_champ']: '' } }));

  const pickConnector = (fnName) =>
    onChange(node.id, (d) => ({
      ...d,
      function: fnName,
      triggerType: triggerTypeFor(fnName),
      config: {
        ...defaultConfigFor(fnName),
        ...(fnName === 'opcua_read' && configDefaults?.opcua?.endpoint ? { endpoint: configDefaults.opcua.endpoint } : {}),
      },
    }));

  const isOpcua = fn === 'opcua_read';
  const isCost = fn === 'calculate_cost';
  const hidden = isOpcua ? new Set(['tags', 'node_id', 'machine']) : new Set();

  return (
    <aside className="w-72 bg-dark-900 border-l border-dark-700 p-4 overflow-y-auto shrink-0">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-blue-400">📋 Configuration</h3>
        <button onClick={onClose} className="text-dark-400 hover:text-white text-lg leading-none">×</button>
      </div>

      {/* Function header: icon + name + category badge + description */}
      <div className={`rounded-lg border ${style.card} p-2.5 mb-4`}>
        <div className="flex items-center gap-2">
          <span className="text-lg">{style.icon}</span>
          <span className="text-sm font-medium text-white">{fn || 'trigger'}</span>
          <span className={`ml-auto text-[10px] px-2 py-0.5 rounded ${style.chip}`}>{category}</span>
        </div>
        {doc.description && <p className="text-[11px] text-dark-300 mt-1.5">{doc.description}</p>}
      </div>

      <Field label="ID du nœud">
        <input value={node.id} disabled className="input opacity-60" />
      </Field>

      {isTrigger ? (
        <Field label="Connecteur" help="Source de données qui déclenche le pipeline.">
          <select value={data.function || ''} onChange={(e) => pickConnector(e.target.value)} className="input">
            {data.function && !connectors.some((c) => c.name === data.function) && (
              <option value={data.function}>{data.function}</option>
            )}
            {connectors.map((c) => (
              <option key={c.name} value={c.name}>{c.name}</option>
            ))}
          </select>
        </Field>
      ) : (
        <Field label="Nom" help="Nom affiché du nœud.">
          <input value={data.name || ''} onChange={(e) => setField('name', e.target.value)} className="input" />
        </Field>
      )}

      {/* OPC-UA machine + tag selection */}
      {isOpcua && <OpcuaTagSelector machines={machines} config={config} setConfigRaw={setConfigRaw} />}

      {/* Config fields (labelled, with help + example) */}
      <div className="mt-4 mb-2 flex items-center justify-between">
        <span className="text-xs font-semibold text-dark-300 uppercase tracking-wider">Paramètres</span>
        <button onClick={addKey} className="text-xs text-blue-400 hover:text-blue-300">+ champ</button>
      </div>

      <div className="space-y-3">
        {Object.entries(config).filter(([k]) => !hidden.has(k)).length === 0 && (
          <p className="text-xs text-dark-500">Aucun paramètre.</p>
        )}
        {Object.entries(config)
          .filter(([k]) => !hidden.has(k))
          .map(([key, value]) => {
            const fd = fieldDoc(fn, key);
            const known = !!doc.fields[key];
            const hasPicker = (fieldPickers[key] || []).length > 0;
            return (
              <div key={key}>
                {known ? (
                  <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{fd.label}</label>
                ) : (
                  <input
                    defaultValue={key}
                    onBlur={(e) => e.target.value !== key && e.target.value && renameKey(key, e.target.value)}
                    className="input text-[11px] mb-1 w-32"
                  />
                )}
                <div className="flex items-center gap-1">
                  <input
                    value={String(value ?? '')}
                    placeholder={fd.example}
                    onChange={(e) => setConfig(key, e.target.value)}
                    className="input flex-1 text-xs"
                  />
                  {hasPicker && (
                    <button onClick={() => setPickerKey(key)} title="Choisir" className="text-dark-400 hover:text-blue-400 px-1">📋</button>
                  )}
                  <button onClick={() => removeKey(key)} title="Supprimer" className="text-dark-500 hover:text-red-400 px-1">×</button>
                </div>
                {fd.help && <p className="text-[10px] text-dark-500 mt-0.5">ⓘ {fd.help}</p>}
              </div>
            );
          })}
      </div>

      {/* Cost live preview */}
      {isCost && <CostPreview config={config} configDefaults={configDefaults} />}

      {!isTrigger && onDelete && (
        <button onClick={() => onDelete(node.id)} className="mt-5 w-full bg-red-600/80 hover:bg-red-600 text-white text-sm px-3 py-1.5 rounded-md transition">
          🗑️ Supprimer le nœud
        </button>
      )}

      {pickerKey && (
        <PickerModal
          title={`Choisir : ${fieldDoc(fn, pickerKey).label}`}
          options={fieldPickers[pickerKey] || []}
          allowCustom
          customLabel={`${pickerKey} personnalisé`}
          onSelect={(o) => { setConfig(pickerKey, o.value); setPickerKey(null); }}
          onClose={() => setPickerKey(null)}
        />
      )}

      <style>{`
        .input { width: 100%; background:#0f172a; border:1px solid #334155; border-radius:6px; padding:6px 8px; color:#e2e8f0; font-size:13px; }
        .input:focus { outline:none; border-color:#3b82f6; }
      `}</style>
    </aside>
  );
}

function Field({ label, help, children }) {
  return (
    <div className="mb-3">
      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{label}</label>
      {children}
      {help && <p className="text-[10px] text-dark-500 mt-0.5">ⓘ {help}</p>}
    </div>
  );
}

// OPC-UA machine dropdown + tag checkbox list with filter, select-all, live values.
function OpcuaTagSelector({ machines, config, setConfigRaw }) {
  const [q, setQ] = useState('');
  const selectedMachine = config.machine || machines[0]?.work_center || '';
  const machine = machines.find((m) => m.work_center === selectedMachine);
  const allTags = machine?.tags || [];
  const tags = allTags.filter((t) => !q || (t.name || '').toLowerCase().includes(q.toLowerCase()));
  const selected = new Set(config.tags || []);

  const applyTags = (arr) => {
    setConfigRaw('tags', arr);
    setConfigRaw('node_id', arr[0] || '');
  };
  const toggle = (id) => {
    const next = new Set(selected);
    next.has(id) ? next.delete(id) : next.add(id);
    applyTags([...next]);
  };

  return (
    <div className="mb-3 border border-dark-700 rounded-md p-2.5">
      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">Machine</label>
      <select value={selectedMachine} onChange={(e) => setConfigRaw('machine', e.target.value)} className="input">
        {machines.length === 0 && <option value="">(aucune machine découverte)</option>}
        {machines.map((m) => (
          <option key={m.work_center} value={m.work_center}>
            {m.work_center} ({(m.tags || []).length} tags)
          </option>
        ))}
      </select>

      <input value={q} onChange={(e) => setQ(e.target.value)} placeholder="🔍 filtrer les tags…" className="input mt-2 text-xs" />
      <div className="flex gap-3 mt-1">
        <button onClick={() => applyTags(allTags.map((t) => t.node_id))} className="text-[11px] text-blue-400 hover:text-blue-300">Tout sélectionner</button>
        <button onClick={() => applyTags([])} className="text-[11px] text-dark-400 hover:text-white">Tout désélectionner</button>
      </div>

      <div className="mt-2 max-h-48 overflow-y-auto space-y-1">
        {tags.length === 0 && <p className="text-xs text-dark-500">Aucun tag (l'agent doit tourner).</p>}
        {tags.map((t) => (
          <label key={t.node_id} className="flex items-center gap-2 text-xs cursor-pointer">
            <input type="checkbox" checked={selected.has(t.node_id)} onChange={() => toggle(t.node_id)} />
            <span className="flex-1 truncate text-dark-200">{t.name || t.node_id}</span>
            <span className="text-dark-500 font-mono">{String(t.value)}</span>
          </label>
        ))}
      </div>
      <div className="text-[11px] text-dark-500 mt-1">{(config.tags || []).length} tag(s) sélectionné(s)</div>
    </div>
  );
}

// Live cost preview for typical durations.
function CostPreview({ config, configDefaults }) {
  const rate = Number(config.hourly_rate ?? configDefaults?.cost?.hourly_rate ?? 85);
  const cur = config.currency || 'EUR';
  const rows = [
    ['30s', 30],
    ['1min', 60],
    ['3min', 180],
    ['5min', 300],
  ];
  return (
    <div className="mt-3 bg-dark-950 border border-dark-700 rounded-md p-2.5">
      <div className="text-[11px] text-dark-400 uppercase tracking-wider mb-1">Aperçu du coût ({rate} €/h)</div>
      {rows.map(([label, s]) => (
        <div key={label} className="flex justify-between text-xs py-0.5">
          <span className="text-dark-300">{label}</span>
          <span className="font-mono text-green-400">{((s / 3600) * rate).toFixed(2)} {cur}</span>
        </div>
      ))}
    </div>
  );
}
