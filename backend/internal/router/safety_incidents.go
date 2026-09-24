package router

import (
	"safetyplatform/internal/constants"
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerIncidentRoutes 安全事件路由。
func (r *Router) registerIncidentRoutes(g *gin.RouterGroup) {
	incidents := g.Group("/incidents")
	incidents.Use(middleware.AuthRequired(r.cfg))
	incidents.GET("", r.incident.List)
	incidents.GET("/overdue-acceptance", r.incident.OverdueAcceptance)
	incidents.GET("/:id", r.incident.Get)
	incidents.POST("", r.incident.Report)
	incidents.POST("/:id/assign", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.incident.Assign)
	incidents.POST("/:id/rectify", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.incident.Rectify)
	incidents.POST("/:id/supervise", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.incident.Supervise)
	incidents.POST("/:id/close", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.incident.Close)
}
