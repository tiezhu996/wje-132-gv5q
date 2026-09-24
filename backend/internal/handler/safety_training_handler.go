package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/dto"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

// SafetyTrainingHandler 安全培训 HTTP 处理器。
type SafetyTrainingHandler struct {
	svc    *service.SafetyTrainingService
	logger *slog.Logger
}

// NewSafetyTrainingHandler 构造安全培训处理器。
func NewSafetyTrainingHandler(svc *service.SafetyTrainingService, logger *slog.Logger) *SafetyTrainingHandler {
	return &SafetyTrainingHandler{svc: svc, logger: logger}
}

// List 培训列表。
func (h *SafetyTrainingHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.svc.List(q.Page, q.PageSize, c.Query("training_type"))
	if err != nil {
		h.wrapError(c, err, "SafetyTraining list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// Get 培训详情。
func (h *SafetyTrainingHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyTraining[id] get: invalid id")
		return
	}
	t, err := h.svc.Get(id)
	if err != nil {
		h.wrapError(c, err, "SafetyTraining get failed")
		return
	}
	OK(c, t)
}

// Create 创建培训。
func (h *SafetyTrainingHandler) Create(c *gin.Context) {
	var req dto.TrainingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyTraining create: "+err.Error())
		return
	}
	t, err := h.svc.Create(req.Topic, req.TrainingType, req.TrainingDate, req.DurationHours,
		req.Trainer, req.Location, req.ContentSummary, req.ParticipantIDs, req.AssessmentMethod)
	if err != nil {
		h.wrapError(c, err, "SafetyTraining[topic="+req.Topic+"] create failed")
		return
	}
	OK(c, t)
}

// Record 记录成绩。
func (h *SafetyTrainingHandler) Record(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyTraining[id] record: invalid id")
		return
	}
	var req dto.TrainingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyTraining[id="+strconv.FormatUint(id, 10)+"] record: "+err.Error())
		return
	}
	t, err := h.svc.Record(id, req.ParticipantIDs, req.PassRate)
	if err != nil {
		h.wrapError(c, err, "SafetyTraining record failed")
		return
	}
	OK(c, t)
}

func (h *SafetyTrainingHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("training handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("training handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
