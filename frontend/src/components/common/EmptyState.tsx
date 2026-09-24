import { Empty } from 'antd'

export default function EmptyState({ text = '暂无数据' }: { text?: string }) {
  return <Empty description={text} />
}
