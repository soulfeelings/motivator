import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Gamepad2, Plus, Play, Vote, Check, Brain, Camera, MessageCircle } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface SocialGame {
  id: string
  name: string
  description?: string
  game_type: string
  status: string
  duration_hours: number
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const statusColors: Record<string, string> = {
  draft: 'bg-gray-500/20 text-gray-400',
  active: 'bg-emerald-500/20 text-emerald-400',
  voting: 'bg-violet-500/20 text-violet-400',
  completed: 'bg-blue-500/20 text-blue-400',
}

const typeColors: Record<string, string> = {
  trivia: 'bg-amber-500/20 text-amber-400',
  photo_challenge: 'bg-pink-500/20 text-pink-400',
  two_truths: 'bg-cyan-500/20 text-cyan-400',
}

const typeLabels: Record<string, string> = {
  trivia: 'Trivia',
  photo_challenge: 'Photo Challenge',
  two_truths: 'Two Truths & a Lie',
}

const gameTypes = [
  { type: 'trivia', icon: Brain, name: 'Trivia', description: 'Test team knowledge with fun quiz questions' },
  { type: 'photo_challenge', icon: Camera, name: 'Photo Challenge', description: 'Share photos around a theme and vote for favorites' },
  { type: 'two_truths', icon: MessageCircle, name: 'Two Truths & a Lie', description: 'Guess which statement is the lie about your teammates' },
]

export default function SocialGames() {
  const navTo = useNavigate()
  const [games, setGames] = useState<SocialGame[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [picking, setPicking] = useState(false)
  const [selectedType, setSelectedType] = useState<string | null>(null)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [duration, setDuration] = useState('24')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null)

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<SocialGame[]>(`/companies/${cid}/social-games`)
        setGames(res ?? [])
      }
    } catch {} finally { setLoading(false) }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId || !selectedType) return
    setError('')
    try {
      const g = await api.post<SocialGame>(`/companies/${companyId}/social-games`, {
        name,
        description: description || undefined,
        game_type: selectedType,
        duration_hours: parseInt(duration),
      })
      setGames(prev => [g, ...prev])
      setPicking(false)
      setSelectedType(null)
      setName('')
      setDescription('')
      setDuration('24')
      setToast({ message: 'Social game created successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create game', type: 'error' })
    }
  }

  async function handleAction(gameId: string, action: string) {
    if (!companyId) return
    try {
      await api.post(`/companies/${companyId}/social-games/${gameId}/${action}`)
      load()
      setToast({ message: `Game ${action} successful`, type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || `Failed to ${action} game`, type: 'error' })
    }
  }

  if (loading) return <CardSkeleton count={3} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Gamepad2 size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">Social Games</h2>
        </div>
        <EmptyState
          icon={Gamepad2}
          title="No company yet"
          description="Create social games like trivia, photo challenges, and two truths & a lie to engage your team."
          action={{ label: 'Create a Company', onClick: () => navTo('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Gamepad2 size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">Social Games</h2>
        </div>
        <button onClick={() => { setPicking(!picking); setSelectedType(null) }} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
          <Plus size={16} /> New Game
        </button>
      </div>

      {error && <p className="text-sm text-red-400 mb-4">{error}</p>}

      {picking && !selectedType && (
        <div className="mb-6">
          <h3 className="text-white font-medium mb-3">Choose a game type</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {gameTypes.map(gt => {
              const Icon = gt.icon
              return (
                <button
                  key={gt.type}
                  onClick={() => setSelectedType(gt.type)}
                  className="bg-gray-900 border border-gray-800 hover:border-violet-500/50 rounded-xl p-5 text-left transition-colors"
                >
                  <div className={`w-10 h-10 rounded-lg flex items-center justify-center mb-3 ${typeColors[gt.type]}`}>
                    <Icon size={20} />
                  </div>
                  <h4 className="text-white font-medium mb-1">{gt.name}</h4>
                  <p className="text-sm text-gray-500">{gt.description}</p>
                </button>
              )
            })}
          </div>
          <button onClick={() => setPicking(false)} className="mt-3 text-sm text-gray-500 hover:text-gray-300">Cancel</button>
        </div>
      )}

      {picking && selectedType && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
          <div className="flex items-center gap-2 mb-2">
            <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${typeColors[selectedType]}`}>{typeLabels[selectedType]}</span>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Game Name</label>
            <input value={name} onChange={e => setName(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="Friday Trivia" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Description</label>
            <textarea value={description} onChange={e => setDescription(e.target.value)} rows={2} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="A fun game for the team..." />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Duration</label>
            <select value={duration} onChange={e => setDuration(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500">
              <option value="24">1 day</option>
              <option value="72">3 days</option>
              <option value="168">1 week</option>
            </select>
          </div>
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">Create</button>
            <button type="button" onClick={() => { setSelectedType(null); setPicking(false) }} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg transition-colors">Cancel</button>
          </div>
        </form>
      )}

      <div className="space-y-3">
        {games.length === 0 && !picking ? (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <Gamepad2 size={32} className="text-violet-400/30 mx-auto mb-3" />
            <p className="text-gray-400">No social games yet. Create one to engage your team!</p>
          </div>
        ) : games.map(g => (
          <div key={g.id} className="bg-gray-900 border border-gray-800 hover:border-gray-700 rounded-xl p-5 group transition-colors">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <Link to={`/social-games/${g.id}`} className="text-white font-medium hover:text-violet-400 transition-colors">{g.name}</Link>
                <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${typeColors[g.game_type] ?? typeColors.trivia}`}>{typeLabels[g.game_type] ?? g.game_type}</span>
                <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${statusColors[g.status] ?? statusColors.draft}`}>{g.status}</span>
              </div>
              <div className="flex items-center gap-1.5">
                {g.status === 'draft' && (
                  <button onClick={() => handleAction(g.id, 'launch')} className="p-1.5 text-emerald-400 hover:bg-emerald-500/20 rounded-lg" title="Launch game"><Play size={14} /></button>
                )}
                {g.status === 'active' && (g.game_type === 'photo_challenge' || g.game_type === 'two_truths') && (
                  <button onClick={() => handleAction(g.id, 'start-voting')} className="p-1.5 text-violet-400 hover:bg-violet-500/20 rounded-lg" title="Start voting"><Vote size={14} /></button>
                )}
                {(g.status === 'active' || g.status === 'voting') && (
                  <button onClick={() => handleAction(g.id, 'complete')} className="p-1.5 text-blue-400 hover:bg-blue-500/20 rounded-lg" title="Complete"><Check size={14} /></button>
                )}
              </div>
            </div>
            <div className="flex gap-4 text-xs text-gray-500">
              <span>{new Date(g.created_at).toLocaleDateString()}</span>
              <span>{g.duration_hours}h duration</span>
            </div>
          </div>
        ))}
      </div>

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
