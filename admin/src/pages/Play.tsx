import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../lib/api'
import { Gamepad2, CheckCircle, XCircle, Send, ThumbsUp } from 'lucide-react'

interface SocialGame {
  id: string
  name: string
  description?: string
  game_type: string
  status: string
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
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function Play() {
  const { gameId } = useParams<{ gameId: string }>()
  const [companyId, setCompanyId] = useState('')
  const [game, setGame] = useState<SocialGame | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Trivia state
  const [questions, setQuestions] = useState<TriviaQuestion[]>([])
  const [currentQ, setCurrentQ] = useState(0)
  const [selected, setSelected] = useState<number | null>(null)
  const [score, setScore] = useState(0)
  const [triviaFinished, setTriviaFinished] = useState(false)

  // Two truths state
  const [statements, setStatements] = useState(['', '', ''])
  const [lieIndex, setLieIndex] = useState(0)
  const [submitted, setSubmitted] = useState(false)
  const [otherSubmissions, setOtherSubmissions] = useState<Submission[]>([])
  const [votedIds, setVotedIds] = useState<Set<string>>(new Set())

  // Photo challenge state
  const [photoUrl, setPhotoUrl] = useState('')
  const [photoSubmitted, setPhotoSubmitted] = useState(false)
  const [gallery, setGallery] = useState<Submission[]>([])

  useEffect(() => { load() }, [gameId])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (!me.memberships?.length) { setError('No company found'); setLoading(false); return }
      const cid = me.memberships[0].company_id
      setCompanyId(cid)

      const g = await api.get<SocialGame>(`/companies/${cid}/social-games/${gameId}`)
      setGame(g)

      if (g.game_type === 'trivia') {
        const qs = await api.get<TriviaQuestion[]>(`/companies/${cid}/social-games/${gameId}/questions`)
        setQuestions(qs ?? [])
      }

      if (g.game_type === 'two_truths' && g.status === 'voting') {
        const subs = await api.get<Submission[]>(`/companies/${cid}/social-games/${gameId}/submissions`)
        setOtherSubmissions(subs ?? [])
      }

      if (g.game_type === 'photo_challenge' && g.status === 'voting') {
        const subs = await api.get<Submission[]>(`/companies/${cid}/social-games/${gameId}/submissions`)
        setGallery(subs ?? [])
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load game')
    } finally {
      setLoading(false)
    }
  }

  function handleTriviaAnswer(idx: number) {
    if (selected !== null) return
    setSelected(idx)
    if (idx === questions[currentQ].correct_index) {
      setScore(prev => prev + 1)
    }
    setTimeout(() => {
      if (currentQ + 1 < questions.length) {
        setCurrentQ(prev => prev + 1)
        setSelected(null)
      } else {
        setTriviaFinished(true)
      }
    }, 1500)
  }

