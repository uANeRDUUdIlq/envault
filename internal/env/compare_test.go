package env

import (
	"testing"
)

func TestCompareIdentical(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(a, b)
	if !r.IsIdentical() {
		t.Fatal("expected identical result")
	}
	if len(r.Identical) != 2 {
		t.Fatalf("expected 2 identical keys, got %d", len(r.Identical))
	}
}

func TestCompareOnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "x"}
	b := map[string]string{"FOO": "bar"}
	r := Compare(a, b)
	if r.IsIdentical() {
		t.Fatal("expected non-identical result")
	}
	if v, ok := r.OnlyInA["ONLY_A"]; !ok || v != "x" {
		t.Fatalf("expected ONLY_A=x in OnlyInA, got %v", r.OnlyInA)
	}
}

func TestCompareOnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "y"}
	r := Compare(a, b)
	if v, ok := r.OnlyInB["ONLY_B"]; !ok || v != "y" {
		t.Fatalf("expected ONLY_B=y in OnlyInB, got %v", r.OnlyInB)
	}
}

func TestCompareDifferentValues(t *testing.T) {
	a := map[string]string{"DB_URL": "postgres://dev", "PORT": "3000"}
	b := map[string]string{"DB_URL": "postgres://prod", "PORT": "3000"}
	r := Compare(a, b)
	pair, ok := r.Different["DB_URL"]
	if !ok {
		t.Fatal("expected DB_URL in Different")
	}
	if pair[0] != "postgres://dev" || pair[1] != "postgres://prod" {
		t.Fatalf("unexpected pair values: %v", pair)
	}
	if len(r.Identical) != 1 || r.Identical[0] != "PORT" {
		t.Fatalf("expected PORT in Identical, got %v", r.Identical)
	}
}

func TestCompareSummaryLines(t *testing.T) {
	a := map[string]string{"A": "1", "C": "old"}
	b := map[string]string{"B": "2", "C": "new"}
	r := Compare(a, b)
	lines := r.SummaryLines()
	if len(lines) == 0 {
		t.Fatal("expected summary lines")
	}
	found := map[string]bool{}
	for _, l := range lines {
		found[l] = true
	}
	if !found["< A=1"] {
		t.Error("expected '< A=1' in summary")
	}
	if !found["> B=2"] {
		t.Error("expected '> B=2' in summary")
	}
	if !found["~ C: old -> new"] {
		t.Error("expected '~ C: old -> new' in summary")
	}
}

func TestCompareEmptyMaps(t *testing.T) {
	r := Compare(map[string]string{}, map[string]string{})
	if !r.IsIdentical() {
		t.Fatal("expected identical for two empty maps")
	}
}
