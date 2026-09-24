import { create } from 'zustand'
import { listTrainings } from '@/api/training'
import type { SafetyTraining } from '@/types'

interface TrainingState {
  list: SafetyTraining[]
  total: number
  fetchList: (params?: Record<string, unknown>) => Promise<void>
}

export const useTrainingStore = create<TrainingState>((set) => ({
  list: [],
  total: 0,
  async fetchList(params = {}) {
    const res: any = await listTrainings(params)
    set({ list: res.data.list, total: res.data.total })
  },
}))
