package env

import (
	"testing"
	"time"
)

func setupRollback(t *testing.T) (*Rollbacker, *SnapshotStore, *Vault) {
	t.Helper()
	dir := t.TempDir()
	ss := NewSnapshotStore(dir + "/snapshots.json")
	v := NewVault()
	v.Set("KEY_A", "old_a")
	v.Set("KEY_B", "old_b")
	v.Set("KEY_C", "old_c")
	return NewRollbacker(ss, v), ss, v
}

func TestRollbackAppliesSnapshot(t *testing.T) {
	rb, ss, v := setupRollback(t)

	snap := Snapshot{
		ID:        "snap-001",
		Namespace: "default",
		Vars:      map[string]string{"KEY_A": "new_a", "KEY_B": "new_b"},
		CreatedAt: time.Now(),
	}
	if err := ss.Add(snap); err != nil {
		t.Fatalf("unexpected error adding snapshot: %v", err)
	}

	res, err := rb.Rollback("snap-001", "alice", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.Get("KEY_A") != "new_a" {
		t.Errorf("expected KEY_A=new_a, got %s", v.Get("KEY_A"))
	}
	if v.Get("KEY_B") != "new_b" {
		t.Errorf("expected KEY_B=new_b, got %s", v.Get("KEY_B"))
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
}

func TestRollbackSkipsKeysNotInSnapshot(t *testing.T) {
	rb, ss, v := setupRollback(t)

	snap := Snapshot{
		ID:   "snap-002",
		Vars: map[string]string{"KEY_A": "new_a"},
	}
	_ = ss.Add(snap)

	res, err := rb.Rollback("snap-002", "bob", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.Get("KEY_B") != "old_b" {
		t.Errorf("KEY_B should be untouched")
	}
	if len(res.Skipped) == 0 {
		t.Error("expected skipped keys when purge=false")
	}
}

func TestRollbackPurgeRemovesExtraKeys(t *testing.T) {
	rb, ss, v := setupRollback(t)

	snap := Snapshot{
		ID:   "snap-003",
		Vars: map[string]string{"KEY_A": "new_a"},
	}
	_ = ss.Add(snap)

	_, err := rb.Rollback("snap-003", "carol", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.Get("KEY_B") != "" {
		t.Errorf("KEY_B should have been purged")
	}
}

func TestRollbackNotFound(t *testing.T) {
	rb, _, _ := setupRollback(t)
	_, err := rb.Rollback("nonexistent", "dave", false)
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRollbackEmptyUser(t *testing.T) {
	rb, ss, _ := setupRollback(t)
	_ = ss.Add(Snapshot{ID: "snap-004", Vars: map[string]string{}})
	_, err := rb.Rollback("snap-004", "", false)
	if err == nil {
		t.Error("expected error for empty user")
	}
}

func TestRollbackSummaryString(t *testing.T) {
	res := &RollbackResult{
		Entry:   RollbackEntry{ToVersion: "snap-005"},
		Applied: []string{"A", "B"},
		Skipped: []string{"C"},
	}
	s := res.RollbackSummary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
