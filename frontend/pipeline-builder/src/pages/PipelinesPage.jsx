import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchPipelines, fetchExamplePipelines, runPipeline } from '../api/client';
import { useStudioStore } from '../store/studioStore';

// Two groups: the pipelines you built (loaded into the engine + KG) and the
// shipped example templates (load one into Compose to start from it).
export default function PipelinesPage() {
  const [mine, setMine] = useState([]);
  const [examples, setExamples] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [runStatus, setRunStatus] = useState(null);
  const requestLoadPipeline = useStudioStore((s) => s.requestLoadPipeline);
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        const [u, ex] = await Promise.all([fetchPipelines(), fetchExamplePipelines().catch(() => ({ pipelines: [] }))]);
        setMine(u.pipelines || []);
        setExamples(ex.pipelines || []);
        setError(null);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const handleLoad = (p) => {
    requestLoadPipeline(p); // pass the full object so examples load too
    navigate('/compose');
  };

  const handleRun = async (id) => {
    setRunStatus({ id, msg: 'Exécution…', type: 'pending' });
    try {
      const res = await runPipeline(id);
      const ok = res.status === 'success';
      const okCount = (res.nodes || []).filter((n) => n.status === 'success').length;
      setRunStatus({ id, type: ok ? 'ok' : 'error', msg: `${res.status} ${okCount}/${(res.nodes || []).length} nœuds OK` });
    } catch (e) {
      setRunStatus({ id, type: 'error', msg: e.message });
    }
  };

  const Card = ({ p, runnable }) => (
    <div className="border border-dark-700 bg-dark-900 rounded-lg p-4">
      <div className="flex items-center gap-3">
        <span className="text-xl">🔧</span>
        <div className="min-w-0 flex-1">
          <div className="font-medium text-white">{p.name}</div>
          <div className="text-xs text-dark-400">{p.description}</div>
          <div className="text-[11px] text-dark-500 mt-0.5 font-mono">
            {(p.nodes || []).length} fonctions · trigger: {p.trigger?.function}
          </div>
        </div>
        <button onClick={() => handleLoad(p)} className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-1.5 rounded-md transition">
          Charger →
        </button>
        {runnable && (
          <button onClick={() => handleRun(p.id)} className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition">
            ▶️ Exécuter
          </button>
        )}
      </div>
      {runStatus && runStatus.id === p.id && (
        <div className={`text-xs mt-2 ${runStatus.type === 'ok' ? 'text-green-400' : runStatus.type === 'error' ? 'text-red-400' : 'text-dark-400'}`}>
          {runStatus.type === 'ok' ? '✅' : runStatus.type === 'error' ? '❌' : '⏳'} {runStatus.msg}
        </div>
      )}
    </div>
  );

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-3xl mx-auto">
        <h2 className="text-xl font-semibold text-white mb-1">📡 Pipelines</h2>
        <p className="text-dark-400 text-sm mb-6">
          Vos pipelines tournent dans le moteur et apparaissent dans le Knowledge Graph. Les modèles
          sont des points de départ à charger dans Compose.
        </p>

        {loading && <p className="text-dark-400 text-sm">Chargement…</p>}
        {error && <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm">❌ {error}</div>}

        {/* Mine */}
        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">Mes pipelines</h3>
        <div className="space-y-3 mb-8">
          {mine.map((p) => <Card key={p.id} p={p} runnable />)}
          {!loading && mine.length === 0 && (
            <p className="text-dark-500 text-sm">Aucun pipeline créé. Construisez-en un dans Compose, ou chargez un modèle ci-dessous.</p>
          )}
        </div>

        {/* Examples */}
        {examples.length > 0 && (
          <>
            <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">📦 Modèles (exemples)</h3>
            <div className="space-y-3">
              {examples.map((p) => <Card key={`ex_${p.id}`} p={p} runnable={false} />)}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
