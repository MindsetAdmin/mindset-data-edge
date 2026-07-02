import { useState, useEffect, useRef } from 'react';
import { fetchStats, fetchKnowledgeGraph, fetchMachines, fetchConfig } from '../api/client';
import { buildEvents, effectiveCost, splitDays, deltaPct } from '../lib/dashboardData';
import { useLiveSocket } from '../lib/useLiveSocket';
import LiveDataPanel from '../components/LiveDataPanel';
import DashboardWidgets from '../components/DashboardWidgets';
import { useStudioStore } from '../store/studioStore';

const FALLBACK_MS = 20000; // safety heartbeat; real-time comes from the WebSocket

export default function DashboardPage() {
  const [stats, setStats] = useState(null);
  const [events, setEvents] = useState([]);
  const [machines, setMachines] = useState([]);
  const [config, setConfig] = useState(null);
  const [lastUpdate, setLastUpdate] = useState(null);
  const [error, setError] = useState(null);
  const refreshRef = useRef(null);
  const debounceRef = useRef(null);

  async function refresh() {
    try {
      const [s, domain, m, c] = await Promise.all([
        fetchStats(),
        fetchKnowledgeGraph('domain'),
        fetchMachines(),
        fetchConfig(),
      ]);
      setStats(s);
      setEvents(buildEvents(domain));
      setMachines(m.machines || []);
      setConfig(c);
      setLastUpdate(new Date());
      setError(null);
    } catch (e) {
      setError(e.message);
    }
  }
  refreshRef.current = refresh;

  // Debounce so a burst of WS messages triggers a single re-fetch.
  const scheduleRefresh = () => {
    clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => refreshRef.current(), 500);
  };

  // Real-time push: any live message triggers a (debounced) refresh.
  const connected = useLiveSocket((msg) => {
    if (msg.type === 'event' || msg.type === 'state' || msg.type === 'tag') scheduleRefresh();
  });

  useEffect(() => {
    refresh();
    const t = setInterval(() => refreshRef.current(), FALLBACK_MS);
    return () => {
      clearInterval(t);
      clearTimeout(debounceRef.current);
    };
  }, []);

  const hourly = stats?.hourly_cost || 85;
  const { today, yesterday } = splitDays(events);

  const stopsToday = today.length;
  const stopsYest = yesterday.length;
  const downtimeToday = today.reduce((s, e) => s + e.duration, 0);
  const downtimeYest = yesterday.reduce((s, e) => s + e.duration, 0);
  const costToday = today.reduce((s, e) => s + effectiveCost(e, hourly), 0);
  const costYest = yesterday.reduce((s, e) => s + effectiveCost(e, hourly), 0);

  // Availability-based TRS proxy over an 8h shift (true OEE needs perf+quality data).
  const SHIFT = 8 * 3600;
  const availability = Math.max(0, Math.min(1, 1 - downtimeToday / SHIFT)) * 100;

  const brokerConnected = stats?.broker_connected;

  // Show only machines wired into a state_machine node in the current session's pipeline.
  const selectedMachines = useStudioStore((s) => s.selectedMachines);
  const filteredMachines = machines
    .filter((m) => m.work_center !== '(autres)')
    .filter((m) => !selectedMachines.length || selectedMachines.includes(m.work_center));

  return (
    <div className="h-full overflow-y-auto p-5">
      <div className="max-w-6xl mx-auto space-y-5">
        {/* Header */}
        <div className="bg-dark-900 border border-dark-700 rounded-lg px-4 py-3 flex flex-wrap items-center gap-x-8 gap-y-1 text-sm">
          
          <span className="text-dark-400">Site : <span className="text-dark-200">{config?.site?.name || config?.site?.id || '—'}</span></span>
          <span className="text-dark-400">
            Statut : <span className={brokerConnected ? 'text-green-400' : 'text-red-400'}>{brokerConnected ? '🟢 Connecté' : '🔴 Déconnecté'}</span>
          </span>
          <span className="text-dark-400">Uptime : <span className="text-dark-200">{fmtDuration(stats?.uptime_seconds)}</span></span>
          <span className={`ml-auto inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-full ${connected ? 'bg-green-500/15 text-green-400' : 'bg-dark-700 text-dark-400'}`}>
            <span className={`w-1.5 h-1.5 rounded-full ${connected ? 'bg-green-400 animate-pulse' : 'bg-dark-500'}`} />
            {connected ? 'LIVE (WebSocket)' : 'hors-ligne'}
          </span>
          <span className="text-dark-400">MàJ : {lastUpdate ? lastUpdate.toLocaleTimeString() : '—'}</span>
        </div>

        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm">
            ❌ {error} — le serveur API tourne-t-il sur :8080 ?
          </div>
        )}

        {/* KPI cards */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
          <Kpi icon="🔴" label="Micro-stops" value={stopsToday} unit="aujourd'hui" delta={deltaPct(stopsToday, stopsYest)} invert />
          <Kpi icon="💰" label="Coût total" value={`${costToday.toFixed(2)} €`} unit="aujourd'hui" delta={deltaPct(costToday, costYest)} invert />
          <Kpi icon="⏱️" label="Temps perdu" value={fmtDuration(downtimeToday)} unit="aujourd'hui" delta={deltaPct(downtimeToday, downtimeYest)} invert />
          <Kpi icon="📈" label="Disponibilité" value={`${availability.toFixed(1)}%`} unit="TRS (dispo.)" />
        </div>

        {/* Widgets pinned via the add_to_dashboard function */}
        <DashboardWidgets />

        {/* Live, user-selected tag data */}
        <LiveDataPanel />

        {/* Events + Machines */}
        <div className="grid lg:grid-cols-2 gap-5">
          <Panel title="📋 Derniers événements">
            {events.length === 0 ? (
              <Empty text="Aucun événement." />
            ) : (
              <div className="divide-y divide-dark-800">
                {events.slice(0, 8).map((e) => (
                  <div key={e.id} className="py-2 flex items-center gap-3 text-sm">
                    <span className="text-dark-500 font-mono text-xs w-16">
                      {e.createdAt ? new Date(e.createdAt).toLocaleTimeString().slice(0, 5) : '—'}
                    </span>
                    <span className="text-white w-24 truncate">{e.workCenter}</span>
                    <span className="text-amber-400 text-xs">Micro-stop {e.duration.toFixed(0)}s</span>
                    <span className="text-dark-400 text-xs ml-auto">{e.cause || '—'}</span>
                    <span className="text-green-400 text-xs font-mono w-20 text-right">{effectiveCost(e, hourly).toFixed(2)} €</span>
                  </div>
                ))}
              </div>
            )}
          </Panel>

          <Panel title="🏭 Statut machines">
            {filteredMachines.length === 0 ? (
              <Empty text={selectedMachines.length === 0
                ? "Configurez une machine dans un pipeline (onglet Compose) pour la voir ici."
                : "Aucune machine sélectionnée active."} />
            ) : (
              <div className="divide-y divide-dark-800">
                {filteredMachines.map((m) => {
                    const running = m.state?.running;
                    const temp = tagValue(m.tags);
                    return (
                      <div key={m.work_center} className="py-2 flex items-center gap-3 text-sm">
                        <span className="text-white w-28 truncate">{m.work_center}</span>
                        <span className={running == null ? 'text-dark-500' : running ? 'text-green-400' : 'text-red-400'}>
                          {running == null ? '⚪ n/a' : running ? '🟢 Running' : '🔴 Stopped'}
                        </span>
                        <span className="text-dark-400 text-xs ml-auto">
                          {temp != null ? `${temp}°C` : `${m.tags.length} tags`}
                        </span>
                      </div>
                    );
                  })}
              </div>
            )}
          </Panel>
        </div>

        {/* Gantt */}
        <Panel title="📊 Timeline machines">
          <Gantt machines={filteredMachines} />
        </Panel>
      </div>
    </div>
  );
}

