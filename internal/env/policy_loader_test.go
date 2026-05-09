package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writePolicyFile(t *testing.T, p *Policy) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	if err := SavePolicy(path, p); err != nil {
		t.Fatalf("SavePolicy: %v", err)
	}
	return path
}

func TestSaveAndLoadPolicy(t *testing.T) {
	orig := samplePolicy()
	path := writePolicyFile(t, orig)

	loaded, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("LoadPolicy: %v", err)
	}

	if loaded.Version != orig.Version {
		t.Errorf("version mismatch: got %s want %s", loaded.Version, orig.Version)
	}
	if len(loaded.Rules) != len(orig.Rules) {
		t.Errorf("rule count mismatch: got %d want %d", len(loaded.Rules), len(orig.Rules))
	}
}

func TestLoadPolicyNotFound(t *testing.T) {
	_, err := LoadPolicy("/nonexistent/policy.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadPolicyInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o600)

	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSavePolicyInvalidPolicy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	bad := &Policy{} // missing version
	if err := SavePolicy(path, bad); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCheckAccessConvenience(t *testing.T) {
	path := writePolicyFile(t, samplePolicy())

	res, err := CheckAccess(path, "admin", "DB_URL")
	if err != nil {
		t.Fatalf("CheckAccess: %v", err)
	}
	if !res.Allowed {
		t.Errorf("expected allow, got deny: %s", res.Reason)
	}
}

func TestCheckAccessDenyConvenience(t *testing.T) {
	path := writePolicyFile(t, samplePolicy())

	res, err := CheckAccess(path, "developer", "DB_PASS")
	if err != nil {
		t.Fatalf("CheckAccess: %v", err)
	}
	if res.Allowed {
		t.Errorf("expected deny, got allow")
	}
}
