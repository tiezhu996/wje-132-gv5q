package service

import (
	"log/slog"
	"time"

	"safetyplatform/internal/constants"
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
