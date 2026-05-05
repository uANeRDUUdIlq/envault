package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/store"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "secrets.json")
}

func TestNewStoreCreatesEmpty(t *testing.T) {
	s, err := store.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Secrets) != 0 {
		t.Errorf("expected empty store, got %d entries", len(s.Secrets))
	}
}

func TestPutAndGet(t *testing.T) {
	s, _ := store.New(tempStorePath(t))
	s.Put("DB_PASSWORD", "encrypted-blob-abc")

	entry, err := s.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if entry.Ciphertext != "encrypted-blob-abc" {
		t.Errorf("expected ciphertext %q, got %q", "encrypted-blob-abc", entry.Ciphertext)
	}
}

func TestGetNotFound(t *testing.T) {
	s, _ := store.New(tempStorePath(t))
	_, err := s.Get("MISSING_KEY")
	if err != store.ErrSecretNotFound {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	s, _ := store.New(tempStorePath(t))
	s.Put("API_KEY", "blob")

	if err := s.Delete("API_KEY"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err := s.Get("API_KEY")
	if err != store.ErrSecretNotFound {
		t.Errorf("expected ErrSecretNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	s, _ := store.New(tempStorePath(t))
	if err := s.Delete("GHOST"); err != store.ErrSecretNotFound {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestSaveAndReload(t *testing.T) {
	path := tempStorePath(t)
	s, _ := store.New(path)
	s.Put("SECRET_ONE", "cipher1")
	s.Put("SECRET_TWO", "cipher2")

	if err := s.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("store file not created: %v", err)
	}

	s2, err := store.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(s2.Secrets) != 2 {
		t.Errorf("expected 2 secrets after reload, got %d", len(s2.Secrets))
	}
	if e, _ := s2.Get("SECRET_ONE"); e.Ciphertext != "cipher1" {
		t.Errorf("reloaded ciphertext mismatch")
	}
}
