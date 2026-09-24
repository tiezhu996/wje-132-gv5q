package repository

import (
	"fmt"

	"safetyplatform/internal/model"

	"gorm.io/gorm"
)

// InspectionItemRepository 检查项仓储。
type InspectionItemRepository struct {
	db *gorm.DB
}

// NewInspectionItemRepository 构造检查项仓储。
func NewInspectionItemRepository(db *gorm.DB) *InspectionItemRepository {
	return &InspectionItemRepository{db: db}
}

// CreateMany 批量创建检查项。
func (r *InspectionItemRepository) CreateMany(items []model.InspectionItem) error {
	return r.CreateManyTx(r.db, items)
}

// CreateManyTx 在事务内批量创建检查项。
func (r *InspectionItemRepository) CreateManyTx(tx *gorm.DB, items []model.InspectionItem) error {
	if len(items) == 0 {
		return nil
	}
	if err := tx.Create(&items).Error; err != nil {
		return fmt.Errorf("create inspection items: %w", err)
	}
	return nil
}

// ListByInspection 查询某检查的检查项。
func (r *InspectionItemRepository) ListByInspection(inspectionID uint64) ([]model.InspectionItem, error) {
	return r.ListByInspectionTx(r.db, inspectionID)
}

// ListByInspectionTx 在事务内查询某检查的检查项。
func (r *InspectionItemRepository) ListByInspectionTx(tx *gorm.DB, inspectionID uint64) ([]model.InspectionItem, error) {
	var list []model.InspectionItem
	if err := tx.Where("inspection_id = ?", inspectionID).Order("id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list inspection items: %w", err)
	}
	return list, nil
}

// Update 更新检查项。
func (r *InspectionItemRepository) Update(item *model.InspectionItem) error {
	return r.UpdateTx(r.db, item)
}

// UpdateTx 在事务内更新检查项。
func (r *InspectionItemRepository) UpdateTx(tx *gorm.DB, item *model.InspectionItem) error {
	if err := tx.Save(item).Error; err != nil {
		return fmt.Errorf("update inspection item: %w", err)
	}
	return nil
}

// DeleteByInspection 删除某检查的全部检查项。
func (r *InspectionItemRepository) DeleteByInspection(inspectionID uint64) error {
	if err := r.db.Where("inspection_id = ?", inspectionID).Delete(&model.InspectionItem{}).Error; err != nil {
		return fmt.Errorf("delete inspection items: %w", err)
	}
	return nil
}
