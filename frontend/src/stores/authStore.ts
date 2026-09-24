import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { login as apiLogin } from '@/api/user'
import { getMe } from '@/api/user'
import type { User } from '@/types'

interface AuthState {
  token: string
  user: User | null
  login: (phone: string, password: string) => Promise<void>
  fetchMe: () => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: '',
      user: null,
      async login(phone: string, password: string) {
        const res: any = await apiLogin({ phone, password })
        set({ token: res.data.token, user: res.data.user })
      },
      async fetchMe() {
        const res: any = await getMe()
        set({ user: res.data })
      },
      logout() {
        set({ token: '', user: null })
      },
    }),
    { name: 'safety-auth' },
  ),
)
