import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchConnections, previewConnection } from '../api/client';
import FieldMapEditor from './FieldMapEditor';

// Key/value grid for sql_query's `params` (named-placeholder values bound by
// :name in the query). Values are static today — the pipeline engine has no
// {{ trigger.x }} templating yet, so this only supports literal values.
function ParamsGrid({ params = {}, onChange }) {
  const { t } = useTranslation();
  const entries = Object.entries(params);

  const add = () => onChange({ ...params, ['param']: '' });
  const remove = (key) => {
    const next = { ...params };
    delete next[key];
    onChange(next);
  };
  const rename = (oldKey, newKey) => {
    if (!newKey || newKey === oldKey) return;
    const next = {};
    for (const [k, v] of Object.entries(params)) next[k === oldKey ? newKey : k] = v;
    onChange(next);
  };
  const setValue = (key, value) => onChange({ ...params, [key]: value });

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <label className="block text-[11px] text-dark-400 uppercase tracking-wider">{t('sqlConfig.params')}</label>
        <button onClick={add} className="text-xs text-blue-400 hover:text-blue-300">+ {t('sqlConfig.param')}</button>
      </div>
      {entries.length === 0 && <p className="text-xs text-dark-500 mb-1">{t('sqlConfig.noParams')}</p>}
      <div className="space-y-1">
        {entries.map(([key, value]) => (
          <div key={key} className="flex items-center gap-1">
            <input defaultValue={key} onBlur={(e) => rename(key, e.target.value)} className="input text-[11px] w-28 font-mono" />
            <span className="text-dark-500 text-xs">=</span>
            <input value={String(value ?? '')} onChange={(e) => setValue(key, e.target.value)} className="input text-xs flex-1 font-mono" />
            <button onClick={() => remove(key)} className="text-dark-500 hover:text-red-400 px-1">×</button>
          </div>
        ))}
      </div>
    </div>
  );
}

// Guided config for sql_query nodes: connection picker, SQL editor, params,
// timeout/limit, canonical type + field_map, and a Preview button that runs
// the query through the same guards as a real execution (capped at 5 rows).
export default function SqlConfigPanel({ config, setConfig, setConfigRaw }) {
  const { t } = useTranslation();
  const [connections, setConnections] = useState([]);
  const [loadingConns, setLoadingConns] = useState(true);
  const [previewing, setPreviewing] = useState(false);
  const [preview, setPreview] = useState(null); // { rows, canonical, canonical_type, row_count, query_ms }
  const [previewError, setPreviewError] = useState(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await fetchConnections();
        setConnections(data.connections || []);
      } catch {
        /* best-effort — the dropdown just shows the raw id if the list fails */
      } finally {
        setLoadingConns(false);
      }
    })();
  }, []);

  const handlePreview = async () => {
    if (!config.connection_id || !config.query) return;
    setPreviewing(true);
    setPreviewError(null);
    setPreview(null);
    try {
      const result = await previewConnection(config.connection_id, {
        query: config.query,
        params: config.params || {},
        limit: 5,
      });
      setPreview(result);
    } catch (err) {
      setPreviewError(err.message);
    } finally {
      setPreviewing(false);
    }
  };

  const previewCols = preview?.rows?.length ? Object.keys(preview.rows[0]) : [];

  return (
    <div className="mb-3 border border-dark-700 rounded-md p-2.5 space-y-1">
      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{t('sqlConfig.connection')}</label>
      <select
        value={config.connection_id || ''}
        onChange={(e) => setConfigRaw('connection_id', e.target.value)}
        className="input mb-2"
      >
        <option value="">{loadingConns ? t('common.loading') : t('sqlConfig.chooseConnection')}</option>
        {connections.map((c) => (
          <option key={c.id} value={c.id}>{c.name} ({c.id})</option>
        ))}
        {config.connection_id && !connections.some((c) => c.id === config.connection_id) && (
          <option value={config.connection_id}>{config.connection_id}</option>
        )}
      </select>
      {connections.length === 0 && !loadingConns && (
        <p className="text-[10px] text-yellow-400/80 mb-2">
          ⚠️ {t('sqlConfig.noConnectionPre')} <span className="font-mono">Connexions</span>
        </p>
      )}

      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{t('sqlConfig.sqlQuery')}</label>
      <textarea
        value={config.query || ''}
        onChange={(e) => setConfigRaw('query', e.target.value)}
        rows={5}
        spellCheck={false}
        placeholder={'SELECT of_number, product_code, status\nFROM work_orders\nWHERE work_center = :work_center\nLIMIT 1'}
        className="input font-mono text-xs resize-y"
      />
      <p className="text-[10px] text-dark-500 mb-2">
        {t('sqlConfig.paramHintPre')} <span className="font-mono">:name</span> {t('sqlConfig.paramHintPost')}
      </p>

      <ParamsGrid params={config.params || {}} onChange={(v) => setConfigRaw('params', v)} />

      <div className="grid grid-cols-2 gap-2 mb-2">
        <label className="block">
          <span className="text-[11px] text-dark-400 uppercase tracking-wider block mb-1">{t('sqlConfig.timeout')}</span>
          <input type="number" value={config.timeout_seconds ?? 30} onChange={(e) => setConfig('timeout_seconds', e.target.value)} className="input text-xs" />
        </label>
        <label className="block">
          <span className="text-[11px] text-dark-400 uppercase tracking-wider block mb-1">{t('sqlConfig.limit')}</span>
          <input type="number" value={config.limit ?? 1000} onChange={(e) => setConfig('limit', e.target.value)} className="input text-xs" />
        </label>
      </div>

      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{t('sqlConfig.canonicalType')}</label>
      <input
        value={config.canonical || ''}
        onChange={(e) => setConfigRaw('canonical', e.target.value)}
        placeholder="work_order"
        className="input text-xs font-mono mb-2"
      />

      <FieldMapEditor fieldMap={config.field_map || {}} onChange={(v) => setConfigRaw('field_map', v)} />

      <button
        onClick={handlePreview}
        disabled={previewing || !config.connection_id || !config.query}
        className="w-full bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-sm px-3 py-1.5 rounded-md transition"
      >
        {previewing ? t('sqlConfig.previewing') : `▶️ ${t('sqlConfig.previewButton')}`}
      </button>

      {previewError && (
        <div className="bg-red-500/15 border border-red-500/40 rounded-md p-2 text-red-400 text-xs mt-2">
          ❌ {previewError}
        </div>
      )}

      {preview && (
        <div className="mt-2 bg-dark-950 border border-dark-700 rounded-md p-2 overflow-x-auto">
          <div className="text-[11px] text-dark-400 mb-1">
            {t('sqlConfig.rowCount', { count: preview.row_count })} · {preview.query_ms}ms
            {preview.canonical_type ? ` · canonical: ${preview.canonical_type}` : ''}
          </div>
          {previewCols.length === 0 ? (
            <p className="text-xs text-dark-500">{t('sqlConfig.noRows')}</p>
          ) : (
            <table className="text-[11px] font-mono w-full">
              <thead>
                <tr>
                  {previewCols.map((c) => (
                    <th key={c} className="text-left text-dark-400 px-1.5 py-0.5 border-b border-dark-700">{c}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {preview.rows.map((row, i) => (
                  <tr key={i}>
                    {previewCols.map((c) => (
                      <td key={c} className="text-dark-200 px-1.5 py-0.5 border-b border-dark-800 whitespace-nowrap">
                        {row[c] === null || row[c] === undefined ? <span className="text-dark-600">null</span> : String(row[c])}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  );
}
