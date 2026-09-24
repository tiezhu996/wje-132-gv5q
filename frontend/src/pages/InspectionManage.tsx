import { useEffect, useState } from 'react'
import { Button, Card, DatePicker, Form, Input, InputNumber, Modal, Select, Space, Table, Checkbox, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { createInspection, executeInspection, getInspection } from '@/api/inspection'
import { useInspectionStore } from '@/stores/inspectionStore'
import StatusBadge from '@/components/common/StatusBadge'
import { formatDate } from '@/utils/dateFormat'
import type { SafetyInspection, InspectionItem } from '@/types'

const typeOptions = [
  { label: '例行检查', value: 'routine' },
  { label: '专项检查', value: 'special' },
  { label: '班前检查', value: 'pre_shift' },
  { label: '应急检查', value: 'emergency' },
]

export default function InspectionManage() {
  const store = useInspectionStore()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [statusFilter, setStatusFilter] = useState('')
  const [open, setOpen] = useState(false)
  const [execId, setExecId] = useState<number>()
  const [items, setItems] = useState<InspectionItem[]>([])
  const [form] = Form.useForm()

  useEffect(() => {
    store.fetchList({ page, page_size: pageSize, status: statusFilter || undefined })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, statusFilter])

  async function onCreate() {
    const values = await form.validateFields()
    await createInspection({
      name: values.name,
      inspection_type: values.inspection_type,
      area: values.area,
      inspection_date: values.inspection_date.format('YYYY-MM-DDTHH:mm:ss'),
      inspector_id: values.inspector_id,
      items: (values.items || []).map((name: string) => ({ item_name: name })),
    })
    message.success('检查计划创建成功')
    setOpen(false)
    form.resetFields()
    setPage(1)
  }

  async function openExecute(id: number) {
    const res: any = await getInspection(id)
    setItems(res.data.items)
    setExecId(id)
  }

  async function onExecute() {
    if (!execId) return
    await executeInspection(execId, items)
    message.success('检查执行完成')
    setExecId(undefined)
    store.fetchList({ page, page_size: pageSize, status: statusFilter || undefined })
  }

  return (
    <Card>
      <Space style={{ marginBottom: 16 }}>
        <Select placeholder="状态" allowClear style={{ width: 130 }} options={[
          { label: '待执行', value: 'scheduled' },
          { label: '执行中', value: 'in_progress' },
          { label: '已完成', value: 'completed' },
          { label: '不合格', value: 'failed' },
        ]} onChange={(v) => { setStatusFilter(v || ''); setPage(1) }} />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>创建检查计划</Button>
      </Space>
      <Table<SafetyInspection>
        rowKey="id"
        dataSource={store.list}
        pagination={{ current: page, pageSize, total: store.total, onChange: (p, ps) => { setPage(p); setPageSize(ps) } }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 70 },
          { title: '名称', dataIndex: 'name' },
          { title: '类型', dataIndex: 'inspection_type', render: (v) => typeOptions.find((o) => o.value === v)?.label || v },
          { title: '区域', dataIndex: 'area' },
          { title: '检查日期', dataIndex: 'inspection_date', render: (v) => formatDate(v) },
          { title: '得分', dataIndex: 'total_score', render: (v) => `${v} 分` },
          { title: '状态', dataIndex: 'status', render: (v) => <StatusBadge status={v} /> },
          {
            title: '操作',
            render: (_, row) =>
              row.status === 'scheduled' || row.status === 'in_progress' ? (
                <Button size="small" type="primary" onClick={() => openExecute(row.id)}>执行检查</Button>
              ) : (
                <a onClick={() => openExecute(row.id)}>查看</a>
              ),
          },
        ]}
      />
      <Modal title="创建检查计划" open={open} onOk={onCreate} onCancel={() => setOpen(false)} width={560}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="检查名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="inspection_type" label="检查类型" rules={[{ required: true }]}><Select options={typeOptions} /></Form.Item>
          <Form.Item name="area" label="区域"><Input /></Form.Item>
          <Form.Item name="inspection_date" label="检查日期" rules={[{ required: true }]}><DatePicker showTime style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="inspector_id" label="检查人ID" rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="items" label="检查项（每行一个）">
            <Select mode="tags" placeholder="输入后回车添加检查项" open={false} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal title="执行检查" open={!!execId} onOk={onExecute} onCancel={() => setExecId(undefined)} width={560}>
        {items.map((item) => (
          <div key={item.id || item.item_name} style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
            <Checkbox checked={item.passed} onChange={(e) => setItems((prev) => prev.map((x) => (x.id === item.id && x.item_name === item.item_name ? { ...x, passed: e.target.checked } : x)))}>
              {item.item_name}
            </Checkbox>
            <Input placeholder="备注" value={item.remark} onChange={(e) => setItems((prev) => prev.map((x) => (x.id === item.id && x.item_name === item.item_name ? { ...x, remark: e.target.value } : x)))} />
          </div>
        ))}
      </Modal>
    </Card>
  )
}
