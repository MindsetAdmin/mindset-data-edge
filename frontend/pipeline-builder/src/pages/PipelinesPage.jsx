import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { fetchPipelines, fetchExamplePipelines, runPipeline, deletePipeline } from '../api/client';
import { useStudioStore } from '../store/studioStore';

// Two groups: the pipelines you built (loaded into the engine + KG) and the
// shipped example templates (load one into Compose to start from it).
export default function PipelinesPage() {
  const { t } = useTranslation();
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
    setRunStatus({ id, msg: t('pipelines.running'), type: 'pending' });
    try {
      const res = await runPipeline(id);
      const ok = res.status === 'success';
      const okCount = (res.nodes || []).filter((n) => n.status === 'success').length;
      setRunStatus({ id, type: ok ? 'ok' : 'error', msg: t('pipelines.nodesOk', { status: res.status, ok: okCount, total: (res.nodes || []).length }) });
    } catch (e) {
      setRunStatus({ id, type: 'error', msg: e.message });
    }
  };

  const handleDelete = async (id) => {
    try {
      await deletePipeline(id);
      setMine((list) => list.filter((p) => p.id !== id));
    } catch (e) {
      setError(e.message);
    }
  };

  // deletable only applies to "Mes pipelines" — the shipped examples are
  // read-only starting points, not files this page should be able to remove.
  const Card = ({ p, runnable, deletable }) => (
    <div className="border border-dark-700 bg-dark-900 rounded-lg p-4">
      <div className="flex items-center gap-3">
        <span className="text-xl">🔧</span>
        <div className="min-w-0 flex-1">
          <div className="font-medium text-white">{p.name}</div>
          <div className="text-xs text-dark-400">{p.description}</div>
          <div className="text-[11px] text-dark-500 mt-0.5 font-mono">
            {t('pipelines.functionCount', { count: (p.nodes || []).length })} · trigger: {p.trigger?.function}
          </div>
        </div>
        <button onClick={() => handleLoad(p)} className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-1.5 rounded-md transition">
          {t('pipelines.load')} →
        </button>
        {runnable && (
          <button onClick={() => handleRun(p.id)} className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition">
            ▶️ {t('common.run')}
          </button>
        )}
        {deletable && (
          <button
            onClick={() => handleDelete(p.id)}
            className="bg-dark-700 hover:bg-red-500/30 text-dark-300 hover:text-red-300 text-sm px-3 py-1.5 rounded-md transition"
          >
            {t('common.delete')}
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
          {t('pipelines.subtitle')}
        </p>

        {loading && <p className="text-dark-400 text-sm">{t('common.loading')}</p>}
        {error && <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm">❌ {error}</div>}

        {/* Mine */}
        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">{t('pipelines.mine')}</h3>
        <div className="space-y-3 mb-8">
          {mine.map((p) => <Card key={p.id} p={p} runnable deletable />)}
          {!loading && mine.length === 0 && (
            <p className="text-dark-500 text-sm">{t('pipelines.noneCreated')}</p>
          )}
        </div>

        {/* Examples */}
        {examples.length > 0 && (
          <>
            <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">📦 {t('pipelines.templates')}</h3>
            <div className="space-y-3">
              {examples.map((p) => <Card key={`ex_${p.id}`} p={p} runnable={false} />)}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
