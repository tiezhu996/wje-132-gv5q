import request from '@/utils/request'

export function login(data: { phone: string; password: string }) {
  return request.post('/auth/login', data)
}

export function register(data: Record<string, unknown>) {
  return request.post('/auth/register', data)
}

export function getMe() {
  return request.get('/users/me')
}

export function updateProfile(data: { name?: string; avatar?: string }) {
  return request.put('/users/me', data)
}
