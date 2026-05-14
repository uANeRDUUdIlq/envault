package env

import (
	"strings"
	"testing"
)

func TestCloneCopiesAllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dest := map[string]string{}
	c := NewCloner(src, dest, CloneOptions{})
	out, res := c.Clone("staging", "prod")
	if res.Copied != 2 || res.Skipped != 0 || res.Overwritten != 0 {
		t.Fatalf("unexpected result: %+v", res)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("expected keys in output, got %v", out)
	}
}

func TestCloneSkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	dest := map[string]string{"A": "old"}
	c := NewCloner(src, dest, CloneOptions{Overwrite: false})
	out, res := c.Clone("staging", "prod")
	if res.Skipped != 1 || res.Copied != 1 {
		t.Fatalf("unexpected result: %+v", res)
	}
	if out["A"] != "old" {
		t.Fatalf("expected A to remain old, got %s", out["A"])
	}
}

func TestCloneOverwritesExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dest := map[string]string{"A": "old"}
	c := NewCloner(src, dest, CloneOptions{Overwrite: true})
	out, res := c.Clone("staging", "prod")
	if res.Overwritten != 1 {
		t.Fatalf("expected 1 overwritten, got %+v", res)
	}
	if out["A"] != "new" {
		t.Fatalf("expected A=new, got %s", out["A"])
	}
}

func TestCloneFilterKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dest := map[string]string{}
	c := NewCloner(src, dest, CloneOptions{FilterKeys: []string{"A", "C"}})
	out, res := c.Clone("staging", "prod")
	if res.Copied != 2 {
		t.Fatalf("expected 2 copied, got %+v", res)
	}
	if _, ok := out["B"]; ok {
		t.Fatal("B should not be in output")
	}
}

func TestCloneSummaryString(t *testing.T) {
	r := CloneResult{
		SourceNamespace: "dev",
		DestNamespace:   "prod",
		Copied:          5,
		Skipped:         1,
		Overwritten:     2,
	}
	s := r.SummaryString()
	if !strings.Contains(s, "dev") || !strings.Contains(s, "prod") {
		t.Fatalf("summary missing namespace info: %s", s)
	}
	if !strings.Contains(s, "5 copied") {
		t.Fatalf("summary missing copied count: %s", s)
	}
}

func TestClonePreservesDestKeys(t *testing.T) {
	src := map[string]string{"A": "1"}
	dest := map[string]string{"Z": "99"}
	c := NewCloner(src, dest, CloneOptions{})
	out, _ := c.Clone("staging", "prod")
	if out["Z"] != "99" {
		t.Fatal("expected dest key Z to be preserved")
	}
}
