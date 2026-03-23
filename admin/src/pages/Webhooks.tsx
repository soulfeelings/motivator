import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { Bell, Plus, Trash2 } from 'lucide-react'
import EmptyState from '../components/EmptyState'
import { CardSkeleton } from '../components/LoadingSkeleton'
import Toast from '../components/Toast'

interface Webhook {
  id: string
  name: string
  url: string
  platform: string
  events: string[]
  is_active: boolean
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const eventOptions = [
  'achievement_completed', 'badge_awarded', 'challenge_created', 'challenge_completed',
  'tournament_started', 'tournament_completed', 'reward_redeemed',
]

export default function Webhooks() {
  const navigate = useNavigate()
  const [webhooks, setWebhooks] = useState<Webhook[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [name, setName] = useState('')
  const [url, setUrl] = useState('')
  const [platform, setPlatform] = useState('slack')
  const [events, setEvents] = useState<string[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [toast, setToast] = useState<{message: string, type: 'success'|'error'} | null>(null)

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<Webhook[]>(`/companies/${cid}/webhooks`)
        setWebhooks(res ?? [])
      }
    } catch {} finally { setLoading(false) }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const w = await api.post<Webhook>(`/companies/${companyId}/webhooks`, { name, url, platform, events })
      setWebhooks(prev => [...prev, w])
      setCreating(false)
      setName(''); setUrl(''); setEvents([])
      setToast({ message: 'Webhook created successfully', type: 'success' })
    } catch (err: any) {
      setError(err.message)
      setToast({ message: err.message || 'Failed to create webhook', type: 'error' })
    }
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Delete this webhook?')) return
    try {
      await api.delete(`/companies/${companyId}/webhooks/${id}`)
      setWebhooks(prev => prev.filter(w => w.id !== id))
      setToast({ message: 'Webhook deleted', type: 'success' })
    } catch (err: any) {
      setToast({ message: err.message || 'Failed to delete webhook', type: 'error' })
    }
  }

  function toggleEvent(event: string) {
    setEvents(prev => prev.includes(event) ? prev.filter(e => e !== event) : [...prev, event])
  }

  if (loading) return <CardSkeleton count={3} />
  if (!companyId) {
    return (
      <div>
        <div className="flex items-center gap-3 mb-6">
          <Bell size={24} className="text-blue-400" />
          <h2 className="text-2xl font-bold text-white">Webhooks</h2>
        </div>
        <EmptyState
          icon={Bell}
          title="No company yet"
          description="Send real-time notifications to Slack or Teams when achievements, badges, and rewards are triggered."
          action={{ label: 'Create a Company', onClick: () => navigate('/company') }}
        />
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Bell size={24} className="text-blue-400" />
          <h2 className="text-2xl font-bold text-white">Webhooks</h2>
          <span className="text-sm text-gray-500">Slack / Teams</span>
        </div>
        <button onClick={() => setCreating(!creating)} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
          <Plus size={16} /> New Webhook
        </button>
      </div>

      {creating && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-xl">
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Name</label>
              <input value={name} onChange={e => setName(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="Sales Channel" required />
            </div>
            <div className="w-36">
              <label className="block text-sm font-medium text-gray-400 mb-1.5">Platform</label>
              <select value={platform} onChange={e => setPlatform(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500">
                <option value="slack">Slack</option>
                <option value="teams">Teams</option>
              </select>
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Webhook URL</label>
            <input value={url} onChange={e => setUrl(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="https://hooks.slack.com/services/..." required />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-2">Events</label>
            <div className="flex flex-wrap gap-2">
              {eventOptions.map(ev => (
                <button key={ev} type="button" onClick={() => toggleEvent(ev)}
                  className={`px-3 py-1.5 text-xs rounded-lg border transition-colors ${
                    events.includes(ev) ? 'bg-violet-600/20 border-violet-500 text-violet-400' : 'bg-gray-800 border-gray-700 text-gray-500 hover:text-gray-300'
                  }`}>{ev}</button>
              ))}
            </div>
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg">Create</button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg">Cancel</button>
          </div>
        </form>
      )}

      {webhooks.length === 0 && !creating ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">No webhooks configured. Add a Slack or Teams webhook to get notifications.</p></div>
      ) : (
        <div className="space-y-3">
          {webhooks.map(w => (
            <div key={w.id} className="bg-gray-900 border border-gray-800 rounded-xl p-5 flex items-center justify-between group">
              <div>
                <div className="flex items-center gap-3 mb-1">
                  <h3 className="text-white font-medium">{w.name}</h3>
                  <span className={`text-xs px-2 py-0.5 rounded-full ${w.platform === 'slack' ? 'bg-green-500/20 text-green-400' : 'bg-blue-500/20 text-blue-400'}`}>{w.platform}</span>
                </div>
                <p className="text-xs text-gray-600 truncate max-w-md">{w.url}</p>
                <div className="flex gap-1.5 mt-2">
                  {w.events.map(ev => (
                    <span key={ev} className="text-xs px-2 py-0.5 bg-gray-800 text-gray-500 rounded">{ev}</span>
                  ))}
                </div>
              </div>
              <button onClick={() => handleDelete(w.id)} className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all"><Trash2 size={16} /></button>
            </div>
          ))}
        </div>
      )}

      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  )
}
