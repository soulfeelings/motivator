import { Building2 } from 'lucide-react'
import CompanyForm from '../components/CompanyForm'

interface SetupProps {
  onComplete: () => void
}

export default function Setup({ onComplete }: SetupProps) {
  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center p-6">
      <div className="w-full max-w-lg">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">Motivator</h1>
          <p className="text-gray-500">Admin Panel</p>
        </div>

        <div className="bg-gray-900/50 border border-gray-800 rounded-2xl p-8">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-lg bg-violet-600/20 flex items-center justify-center">
              <Building2 size={20} className="text-violet-400" />
            </div>
            <div>
              <h2 className="text-lg font-semibold text-white">Create your company</h2>
              <p className="text-sm text-gray-500">Set up your workspace to get started</p>
            </div>
          </div>

          <CompanyForm onSuccess={onComplete} />
        </div>

        <p className="text-center text-xs text-gray-600 mt-6">
          You can invite team members and configure gamification after setup.
        </p>
      </div>
    </div>
  )
}
