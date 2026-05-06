package env

import (
	"testing"
)

func TestAccessSetAndGetPolicy(t *testing.T) {
	s := NewAccessStore()
	p := AccessPolicy{Role: "dev", AllowRead: []string{"*"}, AllowWrite: []string{"APP_"}}
	if err := s.SetPolicy(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.GetPolicy("dev")
	if !ok {
		t.Fatal("expected policy to exist")
	}
	if got.Role != "dev" {
		t.Errorf("got role %q, want %q", got.Role, "dev")
	}
}

func TestAccessSetPolicyEmptyRole(t *testing.T) {
	s := NewAccessStore()
	err := s.SetPolicy(AccessPolicy{Role: ""})
	if err == nil {
		t.Fatal("expected error for empty role")
	}
}

func TestAccessDeletePolicy(t *testing.T) {
	s := NewAccessStore()
	_ = s.SetPolicy(AccessPolicy{Role: "ops", AllowRead: []string{"*"}})
	s.DeletePolicy("ops")
	_, ok := s.GetPolicy("ops")
	if ok {
		t.Fatal("expected policy to be deleted")
	}
}

func TestAccessRoles(t *testing.T) {
	s := NewAccessStore()
	_ = s.SetPolicy(AccessPolicy{Role: "beta"})
	_ = s.SetPolicy(AccessPolicy{Role: "alpha"})
	roles := s.Roles()
	if len(roles) != 2 || roles[0] != "alpha" || roles[1] != "beta" {
		t.Errorf("unexpected roles: %v", roles)
	}
}

func TestCanReadWildcard(t *testing.T) {
	s := NewAccessStore()
	_ = s.SetPolicy(AccessPolicy{Role: "admin", AllowRead: []string{"*"}})
	if !s.CanRead("admin", "SECRET_KEY") {
		t.Error("admin should be able to read any key")
	}
}

func TestCanReadPrefix(t *testing.T) {
	s := NewAccessStore()
	_ = s.SetPolicy(AccessPolicy{Role: "dev", AllowRead: []string{"APP_"}})
	if !s.CanRead("dev", "APP_PORT") {
		t.Error("dev should read APP_ keys")
	}
	if s.CanRead("dev", "DB_PASSWORD") {
		t.Error("dev should NOT read DB_ keys")
	}
}

func TestCanWritePrefix(t *testing.T) {
	s := NewAccessStore()
	_ = s.SetPolicy(AccessPolicy{Role: "ci", AllowRead: []string{"*"}, AllowWrite: []string{"CI_"}})
	if !s.CanWrite("ci", "CI_TOKEN") {
		t.Error("ci should write CI_ keys")
	}
	if s.CanWrite("ci", "APP_SECRET") {
		t.Error("ci should NOT write APP_ keys")
	}
}

func TestCanReadUnknownRole(t *testing.T) {
	s := NewAccessStore()
	if s.CanRead("ghost", "ANY_KEY") {
		t.Error("unknown role should not have read access")
	}
}

func TestFilterReadable(t *testing.T) {
	s := NewAccessStore()
	_ = s.SetPolicy(AccessPolicy{Role: "dev", AllowRead: []string{"APP_"}})
	vars := map[string]string{
		"APP_PORT":    "8080",
		"DB_PASSWORD": "secret",
		"APP_NAME":    "envault",
	}
	filtered := s.FilterReadable("dev", vars)
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered keys, got %d", len(filtered))
	}
	if _, ok := filtered["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be filtered out")
	}
}
