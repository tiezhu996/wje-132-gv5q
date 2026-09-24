//go:build sqlite_integration

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"safetyplatform/internal/config"
	"safetyplatform/internal/constants"
	"safetyplatform/internal/middleware"
	"safetyplatform/internal/model"
	"safetyplatform/internal/repository"
	"safetyplatform/internal/service"
	"safetyplatform/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIncidentRouter(t *testing.T) (*gin.Engine, *gorm.DB, *SafetyIncidentHandler) {
	t.Helper()
	gin.SetMode(gin.ReleaseMode)
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s_http?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.SafetyIncident{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})

	logger := slog.Default()
	incidentRepo := repository.NewSafetyIncidentRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := service.NewSafetyIncidentService(incidentRepo, userRepo, logger)
	h := NewSafetyIncidentHandler(svc, logger)
	cfg := &config.Config{JWTSecret: "test-secret"}

	r := gin.New()
	g := r.Group("/api/v1")
	incidents := g.Group("/incidents")
	incidents.Use(middleware.AuthRequired(cfg))
	incidents.GET("", h.List)
	incidents.GET("/overdue-acceptance", h.OverdueAcceptance)
	incidents.GET("/:id", h.Get)
	incidents.POST("/:id/supervise", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), h.Supervise)
	incidents.POST("/:id/close", middleware.RequireRole(constants.RoleAdmin, constants.RoleSafetyManager), h.Close)
	return r, db, h
}

func token(t *testing.T, userID uint64, role string) string {
	t.Helper()
	tok, err := util.GenerateToken("test-secret", time.Hour, userID, "13800000000", role)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return tok
}

func doJSON(t *testing.T, r *gin.Engine, method, path, tokenStr, body string) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if body != "" {
		buf = bytes.NewBufferString(body)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	if tokenStr != "" {
		req.Header.Set("Authorization", "Bearer "+tokenStr)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type respBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestOverdueAcceptanceAndSuperviseHTTP(t *testing.T) {
	r, db, _ := setupIncidentRouter(t)
	past := time.Now().Add(-3 * time.Hour)
	future := time.Now().Add(3 * time.Hour)
	inc := model.SafetyIncident{
		Title: "电缆破损", Status: constants.IncidentResolved, SeverityLevel: constants.SeverityMajor,
		RectificationDeadline: &past, ReporterID: 1, OccurredAt: time.Now().Add(-time.Hour),
		InvolvedUserIDs: model.JSONList{}, PhotoURLs: model.JSONList{},
	}
	if err := db.Create(&inc).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	// 不进队列的对照记录
	db.Create(&model.SafetyIncident{Title: "未到期", Status: constants.IncidentResolved, SeverityLevel: constants.SeverityFatal,
		RectificationDeadline: &future, ReporterID: 1, OccurredAt: time.Now(),
		InvolvedUserIDs: model.JSONList{}, PhotoURLs: model.JSONList{}})

	// 未登录 -> 401
	w := doJSON(t, r, http.MethodGet, "/api/v1/incidents/overdue-acceptance", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no token status = %d", w.Code)
	}

	// 工人也可以看队列（仅写操作受限）
	w = doJSON(t, r, http.MethodGet, "/api/v1/incidents/overdue-acceptance", token(t, 4, constants.RoleWorker), "")
	if w.Code != http.StatusOK {
		t.Fatalf("worker list status = %d body=%s", w.Code, w.Body.String())
	}
	var rb struct {
		Data struct {
			List []map[string]any `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &rb); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(rb.Data.List) != 1 {
		t.Fatalf("queue len = %d, want 1", len(rb.Data.List))
	}
	if rb.Data.List[0]["id"].(float64) != float64(inc.ID) {
		t.Fatalf("queue item id mismatch")
	}
	if dur, _ := rb.Data.List[0]["overdue_duration"].(float64); dur < 3*3600-60 {
		t.Fatalf("overdue_duration = %v", dur)
	}

	// 工人督办 -> 403
	w = doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/supervise", inc.ID),
		token(t, 4, constants.RoleWorker), `{"note":"快点"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("worker supervise status = %d", w.Code)
	}

	// 缺少说明 -> 400
	w = doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/supervise", inc.ID),
		token(t, 1, constants.RoleAdmin), `{"note":""}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty note status = %d", w.Code)
	}

	// 管理员督办成功
	w = doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/supervise", inc.ID),
		token(t, 1, constants.RoleAdmin), `{"note":"请三日内完成验收"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("supervise status = %d body=%s", w.Code, w.Body.String())
	}
	var okResp respBody
	json.Unmarshal(w.Body.Bytes(), &okResp)
	if okResp.Message != constants.MsgSupervisionSent {
		t.Fatalf("message = %q", okResp.Message)
	}

	// 队列里现在能看到督办说明与时间
	w = doJSON(t, r, http.MethodGet, "/api/v1/incidents/overdue-acceptance", token(t, 1, constants.RoleAdmin), "")
	json.Unmarshal(w.Body.Bytes(), &rb)
	if rb.Data.List[0]["supervision_note"] != "请三日内完成验收" {
		t.Fatalf("supervision_note missing in queue: %v", rb.Data.List[0])
	}
	if rb.Data.List[0]["supervision_at"] == nil {
		t.Fatal("supervision_at missing in queue")
	}

	// 重复督办 -> 409 + 40902
	w = doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/supervise", inc.ID),
		token(t, 1, constants.RoleAdmin), `{"note":"再来一次"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate supervise status = %d", w.Code)
	}
	var conflict respBody
	json.Unmarshal(w.Body.Bytes(), &conflict)
	if conflict.Code != constants.CodeSupervisionConflict {
		t.Fatalf("conflict code = %d, want %d", conflict.Code, constants.CodeSupervisionConflict)
	}

	// 关闭后记录离开队列，且再督办报状态冲突
	w = doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/close", inc.ID),
		token(t, 1, constants.RoleAdmin), "")
	if w.Code != http.StatusOK {
		t.Fatalf("close status = %d body=%s", w.Code, w.Body.String())
	}
	w = doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/supervise", inc.ID),
		token(t, 1, constants.RoleAdmin), `{"note":"关闭后督办"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("supervise after close status = %d", w.Code)
	}
	w = doJSON(t, r, http.MethodGet, "/api/v1/incidents/overdue-acceptance", token(t, 1, constants.RoleAdmin), "")
	json.Unmarshal(w.Body.Bytes(), &rb)
	if len(rb.Data.List) != 0 {
		t.Fatalf("queue after close len = %d, want 0", len(rb.Data.List))
	}
}

// TestSuperviseNotDueHTTP 期限未到的已整改记录不能督办。
func TestSuperviseNotDueHTTP(t *testing.T) {
	r, db, _ := setupIncidentRouter(t)
	future := time.Now().Add(24 * time.Hour)
	inc := model.SafetyIncident{
		Title: "未到期", Status: constants.IncidentResolved, SeverityLevel: constants.SeverityMajor,
		RectificationDeadline: &future, ReporterID: 1, OccurredAt: time.Now(),
		InvolvedUserIDs: model.JSONList{}, PhotoURLs: model.JSONList{},
	}
	db.Create(&inc)
	w := doJSON(t, r, http.MethodPost, fmt.Sprintf("/api/v1/incidents/%d/supervise", inc.ID),
		token(t, 2, constants.RoleSafetyManager), `{"note":"提前督办"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("supervise not-due status = %d body=%s", w.Code, w.Body.String())
	}
}
