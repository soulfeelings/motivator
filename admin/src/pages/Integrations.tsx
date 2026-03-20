import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import { Plug, Plus, Trash2, ArrowRight, Copy, Check } from 'lucide-react'

interface Integration {
  id: string
  provider: string
  name: string
  webhook_secret: string
  is_active: boolean
  created_at: string
}

interface Mapping {
  id: string
  external_event: string
  metric: string
  user_field: string
  transform: { value?: number }
}

interface Event {
  id: string
  external_event: string
  metric?: string
  user_email?: string
  value: number
  processed: boolean
  error?: string
  received_at: string
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

const providers = [
  { id: 'jira', name: 'Jira', color: 'text-blue-400', events: ['issue_created', 'issue_updated', 'issue_resolved'] },
  { id: 'github', name: 'GitHub', color: 'text-gray-300', events: ['push', 'pr_opened', 'pr_merged', 'issue_opened', 'issue_closed'] },
  { id: 'salesforce', name: 'Salesforce', color: 'text-cyan-400', events: ['deal_closed', 'lead_converted', 'opportunity_won'] },
  { id: 'zendesk', name: 'Zendesk', color: 'text-green-400', events: ['ticket_solved', 'ticket_created', 'satisfaction_rated'] },
  { id: 'custom', name: 'Custom', color: 'text-violet-400', events: [] },
]

export default function Integrations() {
  const [integrations, setIntegrations] = useState<Integration[]>([])
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [provider, setProvider] = useState('jira')
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [mappings, setMappings] = useState<Mapping[]>([])
  const [events, setEvents] = useState<Event[]>([])
  const [newEvent, setNewEvent] = useState('')
  const [newMetric, setNewMetric] = useState('')
  const [copied, setCopied] = useState(false)

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const res = await api.get<Integration[]>(`/companies/${cid}/integrations`)
        setIntegrations(res ?? [])
      }
    } catch {} finally { setLoading(false) }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    if (!companyId) return
    setError('')
    try {
      const i = await api.post<Integration>(`/companies/${companyId}/integrations`, { provider, name })
      setIntegrations(prev => [...prev, i])
      setCreating(false)
      setName('')
      setSelectedId(i.id)
      loadDetail(i.id)
    } catch (err: any) { setError(err.message) }
  }

  async function loadDetail(id: string) {
    if (!companyId) return
    const [m, e] = await Promise.all([
      api.get<Mapping[]>(`/companies/${companyId}/integrations/${id}/mappings`),
      api.get<Event[]>(`/companies/${companyId}/integrations/${id}/events`),
    ])
    setMappings(m ?? [])
    setEvents(e ?? [])
  }

  async function addMapping() {
    if (!companyId || !selectedId || !newEvent || !newMetric) return
    const m = await api.post<Mapping>(`/companies/${companyId}/integrations/${selectedId}/mappings`, {
      external_event: newEvent, metric: newMetric, transform: { value: 1 },
    })
    setMappings(prev => [...prev, m])
    setNewEvent('')
    setNewMetric('')
  }

  async function deleteMapping(id: string) {
    if (!companyId || !selectedId) return
    await api.delete(`/companies/${companyId}/integrations/${selectedId}/mappings/${id}`)
    setMappings(prev => prev.filter(m => m.id !== id))
  }

  async function handleDelete(id: string) {
    if (!companyId || !confirm('Delete this integration?')) return
    await api.delete(`/companies/${companyId}/integrations/${id}`)
    setIntegrations(prev => prev.filter(i => i.id !== id))
    if (selectedId === id) { setSelectedId(null); setMappings([]); setEvents([]) }
  }

  function copyWebhookUrl(secret: string) {
    navigator.clipboard.writeText(`${window.location.origin}/api/v1/webhooks/inbound/${secret}`)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const selected = integrations.find(i => i.id === selectedId)
  const providerInfo = providers.find(p => p.id === selected?.provider)

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!companyId) return <div className="bg-gray-900 border border-gray-800 rounded-xl p-8 text-center"><p className="text-gray-400">Create a company first.</p></div>

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Plug size={24} className="text-violet-400" />
          <h2 className="text-2xl font-bold text-white">Integrations</h2>
        </div>
        <button onClick={() => setCreating(!creating)} className="inline-flex items-center gap-2 px-4 py-2 bg-violet-600 hover:bg-violet-500 text-white text-sm font-medium rounded-lg transition-colors">
          <Plus size={16} /> Add Integration
        </button>
      </div>

      {creating && (
        <form onSubmit={handleCreate} className="bg-gray-900 border border-gray-800 rounded-xl p-5 mb-6 space-y-4 max-w-lg">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-2">Provider</label>
            <div className="grid grid-cols-5 gap-2">
              {providers.map(p => (
                <button key={p.id} type="button" onClick={() => setProvider(p.id)}
                  className={`px-3 py-2 text-xs font-medium rounded-lg border transition-colors ${
                    provider === p.id ? 'bg-violet-600/20 border-violet-500 text-violet-400' : 'bg-gray-800 border-gray-700 text-gray-500 hover:text-gray-300'
                  }`}>{p.name}</button>
              ))}
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1.5">Name</label>
            <input value={name} onChange={e => setName(e.target.value)} className="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-violet-500" placeholder="My Jira Project" required />
          </div>
          {error && <p className="text-sm text-red-400">{error}</p>}
          <div className="flex gap-3">
            <button type="submit" className="px-4 py-2.5 bg-violet-600 hover:bg-violet-500 text-white font-medium rounded-lg">Create</button>
            <button type="button" onClick={() => setCreating(false)} className="px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg">Cancel</button>
          </div>
        </form>
      )}

      <div className="flex gap-6">
        <div className="w-72 space-y-2">
          {integrations.map(i => (
            <div key={i.id} onClick={() => { setSelectedId(i.id); loadDetail(i.id) }}
              className={`p-4 rounded-xl border cursor-pointer transition-colors group ${
                selectedId === i.id ? 'bg-gray-900 border-violet-500/50' : 'bg-gray-900 border-gray-800 hover:border-gray-700'
              }`}>
              <div className="flex items-center justify-between">
                <div>
                  <span className={`text-xs font-bold uppercase ${providers.find(p => p.id === i.provider)?.color ?? 'text-gray-400'}`}>{i.provider}</span>
                  <p className="text-white text-sm font-medium">{i.name}</p>
                </div>
                <button onClick={(e) => { e.stopPropagation(); handleDelete(i.id) }} className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100"><Trash2 size={14} /></button>
              </div>
            </div>
          ))}
          {integrations.length === 0 && !creating && <p className="text-gray-500 text-sm text-center py-8">No integrations yet</p>}
        </div>

        {selected && (
          <div className="flex-1 space-y-6">
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h3 className="text-white font-medium mb-3">Webhook URL</h3>
              <div className="flex items-center gap-2">
                <code className="flex-1 text-xs bg-gray-800 px-3 py-2 rounded-lg text-gray-400 truncate">
                  {window.location.origin}/api/v1/webhooks/inbound/{selected.webhook_secret}
                </code>
                <button onClick={() => copyWebhookUrl(selected.webhook_secret)} className="p-2 text-gray-500 hover:text-violet-400">
                  {copied ? <Check size={16} className="text-emerald-400" /> : <Copy size={16} />}
                </button>
              </div>
              <p className="text-xs text-gray-600 mt-2">Paste this URL in your {selected.provider} webhook settings.</p>
            </div>

            <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h3 className="text-white font-medium mb-3">Event Mappings</h3>
              <div className="space-y-2 mb-4">
                {mappings.map(m => (
                  <div key={m.id} className="flex items-center gap-3 p-3 bg-gray-800/50 rounded-lg group">
                    <span className="text-sm text-amber-400 font-mono">{m.external_event}</span>
                    <ArrowRight size={14} className="text-gray-600" />
                    <span className="text-sm text-violet-400 font-mono">{m.metric}</span>
                    <span className="text-xs text-gray-600 ml-auto">+{m.transform?.value ?? 1}</span>
                    <button onClick={() => deleteMapping(m.id)} className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100"><Trash2 size={12} /></button>
                  </div>
                ))}
              </div>
              <div className="flex gap-2">
                <select value={newEvent} onChange={e => setNewEvent(e.target.value)} className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm">
                  <option value="">Event...</option>
                  {(providerInfo?.events ?? []).map(ev => <option key={ev} value={ev}>{ev}</option>)}
                  <option value="custom">Custom...</option>
                </select>
                <input value={newMetric} onChange={e => setNewMetric(e.target.value)} placeholder="metric" className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm" />
                <button onClick={addMapping} disabled={!newEvent || !newMetric} className="px-3 py-2 bg-violet-600 hover:bg-violet-500 disabled:opacity-50 text-white text-sm rounded-lg">Add</button>
              </div>
            </div>

            <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h3 className="text-white font-medium mb-3">Recent Events</h3>
              {events.length === 0 ? (
                <p className="text-gray-500 text-sm">No events received yet. Configure the webhook URL in {selected.provider}.</p>
              ) : (
                <div className="space-y-1 max-h-64 overflow-auto">
                  {events.map(e => (
                    <div key={e.id} className="flex items-center gap-3 text-xs py-2 border-b border-gray-800/50">
                      <span className={`w-2 h-2 rounded-full ${e.processed ? 'bg-emerald-400' : 'bg-red-400'}`} />
                      <span className="text-gray-400">{new Date(e.received_at).toLocaleTimeString()}</span>
                      <span className="text-amber-400 font-mono">{e.external_event}</span>
                      {e.metric && <><ArrowRight size={10} className="text-gray-600" /><span className="text-violet-400 font-mono">{e.metric}</span></>}
                      {e.user_email && <span className="text-gray-500">{e.user_email}</span>}
                      <span className="text-white font-medium ml-auto">+{e.value}</span>
                      {e.error && <span className="text-red-400 truncate max-w-[150px]">{e.error}</span>}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
