package router

import (
	"log/slog"

	"safetyplatform/internal/config"
	"safetyplatform/internal/handler"
	"safetyplatform/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Router 路由装配器。
type Router struct {
	cfg     *config.Config
	db      *gorm.DB
	logger  *slog.Logger
	limiter *middleware.RateLimiter

	user       *handler.UserHandler
	incident   *handler.SafetyIncidentHandler
	inspection *handler.SafetyInspectionHandler
	item       *handler.InspectionItemHandler
	training   *handler.SafetyTrainingHandler
	cert       *handler.WorkerCertificationHandler
	dashboard  *handler.DashboardHandler
	upload     *handler.UploadHandler
	auditLog   *handler.AuditLogHandler
}

// New 构造路由装配器。
func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger,
	user *handler.UserHandler, incident *handler.SafetyIncidentHandler,
	inspection *handler.SafetyInspectionHandler, item *handler.InspectionItemHandler,
	training *handler.SafetyTrainingHandler, cert *handler.WorkerCertificationHandler,
	dashboard *handler.DashboardHandler, upload *handler.UploadHandler, auditLog *handler.AuditLogHandler) *Router {
	return &Router{
		cfg: cfg, db: db, logger: logger,
		limiter: middleware.NewRateLimiter(cfg.RateLimitPerMinute),
		user:    user, incident: incident, inspection: inspection, item: item,
		training: training, cert: cert, dashboard: dashboard, upload: upload, auditLog: auditLog,
	}
}

// Setup 装配全部路由并返回引擎。
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.RequestLogger(r.logger))
	engine.Use(middleware.ErrorHandler(r.logger))
	engine.Use(middleware.CORS(r.cfg))
	engine.Use(middleware.JWTConfig(r.cfg))
	engine.Use(middleware.AuditLog(r.db, r.logger))

	engine.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	engine.GET("/api/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	engine.GET("/api/v1/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	engine.Static("/uploads", r.cfg.UploadDir)

	v1 := engine.Group("/api/v1")
	r.registerAuthRoutes(v1)
	r.registerUserRoutes(v1)
	r.registerIncidentRoutes(v1)
	r.registerInspectionRoutes(v1)
	r.registerItemRoutes(v1)
	r.registerTrainingRoutes(v1)
	r.registerCertRoutes(v1)
	r.registerDashboardRoutes(v1)
	r.registerUploadRoutes(v1)
	r.registerAuditLogRoutes(v1)
	return engine
}

// registerAuthRoutes 登录注册（限流）。
func (r *Router) registerAuthRoutes(g *gin.RouterGroup) {
	auth := g.Group("/auth")
	auth.Use(r.limiter.Limit())
	auth.POST("/register", r.user.Register)
	auth.POST("/login", r.user.Login)
}

// registerUploadRoutes 图片上传。
func (r *Router) registerUploadRoutes(g *gin.RouterGroup) {
	upload := g.Group("/upload")
	upload.Use(middleware.AuthRequired(r.cfg))
	upload.POST("/image", r.upload.UploadImage)
}
