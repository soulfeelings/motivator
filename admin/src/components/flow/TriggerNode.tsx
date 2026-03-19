import { Handle, Position } from '@xyflow/react'
import type { NodeProps } from '@xyflow/react'
import { Zap } from 'lucide-react'

export default function TriggerNode({ data }: NodeProps) {
  return (
    <div className="bg-gray-900 border-2 border-amber-500/50 rounded-xl p-4 min-w-[200px] shadow-lg shadow-amber-500/10">
      <Handle type="source" position={Position.Bottom} className="!bg-amber-500 !w-3 !h-3" />
      <div className="flex items-center gap-2 mb-2">
        <div className="w-7 h-7 rounded-lg bg-amber-500/20 flex items-center justify-center">
          <Zap size={16} className="text-amber-400" />
        </div>
        <span className="text-xs font-bold uppercase tracking-wider text-amber-400">Trigger</span>
      </div>
      <p className="text-white text-sm font-medium">{(data as any).metric || 'Select metric...'}</p>
      {(data as any).description && (
        <p className="text-gray-500 text-xs mt-1">{(data as any).description}</p>
      )}
    </div>
  )
}
