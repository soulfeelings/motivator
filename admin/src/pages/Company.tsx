import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Building2, Calendar, Link, Users, Plus, Pencil } from 'lucide-react'
import CompanyForm from '../components/CompanyForm'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface Company {
  id: string
  name: string
  slug: string
  logo_url?: string
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string; role: string; id: string }>
}

export default function Company() {
  const [company, setCompany] = useState<Company | null>(null)
  const [creating, setCreating] = useState(false)
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [memberCount, setMemberCount] = useState(0)
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null)

  useEffect(() => {
    loadCompany()
  }, [])

  async function loadCompany() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const companyId = me.memberships[0].company_id
        const [c, members] = await Promise.all([
          api.get<Company>(`/companies/${companyId}`),
          api.get<any>(`/companies/${companyId}/members`).catch(() => []),
        ])
        setCompany(c)
        setMemberCount(Array.isArray(members) ? members.length : 0)
      }
    } catch {
      // No company yet
    } finally {
      setLoading(false)
    }
  }

  async function handleUpdate(e: React.FormEvent) {
    e.preventDefault()
    if (!company) return
    setError('')
    try {
      const c = await api.patch<Company>(`/companies/${company.id}`, {
        name: name || company.name,
        slug: slug || company.slug,
      })
      setCompany(c)
      setEditing(false)
      setToast({ message: 'Company updated', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to update', type: 'error' })
    }
  }

  if (loading) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-white mb-6">Company</h2>
        <CardSkeleton count={1} />
      </div>
    )
  }

  if (!company && !creating) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-white mb-6">Company</h2>
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center">
          <p className="text-gray-400 mb-4">No company set up yet.</p>
          <button
            onClick={() => setCreating(true)}
            className="inline-flex items-center gap-2 px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors"
          >
            <Plus size={18} />
            Create Company
          </button>
        </div>
      </div>
    )
  }

  if (creating) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-white mb-6">Create Company</h2>
        <CompanyForm
          onSuccess={(c) => { setCompany(c); setCreating(false); setToast({ message: 'Company created!', type: 'success' }) }}
          onCancel={() => setCreating(false)}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white">Company</h2>
        {!editing && (
          <button
            onClick={() => { setEditing(true); setName(company!.name); setSlug(company!.slug) }}
            className="inline-flex items-center gap-2 px-3 py-2 text-sm text-gray-400 hover:text-white bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors"
          >
            <Pencil size={14} />
            Edit
          </button>
        )}
      </div>

      {editing ? (
        <form onSubmit={handleUpdate} className="bg-gray-900 border border-gray-800 rounded-xl p-6 space-y-4 max-w-lg">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Company Name</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Slug</label>
            <input
              value={slug}
              onChange={(e) => setSlug(e.target.value)}
              className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
            />
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">
              Save Changes
            </button>
            <button type="button" onClick={() => setEditing(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 font-medium rounded-lg transition-colors">
              Cancel
            </button>
          </div>
        </form>
      ) : (
        <div className="space-y-4">
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
            <div className="flex items-center gap-4 mb-6">
              <div className="w-14 h-14 rounded-xl bg-violet-600/20 flex items-center justify-center">
                <Building2 size={28} className="text-violet-400" />
              </div>
              <div>
                <h3 className="text-xl font-semibold text-white">{company!.name}</h3>
                <p className="text-sm text-gray-500">{company!.slug}</p>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div className="flex items-center gap-3 p-3 bg-gray-800/50 rounded-lg">
                <Link size={16} className="text-gray-500" />
                <div>
                  <p className="text-xs text-gray-500">Slug</p>
                  <p className="text-sm text-white">{company!.slug}</p>
                </div>
              </div>
              <div className="flex items-center gap-3 p-3 bg-gray-800/50 rounded-lg">
                <Users size={16} className="text-emerald-400" />
                <div>
                  <p className="text-xs text-gray-500">Members</p>
                  <p className="text-sm text-white">{memberCount}</p>
                </div>
              </div>
              <div className="flex items-center gap-3 p-3 bg-gray-800/50 rounded-lg">
                <Calendar size={16} className="text-gray-500" />
                <div>
                  <p className="text-xs text-gray-500">Created</p>
                  <p className="text-sm text-white">{new Date(company!.created_at).toLocaleDateString()}</p>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
            <h3 className="text-sm font-medium text-gray-400 mb-3">Company ID</h3>
            <code className="text-xs text-gray-500 bg-gray-800 px-3 py-1.5 rounded-lg">{company!.id}</code>
          </div>
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
