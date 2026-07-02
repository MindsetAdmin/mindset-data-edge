// DenseTable — dense list replacing card-per-item patterns.
// 32px row height, monospace for numeric/id columns, zebra stripe on hover.
//
// columns: [{ key, label, align: 'left'|'right'|'center', mono: boolean, width: string }]
// rows:    array of objects keyed by column.key
// getRowKey: (row, index) => string  (optional; defaults to index)
// onRowClick: (row) => void  (optional)
export default function DenseTable({
    columns,
    rows,
    getRowKey,
    onRowClick,
    emptyLabel = 'No data',
    className = '',
}) {
    if (!rows || rows.length === 0) {
        return (
            <div
                className={`text-11 text-text-tertiary italic px-3 py-2 ${className}`}
            >
                {emptyLabel}
            </div>
        );
    }
    return (
        <div className={`overflow-x-auto ${className}`}>
            <table className="w-full border-collapse">
                <thead>
                    <tr className="border-b border-border-subtle">
                        {columns.map((c) => (
                            <th
                                key={c.key}
                                className={`text-11 font-medium text-text-secondary uppercase tracking-wide px-3 py-1.5 ${
                                    c.align === 'right' ? 'text-right' :
                                    c.align === 'center' ? 'text-center' : 'text-left'
                                }`}
                                style={c.width ? { width: c.width } : undefined}
                            >
                                {c.label}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {rows.map((row, i) => (
                        <tr
                            key={getRowKey ? getRowKey(row, i) : i}
                            onClick={onRowClick ? () => onRowClick(row) : undefined}
                            className={`border-b border-border-subtle last:border-b-0 hover:bg-panel-alt transition-colors ${
                                onRowClick ? 'cursor-pointer' : ''
                            }`}
                        >
                            {columns.map((c) => {
                                const val = row[c.key];
                                const displayed = c.render
                                    ? c.render(val, row, i)
                                    : val == null || val === ''
                                    ? '—'
                                    : val;
                                return (
                                    <td
                                        key={c.key}
                                        className={`text-13 text-text-primary px-3 py-2 ${
                                            c.mono ? 'mono tabular' : ''
                                        } ${
                                            c.align === 'right' ? 'text-right' :
                                            c.align === 'center' ? 'text-center' : 'text-left'
                                        }`}
                                    >
                                        {displayed}
                                    </td>
                                );
                            })}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
