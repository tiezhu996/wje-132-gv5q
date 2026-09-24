import { Tag } from 'antd'
import { IncidentStatusText } from '@/constants/incident'

const statusColor: Record<string, string> = {
  reported: 'red',
  investigating: 'orange',
  resolved: 'blue',
  closed: 'green',
  scheduled: 'default',
  in_progress: 'processing',
  completed: 'green',
  failed: 'red',
  pending: 'orange',
  approved: 'green',
  rejected: 'red',
  expired: 'default',
}

export default function StatusBadge({ status }: { status: string }) {
  const text = IncidentStatusText[status] || status
  return <Tag color={statusColor[status] || 'default'}>{text}</Tag>
}
