import { Handle, Position } from 'reactflow';

// The pipeline entry point. Maps to pipeline.Trigger (not into nodes[]).
// Fixed id "trigger" so downstream depends_on references resolve.
export default function TriggerNode({ data, selected }) {
  const topic = data.config?.topic;
  return (
    <div
      className={`rounded-lg border-2 border-cyan-500/70 bg-cyan-500/10 px-3 py-2 min-w-[160px] shadow-lg ${
        selected ? 'ring-2 ring-cyan-300' : ''
      }`}
    >
      <div className="flex items-center gap-2">
        <span className="text-base">⚡</span>
        <div className="min-w-0">
          <div className="text-sm font-medium text-white leading-tight">Trigger</div>
          <div className="text-[10px] text-dark-400 truncate">{data.function}</div>
        </div>
      </div>
      {topic && (
        <div className="mt-1.5">
          <span className="text-[9px] px-1.5 py-0.5 rounded bg-cyan-500/15 text-cyan-300">topic: {topic}</span>
        </div>
      )}
      <Handle type="source" position={Position.Right} className="!w-2 !h-2 !bg-cyan-300" />
    </div>
  );
}
