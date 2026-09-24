package dto

import (
	"time"

	"safetyplatform/internal/model"
)

// IncidentReportRequest 上报事件请求。
type IncidentReportRequest struct {
	Title           string    `json:"title" binding:"required,max=200"`
	Description     string    `json:"description"`
	OccurredAt      time.Time `json:"occurred_at" binding:"required"`
	SiteID          string    `json:"site_id" binding:"max=50"`
	Area            string    `json:"area" binding:"max=100"`
	SeverityLevel   string    `json:"severity_level" binding:"required,oneof=near_miss minor moderate major fatal"`
	Category        string    `json:"category" binding:"max=50"`
	InvolvedUserIDs []string  `json:"involved_user_ids"`
	PhotoURLs       []string  `json:"photo_urls"`
}

// RectificationRequest 整改请求。
type RectificationRequest struct {
	Measures string     `json:"measures" binding:"required"`
	Deadline *time.Time `json:"deadline"`
}

// SupervisionRequest 督办请求。
type SupervisionRequest struct {
	Note string `json:"note" binding:"required,max=500"`
}

// OverdueIncident 逾期验收队列条目，附带逾期时长（秒）。
type OverdueIncident struct {
	model.SafetyIncident
	OverdueDuration int64 `json:"overdue_duration"`
}
