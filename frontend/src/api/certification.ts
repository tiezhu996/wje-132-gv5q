import request from '@/utils/request'

export function listCertifications(params: Record<string, unknown>) {
  return request.get('/certifications', { params })
}

export function listCertsByUser(userId: number) {
  return request.get('/certifications/by-user', { params: { user_id: userId } })
}

export function submitCertification(data: Record<string, unknown>) {
  return request.post('/certifications', data)
}

export function reviewCertification(id: number, status: string) {
  return request.post(`/certifications/${id}/review`, { status })
}
