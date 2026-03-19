import { Handle, Position } from '@xyflow/react'
import type { NodeProps } from '@xyflow/react'
import { GitBranch } from 'lucide-react'

const opLabels: Record<string, string> = { gte: '≥', lte: '≤', eq: '=', gt: '>', lt: '<' }

export default function ConditionNode({ data }: NodeProps) {
  const d = data as any
  return (
    <div className="bg-gray-900 border-2 border-violet-500/50 rounded-xl p-4 min-w-[200px] shadow-lg shadow-violet-500/10">
      <Handle type="target" position={Position.Top} className="!bg-violet-500 !w-3 !h-3" />
      <Handle type="source" position={Position.Bottom} className="!bg-violet-500 !w-3 !h-3" />
      <div className="flex items-center gap-2 mb-2">
        <div className="w-7 h-7 rounded-lg bg-violet-500/20 flex items-center justify-center">
          <GitBranch size={16} className="text-violet-400" />
        </div>
        <span className="text-xs font-bold uppercase tracking-wider text-violet-400">Condition</span>
      </div>
      <p className="text-white text-sm font-medium">
        {d.metric || '?'} {opLabels[d.operator] || d.operator || '?'} {d.threshold ?? '?'}
      </p>
    </div>
  )
}
