import { useEffect, useState } from 'react'
import { Alert, Button, Card, DatePicker, Form, Input, Modal, Select, Space, Table, Tabs, Tag, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { reportIncident } from '@/api/incident'
import { useIncidentStore } from '@/stores/incidentStore'
import { useIncident } from '@/hooks/useIncident'
import RiskLevelTag from '@/components/common/RiskLevelTag'
import StatusBadge from '@/components/common/StatusBadge'
import RoleGuard from '@/components/common/RoleGuard'
import { IncidentCategories, IncidentStatusOptions, SeverityOptions } from '@/constants/incident'
import { formatDateTime, formatOverdueDuration } from '@/utils/dateFormat'
import type { OverdueIncident, SafetyIncident } from '@/types'

type DetailSource = 'list' | 'overdue'

export default function IncidentManage() {
  const store = useIncidentStore()
  const incident = useIncident()
  const [activeTab, setActiveTab] = useState('list')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [filters, setFilters] = useState<Record<string, unknown>>({})
  const [open, setOpen] = useState(false)
  const [detailId, setDetailId] = useState<number>()
  const [detailSource, setDetailSource] = useState<DetailSource>('list')
  const [superviseTarget, setSuperviseTarget] = useState<OverdueIncident>()
  const [sending, setSending] = useState(false)
  const [superviseForm] = Form.useForm()
  const [form] = Form.useForm()

  useEffect(() => {
    store.fetchList({ page, page_size: pageSize, ...filters })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, filters])

  useEffect(() => {
    if (activeTab === 'overdue') store.fetchOverdue()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab])

  useEffect(() => {
    if (detailId) incident.load(detailId)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [detailId])

  // 详情弹窗内执行动作后，按入口刷新对应列表。
  async function refreshSource() {
    if (detailSource === 'overdue') await store.fetchOverdue()
    else await store.fetchList({ page, page_size: pageSize, ...filters })
  }

  function openDetail(id: number, source: DetailSource) {
    setDetailSource(source)
    setDetailId(id)
  }

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

  async function onSupervise() {
    if (!superviseTarget) return
    const values = await superviseForm.validateFields()
    setSending(true)
    try {
      await incident.supervise(superviseTarget.id, values.note)
      await store.fetchOverdue()
      setSuperviseTarget(undefined)
      superviseForm.resetFields()
    } finally {
      setSending(false)
    }
  }

  const overdueColumns = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '标题', dataIndex: 'title' },
    { title: '区域', dataIndex: 'area', width: 120 },
    { title: '风险等级', dataIndex: 'severity_level', width: 100, render: (v: string) => <RiskLevelTag level={v} /> },
    { title: '整改期限', dataIndex: 'rectification_deadline', width: 150, render: (v: string) => formatDateTime(v) },
    {
      title: '逾期时长',
      dataIndex: 'overdue_duration',
      width: 160,
      render: (v: number) => <Tag color="red">{formatOverdueDuration(v)}</Tag>,
    },
    {
      title: '督办说明',
      render: (_: unknown, row: OverdueIncident) =>
        row.supervision_at ? (
          <div>
            <div>{row.supervision_note}</div>
            <small style={{ color: '#999' }}>{formatDateTime(row.supervision_at)}</small>
          </div>
        ) : (
          <span style={{ color: '#999' }}>未督办</span>
        ),
    },
    {
      title: '操作',
      width: 150,
      render: (_: unknown, row: OverdueIncident) => (
        <Space>
          <a onClick={() => openDetail(row.id, 'overdue')}>详情</a>
          <RoleGuard roles={['admin', 'safety_manager']}>
            {!row.supervision_at && <a onClick={() => setSuperviseTarget(row)}>督办</a>}
          </RoleGuard>
        </Space>
      ),
    },
  ]

  return (
    <Card>
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'list',
            label: '事件列表',
            children: (
              <>
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
                      render: (_, row) => <a onClick={() => openDetail(row.id, 'list')}>详情</a>,
                    },
                  ]}
                />
              </>
            ),
          },
          {
            key: 'overdue',
            label: `逾期验收队列${store.overdueList.length ? ` (${store.overdueList.length})` : ''}`,
            children: (
              <>
                <Alert
                  type="warning"
                  showIcon
                  style={{ marginBottom: 16 }}
                  message="以下事件已整改但超过整改期限仍未验收关闭，请尽快督办并落实验收。"
                />
                <Table<OverdueIncident>
                  rowKey="id"
                  dataSource={store.overdueList}
                  pagination={false}
                  columns={overdueColumns}
                />
              </>
            ),
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
      <Modal
        title="发起督办"
        open={!!superviseTarget}
        onOk={onSupervise}
        confirmLoading={sending}
        okText="发出督办"
        onCancel={() => setSuperviseTarget(undefined)}
      >
        <p style={{ color: '#666' }}>
          {superviseTarget?.title}（整改期限：{superviseTarget ? formatDateTime(superviseTarget.rectification_deadline) : '-'}，
          {superviseTarget ? formatOverdueDuration(superviseTarget.overdue_duration) : ''}）
        </p>
        <Form form={superviseForm} layout="vertical">
          <Form.Item name="note" label="督办说明" rules={[{ required: true, message: '请填写督办说明' }, { max: 500 }]}>
            <Input.TextArea rows={4} placeholder="请说明督办要求与验收期限" maxLength={500} showCount />
          </Form.Item>
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
              <p>整改期限：{formatDateTime(cur.rectification_deadline)}</p>
              {cur.supervision_at && (
                <Alert
                  type="warning"
                  style={{ marginBottom: 12 }}
                  message={`督办说明：${cur.supervision_note}`}
                  description={`督办时间：${formatDateTime(cur.supervision_at)}`}
                />
              )}
              <RoleGuard roles={['admin', 'safety_manager']}>
                <Space>
                  {cur.status === 'reported' && <Button type="primary" onClick={async () => { await incident.assign(cur.id); await refreshSource() }}>指派调查</Button>}
                  {cur.status === 'investigating' && <Button onClick={async () => { await incident.rectify(cur.id, '已完成整改，验收合格'); await refreshSource() }}>提交整改</Button>}
                  {cur.status === 'resolved' && <Button danger onClick={async () => { await incident.close(cur.id); await refreshSource() }}>关闭事件</Button>}
                </Space>
              </RoleGuard>
            </div>
          )
        })()}
      </Modal>
    </Card>
  )
}
