import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchStats } from '../api/client';

export default function OverviewPage() {
  const [stats, setStats] = useState(null);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchStats().then(setStats).catch((e) => setError(e.message));
  }, []);

  const cards = [
    { label: 'Pipelines', value: stats?.pipelines, icon: '🔧' },
    { label: 'Fonctions', value: stats?.functions, icon: '⚙️' },
  ];

  const steps = [
    { to: '/connect', icon: '🔌', title: 'Connect', desc: 'Choisir une source de données' },
    { to: '/compose', icon: '⚙️', title: 'Compose', desc: 'Construire le pipeline (drag & drop)' },
    { to: '/pipelines', icon: '📡', title: 'Pipelines', desc: 'Charger / exécuter un pipeline' },
    { to: '/dashboards', icon: '📊', title: 'Dashboards', desc: 'Voir les métriques temps réel' },
    { to: '/kg', icon: '🧠', title: 'Knowledge Graph', desc: 'Explorer les relations' },
  ];

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-2xl font-bold text-white mb-1">🧩 MindSet Data</h1>
        <p className="text-dark-400 text-sm mb-6">
          Construisez des pipelines industriels visuellement, exécutez-les, et explorez le graphe de connaissances.
        </p>

        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm mb-6">
            ❌ {error} — le serveur API tourne-t-il sur :8080 ?
          </div>
        )}

        <div className="grid grid-cols-2 gap-3 mb-8 max-w-md">
          {cards.map((c) => (
            <div key={c.label} className="bg-dark-900 border border-dark-700 rounded-lg p-4 text-center">
              <div className="text-2xl mb-1">{c.icon}</div>
              <div className="text-2xl font-bold text-blue-400">{c.value ?? '—'}</div>
              <div className="text-[11px] text-dark-500 uppercase tracking-wider">{c.label}</div>
            </div>
          ))}
        </div>

        <h2 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-3">Démarrer</h2>
        <div className="grid gap-3 sm:grid-cols-2">
          {steps.map((s) => (
            <button
              key={s.to}
              onClick={() => navigate(s.to)}
              className="text-left bg-dark-900 border border-dark-700 hover:border-blue-500/50 hover:bg-dark-800 rounded-lg p-4 transition flex items-center gap-3"
            >
              <span className="text-2xl">{s.icon}</span>
              <div>
                <div className="font-medium text-white">{s.title}</div>
                <div className="text-xs text-dark-400">{s.desc}</div>
              </div>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
