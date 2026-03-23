import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Award, Plus, Trash2 } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface Badge {
  id: string
  name: string
  description?: string
  icon_url?: string
  xp_reward: number
  coin_reward: number
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function Badges() {
  const navigate = useNavigate()
  const [badges, setBadges] = useState<Badge[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [xpReward, setXpReward] = useState(10)
  const [coinReward, setCoinReward] = useState(5)
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
        const res = await api.get<Badge[]>(`/companies/${cid}/badges`)
        setBadges(res ?? [])
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
      const badge = await api.post<Badge>(`/companies/${companyId}/badges`, {
        name,
        description: description || undefined,
        xp_reward: xpReward,
        coin_reward: coinReward,
      })
      setBadges((prev) => [...prev, badge])
      setCreating(false)
      setName('')
      setDescription('')
      setXpReward(10)
      setCoinReward(5)
      setToast({ message: 'Badge created successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create badge', type: 'error' })
    }
  }

  async function handleDelete(badgeId: string) {
    if (!companyId || !confirm('Delete this badge?')) return
    try {
      await api.delete(`/companies/${companyId}/badges/${badgeId}`)
      setBadges((prev) => prev.filter((b) => b.id !== badgeId))
      setToast({ message: 'Badge deleted', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to delete badge', type: 'error' })
    }
  }

  if (loading) return <CardSkeleton count={3} />
  if (!companyId) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-white mb-6">Badges</h2>
        <EmptyState
          icon={Award}
          title="No company yet"
          description="Create badges to reward your team for hitting milestones and going above and beyond."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white">Badges</h2>
        <button
          onClick={() => setCreating(!creating)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"
        >
          <Plus size={16} />
          New Badge
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
              placeholder="Gold Closer"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Description</label>
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="Awarded for closing 10 deals"
            />
          </div>
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">XP Reward</label>
              <input
                type="number"
                value={xpReward}
                onChange={(e) => setXpReward(Number(e.target.value))}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                min={0}
              />
            </div>
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Coin Reward</label>
              <input
                type="number"
                value={coinReward}
                onChange={(e) => setCoinReward(Number(e.target.value))}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                min={0}
              />
            </div>
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">
              Create Badge
            </button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 font-medium rounded-lg transition-colors">
              Cancel
            </button>
          </div>
        </form>
      )}

      {badges.length === 0 && !creating ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No badges created yet.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {badges.map((badge) => (
            <div key={badge.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 relative group">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-10 h-10 rounded-full bg-violet-600/20 flex items-center justify-center">
                    <Award size={20} className="text-violet-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-medium">{badge.name}</h3>
                    {badge.description && <p className="text-xs text-gray-500">{badge.description}</p>}
                  </div>
                </div>
                <button
                  onClick={() => handleDelete(badge.id)}
                  className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"
                  title="Delete badge"
                >
                  <Trash2 size={14} />
                </button>
              </div>
              <div className="flex gap-4 mt-2">
                <span className="text-xs text-emerald-400">+{badge.xp_reward} XP</span>
                <span className="text-xs text-amber-400">+{badge.coin_reward} Coins</span>
              </div>
            </div>
          ))}
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
