import { useEffect, useState, useCallback } from 'react'
import type { Session } from '@supabase/supabase-js'
import { supabase } from '../lib/supabase'
import { api } from '../lib/api'

interface Membership {
  id: string
  company_id: string
  role: string
}

export function useAuth() {
  const [session, setSession] = useState<Session | null>(null)
  const [loading, setLoading] = useState(true)
  const [memberships, setMemberships] = useState<Membership[] | null>(null)

  const fetchMemberships = useCallback(async () => {
    try {
      const me = await api.get<{ memberships: Membership[] }>('/me')
      setMemberships(me.memberships ?? [])
    } catch {
      setMemberships([])
    }
  }, [])

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session)
      if (session) {
        fetchMemberships().then(() => setLoading(false))
      } else {
        setLoading(false)
      }
    })

    const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session)
      if (session) {
        fetchMemberships()
      } else {
        setMemberships(null)
      }
    })

    return () => subscription.unsubscribe()
  }, [fetchMemberships])

  const signOut = () => supabase.auth.signOut()

  const hasCompany = (memberships?.length ?? 0) > 0

  return { session, loading, signOut, memberships, hasCompany, refreshMemberships: fetchMemberships }
}
