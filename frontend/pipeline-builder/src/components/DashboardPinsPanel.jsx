import { useState, useEffect } from 'react';
import { useLiveSocket } from '../lib/useLiveSocket';
import { fetchDashboardPins } from '../api/client';

// Shows widgets pinned via the add_to_dashboard pipeline function. Each message
// arrives over the WebSocket as {type:'dashboard', data:{label,kind,data,...}};
// we keep the latest per label. Seeds from the server snapshot so pins show even
// if they were published before this panel opened.
export default function DashboardPinsPanel() {
  const [pins, setPins] = useState({}); // label -> {label, kind, data, timestamp_ms}

  useEffect(() => {
    fetchDashboardPins()
      .then((d) => {
        const seed = {};
        (d.pins || []).forEach((p) => p?.label && (seed[p.label] = p));
        setPins((prev) => ({ ...seed, ...prev }));
      })
      .catch(() => {});
  }, []);

  const connected = useLiveSocket((msg) => {
    if (msg.type !== 'dashboard' || !msg.data?.label) return;
    setPins((prev) => ({ ...prev, [msg.data.label]: msg.data }));
  });

  const items = Object.values(pins).sort((a, b) => a.label.localeCompare(b.label));

  const render = (w) => {
    const d = w.data;
    if (w.kind === 'event') {
      return (
        <div className="text-xs text-dark-300">
          {d?.work_center || d?.workCenter || 'event'} · {d?.cause || d?.duration_seconds != null ? `${Math.round(d.duration_seconds || 0)}s` : ''}
        </div>
      );
    }
    // value: show a primitive directly, otherwise the most relevant field
    const v = typeof d === 'object' && d !== null ? (d.value ?? d.total_cost_eur ?? JSON.stringify(d)) : d;
    return <div className="text-2xl font-bold text-blue-400">{String(v)}</div>;
  };

  return (
    <div className="bg-dark-900 border border-dark-700 rounded-lg p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium text-dark-300">📌 Widgets épinglés</h3>
        <span className={`text-[11px] ${connected ? 'text-green-400' : 'text-dark-500'}`}>{connected ? '● live' : '○'}</span>
      </div>

      {items.length === 0 ? (
        <p className="text-sm text-dark-500 py-6 text-center">
          Aucun widget. Ajoutez une fonction <span className="font-mono text-dark-300">add_to_dashboard</span> à un pipeline pour épingler une donnée ou un événement ici.
        </p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
          {items.map((w) => (
            <div key={w.label} className="bg-dark-950 border border-dark-700 rounded-lg p-3">
              <div className="text-[11px] text-dark-500 uppercase tracking-wider mb-1 truncate">{w.label}</div>
              {render(w)}
              <div className="text-[10px] text-dark-600 mt-1">
                {w.timestamp_ms ? new Date(w.timestamp_ms).toLocaleTimeString() : ''}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
