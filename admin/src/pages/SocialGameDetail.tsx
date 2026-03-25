import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api } from '../lib/api'
import { ArrowLeft, Gamepad2, Play, Vote, Check, Plus, Trash2 } from 'lucide-react'
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

interface TriviaQuestion {
  id: string
  question: string
  options: string[]
  correct_index: number
}

interface Submission {
  id: string
  user_id: string
  content: string
  vote_count: number
  created_at: string
}

interface LeaderboardEntry {
  user_id: string
  score: number
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

export default function SocialGameDetail() {
  const { gameId } = useParams<{ gameId: string }>()
  const [game, setGame] = useState<SocialGame | null>(null)
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [tab, setTab] = useState<'content' | 'results'>('content')
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null)

  // Trivia state
  const [questions, setQuestions] = useState<TriviaQuestion[]>([])
  const [newQuestion, setNewQuestion] = useState('')
  const [newOptions, setNewOptions] = useState(['', '', '', ''])
  const [correctIndex, setCorrectIndex] = useState(0)
  const [addingQuestion, setAddingQuestion] = useState(false)

  // Submissions state
  const [submissions, setSubmissions] = useState<Submission[]>([])

  // Results state
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([])
  const [participantCount, setParticipantCount] = useState(0)

  useEffect(() => { load() }, [gameId])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const g = await api.get<SocialGame>(`/companies/${cid}/social-games/${gameId}`)
        setGame(g)

        if (g.game_type === 'trivia') {
          try {
            const qs = await api.get<TriviaQuestion[]>(`/companies/${cid}/social-games/${gameId}/questions`)
            setQuestions(qs ?? [])
          } catch {}
        }

        if (g.game_type === 'photo_challenge' || g.game_type === 'two_truths') {
          if (g.status === 'active' || g.status === 'voting') {
            try {
              const subs = await api.get<Submission[]>(`/companies/${cid}/social-games/${gameId}/submissions`)
              setSubmissions(subs ?? [])
            } catch {}
          }
        }

