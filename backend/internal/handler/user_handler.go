package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/dto"
	"safetyplatform/internal/middleware"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户 HTTP 处理器。
type UserHandler struct {
	svc    *service.UserService
	logger *slog.Logger
}

// NewUserHandler 构造用户处理器。
func NewUserHandler(svc *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

// Register 注册。
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "User register: "+err.Error())
		return
	}
	u, err := h.svc.Register(req.Phone, req.Password, req.Name, req.Role)
	if err != nil {
		h.wrapError(c, err, "User[phone="+req.Phone+"] register failed")
		return
	}
	OK(c, u)
}

// Login 登录。
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "User login: "+err.Error())
		return
	}
	token, u, err := h.svc.Login(c.GetString("jwt_secret"), c.GetInt("jwt_expire_hours"), req.Phone, req.Password)
	if err != nil {
		h.wrapError(c, err, "User[phone="+req.Phone+"] login failed")
		return
	}
	OKWithMessage(c, constants.MsgLoginSuccess, dto.LoginResponse{Token: token, User: u})
}

// Me 当前用户信息。
func (h *UserHandler) Me(c *gin.Context) {
	u, err := h.svc.GetByID(middleware.GetUserID(c))
	if err != nil {
		h.wrapError(c, err, "User me failed")
		return
	}
	OK(c, u)
}

// UpdateProfile 修改资料。
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "User profile update: "+err.Error())
		return
	}
	u, err := h.svc.UpdateProfile(middleware.GetUserID(c), req.Name, req.Avatar)
	if err != nil {
		h.wrapError(c, err, "User profile update failed")
		return
	}
	OK(c, u)
}

// List 用户列表（管理员）。
func (h *UserHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.svc.List(q.Page, q.PageSize)
	if err != nil {
		h.wrapError(c, err, "User list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

func (h *UserHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("user handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("user handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
