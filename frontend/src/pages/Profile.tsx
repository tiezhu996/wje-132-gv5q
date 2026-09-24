import { useEffect, useState } from 'react'
import { Avatar, Button, Card, Form, Input, message, Space } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import AvatarUploader from '@/components/common/AvatarUploader'
import { useUserStore } from '@/stores/userStore'

export default function Profile() {
  const store = useUserStore()
  const [form] = Form.useForm()
  const [avatar, setAvatar] = useState('')

  useEffect(() => {
    store.fetchMe().then((me) => {
      form.setFieldsValue({ name: me.name, phone: me.phone })
      setAvatar(me.avatar || '')
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function onSave() {
    const values = await form.validateFields()
    await store.updateMe({ ...values, avatar })
    message.success('保存成功')
  }

  return (
    <Card style={{ maxWidth: 480 }}>
      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        <Avatar size={72} icon={<UserOutlined />} src={avatar} />
        <AvatarUploader onUploaded={(u) => setAvatar(u)} />
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="姓名"><Input /></Form.Item>
          <Form.Item name="phone" label="手机号"><Input disabled /></Form.Item>
          <Button type="primary" onClick={onSave}>保存</Button>
        </Form>
      </Space>
    </Card>
  )
}
