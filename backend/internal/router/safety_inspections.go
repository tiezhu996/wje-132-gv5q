package router

import (
	"safetyplatform/internal/constants"
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerInspectionRoutes 安全检查路由。
func (r *Router) registerInspectionRoutes(g *gin.RouterGroup) {
	inspections := g.Group("/inspections")
	inspections.Use(middleware.AuthRequired(r.cfg))
	inspections.GET("", r.inspection.List)
	inspections.GET("/:id", r.inspection.Get)
	inspections.POST("", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager, constants.RoleInspector), r.inspection.Create)
	inspections.POST("/:id/execute", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager, constants.RoleInspector), r.inspection.Execute)
	inspections.GET("/:id/report", r.inspection.Report)
}
