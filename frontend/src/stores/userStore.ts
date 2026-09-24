import { create } from 'zustand'
import { getMe, updateProfile } from '@/api/user'
import type { User } from '@/types'

interface UserState {
  me: User | null
  fetchMe: () => Promise<User>
  updateMe: (data: { name?: string; avatar?: string }) => Promise<User>
}

export const useUserStore = create<UserState>((set) => ({
  me: null,
  async fetchMe() {
    const res: any = await getMe()
    set({ me: res.data })
    return res.data as User
  },
  async updateMe(data) {
    const res: any = await updateProfile(data)
    set({ me: res.data })
    return res.data as User
  },
}))
