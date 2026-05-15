package env

import (
	"testing"
)

// TestArchiveRoundtripWithParser verifies Save → reload → Get preserves all vars
// as produced by the env parser.
func TestArchiveRoundtripWithParser(t *testing.T) {
	raw := "DB_URL=postgres://localhost/dev\nSECRET_KEY=hunter2\nDEBUG=true\n"
	vars, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	path := tempArchivePath(t)
	s, err := NewArchiveStore(path)
	if err != nil {
		t.Fatalf("NewArchiveStore: %v", err)
	}
	if err := s.Save("dev-snapshot", vars, "integration test"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// reload from disk
	s2, err := NewArchiveStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := s2.Get("dev-snapshot")
	if !ok {
		t.Fatal("expected archive after reload")
	}
	for k, v := range vars {
		if e.Vars[k] != v {
			t.Errorf("key %s: want %q got %q", k, v, e.Vars[k])
		}
	}
}

// TestArchiveMultipleVersions stores two versions and confirms independent retrieval.
func TestArchiveMultipleVersions(t *testing.T) {
	path := tempArchivePath(t)
	s, _ := NewArchiveStore(path)

	v1 := map[string]string{"APP_ENV": "staging", "PORT": "8080"}
	v2 := map[string]string{"APP_ENV": "production", "PORT": "443", "WORKERS": "4"}

	_ = s.Save("v1", v1, "staging release")
	_ = s.Save("v2", v2, "production release")

	s2, _ := NewArchiveStore(path)

	e1, ok1 := s2.Get("v1")
	e2, ok2 := s2.Get("v2")
	if !ok1 || !ok2 {
		t.Fatal("expected both versions")
	}
	if e1.Vars["APP_ENV"] != "staging" {
		t.Errorf("v1 APP_ENV: want staging got %s", e1.Vars["APP_ENV"])
	}
	if e2.Vars["APP_ENV"] != "production" {
		t.Errorf("v2 APP_ENV: want production got %s", e2.Vars["APP_ENV"])
	}
	if _, exists := e1.Vars["WORKERS"]; exists {
		t.Error("v1 should not have WORKERS key")
	}
	names := s2.List()
	if len(names) != 2 {
		t.Errorf("expected 2 archives, got %d", len(names))
	}
}
