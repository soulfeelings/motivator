import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Shield, UserX, Zap, Coins } from 'lucide-react'

interface Member {
  id: string
  user_id: string
  role: string
  display_name?: string
  job_title?: string
  xp: number
  level: number
  coins: number
  is_active: boolean
  joined_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string; role: string }>
}

const roleBadge: Record<string, string> = {
  owner: 'bg-violet-500/20 text-violet-400',
  admin: 'bg-blue-500/20 text-blue-400',
  manager: 'bg-emerald-500/20 text-emerald-400',
  employee: 'bg-gray-500/20 text-gray-400',
}

export default function Members() {
  const [members, setMembers] = useState<Member[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadMembers()
  }, [])

  async function loadMembers() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<any>(`/companies/${cid}/members`)
        setMembers(res ?? [])
      }
    } catch {
      // no company
    } finally {
      setLoading(false)
    }
  }

  async function deactivate(memberId: string) {
    if (!companyId || !confirm('Remove this member?')) return
    await api.delete(`/companies/${companyId}/members/${memberId}`)
    setMembers((prev) => prev.filter((m) => m.id !== memberId))
  }

  if (loading) return <p className="text-gray-500">Loading...</p>

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white">Members</h2>
        <span className="text-sm text-gray-500">{members.length} total</span>
      </div>

      {members.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No members yet. Send invites to add people.</p>
        </div>
      ) : (
        <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-800">
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Level</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">XP</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Coins</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Joined</th>
                <th className="px-5 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {members.map((m) => (
                <tr key={m.id} className="border-b border-gray-800/50 hover:bg-gray-800/30 transition-colors">
                  <td className="px-5 py-4">
                    <p className="text-sm text-white font-medium">{m.display_name || m.user_id.slice(0, 8)}</p>
                    {m.job_title && <p className="text-xs text-gray-500">{m.job_title}</p>}
                  </td>
                  <td className="px-5 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium ${roleBadge[m.role] ?? roleBadge.employee}`}>
                      <Shield size={12} />
                      {m.role}
                    </span>
                  </td>
                  <td className="px-5 py-4">
                    <span className="inline-flex items-center gap-1 text-sm font-bold text-white">
                      {m.level}
                    </span>
                  </td>
                  <td className="px-5 py-4">
                    <span className="inline-flex items-center gap-1 text-sm text-emerald-400">
                      <Zap size={14} />
                      {m.xp}
                    </span>
                  </td>
                  <td className="px-5 py-4">
                    <span className="inline-flex items-center gap-1 text-sm text-amber-400">
                      <Coins size={14} />
                      {m.coins}
                    </span>
                  </td>
                  <td className="px-5 py-4 text-sm text-gray-500">
                    {new Date(m.joined_at).toLocaleDateString()}
                  </td>
                  <td className="px-5 py-4 text-right">
                    {m.role !== 'owner' && (
                      <button
                        onClick={() => deactivate(m.id)}
                        className="text-gray-500 hover:text-red-400 transition-colors"
                        title="Remove member"
                      >
                        <UserX size={16} />
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
