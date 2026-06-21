import { useState, useEffect, useRef } from 'react';
import {
  ResponsiveContainer,
  ComposedChart,
  Bar,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts';
import { fetchStats, fetchKnowledgeGraph, fetchMachines, fetchConfig } from '../api/client';
import { buildEvents, effectiveCost, splitDays, deltaPct, paretoCauses } from '../lib/dashboardData';

const REFRESH_MS = 5000;

export default function DashboardPage() {
  const [stats, setStats] = useState(null);
  const [events, setEvents] = useState([]);
  const [machines, setMachines] = useState([]);
  const [config, setConfig] = useState(null);
  const [lastUpdate, setLastUpdate] = useState(null);
  const [error, setError] = useState(null);
  const timer = useRef(null);

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

  useEffect(() => {
    refresh();
    timer.current = setInterval(refresh, REFRESH_MS);
    return () => clearInterval(timer.current);
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

  const pareto = paretoCauses(events);
  const connected = stats?.broker_connected;

  return (
    <div className="h-full overflow-y-auto p-5">
      <div className="max-w-6xl mx-auto space-y-5">
        {/* Header */}
        <div className="bg-dark-900 border border-dark-700 rounded-lg px-4 py-3 flex flex-wrap items-center gap-x-8 gap-y-1 text-sm">
          <span className="font-semibold text-white">📊 MindSet Data — Dashboard</span>
          <span className="text-dark-400">Site : <span className="text-dark-200">{config?.site?.name || config?.site?.id || '—'}</span></span>
          <span className="text-dark-400">
            Statut : <span className={connected ? 'text-green-400' : 'text-red-400'}>{connected ? '🟢 Connecté' : '🔴 Déconnecté'}</span>
          </span>
          <span className="text-dark-400">Uptime : <span className="text-dark-200">{fmtDuration(stats?.uptime_seconds)}</span></span>
          <span className="text-dark-400 ml-auto">MàJ : {lastUpdate ? lastUpdate.toLocaleTimeString() : '—'}</span>
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

        {/* Pareto */}
        <Panel title="📈 Pareto des causes">
          {pareto.length === 0 ? (
            <Empty text="Aucun événement enregistré." />
          ) : (
            <ResponsiveContainer width="100%" height={240}>
              <ComposedChart data={pareto}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                <XAxis dataKey="cause" tick={{ fill: '#94a3b8', fontSize: 11 }} />
                <YAxis yAxisId="l" tick={{ fill: '#94a3b8', fontSize: 11 }} />
                <YAxis yAxisId="r" orientation="right" domain={[0, 100]} unit="%" tick={{ fill: '#94a3b8', fontSize: 11 }} />
                <Tooltip contentStyle={tooltipStyle} />
                <Bar yAxisId="l" dataKey="count" fill="#f59e0b" radius={[4, 4, 0, 0]} name="Occurrences" />
                <Line yAxisId="r" type="monotone" dataKey="cumulative" stroke="#38bdf8" strokeWidth={2} name="Cumulé %" />
              </ComposedChart>
            </ResponsiveContainer>
          )}
        </Panel>

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
            {machines.filter((m) => m.work_center !== '(autres)').length === 0 ? (
              <Empty text="Aucune machine découverte." />
            ) : (
              <div className="divide-y divide-dark-800">
                {machines
                  .filter((m) => m.work_center !== '(autres)')
                  .map((m) => {
                    const running = m.state?.running;
                    const temp = tagValue(m.tags, 'temperature');
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
          <Gantt machines={machines.filter((m) => m.work_center !== '(autres)')} />
        </Panel>
      </div>
    </div>
  );
}

// --- helpers / subcomponents ----------------------------------------------

const tooltipStyle = { background: '#0f172a', border: '1px solid #334155', borderRadius: 8, color: '#e2e8f0' };

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
