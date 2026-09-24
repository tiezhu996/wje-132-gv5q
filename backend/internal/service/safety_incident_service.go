package service

import (
	"errors"
	"log/slog"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/dto"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"
)

// SafetyIncidentService 安全事件业务逻辑。
type SafetyIncidentService struct {
	repo            *repository.SafetyIncidentRepository
	userRepo        *repository.UserRepository
	supervisionRepo *repository.IncidentSupervisionRepository
	logger          *slog.Logger
}

// NewSafetyIncidentService 构造安全事件服务。
func NewSafetyIncidentService(repo *repository.SafetyIncidentRepository, userRepo *repository.UserRepository, supervisionRepo *repository.IncidentSupervisionRepository, logger *slog.Logger) *SafetyIncidentService {
	return &SafetyIncidentService{repo: repo, userRepo: userRepo, supervisionRepo: supervisionRepo, logger: logger}
}

// Report 上报事件。
func (s *SafetyIncidentService) Report(reporterID uint64, title, description string, occurredAt time.Time,
	siteID, area, severity, category string, involvedUserIDs, photoURLs []string) (*model.SafetyIncident, error) {
	if !constants.IsValidSeverity(severity) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "SafetyIncident[severity="+severity+"] report: invalid severity")
	}
	i := &model.SafetyIncident{
		Title: title, Description: description, OccurredAt: occurredAt, SiteID: siteID, Area: area,
		SeverityLevel: severity, Category: category,
		InvolvedUserIDs: model.JSONList(involvedUserIDs), PhotoURLs: model.JSONList(photoURLs),
		Status: constants.IncidentReported, ReporterID: reporterID,
	}
	if err := s.repo.Create(i); err != nil {
		s.logger.Error(constants.LogIncidentReportFailed, "error", err.Error())
		return nil, util.Wrap(err, "SafetyIncident[title=%s] report create failed", title)
	}
	s.logger.Info(constants.LogIncidentReportSuccess, "incident_id", i.ID)
	return i, nil
}

// Assign 指派调查。
func (s *SafetyIncidentService) Assign(id uint64) (*model.SafetyIncident, error) {
	i, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] assign find failed", id)
	}
	if i.Status != constants.IncidentReported {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyIncident[id="+u64(id)+"] assign conflict: status="+i.Status)
	}
	i.Status = constants.IncidentInvestigating
	if err := s.repo.Update(i); err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] assign save failed", id)
	}
	s.logger.Info(constants.LogIncidentAssignSuccess, "incident_id", i.ID)
	return i, nil
}

// SubmitRectification 提交整改。
func (s *SafetyIncidentService) SubmitRectification(id uint64, measures string, deadline *time.Time) (*model.SafetyIncident, error) {
	i, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] rectify find failed", id)
	}
	if i.Status != constants.IncidentInvestigating {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyIncident[id="+u64(id)+"] rectify conflict: status="+i.Status)
	}
	i.RectificationMeasures = measures
	i.RectificationDeadline = deadline
	i.Status = constants.IncidentResolved
	if err := s.repo.Update(i); err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] rectify save failed", id)
	}
	s.logger.Info(constants.LogIncidentRectifySuccess, "incident_id", i.ID)
	return i, nil
}

// Close 关闭事件。
func (s *SafetyIncidentService) Close(id uint64) (*model.SafetyIncident, error) {
	i, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] close find failed", id)
	}
	if i.Status != constants.IncidentResolved {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyIncident[id="+u64(id)+"] close conflict: status="+i.Status)
	}
	i.Status = constants.IncidentClosed
	if err := s.repo.Update(i); err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] close save failed", id)
	}
	s.logger.Info(constants.LogIncidentCloseSuccess, "incident_id", i.ID)
	return i, nil
}

// List 分页查询事件。
func (s *SafetyIncidentService) List(page, pageSize int, severity, status string, startDate, endDate *time.Time) ([]model.SafetyIncident, int64, error) {
	return s.repo.List(page, pageSize, severity, status, startDate, endDate)
}

