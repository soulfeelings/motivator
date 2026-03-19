import { Outlet, Link, useLocation } from 'react-router-dom'
import { Building2, Users, Mail, LayoutDashboard, LogOut, Award, Target, Trophy, Swords, Gift } from 'lucide-react'
import { useAuth } from '../hooks/useAuth'

const nav = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/company', label: 'Company', icon: Building2 },
  { to: '/members', label: 'Members', icon: Users },
  { to: '/badges', label: 'Badges', icon: Award },
  { to: '/achievements', label: 'Achievements', icon: Target },
  { to: '/leaderboard', label: 'Leaderboard', icon: Trophy },
  { to: '/challenges', label: 'Challenges', icon: Swords },
  { to: '/rewards', label: 'Rewards', icon: Gift },
  { to: '/invites', label: 'Invites', icon: Mail },
]

export default function Layout() {
  const { pathname } = useLocation()
  const { session, signOut } = useAuth()

  return (
    <div className="flex h-screen bg-gray-950 text-gray-100">
      <aside className="w-64 border-r border-gray-800 flex flex-col">
        <div className="p-6 border-b border-gray-800">
          <h1 className="text-xl font-bold tracking-tight text-white">Motivator</h1>
          <p className="text-xs text-gray-500 mt-1">Admin Panel</p>
        </div>
        <nav className="flex-1 p-4 space-y-1">
          {nav.map(({ to, label, icon: Icon }) => {
            const active = pathname === to
            return (
              <Link
                key={to}
                to={to}
                className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                  active
                    ? 'bg-violet-600/20 text-violet-400'
                    : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/50'
                }`}
              >
                <Icon size={18} />
                {label}
              </Link>
            )
          })}
        </nav>
        <div className="p-4 border-t border-gray-800">
          <div className="text-xs text-gray-500 truncate mb-3">{session?.user?.email}</div>
          <button
            onClick={signOut}
            className="flex items-center gap-2 text-sm text-gray-400 hover:text-red-400 transition-colors"
          >
            <LogOut size={16} />
            Sign out
          </button>
        </div>
      </aside>
      <main className="flex-1 overflow-auto">
        <div className="p-8 max-w-5xl">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
