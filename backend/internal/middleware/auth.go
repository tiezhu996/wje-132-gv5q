package middleware

import (
	"net/http"
	"strings"

	"safetyplatform/internal/config"
	"safetyplatform/internal/constants"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "user_id"
	ctxPhone  = "phone"
	ctxRole   = "role"
)

// AuthRequired 验证 JWT，将用户信息注入 gin.Context。
func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": constants.CodeUnauthorized, "message": constants.MsgUnauthorized, "data": nil})
			return
		}
		claims, err := util.ParseToken(cfg.JWTSecret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": constants.CodeUnauthorized, "message": constants.MsgUnauthorized, "data": nil})
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxPhone, claims.Phone)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

// JWTConfig 将 JWT 配置注入上下文。
func JWTConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwt_secret", cfg.JWTSecret)
		c.Set("jwt_expire_hours", cfg.JWTExpireHours)
		c.Next()
	}
}

// GetUserID 从上下文取用户 ID。
func GetUserID(c *gin.Context) uint64 {
	v, _ := c.Get(ctxUserID)
	id, _ := v.(uint64)
	return id
}

// GetPhone 从上下文取手机号。
func GetPhone(c *gin.Context) string {
	v, _ := c.Get(ctxPhone)
	s, _ := v.(string)
	return s
}

// GetRole 从上下文取角色。
func GetRole(c *gin.Context) string {
	v, _ := c.Get(ctxRole)
	s, _ := v.(string)
	return s
}
