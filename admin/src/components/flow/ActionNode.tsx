import { Handle, Position } from '@xyflow/react'
import type { NodeProps } from '@xyflow/react'
import { Award, Zap, Coins } from 'lucide-react'

export default function ActionNode({ data }: NodeProps) {
  const d = data as any
  return (
    <div className="bg-gray-900 border-2 border-emerald-500/50 rounded-xl p-4 min-w-[200px] shadow-lg shadow-emerald-500/10">
      <Handle type="target" position={Position.Top} className="!bg-emerald-500 !w-3 !h-3" />
      <Handle type="source" position={Position.Bottom} className="!bg-emerald-500 !w-3 !h-3" />
      <div className="flex items-center gap-2 mb-2">
        <div className="w-7 h-7 rounded-lg bg-emerald-500/20 flex items-center justify-center">
          <Award size={16} className="text-emerald-400" />
        </div>
        <span className="text-xs font-bold uppercase tracking-wider text-emerald-400">Action</span>
      </div>
      <p className="text-white text-sm font-medium">{d.label || 'Reward'}</p>
      <div className="flex gap-3 mt-2">
        {d.xp > 0 && (
          <span className="inline-flex items-center gap-1 text-xs text-emerald-400">
            <Zap size={12} /> +{d.xp} XP
          </span>
        )}
        {d.coins > 0 && (
          <span className="inline-flex items-center gap-1 text-xs text-amber-400">
            <Coins size={12} /> +{d.coins}
          </span>
        )}
        {d.badge && (
          <span className="inline-flex items-center gap-1 text-xs text-violet-400">
            <Award size={12} /> Badge
          </span>
        )}
      </div>
    </div>
  )
}
