package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误响应格式。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			status := appErrorStatus(appErr.Code)
			c.JSON(status, gin.H{"code": appErr.Code, "message": appErr.Message, "data": nil})
			return
		}
		logger.Error("unhandled error", "error", err.Error(), "path", c.FullPath())
		c.JSON(http.StatusInternalServerError, gin.H{"code": constants.CodeInternalError, "message": constants.MsgInternalError, "data": nil})
	}
}

func appErrorStatus(code int) int {
	switch code {
	case constants.CodeUnauthorized, constants.CodeInvalidCredentials:
		return http.StatusUnauthorized
	case constants.CodeForbidden:
		return http.StatusForbidden
	case constants.CodeNotFound:
		return http.StatusNotFound
	case constants.CodeConflict, constants.CodeIncidentStatusConflict:
		return http.StatusConflict
	case constants.CodeValidationFailed:
		return http.StatusUnprocessableEntity
	case constants.CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusBadRequest
	}
}
