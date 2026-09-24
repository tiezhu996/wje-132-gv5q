package router

import (
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerDashboardRoutes 仪表盘路由。
func (r *Router) registerDashboardRoutes(g *gin.RouterGroup) {
	dash := g.Group("/dashboard")
	dash.Use(middleware.AuthRequired(r.cfg))
	dash.GET("/stats", r.dashboard.Stats)
}
