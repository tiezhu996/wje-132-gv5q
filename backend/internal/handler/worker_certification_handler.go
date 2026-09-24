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

// WorkerCertificationHandler 人员资质 HTTP 处理器。
type WorkerCertificationHandler struct {
	svc    *service.WorkerCertificationService
	logger *slog.Logger
}

// NewWorkerCertificationHandler 构造人员资质处理器。
func NewWorkerCertificationHandler(svc *service.WorkerCertificationService, logger *slog.Logger) *WorkerCertificationHandler {
	return &WorkerCertificationHandler{svc: svc, logger: logger}
}

// List 资质列表。
func (h *WorkerCertificationHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	expiring := c.Query("expiring") == "true"
	list, total, err := h.svc.List(q.Page, q.PageSize, c.Query("status"), expiring)
	if err != nil {
		h.wrapError(c, err, "WorkerCertification list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// ListByUser 某用户资质。
func (h *WorkerCertificationHandler) ListByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil || userID == 0 {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "WorkerCertification list by user: invalid user_id")
		return
	}
	list, err := h.svc.ListByUser(userID)
	if err != nil {
		h.wrapError(c, err, "WorkerCertification list by user failed")
		return
	}
	OK(c, list)
}

// Submit 提交资质。
func (h *WorkerCertificationHandler) Submit(c *gin.Context) {
	var req dto.CertSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "WorkerCertification submit: "+err.Error())
		return
	}
	cert, err := h.svc.Submit(req.UserID, req.CertType, req.CertNo, req.IssueOrg, req.IssueDate, req.ValidUntil, req.CertPhotoURL)
	if err != nil {
		h.wrapError(c, err, "WorkerCertification submit failed")
		return
	}
	OKWithMessage(c, constants.MsgCertSubmitted, cert)
}

// Review 审核资质。
func (h *WorkerCertificationHandler) Review(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "WorkerCertification[id] review: invalid id")
		return
	}
	var req dto.CertReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "WorkerCertification[id="+strconv.FormatUint(id, 10)+"] review: "+err.Error())
		return
	}
	cert, err := h.svc.Review(id, req.Status)
	if err != nil {
		h.wrapError(c, err, "WorkerCertification review failed")
		return
	}
	OKWithMessage(c, constants.MsgCertReviewed, cert)
}

func (h *WorkerCertificationHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("cert handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("cert handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
