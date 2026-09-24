package service

import (
	"log/slog"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"
)

// WorkerCertificationService 人员资质业务逻辑。
type WorkerCertificationService struct {
	repo     *repository.WorkerCertificationRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewWorkerCertificationService 构造人员资质服务。
func NewWorkerCertificationService(repo *repository.WorkerCertificationRepository, userRepo *repository.UserRepository, logger *slog.Logger) *WorkerCertificationService {
	return &WorkerCertificationService{repo: repo, userRepo: userRepo, logger: logger}
}

// Submit 提交资质。
func (s *WorkerCertificationService) Submit(userID uint64, certType, certNo, issueOrg string, issueDate, validUntil *time.Time, certPhotoURL string) (*model.WorkerCertification, error) {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return nil, util.Wrap(err, "WorkerCertification[user_id=%d] submit: user not found", userID)
	}
	c := &model.WorkerCertification{
		UserID: userID, CertType: certType, CertNo: certNo, IssueOrg: issueOrg,
		IssueDate: issueDate, ValidUntil: validUntil, CertPhotoURL: certPhotoURL,
		Status: constants.CertPending,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, util.Wrap(err, "WorkerCertification[user_id=%d] submit create failed", userID)
	}
	s.logger.Info(constants.LogCertSubmitSuccess, "cert_id", c.ID)
	return c, nil
}

// Review 审核资质。
func (s *WorkerCertificationService) Review(id uint64, status string) (*model.WorkerCertification, error) {
	if status != constants.CertApproved && status != constants.CertRejected {
		return nil, util.NewAppError(constants.CodeValidationFailed, "WorkerCertification[id="+u64(id)+"] review invalid status="+status)
	}
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "WorkerCertification[id=%d] review find failed", id)
	}
	if c.Status != constants.CertPending {
		return nil, util.NewAppError(constants.CodeIncidentStatusConflict, "WorkerCertification[id="+u64(id)+"] review conflict: status="+c.Status)
	}
	c.Status = status
	if err := s.repo.Update(c); err != nil {
		s.logger.Error(constants.LogCertReviewFailed, "error", err.Error())
		return nil, util.Wrap(err, "WorkerCertification[id=%d] review save failed", id)
	}
	s.logger.Info(constants.LogCertReviewSuccess, "cert_id", c.ID, "status", status)
	return c, nil
}

// List 分页查询资质。
func (s *WorkerCertificationService) List(page, pageSize int, status string, expiring bool) ([]model.WorkerCertification, int64, error) {
	return s.repo.List(page, pageSize, status, expiring)
}

// ListByUser 查询某用户资质。
func (s *WorkerCertificationService) ListByUser(userID uint64) ([]model.WorkerCertification, error) {
	return s.repo.ListByUser(userID)
}

// ExpiringSoon 即将过期资质。
func (s *WorkerCertificationService) ExpiringSoon() ([]model.WorkerCertification, error) {
	list, err := s.repo.ExpiringSoon()
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogCertExpiryWarning, "count", len(list))
	return list, nil
}
