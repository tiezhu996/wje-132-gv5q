import { create } from 'zustand'
import { listIncidents } from '@/api/incident'
import type { SafetyIncident } from '@/types'

interface IncidentState {
  list: SafetyIncident[]
  total: number
  fetchList: (params?: Record<string, unknown>) => Promise<void>
}

export const useIncidentStore = create<IncidentState>((set) => ({
  list: [],
  total: 0,
  async fetchList(params = {}) {
    const res: any = await listIncidents(params)
    set({ list: res.data.list, total: res.data.total })
  },
}))
