import request from '@/utils/request'

export function listIncidents(params: Record<string, unknown>) {
  return request.get('/incidents', { params })
}

export function getIncident(id: number) {
  return request.get(`/incidents/${id}`)
}

export function reportIncident(data: Record<string, unknown>) {
  return request.post('/incidents', data)
}

export function assignIncident(id: number) {
  return request.post(`/incidents/${id}/assign`)
}

export function rectifyIncident(id: number, data: { measures: string; deadline?: string }) {
  return request.post(`/incidents/${id}/rectify`, data)
}

export function closeIncident(id: number) {
  return request.post(`/incidents/${id}/close`)
}