// --- helpers / subcomponents ----------------------------------------------


function fmtDuration(seconds) {
  if (seconds == null) return '—';
  const s = Math.round(seconds);
  if (s < 60) return `${s}s`;
  const m = Math.floor(s / 60);
  if (m < 60) return `${m}min ${s % 60}s`;
  const h = Math.floor(m / 60);
  return `${h}h ${m % 60}min`;
}

function tagValue(tags, contains) {
  const t = (tags || []).find((x) => (x.name || '').toLowerCase().includes(contains));
  return t ? t.value : null;
}

function Kpi({ icon, label, value, unit, delta, invert }) {
  let deltaEl = null;
  if (delta != null && isFinite(delta)) {
    const up = delta > 0;
    // For "bad" metrics (stops/cost/downtime), up is red; for good metrics up is green.
    const good = invert ? !up : up;
    deltaEl = (
      <div className={`text-[11px] ${good ? 'text-green-400' : 'text-red-400'}`}>
        {up ? '▲' : '▼'} {Math.abs(delta).toFixed(0)}% vs hier
      </div>
    );
  }
  return (
    <div className="bg-dark-900 border border-dark-700 rounded-lg p-4">
      <div className="text-xl mb-1">{icon}</div>
      <div className="text-2xl font-bold text-blue-400">{value}</div>
      <div className="text-[11px] text-dark-500 uppercase tracking-wider">{label}</div>
      <div className="text-[11px] text-dark-500">{unit}</div>
      {deltaEl}
    </div>
  );
}

