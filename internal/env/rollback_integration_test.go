package env

import (
	"path/filepath"
	"testing"
	"time"
)

func TestRollbackRoundtrip(t *testing.T) {
	dir := t.TempDir()
	ss := NewSnapshotStore(filepath.Join(dir, "snapshots.json"))
	v := NewVault()
	v.Set("DB_HOST", "localhost")
	v.Set("DB_PORT", "5432")
	v.Set("API_KEY", "secret")

	snap := Snapshot{
		ID:        "rollback-snap",
		Namespace: "production",
		Vars:      map[string]string{"DB_HOST": "prod-host", "DB_PORT": "5433"},
		CreatedAt: time.Now(),
	}
	if err := ss.Add(snap); err != nil {
		t.Fatalf("add snapshot: %v", err)
	}

	rb := NewRollbacker(ss, v)
	res, err := rb.Rollback("rollback-snap", "deployer", false)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}

	if v.Get("DB_HOST") != "prod-host" {
		t.Errorf("DB_HOST not updated")
	}
	if v.Get("API_KEY") != "secret" {
		t.Errorf("API_KEY should be preserved (no purge)")
	}

	h := NewRollbackHistory(filepath.Join(dir, "history.json"))
	if err := h.Append(res.Entry); err != nil {
		t.Fatalf("append history: %v", err)
	}

	all, err := h.All()
	if err != nil {
		t.Fatalf("history all: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(all))
	}
	if all[0].RestoredBy != "deployer" {
		t.Errorf("unexpected restored_by: %s", all[0].RestoredBy)
	}
}

func TestRollbackPurgeIntegration(t *testing.T) {
	dir := t.TempDir()
	ss := NewSnapshotStore(filepath.Join(dir, "snapshots.json"))
	v := NewVault()
	v.Set("A", "1")
	v.Set("B", "2")
	v.Set("C", "3")

	_ = ss.Add(Snapshot{
		ID:        "purge-snap",
		Vars:      map[string]string{"A": "10"},
		CreatedAt: time.Now(),
	})

	rb := NewRollbacker(ss, v)
	_, err := rb.Rollback("purge-snap", "admin", true)
	if err != nil {
		t.Fatalf("rollback with purge: %v", err)
	}

	if v.Get("A") != "10" {
		t.Errorf("A should be updated")
	}
	if v.Get("B") != "" || v.Get("C") != "" {
		t.Errorf("B and C should be purged")
	}
}
