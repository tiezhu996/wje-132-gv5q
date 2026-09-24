export interface User {
  id: number
  phone: string
  name: string
  avatar: string
  role: string
  created_at: string
}

export interface SafetyIncident {
  id: number
  title: string
  description: string
  occurred_at: string
  site_id: string
  area: string
  severity_level: string
  category: string
  involved_user_ids: string[]
  photo_urls: string[]
  status: string
  rectification_measures: string
  rectification_deadline: string | null
  reporter_id: number
  created_at: string
}

export interface SafetyInspection {
  id: number
  name: string
  inspection_type: string
  area: string
  inspection_date: string
  inspector_id: number
  total_score: number
  status: string
  issue_count: number
  passed_count: number
  created_at: string
}

export interface InspectionItem {
  id: number
  inspection_id: number
  item_name: string
  passed: boolean
  remark: string
  photo_url: string
}

export interface SafetyTraining {
  id: number
  topic: string
  training_type: string
  training_date: string
  duration_hours: number
  trainer: string
  location: string
  content_summary: string
  participant_ids: string[]
  assessment_method: string
  pass_rate: number
  created_at: string
}

export interface WorkerCertification {
  id: number
  user_id: number
  cert_type: string
  cert_no: string
  issue_org: string
  issue_date: string | null
  valid_until: string | null
  cert_photo_url: string
  status: string
  created_at: string
}

export interface AuditLog {
  id: number
  operator_id: number
  operator_name: string
  action: string
  entity_type: string
  entity_id: string
  detail: string
  ip: string
  created_at: string
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}
