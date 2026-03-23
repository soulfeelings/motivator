import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Trophy, Zap, Coins } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { TableSkeleton } from '../components/LoadingSkeleton'

interface LeaderEntry {
  rank: number
  member_id: string
  user_id: string
  display_name?: string
  xp: number
  level: number
  coins: number
  role: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const medalColors = ['text-yellow-400', 'text-gray-400', 'text-amber-600']

export default function Leaderboard() {
  const navigate = useNavigate()
  const [entries, setEntries] = useState<LeaderEntry[]>([])
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
        const res = await api.get<LeaderEntry[]>(`/companies/${cid}/leaderboard?limit=50`)
        setEntries(res ?? [])
      }
    } catch {
      // no company
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <TableSkeleton rows={8} cols={4} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Trophy size={24} className="text-yellow-400" />
          <h2 className="text-2xl font-bold text-white">Leaderboard</h2>
        </div>
        <EmptyState
          icon={Trophy}
          title="No company yet"
          description="See how your team ranks by XP and track top performers in real time."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <Trophy size={24} className="text-yellow-400" />
        <h2 className="text-2xl font-bold text-white">Leaderboard</h2>
      </div>

      {entries.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No members yet.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {entries.map((e) => (
            <div
              key={e.member_id}
              className={`flex items-center gap-4 p-4 rounded-xl border transition-colors ${
                e.rank <= 3
                  ? 'bg-gray-900 border-gray-700'
                  : 'bg-gray-900/50 border-gray-800/50'
              }`}
            >
              <div className="w-10 text-center">
                {e.rank <= 3 ? (
                  <span className={`text-2xl font-bold ${medalColors[e.rank - 1]}`}>
                    {e.rank}
                  </span>
                ) : (
                  <span className="text-lg text-gray-600 font-medium">{e.rank}</span>
                )}
              </div>

              <div className="flex-1 min-w-0">
                <p className="text-white font-medium truncate">
                  {e.display_name || e.user_id.slice(0, 8)}
                </p>
                <p className="text-xs text-gray-500">Level {e.level}</p>
              </div>

              <div className="flex items-center gap-6">
                <div className="flex items-center gap-1.5 text-emerald-400">
                  <Zap size={16} />
                  <span className="font-bold">{e.xp}</span>
                  <span className="text-xs text-gray-500">XP</span>
                </div>
                <div className="flex items-center gap-1.5 text-amber-400">
                  <Coins size={16} />
                  <span className="font-bold">{e.coins}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
