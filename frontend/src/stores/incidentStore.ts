import { create } from 'zustand'
import { listIncidents, listOverdueAcceptance } from '@/api/incident'
import type { OverdueAcceptanceItem, SafetyIncident } from '@/types'

interface IncidentState {
  list: SafetyIncident[]
  total: number
  overdueList: OverdueAcceptanceItem[]
  fetchList: (params?: Record<string, unknown>) => Promise<void>
  fetchOverdue: () => Promise<void>
}

export const useIncidentStore = create<IncidentState>((set) => ({
  list: [],
  total: 0,
  overdueList: [],
  async fetchList(params = {}) {
    const res: any = await listIncidents(params)
    set({ list: res.data.list, total: res.data.total })
  },
  async fetchOverdue() {
    const res: any = await listOverdueAcceptance()
    set({ overdueList: res.data.list })
  },
}))
