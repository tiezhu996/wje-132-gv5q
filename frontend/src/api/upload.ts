import request from '@/utils/request'

export function uploadImage(file: File) {
  const form = new FormData()
  form.append('file', file)
  return request.post('/upload/image', form, { headers: { 'Content-Type': 'multipart/form-data' } })
}
