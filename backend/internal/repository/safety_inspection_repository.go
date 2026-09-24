package repository

import (
	"errors"
	"fmt"
	"time"

	"safetyplatform/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SafetyInspectionRepository 安全检查仓储。
type SafetyInspectionRepository struct {
	db *gorm.DB
}

// NewSafetyInspectionRepository 构造安全检查仓储。
func NewSafetyInspectionRepository(db *gorm.DB) *SafetyInspectionRepository {
	return &SafetyInspectionRepository{db: db}
}

// Create 创建检查计划。
func (r *SafetyInspectionRepository) Create(i *model.SafetyInspection) error {
	return r.CreateTx(r.db, i)
}

// CreateTx 在事务内创建检查计划。
func (r *SafetyInspectionRepository) CreateTx(tx *gorm.DB, i *model.SafetyInspection) error {
	if err := tx.Create(i).Error; err != nil {
		return fmt.Errorf("create safety inspection: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询检查。
func (r *SafetyInspectionRepository) FindByID(id uint64) (*model.SafetyInspection, error) {
	return r.findByID(r.db, id, false)
}

// FindByIDForUpdate 在事务内锁定检查行。
func (r *SafetyInspectionRepository) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.SafetyInspection, error) {
	return r.findByID(tx, id, true)
}

func (r *SafetyInspectionRepository) findByID(db *gorm.DB, id uint64, forUpdate bool) (*model.SafetyInspection, error) {
	var i model.SafetyInspection
	q := db
	if forUpdate {
		q = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.First(&i, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find safety inspection by id: %w", err)
	}
	return &i, nil
}

// List 分页查询检查计划。
func (r *SafetyInspectionRepository) List(page, pageSize int, status string) ([]model.SafetyInspection, int64, error) {
	var list []model.SafetyInspection
	var total int64
	q := r.db.Model(&model.SafetyInspection{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count inspections: %w", err)
	}
	if err := q.Order("inspection_date DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list inspections: %w", err)
	}
	return list, total, nil
}

// Update 更新检查。
func (r *SafetyInspectionRepository) Update(i *model.SafetyInspection) error {
	return r.UpdateTx(r.db, i)
}

// UpdateTx 在事务内更新检查。
func (r *SafetyInspectionRepository) UpdateTx(tx *gorm.DB, i *model.SafetyInspection) error {
	if err := tx.Save(i).Error; err != nil {
		return fmt.Errorf("update safety inspection: %w", err)
	}
	return nil
}

// Count 统计检查总数。
func (r *SafetyInspectionRepository) Count() (int64, error) {
	var n int64
	if err := r.db.Model(&model.SafetyInspection{}).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count inspections: %w", err)
	}
	return n, nil
}

// CompletedRate 本月完成率。
func (r *SafetyInspectionRepository) CompletedRate() (map[string]float64, error) {
	start := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	var total, completed int64
	if err := r.db.Model(&model.SafetyInspection{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count month inspections: %w", err)
	}
	if err := r.db.Model(&model.SafetyInspection{}).
		Where("created_at >= ? AND created_at < ? AND status = ?", start, end, "completed").Count(&completed).Error; err != nil {
		return nil, fmt.Errorf("count completed: %w", err)
	}
	rate := 0.0
	if total > 0 {
		rate = float64(completed) / float64(total) * 100
	}
	return map[string]float64{"total": float64(total), "completed": float64(completed), "rate": rate}, nil
}
