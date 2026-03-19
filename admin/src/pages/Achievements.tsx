import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Target, Plus, Trash2 } from 'lucide-react'

interface Badge {
  id: string
  name: string
}

interface Achievement {
  id: string
  name: string
  description?: string
  metric: string
  operator: string
  threshold: number
  badge_id?: string
  xp_reward: number
  coin_reward: number
  is_active: boolean
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const operatorLabels: Record<string, string> = {
  gte: '>=',
  lte: '<=',
  eq: '=',
  gt: '>',
  lt: '<',
}

export default function Achievements() {
  const [achievements, setAchievements] = useState<Achievement[]>([])
  const [badges, setBadges] = useState<Badge[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [metric, setMetric] = useState('')
  const [operator, setOperator] = useState('gte')
  const [threshold, setThreshold] = useState(1)
  const [badgeId, setBadgeId] = useState('')
  const [xpReward, setXpReward] = useState(25)
  const [coinReward, setCoinReward] = useState(10)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    load()
  }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const [achRes, badgeRes] = await Promise.all([
          api.get<Achievement[]>(`/companies/${cid}/achievements`),
          api.get<Badge[]>(`/companies/${cid}/badges`),
        ])
        setAchievements(achRes ?? [])
        setBadges(badgeRes ?? [])
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
      const a = await api.post<Achievement>(`/companies/${companyId}/achievements`, {
        name,
        description: description || undefined,
        metric,
        operator,
        threshold,
        badge_id: badgeId || undefined,
        xp_reward: xpReward,
        coin_reward: coinReward,
      })
      setAchievements((prev) => [...prev, a])
      setCreating(false)
      setName('')
      setDescription('')
      setMetric('')
      setThreshold(1)
      setBadgeId('')
    } catch (err: any) {
      setError(err.message)
    }
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Delete this achievement?')) return
    await api.delete(`/companies/${companyId}/achievements/${id}`)
    setAchievements((prev) => prev.filter((a) => a.id !== id))
  }

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!companyId) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
        <p className="text-gray-400">Create a company first.</p>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white">Achievements</h2>
        <button
          onClick={() => setCreating(!creating)}
          className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"
        >
          <Plus size={16} />
          New Achievement
        </button>
      </div>

      {creating && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-xl">
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
              placeholder="Close 10 deals in a month"
            />
          </div>
          <div className="flex gap-3">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Metric</label>
              <input
                value={metric}
                onChange={(e) => setMetric(e.target.value)}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                placeholder="deals_closed"
                required
              />
            </div>
            <div className="w-24">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Op</label>
              <select
                value={operator}
                onChange={(e) => setOperator(e.target.value)}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
              >
                <option value="gte">{"≥"}</option>
                <option value="gt">{">"}</option>
                <option value="eq">{"="}</option>
                <option value="lte">{"≤"}</option>
                <option value="lt">{"<"}</option>
              </select>
            </div>
            <div className="w-28">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Threshold</label>
              <input
                type="number"
                value={threshold}
                onChange={(e) => setThreshold(Number(e.target.value))}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                min={0}
                required
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Award Badge (optional)</label>
            <select
              value={badgeId}
              onChange={(e) => setBadgeId(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
            >
              <option value="">None</option>
              {badges.map((b) => (
                <option key={b.id} value={b.id}>{b.name}</option>
              ))}
            </select>
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
              Create
            </button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 font-medium rounded-lg transition-colors">
              Cancel
            </button>
          </div>
        </form>
      )}

      {achievements.length === 0 && !creating ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No achievements created yet. Define rules to gamify work.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {achievements.map((a) => (
            <div key={a.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 flex items-center justify-between group">
              <div className="flex items-center gap-4">
                <div className="w-10 h-10 rounded-full bg-emerald-600/20 flex items-center justify-center">
                  <Target size={20} className="text-emerald-400" />
                </div>
                <div>
                  <h3 className="text-white font-medium">{a.name}</h3>
                  <p className="text-sm text-gray-500">
                    <span className="text-violet-400 font-mono">{a.metric}</span>
                    {' '}{operatorLabels[a.operator] ?? a.operator}{' '}
                    <span className="text-white font-medium">{a.threshold}</span>
                  </p>
                  {a.description && <p className="text-xs text-gray-600 mt-0.5">{a.description}</p>}
                </div>
              </div>
              <div className="flex items-center gap-4">
                <div className="text-right text-xs space-y-0.5">
                  {a.xp_reward > 0 && <p className="text-emerald-400">+{a.xp_reward} XP</p>}
                  {a.coin_reward > 0 && <p className="text-amber-400">+{a.coin_reward} Coins</p>}
                </div>
                <button
                  onClick={() => handleDelete(a.id)}
                  className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"
                >
                  <Trash2 size={14} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
