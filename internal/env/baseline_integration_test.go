package env

import (
	"path/filepath"
	"testing"
)

func TestBaselineDriftRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bl.json")
	bs, err := NewBaselineStore(path)
	if err != nil {
		t.Fatalf("NewBaselineStore: %v", err)
	}

	initial := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SECRET":  "abc123",
	}
	if err := bs.Set("v1", initial); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Reload to verify persistence
	bs2, err := NewBaselineStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}

	current := map[string]string{
		"DB_HOST": "prod-db.example.com",
		"DB_PORT": "5432",
		"NEW_KEY": "newval",
	}
	dr, err := bs2.Drift("v1", current)
	if err != nil {
		t.Fatalf("Drift: %v", err)
	}

	if !HasChanges(dr) {
		t.Error("expected drift to have changes")
	}
	if len(dr.Updated) != 1 {
		t.Errorf("expected 1 updated key, got %d", len(dr.Updated))
	}
	if len(dr.Removed) != 1 || dr.Removed[0].Key != "SECRET" {
		t.Errorf("expected SECRET removed, got %+v", dr.Removed)
	}
	if len(dr.Added) != 1 || dr.Added[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY added, got %+v", dr.Added)
	}
}

func TestBaselineMultipleSnapshots(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bl.json")
	bs, _ := NewBaselineStore(path)

	_ = bs.Set("dev", map[string]string{"ENV": "development"})
	_ = bs.Set("prod", map[string]string{"ENV": "production", "STRICT": "true"})

	names := bs.List()
	if len(names) != 2 {
		t.Fatalf("expected 2 baselines, got %d", len(names))
	}

	drDev, err := bs.Drift("dev", map[string]string{"ENV": "development"})
	if err != nil {
		t.Fatalf("Drift dev: %v", err)
	}
	if HasChanges(drDev) {
		t.Error("expected no drift for dev baseline")
	}

	drProd, err := bs.Drift("prod", map[string]string{"ENV": "development"})
	if err != nil {
		t.Fatalf("Drift prod: %v", err)
	}
	if !HasChanges(drProd) {
		t.Error("expected drift for prod baseline")
	}
}
