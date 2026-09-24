package service

import (
	"log/slog"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"

	"gorm.io/gorm"
)

// SafetyInspectionService 安全检查业务逻辑。
type SafetyInspectionService struct {
	db       *gorm.DB
	repo     *repository.SafetyInspectionRepository
	itemRepo *repository.InspectionItemRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewSafetyInspectionService 构造安全检查服务。
func NewSafetyInspectionService(db *gorm.DB, repo *repository.SafetyInspectionRepository,
	itemRepo *repository.InspectionItemRepository, userRepo *repository.UserRepository, logger *slog.Logger) *SafetyInspectionService {
	return &SafetyInspectionService{db: db, repo: repo, itemRepo: itemRepo, userRepo: userRepo, logger: logger}
}

// Create 创建检查计划。
func (s *SafetyInspectionService) Create(name, inspectionType, area string, inspectionDate time.Time, inspectorID uint64, items []model.InspectionItem) (*model.SafetyInspection, error) {
	if _, err := s.userRepo.FindByID(inspectorID); err != nil {
		return nil, util.Wrap(err, "SafetyInspection[inspector_id=%d] create: inspector not found", inspectorID)
	}
	ins := &model.SafetyInspection{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		ins.Name = name
		ins.InspectionType = inspectionType
		ins.Area = area
		ins.InspectionDate = inspectionDate
		ins.InspectorID = inspectorID
		ins.Status = constants.InspectionScheduled
		if err := s.repo.CreateTx(tx, ins); err != nil {
			return util.Wrap(err, "SafetyInspection[name=%s] create failed", name)
		}
		for i := range items {
			items[i].InspectionID = ins.ID
		}
		if err := s.itemRepo.CreateManyTx(tx, items); err != nil {
			return util.Wrap(err, "SafetyInspection[id=%d] create items failed", ins.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogInspectionCreateSuccess, "inspection_id", ins.ID)
	return ins, nil
}

// Execute 执行检查：逐项更新检查项，计算得分与状态。
func (s *SafetyInspectionService) Execute(id uint64, items []model.InspectionItem) (*model.SafetyInspection, error) {
	ins := &model.SafetyInspection{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		cur, err := s.repo.FindByIDForUpdate(tx, id)
		if err != nil {
			return util.Wrap(err, "SafetyInspection[id=%d] execute find failed", id)
		}
		if cur.Status != constants.InspectionScheduled && cur.Status != constants.InspectionInProgress {
			return util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyInspection[id="+u64(id)+"] execute conflict: status="+cur.Status)
		}
		existing, err := s.itemRepo.ListByInspectionTx(tx, id)
		if err != nil {
			return util.Wrap(err, "SafetyInspection[id=%d] execute list items failed", id)
		}
		byID := make(map[uint64]*model.InspectionItem, len(existing))
		for i := range existing {
			byID[existing[i].ID] = &existing[i]
		}
		passed, issues := 0, 0
		for _, it := range items {
			exist, ok := byID[it.ID]
			if !ok {
				continue
			}
			exist.Passed = it.Passed
			exist.Remark = it.Remark
			exist.PhotoURL = it.PhotoURL
			if err := s.itemRepo.UpdateTx(tx, exist); err != nil {
				return util.Wrap(err, "SafetyInspection[id=%d] execute item update failed", id)
			}
			if exist.Passed {
				passed++
			} else {
				issues++
			}
		}
		if passed == 0 && issues == 0 {
			for i := range existing {
				if existing[i].Passed {
					passed++
				} else {
					issues++
				}
			}
		}
		total := passed + issues
		score := 0
		if total > 0 {
			score = int(float64(passed) / float64(total) * 100)
		}
		cur.PassedCount = passed
		cur.IssueCount = issues
		cur.TotalScore = score
		if issues == 0 {
			cur.Status = constants.InspectionCompleted
		} else {
			cur.Status = constants.InspectionFailed
		}
		if err := s.repo.UpdateTx(tx, cur); err != nil {
			s.logger.Error(constants.LogInspectionExecuteFailed, "error", err.Error())
			return util.Wrap(err, "SafetyInspection[id=%d] execute save failed", id)
		}
		ins = cur
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogInspectionExecuteSuccess, "inspection_id", ins.ID, "score", ins.TotalScore)
	return ins, nil
}

// List 分页查询检查计划。
func (s *SafetyInspectionService) List(page, pageSize int, status string) ([]model.SafetyInspection, int64, error) {
	return s.repo.List(page, pageSize, status)
}

// GetWithItems 检查详情 + 检查项。
func (s *SafetyInspectionService) GetWithItems(id uint64) (*model.SafetyInspection, []model.InspectionItem, error) {
	ins, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, util.Wrap(err, "SafetyInspection[id=%d] get failed", id)
	}
	items, err := s.itemRepo.ListByInspection(id)
	if err != nil {
		return nil, nil, err
	}
	return ins, items, nil
}

// Report 检查报告数据。
func (s *SafetyInspectionService) Report(id uint64) (*model.SafetyInspection, []model.InspectionItem, error) {
	ins, items, err := s.GetWithItems(id)
	if err != nil {
		return nil, nil, err
	}
	s.logger.Info(constants.LogInspectionReport, "inspection_id", id, "score", ins.TotalScore)
	return ins, items, nil
}

// Stats 检查统计。
func (s *SafetyInspectionService) Stats() (map[string]any, error) {
	total, err := s.repo.Count()
	if err != nil {
		return nil, err
	}
	rate, err := s.repo.CompletedRate()
	if err != nil {
		return nil, err
	}
	return map[string]any{"total": total, "completed_rate": rate["rate"]}, nil
}
