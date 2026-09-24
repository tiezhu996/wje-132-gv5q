package util

import (
	"testing"
	"time"

	"safetyplatform/internal/constants"
)

func TestGenerateTokenAndParse(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(secret, time.Hour, 5, "13800000001", constants.RoleAdmin)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != 5 || claims.Phone != "13800000001" || claims.Role != constants.RoleAdmin {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	for _, tc := range []string{"", "x", "a.b.c"} {
		if _, err := ParseToken("s", tc); err == nil {
			t.Errorf("expected error for %q", tc)
		}
	}
}

func TestFormatters(t *testing.T) {
	if SeverityText(constants.SeverityMajor) != "较大" {
		t.Error("severity text mismatch")
	}
	if IncidentStatusText(constants.IncidentClosed) != "已关闭" {
		t.Error("incident status text mismatch")
	}
	if InspectionStatusText(constants.InspectionCompleted) != "已完成" {
		t.Error("inspection status text mismatch")
	}
	if CertStatusText(constants.CertPending) != "待审核" {
		t.Error("cert status text mismatch")
	}
	if UserRoleText(constants.RoleSafetyManager) != "安全管理员" {
		t.Error("role text mismatch")
	}
}

func TestSeverityValidators(t *testing.T) {
	if !constants.IsValidSeverity(constants.SeverityFatal) {
		t.Error("fatal should be valid")
	}
	if constants.IsValidSeverity("bogus") {
		t.Error("bogus should be invalid")
	}
	if !constants.IsValidIncidentStatus(constants.IncidentResolved) {
		t.Error("resolved should be valid")
	}
}
