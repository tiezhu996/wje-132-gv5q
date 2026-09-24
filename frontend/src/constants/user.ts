// 用户角色枚举（与后端 backend/internal/constants/user.go 保持一致）
export const UserRole = {
  ADMIN: 'admin',
  SAFETY_MANAGER: 'safety_manager',
  INSPECTOR: 'inspector',
  WORKER: 'worker',
} as const

export const UserRoleText: Record<string, string> = {
  [UserRole.ADMIN]: '管理员',
  [UserRole.SAFETY_MANAGER]: '安全管理员',
  [UserRole.INSPECTOR]: '监理',
  [UserRole.WORKER]: '工人',
}
