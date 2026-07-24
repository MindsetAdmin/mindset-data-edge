import { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { ResponsiveContainer, LineChart, Line, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { fetchTags } from '../api/client';
import { useLiveSocket } from '../lib/useLiveSocket';
import PickerModal from './PickerModal';

const COLORS = ['#38bdf8', '#f59e0b', '#34d399', '#a78bfa', '#f87171', '#fbbf24'];
const MAX_POINTS = 60;

// Coerce a tag value to a number for charting (booleans → 0/1; non-numeric → null).
function toNum(v) {
  if (typeof v === 'boolean') return v ? 1 : 0;
  const n = Number(v);
  return Number.isFinite(n) ? n : null;
}

// Live, user-selected tag chart: pick tag(s) and watch their values stream in
// over the WebSocket. Self-contained (own /api/ws subscription).
export default function LiveDataPanel() {
  const { t } = useTranslation();
  const [selected, setSelected] = useState([]); // [{node_id, name}]
  const [available, setAvailable] = useState([]);
  const [rows, setRows] = useState([]); // [{time, <node_id>: number, ...}]
  const [latest, setLatest] = useState({}); // node_id -> raw value
  const [showPicker, setShowPicker] = useState(false);
  const selRef = useRef(selected);
  selRef.current = selected;

  useEffect(() => {
    fetchTags().then((d) => setAvailable(d.tags || [])).catch(() => {});
  }, []);

  const connected = useLiveSocket((msg) => {
    if (msg.type !== 'tag') return;
    const tagMsg = msg.data;
    if (!tagMsg || !selRef.current.some((s) => s.node_id === tagMsg.node_id)) return;
    setLatest((prev) => ({ ...prev, [tagMsg.node_id]: tagMsg.value }));
    const num = toNum(tagMsg.value);
    if (num === null) return;
    setRows((prev) => {
      const row = { time: new Date().toLocaleTimeString().slice(0, 8), [tagMsg.node_id]: num };
      const next = [...prev, row];
      return next.length > MAX_POINTS ? next.slice(next.length - MAX_POINTS) : next;
    });
  });

  const addTag = (o) => {
    setSelected((prev) => (prev.some((s) => s.node_id === o.value) ? prev : [...prev, { node_id: o.value, name: o.label }]));
    setShowPicker(false);
  };
  const removeTag = (id) => setSelected((prev) => prev.filter((s) => s.node_id !== id));

  const pickerOptions = available
    .filter((tg) => !selected.some((s) => s.node_id === tg.node_id))
    .map((tg) => ({
      value: tg.node_id,
      label: tg.name || tg.node_id,
      sub: `${t('builder.value')}: ${tg.value} · ${tg.data_type}`,
      badge: tg.node_id,
      group: (tg.name || '').split('.')[0] || t('liveData.others'),
    }));

  return (
    <div className="bg-dark-900 border border-dark-700 rounded-lg p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium text-dark-300">📡 {t('liveData.title')}</h3>
        <div className="flex items-center gap-3">
          <span className={`text-[11px] ${connected ? 'text-green-400' : 'text-dark-500'}`}>
            {connected ? '● live' : `○ ${t('dashboard.offline')}`}
          </span>
          <button onClick={() => setShowPicker(true)} className="text-xs bg-dark-700 hover:bg-dark-600 text-white px-2.5 py-1 rounded-md transition">
            + tag
          </button>
        </div>
      </div>

      {selected.length === 0 ? (
        <p className="text-sm text-dark-500 py-8 text-center">
          {t('liveData.chooseTags')}
        </p>
      ) : (
        <>
          <div className="flex flex-wrap gap-2 mb-3">
            {selected.map((s, i) => (
              <span key={s.node_id} className="inline-flex items-center gap-2 text-xs px-2 py-1 rounded bg-dark-800 border border-dark-700">
                <span className="w-2 h-2 rounded-full" style={{ background: COLORS[i % COLORS.length] }} />
                {s.name}
                <span className="font-mono text-dark-200">{String(latest[s.node_id] ?? '—')}</span>
                <button onClick={() => removeTag(s.node_id)} className="text-dark-500 hover:text-red-400">×</button>
              </span>
            ))}
          </div>
          <ResponsiveContainer width="100%" height={240}>
            <LineChart data={rows}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
              <XAxis dataKey="time" tick={{ fill: '#94a3b8', fontSize: 10 }} minTickGap={40} />
              <YAxis tick={{ fill: '#94a3b8', fontSize: 11 }} />
              <Tooltip contentStyle={{ background: '#0f172a', border: '1px solid #334155', borderRadius: 8, color: '#e2e8f0' }} />
              {selected.map((s, i) => (
                <Line
                  key={s.node_id}
                  type="monotone"
                  dataKey={s.node_id}
                  name={s.name}
                  stroke={COLORS[i % COLORS.length]}
                  dot={false}
                  isAnimationActive={false}
                  connectNulls
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
          {rows.length === 0 && (
            <p className="text-[11px] text-dark-500 text-center mt-2">{t('liveData.waitingForValues')}</p>
          )}
        </>
      )}

      {showPicker && (
        <PickerModal
          title={t('liveData.pickerTitle')}
          options={pickerOptions}
          allowCustom={false}
          onSelect={addTag}
          onClose={() => setShowPicker(false)}
        />
      )}
    </div>
  );
}
