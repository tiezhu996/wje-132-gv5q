package dto

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,max=20"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Name     string `json:"name" binding:"max=50"`
	Role     string `json:"role" binding:"omitempty,oneof=admin safety_manager inspector worker"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应。
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// UpdateProfileRequest 修改资料请求。
type UpdateProfileRequest struct {
	Name   string `json:"name" binding:"max=50"`
	Avatar string `json:"avatar" binding:"max=255"`
}
