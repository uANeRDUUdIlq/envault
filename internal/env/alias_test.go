package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempAliasPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "aliases.json")
}

func TestAliasSetAndGet(t *testing.T) {
	s, err := NewAliasStore(tempAliasPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("db", "DATABASE_URL"); err != nil {
		t.Fatal(err)
	}
	key, ok := s.Get("db")
	if !ok || key != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL, got %q (ok=%v)", key, ok)
	}
}

func TestAliasPersists(t *testing.T) {
	path := tempAliasPath(t)
	s, _ := NewAliasStore(path)
	_ = s.Set("redis", "REDIS_URL")

	s2, err := NewAliasStore(path)
	if err != nil {
		t.Fatal(err)
	}
	key, ok := s2.Get("redis")
	if !ok || key != "REDIS_URL" {
		t.Fatalf("expected REDIS_URL after reload, got %q", key)
	}
}

func TestAliasDelete(t *testing.T) {
	s, _ := NewAliasStore(tempAliasPath(t))
	_ = s.Set("pg", "POSTGRES_URL")
	if err := s.Delete("pg"); err != nil {
		t.Fatal(err)
	}
	_, ok := s.Get("pg")
	if ok {
		t.Fatal("expected alias to be deleted")
	}
}

func TestAliasDeleteNotFound(t *testing.T) {
	s, _ := NewAliasStore(tempAliasPath(t))
	if err := s.Delete("missing"); err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestAliasInvalidName(t *testing.T) {
	s, _ := NewAliasStore(tempAliasPath(t))
	if err := s.Set("1bad", "SOME_KEY"); err == nil {
		t.Fatal("expected error for invalid alias name")
	}
	if err := s.Set("", "SOME_KEY"); err == nil {
		t.Fatal("expected error for empty alias name")
	}
}

func TestAliasResolveUnknown(t *testing.T) {
	s, _ := NewAliasStore(tempAliasPath(t))
	resolved := s.Resolve("UNKNOWN_KEY")
	if resolved != "UNKNOWN_KEY" {
		t.Fatalf("expected fallback to input, got %q", resolved)
	}
}

func TestAliasList(t *testing.T) {
	s, _ := NewAliasStore(tempAliasPath(t))
	_ = s.Set("a", "KEY_A")
	_ = s.Set("b", "KEY_B")
	list := s.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(list))
	}
}

func TestAliasLoadInvalidJSON(t *testing.T) {
	path := tempAliasPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	_, err := NewAliasStore(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
