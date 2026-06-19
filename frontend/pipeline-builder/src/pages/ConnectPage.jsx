import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchConnectors } from '../api/client';
import { CONNECTOR_TEMPLATES, triggerTypeFor } from '../lib/connectorTemplates';
import { useStudioStore } from '../store/studioStore';

// Step 1 of the workflow: pick a connector. Selecting one stores it and jumps to
// Compose, where it's applied to the trigger (first) node.
export default function ConnectPage() {
  const [connectors, setConnectors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const selectConnector = useStudioStore((s) => s.selectConnector);
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        const data = await fetchConnectors();
        setConnectors(data.connectors || []);
        setError(null);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const handleSelect = (c) => {
    selectConnector(c);
    navigate('/compose');
  };

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-xl font-semibold text-white mb-1">🔌 Choisir un connecteur</h2>
        <p className="text-dark-400 text-sm mb-6">
          Sélectionnez une source de données. Elle devient le nœud d'entrée (trigger) de votre pipeline.
        </p>

        {loading && <p className="text-dark-400 text-sm">Chargement…</p>}
        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm">
            ❌ {error} — le serveur API tourne-t-il sur :8080 ?
          </div>
        )}

        <div className="grid gap-4 sm:grid-cols-2">
          {connectors.map((c) => {
            const template = CONNECTOR_TEMPLATES[c.name] || {};
            const fields = Object.entries(template);
            return (
              <div key={c.name} className="border border-blue-500/30 bg-blue-500/5 rounded-lg p-4 flex flex-col">
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-xl">🔌</span>
                  <h3 className="font-medium text-white">{c.name}</h3>
                  <span className="ml-auto text-[10px] px-2 py-0.5 rounded bg-blue-500/15 text-blue-300 font-mono">
                    {triggerTypeFor(c.name)}
                  </span>
                </div>
                <p className="text-sm text-dark-400 mb-3">{c.description || '—'}</p>

                <div className="text-[11px] text-dark-500 uppercase tracking-wider mb-1">Config par défaut</div>
                {fields.length === 0 ? (
                  <p className="text-xs text-dark-500 mb-3">Aucun paramètre.</p>
                ) : (
                  <div className="space-y-1 mb-3">
                    {fields.map(([k, v]) => (
                      <div key={k} className="flex justify-between text-xs font-mono">
                        <span className="text-dark-300">{k}</span>
                        <span className="text-dark-400">{String(v)}</span>
                      </div>
                    ))}
                  </div>
                )}

                <button
                  onClick={() => handleSelect(c)}
                  className="mt-auto bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-1.5 rounded-md transition"
                >
                  Sélectionner →
                </button>
              </div>
            );
          })}
        </div>

        {!loading && !error && connectors.length === 0 && (
          <p className="text-dark-500 text-sm">Aucun connecteur enregistré.</p>
        )}
      </div>
    </div>
  );
}
