import { Layout as AntLayout, Menu, Avatar, Dropdown, Space } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

const items = [
  { key: '/dashboard', label: '安全概览' },
  { key: '/incidents', label: '事件管理' },
  { key: '/inspections', label: '检查管理' },
  { key: '/trainings', label: '培训管理' },
  { key: '/certifications', label: '资质审核' },
  { key: '/audit-logs', label: '审计日志' },
  { key: '/profile', label: '个人中心' },
]

export default function Layout() {
  const location = useLocation()
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const selected = items.find((it) => location.pathname.startsWith(it.key))?.key || '/dashboard'

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <AntLayout.Header style={{ display: 'flex', alignItems: 'center', gap: 24, background: '#001529' }}>
        <div style={{ color: '#fff', fontSize: 18, fontWeight: 700, cursor: 'pointer' }} onClick={() => navigate('/dashboard')}>
          施工安全管理平台
        </div>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selected]}
          items={items}
          style={{ flex: 1, minWidth: 0 }}
          onClick={({ key }) => navigate(key)}
        />
        <Dropdown
          menu={{ items: [{ key: 'logout', label: '退出登录', onClick: () => { logout(); navigate('/login') } }] }}
        >
          <Space style={{ color: '#fff', cursor: 'pointer' }}>
            <Avatar size="small" icon={<UserOutlined />} src={user?.avatar} />
            {user?.name || user?.phone}
          </Space>
        </Dropdown>
      </AntLayout.Header>
      <AntLayout.Content style={{ padding: 24 }}>
        <Outlet />
      </AntLayout.Content>
    </AntLayout>
  )
}
