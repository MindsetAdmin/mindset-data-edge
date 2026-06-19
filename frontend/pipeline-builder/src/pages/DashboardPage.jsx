import { useState, useEffect } from 'react';
import { fetchStats, fetchKnowledgeGraph } from '../api/client';

export default function DashboardPage() {
  const [stats, setStats] = useState(null);
  const [events, setEvents] = useState([]);
  const [error, setError] = useState(null);
  const [reload, setReload] = useState(0);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [s, g] = await Promise.all([fetchStats(), fetchKnowledgeGraph('domain')]);
        if (cancelled) return;
        setStats(s);
        const evs = (g.nodes || [])
          .filter((n) => n.type === 'Event')
          .map((n) => ({
            id: n.id,
            workCenter: n.properties?.work_center || '—',
            duration: n.properties?.duration_seconds,
            createdAt: n.created_at,
          }))
          .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
        setEvents(evs);
        setError(null);
      } catch (e) {
        if (!cancelled) setError(e.message);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [reload]);

  const downtime = stats?.total_downtime_seconds || 0;
  const cards = [
    { label: 'Micro-arrêts', value: stats?.micro_stops ?? '—', icon: '🔴' },
    { label: 'Downtime total', value: stats ? `${downtime.toFixed(0)} s` : '—', sub: stats ? `${(downtime / 60).toFixed(1)} min` : '', icon: '⏱️' },
    { label: 'Coût estimé', value: stats ? `${(stats.estimated_cost_eur || 0).toFixed(2)} €` : '—', sub: stats ? `@ ${stats.hourly_cost} €/h` : '', icon: '💰' },
    { label: 'Pipelines', value: stats?.pipelines ?? '—', icon: '🔧' },
  ];

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-xl font-semibold text-white">📊 Dashboard</h2>
            <p className="text-dark-400 text-sm">Métriques temps réel des micro-arrêts.</p>
          </div>
          <button
            onClick={() => setReload((r) => r + 1)}
            className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition"
          >
            🔄 Rafraîchir
          </button>
        </div>

        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm mb-6">
            ❌ {error} — le serveur API tourne-t-il sur :8080 ?
          </div>
        )}

        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-8">
          {cards.map((c) => (
            <div key={c.label} className="bg-dark-900 border border-dark-700 rounded-lg p-4">
              <div className="text-xl mb-1">{c.icon}</div>
              <div className="text-2xl font-bold text-blue-400">{c.value}</div>
              <div className="text-[11px] text-dark-500 uppercase tracking-wider">{c.label}</div>
              {c.sub && <div className="text-[11px] text-dark-500 mt-0.5">{c.sub}</div>}
            </div>
          ))}
        </div>

        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-3">Événements récents</h3>
        <div className="bg-dark-900 border border-dark-700 rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-dark-800 text-dark-400 text-xs uppercase tracking-wider">
              <tr>
                <th className="text-left px-4 py-2 font-medium">Machine</th>
                <th className="text-right px-4 py-2 font-medium">Durée</th>
                <th className="text-right px-4 py-2 font-medium">Date</th>
              </tr>
            </thead>
            <tbody>
              {events.map((e) => (
                <tr key={e.id} className="border-t border-dark-800">
                  <td className="px-4 py-2 text-white">{e.workCenter}</td>
                  <td className="px-4 py-2 text-right text-dark-200 font-mono">
                    {e.duration != null ? `${Number(e.duration).toFixed(1)} s` : '—'}
                  </td>
                  <td className="px-4 py-2 text-right text-dark-400">
                    {e.createdAt ? new Date(e.createdAt).toLocaleString() : '—'}
                  </td>
                </tr>
              ))}
              {events.length === 0 && (
                <tr>
                  <td colSpan={3} className="px-4 py-6 text-center text-dark-500">
                    Aucun événement enregistré.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
