package dto

import "time"

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
