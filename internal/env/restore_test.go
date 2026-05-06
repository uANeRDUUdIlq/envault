package env

import (
	"strings"
	"testing"
)

func setupRestorer() (*SnapshotStore, *Restorer) {
	store := NewSnapshotStore()
	return store, NewRestorer(store)
}

func TestRestoreAppliesSnapshot(t *testing.T) {
	store, restorer := setupRestorer()
	vars := map[string]string{"KEY": "old", "KEEP": "same"}
	snap := store.Add("alice", "baseline", vars)

	current := map[string]string{"KEY": "new", "KEEP": "same", "EXTRA": "x"}
	restored, _, err := restorer.Restore(snap.ID, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if restored["KEY"] != "old" {
		t.Errorf("expected KEY=old after restore, got %s", restored["KEY"])
	}
	if _, ok := restored["EXTRA"]; ok {
		t.Error("EXTRA should not exist in restored map")
	}
}

func TestRestoreDiffIsCorrect(t *testing.T) {
	store, restorer := setupRestorer()
	snap := store.Add("bob", "v1", map[string]string{"A": "1", "B": "2"})

	current := map[string]string{"A": "changed", "C": "3"}
	_, result, err := restorer.Restore(snap.ID, current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diff.Added) != 1 || result.Diff.Added["B"] != "2" {
		t.Errorf("expected B added in diff")
	}
	if len(result.Diff.Removed) != 1 {
		t.Errorf("expected C removed in diff")
	}
	if len(result.Diff.Updated) != 1 {
		t.Errorf("expected A updated in diff")
	}
}

func TestRestoreNotFound(t *testing.T) {
	_, restorer := setupRestorer()
	_, _, err := restorer.Restore("ghost", map[string]string{})
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRestoreResultSummaryString(t *testing.T) {
	store, restorer := setupRestorer()
	snap := store.Add("carol", "summary test", map[string]string{"X": "1"})
	_, result, err := restorer.Restore(snap.ID, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	summary := result.SummaryString()
	if !strings.Contains(summary, "carol") {
		t.Errorf("summary missing author: %s", summary)
	}
	if !strings.Contains(summary, "summary test") {
		t.Errorf("summary missing message: %s", summary)
	}
	if !strings.Contains(summary, "added:1") {
		t.Errorf("summary missing added count: %s", summary)
	}
}
