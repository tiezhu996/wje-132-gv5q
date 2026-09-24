package service

import (
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
	repo     *repository.SafetyIncidentRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewSafetyIncidentService 构造安全事件服务。
func NewSafetyIncidentService(repo *repository.SafetyIncidentRepository, userRepo *repository.UserRepository, logger *slog.Logger) *SafetyIncidentService {
	return &SafetyIncidentService{repo: repo, userRepo: userRepo, logger: logger}
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

// OverdueAcceptance 逾期验收队列：仅包含已整改且期限早于当前时刻的事件。
func (s *SafetyIncidentService) OverdueAcceptance() ([]dto.OverdueIncident, error) {
	now := time.Now()
	list, err := s.repo.OverdueAcceptance(now)
	if err != nil {
		return nil, err
	}
	result := make([]dto.OverdueIncident, 0, len(list))
	for i := range list {
		var overdue int64
		if list[i].RectificationDeadline != nil {
			overdue = int64(now.Sub(*list[i].RectificationDeadline).Seconds())
			if overdue < 0 {
				overdue = 0
			}
		}
		result = append(result, dto.OverdueIncident{SafetyIncident: list[i], OverdueDuration: overdue})
	}
	return result, nil
}

// Supervise 对一条已整改且逾期未验收的事件发起一次督办并填写说明。
// 已关闭/期限未到/非已整改的记录不能督办；已督办过的记录返回冲突。
func (s *SafetyIncidentService) Supervise(id, operatorID uint64, note string) (*model.SafetyIncident, error) {
	i, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] supervise find failed", id)
	}
	if i.Status != constants.IncidentResolved {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyIncident[id="+u64(id)+"] supervise conflict: status="+i.Status)
	}
	if i.RectificationDeadline == nil || !i.RectificationDeadline.Before(time.Now()) {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyIncident[id="+u64(id)+"] supervise conflict: deadline not reached")
	}
	if i.SupervisedBy != 0 {
		return nil, util.NewAppError(constants.CodeSupervisionConflict, "SafetyIncident[id="+u64(id)+"] supervise conflict: already supervised")
	}
	rows, err := s.repo.Supervise(id, operatorID, note, time.Now())
	if err != nil {
		return nil, util.Wrap(err, "SafetyIncident[id=%d] supervise save failed", id)
	}
	if rows == 0 {
		// 并发场景下可能已被督办或状态已变更，重新加载确认冲突原因。
		if latest, findErr := s.repo.FindByID(id); findErr == nil && latest.SupervisedBy != 0 {
			s.logger.Warn(constants.LogIncidentSuperviseConflict, "incident_id", id)
			return nil, util.NewAppError(constants.CodeSupervisionConflict, "SafetyIncident[id="+u64(id)+"] supervise conflict: already supervised")
		}
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "SafetyIncident[id="+u64(id)+"] supervise conflict: state changed")
	}
	s.logger.Info(constants.LogIncidentSuperviseSuccess, "incident_id", id, "operator_id", operatorID)
	return s.repo.FindByID(id)
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
