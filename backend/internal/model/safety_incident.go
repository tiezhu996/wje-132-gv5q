package model

import "time"

// SafetyIncident 安全事件实体。
type SafetyIncident struct {
	ID                    uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title                 string     `gorm:"size:200;not null" json:"title"`
	Description           string     `gorm:"type:text" json:"description"`
	OccurredAt            time.Time  `json:"occurred_at"`
	SiteID                string     `gorm:"size:50;not null;default:''" json:"site_id"`
	Area                  string     `gorm:"size:100;not null;default:''" json:"area"`
	SeverityLevel         string     `gorm:"size:30;not null;default:minor;index" json:"severity_level"`
	Category              string     `gorm:"size:50;not null;default:其他" json:"category"`
	InvolvedUserIDs       JSONList   `gorm:"type:json" json:"involved_user_ids"`
	PhotoURLs             JSONList   `gorm:"type:json" json:"photo_urls"`
	Status                string     `gorm:"size:30;not null;default:reported;index" json:"status"`
	RectificationMeasures string     `gorm:"type:text" json:"rectification_measures"`
	RectificationDeadline *time.Time `json:"rectification_deadline"`
	ReporterID            uint64     `gorm:"not null" json:"reporter_id"`
	CreatedAt             time.Time  `json:"created_at"`
}

// TableName 指定表名。
func (SafetyIncident) TableName() string { return "safety_incidents" }
