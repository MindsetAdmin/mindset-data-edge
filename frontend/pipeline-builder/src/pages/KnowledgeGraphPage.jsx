import { useState, useEffect, useMemo } from 'react';
import CytoscapeGraph from '../components/CytoscapeGraph';
import { fetchKnowledgeGraph } from '../api/client';
import { toElements, typesPresent, STYLESHEET, LAYOUTS, NODE_COLORS, FALLBACK_COLOR } from '../lib/kgGraph';

export default function KnowledgeGraphPage() {
  const [kind, setKind] = useState('technical');
  const [reloadKey, setReloadKey] = useState(0);
  const [graph, setGraph] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selected, setSelected] = useState(null);
  const [typeFilter, setTypeFilter] = useState('all');

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        setLoading(true);
        setSelected(null);
        setTypeFilter('all');
        const data = await fetchKnowledgeGraph(kind);
        if (!cancelled) {
          setGraph(data);
          setError(null);
        }
      } catch (err) {
        if (!cancelled) setError(err.message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [kind, reloadKey]);

  const filteredGraph = useMemo(() => {
    if (!graph || typeFilter === 'all') return graph;
    return { nodes: (graph.nodes || []).filter((n) => n.type === typeFilter), edges: graph.edges };
  }, [graph, typeFilter]);

  const elements = useMemo(() => toElements(kind, filteredGraph), [kind, filteredGraph]);
  const legend = useMemo(() => typesPresent(graph), [graph]);
  const layout = LAYOUTS[kind];

  const nodeCount = graph?.nodes?.length || 0;
  const edgeCount = graph?.edges?.length || 0;

  return (
    <div className="h-full flex">
      {/* Sidebar */}
      <aside className="w-60 bg-dark-900 border-r border-dark-700 p-4 shrink-0 flex flex-col gap-4">
        <div>
          <h2 className="text-sm font-semibold text-white mb-2">🧠 Knowledge Graph</h2>
          <div className="flex rounded-md overflow-hidden border border-dark-700 text-xs">
            <button
              onClick={() => setKind('technical')}
              className={`flex-1 px-2 py-1.5 transition ${kind === 'technical' ? 'bg-blue-600 text-white' : 'bg-dark-800 text-dark-300 hover:bg-dark-700'}`}
            >
              Technique
            </button>
            <button
              onClick={() => setKind('domain')}
              className={`flex-1 px-2 py-1.5 transition ${kind === 'domain' ? 'bg-blue-600 text-white' : 'bg-dark-800 text-dark-300 hover:bg-dark-700'}`}
            >
              Domaine
            </button>
          </div>
          <p className="text-[11px] text-dark-500 mt-2">
            {kind === 'technical'
              ? 'Architecture : connexions, topics, pipelines, fonctions.'
              : 'Données : machines, micro-arrêts, causes, coûts.'}
          </p>
        </div>

        <div className="flex gap-2">
          <Stat label="Nœuds" value={nodeCount} />
          <Stat label="Relations" value={edgeCount} />
        </div>

        <div>
          <h3 className="text-[11px] text-dark-400 uppercase tracking-wider mb-2">Filtres</h3>
          <div className="space-y-1.5">
            <button
              onClick={() => setTypeFilter('all')}
              className={`flex items-center gap-2 text-xs w-full px-1.5 py-1 rounded transition ${
                typeFilter === 'all' ? 'bg-dark-700 text-white' : 'text-dark-300 hover:bg-dark-800'
              }`}
            >
              <span className="w-3 h-3 rounded-sm bg-dark-400 border border-dark-900" />
              Tous
            </button>
            {legend.map((t) => (
              <button
                key={t}
                onClick={() => setTypeFilter(t)}
                className={`flex items-center gap-2 text-xs w-full px-1.5 py-1 rounded transition ${
                  typeFilter === t ? 'bg-dark-700 text-white' : 'text-dark-300 hover:bg-dark-800'
                }`}
              >
                <span
                  className="w-3 h-3 rounded-sm border border-dark-900"
                  style={{ background: NODE_COLORS[t] || FALLBACK_COLOR }}
                />
                {t}
              </button>
            ))}
            {legend.length === 0 && <p className="text-xs text-dark-500">—</p>}
          </div>
        </div>

        <button
          onClick={() => setReloadKey((k) => k + 1)}
          className="mt-auto bg-dark-700 hover:bg-dark-600 text-white text-xs px-3 py-1.5 rounded-md transition"
        >
          🔄 Rafraîchir
        </button>
      </aside>

      {/* Graph */}
      <div className="flex-1 relative bg-dark-950 min-w-0">
        {loading && (
          <div className="absolute inset-0 flex items-center justify-center text-dark-400 text-sm z-10">
            Chargement du graphe…
          </div>
        )}
        {error && (
          <div className="absolute inset-0 flex items-center justify-center z-10">
            <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-4 text-red-400 text-sm max-w-sm text-center">
              ❌ {error}
              <p className="text-xs text-dark-400 mt-1">Le serveur API tourne-t-il sur :8080 ?</p>
            </div>
          </div>
        )}
        {!loading && !error && nodeCount === 0 && (
          <div className="absolute inset-0 flex items-center justify-center text-dark-500 text-sm z-10 text-center px-6">
            <div>
              <div className="text-3xl mb-2">📭</div>
              Graphe vide. {kind === 'domain' ? 'Lancez l’agent pour générer des événements.' : 'Aucun pipeline chargé.'}
            </div>
          </div>
        )}
        {nodeCount > 0 && (
          <CytoscapeGraph elements={elements} stylesheet={STYLESHEET} layout={layout} onNodeSelect={setSelected} />
        )}
      </div>

      {/* Details */}
      <aside className="w-64 bg-dark-900 border-l border-dark-700 p-4 shrink-0 overflow-y-auto">
        <h3 className="text-sm font-semibold text-blue-400 mb-3">📋 Détails</h3>
        {!selected ? (
          <p className="text-sm text-dark-500">Cliquez sur un nœud.</p>
        ) : (
          <div className="space-y-3 text-sm">
            <Detail label="ID" value={selected.id} mono />
            <Detail label="Type" value={selected.type} />
            <Detail label="Label" value={selected.label} />
            {selected.properties && Object.keys(selected.properties).length > 0 && (
              <div>
                <div className="text-[11px] text-dark-400 uppercase tracking-wider mb-1">Propriétés</div>
                <div className="space-y-1">
                  {Object.entries(selected.properties).map(([k, v]) => (
                    <div key={k} className="flex justify-between gap-2 text-xs">
                      <span className="text-dark-400">{k}</span>
                      <span className="text-dark-200 font-mono text-right break-all">
                        {typeof v === 'object' ? JSON.stringify(v) : String(v)}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </aside>
    </div>
  );
}

function Stat({ label, value }) {
  return (
    <div className="flex-1 bg-dark-950 border border-dark-700 rounded-md p-2 text-center">
      <div className="text-lg font-bold text-blue-400">{value}</div>
      <div className="text-[10px] text-dark-500 uppercase tracking-wider">{label}</div>
    </div>
  );
}

function Detail({ label, value, mono }) {
  return (
    <div>
      <div className="text-[11px] text-dark-400 uppercase tracking-wider mb-0.5">{label}</div>
      <div className={`text-dark-200 break-all ${mono ? 'font-mono text-xs' : ''}`}>{value}</div>
    </div>
  );
}
