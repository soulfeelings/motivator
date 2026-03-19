import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import FlowEditor from '../components/flow/FlowEditor'
import { ArrowLeft } from 'lucide-react'
import type { Node, Edge } from '@xyflow/react'

interface GamePlan {
  id: string
  name: string
  description?: string
  flow_data: { nodes: Node[]; edges: Edge[] }
  is_active: boolean
}

interface MeResponse {
  memberships: Array<{ company_id: string }>
}

export default function GamePlanEditor() {
  const { planId } = useParams<{ planId: string }>()
  const navigate = useNavigate()
  const [plan, setPlan] = useState<GamePlan | null>(null)
  const [companyId, setCompanyId] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    load()
  }, [planId])

  async function load() {
    try {
      const me = await api.get<MeResponse>('/me')
      if (me.memberships?.length > 0) {
        const cid = me.memberships[0].company_id
        setCompanyId(cid)
        const gp = await api.get<GamePlan>(`/companies/${cid}/game-plans/${planId}`)
        setPlan(gp)
      }
    } catch {
      navigate('/game-plans')
    } finally {
      setLoading(false)
    }
  }

  async function handleSave(nodes: Node[], edges: Edge[]) {
    if (!companyId || !planId) return
    setSaving(true)
    try {
      await api.patch(`/companies/${companyId}/game-plans/${planId}`, {
        flow_data: { nodes, edges },
      })
    } catch (err) {
      console.error('Failed to save:', err)
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!plan) return <p className="text-gray-500">Game plan not found.</p>

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <button
          onClick={() => navigate('/game-plans')}
          className="p-2 text-gray-400 hover:text-white transition-colors"
        >
          <ArrowLeft size={20} />
        </button>
        <div>
          <h2 className="text-2xl font-bold text-white">{plan.name}</h2>
          {plan.description && <p className="text-sm text-gray-500">{plan.description}</p>}
        </div>
        <div className="ml-auto">
          <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium ${
            plan.is_active ? 'bg-emerald-500/20 text-emerald-400' : 'bg-gray-500/20 text-gray-400'
          }`}>
            {plan.is_active ? 'Active' : 'Draft'}
          </span>
        </div>
      </div>
      <FlowEditor
        initialNodes={plan.flow_data.nodes}
        initialEdges={plan.flow_data.edges}
        onSave={handleSave}
        saving={saving}
      />
    </div>
  )
}
