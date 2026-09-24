import dayjs from 'dayjs'

export function formatDateTime(value?: string | Date | null): string {
  if (!value) return '-'
  const d = dayjs(value)
  return d.isValid() ? d.format('YYYY-MM-DD HH:mm') : String(value)
}

export function formatDate(value?: string | Date | null): string {
  if (!value) return '-'
  const d = dayjs(value)
  return d.isValid() ? d.format('YYYY-MM-DD') : String(value)
}

// formatOverdueDuration 将逾期秒数格式化为“逾期 X天 X小时”。
export function formatOverdueDuration(seconds?: number | null): string {
  if (!seconds || seconds <= 0) return '-'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  if (days > 0) return `逾期 ${days} 天 ${hours} 小时`
  const minutes = Math.floor((seconds % 3600) / 60)
  if (hours > 0) return `逾期 ${hours} 小时 ${minutes} 分钟`
  return `逾期 ${Math.max(minutes, 1)} 分钟`
}
