package router

import (
	"safetyplatform/internal/constants"
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerCertRoutes 人员资质路由。
func (r *Router) registerCertRoutes(g *gin.RouterGroup) {
	certs := g.Group("/certifications")
	certs.Use(middleware.AuthRequired(r.cfg))
	certs.GET("", r.cert.List)
	certs.GET("/by-user", r.cert.ListByUser)
	certs.POST("", r.cert.Submit)
	certs.POST("/:id/review", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), r.cert.Review)
}
