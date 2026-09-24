import { Button, Card, Form, Input, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

export default function Login() {
  const navigate = useNavigate()
  const login = useAuthStore((s) => s.login)

  async function onFinish(values: { phone: string; password: string }) {
    try {
      await login(values.phone, values.password)
      message.success('登录成功')
      navigate('/dashboard')
    } catch {
      // 拦截器已提示
    }
  }

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f0f2f5' }}>
      <Card title="建筑施工安全管理平台" style={{ width: 380 }}>
        <Form onFinish={onFinish} initialValues={{ phone: '13800000002', password: 'User@123' }}>
          <Form.Item name="phone" rules={[{ required: true, message: '请输入手机号' }]}>
            <Input placeholder="手机号" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password placeholder="密码" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>登录</Button>
        </Form>
      </Card>
    </div>
  )
}
