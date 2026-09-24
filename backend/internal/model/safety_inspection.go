package model

import "time"

// SafetyInspection 安全检查实体。
type SafetyInspection struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	InspectionType string    `gorm:"size:30;not null;default:routine" json:"inspection_type"`
	Area           string    `gorm:"size:100;not null;default:''" json:"area"`
	InspectionDate time.Time `json:"inspection_date"`
	InspectorID    uint64    `gorm:"not null" json:"inspector_id"`
	TotalScore     int       `gorm:"not null;default:0" json:"total_score"`
	Status         string    `gorm:"size:30;not null;default:scheduled;index" json:"status"`
	IssueCount     int       `gorm:"not null;default:0" json:"issue_count"`
	PassedCount    int       `gorm:"not null;default:0" json:"passed_count"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName 指定表名。
func (SafetyInspection) TableName() string { return "safety_inspections" }
