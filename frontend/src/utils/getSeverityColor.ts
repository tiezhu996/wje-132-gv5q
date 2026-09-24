// 严重等级颜色映射（与 RiskLevelTag 共用）
export const SeverityColor: Record<string, string> = {
  near_miss: 'blue',
  minor: 'green',
  moderate: 'orange',
  major: 'volcano',
  fatal: 'red',
}

export const SeverityBgColor: Record<string, string> = {
  near_miss: '#e6f4ff',
  minor: '#f6ffed',
  moderate: '#fff7e6',
  major: '#fff2e8',
  fatal: '#fff1f0',
}

export function getSeverityColor(level: string): string {
  return SeverityColor[level] || 'default'
}
