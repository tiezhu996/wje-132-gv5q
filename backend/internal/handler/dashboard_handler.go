package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘 HTTP 处理器。
type DashboardHandler struct {
	svc    *service.DashboardService
	logger *slog.Logger
}

// NewDashboardHandler 构造仪表盘处理器。
func NewDashboardHandler(svc *service.DashboardService, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{svc: svc, logger: logger}
}

// Stats 仪表盘统计数据。
func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		h.wrapError(c, err, "Dashboard stats failed")
		return
	}
	OK(c, stats)
}

func (h *DashboardHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		h.logger.Warn("dashboard handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("dashboard handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
