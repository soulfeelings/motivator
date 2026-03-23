import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Send, X, Clock, Check, Ban } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { TableSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface Invite {
  id: string
  email: string
  role: string
  status: string
  token: string
  expires_at: string
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string; role: string }>
}

const statusStyle: Record<string, { icon: typeof Clock; class: string }> = {
  pending: { icon: Clock, class: 'text-amber-400' },
  accepted: { icon: Check, class: 'text-emerald-400' },
  expired: { icon: Ban, class: 'text-gray-500' },
  revoked: { icon: X, class: 'text-red-400' },
}

export default function Invites() {
  const navigate = useNavigate()
  const [invites, setInvites] = useState<Invite[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [email, setEmail] = useState('')
  const [role, setRole] = useState('employee')
  const [sending, setSending] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
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
        const res = await api.get<any>(`/companies/${cid}/invites`)
        setInvites(res ?? [])
      }
    } catch {
      // no company
    } finally {
      setLoading(false)
    }
  }

  async function handleSend(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    setSending(true)
    try {
      const invite = await api.post<Invite>(`/companies/${companyId}/invites`, { email, role })
      setInvites((prev) => [invite, ...prev])
      setEmail('')
      setToast({ message: `Invite sent to ${email}`, type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to send invite', type: 'error' })
    } finally {
      setSending(false)
    }
  }

  async function revoke(inviteId: string) {
    if (!companyId || !confirm('Revoke this invite?')) return
    try {
      await api.delete(`/companies/${companyId}/invites/${inviteId}`)
      setInvites((prev) => prev.map((i) => (i.id === inviteId ? { ...i, status: 'revoked' } : i)))
      setToast({ message: 'Invite revoked', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to revoke invite', type: 'error' })
    }
  }

  if (loading) return <TableSkeleton rows={5} cols={5} />
  if (!companyId) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-white mb-6">Invites</h2>
        <EmptyState
          icon={Send}
          title="No company yet"
          description="Invite team members by email so they can join your company and start earning XP."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-white mb-6">Invites</h2>

      <form onSubmit={handleSend} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6">
        <div className="flex gap-3 items-end">
          <div className="flex-1">
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
          <button
            type="submit"
            disabled={sending}
            className="inline-flex items-center gap-2 px-4 py-2.5 bg-violet-600 hover:bg-violet-500 disabled:opacity-50 text-white font-medium rounded-lg transition-colors"
          >
            <Send size={16} />
            {sending ? 'Sending...' : 'Send'}
          </button>
        </div>
        {error && <p className="text-sm text-red-400 mt-3">{error}</p>}
      </form>

      {invites.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400">No invites sent yet.</p>
        </div>
      ) : (
        <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-800">
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th className="text-left px-5 py-3 text-xs font-medium text-gray-500 uppercase tracking-wider">Sent</th>
                <th className="px-5 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {invites.map((inv) => {
                const s = statusStyle[inv.status] ?? statusStyle.pending
                const Icon = s.icon
                return (
                  <tr key={inv.id} className="border-b border-gray-800/50 hover:bg-gray-800/30 transition-colors">
                    <td className="px-5 py-4 text-sm text-white">{inv.email}</td>
                    <td className="px-5 py-4 text-sm text-gray-400 capitalize">{inv.role}</td>
                    <td className="px-5 py-4">
                      <span className={`inline-flex items-center gap-1.5 text-sm ${s.class}`}>
                        <Icon size={14} />
                        {inv.status}
                      </span>
                    </td>
                    <td className="px-5 py-4 text-sm text-gray-500">
                      {new Date(inv.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-5 py-4 text-right">
                      {inv.status === 'pending' && (
                        <button
                          onClick={() => revoke(inv.id)}
                          className="text-gray-500 hover:text-red-400 transition-colors"
                          title="Revoke invite"
                        >
                          <X size={16} />
                        </button>
                      )}
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
