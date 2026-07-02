// TimeRangeSelector — Grafana-style segmented control for time ranges.
// Standard values: 5m, 15m, 1h, 6h, 24h, 7d.
const DEFAULT_RANGES = [
    { value: '5m', label: '5m' },
    { value: '15m', label: '15m' },
    { value: '1h', label: '1h' },
    { value: '6h', label: '6h' },
    { value: '24h', label: '24h' },
    { value: '7d', label: '7d' },
];

export default function TimeRangeSelector({
    value,
    onChange,
    ranges = DEFAULT_RANGES,
    className = '',
}) {
    return (
        <div
            className={`inline-flex items-center bg-panel border border-border-subtle rounded overflow-hidden ${className}`}
            role="group"
            aria-label="Time range"
        >
            {ranges.map((r) => {
                const active = r.value === value;
                return (
                    <button
                        key={r.value}
                        type="button"
                        onClick={() => onChange && onChange(r.value)}
                        className={`px-2.5 py-1 mono text-11 transition-colors border-r border-border-subtle last:border-r-0 ${
                            active
                                ? 'bg-elevated text-text-primary'
                                : 'text-text-tertiary hover:text-text-primary hover:bg-panel-alt'
                        }`}
                    >
                        {r.label}
                    </button>
                );
            })}
        </div>
    );
}
