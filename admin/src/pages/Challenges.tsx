import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Swords, Clock, Check, X } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'

interface Challenge {
  id: string
  challenger_id: string
  opponent_id: string
  metric: string
  target: number
  wager: number
  status: string
  challenger_score: number
  opponent_score: number
  winner_id?: string
  xp_reward: number
  deadline: string
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string; id: string }>
}

const statusStyle: Record<string, { icon: typeof Clock; class: string }> = {
  pending: { icon: Clock, class: 'text-amber-400' },
  active: { icon: Swords, class: 'text-violet-400' },
  completed: { icon: Check, class: 'text-emerald-400' },
  declined: { icon: X, class: 'text-gray-500' },
}

export default function Challenges() {
  const navigate = useNavigate()
  const [challenges, setChallenges] = useState<Challenge[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
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
        const res = await api.get<Challenge[]>(`/companies/${cid}/challenges`)
        setChallenges(res ?? [])
      }
    } catch {
      // no company
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <CardSkeleton count={4} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Swords size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">1v1 Challenges</h2>
        </div>
        <EmptyState
          icon={Swords}
          title="No company yet"
          description="View 1v1 challenges where employees compete head-to-head on specific metrics."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <Swords size={24} className="text-violet-400" />
        <h2 className="text-2xl font-bold text-white">1v1 Challenges</h2>
        <span className="text-sm text-gray-500 ml-auto">{challenges.length} total</span>
      </div>

      {challenges.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No challenges yet. Employees can create challenges from the mobile app.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {challenges.map((ch) => {
            const s = statusStyle[ch.status] ?? statusStyle.pending
            const Icon = s.icon
            return (
              <div key={ch.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    <span className={`inline-flex items-center gap-1.5 text-sm font-medium ${s.class}`}>
                      <Icon size={16} />
                      {ch.status}
                    </span>
                    <span className="text-xs text-gray-600">|</span>
                    <span className="text-sm text-violet-400 font-mono">{ch.metric}</span>
                    <span className="text-xs text-gray-500">target: {ch.target}</span>
                  </div>
                  <div className="text-xs text-gray-500">
                    {new Date(ch.created_at).toLocaleDateString()}
                  </div>
                </div>
                <div className="flex items-center gap-6">
                  <div className="flex-1 bg-gray-800/50 rounded-lg p-3 text-center">
                    <p className="text-xs text-gray-500 mb-1">Challenger</p>
                    <p className="text-white font-medium">{ch.challenger_id.slice(0, 8)}</p>
                    <p className="text-lg font-bold text-emerald-400">{ch.challenger_score}</p>
                  </div>
                  <span className="text-gray-600 font-bold">VS</span>
                  <div className="flex-1 bg-gray-800/50 rounded-lg p-3 text-center">
                    <p className="text-xs text-gray-500 mb-1">Opponent</p>
                    <p className="text-white font-medium">{ch.opponent_id.slice(0, 8)}</p>
                    <p className="text-lg font-bold text-emerald-400">{ch.opponent_score}</p>
                  </div>
                </div>
                {ch.wager > 0 && (
                  <p className="text-xs text-amber-400 mt-2">Wager: {ch.wager} coins | XP: +{ch.xp_reward}</p>
                )}
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
