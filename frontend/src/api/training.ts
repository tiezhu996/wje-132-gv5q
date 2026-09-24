import request from '@/utils/request'

export function listTrainings(params: Record<string, unknown>) {
  return request.get('/trainings', { params })
}

export function getTraining(id: number) {
  return request.get(`/trainings/${id}`)
}

export function createTraining(data: Record<string, unknown>) {
  return request.post('/trainings', data)
}

export function recordTraining(id: number, data: { participant_ids?: string[]; pass_rate: number }) {
  return request.post(`/trainings/${id}/record`, data)
}
