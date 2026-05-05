package env

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	input := `
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" || entries[0].Value != "localhost" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParseSkipsCommentsAndBlanks(t *testing.T) {
	input := `# This is a comment

API_KEY=secret
# another comment
DEBUG=true
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestParseStripsQuotes(t *testing.T) {
	input := `SECRET="hello world"
TOKEN='bearer abc'
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", entries[0].Value)
	}
	if entries[1].Value != "bearer abc" {
		t.Errorf("expected 'bearer abc', got %q", entries[1].Value)
	}
}

func TestParseMissingEquals(t *testing.T) {
	input := "BADLINE\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParseEmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSerializeRoundtrip(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	output := Serialize(entries)
	parsed, err := Parse(strings.NewReader(output))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(parsed) != len(entries) {
		t.Fatalf("expected %d entries after roundtrip, got %d", len(entries), len(parsed))
	}
	for i, e := range entries {
		if parsed[i].Key != e.Key || parsed[i].Value != e.Value {
			t.Errorf("entry %d mismatch: want %+v, got %+v", i, e, parsed[i])
		}
	}
}
