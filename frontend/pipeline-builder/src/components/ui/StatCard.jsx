// StatCard — dense KPI card in the Grafana stat-panel style.
// Big monospace number, small label, optional delta comparison.
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

export default function StatCard({
    label,
    value,
    unit,
    delta,          // number or null — positive = up, negative = down
    deltaLabel,     // e.g. "vs yesterday"
    deltaGoodDirection = 'up', // 'up' | 'down' — determines color semantics
    sparkline,      // optional inline sparkline component
}) {
    const hasDelta = typeof delta === 'number' && Number.isFinite(delta);
    let deltaColor = 'text-text-tertiary';
    let DeltaIcon = Minus;
    if (hasDelta) {
        if (delta > 0) {
            DeltaIcon = TrendingUp;
            deltaColor = deltaGoodDirection === 'up' ? 'text-status-running' : 'text-status-stopped';
        } else if (delta < 0) {
            DeltaIcon = TrendingDown;
            deltaColor = deltaGoodDirection === 'up' ? 'text-status-stopped' : 'text-status-running';
        }
    }

    return (
        <div className="bg-panel border border-border-subtle hover:border-border-strong transition-colors rounded p-3 flex flex-col gap-1.5 min-w-0">
            <div className="text-11 text-text-secondary uppercase tracking-wide truncate">
                {label}
            </div>
            <div className="flex items-baseline gap-1.5 min-w-0">
                <span className="mono text-20 font-medium text-text-primary tabular truncate">
                    {value}
                </span>
                {unit && (
                    <span className="text-13 text-text-tertiary shrink-0">{unit}</span>
                )}
            </div>
            {hasDelta && (
                <div className={`flex items-center gap-1 text-11 ${deltaColor}`}>
                    <DeltaIcon size={12} strokeWidth={1.5} />
                    <span className="mono tabular">
                        {delta > 0 ? '+' : ''}
                        {delta}
                    </span>
                    {deltaLabel && (
                        <span className="text-text-tertiary">{deltaLabel}</span>
                    )}
                </div>
            )}
            {sparkline && <div className="h-6 mt-1">{sparkline}</div>}
        </div>
    );
}
