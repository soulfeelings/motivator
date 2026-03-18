import { useEffect, useState } from 'react'
import { Building2, Users, Mail } from 'lucide-react'
import { api } from '../lib/api'

interface MeResponse {
  user_id: string
  email: string
  memberships: Array<{
    id: string
    company_id: string
    role: string
  }>
}

export default function Dashboard() {
  const [me, setMe] = useState<MeResponse | null>(null)

  useEffect(() => {
    api.get<MeResponse>('/me').then(setMe).catch(() => {})
  }, [])

  const stats = [
    { label: 'Companies', value: me?.memberships?.length ?? 0, icon: Building2, color: 'text-violet-400' },
    { label: 'Members', value: '—', icon: Users, color: 'text-emerald-400' },
    { label: 'Pending Invites', value: '—', icon: Mail, color: 'text-amber-400' },
  ]

  return (
    <div>
      <h2 className="text-2xl font-bold text-white mb-6">Dashboard</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
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
    </div>
  )
}
