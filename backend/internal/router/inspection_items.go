package router

import (
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerItemRoutes 检查项路由。
func (r *Router) registerItemRoutes(g *gin.RouterGroup) {
	items := g.Group("/inspection-items")
	items.Use(middleware.AuthRequired(r.cfg))
	items.GET("/by-inspection/:id", r.item.ListByInspection)
}
