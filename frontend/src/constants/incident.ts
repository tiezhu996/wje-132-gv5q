// 严重等级/事件状态枚举（与后端 backend/internal/constants/incident.go 保持一致）
export const SeverityLevel = {
  NEAR_MISS: 'near_miss',
  MINOR: 'minor',
  MODERATE: 'moderate',
  MAJOR: 'major',
  FATAL: 'fatal',
} as const

export const SeverityText: Record<string, string> = {
  [SeverityLevel.NEAR_MISS]: '未遂',
  [SeverityLevel.MINOR]: '轻微',
  [SeverityLevel.MODERATE]: '一般',
  [SeverityLevel.MAJOR]: '较大',
  [SeverityLevel.FATAL]: '重大',
}

export const IncidentStatus = {
  REPORTED: 'reported',
  INVESTIGATING: 'investigating',
  RESOLVED: 'resolved',
  CLOSED: 'closed',
} as const

export const IncidentStatusText: Record<string, string> = {
  [IncidentStatus.REPORTED]: '已上报',
  [IncidentStatus.INVESTIGATING]: '调查中',
  [IncidentStatus.RESOLVED]: '已整改',
  [IncidentStatus.CLOSED]: '已关闭',
}

export const IncidentCategories = ['坠落', '触电', '物体打击', '坍塌', '机械伤害', '其他']
export const SeverityOptions = Object.entries(SeverityText).map(([value, label]) => ({ label, value }))
export const IncidentStatusOptions = Object.entries(IncidentStatusText).map(([value, label]) => ({ label, value }))
