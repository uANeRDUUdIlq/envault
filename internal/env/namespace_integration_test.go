package env

import (
	"path/filepath"
	"testing"
)

// TestNamespaceRoundtrip verifies that namespaces survive a full save/reload cycle
// with all fields intact.
func TestNamespaceRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ns.json")

	s1, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatal(err)
	}
	namespaces := []Namespace{
		{Name: "production", Description: "Live traffic", Meta: map[string]string{"region": "us-east-1"}},
		{Name: "staging", Description: "Pre-release", Meta: map[string]string{"region": "us-west-2"}},
		{Name: "dev", Description: "Local development"},
	}
	for _, ns := range namespaces {
		if err := s1.Add(ns); err != nil {
			t.Fatalf("add %q: %v", ns.Name, err)
		}
	}

	s2, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, ns := range namespaces {
		got, ok := s2.Get(ns.Name)
		if !ok {
			t.Errorf("namespace %q missing after reload", ns.Name)
			continue
		}
		if got.Description != ns.Description {
			t.Errorf("%s description: got %q want %q", ns.Name, got.Description, ns.Description)
		}
		for k, v := range ns.Meta {
			if got.Meta[k] != v {
				t.Errorf("%s meta[%s]: got %q want %q", ns.Name, k, got.Meta[k], v)
			}
		}
	}
}

// TestNamespaceDeleteAndReload ensures deletions persist across store reloads.
func TestNamespaceDeleteAndReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ns.json")

	s1, _ := NewNamespaceStore(path)
	_ = s1.Add(Namespace{Name: "alpha"})
	_ = s1.Add(Namespace{Name: "beta"})
	_ = s1.Delete("alpha")

	s2, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s2.Get("alpha"); ok {
		t.Error("deleted namespace 'alpha' should not appear after reload")
	}
	if _, ok := s2.Get("beta"); !ok {
		t.Error("namespace 'beta' should still exist after reload")
	}
}
