//go:build sqlite_integration

package service

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.SafetyIncident{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	return db
}

func newIncidentService(db *gorm.DB) *SafetyIncidentService {
	return NewSafetyIncidentService(
		repository.NewSafetyIncidentRepository(db),
		repository.NewUserRepository(db),
		slog.Default(),
	)
}

func mkIncident(t *testing.T, db *gorm.DB, status string, deadline *time.Time, supervisedBy uint64) model.SafetyIncident {
	t.Helper()
	i := model.SafetyIncident{
		Title: "t", Status: status, SeverityLevel: constants.SeverityMajor,
		RectificationDeadline: deadline, SupervisedBy: supervisedBy, ReporterID: 1,
		OccurredAt: time.Now().Add(-time.Hour),
	}
	if err := db.Create(&i).Error; err != nil {
		t.Fatalf("create incident: %v", err)
	}
	return i
}

func TestOverdueAcceptanceQueue(t *testing.T) {
	db := newTestDB(t)
	svc := newIncidentService(db)
	past := time.Now().Add(-48 * time.Hour)
	past2 := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	overdueMajor := mkIncident(t, db, constants.IncidentResolved, &past, 0)  // 逾期 48h major
	overdueFatal := mkIncident(t, db, constants.IncidentResolved, &past2, 0) // 逾期 24h fatal
	db.Model(&model.SafetyIncident{}).Where("id = ?", overdueFatal.ID).
		Update("severity_level", constants.SeverityFatal)
	_ = mkIncident(t, db, constants.IncidentClosed, &past, 0)              // 已关闭 -> 不进队列
	_ = mkIncident(t, db, constants.IncidentResolved, &future, 0)          // 期限未到 -> 不进队列
	_ = mkIncident(t, db, constants.IncidentResolved, nil, 0)              // 无期限 -> 不进队列
	_ = mkIncident(t, db, constants.IncidentInvestigating, &past, 0)       // 非已整改 -> 不进队列
	supervised := mkIncident(t, db, constants.IncidentResolved, &past2, 9) // 逾期 24h，已督办，仍在队列
	note, sat := "请立即整改", past2.Add(time.Hour)
	db.Model(&model.SafetyIncident{}).Where("id = ?", supervised.ID).
		Updates(map[string]any{"supervision_note": note, "supervision_at": sat})

	list, err := svc.OverdueAcceptance()
	if err != nil {
		t.Fatalf("OverdueAcceptance: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("queue size = %d, want 3", len(list))
	}
	// 第一排序键：逾期时长从长到短（deadline ASC）
	if list[0].ID != overdueMajor.ID {
		t.Fatalf("first = %d, want longest-overdue %d", list[0].ID, overdueMajor.ID)
	}
	// 同为 24h 时风险等级高者（fatal）排在已督办的记录之前
	if list[1].ID != overdueFatal.ID || list[2].ID != supervised.ID {
		t.Fatalf("order = [%d %d %d], want [%d %d %d]",
			list[0].ID, list[1].ID, list[2].ID, overdueMajor.ID, overdueFatal.ID, supervised.ID)
	}
	if list[0].OverdueDuration < 48*3600-60 {
		t.Fatalf("overdue duration too small: %d", list[0].OverdueDuration)
	}
	var foundSupervised bool
	for _, item := range list {
		if item.ID == supervised.ID {
			foundSupervised = true
			if item.SupervisionNote != note {
				t.Fatalf("supervision note = %q", item.SupervisionNote)
			}
		}
	}
	if !foundSupervised {
		t.Fatal("supervised record should remain in queue")
	}
}

func TestSuperviseOnceAndConflicts(t *testing.T) {
	db := newTestDB(t)
	svc := newIncidentService(db)
	past := time.Now().Add(-2 * time.Hour)
	future := time.Now().Add(2 * time.Hour)

	target := mkIncident(t, db, constants.IncidentResolved, &past, 0)
	closed := mkIncident(t, db, constants.IncidentClosed, &past, 0)
	notDue := mkIncident(t, db, constants.IncidentResolved, &future, 0)
	investigating := mkIncident(t, db, constants.IncidentInvestigating, &past, 0)
	noDeadline := mkIncident(t, db, constants.IncidentResolved, nil, 0)

	updated, err := svc.Supervise(target.ID, 7, "请尽快验收")
	if err != nil {
		t.Fatalf("first supervise: %v", err)
	}
	if updated.SupervisionNote != "请尽快验收" || updated.SupervisedBy != 7 || updated.SupervisionAt == nil {
		t.Fatalf("supervision fields not saved: %+v", updated)
	}

	// 重复督办 -> 40902
	_, err = svc.Supervise(target.ID, 7, "再次督办")
	if appErr, ok := err.(*util.AppError); !ok || appErr.Code != constants.CodeSupervisionConflict {
		t.Fatalf("duplicate supervise err = %v, want CodeSupervisionConflict", err)
	}

	// 已关闭 / 调查中 -> 状态冲突 40901
	for _, id := range []uint64{closed.ID, investigating.ID} {
		_, err = svc.Supervise(id, 7, "x")
		if appErr, ok := err.(*util.AppError); !ok || appErr.Code != constants.CodeIncidentStatusConflict {
			t.Fatalf("id=%d err = %v, want status conflict", id, err)
		}
	}
	// 期限未到 / 无期限 -> 状态冲突 40901
	for _, id := range []uint64{notDue.ID, noDeadline.ID} {
		_, err = svc.Supervise(id, 7, "x")
		if appErr, ok := err.(*util.AppError); !ok || appErr.Code != constants.CodeIncidentStatusConflict {
			t.Fatalf("id=%d err = %v, want not-due conflict", id, err)
		}
	}

	// 督办后事件仍可关闭（原有操作继续使用）
	if _, err := svc.Close(target.ID); err != nil {
		t.Fatalf("close after supervise: %v", err)
	}
	got, _ := svc.Get(target.ID)
	if got.Status != constants.IncidentClosed {
		t.Fatalf("status = %s, want closed", got.Status)
	}
}

func TestSuperviseConcurrentOnlyOneWins(t *testing.T) {
	db := newTestDB(t)
	svc := newIncidentService(db)
	past := time.Now().Add(-time.Hour)
	target := mkIncident(t, db, constants.IncidentResolved, &past, 0)

	// 模拟另一个请求已督办
	rows, err := repository.NewSafetyIncidentRepository(db).Supervise(target.ID, 1, "race", time.Now())
	if err != nil || rows != 1 {
		t.Fatalf("seed supervise rows=%d err=%v", rows, err)
	}
	// 服务层先 FindByID 看到 supervised_by=1，应直接报重复督办
	_, err = svc.Supervise(target.ID, 2, "late")
	if appErr, ok := err.(*util.AppError); !ok || appErr.Code != constants.CodeSupervisionConflict {
		t.Fatalf("concurrent supervise err = %v, want conflict", err)
	}
}
