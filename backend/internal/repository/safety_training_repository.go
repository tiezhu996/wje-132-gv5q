package repository

import (
	"errors"
	"fmt"
	"time"

	"safetyplatform/internal/model"

	"gorm.io/gorm"
)

// SafetyTrainingRepository 安全培训仓储。
type SafetyTrainingRepository struct {
	db *gorm.DB
}

// NewSafetyTrainingRepository 构造安全培训仓储。
func NewSafetyTrainingRepository(db *gorm.DB) *SafetyTrainingRepository {
	return &SafetyTrainingRepository{db: db}
}

// Create 创建培训。
func (r *SafetyTrainingRepository) Create(t *model.SafetyTraining) error {
	if err := r.db.Create(t).Error; err != nil {
		return fmt.Errorf("create safety training: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询培训。
func (r *SafetyTrainingRepository) FindByID(id uint64) (*model.SafetyTraining, error) {
	var t model.SafetyTraining
	if err := r.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find safety training by id: %w", err)
	}
	return &t, nil
}

// List 分页查询培训。
func (r *SafetyTrainingRepository) List(page, pageSize int, trainingType string) ([]model.SafetyTraining, int64, error) {
	var list []model.SafetyTraining
	var total int64
	q := r.db.Model(&model.SafetyTraining{})
	if trainingType != "" {
		q = q.Where("training_type = ?", trainingType)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count trainings: %w", err)
	}
	if err := q.Order("training_date DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list trainings: %w", err)
	}
	return list, total, nil
}

// Update 更新培训。
func (r *SafetyTrainingRepository) Update(t *model.SafetyTraining) error {
	if err := r.db.Save(t).Error; err != nil {
		return fmt.Errorf("update safety training: %w", err)
	}
	return nil
}

// CompletedRate 本月培训完成率（按有成绩记录计）。
func (r *SafetyTrainingRepository) CompletedRate() (map[string]float64, error) {
	start := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	var total, done int64
	if err := r.db.Model(&model.SafetyTraining{}).Where("training_date >= ? AND training_date < ?", start, end).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count month trainings: %w", err)
	}
	if err := r.db.Model(&model.SafetyTraining{}).Where("training_date >= ? AND training_date < ? AND pass_rate > 0", start, end).Count(&done).Error; err != nil {
		return nil, fmt.Errorf("count done trainings: %w", err)
	}
	rate := 0.0
	if total > 0 {
		rate = float64(done) / float64(total) * 100
	}
	return map[string]float64{"total": float64(total), "done": float64(done), "rate": rate}, nil
}
