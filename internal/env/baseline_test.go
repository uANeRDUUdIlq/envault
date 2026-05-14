package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempBaselinePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "baselines.json")
}

func TestBaselineSetAndGet(t *testing.T) {
	bs, err := NewBaselineStore(tempBaselinePath(t))
	if err != nil {
		t.Fatalf("NewBaselineStore: %v", err)
	}
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := bs.Set("prod", vars); err != nil {
		t.Fatalf("Set: %v", err)
	}
	b, err := bs.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if b.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", b.Vars["FOO"])
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestBaselinePersists(t *testing.T) {
	path := tempBaselinePath(t)
	bs, _ := NewBaselineStore(path)
	_ = bs.Set("staging", map[string]string{"KEY": "val"})

	bs2, err := NewBaselineStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	b, err := bs2.Get("staging")
	if err != nil {
		t.Fatalf("Get after reload: %v", err)
	}
	if b.Vars["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", b.Vars["KEY"])
	}
}

func TestBaselineGetNotFound(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	_, err := bs.Get("missing")
	if err == nil {
		t.Error("expected error for missing baseline")
	}
}

func TestBaselineDelete(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	_ = bs.Set("dev", map[string]string{"X": "1"})
	if err := bs.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := bs.Get("dev"); err == nil {
		t.Error("expected error after delete")
	}
}

func TestBaselineDeleteNotFound(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	if err := bs.Delete("nope"); err == nil {
		t.Error("expected error deleting non-existent baseline")
	}
}

func TestBaselineList(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	_ = bs.Set("a", map[string]string{})
	_ = bs.Set("b", map[string]string{})
	names := bs.List()
	if len(names) != 2 {
		t.Errorf("expected 2 baselines, got %d", len(names))
	}
}

func TestBaselineDrift(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	_ = bs.Set("prod", map[string]string{"FOO": "old", "GONE": "bye"})
	current := map[string]string{"FOO": "new", "ADDED": "yes"}
	dr, err := bs.Drift("prod", current)
	if err != nil {
		t.Fatalf("Drift: %v", err)
	}
	if len(dr.Updated) != 1 || dr.Updated[0].Key != "FOO" {
		t.Errorf("expected FOO in Updated, got %+v", dr.Updated)
	}
	if len(dr.Added) != 1 || dr.Added[0].Key != "ADDED" {
		t.Errorf("expected ADDED in Added, got %+v", dr.Added)
	}
	if len(dr.Removed) != 1 || dr.Removed[0].Key != "GONE" {
		t.Errorf("expected GONE in Removed, got %+v", dr.Removed)
	}
}

func TestBaselineSetEmptyName(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	if err := bs.Set("", map[string]string{"K": "v"}); err == nil {
		t.Error("expected error for empty baseline name")
	}
}

func TestBaselineIsolatesVarsCopy(t *testing.T) {
	bs, _ := NewBaselineStore(tempBaselinePath(t))
	vars := map[string]string{"KEY": "original"}
	_ = bs.Set("snap", vars)
	vars["KEY"] = "mutated"
	b, _ := bs.Get("snap")
	if b.Vars["KEY"] != "original" {
		t.Error("baseline should not reflect mutation of original map")
	}
}

func TestBaselineStoreCreatesFile(t *testing.T) {
	path := tempBaselinePath(t)
	bs, _ := NewBaselineStore(path)
	_ = bs.Set("x", map[string]string{"A": "1"})
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
