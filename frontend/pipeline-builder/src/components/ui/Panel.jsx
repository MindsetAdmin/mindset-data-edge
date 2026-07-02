// Panel — Grafana-style panel container.
// Every dashboard section, form group, list wrapper uses this.
// Standard structure: header (title + toolbar + actions) → body.
import { MoreHorizontal } from 'lucide-react';

export default function Panel({
    title,
    subtitle,
    toolbar,
    actions,
    loading = false,
    error = null,
    noPadding = false,
    className = '',
    children,
}) {
    return (
        <div
            className={`bg-panel border border-border-subtle hover:border-border-strong transition-colors rounded overflow-hidden ${className}`}
        >
            {(title || subtitle || toolbar || actions) && (
                <div className="flex items-center gap-3 px-3 py-2 border-b border-border-subtle">
                    <div className="flex-1 min-w-0">
                        {title && (
                            <h3 className="text-13 font-medium text-text-primary truncate">
                                {title}
                            </h3>
                        )}
                        {subtitle && (
                            <p className="text-11 text-text-tertiary truncate mt-0.5">
                                {subtitle}
                            </p>
                        )}
                    </div>
                    {toolbar && <div className="flex items-center gap-2">{toolbar}</div>}
                    {actions && (
                        <button
                            type="button"
                            className="p-1 text-text-tertiary hover:text-text-primary transition-colors"
                            title="Actions"
                        >
                            <MoreHorizontal size={16} strokeWidth={1.5} />
                        </button>
                    )}
                </div>
            )}
            <div className={noPadding ? '' : 'p-3'}>
                {loading && (
                    <div className="text-11 text-text-tertiary italic">Loading…</div>
                )}
                {error && (
                    <div className="text-11 text-status-stopped">{String(error)}</div>
                )}
                {!loading && !error && children}
            </div>
        </div>
    );
}
