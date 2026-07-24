import { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

// Generic searchable picker modal. Options: { value, label, sub, badge, group }.
// The user CHOOSES from available options (search + select) — never free-types.
export default function PickerModal({ title, options = [], onSelect, onClose, allowCustom = false, customLabel }) {
  const { t } = useTranslation();
  const [query, setQuery] = useState('');
  const [custom, setCustom] = useState('');
  const resolvedCustomLabel = customLabel || t('pickerModal.defaultCustomLabel');

  const filtered = useMemo(() => {
    const q = query.toLowerCase().trim();
    if (!q) return options;
    return options.filter((o) =>
      [o.label, o.value, o.sub, o.group].filter(Boolean).some((s) => String(s).toLowerCase().includes(q))
    );
  }, [options, query]);

  const groups = useMemo(() => {
    const m = new Map();
    filtered.forEach((o) => {
      const g = o.group || '';
      if (!m.has(g)) m.set(g, []);
      m.get(g).push(o);
    });
    return [...m.entries()];
  }, [filtered]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" onClick={onClose}>
      <div
        className="bg-dark-900 border border-dark-700 rounded-xl w-full max-w-lg max-h-[80vh] flex flex-col shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between px-4 py-3 border-b border-dark-700">
          <h3 className="text-sm font-semibold text-white">{title}</h3>
          <button onClick={onClose} className="text-dark-400 hover:text-white text-lg leading-none">×</button>
        </div>

        <div className="p-4 border-b border-dark-700">
          <input
            autoFocus
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={`🔍 ${t('common.search')}`}
            className="w-full bg-dark-950 border border-dark-700 rounded-md px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
          />
        </div>

        <div className="flex-1 overflow-y-auto p-3 space-y-3">
          {groups.length === 0 && <p className="text-sm text-dark-500 px-1">{t('pickerModal.noResults')}</p>}
          {groups.map(([group, items]) => (
            <div key={group}>
              {group && <div className="text-[11px] text-dark-400 uppercase tracking-wider mb-1 px-1">{group}</div>}
              <div className="space-y-1">
                {items.map((o) => (
                  <button
                    key={o.value}
                    onClick={() => onSelect(o)}
                    className="w-full text-left bg-dark-800 hover:bg-dark-700 border border-dark-700 rounded-md px-3 py-2 transition flex items-center gap-2"
                  >
                    <div className="min-w-0 flex-1">
                      <div className="text-sm text-white truncate">{o.label}</div>
                      {o.sub && <div className="text-[11px] text-dark-400 truncate">{o.sub}</div>}
                    </div>
                    {o.badge && (
                      <span className="text-[10px] px-2 py-0.5 rounded bg-blue-500/15 text-blue-300 font-mono shrink-0">
                        {o.badge}
                      </span>
                    )}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>

        {allowCustom && (
          <div className="p-3 border-t border-dark-700 flex gap-2">
            <input
              value={custom}
              onChange={(e) => setCustom(e.target.value)}
              placeholder={resolvedCustomLabel}
              className="flex-1 bg-dark-950 border border-dark-700 rounded-md px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
            />
            <button
              onClick={() => custom.trim() && onSelect({ value: custom.trim(), label: custom.trim() })}
              className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-2 rounded-md transition"
            >
              {t('pickerModal.use')}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
