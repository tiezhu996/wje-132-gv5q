import { useEffect, useState } from 'react'
import { Button, Card, DatePicker, Form, Input, Modal, Select, Space, Table, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { reportIncident } from '@/api/incident'
import { useIncidentStore } from '@/stores/incidentStore'
import { useIncident } from '@/hooks/useIncident'
import RiskLevelTag from '@/components/common/RiskLevelTag'
import StatusBadge from '@/components/common/StatusBadge'
import RoleGuard from '@/components/common/RoleGuard'
import { IncidentCategories, IncidentStatusOptions, SeverityOptions } from '@/constants/incident'
import { formatDateTime } from '@/utils/dateFormat'
import type { SafetyIncident } from '@/types'

export default function IncidentManage() {
  const store = useIncidentStore()
  const incident = useIncident()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [filters, setFilters] = useState<Record<string, unknown>>({})
  const [open, setOpen] = useState(false)
  const [detailId, setDetailId] = useState<number>()
  const [form] = Form.useForm()

  useEffect(() => {
    store.fetchList({ page, page_size: pageSize, ...filters })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, filters])

  useEffect(() => {
    if (detailId) incident.load(detailId)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [detailId])

  async function onReport() {
    const values = await form.validateFields()
    await reportIncident({
      title: values.title,
      description: values.description,
      occurred_at: values.occurred_at ? values.occurred_at.format('YYYY-MM-DDTHH:mm:ss') : undefined,
      site_id: values.site_id,
      area: values.area,
      severity_level: values.severity_level,
      category: values.category,
    })
    message.success('事件上报成功')
    setOpen(false)
    form.resetFields()
    setPage(1)
  }

  return (
    <Card>
      <Space style={{ marginBottom: 16 }}>
        <Select placeholder="严重等级" allowClear style={{ width: 130 }} options={SeverityOptions} onChange={(v) => { setFilters({ severity: v }); setPage(1) }} />
        <Select placeholder="状态" allowClear style={{ width: 130 }} options={IncidentStatusOptions} onChange={(v) => { setFilters({ status: v }); setPage(1) }} />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>上报事件</Button>
      </Space>
      <Table<SafetyIncident>
        rowKey="id"
        dataSource={store.list}
        pagination={{ current: page, pageSize, total: store.total, onChange: (p, ps) => { setPage(p); setPageSize(ps) } }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 70 },
          { title: '标题', dataIndex: 'title' },
          { title: '区域', dataIndex: 'area' },
          { title: '风险等级', dataIndex: 'severity_level', render: (v) => <RiskLevelTag level={v} /> },
          { title: '分类', dataIndex: 'category' },
          { title: '状态', dataIndex: 'status', render: (v) => <StatusBadge status={v} /> },
          { title: '发生时间', dataIndex: 'occurred_at', render: (v) => formatDateTime(v) },
          {
            title: '操作',
            render: (_, row) => <a onClick={() => setDetailId(row.id)}>详情</a>,
          },
        ]}
      />
      <Modal title="上报安全事件" open={open} onOk={onReport} onCancel={() => setOpen(false)} width={560}>
        <Form form={form} layout="vertical">
          <Form.Item name="title" label="事件标题" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="description" label="事件描述"><Input.TextArea rows={3} /></Form.Item>
          <Form.Item name="occurred_at" label="发生时间" rules={[{ required: true }]}><DatePicker showTime style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="site_id" label="工地编号"><Input /></Form.Item>
          <Form.Item name="area" label="区域"><Input /></Form.Item>
          <Form.Item name="severity_level" label="严重等级" rules={[{ required: true }]}><Select options={SeverityOptions} /></Form.Item>
          <Form.Item name="category" label="分类" rules={[{ required: true }]}><Select options={IncidentCategories.map((c) => ({ label: c, value: c }))} /></Form.Item>
        </Form>
      </Modal>
      <Modal title="事件详情" open={!!detailId} onCancel={() => setDetailId(undefined)} footer={null} width={640}>
        {incident.incident && (() => {
          const cur = incident.incident
          return (
            <div>
              <p><b>{cur.title}</b> <RiskLevelTag level={cur.severity_level} /> <StatusBadge status={cur.status} /></p>
              <p>{cur.description}</p>
              <p>区域：{cur.area} / 分类：{cur.category}</p>
              <p>整改措施：{cur.rectification_measures || '-'}</p>
              <RoleGuard roles={['admin', 'safety_manager']}>
                <Space>
                  {cur.status === 'reported' && <Button type="primary" onClick={() => incident.assign(cur.id)}>指派调查</Button>}
                  {cur.status === 'investigating' && <Button onClick={() => incident.rectify(cur.id, '已完成整改，验收合格')}>提交整改</Button>}
                  {cur.status === 'resolved' && <Button danger onClick={() => incident.close(cur.id)}>关闭事件</Button>}
                </Space>
              </RoleGuard>
            </div>
          )
        })()}
      </Modal>
    </Card>
  )
}
