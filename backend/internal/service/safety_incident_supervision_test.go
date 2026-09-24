package service

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"safetyplatform/internal/constants"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDBCounter int64

func newTestIncidentService(t *testing.T) (*SafetyIncidentService, *gorm.DB) {
	t.Helper()
	dsn := fmt.Sprintf("file:incident_test_%d?mode=memory&cache=shared", atomic.AddInt64(&testDBCounter, 1))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.SafetyIncident{}, &model.IncidentSupervision{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	incidentRepo := repository.NewSafetyIncidentRepository(db)
	userRepo := repository.NewUserRepository(db)
	supervisionRepo := repository.NewIncidentSupervisionRepository(db)
	return NewSafetyIncidentService(incidentRepo, userRepo, supervisionRepo, logger), db
}

func seedIncident(t *testing.T, db *gorm.DB, status string, deadline *time.Time, severity string, occurredOffsetDays int) uint64 {
	t.Helper()
	i := &model.SafetyIncident{
		Title: "t-" + status + "-" + severity, OccurredAt: time.Now().AddDate(0, 0, occurredOffsetDays),
		SeverityLevel: severity, Status: status, RectificationDeadline: deadline, ReporterID: 1,
	}
	if err := db.Create(i).Error; err != nil {
		t.Fatalf("seed incident: %v", err)
	}
	return i.ID
}

func TestOverdueAcceptanceQueue(t *testing.T) {
	svc, db := newTestIncidentService(t)
	now := time.Now()
	past := now.Add(-48 * time.Hour)
	older := now.Add(-72 * time.Hour)
	future := now.Add(24 * time.Hour)

	// 已整改逾期（moderate，逾期 2 天）
	idOverdueModerate := seedIncident(t, db, constants.IncidentResolved, &past, constants.SeverityModerate, -3)
	// 已整改逾期且更早、更严重（major，逾期 3 天）—— 应排第一
	idOverdueMajor := seedIncident(t, db, constants.IncidentResolved, &older, constants.SeverityMajor, -4)
	// 已整改但期限未到 —— 不能进队列
	seedIncident(t, db, constants.IncidentResolved, &future, constants.SeverityFatal, -1)
	// 已关闭且逾期 —— 不能进队列
	seedIncident(t, db, constants.IncidentClosed, &past, constants.SeverityFatal, -5)
	// 已整改但无期限 —— 不能进队列
	seedIncident(t, db, constants.IncidentResolved, nil, constants.SeverityMajor, -2)

	items, err := svc.OverdueAcceptanceQueue()
	if err != nil {
		t.Fatalf("queue: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 overdue items, got %d", len(items))
	}
	if items[0].ID != idOverdueMajor {
		t.Fatalf("expected longest overdue major first, got id=%d", items[0].ID)
	}
	if items[1].ID != idOverdueModerate {
		t.Fatalf("expected moderate second, got id=%d", items[1].ID)
	}
	if items[0].OverdueSeconds <= items[1].OverdueSeconds {
		t.Fatalf("overdue seconds must be descending: %d vs %d", items[0].OverdueSeconds, items[1].OverdueSeconds)
	}
	if items[0].Supervision != nil {
		t.Fatalf("expected no supervision initially")
	}
}

func TestOverdueQueueSeverityTieBreak(t *testing.T) {
	svc, db := newTestIncidentService(t)
	// 同一期限：风险高者排前
	deadline := time.Now().Add(-24 * time.Hour)
	idMinor := seedIncident(t, db, constants.IncidentResolved, &deadline, constants.SeverityMinor, -2)
	idFatal := seedIncident(t, db, constants.IncidentResolved, &deadline, constants.SeverityFatal, -2)
	_ = idMinor

	items, err := svc.OverdueAcceptanceQueue()
	if err != nil {
		t.Fatalf("queue: %v", err)
	}
	if len(items) != 2 || items[0].ID != idFatal {
		t.Fatalf("expected fatal first on tie, got %+v", items)
	}
}

func TestSuperviseSuccessAndConflict(t *testing.T) {
	svc, db := newTestIncidentService(t)
	past := time.Now().Add(-24 * time.Hour)
	id := seedIncident(t, db, constants.IncidentResolved, &past, constants.SeverityMajor, -2)

	sup, err := svc.Supervise(id, 1, "请立即组织验收")
	if err != nil {
		t.Fatalf("first supervise: %v", err)
	}
	if sup.Note != "请立即组织验收" {
		t.Fatalf("unexpected note: %s", sup.Note)
	}

	// 队列条目应带督办信息
	items, err := svc.OverdueAcceptanceQueue()
	if err != nil {
		t.Fatalf("queue: %v", err)
	}
	if len(items) != 1 || items[0].Supervision == nil || items[0].Supervision.Note != "请立即组织验收" {
		t.Fatalf("expected supervision attached to queue item: %+v", items)
	}

	// 重复督办 → 冲突
	_, err = svc.Supervise(id, 1, "再次督办")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeIncidentSupervisionDuplicate {
		t.Fatalf("expected duplicate conflict, got %v", err)
	}
}

func TestSuperviseGuardClosedAndNotDue(t *testing.T) {
	svc, db := newTestIncidentService(t)
	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	idClosed := seedIncident(t, db, constants.IncidentClosed, &past, constants.SeverityMajor, -5)
	if _, err := svc.Supervise(idClosed, 1, "x"); err == nil {
		t.Fatal("closed incident must not be supervisable")
	} else if appErr, ok := err.(*util.AppError); !ok || appErr.Code != constants.CodeIncidentStatusConflict {
		t.Fatalf("expected status conflict for closed, got %v", err)
	}

	idFuture := seedIncident(t, db, constants.IncidentResolved, &future, constants.SeverityMajor, -1)
	if _, err := svc.Supervise(idFuture, 1, "x"); err == nil {
		t.Fatal("not-due incident must not be supervisable")
	} else if appErr, ok := err.(*util.AppError); !ok || appErr.Code != constants.CodeIncidentNotOverdue {
		t.Fatalf("expected not-overdue conflict, got %v", err)
	}

	idNoDeadline := seedIncident(t, db, constants.IncidentResolved, nil, constants.SeverityMajor, -1)
	if _, err := svc.Supervise(idNoDeadline, 1, "x"); err == nil {
		t.Fatal("incident without deadline must not be supervisable")
	}

	if _, err := svc.Supervise(99999, 1, "x"); err == nil {
		t.Fatal("missing incident must error")
	}
}
