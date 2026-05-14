package env

import (
	"path/filepath"
	"testing"
	"time"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "rollback_history.json")
}

func TestRollbackHistoryAppendAndAll(t *testing.T) {
	h := NewRollbackHistory(tempHistoryPath(t))

	e1 := RollbackEntry{Timestamp: time.Now().Add(-2 * time.Minute), ToVersion: "v1", RestoredBy: "alice"}
	e2 := RollbackEntry{Timestamp: time.Now().Add(-1 * time.Minute), ToVersion: "v2", RestoredBy: "bob"}

	if err := h.Append(e1); err != nil {
		t.Fatalf("append e1: %v", err)
	}
	if err := h.Append(e2); err != nil {
		t.Fatalf("append e2: %v", err)
	}

	all, err := h.All()
	if err != nil {
		t.Fatalf("all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Sorted descending: e2 should be first
	if all[0].ToVersion != "v2" {
		t.Errorf("expected v2 first, got %s", all[0].ToVersion)
	}
}

func TestRollbackHistoryEmptyFile(t *testing.T) {
	h := NewRollbackHistory(tempHistoryPath(t))
	all, err := h.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected empty history")
	}
}

func TestRollbackHistorySince(t *testing.T) {
	h := NewRollbackHistory(tempHistoryPath(t))
	base := time.Now().Add(-10 * time.Minute)

	_ = h.Append(RollbackEntry{Timestamp: base.Add(1 * time.Minute), ToVersion: "v1"})
	_ = h.Append(RollbackEntry{Timestamp: base.Add(5 * time.Minute), ToVersion: "v2"})
	_ = h.Append(RollbackEntry{Timestamp: base.Add(9 * time.Minute), ToVersion: "v3"})

	results, err := h.Since(base.Add(4 * time.Minute))
	if err != nil {
		t.Fatalf("since: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results after cutoff, got %d", len(results))
	}
}

func TestRollbackHistoryPersists(t *testing.T) {
	path := tempHistoryPath(t)
	h1 := NewRollbackHistory(path)
	_ = h1.Append(RollbackEntry{Timestamp: time.Now(), ToVersion: "v99", RestoredBy: "carol"})

	h2 := NewRollbackHistory(path)
	all, err := h2.All()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(all) != 1 || all[0].ToVersion != "v99" {
		t.Errorf("persisted entry not found")
	}
}
