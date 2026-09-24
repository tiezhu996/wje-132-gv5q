package service

import (
	"testing"

	"safetyplatform/internal/constants"
)

func TestU64(t *testing.T) {
	if u64(0) != "0" || u64(123) != "123" {
		t.Error("u64 mismatch")
	}
}

func TestContains(t *testing.T) {
	if !contains(constants.UserRoleValues, constants.RoleInspector) {
		t.Error("inspector should be in roles")
	}
	if contains(constants.UserRoleValues, "bogus") {
		t.Error("bogus should not be in roles")
	}
}

func TestIncidentCategories(t *testing.T) {
	found := false
	for _, c := range constants.IncidentCategories {
		if c == "坠落" {
			found = true
		}
	}
	if !found {
		t.Error("坠落 should be in categories")
	}
}
