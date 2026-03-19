import { Zap, GitBranch, Award } from 'lucide-react'
import type { DragEvent } from 'react'

const nodeTypes = [
  { type: 'trigger', label: 'Trigger', icon: Zap, color: 'border-amber-500/30 hover:border-amber-500/60 text-amber-400', desc: 'Metric event' },
  { type: 'condition', label: 'Condition', icon: GitBranch, color: 'border-violet-500/30 hover:border-violet-500/60 text-violet-400', desc: 'Check threshold' },
  { type: 'action', label: 'Action', icon: Award, color: 'border-emerald-500/30 hover:border-emerald-500/60 text-emerald-400', desc: 'Award reward' },
]

export default function NodeSidebar() {
  function onDragStart(e: DragEvent, nodeType: string) {
    e.dataTransfer.setData('application/reactflow', nodeType)
    e.dataTransfer.effectAllowed = 'move'
  }

  return (
    <div className="w-56 border-r border-gray-800 bg-gray-950 p-4 space-y-3">
      <h3 className="text-xs font-bold uppercase tracking-wider text-gray-500 mb-4">Drag to canvas</h3>
      {nodeTypes.map(({ type, label, icon: Icon, color, desc }) => (
        <div
          key={type}
          draggable
          onDragStart={(e) => onDragStart(e, type)}
          className={`border-2 rounded-xl p-3 cursor-grab active:cursor-grabbing transition-colors ${color}`}
        >
          <div className="flex items-center gap-2">
            <Icon size={18} />
            <div>
              <p className="text-sm font-medium text-white">{label}</p>
              <p className="text-xs text-gray-500">{desc}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
