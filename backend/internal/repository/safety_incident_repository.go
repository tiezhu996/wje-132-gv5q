package repository

import (
	"errors"
	"fmt"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"

	"gorm.io/gorm"
)

// SafetyIncidentRepository 安全事件仓储。
type SafetyIncidentRepository struct {
	db *gorm.DB
}

// NewSafetyIncidentRepository 构造安全事件仓储。
func NewSafetyIncidentRepository(db *gorm.DB) *SafetyIncidentRepository {
	return &SafetyIncidentRepository{db: db}
}

// Create 创建事件。
func (r *SafetyIncidentRepository) Create(i *model.SafetyIncident) error {
	if err := r.db.Create(i).Error; err != nil {
		return fmt.Errorf("create safety incident: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询事件。
func (r *SafetyIncidentRepository) FindByID(id uint64) (*model.SafetyIncident, error) {
	var i model.SafetyIncident
	if err := r.db.First(&i, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find safety incident by id: %w", err)
	}
	return &i, nil
}

// List 分页查询事件，支持严重等级/状态/时间筛选。
func (r *SafetyIncidentRepository) List(page, pageSize int, severity, status string, startDate, endDate *time.Time) ([]model.SafetyIncident, int64, error) {
	var list []model.SafetyIncident
	var total int64
	q := r.db.Model(&model.SafetyIncident{})
	if severity != "" {
		q = q.Where("severity_level = ?", severity)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if startDate != nil {
		q = q.Where("occurred_at >= ?", startDate)
	}
	if endDate != nil {
		q = q.Where("occurred_at <= ?", endDate)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count incidents: %w", err)
	}
	if err := q.Order("occurred_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list incidents: %w", err)
	}
	return list, total, nil
}

// Update 更新事件。
func (r *SafetyIncidentRepository) Update(i *model.SafetyIncident) error {
	if err := r.db.Save(i).Error; err != nil {
		return fmt.Errorf("update safety incident: %w", err)
	}
	return nil
}

// Trend30 近 30 天事件趋势。
func (r *SafetyIncidentRepository) Trend30() ([]map[string]any, error) {
	start := time.Now().AddDate(0, 0, -29)
	var rows []map[string]any
	if err := r.db.Model(&model.SafetyIncident{}).
		Select("DATE(occurred_at) AS day, COUNT(*) AS cnt").
		Where("occurred_at >= ?", start).
		Group("DATE(occurred_at)").Order("day ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("incident trend: %w", err)
	}
	return rows, nil
}

// SeverityDistribution 严重等级分布。
func (r *SafetyIncidentRepository) SeverityDistribution() ([]map[string]any, error) {
	var rows []map[string]any
	if err := r.db.Model(&model.SafetyIncident{}).
		Select("severity_level, COUNT(*) AS cnt").Group("severity_level").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("severity distribution: %w", err)
	}
	return rows, nil
}

// OverdueAcceptance 逾期验收队列：已整改（resolved）且整改期限早于当前时刻。
// 排序：逾期时长从长到短（期限升序），风险等级从高到低。
func (r *SafetyIncidentRepository) OverdueAcceptance(now time.Time) ([]model.SafetyIncident, error) {
	var list []model.SafetyIncident
	severityOrder := "CASE severity_level " +
		"WHEN '" + constants.SeverityFatal + "' THEN 5 " +
		"WHEN '" + constants.SeverityMajor + "' THEN 4 " +
		"WHEN '" + constants.SeverityModerate + "' THEN 3 " +
		"WHEN '" + constants.SeverityMinor + "' THEN 2 " +
		"WHEN '" + constants.SeverityNearMiss + "' THEN 1 " +
		"ELSE 0 END DESC"
	if err := r.db.Where("status = ? AND rectification_deadline IS NOT NULL AND rectification_deadline < ?",
		constants.IncidentResolved, now).
		Order("rectification_deadline ASC").Order(severityOrder).
		Find(&list).Error; err != nil {
		return nil, fmt.Errorf("overdue acceptance list: %w", err)
	}
	return list, nil
}

// Supervise 对事件发起一次督办。仅当事件处于已整改且尚未督办时更新成功，
// 返回 rows==0 表示记录不存在、状态已变更或已被督办（调用方需区分冲突原因）。
func (r *SafetyIncidentRepository) Supervise(id uint64, operatorID uint64, note string, now time.Time) (int64, error) {
	res := r.db.Model(&model.SafetyIncident{}).
		Where("id = ? AND status = ? AND supervised_by = 0", id, constants.IncidentResolved).
		Updates(map[string]any{
			"supervision_note": note,
			"supervision_at":   now,
			"supervised_by":    operatorID,
		})
	if res.Error != nil {
		return 0, fmt.Errorf("supervise safety incident: %w", res.Error)
	}
	return res.RowsAffected, nil
}

// PendingRectification 待整改事件。
func (r *SafetyIncidentRepository) PendingRectification() ([]model.SafetyIncident, error) {
	var list []model.SafetyIncident
	if err := r.db.Where("status IN ?", []string{"reported", "investigating"}).
		Order("rectification_deadline ASC").Limit(10).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("pending rectification: %w", err)
	}
	return list, nil
}
