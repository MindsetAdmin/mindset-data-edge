import { useState, useEffect, useRef } from 'react';
import { ResponsiveContainer, LineChart, Line, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { Settings, X, Plus, LineChart as LineIcon, BarChart3, Gauge, Hash, Factory } from 'lucide-react';
import { useLiveSocket } from '../lib/useLiveSocket';
import { fetchDashboardPins } from '../api/client';
import Panel from './ui/Panel';
import StatusDot from './ui/StatusDot';

// Interactive dashboard widgets fed by add_to_dashboard pins (mindset/dashboard/#).
// 2026-07-01 redesign: emoji-free, MindSet design tokens, precise typography.

const STORE_KEY = 'mindset_widgets_v1';
const MAX_POINTS = 1000;
const RANGES = { '1m': 60e3, '5m': 300e3, '1h': 3600e3, '4h': 4 * 3600e3, '24h': 24 * 3600e3 };
const CHART_TYPES = [
    ['line', 'Line'],
    ['bar', 'Bar'],
    ['gauge', 'Gauge'],
    ['value', 'Value'],
    ['status', 'Status'],
];
const CHART_ICONS = {
    line: LineIcon,
    bar: BarChart3,
    gauge: Gauge,
    value: Hash,
    status: Factory,
};

// Recharts colors (CSS variable fallbacks for inline styles)
const CHART_COLOR = '#E5A445';          // brand accent
const CHART_GRID = '#2A2A31';           // border-subtle
const CHART_AXIS = '#6E6E7A';           // text-tertiary
const TOOLTIP_STYLE = {
    background: '#131316',
    border: '1px solid #2A2A31',
    borderRadius: 4,
    color: '#E8E8ED',
    fontSize: 11,
    fontFamily: 'Inter, system-ui, sans-serif',
};

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
    const [history, setHistory] = useState({});
    const [widgets, setWidgets] = useState(() => {
        try { return JSON.parse(sessionStorage.getItem(STORE_KEY)) || []; } catch { return []; }
    });
    const [showAdd, setShowAdd] = useState(false);
    const [now, setNow] = useState(Date.now());
    const histRef = useRef(history);
    histRef.current = history;

    useEffect(() => { sessionStorage.setItem(STORE_KEY, JSON.stringify(widgets)); }, [widgets]);
    useEffect(() => { const t = setInterval(() => setNow(Date.now()), 2000); return () => clearInterval(t); }, []);

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
        <Panel
            title="Pinned widgets"
            toolbar={
                <>
                    <StatusDot
                        state={connected ? 'running' : 'idle'}
                        pulse={connected}
                        label={connected ? 'Live' : 'Offline'}
                    />
                    <button
                        onClick={() => setShowAdd(true)}
                        className="inline-flex items-center gap-1 text-11 text-text-primary bg-accent hover:bg-[#c98d33] transition-colors rounded px-2 py-1"
                    >
                        <Plus size={12} strokeWidth={2} />
                        <span>Add</span>
                    </button>
                </>
            }
        >
            {widgets.length === 0 ? (
                <p className="text-13 text-text-tertiary py-8 text-center">
                    No widget yet. Click <span className="mono text-text-secondary">Add</span> to visualize a value emitted by{' '}
                    <span className="mono text-text-secondary">add_to_dashboard</span>.
                </p>
            ) : (
                <div className="grid sm:grid-cols-2 gap-3">
                    {widgets.map((w) => (
                        <Widget
                            key={w.id}
                            w={w}
                            series={history[w.label] || []}
                            now={now}
                            onRemove={() => removeWidget(w.id)}
                            onUpdate={(patch) => updateWidget(w.id, patch)}
                        />
                    ))}
                </div>
            )}

            {showAdd && <AddWidget sources={sources} onAdd={addWidget} onClose={() => setShowAdd(false)} />}
        </Panel>
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
    const Icon = CHART_ICONS[w.type] || LineIcon;

    return (
        <div className="bg-panel-alt border border-border-subtle hover:border-border-strong transition-colors rounded p-3">
            <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-1.5 text-13 font-medium text-text-primary truncate min-w-0">
                    <Icon size={13} strokeWidth={1.5} className="text-text-tertiary shrink-0" />
                    <span className="truncate">{w.label}</span>
                </div>
                <div className="flex items-center gap-0.5 shrink-0">
                    <button
                        onClick={() => setCfg((c) => !c)}
                        title="Configure"
                        className="p-1 text-text-tertiary hover:text-text-primary transition-colors"
                    >
                        <Settings size={12} strokeWidth={1.5} />
                    </button>
                    <button
                        onClick={onRemove}
                        title="Remove"
                        className="p-1 text-text-tertiary hover:text-status-stopped transition-colors"
                    >
                        <X size={12} strokeWidth={1.5} />
                    </button>
                </div>
            </div>

            {cfg && (
                <div className="flex flex-wrap gap-2 mb-2 text-11">
                    <select
                        value={w.type}
                        onChange={(e) => onUpdate({ type: e.target.value })}
                        className="bg-panel border border-border-subtle hover:border-border-strong rounded px-1.5 py-1 text-text-primary transition-colors"
                    >
                        {CHART_TYPES.map(([v, l]) => <option key={v} value={v}>{l}</option>)}
                    </select>
                    <select
                        value={w.range}
                        onChange={(e) => onUpdate({ range: e.target.value })}
                        className="bg-panel border border-border-subtle hover:border-border-strong rounded px-1.5 py-1 text-text-primary mono transition-colors"
                    >
                        {Object.keys(RANGES).map((r) => <option key={r} value={r}>{r}</option>)}
                    </select>
                </div>
            )}

            <WidgetBody type={w.type} data={data} stats={s} lastPin={lastPin} />

            {s && w.type !== 'status' && (
                <div className="flex flex-wrap gap-x-3 gap-y-0.5 mt-2 text-11 text-text-muted">
                    {s.last != null && <span>Last: <span className="mono tabular text-text-secondary">{fmt(s.last)}</span></span>}
                    {s.min != null && <span>Min: <span className="mono tabular">{fmt(s.min)}</span></span>}
                    {s.max != null && <span>Max: <span className="mono tabular">{fmt(s.max)}</span></span>}
                    {s.avg != null && <span>Avg: <span className="mono tabular">{fmt(s.avg)}</span></span>}
                    <span>N: <span className="mono tabular">{s.count}</span></span>
                </div>
            )}
        </div>
    );
}

