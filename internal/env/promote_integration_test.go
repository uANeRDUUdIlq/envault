package env

import (
	"testing"
)

func TestPromoteRoundtrip(t *testing.T) {
	path := tempNamespacePath(t)
	ns, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatalf("NewNamespaceStore: %v", err)
	}
	for _, name := range []string{"dev", "staging", "production"} {
		if err := ns.Add(name); err != nil {
			t.Fatalf("add %s: %v", name, err)
		}
	}
	_ = ns.SetVars("dev", map[string]string{"APP_ENV": "dev", "DB_URL": "localhost"})

	p := NewPromoter(ns)

	// dev -> staging
	res1, err := p.Promote("dev", "staging", PromoteOptions{})
	if err != nil {
		t.Fatalf("dev->staging: %v", err)
	}
	if len(res1.Copied) != 2 {
		t.Errorf("expected 2 copied dev->staging, got %d", len(res1.Copied))
	}

	// staging -> production
	res2, err := p.Promote("staging", "production", PromoteOptions{Keys: []string{"DB_URL"}})
	if err != nil {
		t.Fatalf("staging->production: %v", err)
	}
	if len(res2.Copied) != 1 {
		t.Errorf("expected 1 copied staging->production, got %d", len(res2.Copied))
	}

	vars, _ := ns.GetVars("production")
	if vars["DB_URL"] != "localhost" {
		t.Errorf("expected DB_URL=localhost in production, got %q", vars["DB_URL"])
	}
	if _, ok := vars["APP_ENV"]; ok {
		t.Error("APP_ENV should not be in production (filtered out)")
	}
}

func TestPromoteChainOverwrite(t *testing.T) {
	path := tempNamespacePath(t)
	ns, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatalf("NewNamespaceStore: %v", err)
	}
	for _, name := range []string{"staging", "production"} {
		_ = ns.Add(name)
	}
	_ = ns.SetVars("staging", map[string]string{"SECRET": "v2"})
	_ = ns.SetVars("production", map[string]string{"SECRET": "v1"})

	p := NewPromoter(ns)
	res, err := p.Promote("staging", "production", PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("promote: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
	vars, _ := ns.GetVars("production")
	if vars["SECRET"] != "v2" {
		t.Errorf("expected SECRET=v2, got %q", vars["SECRET"])
	}
}
