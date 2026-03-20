import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './hooks/useAuth'
import Layout from './components/Layout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Company from './pages/Company'
import Members from './pages/Members'
import Badges from './pages/Badges'
import Achievements from './pages/Achievements'
import Leaderboard from './pages/Leaderboard'
import GamePlans from './pages/GamePlans'
import GamePlanEditor from './pages/GamePlanEditor'
import Teams from './pages/Teams'
import Challenges from './pages/Challenges'
import Rewards from './pages/Rewards'
import Quests from './pages/Quests'
import Tournaments from './pages/Tournaments'
import Analytics from './pages/Analytics'
import Integrations from './pages/Integrations'
import Webhooks from './pages/Webhooks'
import Invites from './pages/Invites'
import DocsPage from './pages/Docs'

function App() {
  const { session, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <p className="text-gray-500">Loading...</p>
      </div>
    )
  }

  if (!session) {
    return <Login />
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Dashboard />} />
          <Route path="/company" element={<Company />} />
          <Route path="/members" element={<Members />} />
          <Route path="/badges" element={<Badges />} />
          <Route path="/achievements" element={<Achievements />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
          <Route path="/game-plans" element={<GamePlans />} />
          <Route path="/game-plans/:planId" element={<GamePlanEditor />} />
          <Route path="/teams" element={<Teams />} />
          <Route path="/challenges" element={<Challenges />} />
          <Route path="/rewards" element={<Rewards />} />
          <Route path="/quests" element={<Quests />} />
          <Route path="/tournaments" element={<Tournaments />} />
          <Route path="/analytics" element={<Analytics />} />
          <Route path="/integrations" element={<Integrations />} />
          <Route path="/webhooks" element={<Webhooks />} />
          <Route path="/invites" element={<Invites />} />
          <Route path="/docs" element={<DocsPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
