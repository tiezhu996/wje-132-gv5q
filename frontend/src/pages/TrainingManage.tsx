import { useEffect, useState } from 'react'
import { Button, Card, DatePicker, Form, Input, InputNumber, Modal, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { createTraining, recordTraining } from '@/api/training'
import { useTrainingStore } from '@/stores/trainingStore'
import { formatDate } from '@/utils/dateFormat'
import type { SafetyTraining } from '@/types'

const typeOptions = [
  { label: '入场教育', value: 'induction' },
  { label: '常规培训', value: 'regular' },
  { label: '专项培训', value: 'special' },
  { label: '应急演练', value: 'emergency' },
]

export default function TrainingManage() {
  const store = useTrainingStore()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [open, setOpen] = useState(false)
  const [recordId, setRecordId] = useState<number>()
  const [passRate, setPassRate] = useState(100)
  const [form] = Form.useForm()

  useEffect(() => {
    store.fetchList({ page, page_size: pageSize })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize])

  async function onCreate() {
    const values = await form.validateFields()
    await createTraining({
      topic: values.topic,
      training_type: values.training_type,
      training_date: values.training_date.format('YYYY-MM-DDTHH:mm:ss'),
      duration_hours: values.duration_hours,
      trainer: values.trainer,
      location: values.location,
      content_summary: values.content_summary,
      assessment_method: values.assessment_method,
    })
    message.success('培训创建成功')
    setOpen(false)
    form.resetFields()
    setPage(1)
  }

  async function onRecord() {
    if (!recordId) return
    await recordTraining(recordId, { pass_rate: passRate })
    message.success('成绩已记录')
    setRecordId(undefined)
    store.fetchList({ page, page_size: pageSize })
  }

  return (
    <Card>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>创建培训</Button>
      </Space>
      <Table<SafetyTraining>
        rowKey="id"
        dataSource={store.list}
        pagination={{ current: page, pageSize, total: store.total, onChange: (p, ps) => { setPage(p); setPageSize(ps) } }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 70 },
          { title: '主题', dataIndex: 'topic' },
          { title: '类型', dataIndex: 'training_type', render: (v) => typeOptions.find((o) => o.value === v)?.label || v },
          { title: '培训日期', dataIndex: 'training_date', render: (v) => formatDate(v) },
          { title: '讲师', dataIndex: 'trainer' },
          { title: '地点', dataIndex: 'location' },
          { title: '考核方式', dataIndex: 'assessment_method' },
          { title: '通过率', dataIndex: 'pass_rate', render: (v) => `${v}%` },
          {
            title: '操作',
            render: (_, row) => (
              <Button size="small" type="primary" onClick={() => setRecordId(row.id)}>记录成绩</Button>
            ),
          },
        ]}
      />
      <Modal title="创建培训" open={open} onOk={onCreate} onCancel={() => setOpen(false)} width={560}>
        <Form form={form} layout="vertical">
          <Form.Item name="topic" label="培训主题" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="training_type" label="培训类型" rules={[{ required: true }]}><Select options={typeOptions} /></Form.Item>
          <Form.Item name="training_date" label="培训日期" rules={[{ required: true }]}><DatePicker showTime style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="duration_hours" label="时长(小时)"><InputNumber min={1} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="trainer" label="讲师"><Input /></Form.Item>
          <Form.Item name="location" label="地点"><Input /></Form.Item>
          <Form.Item name="content_summary" label="内容摘要"><Input.TextArea rows={3} /></Form.Item>
          <Form.Item name="assessment_method" label="考核方式"><Select options={['笔试', '实操', '口试'].map((v) => ({ label: v, value: v }))} /></Form.Item>
        </Form>
      </Modal>
      <Modal title="记录培训成绩" open={!!recordId} onOk={onRecord} onCancel={() => setRecordId(undefined)} width={420}>
        <div>
          <div style={{ marginBottom: 8 }}>通过率（%）：</div>
          <InputNumber min={0} max={100} value={passRate} onChange={(v) => setPassRate(v || 0)} style={{ width: '100%' }} />
        </div>
      </Modal>
    </Card>
  )
}
