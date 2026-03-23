import type { Node } from '@xyflow/react'
import { X, Zap, GitBranch, Award } from 'lucide-react'

interface NodeEditPanelProps {
  node: Node
  onUpdate: (nodeId: string, data: Record<string, unknown>) => void
  onClose: () => void
}

const inputClass = 'w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500'
const labelClass = 'block text-xs font-medium text-gray-400 mb-1.5'
const selectClass = 'w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm focus:outline-none focus:ring-2 focus:ring-violet-500'

const nodeConfig = {
  trigger: { label: 'Trigger', icon: Zap, color: 'text-amber-400', border: 'border-amber-500/30' },
  condition: { label: 'Condition', icon: GitBranch, color: 'text-violet-400', border: 'border-violet-500/30' },
  action: { label: 'Action', icon: Award, color: 'text-emerald-400', border: 'border-emerald-500/30' },
}

export default function NodeEditPanel({ node, onUpdate, onClose }: NodeEditPanelProps) {
  const config = nodeConfig[node.type as keyof typeof nodeConfig]
  if (!config) return null

  const Icon = config.icon
  const data = node.data as Record<string, unknown>

  function update(field: string, value: unknown) {
    onUpdate(node.id, { ...data, [field]: value })
  }

  return (
    <div className="w-72 border-l border-gray-800 bg-gray-950 flex flex-col">
      <div className={`flex items-center justify-between p-4 border-b ${config.border}`}>
        <div className="flex items-center gap-2">
          <Icon size={16} className={config.color} />
          <span className="text-sm font-semibold text-white">{config.label}</span>
        </div>
        <button onClick={onClose} className="text-gray-500 hover:text-gray-300 transition-colors">
          <X size={16} />
        </button>
      </div>

      <div className="p-4 space-y-4 overflow-y-auto flex-1">
        {node.type === 'trigger' && (
          <>
            <div>
              <label className={labelClass}>Metric</label>
              <input
                className={inputClass}
                value={(data.metric as string) ?? ''}
                onChange={(e) => update('metric', e.target.value)}
                placeholder="e.g. deals_closed"
              />
            </div>
            <div>
              <label className={labelClass}>Description</label>
              <input
                className={inputClass}
                value={(data.description as string) ?? ''}
                onChange={(e) => update('description', e.target.value)}
                placeholder="When this event fires"
              />
            </div>
          </>
        )}

        {node.type === 'condition' && (
          <>
            <div>
              <label className={labelClass}>Metric</label>
              <input
                className={inputClass}
                value={(data.metric as string) ?? ''}
                onChange={(e) => update('metric', e.target.value)}
                placeholder="e.g. deals_closed"
              />
            </div>
            <div>
              <label className={labelClass}>Operator</label>
              <select
                className={selectClass}
                value={(data.operator as string) ?? 'gte'}
                onChange={(e) => update('operator', e.target.value)}
              >
                <option value="gte">≥ Greater or equal</option>
                <option value="gt">&gt; Greater than</option>
                <option value="eq">= Equal to</option>
                <option value="lte">≤ Less or equal</option>
                <option value="lt">&lt; Less than</option>
              </select>
            </div>
            <div>
              <label className={labelClass}>Threshold</label>
              <input
                type="number"
                className={inputClass}
                value={(data.threshold as number) ?? 0}
                onChange={(e) => update('threshold', Number(e.target.value))}
                placeholder="10"
              />
            </div>
          </>
        )}

        {node.type === 'action' && (
          <>
            <div>
              <label className={labelClass}>Label</label>
              <input
                className={inputClass}
                value={(data.label as string) ?? ''}
                onChange={(e) => update('label', e.target.value)}
                placeholder="e.g. Award Reward"
              />
            </div>
            <div>
              <label className={labelClass}>XP Reward</label>
              <input
                type="number"
                className={inputClass}
                value={(data.xp as number) ?? 0}
                onChange={(e) => update('xp', Number(e.target.value))}
                placeholder="50"
              />
            </div>
            <div>
              <label className={labelClass}>Coins Reward</label>
              <input
                type="number"
                className={inputClass}
                value={(data.coins as number) ?? 0}
                onChange={(e) => update('coins', Number(e.target.value))}
                placeholder="25"
              />
            </div>
            <div className="flex items-center gap-3">
              <input
                type="checkbox"
                id="badge-toggle"
                checked={(data.badge as boolean) ?? false}
                onChange={(e) => update('badge', e.target.checked)}
                className="w-4 h-4 rounded border-gray-700 bg-gray-800 text-violet-500 focus:ring-violet-500"
              />
              <label htmlFor="badge-toggle" className="text-sm text-gray-400">Award badge</label>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
