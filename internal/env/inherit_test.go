package env

import (
	"strings"
	"testing"
)

func TestInheritAllKeys(t *testing.T) {
	parent := map[string]string{"A": "1", "B": "2"}
	child := map[string]string{"C": "3"}

	out, res := Inherit(parent, child, InheritOptions{})

	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Fatalf("unexpected result: %v", out)
	}
	if len(res.Inherited) != 2 {
		t.Fatalf("expected 2 inherited, got %d", len(res.Inherited))
	}
	if len(res.Skipped) != 0 {
		t.Fatalf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestInheritSkipsExistingWithoutOverwrite(t *testing.T) {
	parent := map[string]string{"A": "parent_a", "B": "parent_b"}
	child := map[string]string{"A": "child_a"}

	out, res := Inherit(parent, child, InheritOptions{Overwrite: false})

	if out["A"] != "child_a" {
		t.Fatalf("expected child_a, got %s", out["A"])
	}
	if out["B"] != "parent_b" {
		t.Fatalf("expected parent_b, got %s", out["B"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Fatalf("expected A skipped, got %v", res.Skipped)
	}
}

func TestInheritOverwritesExisting(t *testing.T) {
	parent := map[string]string{"A": "parent_a"}
	child := map[string]string{"A": "child_a"}

	out, res := Inherit(parent, child, InheritOptions{Overwrite: true})

	if out["A"] != "parent_a" {
		t.Fatalf("expected parent_a, got %s", out["A"])
	}
	if len(res.Inherited) != 1 {
		t.Fatalf("expected 1 inherited, got %d", len(res.Inherited))
	}
}

func TestInheritFilterKeys(t *testing.T) {
	parent := map[string]string{"A": "1", "B": "2", "C": "3"}
	child := map[string]string{}

	out, res := Inherit(parent, child, InheritOptions{Keys: []string{"A", "C"}})

	if _, ok := out["B"]; ok {
		t.Fatal("B should not be inherited")
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Fatalf("unexpected result: %v", out)
	}
	if len(res.Inherited) != 2 {
		t.Fatalf("expected 2 inherited, got %d", len(res.Inherited))
	}
}

func TestInheritDoesNotMutateInputs(t *testing.T) {
	parent := map[string]string{"A": "1"}
	child := map[string]string{"B": "2"}

	_, _ = Inherit(parent, child, InheritOptions{})

	if _, ok := child["A"]; ok {
		t.Fatal("child should not be mutated")
	}
}

func TestInheritSummaryString(t *testing.T) {
	res := InheritResult{
		Inherited: []string{"A", "B"},
		Skipped:   []string{"C"},
	}
	s := res.SummaryString()
	if !strings.Contains(s, "2") || !strings.Contains(s, "1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}
