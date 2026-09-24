package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

// InspectionItemHandler 检查项 HTTP 处理器。
type InspectionItemHandler struct {
	svc    *service.SafetyInspectionService
	logger *slog.Logger
}

// NewInspectionItemHandler 构造检查项处理器。
func NewInspectionItemHandler(svc *service.SafetyInspectionService, logger *slog.Logger) *InspectionItemHandler {
	return &InspectionItemHandler{svc: svc, logger: logger}
}

// ListByInspection 查询某检查的检查项。
func (h *InspectionItemHandler) ListByInspection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "InspectionItem list: invalid inspection id")
		return
	}
	_, items, err := h.svc.GetWithItems(id)
	if err != nil {
		h.wrapError(c, err, "InspectionItem list failed")
		return
	}
	OK(c, items)
}

func (h *InspectionItemHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		h.logger.Warn("inspection item handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("inspection item handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
