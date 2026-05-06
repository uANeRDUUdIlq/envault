package env

import (
	"testing"
)

func TestMergeNoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	ours := map[string]string{"A": "1", "B": "updated"}
	theirs := map[string]string{"A": "1", "B": "2", "C": "3"}

	r := Merge(base, ours, theirs, MergeStrategyOurs)

	if r.Merged["A"] != "1" {
		t.Errorf("expected A=1, got %s", r.Merged["A"])
	}
	if r.Merged["B"] != "updated" {
		t.Errorf("expected B=updated (ours changed), got %s", r.Merged["B"])
	}
	if r.Merged["C"] != "3" {
		t.Errorf("expected C=3 (theirs added), got %s", r.Merged["C"])
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
}

func TestMergeConflictResolvesOurs(t *testing.T) {
	base := map[string]string{"KEY": "base"}
	ours := map[string]string{"KEY": "local"}
	theirs := map[string]string{"KEY": "remote"}

	r := Merge(base, ours, theirs, MergeStrategyOurs)

	if r.Merged["KEY"] != "local" {
		t.Errorf("expected local value, got %s", r.Merged["KEY"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0].Resolved != "local" {
		t.Errorf("expected 1 conflict resolved to 'local'")
	}
}

func TestMergeConflictResolvesTheirs(t *testing.T) {
	base := map[string]string{"KEY": "base"}
	ours := map[string]string{"KEY": "local"}
	theirs := map[string]string{"KEY": "remote"}

	r := Merge(base, ours, theirs, MergeStrategyTheirs)

	if r.Merged["KEY"] != "remote" {
		t.Errorf("expected remote value, got %s", r.Merged["KEY"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0].Resolved != "remote" {
		t.Errorf("expected 1 conflict resolved to 'remote'")
	}
}

func TestMergeTheyAddedWeNeverHad(t *testing.T) {
	base := map[string]string{}
	ours := map[string]string{}
	theirs := map[string]string{"NEW": "value"}

	r := Merge(base, ours, theirs, MergeStrategyOurs)

	if r.Merged["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %s", r.Merged["NEW"])
	}
	if len(r.Added) != 1 || r.Added[0] != "NEW" {
		t.Errorf("expected NEW in Added list")
	}
}

func TestMergeWeDeletedKey(t *testing.T) {
	base := map[string]string{"OLD": "v"}
	ours := map[string]string{}
	theirs := map[string]string{"OLD": "v"}

	r := Merge(base, ours, theirs, MergeStrategyOurs)

	if _, ok := r.Merged["OLD"]; ok {
		t.Error("expected OLD to be absent after our deletion")
	}
}

func TestMergeEmptyMaps(t *testing.T) {
	r := Merge(
		map[string]string{},
		map[string]string{},
		map[string]string{},
		MergeStrategyOurs,
	)
	if len(r.Merged) != 0 {
		t.Errorf("expected empty merged map")
	}
}