// Get 事件详情。
func (s *SafetyIncidentService) Get(id uint64) (*model.SafetyIncident, error) {
	return s.repo.FindByID(id)
}

// Trend30 近 30 天趋势。
func (s *SafetyIncidentService) Trend30() ([]map[string]any, error) {
	return s.repo.Trend30()
}

// SeverityDistribution 严重等级分布。
func (s *SafetyIncidentService) SeverityDistribution() ([]map[string]any, error) {
	return s.repo.SeverityDistribution()
}

// PendingRectification 待整改列表。
func (s *SafetyIncidentService) PendingRectification() ([]model.SafetyIncident, error) {
	return s.repo.PendingRectification()
}

// OverdueAcceptanceQueue 逾期验收队列：仅包含已整改且整改期限早于当前时刻的事件，
// 按逾期时长（期限越早越靠前）和风险等级（fatal>major>moderate>minor>near_miss）排序。
func (s *SafetyIncidentService) OverdueAcceptanceQueue() ([]dto.OverdueAcceptanceItem, error) {
	now := time.Now()
	list, err := s.repo.ListOverdueResolved(now, constants.SeverityRankOrder)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident overdue acceptance queue failed")
	}
	ids := make([]uint64, 0, len(list))
	for i := range list {
		ids = append(ids, list[i].ID)
	}
	supMap, err := s.supervisionRepo.MapByIncidentIDs(ids)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident overdue acceptance queue load supervisions failed")
	}
	items := make([]dto.OverdueAcceptanceItem, 0, len(list))
	for i := range list {
		sup, ok := supMap[list[i].ID]
		if !ok {
			items = append(items, dto.NewOverdueAcceptanceItem(list[i], now, nil))
			continue
		}
		supCopy := sup
		items = append(items, dto.NewOverdueAcceptanceItem(list[i], now, &supCopy))
	}
	return items, nil
}

// Supervise 对一条已整改逾期事件发起督办，每条事件仅允许督办一次。
// 已关闭、未整改或期限未到的事件不能督办；重复督办返回冲突。
func (s *SafetyIncidentService) Supervise(incidentID, operatorID uint64, note string) (*model.IncidentSupervision, error) {
	i, err := s.repo.FindByID(incidentID)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] supervise find failed", incidentID)
	}
	if i.Status != constants.IncidentResolved {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict,
			"SafetyIncident[id="+u64(incidentID)+"] supervise conflict: status="+i.Status+" (only resolved accepted)")
	}
	if i.RectificationDeadline == nil || !i.RectificationDeadline.Before(time.Now()) {
		return nil, util.NewAppError(constants.CodeIncidentNotOverdue,
			"SafetyIncident[id="+u64(incidentID)+"] supervise conflict: rectification deadline not reached")
	}
	operatorName := ""
	if u, uErr := s.userRepo.FindByID(operatorID); uErr == nil {
		operatorName = u.Name
	}
	sup := &model.IncidentSupervision{
		IncidentID:   incidentID,
		Note:         note,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		CreatedAt:    time.Now(),
	}
	if err := s.supervisionRepo.Create(sup); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			s.logger.Warn(constants.LogIncidentSuperviseConflict, "incident_id", incidentID, "operator_id", operatorID)
			return nil, util.NewAppError(constants.CodeIncidentSupervisionDuplicate,
				"SafetyIncident[id="+u64(incidentID)+"] supervise conflict: already supervised")
		}
		return nil, util.Wrap(err, "SafetyIncident[id=%d] supervise save failed", incidentID)
	}
	s.logger.Info(constants.LogIncidentSuperviseSuccess, "incident_id", incidentID, "operator_id", operatorID)
	return sup, nil
}

func u64(v uint64) string {
	if v == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	return string(b[i:])
}
