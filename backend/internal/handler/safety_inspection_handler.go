package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/dto"
	"safetyplatform/internal/model"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

// SafetyInspectionHandler 安全检查 HTTP 处理器。
type SafetyInspectionHandler struct {
	svc    *service.SafetyInspectionService
	logger *slog.Logger
}

// NewSafetyInspectionHandler 构造安全检查处理器。
func NewSafetyInspectionHandler(svc *service.SafetyInspectionService, logger *slog.Logger) *SafetyInspectionHandler {
	return &SafetyInspectionHandler{svc: svc, logger: logger}
}

// List 检查计划列表。
func (h *SafetyInspectionHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.svc.List(q.Page, q.PageSize, c.Query("status"))
	if err != nil {
		h.wrapError(c, err, "SafetyInspection list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// Get 检查详情。
func (h *SafetyInspectionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyInspection[id] get: invalid id")
		return
	}
	ins, items, err := h.svc.GetWithItems(id)
	if err != nil {
		h.wrapError(c, err, "SafetyInspection get failed")
		return
	}
	OK(c, gin.H{"inspection": ins, "items": items})
}

// Create 创建检查计划。
func (h *SafetyInspectionHandler) Create(c *gin.Context) {
	var req dto.InspectionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyInspection create: "+err.Error())
		return
	}
	items := make([]model.InspectionItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, model.InspectionItem{ItemName: it.ItemName, Passed: it.Passed, Remark: it.Remark, PhotoURL: it.PhotoURL})
	}
	ins, err := h.svc.Create(req.Name, req.InspectionType, req.Area, req.InspectionDate, req.InspectorID, items)
	if err != nil {
		h.wrapError(c, err, "SafetyInspection[name="+req.Name+"] create failed")
		return
	}
	OK(c, ins)
}

// Execute 执行检查。
func (h *SafetyInspectionHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyInspection[id] execute: invalid id")
		return
	}
	var req dto.InspectionExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyInspection[id="+strconv.FormatUint(id, 10)+"] execute: "+err.Error())
		return
	}
	items := make([]model.InspectionItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, model.InspectionItem{ID: it.ID, ItemName: it.ItemName, Passed: it.Passed, Remark: it.Remark, PhotoURL: it.PhotoURL})
	}
	ins, err := h.svc.Execute(id, items)
	if err != nil {
		h.wrapError(c, err, "SafetyInspection execute failed")
		return
	}
	OKWithMessage(c, constants.MsgInspectionCompleted, ins)
}

// Report 检查报告。
func (h *SafetyInspectionHandler) Report(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "SafetyInspection[id] report: invalid id")
		return
	}
	ins, items, err := h.svc.Report(id)
	if err != nil {
		h.wrapError(c, err, "SafetyInspection report failed")
		return
	}
	OK(c, gin.H{"inspection": ins, "items": items})
}

func (h *SafetyInspectionHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("inspection handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("inspection handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
