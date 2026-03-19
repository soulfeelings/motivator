import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Trophy, Plus, Trash2, Play, Check } from 'lucide-react'

interface Tournament {
  id: string
  name: string
  description?: string
  season?: string
  metric: string
  prize_pool: number
  status: string
  participant_count: number
  starts_at: string
  ends_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const statusColors: Record<string, string> = {
  draft: 'bg-gray-500/20 text-gray-400',
  registration: 'bg-blue-500/20 text-blue-400',
  active: 'bg-violet-500/20 text-violet-400',
  completed: 'bg-emerald-500/20 text-emerald-400',
}

export default function Tournaments() {
  const [tournaments, setTournaments] = useState<Tournament[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [metric, setMetric] = useState('')
  const [season, setSeason] = useState('')
  const [prizePool, setPrizePool] = useState(500)
  const [startsAt, setStartsAt] = useState('')
  const [endsAt, setEndsAt] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<Tournament[]>(`/companies/${cid}/tournaments`)
        setTournaments(res ?? [])
      }
    } catch {} finally { setLoading(false) }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const t = await api.post<Tournament>(`/companies/${companyId}/tournaments`, {
        name, metric, season: season || undefined, prize_pool: prizePool,
        starts_at: new Date(startsAt).toISOString(), ends_at: new Date(endsAt).toISOString(),
      })
      setTournaments(prev => [t, ...prev])
      setCreating(false)
      setName(''); setMetric(''); setSeason('')
    } catch (err: any) { setError(err.message) }
  }

  async function updateStatus(id: string, status: string) {
    if (!companyId) return
    if (status === 'completed') {
      await api.post(`/companies/${companyId}/tournaments/${id}/complete`)
    } else {
      await api.patch(`/companies/${companyId}/tournaments/${id}/status`, { status })
    }
    load()
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Delete this tournament?')) return
    await api.delete(`/companies/${companyId}/tournaments/${id}`)
    setTournaments(prev => prev.filter(t => t.id !== id))
  }

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!companyId) return <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">Create a company first.</p></div>

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Trophy size={24} className="text-amber-400" />
          <h2 className="text-2xl font-bold text-white">Tournaments</h2>
        </div>
        <button onClick={() => setCreating(!creating)} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
          <Plus size={16} /> New Tournament
        </button>
      </div>

      {creating && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-xl">
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Name</label>
              <input value={name} onChange={e => setName(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="Q1 Sales Championship" required />
            </div>
            <div className="w-32">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Season</label>
              <input value={season} onChange={e => setSeason(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="Q1 2026" />
            </div>
          </div>
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Metric</label>
              <input value={metric} onChange={e => setMetric(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="deals_closed" required />
            </div>
            <div className="w-32">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Prize Pool</label>
              <input type="number" value={prizePool} onChange={e => setPrizePool(Number(e.target.value))} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" min={0} />
            </div>
          </div>
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Starts</label>
              <input type="datetime-local" value={startsAt} onChange={e => setStartsAt(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" required />
            </div>
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Ends</label>
              <input type="datetime-local" value={endsAt} onChange={e => setEndsAt(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" required />
            </div>
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg">Create</button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg">Cancel</button>
          </div>
        </form>
      )}

      {tournaments.length === 0 && !creating ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">No tournaments yet.</p></div>
      ) : (
        <div className="space-y-3">
          {tournaments.map(t => (
            <div key={t.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 group">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <h3 className="text-white font-medium">{t.name}</h3>
                  {t.season && <span className="text-xs text-gray-500">{t.season}</span>}
                  <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${statusColors[t.status] ?? statusColors.draft}`}>{t.status}</span>
                </div>
                <div className="flex items-center gap-2">
                  {t.status === 'draft' && (
                    <button onClick={() => updateStatus(t.id, 'registration')} className="p-1.5 text-blue-400 hover:bg-blue-500/20 rounded-lg" title="Open registration"><Play size={14} /></button>
                  )}
                  {t.status === 'registration' && (
                    <button onClick={() => updateStatus(t.id, 'active')} className="p-1.5 text-violet-400 hover:bg-violet-500/20 rounded-lg" title="Start"><Play size={14} /></button>
                  )}
                  {t.status === 'active' && (
                    <button onClick={() => updateStatus(t.id, 'completed')} className="p-1.5 text-emerald-400 hover:bg-emerald-500/20 rounded-lg" title="Complete"><Check size={14} /></button>
                  )}
                  <button onClick={() => handleDelete(t.id)} className="p-1.5 text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"><Trash2 size={14} /></button>
                </div>
              </div>
              <div className="flex gap-6 text-sm text-gray-500">
                <span className="text-violet-400 font-mono">{t.metric}</span>
                <span>{t.participant_count} participants</span>
                <span>{t.prize_pool} coin prize pool</span>
                <span>{new Date(t.starts_at).toLocaleDateString()} — {new Date(t.ends_at).toLocaleDateString()}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
