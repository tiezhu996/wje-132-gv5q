package model

import "time"

// IncidentSupervision 事件督办记录：一条已整改逾期事件仅允许一次督办。
type IncidentSupervision struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	IncidentID   uint64    `gorm:"not null;uniqueIndex:uk_supervision_incident" json:"incident_id"`
	Note         string    `gorm:"type:text;not null" json:"note"`
	OperatorID   uint64    `gorm:"not null;default:0" json:"operator_id"`
	OperatorName string    `gorm:"size:50;not null;default:''" json:"operator_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名。
func (IncidentSupervision) TableName() string { return "incident_supervisions" }
