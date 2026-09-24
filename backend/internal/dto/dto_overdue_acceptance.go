package dto

import (
	"time"

	"safetyplatform/internal/model"
)

// SupervisionView 督办信息。
type SupervisionView struct {
	Note         string    `json:"note"`
	OperatorID   uint64    `json:"operator_id"`
	OperatorName string    `json:"operator_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// OverdueAcceptanceItem 逾期验收队列条目。
type OverdueAcceptanceItem struct {
	model.SafetyIncident
	OverdueSeconds int64            `json:"overdue_seconds"`
	Supervision    *SupervisionView `json:"supervision"`
}

// NewOverdueAcceptanceItem 组装逾期队列条目。
func NewOverdueAcceptanceItem(i model.SafetyIncident, now time.Time, sup *model.IncidentSupervision) OverdueAcceptanceItem {
	item := OverdueAcceptanceItem{
		SafetyIncident: i,
		OverdueSeconds: int64(now.Sub(*i.RectificationDeadline).Seconds()),
	}
	if sup != nil {
		item.Supervision = &SupervisionView{
			Note:         sup.Note,
			OperatorID:   sup.OperatorID,
			OperatorName: sup.OperatorName,
			CreatedAt:    sup.CreatedAt,
		}
	}
	return item
}
