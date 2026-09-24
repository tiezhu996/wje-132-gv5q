import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic, Table } from 'antd'
import ReactECharts from 'echarts-for-react'
import { getDashboardStats } from '@/api/dashboard'
import RiskLevelTag from '@/components/common/RiskLevelTag'
import StatusBadge from '@/components/common/StatusBadge'
import { SeverityText } from '@/constants/incident'
import { formatDate } from '@/utils/dateFormat'
import type { SafetyIncident } from '@/types'

interface Stats {
  trend: { day: string; cnt: number }[]
  severity_distribution: { severity_level: string; cnt: number }[]
  pending_rectification: SafetyIncident[]
  inspection: { total: number; completed_rate: number }
  training_completed_rate: number
  expiring_certs: unknown[]
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats | null>(null)

  useEffect(() => {
    getDashboardStats().then((res: any) => setStats(res.data))
  }, [])

  const trendOption = {
    title: { text: '近 30 天事件趋势' },
    xAxis: { type: 'category', data: stats?.trend.map((t) => t.day) || [] },
    yAxis: { type: 'value' },
    series: [{ name: '事件数', type: 'line', smooth: true, data: stats?.trend.map((t) => t.cnt) || [] }],
  }

  const pieOption = {
    title: { text: '风险等级分布' },
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: '60%',
      data: stats?.severity_distribution.map((d) => ({ name: SeverityText[d.severity_level] || d.severity_level, value: d.cnt })) || [],
    }],
  }

  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}><Card><Statistic title="检查计划总数" value={stats?.inspection.total || 0} /></Card></Col>
        <Col span={6}><Card><Statistic title="检查完成率" value={stats?.inspection.completed_rate || 0} suffix="%" /></Card></Col>
        <Col span={6}><Card><Statistic title="本月培训完成率" value={stats?.training_completed_rate || 0} suffix="%" /></Card></Col>
        <Col span={6}><Card><Statistic title="即将过期资质" value={stats?.expiring_certs?.length || 0} /></Card></Col>
      </Row>
      <Row gutter={16}>
        <Col span={14}><Card><ReactECharts option={trendOption} style={{ height: 320 }} /></Card></Col>
        <Col span={10}><Card><ReactECharts option={pieOption} style={{ height: 320 }} /></Card></Col>
      </Row>
      <Card title="待整改事件" style={{ marginTop: 16 }}>
        <Table<SafetyIncident>
          rowKey="id"
          dataSource={stats?.pending_rectification || []}
          pagination={false}
          columns={[
            { title: '标题', dataIndex: 'title' },
            { title: '区域', dataIndex: 'area' },
            { title: '风险等级', dataIndex: 'severity_level', render: (v) => <RiskLevelTag level={v} /> },
            { title: '状态', dataIndex: 'status', render: (v) => <StatusBadge status={v} /> },
            { title: '整改期限', dataIndex: 'rectification_deadline', render: (v) => formatDate(v) },
          ]}
        />
      </Card>
    </div>
  )
}
