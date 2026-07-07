import { useState, useEffect, useMemo } from 'react';
import { RefreshCw, Network, Inbox, AlertCircle } from 'lucide-react';
import ForceGraph, { typesPresent, NODE_COLORS, FALLBACK_COLOR } from '../components/ForceGraph';
import { fetchKG } from '../api/client';

// Unified KG page (2026-07-02 merge). Single Cytoscape view over the merged
// graph. Category filter (business / platform / all) replaces the old
// Technique/Domaine toggle.
const CATEGORIES = [
    { value: 'all', label: 'All' },
    { value: 'business', label: 'Business' },
    { value: 'platform', label: 'Platform' },
];

export default function KnowledgeGraphPage() {
    const [category, setCategory] = useState('all');
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
                const data = await fetchKG(category);
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
    }, [category, reloadKey]);

    const filteredGraph = useMemo(() => {
        if (!graph || typeFilter === 'all') return graph;
        return { nodes: (graph.nodes || []).filter((n) => n.type === typeFilter), edges: graph.edges };
    }, [graph, typeFilter]);

    const legend = useMemo(() => typesPresent(graph), [graph]);

    const nodeCount = graph?.nodes?.length || 0;
    const edgeCount = graph?.edges?.length || 0;

    return (
        <div className="h-full flex bg-canvas text-text-primary">
            {/* Sidebar */}
            <aside className="w-60 bg-panel border-r border-border-subtle p-4 shrink-0 flex flex-col gap-4">
                <div>
                    <h2 className="inline-flex items-center gap-1.5 text-15 font-medium text-text-primary mb-2">
                        <Network size={14} strokeWidth={1.5} />
                        <span>Knowledge Graph</span>
                    </h2>
                    <div className="inline-flex w-full rounded border border-border-subtle overflow-hidden">
                        {CATEGORIES.map((c) => (
                            <button
                                key={c.value}
                                onClick={() => setCategory(c.value)}
                                className={`flex-1 px-2 py-1 text-11 border-r border-border-subtle last:border-r-0 transition-colors ${
                                    category === c.value
                                        ? 'bg-elevated text-text-primary'
                                        : 'bg-panel text-text-tertiary hover:text-text-primary hover:bg-panel-alt'
                                }`}
                            >
                                {c.label}
                            </button>
                        ))}
                    </div>
                    <p className="text-11 text-text-tertiary mt-2 leading-relaxed">
                        {category === 'business' && 'Site fingerprint — Equipment, Events, Causes, Costs.'}
                        {category === 'platform' && 'Platform topology — Pipelines, Functions, Topics, Connections, Dashboards.'}
                        {category === 'all' && 'Unified view — business events + platform wiring in one graph.'}
                    </p>
                </div>

                <div className="flex gap-2">
                    <Stat label="Nodes" value={nodeCount} />
                    <Stat label="Edges" value={edgeCount} />
                </div>

                <div>
                    <h3 className="text-11 text-text-secondary uppercase tracking-wide mb-2">Type filter</h3>
                    <div className="space-y-0.5">
                        <button
                            onClick={() => setTypeFilter('all')}
                            className={`flex items-center gap-2 text-13 w-full px-1.5 py-1 rounded transition-colors ${
                                typeFilter === 'all'
                                    ? 'bg-elevated text-text-primary'
                                    : 'text-text-secondary hover:text-text-primary hover:bg-panel-alt'
                            }`}
                        >
                            <span className="w-2.5 h-2.5 rounded-sm bg-text-tertiary" />
                            All
                        </button>
                        {legend.map((t) => (
                            <button
                                key={t}
                                onClick={() => setTypeFilter(t)}
                                className={`flex items-center gap-2 text-13 w-full px-1.5 py-1 rounded transition-colors ${
                                    typeFilter === t
                                        ? 'bg-elevated text-text-primary'
                                        : 'text-text-secondary hover:text-text-primary hover:bg-panel-alt'
                                }`}
                            >
                                <span
                                    className="w-2.5 h-2.5 rounded-sm"
                                    style={{ background: NODE_COLORS[t] || FALLBACK_COLOR }}
                                />
                                {t}
                            </button>
                        ))}
                        {legend.length === 0 && <p className="text-11 text-text-muted italic">—</p>}
                    </div>
                </div>

                <button
                    onClick={() => setReloadKey((k) => k + 1)}
                    className="mt-auto inline-flex items-center justify-center gap-1.5 bg-panel-alt hover:bg-elevated border border-border-subtle text-text-primary text-11 px-3 py-1.5 rounded transition-colors"
                >
                    <RefreshCw size={12} strokeWidth={1.5} />
                    <span>Refresh</span>
                </button>
            </aside>

            {/* Graph */}
            <div className="flex-1 relative bg-canvas min-w-0">
                {loading && (
                    <div className="absolute inset-0 flex items-center justify-center text-text-tertiary text-13 z-10 italic">
                        Loading graph…
                    </div>
                )}
                {error && (
                    <div className="absolute inset-0 flex items-center justify-center z-10">
                        <div className="bg-panel border border-status-stopped rounded p-4 text-status-stopped text-13 max-w-sm text-center inline-flex flex-col items-center gap-2">
                            <AlertCircle size={20} strokeWidth={1.5} />
                            <span>{error}</span>
                            <p className="text-11 text-text-tertiary mt-1">Is the API server running on :8080?</p>
                        </div>
                    </div>
                )}
                {!loading && !error && nodeCount === 0 && (
                    <div className="absolute inset-0 flex items-center justify-center z-10 text-center px-6">
                        <div className="flex flex-col items-center gap-2 text-text-tertiary">
                            <Inbox size={28} strokeWidth={1.5} />
                            <p className="text-13">
                                Empty graph.{' '}
                                {category === 'business'
                                    ? 'Start the agent to generate events.'
                                    : category === 'platform'
                                    ? 'No pipelines loaded.'
                                    : 'No data yet.'}
                            </p>
                        </div>
                    </div>
                )}
                {nodeCount > 0 && (
                    <ForceGraph graph={filteredGraph} onNodeSelect={setSelected} />
                )}
            </div>

            {/* Details */}
            <aside className="w-64 bg-panel border-l border-border-subtle p-4 shrink-0 overflow-y-auto">
                <h3 className="text-13 font-medium text-text-primary mb-3">Details</h3>
                {!selected ? (
                    <p className="text-13 text-text-tertiary italic">Click a node.</p>
                ) : (
                    <div className="space-y-3 text-13">
                        <Detail label="ID" value={selected.id} mono />
                        <Detail label="Category" value={selected.category || '—'} />
                        <Detail label="Type" value={selected.type} />
                        <Detail label="Label" value={selected.label} />
                        {selected.properties && Object.keys(selected.properties).length > 0 && (
                            <div>
                                <div className="text-11 text-text-secondary uppercase tracking-wide mb-1">Properties</div>
                                <div className="space-y-1">
                                    {Object.entries(selected.properties).map(([k, v]) => (
                                        <div key={k} className="flex justify-between gap-2 text-11">
                                            <span className="text-text-tertiary">{k}</span>
                                            <span className="mono text-text-primary text-right break-all">
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
        <div className="flex-1 bg-panel-alt border border-border-subtle rounded p-2 text-center">
            <div className="mono text-15 font-medium text-text-primary tabular">{value}</div>
            <div className="text-11 text-text-secondary uppercase tracking-wide">{label}</div>
        </div>
    );
}

function Detail({ label, value, mono }) {
    return (
        <div>
            <div className="text-11 text-text-secondary uppercase tracking-wide mb-0.5">{label}</div>
            <div className={`text-text-primary break-all ${mono ? 'mono text-11' : 'text-13'}`}>{value}</div>
        </div>
    );
}
