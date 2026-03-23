import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Workflow, Plus, Trash2, Play, Pause, Pencil } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface GamePlan {
  id: string
  name: string
  description?: string
  is_active: boolean
  flow_data: { nodes: unknown[]; edges: unknown[] }
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function GamePlans() {
  const navigate = useNavigate()
  const [plans, setPlans] = useState<GamePlan[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [toast, setToast] = useState<{message: string, type: 'success'|'error'} | null>(null)

  useEffect(() => {
    load()
  }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<GamePlan[]>(`/companies/${cid}/game-plans`)
        setPlans(res ?? [])
      }
    } catch {
      // no company
    } finally {
      setLoading(false)
    }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const plan = await api.post<GamePlan>(`/companies/${companyId}/game-plans`, {
        name,
        description: description || undefined,
      })
      setToast({ message: 'Game plan created', type: 'success' })
      navigate(`/game-plans/${plan.id}`)
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create game plan', type: 'error' })
    }
  }

  async function toggleActive(plan: GamePlan) {
    if (!companyId) return
    try {
      const endpoint = plan.is_active ? 'deactivate' : 'activate'
      await api.post(`/companies/${companyId}/game-plans/${plan.id}/${endpoint}`)
      setPlans((prev) => prev.map((p) => (p.id === plan.id ? { ...p, is_active: !p.is_active } : p)))
      setToast({ message: plan.is_active ? 'Game plan deactivated' : 'Game plan activated', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to update game plan', type: 'error' })
    }
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Delete this game plan?')) return
    try {
      await api.delete(`/companies/${companyId}/game-plans/${id}`)
      setPlans((prev) => prev.filter((p) => p.id !== id))
      setToast({ message: 'Game plan deleted', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to delete game plan', type: 'error' })
    }
  }

  if (loading) return <CardSkeleton count={3} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Workflow size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">Game Plans</h2>
        </div>
        <EmptyState
          icon={Workflow}
          title="No company yet"
          description="Design visual gamification flows that connect triggers, rules, and rewards together."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Workflow size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">Game Plans</h2>
        </div>
        <button
          onClick={() => setCreating(!creating)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"
        >
          <Plus size={16} />
          New Game Plan
        </button>
      </div>

      {creating && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Name</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="Sales Gamification Q1"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Description</label>
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="Gamification rules for the sales team"
            />
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">
              Create & Edit
            </button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 font-medium rounded-lg transition-colors">
              Cancel
            </button>
          </div>
        </form>
      )}

      {plans.length === 0 && !creating ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No game plans yet. Create one to design your gamification flow.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {plans.map((plan) => (
            <div key={plan.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 flex items-center justify-between group">
              <div className="flex items-center gap-4">
                <div className={`w-3 h-3 rounded-full ${plan.is_active ? 'bg-emerald-400' : 'bg-gray-600'}`} />
                <div>
                  <h3 className="text-white font-medium">{plan.name}</h3>
                  <p className="text-xs text-gray-500">
                    {plan.flow_data?.nodes?.length ?? 0} nodes · {plan.flow_data?.edges?.length ?? 0} connections
                    {plan.description && ` · ${plan.description}`}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => navigate(`/game-plans/${plan.id}`)}
                  className="p-2 text-gray-500 hover:text-violet-400 transition-colors"
                  title="Edit flow"
                >
                  <Pencil size={16} />
                </button>
                <button
                  onClick={() => toggleActive(plan)}
                  className={`p-2 transition-colors ${plan.is_active ? 'text-emerald-400 hover:text-amber-400' : 'text-gray-500 hover:text-emerald-400'}`}
                  title={plan.is_active ? 'Deactivate' : 'Activate'}
                >
                  {plan.is_active ? <Pause size={16} /> : <Play size={16} />}
                </button>
                <button
                  onClick={() => handleDelete(plan.id)}
                  className="p-2 text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"
                  title="Delete"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
