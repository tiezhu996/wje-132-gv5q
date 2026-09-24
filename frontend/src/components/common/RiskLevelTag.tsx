import { Tag } from 'antd'
import { SeverityText } from '@/constants/incident'
import { getSeverityColor } from '@/utils/getSeverityColor'

export default function RiskLevelTag({ level }: { level: string }) {
  return <Tag color={getSeverityColor(level)}>{SeverityText[level] || level}</Tag>
}
