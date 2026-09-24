package util

import (
	"fmt"
	"time"

	"safetyplatform/internal/constants"
)

// formatters.go 同时提供日期格式化、严重等级文本、事件状态文本、检查状态文本、培训类型文本、资质状态文本。

// FormatDateTime 格式化日期时间。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate 格式化日期。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// SeverityText 严重等级文本。
func SeverityText(level string) string {
	switch level {
	case constants.SeverityNearMiss:
		return "未遂"
	case constants.SeverityMinor:
		return "轻微"
	case constants.SeverityModerate:
		return "一般"
	case constants.SeverityMajor:
		return "较大"
	case constants.SeverityFatal:
		return "重大"
	default:
		return level
	}
}

// IncidentStatusText 事件状态文本。
func IncidentStatusText(s string) string {
	switch s {
	case constants.IncidentReported:
		return "已上报"
	case constants.IncidentInvestigating:
		return "调查中"
	case constants.IncidentResolved:
		return "已整改"
	case constants.IncidentClosed:
		return "已关闭"
	default:
		return s
	}
}

// InspectionStatusText 检查状态文本。
func InspectionStatusText(s string) string {
	switch s {
	case constants.InspectionScheduled:
		return "待执行"
	case constants.InspectionInProgress:
		return "执行中"
	case constants.InspectionCompleted:
		return "已完成"
	case constants.InspectionFailed:
		return "不合格"
	default:
		return s
	}
}

// TrainingTypeText 培训类型文本。
func TrainingTypeText(t string) string {
	switch t {
	case constants.TrainingInduction:
		return "入场教育"
	case constants.TrainingRegular:
		return "常规培训"
	case constants.TrainingSpecial:
		return "专项培训"
	case constants.TrainingEmergency:
		return "应急演练"
	default:
		return t
	}
}

// CertStatusText 资质状态文本。
func CertStatusText(s string) string {
	switch s {
	case constants.CertPending:
		return "待审核"
	case constants.CertApproved:
		return "已通过"
	case constants.CertRejected:
		return "已拒绝"
	case constants.CertExpired:
		return "已过期"
	default:
		return s
	}
}

// UserRoleText 角色文本。
func UserRoleText(r string) string {
	switch r {
	case constants.RoleAdmin:
		return "管理员"
	case constants.RoleSafetyManager:
		return "安全管理员"
	case constants.RoleInspector:
		return "监理"
	case constants.RoleWorker:
		return "工人"
	default:
		return r
	}
}

// FormatScore 格式化分数文本。
func FormatScore(score int) string {
	return fmt.Sprintf("%d 分", score)
}
