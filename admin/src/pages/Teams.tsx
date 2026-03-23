import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Users2, Plus, Trash2, Swords } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface Team {
  id: string
  name: string
  description?: string
  color: string
  member_count: number
  created_at: string
}

interface TeamBattle {
  id: string
  team_a_id: string
  team_b_id: string
  metric: string
  target: number
  team_a_score: number
  team_b_score: number
  status: string
  winner_id?: string
  xp_reward: number
  coin_reward: number
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function Teams() {
  const navigate = useNavigate()
  const [teams, setTeams] = useState<Team[]>([])
  const [battles, setBattles] = useState<TeamBattle[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [creatingBattle, setCreatingBattle] = useState(false)
  const [name, setName] = useState('')
  const [color, setColor] = useState('#8b5cf6')
  const [battleTeamA, setBattleTeamA] = useState('')
  const [battleTeamB, setBattleTeamB] = useState('')
  const [battleMetric, setBattleMetric] = useState('')
  const [battleTarget, setBattleTarget] = useState(100)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [tab, setTab] = useState<'teams' | 'battles'>('teams')
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
        const [t, b] = await Promise.all([
          api.get<Team[]>(`/companies/${cid}/teams`),
          api.get<TeamBattle[]>(`/companies/${cid}/team-battles`),
        ])
        setTeams(t ?? [])
        setBattles(b ?? [])
      }
    } catch {
      // no company
    } finally {
      setLoading(false)
    }
  }

  async function handleCreateTeam(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const team = await api.post<Team>(`/companies/${companyId}/teams`, { name, color })
      setTeams((prev) => [...prev, team])
      setCreating(false)
      setName('')
      setToast({ message: 'Team created successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create team', type: 'error' })
    }
  }

  async function handleCreateBattle(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const battle = await api.post<TeamBattle>(`/companies/${companyId}/team-battles`, {
        team_a_id: battleTeamA,
        team_b_id: battleTeamB,
        metric: battleMetric,
        target: battleTarget,
      })
      setBattles((prev) => [battle, ...prev])
      setCreatingBattle(false)
      setToast({ message: 'Battle created successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create battle', type: 'error' })
    }
  }

  async function handleDeleteTeam(id: string) {
    if (!companyId || !confirm('Delete this team?')) return
    try {
      await api.delete(`/companies/${companyId}/teams/${id}`)
      setTeams((prev) => prev.filter((t) => t.id !== id))
      setToast({ message: 'Team deleted', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to delete team', type: 'error' })
    }
  }

  const teamName = (id: string) => teams.find((t) => t.id === id)?.name ?? id.slice(0, 8)

  if (loading) return <CardSkeleton count={3} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Users2 size={24} className="text-blue-400" />
          <h2 className="text-2xl font-bold text-white">Teams</h2>
        </div>
        <EmptyState
          icon={Users2}
          title="No company yet"
          description="Organize employees into teams and launch head-to-head team battles."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Users2 size={24} className="text-blue-400" />
          <h2 className="text-2xl font-bold text-white">Teams</h2>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex bg-gray-800 rounded-lg p-0.5">
            <button onClick={() => setTab('teams')} className={`px-3 py-1.5 text-sm rounded-md transition-colors ${tab === 'teams' ? 'bg-gray-700 text-white' : 'text-gray-400'}`}>Teams</button>
            <button onClick={() => setTab('battles')} className={`px-3 py-1.5 text-sm rounded-md transition-colors ${tab === 'battles' ? 'bg-gray-700 text-white' : 'text-gray-400'}`}>Battles ({battles.length})</button>
          </div>
          {tab === 'teams' && (
            <button onClick={() => setCreating(!creating)} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
              <Plus size={16} /> New Team
            </button>
          )}
          {tab === 'battles' && teams.length >= 2 && (
            <button onClick={() => setCreatingBattle(!creatingBattle)} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
              <Swords size={16} /> New Battle
            </button>
          )}
        </div>
      </div>

      {error && <p className="text-sm text-red-400 mb-4">{error}</p>}

      {tab === 'teams' && (
        <>
          {creating && (
            <form onSubmit={handleCreateTeam} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
              <div className="flex gap-4">
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Team Name</label>
                  <input value={name} onChange={(e) => setName(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="Alpha Team" required />
                </div>
                <div className="w-24">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Color</label>
                  <input type="color" value={color} onChange={(e) => setColor(e.target.value)} className="w-full h-[42px] bg-gray-800 border border-gray-700 rounded-lg cursor-pointer" />
                </div>
              </div>
              <div className="flex gap-3">
                <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg">Create</button>
                <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg">Cancel</button>
              </div>
            </form>
          )}

          {teams.length === 0 && !creating ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">No teams yet.</p></div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {teams.map((t) => (
                <div key={t.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 relative group">
                  <div className="flex items-center gap-3 mb-2">
                    <div className="w-4 h-4 rounded-full" style={{ backgroundColor: t.color }} />
                    <h3 className="text-white font-medium">{t.name}</h3>
                    <button onClick={() => handleDeleteTeam(t.id)} className="ml-auto text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"><Trash2 size={14} /></button>
                  </div>
                  <p className="text-sm text-gray-500">{t.member_count} members</p>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {tab === 'battles' && (
        <>
          {creatingBattle && (
            <form onSubmit={handleCreateBattle} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
              <div className="flex gap-4">
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Team A</label>
                  <select value={battleTeamA} onChange={(e) => setBattleTeamA(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white" required>
                    <option value="">Select...</option>
                    {teams.map((t) => <option key={t.id} value={t.id}>{t.name}</option>)}
                  </select>
                </div>
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Team B</label>
                  <select value={battleTeamB} onChange={(e) => setBattleTeamB(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white" required>
                    <option value="">Select...</option>
                    {teams.filter((t) => t.id !== battleTeamA).map((t) => <option key={t.id} value={t.id}>{t.name}</option>)}
                  </select>
                </div>
              </div>
              <div className="flex gap-4">
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Metric</label>
                  <input value={battleMetric} onChange={(e) => setBattleMetric(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white" placeholder="deals_closed" required />
                </div>
                <div className="w-28">
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">Target</label>
                  <input type="number" value={battleTarget} onChange={(e) => setBattleTarget(Number(e.target.value))} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white" min={1} />
                </div>
              </div>
              <div className="flex gap-3">
                <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg">Create Battle</button>
                <button type="button" onClick={() => setCreatingBattle(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg">Cancel</button>
              </div>
            </form>
          )}

          {battles.length === 0 && !creatingBattle ? (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">No team battles yet.</p></div>
          ) : (
            <div className="space-y-3">
              {battles.map((b) => (
                <div key={b.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                  <div className="flex items-center justify-between mb-3">
                    <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${
                      b.status === 'completed' ? 'bg-emerald-500/20 text-emerald-400' :
                      b.status === 'active' ? 'bg-violet-500/20 text-violet-400' : 'bg-amber-500/20 text-amber-400'
                    }`}>{b.status}</span>
                    <span className="text-xs text-gray-500">{b.metric} | target: {b.target}</span>
                  </div>
                  <div className="flex items-center gap-6">
                    <div className="flex-1 text-center">
                      <p className="text-sm text-gray-400">{teamName(b.team_a_id)}</p>
                      <p className="text-2xl font-bold text-white">{b.team_a_score}</p>
                    </div>
                    <span className="text-gray-600 font-bold">VS</span>
                    <div className="flex-1 text-center">
                      <p className="text-sm text-gray-400">{teamName(b.team_b_id)}</p>
                      <p className="text-2xl font-bold text-white">{b.team_b_score}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
