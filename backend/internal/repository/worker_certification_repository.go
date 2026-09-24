package repository

import (
	"errors"
	"fmt"
	"time"

	"safetyplatform/internal/model"

	"gorm.io/gorm"
)

// WorkerCertificationRepository 人员资质仓储。
type WorkerCertificationRepository struct {
	db *gorm.DB
}

// NewWorkerCertificationRepository 构造人员资质仓储。
func NewWorkerCertificationRepository(db *gorm.DB) *WorkerCertificationRepository {
	return &WorkerCertificationRepository{db: db}
}

// Create 创建资质。
func (r *WorkerCertificationRepository) Create(c *model.WorkerCertification) error {
	if err := r.db.Create(c).Error; err != nil {
		return fmt.Errorf("create worker certification: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询资质。
func (r *WorkerCertificationRepository) FindByID(id uint64) (*model.WorkerCertification, error) {
	var c model.WorkerCertification
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find worker certification by id: %w", err)
	}
	return &c, nil
}

// List 分页查询资质，支持状态/过期筛选。
func (r *WorkerCertificationRepository) List(page, pageSize int, status string, expiring bool) ([]model.WorkerCertification, int64, error) {
	var list []model.WorkerCertification
	var total int64
	q := r.db.Model(&model.WorkerCertification{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if expiring {
		q = q.Where("valid_until IS NOT NULL AND valid_until <= ?", time.Now().AddDate(0, 1, 0))
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count certifications: %w", err)
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list certifications: %w", err)
	}
	return list, total, nil
}

// ListByUser 查询某用户的资质。
func (r *WorkerCertificationRepository) ListByUser(userID uint64) ([]model.WorkerCertification, error) {
	var list []model.WorkerCertification
	if err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list certifications by user: %w", err)
	}
	return list, nil
}

// Update 更新资质。
func (r *WorkerCertificationRepository) Update(c *model.WorkerCertification) error {
	if err := r.db.Save(c).Error; err != nil {
		return fmt.Errorf("update worker certification: %w", err)
	}
	return nil
}

// ExpiringSoon 即将过期资质。
func (r *WorkerCertificationRepository) ExpiringSoon() ([]model.WorkerCertification, error) {
	var list []model.WorkerCertification
	from := time.Now()
	to := time.Now().AddDate(0, 1, 0)
	if err := r.db.Where("valid_until IS NOT NULL AND valid_until >= ? AND valid_until <= ?", from, to).
		Order("valid_until ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list expiring certifications: %w", err)
	}
	return list, nil
}
