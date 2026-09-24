package dto

import "time"

// InspectionCreateRequest 创建检查计划请求。
type InspectionCreateRequest struct {
	Name           string              `json:"name" binding:"required,max=200"`
	InspectionType string              `json:"inspection_type" binding:"required,oneof=routine special pre_shift emergency"`
	Area           string              `json:"area" binding:"max=100"`
	InspectionDate time.Time           `json:"inspection_date" binding:"required"`
	InspectorID    uint64              `json:"inspector_id" binding:"required"`
	Items          []InspectionItemDTO `json:"items"`
}

// InspectionItemDTO 检查项 DTO。
type InspectionItemDTO struct {
	ID       uint64 `json:"id"`
	ItemName string `json:"item_name" binding:"required,max=200"`
	Passed   bool   `json:"passed"`
	Remark   string `json:"remark" binding:"max=255"`
	PhotoURL string `json:"photo_url" binding:"max=255"`
}

// InspectionExecuteRequest 执行检查请求。
type InspectionExecuteRequest struct {
	Items []InspectionItemDTO `json:"items" binding:"required"`
}
