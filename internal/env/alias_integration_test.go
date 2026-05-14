package env

import (
	"testing"
)

// TestAliasRoundtripWithResolve verifies that aliases saved and reloaded
// correctly resolve to their canonical keys through the full store lifecycle.
func TestAliasRoundtripWithResolve(t *testing.T) {
	path := tempAliasPath(t)

	// Phase 1: write aliases.
	s, err := NewAliasStore(path)
	if err != nil {
		t.Fatal(err)
	}
	pairs := map[string]string{
		"db":    "DATABASE_URL",
		"redis": "REDIS_URL",
		"port":  "APP_PORT",
	}
	for name, key := range pairs {
		if err := s.Set(name, key); err != nil {
			t.Fatalf("Set(%q): %v", name, err)
		}
	}

	// Phase 2: reload and verify resolution.
	s2, err := NewAliasStore(path)
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range pairs {
		got := s2.Resolve(name)
		if got != want {
			t.Errorf("Resolve(%q) = %q, want %q", name, got, want)
		}
	}
}

// TestAliasDeleteAndReload ensures deletions are reflected after a store reload.
func TestAliasDeleteAndReload(t *testing.T) {
	path := tempAliasPath(t)

	s, _ := NewAliasStore(path)
	_ = s.Set("svc", "SERVICE_URL")
	_ = s.Set("key", "API_KEY")
	_ = s.Delete("svc")

	s2, err := NewAliasStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s2.Get("svc"); ok {
		t.Error("expected 'svc' alias to be absent after reload")
	}
	if k, ok := s2.Get("key"); !ok || k != "API_KEY" {
		t.Errorf("expected 'key' alias to remain, got %q ok=%v", k, ok)
	}
}
