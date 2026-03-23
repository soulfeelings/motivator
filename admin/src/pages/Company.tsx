import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Plus } from 'lucide-react'
import CompanyForm from '../components/CompanyForm'
import { CardSkeleton } from '../components/LoadingSkeleton'

interface Company {
  id: string
  name: string
  slug: string
  logo_url?: string
  created_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string; role: string }>
}

export default function Company() {
  const [company, setCompany] = useState<Company | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadCompany()
  }, [])

  async function loadCompany() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const c = await api.get<Company>(`/companies/${me.memberships[0].company_id}`)
        setCompany(c)
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
      const c = await api.patch<Company>(`/companies/${company.id}`, { name, slug })
      setCompany(c)
    } catch (err: any) {
      setError(err.message)
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
          onSuccess={(c) => { setCompany(c); setCreating(false) }}
          onCancel={() => setCreating(false)}
        />
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-white mb-6">Company</h2>
      <form onSubmit={handleUpdate} className="bg-gray-900 border border-gray-800 rounded-xl p-6 space-y-4 max-w-lg">
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1.5">Company Name</label>
          <input
            defaultValue={company!.name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1.5">Slug</label>
          <input
            defaultValue={company!.slug}
            onChange={(e) => setSlug(e.target.value)}
            className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1.5">Created</label>
          <p className="text-sm text-gray-500">{new Date(company!.created_at).toLocaleDateString()}</p>
        </div>
        {error && <p className="text-sm text-red-400">{error}</p>}
        <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg transition-colors">
          Save Changes
        </button>
      </form>
    </div>
  )
}
