package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempCheckpointPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoints.json")
}

func TestCheckpointSaveAndGet(t *testing.T) {
	cs, err := NewCheckpointStore(tempCheckpointPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := cs.Save("v1", vars, "initial checkpoint"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	cp, ok := cs.Get("v1")
	if !ok {
		t.Fatal("expected checkpoint to exist")
	}
	if cp.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", cp.Vars["FOO"])
	}
	if cp.Message != "initial checkpoint" {
		t.Errorf("unexpected message: %q", cp.Message)
	}
}

func TestCheckpointPersists(t *testing.T) {
	path := tempCheckpointPath(t)
	cs, _ := NewCheckpointStore(path)
	_ = cs.Save("v1", map[string]string{"KEY": "val"}, "")

	cs2, err := NewCheckpointStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if _, ok := cs2.Get("v1"); !ok {
		t.Error("checkpoint not found after reload")
	}
}

func TestCheckpointGetNotFound(t *testing.T) {
	cs, _ := NewCheckpointStore(tempCheckpointPath(t))
	if _, ok := cs.Get("missing"); ok {
		t.Error("expected not found")
	}
}

func TestCheckpointDelete(t *testing.T) {
	cs, _ := NewCheckpointStore(tempCheckpointPath(t))
	_ = cs.Save("v1", map[string]string{"A": "1"}, "")
	if err := cs.Delete("v1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, ok := cs.Get("v1"); ok {
		t.Error("expected checkpoint to be deleted")
	}
}

func TestCheckpointDeleteNotFound(t *testing.T) {
	cs, _ := NewCheckpointStore(tempCheckpointPath(t))
	if err := cs.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent checkpoint")
	}
}

func TestCheckpointList(t *testing.T) {
	cs, _ := NewCheckpointStore(tempCheckpointPath(t))
	_ = cs.Save("a", map[string]string{}, "")
	_ = cs.Save("b", map[string]string{}, "")
	names := cs.List()
	if len(names) != 2 {
		t.Errorf("expected 2 checkpoints, got %d", len(names))
	}
}

func TestCheckpointEmptyName(t *testing.T) {
	cs, _ := NewCheckpointStore(tempCheckpointPath(t))
	if err := cs.Save("", map[string]string{}, ""); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestCheckpointVarsCopied(t *testing.T) {
	cs, _ := NewCheckpointStore(tempCheckpointPath(t))
	vars := map[string]string{"X": "1"}
	_ = cs.Save("cp", vars, "")
	vars["X"] = "mutated"
	cp, _ := cs.Get("cp")
	if cp.Vars["X"] != "1" {
		t.Error("checkpoint should not reflect mutation of original map")
	}
}

func TestCheckpointMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no", "such", "file.json")
	_, err := NewCheckpointStore(path)
	if err == nil {
		t.Error("expected error for unreadable path")
	}
}

func TestCheckpointInvalidJSON(t *testing.T) {
	path := tempCheckpointPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0600)
	_, err := NewCheckpointStore(path)
	if err == nil {
		t.Error("expected parse error")
	}
}
