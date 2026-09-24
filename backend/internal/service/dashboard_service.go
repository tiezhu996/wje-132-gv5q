package service

import (
	"log/slog"

	"safetyplatform/internal/constants"
)

// DashboardService 仪表盘统计业务逻辑。
type DashboardService struct {
	incidentSvc   *SafetyIncidentService
	inspectionSvc *SafetyInspectionService
	trainingSvc   *SafetyTrainingService
	certSvc       *WorkerCertificationService
	logger        *slog.Logger
}

// NewDashboardService 构造仪表盘服务。
func NewDashboardService(incidentSvc *SafetyIncidentService, inspectionSvc *SafetyInspectionService,
	trainingSvc *SafetyTrainingService, certSvc *WorkerCertificationService, logger *slog.Logger) *DashboardService {
	return &DashboardService{incidentSvc: incidentSvc, inspectionSvc: inspectionSvc, trainingSvc: trainingSvc, certSvc: certSvc, logger: logger}
}

// Stats 汇总仪表盘数据。
func (s *DashboardService) Stats() (map[string]any, error) {
	trend, err := s.incidentSvc.Trend30()
	if err != nil {
		return nil, err
	}
	distribution, err := s.incidentSvc.SeverityDistribution()
	if err != nil {
		return nil, err
	}
	pending, err := s.incidentSvc.PendingRectification()
	if err != nil {
		return nil, err
	}
	inspectionStats, err := s.inspectionSvc.Stats()
	if err != nil {
		return nil, err
	}
	trainingRate, err := s.trainingSvc.CompletedRate()
	if err != nil {
		return nil, err
	}
	expiring, err := s.certSvc.ExpiringSoon()
	if err != nil {
		return nil, err
	}
	result := map[string]any{
		"trend":                   trend,
		"severity_distribution":   distribution,
		"pending_rectification":   pending,
		"inspection":              inspectionStats,
		"training_completed_rate": trainingRate["rate"],
		"expiring_certs":          expiring,
	}
	s.logger.Info(constants.LogDashboardStats, "data", result)
	return result, nil
}
