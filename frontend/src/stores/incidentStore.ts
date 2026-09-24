import { create } from 'zustand'
import { listIncidents, listOverdueIncidents } from '@/api/incident'
import type { OverdueIncident, SafetyIncident } from '@/types'

interface IncidentState {
  list: SafetyIncident[]
  total: number
  overdueList: OverdueIncident[]
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
    const res: any = await listOverdueIncidents()
    set({ overdueList: res.data.list })
  },
}))
