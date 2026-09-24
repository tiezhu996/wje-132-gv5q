import { useEffect, useState } from 'react'
import { Button, Card, Modal, Select, Space, Table, message } from 'antd'
import { listCertifications, reviewCertification } from '@/api/certification'
import StatusBadge from '@/components/common/StatusBadge'
import UserAvatar from '@/components/common/UserAvatar'
import RoleGuard from '@/components/common/RoleGuard'
import { formatDate } from '@/utils/dateFormat'
import type { WorkerCertification } from '@/types'

export default function CertReview() {
  const [list, setList] = useState<WorkerCertification[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [statusFilter, setStatusFilter] = useState('')
  const [expiring, setExpiring] = useState(false)
  const [detail, setDetail] = useState<WorkerCertification | null>(null)

  async function load() {
    const res: any = await listCertifications({ page, page_size: pageSize, status: statusFilter || undefined, expiring })
    setList(res.data.list)
    setTotal(res.data.total)
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize, statusFilter, expiring])

  async function review(id: number, status: string) {
    await reviewCertification(id, status)
    message.success('审核完成')
    load()
  }

  return (
    <Card>
      <Space style={{ marginBottom: 16 }}>
        <Select placeholder="状态" allowClear style={{ width: 130 }} options={[
          { label: '待审核', value: 'pending' },
          { label: '已通过', value: 'approved' },
          { label: '已拒绝', value: 'rejected' },
          { label: '已过期', value: 'expired' },
        ]} onChange={(v) => { setStatusFilter(v || ''); setPage(1) }} />
        <Button type={expiring ? 'primary' : 'default'} onClick={() => setExpiring((v) => !v)}>过期预警</Button>
      </Space>
      <Table<WorkerCertification>
        rowKey="id"
        dataSource={list}
        pagination={{ current: page, pageSize, total, onChange: (p, ps) => { setPage(p); setPageSize(ps) } }}
        columns={[
          { title: 'ID', dataIndex: 'id', width: 70 },
          { title: '用户ID', dataIndex: 'user_id', width: 90 },
          { title: '证书类型', dataIndex: 'cert_type' },
          { title: '证书编号', dataIndex: 'cert_no' },
          { title: '发证机构', dataIndex: 'issue_org' },
          { title: '有效期至', dataIndex: 'valid_until', render: (v) => formatDate(v) },
          { title: '状态', dataIndex: 'status', render: (v) => <StatusBadge status={v} /> },
          {
            title: '操作',
            render: (_, row) => (
              <>
                <a onClick={() => setDetail(row)}>详情</a>
                <RoleGuard roles={['admin', 'safety_manager']}>
                  {row.status === 'pending' && (
                    <>
                      <Button size="small" type="primary" style={{ marginLeft: 8 }} onClick={() => review(row.id, 'approved')}>通过</Button>
                      <Button size="small" danger style={{ marginLeft: 8 }} onClick={() => review(row.id, 'rejected')}>拒绝</Button>
                    </>
                  )}
                </RoleGuard>
              </>
            ),
          },
        ]}
      />
      <Modal title="资质详情" open={!!detail} onCancel={() => setDetail(null)} footer={null}>
        {detail && (
          <div>
            <p><UserAvatar name={`用户 #${detail.user_id}`} /></p>
            <p>证书类型：{detail.cert_type}</p>
            <p>证书编号：{detail.cert_no}</p>
            <p>发证机构：{detail.issue_org}</p>
            <p>发证日期：{formatDate(detail.issue_date)}</p>
            <p>有效期至：{formatDate(detail.valid_until)}</p>
            {detail.cert_photo_url && <p><img src={detail.cert_photo_url} alt="证书照片" style={{ maxWidth: 300 }} /></p>}
          </div>
        )}
      </Modal>
    </Card>
  )
}
