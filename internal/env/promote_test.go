package env

import (
	"testing"
)

func setupPromoter(t *testing.T) (*Promoter, *NamespaceStore) {
	t.Helper()
	path := tempNamespacePath(t)
	ns, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatalf("NewNamespaceStore: %v", err)
	}
	if err := ns.Add("staging"); err != nil {
		t.Fatalf("add staging: %v", err)
	}
	if err := ns.Add("production"); err != nil {
		t.Fatalf("add production: %v", err)
	}
	return NewPromoter(ns), ns
}

func TestPromoteCopiesAllKeys(t *testing.T) {
	p, ns := setupPromoter(t)
	_ = ns.SetVars("staging", map[string]string{"FOO": "bar", "BAZ": "qux"})

	res, err := p.Promote("staging", "production", PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestPromoteSkipsExistingWithoutOverwrite(t *testing.T) {
	p, ns := setupPromoter(t)
	_ = ns.SetVars("staging", map[string]string{"FOO": "new"})
	_ = ns.SetVars("production", map[string]string{"FOO": "old"})

	res, err := p.Promote("staging", "production", PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	vars, _ := ns.GetVars("production")
	if vars["FOO"] != "old" {
		t.Errorf("expected FOO=old, got %q", vars["FOO"])
	}
}

func TestPromoteOverwritesExisting(t *testing.T) {
	p, ns := setupPromoter(t)
	_ = ns.SetVars("staging", map[string]string{"FOO": "new"})
	_ = ns.SetVars("production", map[string]string{"FOO": "old"})

	res, err := p.Promote("staging", "production", PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
	vars, _ := ns.GetVars("production")
	if vars["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", vars["FOO"])
	}
}

func TestPromoteFilterKeys(t *testing.T) {
	p, ns := setupPromoter(t)
	_ = ns.SetVars("staging", map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"})

	res, err := p.Promote("staging", "production", PromoteOptions{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	vars, _ := ns.GetVars("production")
	if _, ok := vars["BAR"]; ok {
		t.Error("BAR should not have been promoted")
	}
}

func TestPromoteSummaryString(t *testing.T) {
	res := &PromoteResult{
		SourceNamespace: "staging",
		TargetNamespace: "production",
		Copied:          []string{"A", "B"},
		Skipped:         []string{"C"},
		Overwritten:     []string{},
	}
	s := res.SummaryString()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
