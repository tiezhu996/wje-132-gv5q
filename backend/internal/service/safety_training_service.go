package service

import (
	"log/slog"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"
)

// SafetyTrainingService 安全培训业务逻辑。
type SafetyTrainingService struct {
	repo     *repository.SafetyTrainingRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewSafetyTrainingService 构造安全培训服务。
func NewSafetyTrainingService(repo *repository.SafetyTrainingRepository, userRepo *repository.UserRepository, logger *slog.Logger) *SafetyTrainingService {
	return &SafetyTrainingService{repo: repo, userRepo: userRepo, logger: logger}
}

// Create 创建培训。
func (s *SafetyTrainingService) Create(topic, trainingType string, trainingDate time.Time, durationHours int,
	trainer, location, contentSummary string, participantIDs []string, assessmentMethod string) (*model.SafetyTraining, error) {
	t := &model.SafetyTraining{
		Topic: topic, TrainingType: trainingType, TrainingDate: trainingDate, DurationHours: durationHours,
		Trainer: trainer, Location: location, ContentSummary: contentSummary,
		ParticipantIDs: model.JSONList(participantIDs), AssessmentMethod: assessmentMethod,
	}
	if err := s.repo.Create(t); err != nil {
		return nil, util.Wrap(err, "SafetyTraining[topic=%s] create failed", topic)
	}
	s.logger.Info(constants.LogTrainingCreateSuccess, "training_id", t.ID)
	return t, nil
}

// Record 记录培训成绩（签到 + 通过率）。
func (s *SafetyTrainingService) Record(id uint64, participantIDs []string, passRate float64) (*model.SafetyTraining, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "SafetyTraining[id=%d] record find failed", id)
	}
	if participantIDs != nil {
		t.ParticipantIDs = model.JSONList(participantIDs)
	}
	if passRate >= 0 {
		t.PassRate = passRate
	}
	if err := s.repo.Update(t); err != nil {
		return nil, util.Wrap(err, "SafetyTraining[id=%d] record save failed", id)
	}
	s.logger.Info(constants.LogTrainingRecordSuccess, "training_id", t.ID, "pass_rate", passRate)
	return t, nil
}

// List 分页查询培训。
func (s *SafetyTrainingService) List(page, pageSize int, trainingType string) ([]model.SafetyTraining, int64, error) {
	return s.repo.List(page, pageSize, trainingType)
}

// Get 培训详情。
func (s *SafetyTrainingService) Get(id uint64) (*model.SafetyTraining, error) {
	return s.repo.FindByID(id)
}

// CompletedRate 本月培训完成率。
func (s *SafetyTrainingService) CompletedRate() (map[string]float64, error) {
	return s.repo.CompletedRate()
}
