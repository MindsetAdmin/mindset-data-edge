// Non-interactive background band (ENTRÉE / CŒUR / SORTIE) that pans & zooms
// with the canvas. Rendered behind the real nodes.
export default function ZoneNode({ data }) {
  return (
    <div
      className={`w-full h-full rounded-xl border-2 border-dashed ${data.className}`}
      style={{ pointerEvents: 'none' }}
    >
      <div className="text-[11px] font-semibold uppercase tracking-widest px-3 py-1.5 opacity-70">
        {data.label}
      </div>
    </div>
  );
}
