package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempArchivePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "archive.json")
}

func TestArchiveSaveAndGet(t *testing.T) {
	s, err := NewArchiveStore(tempArchivePath(t))
	if err != nil {
		t.Fatalf("NewArchiveStore: %v", err)
	}
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := s.Save("v1", vars, "initial"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	e, ok := s.Get("v1")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", e.Vars["FOO"])
	}
	if e.Note != "initial" {
		t.Errorf("expected note 'initial', got %s", e.Note)
	}
	if e.ArchivedAt.IsZero() {
		t.Error("expected non-zero ArchivedAt")
	}
}

func TestArchivePersists(t *testing.T) {
	path := tempArchivePath(t)
	s, _ := NewArchiveStore(path)
	_ = s.Save("snap", map[string]string{"K": "V"}, "")

	s2, err := NewArchiveStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if _, ok := s2.Get("snap"); !ok {
		t.Error("expected archive to persist after reload")
	}
}

func TestArchiveList(t *testing.T) {
	s, _ := NewArchiveStore(tempArchivePath(t))
	_ = s.Save("beta", map[string]string{}, "")
	_ = s.Save("alpha", map[string]string{}, "")
	names := s.List()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected list: %v", names)
	}
}

func TestArchiveDelete(t *testing.T) {
	s, _ := NewArchiveStore(tempArchivePath(t))
	_ = s.Save("tmp", map[string]string{"X": "1"}, "")
	if err := s.Delete("tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, ok := s.Get("tmp"); ok {
		t.Error("expected entry to be deleted")
	}
}

func TestArchiveDeleteNotFound(t *testing.T) {
	s, _ := NewArchiveStore(tempArchivePath(t))
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error for missing archive")
	}
}

func TestArchiveSaveEmptyName(t *testing.T) {
	s, _ := NewArchiveStore(tempArchivePath(t))
	if err := s.Save("", map[string]string{}, ""); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestArchiveDoesNotMutateInput(t *testing.T) {
	s, _ := NewArchiveStore(tempArchivePath(t))
	vars := map[string]string{"A": "1"}
	_ = s.Save("x", vars, "")
	vars["A"] = "mutated"
	e, _ := s.Get("x")
	if e.Vars["A"] != "1" {
		t.Error("archive should not reflect mutations to original map")
	}
}

func TestArchiveLoadInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0600)
	if _, err := NewArchiveStore(path); err == nil {
		t.Error("expected error on invalid JSON")
	}
}
