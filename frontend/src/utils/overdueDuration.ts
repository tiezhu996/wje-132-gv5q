// 将逾期秒数格式化为“X天X小时”（不足 1 小时按分钟展示）。
export function formatOverdueDuration(seconds: number): string {
  if (seconds <= 0) return '-'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  if (days > 0) return `${days}天${hours}小时`
  if (hours > 0) return `${hours}小时${minutes}分钟`
  return `${Math.max(minutes, 1)}分钟`
}