function Panel({ title, children }) {
  return (
    <div className="bg-dark-900 border border-dark-700 rounded-lg p-4">
      <h3 className="text-sm font-medium text-dark-300 mb-3">{title}</h3>
      {children}
    </div>
  );
}

function Empty({ text }) {
  return <p className="text-sm text-dark-500 py-6 text-center">{text}</p>;
}

// Gantt: reconstructs Running/Stopped segments from each machine's state history.
function Gantt({ machines }) {
  const now = Date.now();
  const withState = machines.filter((m) => m.state);
  if (withState.length === 0) {
    return <Empty text="Pas d'historique de transitions (l'agent doit tourner)." />;
  }

  // Window start = earliest transition, or 1h ago.
  let start = now - 3600 * 1000;
  withState.forEach((m) =>
    (m.state.history || []).forEach((h) => {
      const t = new Date(h.at).getTime();
      if (t < start) start = t;
    })
  );
  const span = Math.max(now - start, 1);

  return (
    <div className="space-y-2">
      {withState.map((m) => {
        const segs = buildSegments(m.state, start, now);
        return (
          <div key={m.work_center} className="flex items-center gap-2">
            <span className="text-xs text-dark-300 w-24 truncate">{m.work_center}</span>
            <div className="flex-1 h-5 rounded overflow-hidden flex bg-dark-950">
              {segs.map((s, i) => (
                <div
                  key={i}
                  className={s.running ? 'bg-green-500/70' : 'bg-red-500/70'}
                  style={{ width: `${((s.to - s.from) / span) * 100}%` }}
                  title={`${s.running ? 'Running' : 'Stopped'} ${Math.round((s.to - s.from) / 1000)}s`}
                />
              ))}
            </div>
          </div>
        );
      })}
      <div className="flex gap-4 text-[11px] text-dark-400 pt-1">
        <span>🟢 Running</span>
        <span>🔴 Stopped</span>
        <span className="ml-auto">{new Date(start).toLocaleTimeString().slice(0, 5)} → {new Date(now).toLocaleTimeString().slice(0, 5)}</span>
      </div>
    </div>
  );
}

function buildSegments(state, start, now) {
  const history = (state.history || []).slice().sort((a, b) => new Date(a.at) - new Date(b.at));
  const segs = [];
  let t0 = start;
  let cur = history.length ? history[0].from : state.running;
  history.forEach((h) => {
    const at = new Date(h.at).getTime();
    if (at > t0) {
      segs.push({ from: t0, to: at, running: cur });
      t0 = at;
    }
    cur = h.to;
  });
  segs.push({ from: t0, to: now, running: cur });
  return segs.filter((s) => s.to > s.from);
}
