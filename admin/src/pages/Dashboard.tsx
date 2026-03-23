import { useEffect, useState } from 'react'
import { Users, Mail, Swords, Award, Trophy, Target } from 'lucide-react'
import { api } from '../lib/api'
import { StatSkeleton } from '../components/LoadingSkeleton'

interface MeResponse {
  user_id: string
  email: string
  memberships: Array<{
    id: string
    company_id: string
    role: string
  }>
}

interface DashboardData {
  memberCount: number
  pendingInvites: number
  activeChallenges: number
  badgeCount: number
  achievementCount: number
  topPlayer: { user_id: string; xp: number; level: number } | null
}

export default function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboard()
  }, [])

  async function loadDashboard() {
    try {
      const meData = await api.get<MeResponse>('/me')

      if (meData.memberships?.length > 0) {
        const companyId = meData.memberships[0].company_id

        const [members, invites, challenges, badges, achievements, leaderboard] = await Promise.all([
          api.get<any>(`/companies/${companyId}/members`).catch(() => []),
          api.get<any>(`/companies/${companyId}/invites`).catch(() => ({ meta: { total: 0 } })),
          api.get<any>(`/companies/${companyId}/challenges`).catch(() => []),
          api.get<any>(`/companies/${companyId}/badges`).catch(() => []),
          api.get<any>(`/companies/${companyId}/achievements`).catch(() => []),
          api.get<any>(`/companies/${companyId}/leaderboard`).catch(() => []),
        ])

        const memberList = Array.isArray(members) ? members : []
        const challengeList = Array.isArray(challenges) ? challenges : []
        const badgeList = Array.isArray(badges) ? badges : []
        const achievementList = Array.isArray(achievements) ? achievements : []
        const leaderboardList = Array.isArray(leaderboard) ? leaderboard : []
        const pendingCount = invites?.meta?.total ?? (Array.isArray(invites) ? invites.filter((i: any) => i.status === 'pending').length : 0)

        setData({
          memberCount: memberList.length,
          pendingInvites: pendingCount,
          activeChallenges: challengeList.filter((c: any) => c.status === 'pending' || c.status === 'active').length,
          badgeCount: badgeList.length,
          achievementCount: achievementList.length,
          topPlayer: leaderboardList.length > 0 ? leaderboardList[0] : null,
        })
      }
    } catch {
      // Silent fail
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-white mb-6">Dashboard</h2>
        <StatSkeleton />
      </div>
    )
  }

  const stats = [
    { label: 'Members', value: data?.memberCount ?? 0, icon: Users, color: 'text-emerald-400' },
    { label: 'Pending Invites', value: data?.pendingInvites ?? 0, icon: Mail, color: 'text-amber-400' },
    { label: 'Active Challenges', value: data?.activeChallenges ?? 0, icon: Swords, color: 'text-red-400' },
    { label: 'Badges', value: data?.badgeCount ?? 0, icon: Award, color: 'text-violet-400' },
    { label: 'Achievements', value: data?.achievementCount ?? 0, icon: Target, color: 'text-blue-400' },
  ]

  return (
    <div>
      <h2 className="text-2xl font-bold text-white mb-6">Dashboard</h2>

      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4 mb-8">
        {stats.map(({ label, value, icon: Icon, color }) => (
          <div key={label} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm text-gray-500">{label}</span>
              <Icon size={18} className={color} />
            </div>
            <p className="text-2xl font-bold text-white">{value}</p>
          </div>
        ))}
      </div>

      {data?.topPlayer && (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 mb-6">
          <h3 className="text-sm font-medium text-gray-500 mb-4 flex items-center gap-2">
            <Trophy size={16} className="text-amber-400" />
            Top Player
          </h3>
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-full bg-amber-500/20 flex items-center justify-center">
              <Trophy size={20} className="text-amber-400" />
            </div>
            <div>
              <p className="text-white font-medium">{data.topPlayer.user_id.slice(0, 8)}</p>
              <p className="text-sm text-gray-500">Level {data.topPlayer.level} · {data.topPlayer.xp} XP</p>
            </div>
          </div>
        </div>
      )}

      {(data?.memberCount ?? 0) === 0 && (
        <div className="bg-gray-900/50 border border-dashed border-gray-700 rounded-xl p-6 text-center">
          <p className="text-gray-400 mb-1">Your workspace is ready!</p>
          <p className="text-sm text-gray-600">Invite team members and create badges to get started.</p>
        </div>
      )}
    </div>
  )
}
