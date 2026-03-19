import type { LucideIcon } from 'lucide-react'

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description?: string
  action?: { label: string; onClick: () => void }
}

export default function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-12 text-center">
      {Icon && (
        <div className="w-12 h-12 mx-auto mb-4 rounded-full bg-gray-800 flex items-center justify-center">
          <Icon size={24} className="text-gray-500" />
        </div>
      )}
      <h3 className="text-white font-medium mb-1">{title}</h3>
      {description && <p className="text-sm text-gray-500 mb-4">{description}</p>}
      {action && (
        <button
          onClick={action.onClick}
          className="px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"
        >
          {action.label}
        </button>
      )}
    </div>
  )
}
