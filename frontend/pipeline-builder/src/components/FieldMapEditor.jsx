// Editor for sql_query's field_map: canonical field name -> either a raw
// column name (simple) or an enum translation {from, values} (per
// docs/mysql_connector.md §6b). Mirrors the generic key/value editor in
// NodeConfigPanel but understands the two-shape value.

import { useTranslation } from 'react-i18next';

function isEnumSpec(v) {
  return v && typeof v === 'object' && !Array.isArray(v);
}

export default function FieldMapEditor({ fieldMap = {}, onChange }) {
  const { t } = useTranslation();
  const entries = Object.entries(fieldMap);

  const update = (next) => onChange(next);

  const addField = () => update({ ...fieldMap, ['canonical_field']: '' });

  const removeField = (key) => {
    const next = { ...fieldMap };
    delete next[key];
    update(next);
  };

  const renameField = (oldKey, newKey) => {
    if (!newKey || newKey === oldKey) return;
    const next = {};
    for (const [k, v] of Object.entries(fieldMap)) next[k === oldKey ? newKey : k] = v;
    update(next);
  };

  const setSimple = (key, column) => update({ ...fieldMap, [key]: column });

  const setMode = (key, mode) => {
    const current = fieldMap[key];
    if (mode === 'enum') {
      const from = isEnumSpec(current) ? current.from : (typeof current === 'string' ? current : '');
      update({ ...fieldMap, [key]: { from, values: isEnumSpec(current) ? current.values || {} : {} } });
    } else {
      const column = isEnumSpec(current) ? current.from : (typeof current === 'string' ? current : '');
      update({ ...fieldMap, [key]: column });
    }
  };

  const setEnumFrom = (key, from) => {
    const current = fieldMap[key];
    update({ ...fieldMap, [key]: { from, values: isEnumSpec(current) ? current.values || {} : {} } });
  };

  const setEnumValue = (key, rawValue, canonicalValue) => {
    const current = fieldMap[key];
    const values = { ...(isEnumSpec(current) ? current.values || {} : {}), [rawValue]: canonicalValue };
    update({ ...fieldMap, [key]: { from: isEnumSpec(current) ? current.from : '', values } });
  };

  const renameEnumValue = (key, oldRaw, newRaw) => {
    if (!newRaw || newRaw === oldRaw) return;
    const current = fieldMap[key];
    const values = { ...(isEnumSpec(current) ? current.values || {} : {}) };
    const v = values[oldRaw];
    delete values[oldRaw];
    values[newRaw] = v;
    update({ ...fieldMap, [key]: { from: isEnumSpec(current) ? current.from : '', values } });
  };

  const removeEnumValue = (key, rawValue) => {
    const current = fieldMap[key];
    const values = { ...(isEnumSpec(current) ? current.values || {} : {}) };
    delete values[rawValue];
    update({ ...fieldMap, [key]: { from: isEnumSpec(current) ? current.from : '', values } });
  };

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <label className="block text-[11px] text-dark-400 uppercase tracking-wider">
          {t('fieldMap.title')}
        </label>
        <button onClick={addField} className="text-xs text-blue-400 hover:text-blue-300">+ {t('nodeConfig.field')}</button>
      </div>
      <p className="text-[10px] text-dark-500 mb-2">
        {t('fieldMap.hint')}
      </p>

      {entries.length === 0 && <p className="text-xs text-dark-500 mb-2">{t('fieldMap.noMapping')}</p>}

      <div className="space-y-2">
        {entries.map(([key, value]) => {
          const enumMode = isEnumSpec(value);
          return (
            <div key={key} className="border border-dark-700 rounded-md p-2 bg-dark-950">
              <div className="flex items-center gap-1 mb-1">
                <input
                  defaultValue={key}
                  onBlur={(e) => renameField(key, e.target.value)}
                  className="input text-[11px] flex-1 font-mono"
                  placeholder="of_number"
                />
                <div className="flex text-[10px] rounded overflow-hidden border border-dark-600">
                  <button
                    onClick={() => setMode(key, 'simple')}
                    className={`px-1.5 py-1 ${!enumMode ? 'bg-blue-600 text-white' : 'bg-dark-800 text-dark-400'}`}
                  >
                    {t('fieldMap.simple')}
                  </button>
                  <button
                    onClick={() => setMode(key, 'enum')}
                    className={`px-1.5 py-1 ${enumMode ? 'bg-blue-600 text-white' : 'bg-dark-800 text-dark-400'}`}
                  >
                    {t('fieldMap.enum')}
                  </button>
                </div>
                <button onClick={() => removeField(key)} title={t('common.delete')} className="text-dark-500 hover:text-red-400 px-1">×</button>
              </div>

              {!enumMode ? (
                <input
                  value={typeof value === 'string' ? value : ''}
                  placeholder={t('fieldMap.sourceColumnEx1')}
                  onChange={(e) => setSimple(key, e.target.value)}
                  className="input text-xs font-mono"
                />
              ) : (
                <div className="pl-2 border-l-2 border-dark-700">
                  <input
                    value={value.from || ''}
                    placeholder={t('fieldMap.sourceColumnEx2')}
                    onChange={(e) => setEnumFrom(key, e.target.value)}
                    className="input text-xs font-mono mb-1"
                  />
                  <div className="space-y-1">
                    {Object.entries(value.values || {}).map(([raw, canonical]) => (
                      <div key={raw} className="flex items-center gap-1">
                        <input
                          defaultValue={raw}
                          onBlur={(e) => renameEnumValue(key, raw, e.target.value)}
                          className="input text-[11px] w-16 font-mono"
                        />
                        <span className="text-dark-500 text-xs">→</span>
                        <input
                          value={canonical}
                          onChange={(e) => setEnumValue(key, raw, e.target.value)}
                          className="input text-[11px] flex-1 font-mono"
                          placeholder="RUNNING"
                        />
                        <button onClick={() => removeEnumValue(key, raw)} className="text-dark-500 hover:text-red-400 px-1">×</button>
                      </div>
                    ))}
                    <button
                      onClick={() => setEnumValue(key, 'new_value', '')}
                      className="text-[11px] text-blue-400 hover:text-blue-300"
                    >
                      + {t('fieldMap.addValue')}
                    </button>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
