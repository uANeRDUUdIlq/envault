package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPinPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestPinAndGet(t *testing.T) {
	ps, err := NewPinStore(tempPinPath(t))
	if err != nil {
		t.Fatalf("NewPinStore: %v", err)
	}
	if err := ps.Pin("production", "snap-001", "alice", "stable release"); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	e, ok := ps.Get("production")
	if !ok {
		t.Fatal("expected pin to exist")
	}
	if e.SnapshotID != "snap-001" || e.PinnedBy != "alice" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestPinPersists(t *testing.T) {
	path := tempPinPath(t)
	ps, _ := NewPinStore(path)
	_ = ps.Pin("staging", "snap-002", "bob", "")

	ps2, err := NewPinStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := ps2.Get("staging")
	if !ok || e.SnapshotID != "snap-002" {
		t.Errorf("pin not persisted: %+v", e)
	}
}

func TestUnpin(t *testing.T) {
	ps, _ := NewPinStore(tempPinPath(t))
	_ = ps.Pin("dev", "snap-003", "carol", "")
	if err := ps.Unpin("dev"); err != nil {
		t.Fatalf("Unpin: %v", err)
	}
	if _, ok := ps.Get("dev"); ok {
		t.Error("expected pin to be removed")
	}
}

func TestUnpinNotFound(t *testing.T) {
	ps, _ := NewPinStore(tempPinPath(t))
	if err := ps.Unpin("nonexistent"); err == nil {
		t.Error("expected error for missing pin")
	}
}

func TestPinList(t *testing.T) {
	ps, _ := NewPinStore(tempPinPath(t))
	_ = ps.Pin("env-a", "snap-a", "alice", "")
	_ = ps.Pin("env-b", "snap-b", "bob", "")
	list := ps.List()
	if len(list) != 2 {
		t.Errorf("expected 2 pins, got %d", len(list))
	}
}

func TestPinEmptyEnvironment(t *testing.T) {
	ps, _ := NewPinStore(tempPinPath(t))
	if err := ps.Pin("", "snap-x", "alice", ""); err == nil {
		t.Error("expected error for empty environment")
	}
}

func TestPinEmptySnapshotID(t *testing.T) {
	ps, _ := NewPinStore(tempPinPath(t))
	if err := ps.Pin("prod", "", "alice", ""); err == nil {
		t.Error("expected error for empty snapshot ID")
	}
}

func TestNewPinStoreInvalidJSON(t *testing.T) {
	path := tempPinPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	if _, err := NewPinStore(path); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
