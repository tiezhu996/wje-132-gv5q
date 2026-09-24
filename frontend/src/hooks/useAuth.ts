import { useAuthStore } from '@/stores/authStore'

export function useAuth() {
  const auth = useAuthStore()
  const isLoggedIn = !!auth.token
  const role = auth.user?.role || ''
  const isAdmin = role === 'admin'
  const isManager = role === 'admin' || role === 'safety_manager'
  return { auth, isLoggedIn, role, isAdmin, isManager }
}