        if (g.status === 'completed') {
          try {
            const res = await api.get<{ leaderboard: LeaderboardEntry[]; participant_count: number }>(`/companies/${cid}/social-games/${gameId}/results`)
            setLeaderboard(res.leaderboard ?? [])
            setParticipantCount(res.participant_count ?? 0)
          } catch {}
          setTab('results')
        }
      }
    } catch {} finally { setLoading(false) }
  }

  async function handleAction(action: string) {
    if (!companyId || !gameId) return
    try {
      await api.post(`/companies/${companyId}/social-games/${gameId}/${action}`)
      load()
      setToast({ message: `Game ${action} successful`, type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || `Failed to ${action} game`, type: 'error' })
    }
  }

  async function handleAddQuestion(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId || !gameId) return
    try {
      const q = await api.post<TriviaQuestion>(`/companies/${companyId}/social-games/${gameId}/questions`, {
        question: newQuestion,
        options: newOptions,
        correct_index: correctIndex,
      })
      setQuestions(prev => [...prev, q])
      setNewQuestion('')
      setNewOptions(['', '', '', ''])
      setCorrectIndex(0)
      setAddingQuestion(false)
      setToast({ message: 'Question added', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to add question', type: 'error' })
    }
  }

  async function handleDeleteQuestion(qId: string) {
    if (!companyId || !gameId || !confirm('Delete this question?')) return
    try {
      await api.delete(`/companies/${companyId}/social-games/${gameId}/questions/${qId}`)
      setQuestions(prev => prev.filter(q => q.id !== qId))
      setToast({ message: 'Question deleted', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to delete question', type: 'error' })
    }
  }

  if (loading) return <CardSkeleton count={3} />
  if (!game) {
    return (
      <div>
        <Link to="/social-games" className="text-sm text-gray-500 hover:text-gray-300 flex items-center gap-1 mb-4"><ArrowLeft size={14} /> Back to Social Games</Link>
        <p className="text-gray-400">Game not found.</p>
      </div>
    )
  }

  return (
    <div>
      <Link to="/social-games" className="text-sm text-gray-500 hover:text-gray-300 flex items-center gap-1 mb-4"><ArrowLeft size={14} /> Back to Social Games</Link>

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Gamepad2 size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">{game.name}</h2>
          <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${typeColors[game.game_type] ?? typeColors.trivia}`}>{typeLabels[game.game_type] ?? game.game_type}</span>
          <span className={`text-xs font-bold uppercase px-2.5 py-1 rounded-full ${statusColors[game.status] ?? statusColors.draft}`}>{game.status}</span>
        </div>
        <div className="flex items-center gap-2">
          {game.status === 'draft' && (
            <button onClick={() => handleAction('launch')} className="inline-flex items-center gap-2 px-3 py-2 bg-emerald-600 hover:bg-emerald-500 text-white text-sm font-medium rounded-lg transition-colors"><Play size={14} /> Launch</button>
          )}
          {game.status === 'active' && (game.game_type === 'photo_challenge' || game.game_type === 'two_truths') && (
            <button onClick={() => handleAction('start-voting')} className="inline-flex items-center gap-2 px-3 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"><Vote size={14} /> Start Voting</button>
          )}
          {(game.status === 'active' || game.status === 'voting') && (
            <button onClick={() => handleAction('complete')} className="inline-flex items-center gap-2 px-3 py-2 bg-blue-600 hover:bg-blue-500 text-white text-sm font-medium rounded-lg transition-colors"><Check size={14} /> Complete</button>
          )}
        </div>
      </div>

      {game.description && <p className="text-gray-400 mb-6">{game.description}</p>}

      {game.status === 'completed' && (
        <div className="flex gap-2 mb-6">
          <button onClick={() => setTab('content')} className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors ${tab === 'content' ? 'bg-violet-600/20 text-violet-400' : 'text-gray-400 hover:text-gray-200'}`}>Content</button>
          <button onClick={() => setTab('results')} className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors ${tab === 'results' ? 'bg-violet-600/20 text-violet-400' : 'text-gray-400 hover:text-gray-200'}`}>Results</button>
        </div>
      )}

      {tab === 'content' && game.game_type === 'trivia' && (
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-white font-medium">Questions ({questions.length})</h3>
            {game.status === 'draft' && (
              <button onClick={() => setAddingQuestion(!addingQuestion)} className="inline-flex items-center gap-2 px-3 py-1.5 bg-violet-600 hover:bg-violet-500 text-white text-sm rounded-lg transition-colors"><Plus size={14} /> Add Question</button>
            )}
          </div>

          {addingQuestion && (
            <form onSubmit={handleAddQuestion} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-4 space-y-4 max-w-lg">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1.5">Question</label>
                <input value={newQuestion} onChange={e => setNewQuestion(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" required />
              </div>
              {newOptions.map((opt, i) => (
                <div key={i}>
                  <label className="block text-sm font-medium text-gray-400 mb-1.5">
                    Option {i + 1} {i === correctIndex && <span className="text-emerald-400">(correct)</span>}
                  </label>
                  <div className="flex gap-2">
                    <input value={opt} onChange={e => { const o = [...newOptions]; o[i] = e.target.value; setNewOptions(o) }} className="flex-1 px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" required />
                    <button type="button" onClick={() => setCorrectIndex(i)} className={`px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${i === correctIndex ? 'bg-emerald-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}>
                      {i === correctIndex ? 'Correct' : 'Set'}
                    </button>
                  </div>
                </div>
              ))}
              <div className="flex gap-3">
                <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">Add</button>
                <button type="button" onClick={() => setAddingQuestion(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg transition-colors">Cancel</button>
              </div>
            </form>
          )}

          <div className="space-y-3">
            {questions.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 text-center">
                <p className="text-gray-500">No questions yet. Add some before launching!</p>
              </div>
            ) : questions.map((q, idx) => (
              <div key={q.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4 group">
                <div className="flex items-start justify-between mb-2">
                  <p className="text-white font-medium"><span className="text-gray-500 mr-2">Q{idx + 1}.</span>{q.question}</p>
                  {game.status === 'draft' && (
                    <button onClick={() => handleDeleteQuestion(q.id)} className="p-1.5 text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"><Trash2 size={14} /></button>
                  )}
                </div>
                <div className="grid grid-cols-2 gap-2">
                  {q.options.map((opt, oi) => (
                    <div key={oi} className={`text-sm px-3 py-1.5 rounded-lg ${oi === q.correct_index ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-gray-800/50 text-gray-400'}`}>
                      {opt}
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {tab === 'content' && (game.game_type === 'photo_challenge' || game.game_type === 'two_truths') && (
        <div>
          <h3 className="text-white font-medium mb-4">Submissions ({submissions.length})</h3>
          <div className="space-y-3">
            {submissions.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 text-center">
                <p className="text-gray-500">No submissions yet.</p>
              </div>
            ) : submissions.map(s => (
              <div key={s.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs text-gray-500">User: {s.user_id.slice(0, 8)}...</span>
                  {s.vote_count > 0 && <span className="text-xs text-violet-400">{s.vote_count} votes</span>}
                </div>
                <p className="text-white text-sm">{s.content}</p>
                <span className="text-xs text-gray-600">{new Date(s.created_at).toLocaleString()}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {tab === 'results' && (
        <div>
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-4">
            <p className="text-gray-400 text-sm">Participants</p>
            <p className="text-2xl font-bold text-white">{participantCount}</p>
          </div>
          <h3 className="text-white font-medium mb-3">Leaderboard</h3>
          <div className="space-y-2">
            {leaderboard.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 text-center">
                <p className="text-gray-500">No results available.</p>
              </div>
            ) : leaderboard.map((entry, i) => (
              <div key={entry.user_id} className="bg-gray-900 border border-gray-800 rounded-xl p-4 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className={`w-7 h-7 flex items-center justify-center rounded-full text-sm font-bold ${i === 0 ? 'bg-amber-500/20 text-amber-400' : i === 1 ? 'bg-gray-400/20 text-gray-300' : i === 2 ? 'bg-orange-500/20 text-orange-400' : 'bg-gray-800 text-gray-500'}`}>
                    {i + 1}
                  </span>
                  <span className="text-white text-sm">{entry.user_id.slice(0, 8)}...</span>
                </div>
                <span className="text-violet-400 font-medium">{entry.score} pts</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
