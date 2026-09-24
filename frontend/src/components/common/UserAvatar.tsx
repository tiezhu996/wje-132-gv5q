import { Avatar } from 'antd'
import { UserOutlined } from '@ant-design/icons'

export default function UserAvatar({ name, avatar, size = 28 }: { name?: string; avatar?: string; size?: number }) {
  return (
    <span>
      <Avatar size={size} icon={<UserOutlined />} src={avatar} style={{ marginRight: 6 }} />
      {name || '-'}
    </span>
  )
}
