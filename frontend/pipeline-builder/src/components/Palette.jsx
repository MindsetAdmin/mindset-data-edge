import { useTranslation } from 'react-i18next';
import { typeStyle, getCategory } from '../lib/functionMeta';

// Left palette: draggable function cards. Dragging sets the function JSON on the
// dataTransfer; BuilderPage's onDrop reads it and creates a ReactFlow node.
export default function Palette({ functions, loading, error }) {
  const { t } = useTranslation();
  const onDragStart = (event, fn) => {
    event.dataTransfer.setData('application/reactflow', JSON.stringify(fn));
    event.dataTransfer.effectAllowed = 'move';
  };

  return (
    <aside className="w-72 bg-dark-900 border-r border-dark-700 p-4 overflow-y-auto shrink-0">
      <h2 className="text-xs font-semibold text-dark-400 uppercase tracking-wider mb-4">
        📦 {t('palette.title')}
      </h2>

      {loading && (
        <div className="text-center text-dark-400 py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4" />
          <p className="text-sm">{t('palette.loadingFunctions')}</p>
        </div>
      )}

      {error && (
        <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm">
          ❌ {error}
          <p className="text-xs text-dark-400 mt-1">{t('overview.serverDownHint')}</p>
        </div>
      )}

      {!loading && !error && groupByCategory(functions).map(([category, items]) => (
        <div key={category} className="mb-5">
          <h3 className="text-xs font-medium text-dark-300 mb-2">{category}</h3>
          <div className="space-y-2">
            {items.map((fn) => {
              const s = typeStyle(fn.type);
              return (
                <div
                  key={fn.name}
                  draggable
                  onDragStart={(e) => onDragStart(e, fn)}
                  className={`border rounded-lg p-2.5 cursor-grab active:cursor-grabbing hover:bg-dark-800 transition ${s.card}`}
                >
                  <div className="flex items-center gap-2">
                    <span className="text-lg">{s.icon}</span>
                    <div className="min-w-0">
                      <div className="font-medium text-sm text-white truncate">{fn.name}</div>
                      <div className="text-[11px] text-dark-400 truncate">{fn.description || '—'}</div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      ))}

      <p className="text-[11px] text-dark-600 mt-6 leading-relaxed">
        {t('palette.hint')}
      </p>
    </aside>
  );
}

function groupByCategory(functions) {
  const groups = new Map();
  (functions || []).forEach((fn) => {
    const cat = getCategory(fn.type);
    if (!groups.has(cat)) groups.set(cat, []);
    groups.get(cat).push(fn);
  });
  return Array.from(groups.entries());
}
