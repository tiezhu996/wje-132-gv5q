package model

import "time"

// SafetyTraining 安全培训实体。
type SafetyTraining struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Topic            string    `gorm:"size:200;not null" json:"topic"`
	TrainingType     string    `gorm:"size:30;not null;default:regular" json:"training_type"`
	TrainingDate     time.Time `json:"training_date"`
	DurationHours    int       `gorm:"not null;default:1" json:"duration_hours"`
	Trainer          string    `gorm:"size:50;not null;default:''" json:"trainer"`
	Location         string    `gorm:"size:255;not null;default:''" json:"location"`
	ContentSummary   string    `gorm:"type:text" json:"content_summary"`
	ParticipantIDs   JSONList  `gorm:"type:json" json:"participant_ids"`
	AssessmentMethod string    `gorm:"size:30;not null;default:笔试" json:"assessment_method"`
	PassRate         float64   `gorm:"type:decimal(5,2);not null;default:0" json:"pass_rate"`
	CreatedAt        time.Time `json:"created_at"`
}

// TableName 指定表名。
func (SafetyTraining) TableName() string { return "safety_trainings" }
