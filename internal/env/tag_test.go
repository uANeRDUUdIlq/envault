package env

import (
	"testing"
)

func TestTagAddAndGet(t *testing.T) {
	ts := NewTagStore()
	err := ts.Add("secrets", []string{"DB_PASS", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tag, err := ts.Get("secrets")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "secrets" {
		t.Errorf("expected name 'secrets', got %q", tag.Name)
	}
	if len(tag.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(tag.Keys))
	}
}

func TestTagAddEmptyName(t *testing.T) {
	ts := NewTagStore()
	err := ts.Add("", []string{"FOO"})
	if err == nil {
		t.Fatal("expected error for empty tag name")
	}
}

func TestTagAddNoKeys(t *testing.T) {
	ts := NewTagStore()
	err := ts.Add("empty", []string{})
	if err == nil {
		t.Fatal("expected error for empty keys slice")
	}
}

func TestTagGetNotFound(t *testing.T) {
	ts := NewTagStore()
	_, err := ts.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestTagDelete(t *testing.T) {
	ts := NewTagStore()
	_ = ts.Add("infra", []string{"HOST", "PORT"})
	if err := ts.Delete("infra"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ts.Delete("infra"); err == nil {
		t.Fatal("expected error deleting non-existent tag")
	}
}

func TestTagList(t *testing.T) {
	ts := NewTagStore()
	_ = ts.Add("z-tag", []string{"Z"})
	_ = ts.Add("a-tag", []string{"A"})
	_ = ts.Add("m-tag", []string{"M"})
	names := ts.List()
	if len(names) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(names))
	}
	if names[0] != "a-tag" || names[1] != "m-tag" || names[2] != "z-tag" {
		t.Errorf("tags not sorted: %v", names)
	}
}

func TestTagFilter(t *testing.T) {
	ts := NewTagStore()
	_ = ts.Add("db", []string{"DB_HOST", "DB_PASS"})
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"APP_ENV": "prod",
	}
	result, err := ts.Filter("db", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 filtered vars, got %d", len(result))
	}
	if _, ok := result["APP_ENV"]; ok {
		t.Error("APP_ENV should not be in filtered result")
	}
}

func TestTagFilterMissingTag(t *testing.T) {
	ts := NewTagStore()
	_, err := ts.Filter("nope", map[string]string{"A": "1"})
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestTagDedupKeys(t *testing.T) {
	ts := NewTagStore()
	_ = ts.Add("dup", []string{"FOO", "FOO", "BAR", " "})
	tag, _ := ts.Get("dup")
	if len(tag.Keys) != 2 {
		t.Errorf("expected 2 deduped keys, got %d: %v", len(tag.Keys), tag.Keys)
	}
}
