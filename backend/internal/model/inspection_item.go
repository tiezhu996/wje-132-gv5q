package model

// InspectionItem 检查项实体。
type InspectionItem struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	InspectionID uint64 `gorm:"not null;index" json:"inspection_id"`
	ItemName     string `gorm:"size:200;not null" json:"item_name"`
	Passed       bool   `gorm:"not null;default:false" json:"passed"`
	Remark       string `gorm:"size:255;not null;default:''" json:"remark"`
	PhotoURL     string `gorm:"size:255;not null;default:''" json:"photo_url"`
}

// TableName 指定表名。
func (InspectionItem) TableName() string { return "inspection_items" }
