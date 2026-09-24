package constants

// 统一错误码。
const (
	CodeOK                           = 0
	CodeBadRequest                   = 40000
	CodeUnauthorized                 = 40100
	CodeForbidden                    = 40300
	CodeNotFound                     = 40400
	CodeConflict                     = 40900
	CodeTooManyRequests              = 42900
	CodeValidationFailed             = 42200
	CodeInternalError                = 50000
	CodeUserExists                   = 40001
	CodeInvalidCredentials           = 40101
	CodeIncidentStatusConflict       = 40901
	CodeIncidentSupervisionDuplicate = 40902
	CodeIncidentNotOverdue           = 40903
	CodeUploadTooLarge               = 41300
	CodeUnsupportedType              = 41500
)