function WidgetBody({ type, data, stats: s, lastPin }) {
    if (type === 'status') {
        const running = lastPin
            ? (typeof lastPin.data === 'boolean' ? lastPin.data : lastPin.data?.value > 0 || lastPin.data?.value === true)
            : null;
        const state = running == null ? 'idle' : running ? 'running' : 'stopped';
        const label = running == null ? 'No signal' : running ? 'Running' : 'Stopped';
        return (
            <div className="py-4 flex flex-col items-center gap-2">
                <StatusDot state={state} pulse={running === true} size={12} label={null} />
                <div className={`text-15 font-medium ${
                    running == null ? 'text-text-tertiary' :
                    running ? 'text-status-running' : 'text-status-stopped'
                }`}>
                    {label}
                </div>
                {lastPin?.timestamp_ms && (
                    <div className="text-11 text-text-muted mono">
                        {new Date(lastPin.timestamp_ms).toLocaleTimeString()}
                    </div>
                )}
            </div>
        );
    }
    if (type === 'value') {
        return (
            <div className="py-4 text-center">
                <div className="mono text-20 font-medium text-text-primary tabular">
                    {s?.last != null ? fmt(s.last) : '—'}
                </div>
            </div>
        );
    }
    if (type === 'gauge') {
        const min = s?.min ?? 0, max = s?.max ?? 1, last = s?.last ?? 0;
        const pct = max > min ? Math.max(0, Math.min(100, ((last - min) / (max - min)) * 100)) : 0;
        return (
            <div className="py-2">
                <div className="text-center mono text-20 font-medium text-text-primary tabular mb-2">
                    {s?.last != null ? fmt(s.last) : '—'}
                </div>
                <div className="h-1.5 bg-elevated rounded overflow-hidden">
                    <div className="h-full bg-accent transition-all" style={{ width: `${pct}%` }} />
                </div>
                <div className="flex justify-between text-11 text-text-muted mono mt-1">
                    <span>{fmt(min)}</span><span>{fmt(max)}</span>
                </div>
            </div>
        );
    }
    if (data.length === 0) {
        return <p className="text-11 text-text-tertiary py-6 text-center italic">Awaiting data…</p>;
    }
    const Chart = type === 'bar' ? BarChart : LineChart;
    const Series = type === 'bar'
        ? <Bar dataKey="value" fill={CHART_COLOR} radius={[2, 2, 0, 0]} isAnimationActive={false} />
        : <Line type="monotone" dataKey="value" stroke={CHART_COLOR} strokeWidth={1.5} dot={false} isAnimationActive={false} />;
    return (
        <ResponsiveContainer width="100%" height={150}>
            <Chart data={data}>
                <CartesianGrid strokeDasharray="2 4" stroke={CHART_GRID} />
                <XAxis dataKey="time" tick={{ fill: CHART_AXIS, fontSize: 10, fontFamily: 'JetBrains Mono, monospace' }} minTickGap={40} axisLine={{ stroke: CHART_GRID }} tickLine={{ stroke: CHART_GRID }} />
                <YAxis tick={{ fill: CHART_AXIS, fontSize: 10, fontFamily: 'JetBrains Mono, monospace' }} width={36} axisLine={{ stroke: CHART_GRID }} tickLine={{ stroke: CHART_GRID }} />
                <Tooltip contentStyle={TOOLTIP_STYLE} labelStyle={{ color: '#A8A8B2' }} />
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
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4"
            onClick={onClose}
        >
            <div
                className="bg-panel border border-border-subtle rounded w-full max-w-md p-5 shadow-2xl"
                onClick={(e) => e.stopPropagation()}
            >
                <h3 className="text-15 font-medium text-text-primary mb-4">Add widget</h3>
                {sources.length === 0 ? (
                    <p className="text-13 text-text-tertiary mb-3">
                        No source yet. Run a pipeline with{' '}
                        <span className="mono text-text-secondary">add_to_dashboard</span> first.
                    </p>
                ) : (
                    <>
                        <Field label="Source">
                            <select
                                value={label}
                                onChange={(e) => setLabel(e.target.value)}
                                className="w-full bg-canvas border border-border-subtle hover:border-border-strong rounded px-3 py-1.5 text-13 text-text-primary transition-colors"
                            >
                                {sources.map((s) => <option key={s} value={s}>{s}</option>)}
                            </select>
                        </Field>
                        <Field label="Chart type">
                            <select
                                value={type}
                                onChange={(e) => setType(e.target.value)}
                                className="w-full bg-canvas border border-border-subtle hover:border-border-strong rounded px-3 py-1.5 text-13 text-text-primary transition-colors"
                            >
                                {CHART_TYPES.map(([v, l]) => <option key={v} value={v}>{l}</option>)}
                            </select>
                        </Field>
                        <Field label="Time range">
                            <select
                                value={range}
                                onChange={(e) => setRange(e.target.value)}
                                className="w-full bg-canvas border border-border-subtle hover:border-border-strong rounded px-3 py-1.5 text-13 text-text-primary mono transition-colors"
                            >
                                {Object.keys(RANGES).map((r) => <option key={r} value={r}>{r}</option>)}
                            </select>
                        </Field>
                    </>
                )}
                <div className="flex justify-end gap-2 mt-5">
                    <button
                        onClick={onClose}
                        className="text-13 text-text-secondary hover:text-text-primary px-3 py-1.5 transition-colors"
                    >
                        Cancel
                    </button>
                    {sources.length > 0 && (
                        <button
                            onClick={() => onAdd({ label, type, range })}
                            className="bg-accent hover:bg-[#c98d33] text-text-primary text-13 px-3 py-1.5 rounded transition-colors"
                        >
                            Add
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}

function Field({ label, children }) {
    return (
        <div className="mb-3">
            <label className="block text-11 text-text-secondary uppercase tracking-wide mb-1">{label}</label>
            {children}
        </div>
    );
}

function fmt(v) {
    if (typeof v !== 'number') return String(v ?? '—');
    return Number.isInteger(v) ? String(v) : v.toFixed(2);
}
