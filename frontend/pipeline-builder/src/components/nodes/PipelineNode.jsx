import { Handle, Position } from 'reactflow';
import { typeStyle } from '../../lib/functionMeta';

// Custom ReactFlow node for a pipeline step (transform/calculate/condition/output/connector).
export default function PipelineNode({ data, selected }) {
  const s = typeStyle(data.type);
  const configKeys = Object.keys(data.config || {});

  return (
    <div
      className={`relative rounded-lg border-2 bg-dark-800 px-3 py-2 min-w-[160px] shadow-lg ${s.border} ${
        selected ? 'ring-2 ring-blue-400' : ''
      }`}
    >
      <Handle type="target" position={Position.Left} className="!w-2 !h-2 !bg-dark-300" />

      {data.runStatus && (
        <span
          className="absolute -top-1.5 -right-1.5 text-[10px]"
          title={`run: ${data.runStatus}`}
        >
          {data.runStatus === 'success' ? '✅' : data.runStatus === 'failed' ? '❌' : '⏳'}
        </span>
      )}

      <div className="flex items-center gap-2">
        <span className="text-base">{s.icon}</span>
        <div className="min-w-0">
          <div className="text-sm font-medium text-white leading-tight truncate">{data.name}</div>
          <div className="text-[10px] text-dark-400 truncate">{data.function}</div>
        </div>
      </div>

      {configKeys.length > 0 && (
        <div className="mt-1.5 flex flex-wrap gap-1">
          {configKeys.slice(0, 3).map((k) => (
            <span key={k} className={`text-[9px] px-1.5 py-0.5 rounded ${s.chip}`}>
              {k}: {String(data.config[k])}
            </span>
          ))}
        </div>
      )}

      <Handle type="source" position={Position.Right} className="!w-2 !h-2 !bg-dark-300" />
    </div>
  );
}
