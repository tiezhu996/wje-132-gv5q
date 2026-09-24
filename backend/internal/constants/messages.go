package constants

// messages.go 同时承载前端提示文案、后端返回文案与日志文案。
const (
	MsgOK                  = "ok"
	MsgUnauthorized        = "未登录或登录已过期"
	MsgForbidden           = "没有权限执行该操作"
	MsgNotFound            = "资源不存在"
	MsgTooManyRequests     = "请求过于频繁，请稍后再试"
	MsgInternalError       = "服务器内部错误"
	MsgPhoneExists         = "手机号已注册"
	MsgInvalidCredentials  = "手机号或密码错误"
	MsgIncidentReported    = "事件上报成功"
	MsgIncidentClosed      = "事件已关闭"
	MsgInspectionCompleted = "检查已完成"
	MsgCertSubmitted       = "资质提交成功"
	MsgCertReviewed        = "审核完成"
	MsgUploadTooLarge      = "上传文件过大"
	MsgUnsupportedFileType = "不支持的文件类型"
	MsgLoginSuccess        = "登录成功"
)
