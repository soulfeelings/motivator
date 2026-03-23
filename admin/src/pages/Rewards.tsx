import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Gift, Plus, Trash2, Check, Clock } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface Reward {
  id: string
  name: string
  description?: string
  cost_coins: number
  stock?: number
  is_active: boolean
}

interface Redemption {
  id: string
  membership_id: string
  reward_id: string
  coins_spent: number
  status: string
  redeemed_at: string
  reward?: Reward
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function Rewards() {
  const navigate = useNavigate()
  const [rewards, setRewards] = useState<Reward[]>([])
  const [redemptions, setRedemptions] = useState<Redemption[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [costCoins, setCostCoins] = useState(100)
  const [stock, setStock] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [tab, setTab] = useState<'rewards' | 'redemptions'>('rewards')
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
        const [rw, rd] = await Promise.all([
          api.get<Reward[]>(`/companies/${cid}/rewards`),
          api.get<Redemption[]>(`/companies/${cid}/rewards/redemptions`),
        ])
        setRewards(rw ?? [])
        setRedemptions(rd ?? [])
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
      const rw = await api.post<Reward>(`/companies/${companyId}/rewards`, {
        name,
        description: description || undefined,
        cost_coins: costCoins,
        stock: stock ? Number(stock) : undefined,
      })
      setRewards((prev) => [...prev, rw])
      setCreating(false)
      setName('')
      setDescription('')
      setCostCoins(100)
      setStock('')
      setToast({ message: 'Reward created successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create reward', type: 'error' })
    }
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Remove this reward?')) return
    try {
      await api.delete(`/companies/${companyId}/rewards/${id}`)
      setRewards((prev) => prev.filter((r) => r.id !== id))
      setToast({ message: 'Reward removed', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to remove reward', type: 'error' })
    }
  }

  async function handleFulfill(redemptionId: string) {
    if (!companyId) return
    try {
      await api.post(`/companies/${companyId}/rewards/redemptions/${redemptionId}/fulfill`)
      setRedemptions((prev) =>
        prev.map((r) => (r.id === redemptionId ? { ...r, status: 'fulfilled' } : r))
      )
      setToast({ message: 'Redemption fulfilled', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to fulfill redemption', type: 'error' })
    }
  }

  if (loading) return <CardSkeleton count={3} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Gift size={24} className="text-amber-400" />
          <h2 className="text-2xl font-bold text-white">Reward Store</h2>
        </div>
        <EmptyState
          icon={Gift}
          title="No company yet"
          description="Set up a reward store where employees spend earned coins on real perks and prizes."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Gift size={24} className="text-amber-400" />
          <h2 className="text-2xl font-bold text-white">Reward Store</h2>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex bg-gray-800 rounded-lg p-0.5">
            <button
              onClick={() => setTab('rewards')}
              className={`px-3 py-1.5 text-sm rounded-md transition-colors ${tab === 'rewards' ? 'bg-gray-700 text-white' : 'text-gray-400'}`}
            >
              Rewards
            </button>
            <button
              onClick={() => setTab('redemptions')}
              className={`px-3 py-1.5 text-sm rounded-md transition-colors ${tab === 'redemptions' ? 'bg-gray-700 text-white' : 'text-gray-400'}`}
            >
              Redemptions ({redemptions.length})
            </button>
          </div>
          {tab === 'rewards' && (
            <button
              onClick={() => setCreating(!creating)}
              className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"
            >
              <Plus size={16} />
              New Reward
            </button>
          )}
        </div>
      </div>

      {tab === 'rewards' && (
        <>
          {creating && (
            <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1.5">Name</label>
                <input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                  placeholder="Extra PTO Day"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1.5">Description</label>
                <input
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                  placeholder="One extra paid day off"
                />
              </div>
              <div className="flex gap-4">
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Cost (coins)</label>
                  <input
                    type="number"
                    value={costCoins}
                    onChange={(e) => setCostCoins(Number(e.target.value))}
                    className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                    min={1}
                    required
                  />
                </div>
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Stock (empty = unlimited)</label>
                  <input
                    type="number"
                    value={stock}
                    onChange={(e) => setStock(e.target.value)}
                    className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                    placeholder="∞"
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

          {rewards.length === 0 && !creating ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
              <p className="text-gray-400">No rewards in the store yet.</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {rewards.map((rw) => (
                <div key={rw.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 relative group">
                  <div className="flex items-start justify-between mb-2">
                    <h3 className="text-white font-medium">{rw.name}</h3>
                    <button
                      onClick={() => handleDelete(rw.id)}
                      className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"
                    >
                      <Trash2 size={14} />
                    </button>
                  </div>
                  {rw.description && <p className="text-xs text-gray-500 mb-3">{rw.description}</p>}
                  <div className="flex items-center justify-between">
                    <span className="text-amber-400 font-bold">{rw.cost_coins} coins</span>
                    <span className="text-xs text-gray-500">
                      {rw.stock != null ? `${rw.stock} left` : 'Unlimited'}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {tab === 'redemptions' && (
        <>
          {redemptions.length === 0 ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
              <p className="text-gray-400">No redemptions yet.</p>
            </div>
          ) : (
            <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-800">
                    <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase">Reward</th>
                    <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase">Member</th>
                    <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase">Coins</th>
                    <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase">Date</th>
                    <th className="px-5 py-3"></th>
                  </tr>
                </thead>
                <tbody>
                  {redemptions.map((rd) => (
                    <tr key={rd.id} className="border-b border-gray-800/50 hover:bg-gray-800/30">
                      <td className="px-5 py-4 text-sm text-white">{rd.reward?.name ?? rd.reward_id.slice(0, 8)}</td>
                      <td className="px-5 py-4 text-sm text-gray-400">{rd.membership_id.slice(0, 8)}</td>
                      <td className="px-5 py-4 text-sm text-amber-400">{rd.coins_spent}</td>
                      <td className="px-5 py-4">
                        <span className={`inline-flex items-center gap-1 text-xs font-medium ${
                          rd.status === 'fulfilled' ? 'text-emerald-400' :
                          rd.status === 'pending' ? 'text-amber-400' : 'text-gray-400'
                        }`}>
                          {rd.status === 'fulfilled' ? <Check size={12} /> : <Clock size={12} />}
                          {rd.status}
                        </span>
                      </td>
                      <td className="px-5 py-4 text-sm text-gray-500">{new Date(rd.redeemed_at).toLocaleDateString()}</td>
                      <td className="px-5 py-4 text-right">
                        {rd.status === 'pending' && (
                          <button
                            onClick={() => handleFulfill(rd.id)}
                            className="text-xs px-3 py-1 bg-emerald-600/20 text-emerald-400 rounded-lg hover:bg-emerald-600/30 transition-colors"
                          >
                            Fulfill
                          </button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
