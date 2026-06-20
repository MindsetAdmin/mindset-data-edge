import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchPipelines, runPipeline } from '../api/client';
import { useStudioStore } from '../store/studioStore';

// Pre-defined pipelines: load one onto the canvas (1-click) or run it directly.
export default function PipelinesPage() {
  const [pipelines, setPipelines] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [runStatus, setRunStatus] = useState(null);
  const [query, setQuery] = useState('');
  const requestLoadPipeline = useStudioStore((s) => s.requestLoadPipeline);
  const navigate = useNavigate();

  const visible = pipelines.filter((p) => {
    const q = query.toLowerCase().trim();
    if (!q) return true;
    return [p.name, p.id, p.description].filter(Boolean).some((s) => s.toLowerCase().includes(q));
  });

  useEffect(() => {
    (async () => {
      try {
        const data = await fetchPipelines();
        setPipelines(data.pipelines || []);
        setError(null);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const handleLoad = (id) => {
    requestLoadPipeline(id);
    navigate('/compose');
  };

  const handleRun = async (id) => {
    setRunStatus({ id, msg: 'Exécution…', type: 'pending' });
    try {
      const res = await runPipeline(id);
      const ok = res.status === 'success';
      const okCount = (res.nodes || []).filter((n) => n.status === 'success').length;
      setRunStatus({
        id,
        type: ok ? 'ok' : 'error',
        msg: `${res.status} — ${okCount}/${(res.nodes || []).length} nœuds OK`,
      });
    } catch (e) {
      setRunStatus({ id, type: 'error', msg: e.message });
    }
  };

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-3xl mx-auto">
        <h2 className="text-xl font-semibold text-white mb-1">📡 Pipelines pré-définis</h2>
        <p className="text-dark-400 text-sm mb-6">Chargez un pipeline en un clic, ou exécutez-le directement.</p>

        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="🔍 Rechercher une pipeline…"
          className="w-full bg-dark-950 border border-dark-700 rounded-md px-3 py-2 text-sm text-white mb-4 focus:outline-none focus:border-blue-500"
        />

        {loading && <p className="text-dark-400 text-sm">Chargement…</p>}
        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm">❌ {error}</div>
        )}

        <div className="space-y-3">
          {visible.map((p) => (
            <div key={p.id} className="border border-dark-700 bg-dark-900 rounded-lg p-4">
              <div className="flex items-center gap-3">
                <span className="text-xl">🔧</span>
                <div className="min-w-0 flex-1">
                  <div className="font-medium text-white">{p.name}</div>
                  <div className="text-xs text-dark-400">{p.description || '—'}</div>
                  <div className="text-[11px] text-dark-500 mt-0.5 font-mono">
                    {(p.nodes || []).length} fonctions · trigger: {p.trigger?.function || '—'}
                  </div>
                </div>
                <button
                  onClick={() => handleLoad(p.id)}
                  className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-1.5 rounded-md transition"
                >
                  Charger →
                </button>
                <button
                  onClick={() => handleRun(p.id)}
                  className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition"
                >
                  ▶️ Exécuter
                </button>
              </div>
              {runStatus && runStatus.id === p.id && (
                <div
                  className={`text-xs mt-2 ${
                    runStatus.type === 'ok' ? 'text-green-400' : runStatus.type === 'error' ? 'text-red-400' : 'text-dark-400'
                  }`}
                >
                  {runStatus.type === 'ok' ? '✅' : runStatus.type === 'error' ? '❌' : '⏳'} {runStatus.msg}
                </div>
              )}
            </div>
          ))}
        </div>

        {!loading && !error && pipelines.length === 0 && (
          <p className="text-dark-500 text-sm">Aucun pipeline. Créez-en un dans Compose.</p>
        )}
      </div>
    </div>
  );
}
