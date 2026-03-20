import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Heart, Plus, Trash2, Play, Eye, Vote, Check } from 'lucide-react'

interface Quest {
  id: string
  name: string
  description?: string
  status: string
  xp_reward: number
  coin_reward: number
  bonus_xp: number
  bonus_coins: number
  pair_count: number
  sent_count: number
  deadline: string
  created_at: string
}

interface QuestPair {
  id: string
  sender_id: string
  receiver_id: string
  message?: string
  sent_at?: string
  vote_count: number
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const statusColors: Record<string, string> = {
  draft: 'bg-gray-500/20 text-gray-400',
  active: 'bg-emerald-500/20 text-emerald-400',
  voting: 'bg-amber-500/20 text-amber-400',
  revealed: 'bg-violet-500/20 text-violet-400',
  completed: 'bg-blue-500/20 text-blue-400',
}

export default function Quests() {
  const [quests, setQuests] = useState<Quest[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('Secret Motivator')
  const [description, setDescription] = useState('Send an anonymous positive message to a random colleague!')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [pairs, setPairs] = useState<QuestPair[]>([])

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<Quest[]>(`/companies/${cid}/quests`)
        setQuests(res ?? [])
      }
    } catch {} finally { setLoading(false) }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const q = await api.post<Quest>(`/companies/${companyId}/quests`, {
        name, description: description || undefined,
      })
      setQuests(prev => [q, ...prev])
      setCreating(false)
    } catch (err: any) { setError(err.message) }
  }

  async function handleAction(questId: string, action: string) {
    if (!companyId) return
    try {
      await api.post(`/companies/${companyId}/quests/${questId}/${action}`)
      load()
    } catch (err: any) { setError(err.message) }
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Delete this quest?')) return
    await api.delete(`/companies/${companyId}/quests/${id}`)
    setQuests(prev => prev.filter(q => q.id !== id))
    if (selectedId === id) { setSelectedId(null); setPairs([]) }
  }

  async function loadPairs(questId: string) {
    if (!companyId) return
    setSelectedId(questId)
    const res = await api.get<QuestPair[]>(`/companies/${companyId}/quests/${questId}/pairs`)
    setPairs(res ?? [])
  }

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!companyId) return <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">Create a company first.</p></div>

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Heart size={24} className="text-pink-400" />
          <h2 className="text-2xl font-bold text-white">Secret Motivator</h2>
        </div>
        <button onClick={() => setCreating(!creating)} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
          <Plus size={16} /> New Quest
        </button>
      </div>

      {error && <p className="text-sm text-red-400 mb-4">{error}</p>}

      {creating && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Quest Name</label>
            <input value={name} onChange={e => setName(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Description</label>
            <textarea value={description} onChange={e => setDescription(e.target.value)} rows={2} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" />
          </div>
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg">Create</button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg">Cancel</button>
          </div>
        </form>
      )}

      <div className="flex gap-6">
        <div className="flex-1 space-y-3">
          {quests.length === 0 && !creating ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
              <Heart size={32} className="text-pink-400/30 mx-auto mb-3" />
              <p className="text-gray-400">No quests yet. Create one to spread positivity!</p>
            </div>
          ) : quests.map(q => (
            <div key={q.id} onClick={() => loadPairs(q.id)}
              className={`bg-gray-900 border rounded-xl p-5 cursor-pointer group transition-colors ${selectedId === q.id ? 'border-violet-500/50' : 'border-gray-800 hover:border-gray-700'}`}>
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <h3 className="text-white font-medium">{q.name}</h3>
                  <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${statusColors[q.status] ?? statusColors.draft}`}>{q.status}</span>
                </div>
                <div className="flex items-center gap-1.5">
                  {q.status === 'draft' && (
                    <button onClick={e => { e.stopPropagation(); handleAction(q.id, 'start') }} className="p-1.5 text-emerald-400 hover:bg-emerald-500/20 rounded-lg" title="Start (pair members)"><Play size={14} /></button>
                  )}
                  {q.status === 'active' && (
                    <button onClick={e => { e.stopPropagation(); handleAction(q.id, 'voting') }} className="p-1.5 text-amber-400 hover:bg-amber-500/20 rounded-lg" title="Start voting"><Vote size={14} /></button>
                  )}
                  {q.status === 'voting' && (
                    <button onClick={e => { e.stopPropagation(); handleAction(q.id, 'reveal') }} className="p-1.5 text-violet-400 hover:bg-violet-500/20 rounded-lg" title="Reveal senders"><Eye size={14} /></button>
                  )}
                  {q.status === 'revealed' && (
                    <button onClick={e => { e.stopPropagation(); handleAction(q.id, 'complete') }} className="p-1.5 text-blue-400 hover:bg-blue-500/20 rounded-lg" title="Complete & award"><Check size={14} /></button>
                  )}
                  <button onClick={e => { e.stopPropagation(); handleDelete(q.id) }} className="p-1.5 text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100"><Trash2 size={14} /></button>
                </div>
              </div>
              <div className="flex gap-4 text-xs text-gray-500">
                <span>{q.pair_count} pairs</span>
                <span>{q.sent_count} messages sent</span>
                <span>+{q.xp_reward} XP / +{q.coin_reward} coins</span>
                <span>Bonus: +{q.bonus_xp} XP</span>
              </div>
            </div>
          ))}
        </div>

        {selectedId && pairs.length > 0 && (
          <div className="w-80 bg-gray-900 border border-gray-800 rounded-xl p-5 h-fit">
            <h3 className="text-white font-medium mb-4">Pairs & Messages</h3>
            <div className="space-y-3 max-h-96 overflow-auto">
              {pairs.map(p => (
                <div key={p.id} className="p-3 bg-gray-800/50 rounded-lg">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs text-gray-500">{p.sender_id.slice(0, 8)} → {p.receiver_id.slice(0, 8)}</span>
                    {p.vote_count > 0 && <span className="text-xs text-pink-400">{p.vote_count} votes</span>}
                  </div>
                  {p.message ? (
                    <p className="text-sm text-white italic">"{p.message}"</p>
                  ) : (
                    <p className="text-xs text-gray-600">Not sent yet</p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
