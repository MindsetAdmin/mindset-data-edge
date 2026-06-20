// Right-hand panel to edit the selected node: its name (steps) or trigger
// fields, plus arbitrary config key/value pairs. Values are coerced to
// number/boolean where possible so the saved YAML keeps proper types.

import { useState } from 'react';
import PickerModal from './PickerModal';
import { defaultConfigFor, triggerTypeFor } from '../lib/connectorTemplates';

function coerce(v) {
  if (v === '' || v == null) return '';
  if (v === 'true') return true;
  if (v === 'false') return false;
  const n = Number(v);
  return v.trim() !== '' && !Number.isNaN(n) ? n : v;
}

export default function NodeConfigPanel({ node, connectors = [], fieldPickers = {}, onChange, onDelete, onClose }) {
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
  const config = data.config || {};

  const setField = (key, value) => onChange(node.id, (d) => ({ ...d, [key]: value }));
  const setConfig = (key, value) =>
    onChange(node.id, (d) => ({ ...d, config: { ...d.config, [key]: coerce(value) } }));
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
  const addKey = () =>
    onChange(node.id, (d) => ({ ...d, config: { ...d.config, ['nouveau_champ']: '' } }));

  // Picking a connector on the trigger swaps the function, trigger type, and
  // resets config to that connector's default template.
  const pickConnector = (fnName) =>
    onChange(node.id, (d) => ({
      ...d,
      function: fnName,
      triggerType: triggerTypeFor(fnName),
      config: defaultConfigFor(fnName),
    }));

  return (
    <aside className="w-72 bg-dark-900 border-l border-dark-700 p-4 overflow-y-auto shrink-0">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-semibold text-blue-400">📋 Configuration</h3>
        <button onClick={onClose} className="text-dark-400 hover:text-white text-lg leading-none">×</button>
      </div>

      <Field label="ID du nœud">
        <input value={node.id} disabled className="input opacity-60" />
      </Field>

      {isTrigger ? (
        <>
          <Field label="Connecteur">
            <select value={data.function || ''} onChange={(e) => pickConnector(e.target.value)} className="input">
              {data.function && !connectors.some((c) => c.name === data.function) && (
                <option value={data.function}>{data.function}</option>
              )}
              {connectors.map((c) => (
                <option key={c.name} value={c.name}>{c.name}</option>
              ))}
            </select>
          </Field>
          <Field label="Type de trigger">
            <input value={data.triggerType || ''} onChange={(e) => setField('triggerType', e.target.value)} className="input" />
          </Field>
        </>
      ) : (
        <>
          <Field label="Nom">
            <input value={data.name || ''} onChange={(e) => setField('name', e.target.value)} className="input" />
          </Field>
          <Field label="Type / Fonction">
            <input value={`${data.type} · ${data.function}`} disabled className="input opacity-60" />
          </Field>
        </>
      )}

      <div className="mt-4 mb-2 flex items-center justify-between">
        <span className="text-xs font-semibold text-dark-300 uppercase tracking-wider">Config</span>
        <button onClick={addKey} className="text-xs text-blue-400 hover:text-blue-300">+ champ</button>
      </div>

      <div className="space-y-2">
        {Object.keys(config).length === 0 && (
          <p className="text-xs text-dark-500">Aucun paramètre. Ajoutez un champ.</p>
        )}
        {Object.entries(config).map(([key, value]) => {
          const hasPicker = (fieldPickers[key] || []).length > 0;
          return (
            <div key={key} className="flex items-center gap-1">
              <input
                defaultValue={key}
                onBlur={(e) => e.target.value !== key && e.target.value && renameKey(key, e.target.value)}
                className="input w-20 text-xs"
              />
              <input
                value={String(value ?? '')}
                onChange={(e) => setConfig(key, e.target.value)}
                className="input flex-1 text-xs"
              />
              {hasPicker && (
                <button
                  onClick={() => setPickerKey(key)}
                  title="Choisir parmi les options disponibles"
                  className="text-dark-400 hover:text-blue-400 px-1"
                >
                  📋
                </button>
              )}
              <button onClick={() => removeKey(key)} className="text-dark-500 hover:text-red-400 px-1">×</button>
            </div>
          );
        })}
      </div>

      {!isTrigger && onDelete && (
        <button
          onClick={() => onDelete(node.id)}
          className="mt-5 w-full bg-red-600/80 hover:bg-red-600 text-white text-sm px-3 py-1.5 rounded-md transition"
        >
          🗑️ Supprimer le nœud
        </button>
      )}

      {pickerKey && (
        <PickerModal
          title={`Choisir : ${pickerKey}`}
          options={fieldPickers[pickerKey] || []}
          allowCustom
          customLabel={`${pickerKey} personnalisé`}
          onSelect={(o) => {
            setConfig(pickerKey, o.value);
            setPickerKey(null);
          }}
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

function Field({ label, children }) {
  return (
    <div className="mb-3">
      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{label}</label>
      {children}
    </div>
  );
}
