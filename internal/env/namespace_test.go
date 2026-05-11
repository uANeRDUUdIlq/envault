package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempNamespacePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "namespaces.json")
}

func TestNamespaceAddAndGet(t *testing.T) {
	store, err := NewNamespaceStore(tempNamespacePath(t))
	if err != nil {
		t.Fatal(err)
	}
	ns := Namespace{Name: "production", Description: "prod env"}
	if err := store.Add(ns); err != nil {
		t.Fatal(err)
	}
	got, ok := store.Get("production")
	if !ok {
		t.Fatal("expected namespace to exist")
	}
	if got.Description != "prod env" {
		t.Errorf("description: got %q want %q", got.Description, "prod env")
	}
}

func TestNamespaceAddDuplicate(t *testing.T) {
	store, _ := NewNamespaceStore(tempNamespacePath(t))
	_ = store.Add(Namespace{Name: "staging"})
	if err := store.Add(Namespace{Name: "staging"}); err == nil {
		t.Error("expected error for duplicate namespace")
	}
}

func TestNamespaceInvalidName(t *testing.T) {
	store, _ := NewNamespaceStore(tempNamespacePath(t))
	invalidNames := []string{"", "1bad", "has space", "way-too-long-name-that-exceeds-sixty-four-characters-limit-here-x"}
	for _, name := range invalidNames {
		if err := store.Add(Namespace{Name: name}); err == nil {
			t.Errorf("expected error for name %q", name)
		}
	}
}

func TestNamespaceList(t *testing.T) {
	store, _ := NewNamespaceStore(tempNamespacePath(t))
	_ = store.Add(Namespace{Name: "alpha"})
	_ = store.Add(Namespace{Name: "beta"})
	list := store.List()
	if len(list) != 2 {
		t.Errorf("expected 2 namespaces, got %d", len(list))
	}
}

func TestNamespaceDelete(t *testing.T) {
	store, _ := NewNamespaceStore(tempNamespacePath(t))
	_ = store.Add(Namespace{Name: "dev"})
	if err := store.Delete("dev"); err != nil {
		t.Fatal(err)
	}
	if _, ok := store.Get("dev"); ok {
		t.Error("expected namespace to be deleted")
	}
}

func TestNamespaceDeleteNotFound(t *testing.T) {
	store, _ := NewNamespaceStore(tempNamespacePath(t))
	if err := store.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent namespace")
	}
}

func TestNamespacePersists(t *testing.T) {
	path := tempNamespacePath(t)
	s1, _ := NewNamespaceStore(path)
	_ = s1.Add(Namespace{Name: "ci", Description: "CI environment", Meta: map[string]string{"owner": "team"}})

	s2, err := NewNamespaceStore(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := s2.Get("ci")
	if !ok {
		t.Fatal("expected namespace to persist")
	}
	if got.Meta["owner"] != "team" {
		t.Errorf("meta owner: got %q want %q", got.Meta["owner"], "team")
	}
}

func TestNamespaceStoreGetNotFound(t *testing.T) {
	store, _ := NewNamespaceStore(tempNamespacePath(t))
	if _, ok := store.Get("missing"); ok {
		t.Error("expected not found")
	}
}

func TestNamespaceStoreInvalidJSON(t *testing.T) {
	path := tempNamespacePath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	if _, err := NewNamespaceStore(path); err == nil {
		t.Error("expected error on invalid JSON")
	}
}
