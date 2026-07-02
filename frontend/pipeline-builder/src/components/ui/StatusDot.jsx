// StatusDot — the small colored circle that replaces ● / 📌 / ⚠️ emojis.
// State drives color; optional `pulse` adds a subtle animation.
const STATE_COLORS = {
    running: 'bg-status-running',
    stopped: 'bg-status-stopped',
    warn: 'bg-status-warn',
    info: 'bg-status-info',
    idle: 'bg-status-idle',
};

const STATE_LABELS = {
    running: 'Running',
    stopped: 'Stopped',
    warn: 'Warning',
    info: 'Info',
    idle: 'Idle',
};

export default function StatusDot({
    state = 'idle',
    pulse = false,
    label,             // optional inline text next to the dot
    size = 8,          // 8 default; use 6 for very dense contexts
}) {
    const color = STATE_COLORS[state] || STATE_COLORS.idle;
    const text = label ?? STATE_LABELS[state];
    return (
        <span className="inline-flex items-center gap-1.5">
            <span
                className={`inline-block rounded-full shrink-0 ${color} ${
                    pulse ? 'animate-pulse' : ''
                }`}
                style={{ width: size, height: size }}
                aria-label={text}
            />
            {label !== null && text && (
                <span className="text-11 text-text-secondary">{text}</span>
            )}
        </span>
    );
}
