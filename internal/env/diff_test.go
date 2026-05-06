package env

import (
	"testing"
)

func TestDiffAdded(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar", "BAZ": "qux"}

	changes := Diff(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != ChangeAdded || changes[0].Key != "BAZ" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffRemoved(t *testing.T) {
	old := map[string]string{"FOO": "bar", "OLD_KEY": "value"}
	new := map[string]string{"FOO": "bar"}

	changes := Diff(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != ChangeRemoved || changes[0].Key != "OLD_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffUpdated(t *testing.T) {
	old := map[string]string{"FOO": "old_value"}
	new := map[string]string{"FOO": "new_value"}

	changes := Diff(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Type != ChangeUpdated || c.OldVal != "old_value" || c.NewVal != "new_value" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiffNoChanges(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	new := map[string]string{"FOO": "bar", "BAZ": "qux"}

	changes := Diff(old, new)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestHasChanges(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "2"}

	if !HasChanges(old, new) {
		t.Error("expected HasChanges to return true")
	}

	if HasChanges(old, old) {
		t.Error("expected HasChanges to return false for identical maps")
	}
}

func TestChangeString(t *testing.T) {
	cases := []struct {
		change Change
		want   string
	}{
		{Change{Key: "X", Type: ChangeAdded, NewVal: "1"}, "+ X=1"},
		{Change{Key: "X", Type: ChangeRemoved, OldVal: "1"}, "- X=1"},
		{Change{Key: "X", Type: ChangeUpdated, OldVal: "1", NewVal: "2"}, "~ X: 1 -> 2"},
	}
	for _, tc := range cases {
		if got := tc.change.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}
