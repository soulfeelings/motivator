import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { BarChart3, Users, Zap, Coins, Award, Target, Swords, Gift } from 'lucide-react'
import { StatSkeleton } from '../components/LoadingSkeleton'

interface Overview {
  total_members: number
  active_members: number
  total_xp_awarded: number
  total_coins_awarded: number
  total_coins_spent: number
  total_badges_awarded: number
  total_achievements_completed: number
  total_challenges: number
  total_redemptions: number
}

interface TopPerformer {
  membership_id: string
  display_name?: string
  xp: number
  level: number
  badges: number
  achievements: number
}

interface AchievementStat {
  name: string
  metric: string
  completions: number
}

interface ChallengeStat {
  total_challenges: number
  completed: number
  active: number
  pending: number
  avg_xp_reward: number
}

interface RewardStat {
  name: string
  cost_coins: number
  total_redeemed: number
  total_coins_spent: number
}

interface XPDist {
  level: number
  count: number
}

interface Dashboard {
  overview: Overview
  top_performers: TopPerformer[]
  achievement_stats: AchievementStat[]
  challenge_stats: ChallengeStat
  reward_stats: RewardStat[]
  xp_distribution: XPDist[]
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function Analytics() {
  const [data, setData] = useState<Dashboard | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        const d = await api.get<Dashboard>(`/companies/${cid}/analytics`)
        setData(d)
      }
    } catch {} finally { setLoading(false) }
  }

  if (loading) return <StatSkeleton />
  if (!data) return <p className="text-gray-500">No data available.</p>

  const o = data.overview
  const maxXP = Math.max(...(data.top_performers?.map(p => p.xp) ?? [1]), 1)
  const maxCompletions = Math.max(...(data.achievement_stats?.map(a => a.completions) ?? [1]), 1)
  const maxRedeemed = Math.max(...(data.reward_stats?.map(r => r.total_redeemed) ?? [1]), 1)
  const maxLevelCount = Math.max(...(data.xp_distribution?.map(d => d.count) ?? [1]), 1)

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <BarChart3 size={24} className="text-violet-400" />
        <h2 className="text-2xl font-bold text-white">Analytics</h2>
      </div>

      {/* Overview Cards */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3 mb-8">
        {[
          { label: 'Members', value: o.total_members, sub: `${o.active_members} active`, icon: Users, color: 'text-blue-400' },
          { label: 'Total XP', value: o.total_xp_awarded.toLocaleString(), icon: Zap, color: 'text-emerald-400' },
          { label: 'Coins In', value: o.total_coins_awarded.toLocaleString(), sub: `${o.total_coins_spent.toLocaleString()} spent`, icon: Coins, color: 'text-amber-400' },
          { label: 'Badges', value: o.total_badges_awarded, icon: Award, color: 'text-violet-400' },
          { label: 'Achievements', value: o.total_achievements_completed, icon: Target, color: 'text-emerald-400' },
        ].map(({ label, value, sub, icon: Icon, color }) => (
          <div key={label} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs text-gray-500">{label}</span>
              <Icon size={16} className={color} />
            </div>
            <p className="text-xl font-bold text-white">{value}</p>
            {sub && <p className="text-xs text-gray-600 mt-0.5">{sub}</p>}
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Performers */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <h3 className="text-white font-medium mb-4">Top Performers</h3>
          <div className="space-y-3">
            {(data.top_performers ?? []).slice(0, 8).map((p, i) => (
              <div key={p.membership_id} className="flex items-center gap-3">
                <span className={`w-6 text-center text-sm font-bold ${i < 3 ? 'text-amber-400' : 'text-gray-600'}`}>{i + 1}</span>
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-sm text-white">{p.display_name ?? p.membership_id.slice(0, 8)}</span>
                    <span className="text-xs text-emerald-400">{p.xp} XP</span>
                  </div>
                  <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden">
                    <div className="h-full bg-emerald-500/60 rounded-full" style={{ width: `${(p.xp / maxXP) * 100}%` }} />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Achievement Completion Rates */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <h3 className="text-white font-medium mb-4">Achievement Completions</h3>
          <div className="space-y-3">
            {(data.achievement_stats ?? []).slice(0, 8).map(a => (
              <div key={a.name}>
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm text-white">{a.name}</span>
                  <span className="text-xs text-violet-400">{a.completions}</span>
                </div>
                <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden">
                  <div className="h-full bg-violet-500/60 rounded-full" style={{ width: `${(a.completions / maxCompletions) * 100}%` }} />
                </div>
              </div>
            ))}
            {(!data.achievement_stats || data.achievement_stats.length === 0) && <p className="text-gray-500 text-sm">No achievements yet</p>}
          </div>
        </div>

        {/* Challenge Stats */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <div className="flex items-center gap-2 mb-4">
            <Swords size={18} className="text-violet-400" />
            <h3 className="text-white font-medium">Challenges</h3>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {[
              { label: 'Total', value: data.challenge_stats.total_challenges, color: 'text-white' },
              { label: 'Completed', value: data.challenge_stats.completed, color: 'text-emerald-400' },
              { label: 'Active', value: data.challenge_stats.active, color: 'text-violet-400' },
              { label: 'Avg XP', value: data.challenge_stats.avg_xp_reward, color: 'text-amber-400' },
            ].map(s => (
              <div key={s.label}>
                <p className="text-xs text-gray-500">{s.label}</p>
                <p className={`text-lg font-bold ${s.color}`}>{s.value}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Reward Stats */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <div className="flex items-center gap-2 mb-4">
            <Gift size={18} className="text-amber-400" />
            <h3 className="text-white font-medium">Reward Redemptions</h3>
          </div>
          <div className="space-y-3">
            {(data.reward_stats ?? []).map(r => (
              <div key={r.name}>
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm text-white">{r.name}</span>
                  <span className="text-xs text-amber-400">{r.total_redeemed} redeemed</span>
                </div>
                <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden">
                  <div className="h-full bg-amber-500/60 rounded-full" style={{ width: `${(r.total_redeemed / maxRedeemed) * 100}%` }} />
                </div>
              </div>
            ))}
            {(!data.reward_stats || data.reward_stats.length === 0) && <p className="text-gray-500 text-sm">No redemptions yet</p>}
          </div>
        </div>

        {/* XP Level Distribution */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5 lg:col-span-2">
          <h3 className="text-white font-medium mb-4">Level Distribution</h3>
          <div className="flex items-end gap-2 h-32">
            {(data.xp_distribution ?? []).map(d => (
              <div key={d.level} className="flex-1 flex flex-col items-center gap-1">
                <span className="text-xs text-gray-500">{d.count}</span>
                <div className="w-full bg-violet-500/40 rounded-t" style={{ height: `${(d.count / maxLevelCount) * 100}%`, minHeight: '4px' }} />
                <span className="text-xs text-gray-600">L{d.level}</span>
              </div>
            ))}
            {(!data.xp_distribution || data.xp_distribution.length === 0) && (
              <p className="text-gray-500 text-sm w-full text-center">No level data yet</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
