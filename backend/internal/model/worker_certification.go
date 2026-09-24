package model

import "time"

// WorkerCertification 人员资质实体。
type WorkerCertification struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"not null;index" json:"user_id"`
	CertType     string     `gorm:"size:50;not null;default:''" json:"cert_type"`
	CertNo       string     `gorm:"size:50;not null;default:''" json:"cert_no"`
	IssueOrg     string     `gorm:"size:100;not null;default:''" json:"issue_org"`
	IssueDate    *time.Time `json:"issue_date"`
	ValidUntil   *time.Time `json:"valid_until"`
	CertPhotoURL string     `gorm:"size:255;not null;default:''" json:"cert_photo_url"`
	Status       string     `gorm:"size:30;not null;default:pending;index" json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TableName 指定表名。
func (WorkerCertification) TableName() string { return "worker_certifications" }
