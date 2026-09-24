package constants

// SeverityLevel 严重等级枚举。
const (
	SeverityNearMiss = "near_miss"
	SeverityMinor    = "minor"
	SeverityModerate = "moderate"
	SeverityMajor    = "major"
	SeverityFatal    = "fatal"
)

// SeverityValues 全部严重等级值。
var SeverityValues = []string{SeverityNearMiss, SeverityMinor, SeverityModerate, SeverityMajor, SeverityFatal}

// IncidentStatus 事件状态枚举。
const (
	IncidentReported      = "reported"
	IncidentInvestigating = "investigating"
	IncidentResolved      = "resolved"
	IncidentClosed        = "closed"
)

// IncidentStatusValues 全部事件状态值。
var IncidentStatusValues = []string{IncidentReported, IncidentInvestigating, IncidentResolved, IncidentClosed}

// InspectionStatus 检查状态枚举。
const (
	InspectionScheduled  = "scheduled"
	InspectionInProgress = "in_progress"
	InspectionCompleted  = "completed"
	InspectionFailed     = "failed"
)

// InspectionType 检查类型枚举。
const (
	InspectionRoutine   = "routine"
	InspectionSpecial   = "special"
	InspectionPreShift  = "pre_shift"
	InspectionEmergency = "emergency"
)

// TrainingType 培训类型枚举。
const (
	TrainingInduction = "induction"
	TrainingRegular   = "regular"
	TrainingSpecial   = "special"
	TrainingEmergency = "emergency"
)

// CertStatus 资质状态枚举。
const (
	CertPending  = "pending"
	CertApproved = "approved"
	CertRejected = "rejected"
	CertExpired  = "expired"
)

// IsValidSeverity 校验严重等级。
func IsValidSeverity(s string) bool {
	for _, v := range SeverityValues {
		if v == s {
			return true
		}
	}
	return false
}

// IsValidIncidentStatus 校验事件状态。
func IsValidIncidentStatus(s string) bool {
	for _, v := range IncidentStatusValues {
		if v == s {
			return true
		}
	}
	return false
}

// IncidentCategories 事件分类。
var IncidentCategories = []string{"坠落", "触电", "物体打击", "坍塌", "机械伤害", "其他"}
