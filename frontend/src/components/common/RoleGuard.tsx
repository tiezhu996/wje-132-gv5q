import type { ReactNode } from 'react'
import { useAuthStore } from '@/stores/authStore'

export default function RoleGuard({ roles, children }: { roles: string[]; children: ReactNode }) {
  const role = useAuthStore((s) => s.user?.role || '')
  if (!roles.includes(role)) return null
  return <>{children}</>
}
