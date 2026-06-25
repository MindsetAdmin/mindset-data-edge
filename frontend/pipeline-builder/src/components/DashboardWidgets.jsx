import { useState, useEffect, useRef } from 'react';
import { ResponsiveContainer, LineChart, Line, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { useLiveSocket } from '../lib/useLiveSocket';
import { fetchDashboardPins } from '../api/client';

// Interactive dashboard widgets fed by add_to_dashboard pins (mindset/dashboard/#).
// Each pin label is a data source; the user adds widgets (line/bar/gauge/value/
// status), picks a time range, and sees live stats. Config persists in localStorage.

const STORE_KEY = 'mindset_widgets_v1';
const MAX_POINTS = 1000;
const RANGES = { '1m': 60e3, '5m': 300e3, '1h': 3600e3, '4h': 4 * 3600e3, '24h': 24 * 3600e3 };
const CHART_TYPES = [
  ['line', 'Ligne'],
  ['bar', 'Barres'],
  ['gauge', 'Jauge'],
  ['value', 'Valeur'],
  ['status', 'Statut'],
];
const tooltipStyle = { background: '#0f172a', border: '1px solid #334155', borderRadius: 8, color: '#e2e8f0' };

// Extract a numeric value from a pin payload (never expose raw JSON).
function toValue(pin) {
  const d = pin?.data;
  if (typeof d === 'number') return d;
  if (typeof d === 'boolean') return d ? 1 : 0;
  if (d && typeof d === 'object') {
    if (typeof d.value === 'number') return d.value;
    if (typeof d.value === 'boolean') return d.value ? 1 : 0;
    for (const k of ['total_cost_eur', 'cost_eur', 'duration_seconds', 'rate_per_sec', 'temperature', 'speed']) {
      if (typeof d[k] === 'number') return d[k];
    }
    const n = Number(d.value);
    if (Number.isFinite(n)) return n;
  }
  const n = Number(d);
  return Number.isFinite(n) ? n : null;
}
function isStatusPin(pin) {
  return pin?.kind === 'status' || typeof pin?.data === 'boolean' || typeof pin?.data?.value === 'boolean';
}

export default function DashboardWidgets() {
  const [history, setHistory] = useState({}); // label -> [{t, v, status, pin}]
  const [widgets, setWidgets] = useState(() => {
    try { return JSON.parse(localStorage.getItem(STORE_KEY)) || []; } catch { return []; }
  });
  const [showAdd, setShowAdd] = useState(false);
  const [now, setNow] = useState(Date.now());
  const histRef = useRef(history);
  histRef.current = history;

  // Persist widgets.
  useEffect(() => { localStorage.setItem(STORE_KEY, JSON.stringify(widgets)); }, [widgets]);

  // Tick so time-window views refresh even between messages.
  useEffect(() => { const t = setInterval(() => setNow(Date.now()), 2000); return () => clearInterval(t); }, []);

  // Seed from the server snapshot.
  useEffect(() => {
    fetchDashboardPins().then((d) => {
      const seed = {};
      (d.pins || []).forEach((p) => {
        if (!p?.label) return;
        seed[p.label] = [{ t: p.timestamp_ms || Date.now(), v: toValue(p), status: isStatusPin(p) ? toValue(p) > 0 : null, pin: p }];
      });
      setHistory((prev) => ({ ...seed, ...prev }));
    }).catch(() => {});
  }, []);

  // Live updates.
  const connected = useLiveSocket((msg) => {
    if (msg.type !== 'dashboard' || !msg.data?.label) return;
    const p = msg.data;
    setHistory((prev) => {
      const arr = (prev[p.label] || []).concat({
        t: p.timestamp_ms || Date.now(),
        v: toValue(p),
        status: isStatusPin(p) ? toValue(p) > 0 : null,
        pin: p,
      });
      return { ...prev, [p.label]: arr.length > MAX_POINTS ? arr.slice(arr.length - MAX_POINTS) : arr };
    });
  });

  const sources = Object.keys(history).sort();
  const addWidget = (w) => { setWidgets((ws) => [...ws, { id: Math.random().toString(36).slice(2, 8), range: '1h', ...w }]); setShowAdd(false); };
  const removeWidget = (id) => setWidgets((ws) => ws.filter((w) => w.id !== id));
  const updateWidget = (id, patch) => setWidgets((ws) => ws.map((w) => (w.id === id ? { ...w, ...patch } : w)));

  return (
    <div className="bg-dark-900 border border-dark-700 rounded-lg p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium text-dark-300">📌 Widgets épinglés</h3>
        <div className="flex items-center gap-3">
          <span className={`text-[11px] ${connected ? 'text-green-400' : 'text-dark-500'}`}>{connected ? '● live' : '○'}</span>
          <button onClick={() => setShowAdd(true)} className="text-xs bg-blue-600 hover:bg-blue-500 text-white px-2.5 py-1 rounded-md transition">+ Ajouter</button>
        </div>
      </div>

      {widgets.length === 0 ? (
        <p className="text-sm text-dark-500 py-8 text-center">
          Aucun widget. Cliquez « + Ajouter » pour visualiser une donnée publiée par <span className="font-mono text-dark-300">add_to_dashboard</span>.
        </p>
      ) : (
        <div className="grid sm:grid-cols-2 gap-3">
          {widgets.map((w) => (
            <Widget key={w.id} w={w} series={history[w.label] || []} now={now} onRemove={() => removeWidget(w.id)} onUpdate={(patch) => updateWidget(w.id, patch)} />
          ))}
        </div>
      )}

      {showAdd && <AddWidget sources={sources} onAdd={addWidget} onClose={() => setShowAdd(false)} />}
    </div>
  );
}

function stats(points) {
  if (points.length === 0) return null;
  const vals = points.map((p) => p.v).filter((v) => typeof v === 'number' && Number.isFinite(v));
  if (vals.length === 0) return { count: points.length };
  const min = Math.min(...vals), max = Math.max(...vals);
  const avg = vals.reduce((a, b) => a + b, 0) / vals.length;
  return { min, max, avg, last: vals[vals.length - 1], count: points.length };
}

function Widget({ w, series, now, onRemove, onUpdate }) {
  const [cfg, setCfg] = useState(false);
  const cutoff = now - (RANGES[w.range] || RANGES['1h']);
  const pts = series.filter((p) => p.t >= cutoff);
  const data = pts.map((p) => ({ t: p.t, time: new Date(p.t).toLocaleTimeString().slice(0, 8), value: p.v }));
  const s = stats(pts);
  const lastPin = series[series.length - 1]?.pin;

  return (
    <div className="bg-dark-950 border border-dark-700 rounded-lg p-3">
      <div className="flex items-center justify-between mb-2">
        <div className="text-sm font-medium text-white truncate">{icon(w.type)} {w.label}</div>
        <div className="flex items-center gap-1">
          <button onClick={() => setCfg((c) => !c)} title="Configurer" className="text-dark-400 hover:text-white text-xs px-1">⚙️</button>
          <button onClick={onRemove} title="Fermer" className="text-dark-400 hover:text-red-400 text-sm px-1">✕</button>
        </div>
      </div>

      {cfg && (
        <div className="flex flex-wrap gap-2 mb-2 text-[11px]">
          <select value={w.type} onChange={(e) => onUpdate({ type: e.target.value })} className="bg-dark-900 border border-dark-700 rounded px-1.5 py-1 text-dark-200">
            {CHART_TYPES.map(([v, l]) => <option key={v} value={v}>{l}</option>)}
          </select>
          <select value={w.range} onChange={(e) => onUpdate({ range: e.target.value })} className="bg-dark-900 border border-dark-700 rounded px-1.5 py-1 text-dark-200">
            {Object.keys(RANGES).map((r) => <option key={r} value={r}>{r}</option>)}
          </select>
        </div>
      )}

      <WidgetBody type={w.type} data={data} stats={s} lastPin={lastPin} />

      {s && w.type !== 'status' && (
        <div className="flex flex-wrap gap-x-3 gap-y-0.5 mt-2 text-[10px] text-dark-500">
          {s.last != null && <span>Dernier: <span className="text-dark-300">{fmt(s.last)}</span></span>}
          {s.min != null && <span>Min: {fmt(s.min)}</span>}
          {s.max != null && <span>Max: {fmt(s.max)}</span>}
          {s.avg != null && <span>Moy: {fmt(s.avg)}</span>}
          <span>N: {s.count}</span>
        </div>
      )}
    </div>
  );
}

function WidgetBody({ type, data, stats: s, lastPin }) {
  if (type === 'status') {
    const running = lastPin ? (typeof lastPin.data === 'boolean' ? lastPin.data : lastPin.data?.value > 0 || lastPin.data?.value === true) : null;
    return (
      <div className="py-3 text-center">
        <div className={`text-lg font-semibold ${running == null ? 'text-dark-500' : running ? 'text-green-400' : 'text-red-400'}`}>
          {running == null ? '⚪ n/a' : running ? '🟢 Running' : '🔴 Stopped'}
        </div>
        <div className="text-[10px] text-dark-500 mt-1">{lastPin?.timestamp_ms ? new Date(lastPin.timestamp_ms).toLocaleTimeString() : ''}</div>
      </div>
    );
  }
  if (type === 'value') {
    return <div className="py-4 text-center text-3xl font-bold text-blue-400">{s?.last != null ? fmt(s.last) : '—'}</div>;
  }
  if (type === 'gauge') {
    const min = s?.min ?? 0, max = s?.max ?? 1, last = s?.last ?? 0;
    const pct = max > min ? Math.max(0, Math.min(100, ((last - min) / (max - min)) * 100)) : 0;
    return (
      <div className="py-3">
        <div className="text-center text-2xl font-bold text-blue-400 mb-2">{s?.last != null ? fmt(s.last) : '—'}</div>
        <div className="h-2 bg-dark-800 rounded-full overflow-hidden"><div className="h-full bg-blue-500" style={{ width: `${pct}%` }} /></div>
        <div className="flex justify-between text-[10px] text-dark-500 mt-0.5"><span>{fmt(min)}</span><span>{fmt(max)}</span></div>
      </div>
    );
  }
  if (data.length === 0) return <p className="text-xs text-dark-500 py-6 text-center">En attente de données…</p>;
  const Chart = type === 'bar' ? BarChart : LineChart;
  const Series = type === 'bar'
    ? <Bar dataKey="value" fill="#38bdf8" radius={[3, 3, 0, 0]} isAnimationActive={false} />
    : <Line type="monotone" dataKey="value" stroke="#38bdf8" strokeWidth={2} dot={false} isAnimationActive={false} />;
  return (
    <ResponsiveContainer width="100%" height={150}>
      <Chart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
        <XAxis dataKey="time" tick={{ fill: '#94a3b8', fontSize: 9 }} minTickGap={40} />
        <YAxis tick={{ fill: '#94a3b8', fontSize: 10 }} width={32} />
        <Tooltip contentStyle={tooltipStyle} />
        {Series}
      </Chart>
    </ResponsiveContainer>
  );
}

function AddWidget({ sources, onAdd, onClose }) {
  const [label, setLabel] = useState(sources[0] || '');
  const [type, setType] = useState('line');
  const [range, setRange] = useState('1h');
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" onClick={onClose}>
      <div className="bg-dark-900 border border-dark-700 rounded-xl w-full max-w-md p-5 shadow-2xl" onClick={(e) => e.stopPropagation()}>
        <h3 className="text-white font-semibold mb-3">➕ Ajouter un widget</h3>
        {sources.length === 0 ? (
          <p className="text-sm text-dark-400 mb-3">Aucune source. Exécutez un pipeline avec <span className="font-mono">add_to_dashboard</span> d'abord.</p>
        ) : (
          <>
            <Label text="Donnée">
              <select value={label} onChange={(e) => setLabel(e.target.value)} className="w-full bg-dark-950 border border-dark-700 rounded-md px-3 py-2 text-sm text-white">
                {sources.map((s) => <option key={s} value={s}>{s}</option>)}
              </select>
            </Label>
            <Label text="Type de graphique">
              <select value={type} onChange={(e) => setType(e.target.value)} className="w-full bg-dark-950 border border-dark-700 rounded-md px-3 py-2 text-sm text-white">
                {CHART_TYPES.map(([v, l]) => <option key={v} value={v}>{l}</option>)}
              </select>
            </Label>
            <Label text="Plage de temps">
              <select value={range} onChange={(e) => setRange(e.target.value)} className="w-full bg-dark-950 border border-dark-700 rounded-md px-3 py-2 text-sm text-white">
                {Object.keys(RANGES).map((r) => <option key={r} value={r}>{r}</option>)}
              </select>
            </Label>
          </>
        )}
        <div className="flex justify-end gap-2 mt-4">
          <button onClick={onClose} className="text-dark-400 hover:text-white text-sm px-3 py-1.5">Annuler</button>
          {sources.length > 0 && (
            <button onClick={() => onAdd({ label, type, range })} className="bg-blue-600 hover:bg-blue-500 text-white text-sm px-3 py-1.5 rounded-md">Ajouter</button>
          )}
        </div>
      </div>
    </div>
  );
}

function Label({ text, children }) {
  return (
    <div className="mb-3">
      <label className="block text-[11px] text-dark-400 uppercase tracking-wider mb-1">{text}</label>
      {children}
    </div>
  );
}

function icon(type) {
  return { line: '📈', bar: '📊', gauge: '🎛️', value: '🔢', status: '🏭' }[type] || '📊';
}
function fmt(v) {
  if (typeof v !== 'number') return String(v ?? '—');
  return Number.isInteger(v) ? String(v) : v.toFixed(2);
}
