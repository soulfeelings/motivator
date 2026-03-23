import { useState } from 'react'
import { api } from '../lib/api'

interface Company {
  id: string
  name: string
  slug: string
  logo_url?: string
  created_at: string
}

interface CompanyFormProps {
  onSuccess: (company: Company) => void
  onCancel?: () => void
}

export default function CompanyForm({ onSuccess, onCancel }: CompanyFormProps) {
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setSubmitting(true)
    try {
      const c = await api.post<Company>('/companies', { name, slug })
      onSuccess(c)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 border border-gray-800 rounded-xl p-6 space-y-4 max-w-lg">
      <div>
        <label className="block text-sm font-medium text-gray-400 mb-1.5">Company Name</label>
        <input
          value={name}
          onChange={(e) => { setName(e.target.value); setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')) }}
          className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          placeholder="Acme Corp"
          required
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-400 mb-1.5">Slug</label>
        <input
          value={slug}
          onChange={(e) => setSlug(e.target.value)}
          className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
          placeholder="acme-corp"
          required
        />
      </div>
      {error && <p className="text-sm text-red-400">{error}</p>}
      <div className="flex gap-3">
        <button
          type="submit"
          disabled={submitting}
          className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
        >
          {submitting ? 'Creating...' : 'Create'}
        </button>
        {onCancel && (
          <button type="button" onClick={onCancel} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 font-medium rounded-lg transition-colors">
            Cancel
          </button>
        )}
      </div>
    </form>
  )
}
