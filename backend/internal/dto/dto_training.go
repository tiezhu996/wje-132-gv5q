package dto

import "time"

// TrainingCreateRequest 创建培训请求。
type TrainingCreateRequest struct {
	Topic            string    `json:"topic" binding:"required,max=200"`
	TrainingType     string    `json:"training_type" binding:"required,oneof=induction regular special emergency"`
	TrainingDate     time.Time `json:"training_date" binding:"required"`
	DurationHours    int       `json:"duration_hours" binding:"min=1"`
	Trainer          string    `json:"trainer" binding:"max=50"`
	Location         string    `json:"location" binding:"max=255"`
	ContentSummary   string    `json:"content_summary"`
	ParticipantIDs   []string  `json:"participant_ids"`
	AssessmentMethod string    `json:"assessment_method" binding:"max=30"`
}

// TrainingRecordRequest 记录成绩请求。
type TrainingRecordRequest struct {
	ParticipantIDs []string `json:"participant_ids"`
	PassRate       float64  `json:"pass_rate" binding:"min=0,max=100"`
}
