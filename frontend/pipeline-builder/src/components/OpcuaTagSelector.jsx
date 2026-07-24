import { useState, useMemo, Fragment } from 'react';
import { useTranslation } from 'react-i18next';
import { opcuaDiscover, opcuaSubscribe } from '../api/client';

// Discovered-tag table — ISA-95 selection only (Entry 125: Brut/Les deux and
// the Type column removed from display; raw storage is still what always
// happens underneath, this table just no longer offers it as a routing choice).
// Plus, since Entry 124, a preview of each tag's auto-computed ISA-95 mapping
// (Area/WorkCenter/WorkUnit/Tag) with a confidence score — editable before
// Apply, so a wrong auto-guess can be corrected instead of just accepted or
// rejected after the fact.
// Governance: only ISA-95-selected tags are published to mindset/site/# and
// can be used by functions.
export default function OpcuaTagSelector({ onApplied }) {
  const { t: tr } = useTranslation();
  const [tags, setTags] = useState([]);
  const [selections, setSelections] = useState({}); // node_id -> mode
  const [overrides, setOverrides] = useState({}); // node_id -> { area, work_center, work_unit, tag_name }
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
      const discoveredTags = data.tags || [];
      setTags(discoveredTags);
      // Pre-fill the editable ISA-95 fields with the server's own guess, so
      // editing means changing a value, not typing one from scratch.
      const initial = {};
      for (const t of discoveredTags) {
        initial[t.node_id] = {
          area: t.area || '',
          work_center: t.work_center || '',
          work_unit: t.work_unit || '',
          tag_name: t.tag_name || '',
        };
      }
      setOverrides(initial);
      setDiscovered(true);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const setOverrideField = (nodeId, field, value) =>
    setOverrides((o) => ({ ...o, [nodeId]: { ...o[nodeId], [field]: value } }));

  // Only send a field if it actually differs from the server's own guess —
  // keeps the payload minimal and makes "untouched" vs. "corrected" visible
  // server-side (an empty override field means "use the mapper's guess").
  const overridePayload = (t) => {
    const o = overrides[t.node_id] || {};
    const diff = {};
    if (o.area && o.area !== t.area) diff.area = o.area;
    if (o.work_center && o.work_center !== t.work_center) diff.work_center = o.work_center;
    if (o.work_unit && o.work_unit !== t.work_unit) diff.work_unit = o.work_unit;
    if (o.tag_name && o.tag_name !== t.tag_name) diff.tag_name = o.tag_name;
    return diff;
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
    const tagById = new Map(tags.map((t) => [t.node_id, t]));
    const payload = Object.entries(selections).map(([node_id, mode]) => {
      const t = tagById.get(node_id);
      return { node_id, mode, ...(t ? overridePayload(t) : {}) };
    });
    if (payload.length === 0) {
      setError(tr('opcuaTags.selectAtLeastOne'));
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
        <h3 className="font-medium text-white">{tr('opcuaTags.title')}</h3>
        <span className="text-xs text-dark-400">{discovered ? `${tags.length} tag(s)` : ''}</span>
        <button
          onClick={handleDiscover}
          disabled={loading}
          className="ml-auto bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-xs px-3 py-1.5 rounded-md transition"
        >
          {loading ? tr('opcuaTags.discovering') : discovered ? `↻ ${tr('opcuaTags.rediscover')}` : `🔍 ${tr('opcuaTags.discover')}`}
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
              placeholder={`🔍 ${tr('opcuaTags.filterByName')}`}
              value={filterName}
              onChange={(e) => setFilterName(e.target.value)}
              className="bg-dark-800 border border-dark-600 rounded-md px-2 py-1 text-sm text-white flex-1 min-w-[140px]"
            />
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="bg-dark-800 border border-dark-600 rounded-md px-2 py-1 text-sm text-white"
            >
              {types.map((ty) => (
                <option key={ty} value={ty}>{ty}</option>
              ))}
            </select>
            <button
              onClick={() => bulk('isa95')}
              className="text-xs px-2 py-1 rounded-md bg-dark-700 hover:bg-dark-600 text-white"
            >
              {tr('opcuaTags.allIsa95')}
            </button>
            <button
              onClick={() => setSelections({})}
              className="text-xs px-2 py-1 rounded-md bg-dark-700 hover:bg-dark-600 text-white"
            >
              {tr('opcuaTags.clear')}
            </button>
          </div>

          {/* Tag table */}
          <div className="max-h-[42vh] overflow-y-auto border border-dark-700 rounded-md">
            <table className="w-full text-sm">
              <thead className="sticky top-0 bg-dark-800 text-dark-400 text-xs">
                <tr>
                  <th className="text-left px-2 py-1.5">Tag</th>
                  <th className="text-left px-2 py-1.5">{tr('opcuaTags.value')}</th>
                  <th className="text-center px-2 py-1.5">ISA-95</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((t) => {
                  const mode = selections[t.node_id];
                  const o = overrides[t.node_id] || {};
                  const isaActive = mode === 'isa95';
                  const pct = Math.round((t.confidence ?? 0) * 100);
                  const confColor = pct >= 70 ? 'text-status-running' : 'text-status-warn';
                  return (
                    <Fragment key={t.node_id}>
                    <tr className="border-t border-dark-700/60 hover:bg-dark-800/40">
                      <td className="px-2 py-1.5">
                        <div className="text-white">{t.name}</div>
                        <div className="text-[10px] text-dark-500 font-mono">{t.node_id}</div>
                      </td>
                      <td className="px-2 py-1.5 text-dark-400 font-mono text-xs">
                        {t.value === null || t.value === undefined ? '' : String(t.value)}
                      </td>
                      <td className="text-center px-2 py-1.5">
                        <input
                          type="checkbox"
                          checked={isaActive}
                          onChange={() => setMode(t.node_id, 'isa95')}
                          className="accent-blue-500 cursor-pointer"
                          title={tr('opcuaTags.publishIsa95Title')}
                        />
                      </td>
                    </tr>
                    {/* ISA-95 preview + edit (Entry 124) — shown for every discovered
                        tag so confidence is visible before deciding a mode; inputs
                        only matter once ISA-95/Les deux is selected, so they're
                        dimmed (not disabled — still editable in advance) otherwise. */}
                    <tr className={`border-t-0 ${isaActive ? '' : 'opacity-50'}`}>
                      <td colSpan={3} className="px-2 pb-2">
                        <div className="flex flex-wrap items-center gap-1.5 text-[11px]">
                          <span
                            className={`font-mono ${confColor}`}
                            title={tr('opcuaTags.confidenceTitle')}
                          >
                            {pct}%
                          </span>
                          {t.pending && (
                            <span className="text-status-warn" title={tr('opcuaTags.pendingTitle')}>
                              ⏸
                            </span>
                          )}
                          <span className="text-dark-500 font-mono">{t.site}</span>
                          <span className="text-dark-600">/</span>
                          <input
                            value={o.area}
                            onChange={(e) => setOverrideField(t.node_id, 'area', e.target.value)}
                            placeholder="area"
                            className="bg-dark-800 border border-dark-600 rounded px-1 py-0.5 text-white font-mono w-24"
                          />
                          <span className="text-dark-600">/</span>
                          <input
                            value={o.work_center}
                            onChange={(e) => setOverrideField(t.node_id, 'work_center', e.target.value)}
                            placeholder="work_center"
                            className="bg-dark-800 border border-dark-600 rounded px-1 py-0.5 text-white font-mono w-28"
                          />
                          {(o.work_unit || t.work_unit) && (
                            <>
                              <span className="text-dark-600">/</span>
                              <input
                                value={o.work_unit}
                                onChange={(e) => setOverrideField(t.node_id, 'work_unit', e.target.value)}
                                placeholder="work_unit"
                                className="bg-dark-800 border border-dark-600 rounded px-1 py-0.5 text-white font-mono w-24"
                              />
                            </>
                          )}
                          <span className="text-dark-600">/</span>
                          <input
                            value={o.tag_name}
                            onChange={(e) => setOverrideField(t.node_id, 'tag_name', e.target.value)}
                            placeholder="tag_name"
                            className="bg-dark-800 border border-dark-600 rounded px-1 py-0.5 text-white font-mono w-28"
                          />
                        </div>
                      </td>
                    </tr>
                    </Fragment>
                  );
                })}
                {filtered.length === 0 && (
                  <tr>
                    <td colSpan={3} className="text-center text-dark-500 py-4 text-xs">
                      {tr('opcuaTags.noneMatchFilter')}
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* Legend + apply */}
          <div className="flex items-center gap-3 mt-3">
            <p className="text-[11px] text-dark-500">
              <span className="text-blue-300">ISA-95</span> {tr('opcuaTags.isa95Legend')}
            </p>
            <button
              onClick={handleApply}
              disabled={applying || selectedCount === 0}
              className="ml-auto bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white text-sm px-4 py-1.5 rounded-md transition"
            >
              {applying ? tr('opcuaTags.applying') : `✅ ${tr('opcuaTags.apply')} (${selectedCount})`}
            </button>
          </div>
        </>
      )}

      {discovered && tags.length === 0 && !loading && (
        <p className="text-dark-500 text-sm">{tr('opcuaTags.noneDiscovered')}</p>
      )}
    </div>
  );
}
