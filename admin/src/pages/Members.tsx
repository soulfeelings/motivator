import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Shield, UserX, Zap, Coins, UserPlus } from 'lucide-react'
import { TableSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

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
  const [adding, setAdding] = useState(false)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [role, setRole] = useState('employee')
  const [error, setError] = useState('')
  const [addLoading, setAddLoading] = useState(false)
  const [toast, setToast] = useState<{message: string, type: 'success'|'error'} | null>(null)

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

  async function handleAdd(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    setAddLoading(true)
    try {
      await api.post(`/companies/${companyId}/members/add`, {
        email,
        password,
        role,
        display_name: displayName || undefined,
      })
      setAdding(false)
      setEmail('')
      setPassword('')
      setDisplayName('')
      setRole('employee')
      loadMembers()
      setToast({ message: 'Member added successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to add member', type: 'error' })
    } finally {
      setAddLoading(false)
    }
  }

  async function deactivate(memberId: string) {
    if (!companyId || !confirm('Remove this member?')) return
    try {
      await api.delete(`/companies/${companyId}/members/${memberId}`)
      setMembers((prev) => prev.filter((m) => m.id !== memberId))
      setToast({ message: 'Member removed', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to remove member', type: 'error' })
    }
  }

  if (loading) return <TableSkeleton rows={5} cols={6} />

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white">Members</h2>
        <div className="flex items-center gap-3">
          <span className="text-sm text-gray-500">{members.length} total</span>
          <button
            onClick={() => setAdding(!adding)}
            className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors"
          >
            <UserPlus size={16} />
            Add Member
          </button>
        </div>
      </div>

      {adding && (
        <form onSubmit={handleAdd} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="employee@company.com"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
              placeholder="Min 6 characters"
              minLength={6}
              required
            />
          </div>
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Display Name</label>
              <input
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
                placeholder="John Doe"
              />
            </div>
            <div className="w-40">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Role</label>
              <select
                value={role}
                onChange={(e) => setRole(e.target.value)}
                className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
              >
                <option value="employee">Employee</option>
                <option value="manager">Manager</option>
                <option value="admin">Admin</option>
              </select>
            </div>
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button
              type="submit"
              disabled={addLoading}
              className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 disabled:opacity-50 text-white font-medium rounded-lg transition-colors"
            >
              {addLoading ? 'Creating...' : 'Create & Add'}
            </button>
            <button type="button" onClick={() => { setAdding(false); setError('') }} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 font-medium rounded-lg transition-colors">
              Cancel
            </button>
          </div>
          <p className="text-xs text-gray-600">Creates a Supabase Auth account and adds them to the company in one step.</p>
        </form>
      )}

      {members.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No members yet. Add your first team member above.</p>
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
                    <span className="inline-flex items-center gap-1 text-sm font-bold text-white">{m.level}</span>
                  </td>
                  <td className="px-5 py-4">
                    <span className="inline-flex items-center gap-1 text-sm text-emerald-400">
                      <Zap size={14} /> {m.xp}
                    </span>
                  </td>
                  <td className="px-5 py-4">
                    <span className="inline-flex items-center gap-1 text-sm text-amber-400">
                      <Coins size={14} /> {m.coins}
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

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
