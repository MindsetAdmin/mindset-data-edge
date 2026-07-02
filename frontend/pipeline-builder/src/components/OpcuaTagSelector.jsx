import { useState, useMemo } from 'react';
import { opcuaDiscover, opcuaSubscribe } from '../api/client';

const MODES = [
  { key: 'raw', label: 'Brut', hint: 'stockage uniquement' },
  { key: 'isa95', label: 'ISA-95', hint: 'utilisable dans les fonctions' },
  { key: 'both', label: 'Les deux', hint: 'stockage + fonctions' },
];

// Discovered-tag table with per-row data-flow choice (Raw / ISA-95 / Both).
// Governance: only ISA-95 / Both tags are published to mindset/site/# and can be
// used by functions; Raw tags are storage-only.
export default function OpcuaTagSelector({ onApplied }) {
  const [tags, setTags] = useState([]);
  const [selections, setSelections] = useState({}); // node_id -> mode
  const [filterName, setFilterName] = useState('');
  const [filterType, setFilterType] = useState('All');
  const [loading, setLoading] = useState(false);
  const [applying, setApplying] = useState(false);
  const [error, setError] = useState(null);
  const [discovered, setDiscovered] = useState(false);

  const types = useMemo(() => {
    const s = new Set(tags.map((t) => t.data_type).filter(Boolean));
    return ['All', ...Array.from(s).sort()];
  }, [tags]);

  const filtered = useMemo(() => {
    const q = filterName.trim().toLowerCase();
    return tags.filter(
      (t) =>
        (filterType === 'All' || t.data_type === filterType) &&
        (q === '' || (t.name || '').toLowerCase().includes(q))
    );
  }, [tags, filterName, filterType]);

  const handleDiscover = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await opcuaDiscover();
      setTags(data.tags || []);
      setDiscovered(true);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const setMode = (nodeId, mode) =>
    setSelections((s) => {
      const next = { ...s };
      if (next[nodeId] === mode) delete next[nodeId]; // toggle off
      else next[nodeId] = mode;
      return next;
    });

  const bulk = (mode) =>
    setSelections(() => {
      const next = {};
      for (const t of filtered) next[t.node_id] = mode;
      return next;
    });

  const selectedCount = Object.keys(selections).length;

  const handleApply = async () => {
    const payload = Object.entries(selections).map(([node_id, mode]) => ({ node_id, mode }));
    if (payload.length === 0) {
      setError('Sélectionnez au moins un tag (Brut, ISA-95 ou Les deux).');
      return;
    }
    setApplying(true);
    setError(null);
    try {
      await opcuaSubscribe(payload);
      onApplied?.(payload);
    } catch (err) {
      setError(err.message);
    } finally {
      setApplying(false);
    }
  };

  return (
    <div className="border border-dark-700 bg-dark-900 rounded-lg p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-xl">📋</span>
        <h3 className="font-medium text-white">Tags découverts</h3>
        <span className="text-xs text-dark-400">{discovered ? `${tags.length} tag(s)` : ''}</span>
        <button
          onClick={handleDiscover}
          disabled={loading}
          className="ml-auto bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-xs px-3 py-1.5 rounded-md transition"
        >
          {loading ? 'Découverte…' : discovered ? '↻ Re-découvrir' : '🔍 Découvrir les tags'}
        </button>
      </div>

      {error && (
        <div className="bg-red-500/15 border border-red-500/40 rounded-md p-2 text-red-400 text-xs mb-3">
          ❌ {error}
        </div>
      )}

      {discovered && tags.length > 0 && (
        <>
          {/* Filters + bulk */}
          <div className="flex flex-wrap items-center gap-2 mb-3">
            <input
              placeholder="🔍 Filtrer par nom…"
              value={filterName}
              onChange={(e) => setFilterName(e.target.value)}
              className="bg-dark-800 border border-dark-600 rounded-md px-2 py-1 text-sm text-white flex-1 min-w-[140px]"
            />
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="bg-dark-800 border border-dark-600 rounded-md px-2 py-1 text-sm text-white"
            >
              {types.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
            <button
              onClick={() => bulk('isa95')}
              className="text-xs px-2 py-1 rounded-md bg-dark-700 hover:bg-dark-600 text-white"
            >
              Tout en ISA-95
            </button>
            <button
              onClick={() => setSelections({})}
              className="text-xs px-2 py-1 rounded-md bg-dark-700 hover:bg-dark-600 text-white"
            >
              Effacer
            </button>
          </div>

          {/* Tag table */}
          <div className="max-h-[42vh] overflow-y-auto border border-dark-700 rounded-md">
            <table className="w-full text-sm">
              <thead className="sticky top-0 bg-dark-800 text-dark-400 text-xs">
                <tr>
                  <th className="text-left px-2 py-1.5">Tag</th>
                  <th className="text-left px-2 py-1.5">Type</th>
                  <th className="text-left px-2 py-1.5">Valeur</th>
                  <th className="text-center px-2 py-1.5">Brut</th>
                  <th className="text-center px-2 py-1.5">ISA-95</th>
                  <th className="text-center px-2 py-1.5">Les deux</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((t) => {
                  const mode = selections[t.node_id];
                  return (
                    <tr key={t.node_id} className="border-t border-dark-700/60 hover:bg-dark-800/40">
                      <td className="px-2 py-1.5">
                        <div className="text-white">{t.name}</div>
                        <div className="text-[10px] text-dark-500 font-mono">{t.node_id}</div>
                      </td>
                      <td className="px-2 py-1.5 text-dark-300 font-mono text-xs">{t.data_type}</td>
                      <td className="px-2 py-1.5 text-dark-400 font-mono text-xs">
                        {t.value === null || t.value === undefined ? '' : String(t.value)}
                      </td>
                      {MODES.map((m) => (
                        <td key={m.key} className="text-center px-2 py-1.5">
                          <input
                            type="radio"
                            name={`mode-${t.node_id}`}
                            checked={mode === m.key}
                            onChange={() => setMode(t.node_id, m.key)}
                            className="accent-blue-500 cursor-pointer"
                            title={m.hint}
                          />
                        </td>
                      ))}
                    </tr>
                  );
                })}
                {filtered.length === 0 && (
                  <tr>
                    <td colSpan={6} className="text-center text-dark-500 py-4 text-xs">
                      Aucun tag ne correspond au filtre.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* Legend + apply */}
          <div className="flex items-center gap-3 mt-3">
            <p className="text-[11px] text-dark-500">
              <span className="text-dark-300">Brut</span> = stockage seul ·{' '}
              <span className="text-blue-300">ISA-95</span> = utilisable dans les fonctions ·{' '}
              <span className="text-blue-300">Les deux</span> = stockage + fonctions
            </p>
            <button
              onClick={handleApply}
              disabled={applying || selectedCount === 0}
              className="ml-auto bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white text-sm px-4 py-1.5 rounded-md transition"
            >
              {applying ? 'Application…' : `✅ Appliquer (${selectedCount})`}
            </button>
          </div>
        </>
      )}

      {discovered && tags.length === 0 && !loading && (
        <p className="text-dark-500 text-sm">Aucun tag découvert sur ce serveur.</p>
      )}
    </div>
  );
}
