import request from '@/utils/request'

export function listInspections(params: Record<string, unknown>) {
  return request.get('/inspections', { params })
}

export function getInspection(id: number) {
  return request.get(`/inspections/${id}`)
}

export function createInspection(data: Record<string, unknown>) {
  return request.post('/inspections', data)
}

export function executeInspection(id: number, items: unknown[]) {
  return request.post(`/inspections/${id}/execute`, { items })
}

export function getInspectionReport(id: number) {
  return request.get(`/inspections/${id}/report`)
}
