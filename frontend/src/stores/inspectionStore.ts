import { create } from 'zustand'
import { listInspections } from '@/api/inspection'
import type { SafetyInspection } from '@/types'

interface InspectionState {
  list: SafetyInspection[]
  total: number
  fetchList: (params?: Record<string, unknown>) => Promise<void>
}

export const useInspectionStore = create<InspectionState>((set) => ({
  list: [],
  total: 0,
  async fetchList(params = {}) {
    const res: any = await listInspections(params)
    set({ list: res.data.list, total: res.data.total })
  },
}))
