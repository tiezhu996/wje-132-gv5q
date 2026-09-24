package router

import (
	"safetyplatform/internal/constants"
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerTrainingRoutes 安全培训路由。
func (r *Router) registerTrainingRoutes(g *gin.RouterGroup) {
	trainings := g.Group("/trainings")
	trainings.Use(middleware.AuthRequired(r.cfg))
	trainings.GET("", r.training.List)
	trainings.GET("/:id", r.training.Get)
	trainings.POST("", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.training.Create)
	trainings.POST("/:id/record", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.training.Record)
}