  async function handleTwoTruthsSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      await api.post(`/companies/${companyId}/social-games/${gameId}/submit`, {
        statements,
        lie_index: lieIndex,
      })
      setSubmitted(true)
    } catch (err: any) {
      setError(err.message || 'Failed to submit')
    }
  }

  async function handlePhotoSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      await api.post(`/companies/${companyId}/social-games/${gameId}/submit`, {
        content: photoUrl,
      })
      setPhotoSubmitted(true)
    } catch (err: any) {
      setError(err.message || 'Failed to submit')
    }
  }

  async function handleVote(submissionId: string) {
    try {
      await api.post(`/companies/${companyId}/social-games/${gameId}/vote`, { submission_id: submissionId, vote_value: 1 })
      setVotedIds(prev => new Set(prev).add(submissionId))
    } catch (err: any) {
      setError(err.message || 'Failed to vote')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="animate-pulse text-center">
          <Gamepad2 size={32} className="text-violet-400 mx-auto mb-3" />
          <p className="text-gray-500 text-sm">Loading game...</p>
        </div>
      </div>
    )
  }

  if (error && !game) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <Gamepad2 size={32} className="text-gray-600 mx-auto mb-3" />
          <p className="text-red-400 text-sm">{error}</p>
        </div>
      </div>
    )
  }

  if (!game) return null

  return (
    <div className="min-h-screen bg-gray-950">
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="text-center mb-8">
          <h1 className="text-lg font-bold text-violet-400 tracking-tight mb-4">Motivator</h1>
          <h2 className="text-2xl font-bold text-white mb-2">{game.name}</h2>
          {game.description && <p className="text-gray-400">{game.description}</p>}
        </div>

        {error && <p className="text-sm text-red-400 text-center mb-4">{error}</p>}

        {/* Trivia */}
        {game.game_type === 'trivia' && !triviaFinished && questions.length > 0 && (
          <div>
            <div className="flex items-center justify-between mb-4">
              <span className="text-sm text-gray-500">Question {currentQ + 1} of {questions.length}</span>
              <span className="text-sm text-violet-400">{score} correct</span>
            </div>
            <div className="w-full bg-gray-800 rounded-full h-1.5 mb-6">
              <div className="bg-violet-600 h-1.5 rounded-full transition-all" style={{ width: `${((currentQ + 1) / questions.length) * 100}%` }} />
            </div>
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 mb-4">
              <p className="text-white text-lg font-medium">{questions[currentQ].question}</p>
            </div>
            <div className="space-y-3">
              {questions[currentQ].options.map((opt, i) => {
                let classes = 'w-full text-left px-5 py-4 rounded-xl text-sm font-medium transition-all '
                if (selected === null) {
                  classes += 'bg-gray-900 border border-gray-800 text-white hover:border-violet-500/50'
                } else if (i === questions[currentQ].correct_index) {
                  classes += 'bg-emerald-500/10 border border-emerald-500/30 text-emerald-400'
                } else if (i === selected) {
                  classes += 'bg-red-500/10 border border-red-500/30 text-red-400'
                } else {
                  classes += 'bg-gray-900 border border-gray-800 text-gray-600'
                }
                return (
                  <button key={i} onClick={() => handleTriviaAnswer(i)} disabled={selected !== null} className={classes}>
                    <div className="flex items-center justify-between">
                      <span>{opt}</span>
                      {selected !== null && i === questions[currentQ].correct_index && <CheckCircle size={18} />}
                      {selected !== null && i === selected && i !== questions[currentQ].correct_index && <XCircle size={18} />}
                    </div>
                  </button>
                )
              })}
            </div>
          </div>
        )}

        {game.game_type === 'trivia' && triviaFinished && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <CheckCircle size={48} className="text-emerald-400 mx-auto mb-4" />
            <h3 className="text-2xl font-bold text-white mb-2">Game Over!</h3>
            <p className="text-gray-400 mb-4">You got <span className="text-violet-400 font-bold">{score}</span> out of <span className="text-white font-bold">{questions.length}</span> correct</p>
            <div className="w-32 mx-auto bg-gray-800 rounded-full h-3">
              <div className="bg-violet-600 h-3 rounded-full" style={{ width: `${(score / questions.length) * 100}%` }} />
            </div>
          </div>
        )}

        {game.game_type === 'trivia' && questions.length === 0 && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <p className="text-gray-500">This trivia has no questions yet. Check back later!</p>
          </div>
        )}

        {/* Two Truths */}
        {game.game_type === 'two_truths' && game.status === 'active' && !submitted && (
          <form onSubmit={handleTwoTruthsSubmit} className="space-y-4">
            <p className="text-gray-400 text-sm text-center mb-4">Enter two truths and one lie about yourself. Others will try to guess the lie!</p>
            {statements.map((s, i) => (
              <div key={i}>
                <label className="block text-sm font-medium text-gray-400 mb-1.5">
                  Statement {i + 1} {i === lieIndex && <span className="text-pink-400">(the lie)</span>}
                </label>
                <div className="flex gap-2">
                  <input
                    value={s}
                    onChange={e => { const arr = [...statements]; arr[i] = e.target.value; setStatements(arr) }}
                    className="flex-1 px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                    placeholder={i === lieIndex ? 'This one is the lie...' : 'Enter a truth...'}
                    required
                  />
                  <button type="button" onClick={() => setLieIndex(i)} className={`px-3 py-2.5 rounded-lg text-xs font-medium transition-colors ${i === lieIndex ? 'bg-pink-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700'}`}>
                    {i === lieIndex ? 'Lie' : 'Set Lie'}
                  </button>
                </div>
              </div>
            ))}
            <button type="submit" className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">
              <Send size={16} /> Submit
            </button>
          </form>
        )}

        {game.game_type === 'two_truths' && game.status === 'active' && submitted && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <CheckCircle size={48} className="text-emerald-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">Submitted!</h3>
            <p className="text-gray-400">Voting will start soon. Check back when it does!</p>
          </div>
        )}

        {game.game_type === 'two_truths' && game.status === 'voting' && (
          <div className="space-y-4">
            <p className="text-gray-400 text-sm text-center mb-4">Guess which statement is the lie!</p>
            {otherSubmissions.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 text-center">
                <p className="text-gray-500">No submissions to vote on.</p>
              </div>
            ) : otherSubmissions.map(sub => {
              let parsed: string[] = []
              try { parsed = JSON.parse(sub.content) } catch { parsed = [sub.content] }
              return (
                <div key={sub.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                  <p className="text-xs text-gray-500 mb-3">User {sub.user_id.slice(0, 8)}...</p>
                  <div className="space-y-2">
                    {parsed.map((stmt, si) => (
                      <button
                        key={si}
                        onClick={() => handleVote(sub.id)}
                        disabled={votedIds.has(sub.id)}
                        className={`w-full text-left px-4 py-3 rounded-lg text-sm transition-colors ${votedIds.has(sub.id) ? 'bg-gray-800/50 text-gray-500' : 'bg-gray-800 text-white hover:bg-violet-600/20 hover:text-violet-400'}`}
                      >
                        {stmt}
                      </button>
                    ))}
                  </div>
                </div>
              )
            })}
          </div>
        )}

        {/* Photo Challenge */}
        {game.game_type === 'photo_challenge' && game.status === 'active' && !photoSubmitted && (
          <form onSubmit={handlePhotoSubmit} className="space-y-4">
            <p className="text-gray-400 text-sm text-center mb-4">Share a photo URL for the challenge!</p>
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Photo URL</label>
              <input
                value={photoUrl}
                onChange={e => setPhotoUrl(e.target.value)}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                placeholder="https://example.com/my-photo.jpg"
                required
              />
            </div>
            <button type="submit" className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">
              <Send size={16} /> Submit Photo
            </button>
          </form>
        )}

        {game.game_type === 'photo_challenge' && game.status === 'active' && photoSubmitted && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <CheckCircle size={48} className="text-emerald-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">Photo Submitted!</h3>
            <p className="text-gray-400">Voting will start soon.</p>
          </div>
        )}

        {game.game_type === 'photo_challenge' && game.status === 'voting' && (
          <div className="space-y-4">
            <p className="text-gray-400 text-sm text-center mb-4">Vote for your favorite photos!</p>
            {gallery.length === 0 ? (
              <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 text-center">
                <p className="text-gray-500">No photos submitted.</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {gallery.map(sub => (
                  <div key={sub.id} className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
                    <div className="aspect-video bg-gray-800 flex items-center justify-center">
                      <img src={sub.content} alt="Submission" className="w-full h-full object-cover" onError={e => { (e.target as HTMLImageElement).style.display = 'none' }} />
                    </div>
                    <div className="p-3 flex items-center justify-between">
                      <span className="text-xs text-gray-500">{sub.vote_count} votes</span>
                      <button
                        onClick={() => handleVote(sub.id)}
                        disabled={votedIds.has(sub.id)}
                        className={`inline-flex items-center gap-1 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${votedIds.has(sub.id) ? 'bg-violet-600/20 text-violet-400' : 'bg-gray-800 text-gray-400 hover:bg-violet-600 hover:text-white'}`}
                      >
                        <ThumbsUp size={14} /> {votedIds.has(sub.id) ? 'Voted' : 'Vote'}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {game.status === 'completed' && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <CheckCircle size={48} className="text-blue-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">Game Complete!</h3>
            <p className="text-gray-400">Thanks for playing. Results are in!</p>
          </div>
        )}

        {game.status === 'draft' && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
            <Gamepad2 size={48} className="text-gray-600 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">Coming Soon</h3>
            <p className="text-gray-400">This game hasn't started yet. Check back later!</p>
          </div>
        )}
      </div>
    </div>
  )
}
