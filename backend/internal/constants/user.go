package constants

// UserRole 用户角色枚举。
const (
	RoleAdmin         = "admin"
	RoleSafetyManager = "safety_manager"
	RoleInspector     = "inspector"
	RoleWorker        = "worker"
)

// UserRoleValues 全部角色值。
var UserRoleValues = []string{RoleAdmin, RoleSafetyManager, RoleInspector, RoleWorker}
