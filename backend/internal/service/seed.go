package service

import (
	"log/slog"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedService 启动时幂等写入预置数据。
type SeedService struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewSeedService 构造种子服务。
func NewSeedService(db *gorm.DB, logger *slog.Logger) *SeedService {
	return &SeedService{db: db, logger: logger}
}

// Seed 当 users 表为空时写入种子数据。
func (s *SeedService) Seed() error {
	var count int64
	if err := s.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("User@123"), bcrypt.DefaultCost)
	users := []model.User{
		{Phone: "13800000001", PasswordHash: string(adminHash), Name: "系统管理员", Role: constants.RoleAdmin},
		{Phone: "13800000002", PasswordHash: string(userHash), Name: "王安全", Role: constants.RoleSafetyManager},
		{Phone: "13800000003", PasswordHash: string(userHash), Name: "李监理", Role: constants.RoleInspector},
		{Phone: "13800000004", PasswordHash: string(userHash), Name: "赵工", Role: constants.RoleWorker},
	}
	for i := range users {
		if err := s.db.Create(&users[i]).Error; err != nil {
			return err
		}
	}
	incidents := []model.SafetyIncident{
		{Title: "脚手架扣件松动", Description: "三层东侧脚手架扣件松动，存在坠落风险。", OccurredAt: time.Now().AddDate(0, 0, -1), SiteID: "SITE-A", Area: "三层东侧", SeverityLevel: constants.SeverityMajor, Category: "坠落", InvolvedUserIDs: model.JSONList{"4"}, Status: constants.IncidentInvestigating, ReporterID: 3},
		{Title: "临时用电电缆破损", Description: "二级配电箱电缆绝缘层破损。", OccurredAt: time.Now().AddDate(0, 0, -2), SiteID: "SITE-A", Area: "加工区", SeverityLevel: constants.SeverityModerate, Category: "触电", InvolvedUserIDs: model.JSONList{"4"}, Status: constants.IncidentResolved, ReporterID: 3},
		{Title: "高处坠物未遂", Description: "塔吊吊运时构件滑落未造成伤害。", OccurredAt: time.Now().AddDate(0, 0, -5), SiteID: "SITE-A", Area: "吊装区", SeverityLevel: constants.SeverityMinor, Category: "物体打击", InvolvedUserIDs: model.JSONList{"4"}, Status: constants.IncidentClosed, ReporterID: 2},
	}
	for i := range incidents {
		if err := s.db.Create(&incidents[i]).Error; err != nil {
			return err
		}
	}
	inspections := []model.SafetyInspection{
		{Name: "8月例行安全检查", InspectionType: constants.InspectionRoutine, Area: "全工地", InspectionDate: time.Now().AddDate(0, 0, -1), InspectorID: 3, TotalScore: 92, Status: constants.InspectionCompleted, IssueCount: 1, PassedCount: 2},
		{Name: "高处作业专项检查", InspectionType: constants.InspectionSpecial, Area: "三层作业面", InspectionDate: time.Now().AddDate(0, 0, 1), InspectorID: 3, Status: constants.InspectionScheduled},
	}
	for i := range inspections {
		if err := s.db.Create(&inspections[i]).Error; err != nil {
			return err
		}
	}
	items := []model.InspectionItem{
		{InspectionID: 1, ItemName: "安全帽佩戴", Passed: true},
		{InspectionID: 1, ItemName: "临边防护栏杆", Passed: false, Remark: "东侧栏杆缺失"},
		{InspectionID: 1, ItemName: "消防器材齐全", Passed: true},
	}
	for i := range items {
		if err := s.db.Create(&items[i]).Error; err != nil {
			return err
		}
	}
	trainings := []model.SafetyTraining{
		{Topic: "新员工入场安全培训", TrainingType: constants.TrainingInduction, TrainingDate: time.Now().AddDate(0, 0, -3), DurationHours: 4, Trainer: "王安全", Location: "培训室A", ContentSummary: "入场安全须知与应急疏散。", ParticipantIDs: model.JSONList{"4"}, AssessmentMethod: "笔试", PassRate: 95},
		{Topic: "高空作业专项培训", TrainingType: constants.TrainingSpecial, TrainingDate: time.Now().AddDate(0, 0, 2), DurationHours: 2, Trainer: "王安全", Location: "培训室B", ContentSummary: "高空作业规范与防护用品使用。", AssessmentMethod: "实操"},
	}
	for i := range trainings {
		if err := s.db.Create(&trainings[i]).Error; err != nil {
			return err
		}
	}
	now := time.Now()
	certs := []model.WorkerCertification{
		{UserID: 4, CertType: "特种作业证", CertNo: "TZ20260001", IssueOrg: "市应急管理局", IssueDate: &now, ValidUntil: &now, Status: constants.CertApproved},
		{UserID: 3, CertType: "安全员证", CertNo: "AQ20260002", IssueOrg: "市住建局", IssueDate: &now, ValidUntil: &now, Status: constants.CertPending},
	}
	for i := range certs {
		if err := s.db.Create(&certs[i]).Error; err != nil {
			return err
		}
	}
	s.logger.Info("seed data created")
	return nil
}
