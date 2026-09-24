package repository

import (
	"errors"
	"fmt"
	"strings"

	"safetyplatform/internal/model"

	"gorm.io/gorm"
)

// IncidentSupervisionRepository 事件督办仓储。
type IncidentSupervisionRepository struct {
	db *gorm.DB
}

// NewIncidentSupervisionRepository 构造事件督办仓储。
func NewIncidentSupervisionRepository(db *gorm.DB) *IncidentSupervisionRepository {
	return &IncidentSupervisionRepository{db: db}
}

// Create 创建督办记录；同一事件重复督办时返回 ErrDuplicate。
func (r *IncidentSupervisionRepository) Create(s *model.IncidentSupervision) error {
	if err := r.db.Create(s).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || isDupErr(err) {
			return ErrDuplicate
		}
		return fmt.Errorf("create incident supervision: %w", err)
	}
	return nil
}

// isDupErr 兼容 MySQL（Duplicate entry）与 SQLite（UNIQUE constraint failed）的唯一键报错。
func isDupErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate entry") || strings.Contains(msg, "unique constraint failed")
}

// FindByIncidentID 按事件 ID 查询督办记录，未督办时返回 ErrNotFound。
func (r *IncidentSupervisionRepository) FindByIncidentID(incidentID uint64) (*model.IncidentSupervision, error) {
	var s model.IncidentSupervision
	if err := r.db.Where("incident_id = ?", incidentID).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find incident supervision by incident id: %w", err)
	}
	return &s, nil
}

// MapByIncidentIDs 批量按事件 ID 查询督办记录，key 为 incident_id。
func (r *IncidentSupervisionRepository) MapByIncidentIDs(incidentIDs []uint64) (map[uint64]model.IncidentSupervision, error) {
	result := make(map[uint64]model.IncidentSupervision)
	if len(incidentIDs) == 0 {
		return result, nil
	}
	var list []model.IncidentSupervision
	if err := r.db.Where("incident_id IN ?", incidentIDs).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list incident supervisions by incident ids: %w", err)
	}
	for i := range list {
		result[list[i].IncidentID] = list[i]
	}
	return result, nil
}
