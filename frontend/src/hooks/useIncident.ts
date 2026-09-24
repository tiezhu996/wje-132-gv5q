import { useState } from 'react'
import { message } from 'antd'
import { assignIncident, rectifyIncident, closeIncident, getIncident } from '@/api/incident'
import type { SafetyIncident } from '@/types'

export function useIncident() {
  const [loading, setLoading] = useState(false)
  const [incident, setIncident] = useState<SafetyIncident | null>(null)

  async function load(id: number) {
    setLoading(true)
    try {
      const res: any = await getIncident(id)
      setIncident(res.data)
    } finally {
      setLoading(false)
    }
  }

  async function assign(id: number) {
    await assignIncident(id)
    message.success('已指派调查')
    await load(id)
  }

  async function rectify(id: number, measures: string, deadline?: string) {
    await rectifyIncident(id, { measures, deadline })
    message.success('整改已提交')
    await load(id)
  }

  async function close(id: number) {
    await closeIncident(id)
    message.success('事件已关闭')
    await load(id)
  }

  return { loading, incident, load, assign, rectify, close }
}
