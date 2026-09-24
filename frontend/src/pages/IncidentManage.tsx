import { useEffect, useState } from 'react'
import { Alert, Button, Card, DatePicker, Form, Input, Modal, Select, Space, Table, Tabs, Tag, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { reportIncident } from '@/api/incident'
import { useIncidentStore } from '@/stores/incidentStore'
import { useIncident } from '@/hooks/useIncident'
import RiskLevelTag from '@/components/common/RiskLevelTag'
import StatusBadge from '@/components/common/StatusBadge'
import RoleGuard from '@/components/common/RoleGuard'
import EmptyState from '@/components/common/EmptyState'
import { IncidentCategories, IncidentStatusOptions, SeverityOptions } from '@/constants/incident'
import { formatDateTime } from '@/utils/dateFormat'
import { formatOverdueDuration } from '@/utils/overdueDuration'
import type { OverdueAcceptanceItem, SafetyIncident } from '@/types'

export default function IncidentManage() {
  const store = useIncidentStore()
  const incident = useIncident()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [filters, setFilters] = useState<Record<string, unknown>>({})
  const [open, setOpen] = useState(false)
  const [detailId, setDetailId] = useState<number>()
  const [tab, setTab] = useState('list')
  const [superviseTarget, setSuperviseTarget] = useState<OverdueAcceptanceItem>()
  const [superviseLoading, setSuperviseLoading] = useState(false)
  const [superviseForm] = Form.useForm()
  const [form] = Form.useForm()

  useEffect(() => {
    store.fetchList({ page, page_size: pageSize, ...filters })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, filters])

  useEffect(() => {
    if (tab === 'overdue') store.fetchOverdue()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab])

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

  async function onSupervise() {
    if (!superviseTarget) return
    const values = await superviseForm.validateFields()
    setSuperviseLoading(true)
    try {
      await incident.supervise(superviseTarget.id, values.note)
      setSuperviseTarget(undefined)
      superviseForm.resetFields()
      await store.fetchOverdue()
    } catch {
      // 重复督办等冲突提示已由 request 拦截器统一弹出
    } finally {
      setSuperviseLoading(false)
    }
  }

  async function onClose(id: number) {
    await incident.close(id)
    await Promise.all([
      store.fetchList({ page, page_size: pageSize, ...filters }),
      tab === 'overdue' ? store.fetchOverdue() : Promise.resolve(),
    ])
  }

  const detailModal = (
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
                {cur.status === 'resolved' && <Button danger onClick={() => onClose(cur.id)}>关闭事件</Button>}
              </Space>
            </RoleGuard>
          </div>
        )
      })()}
    </Modal>
  )

  return (
    <Card>
      <Tabs
        activeKey={tab}
        onChange={setTab}
        items={[
          {
            key: 'list',
            label: '事件列表',
            children: (
              <>
                <Space style={{ marginBottom: 16 }}>
                  <Select placeholder="严重等级" allowClear style={{ width: 130 }} options={SeverityOptions} onChange={(v) => { setFilters((f) => ({ ...f, severity: v })); setPage(1) }} />
                  <Select placeholder="状态" allowClear style={{ width: 130 }} options={IncidentStatusOptions} onChange={(v) => { setFilters((f) => ({ ...f, status: v })); setPage(1) }} />
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
              </>
            ),
          },
          {
            key: 'overdue',
            label: <span>逾期验收{store.overdueList.length > 0 && <Tag color="red" style={{ marginInlineStart: 6 }}>{store.overdueList.length}</Tag>}</span>,
            children: (
              <>
                <Alert
                  type="warning"
                  showIcon
                  style={{ marginBottom: 16 }}
                  message="逾期验收队列仅列入已整改且整改期限已过的事件，按逾期时长与风险等级排序；每条记录仅可发起一次督办。"
                />
                <Table<OverdueAcceptanceItem>
                  rowKey="id"
                  dataSource={store.overdueList}
                  pagination={false}
                  locale={{ emptyText: <EmptyState text="暂无逾期未验收的整改记录" /> }}
                  columns={[
                    { title: 'ID', dataIndex: 'id', width: 70 },
                    { title: '标题', dataIndex: 'title', render: (v, row) => <a onClick={() => setDetailId(row.id)}>{v}</a> },
                    { title: '区域', dataIndex: 'area' },
                    { title: '风险等级', dataIndex: 'severity_level', render: (v) => <RiskLevelTag level={v} /> },
                    { title: '整改期限', dataIndex: 'rectification_deadline', render: (v) => formatDateTime(v) },
                    {
                      title: '逾期时长',
                      dataIndex: 'overdue_seconds',
                      width: 120,
                      render: (v: number) => <Tag color="red">{formatOverdueDuration(v)}</Tag>,
                      sorter: (a, b) => a.overdue_seconds - b.overdue_seconds,
                      defaultSortOrder: 'descend',
                    },
                    {
                      title: '督办说明',
                      dataIndex: 'supervision',
                      render: (v: OverdueAcceptanceItem['supervision']) => v
                        ? (
                          <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                            <span>{v.note}</span>
                            <span style={{ color: '#999', fontSize: 12 }}>
                              {v.operator_name ? `${v.operator_name} · ` : ''}{formatDateTime(v.created_at)}
                            </span>
                          </div>
                        )
                        : <Tag>未督办</Tag>,
                    },
                    {
                      title: '操作',
                      width: 120,
                      render: (_, row) => (
                        <RoleGuard roles={['admin', 'safety_manager']}>
                          {row.supervision
                            ? <span style={{ color: '#999' }}>已督办</span>
                            : <a onClick={() => { setSuperviseTarget(row); superviseForm.resetFields() }}>发起督办</a>}
                        </RoleGuard>
                      ),
                    },
                  ]}
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
        title={superviseTarget ? `发起督办：${superviseTarget.title}` : '发起督办'}
        open={!!superviseTarget}
        onOk={onSupervise}
        confirmLoading={superviseLoading}
        onCancel={() => setSuperviseTarget(undefined)}
        okText="提交督办"
      >
        <Form form={superviseForm} layout="vertical">
          <Form.Item name="note" label="督办说明" rules={[{ required: true, message: '请填写督办说明' }, { max: 500, message: '督办说明不超过 500 字' }]}>
            <Input.TextArea rows={4} placeholder="请填写督办要求与说明，提交后将通知相关责任人尽快验收" maxLength={500} showCount />
          </Form.Item>
        </Form>
      </Modal>
      {detailModal}
    </Card>
  )
}
