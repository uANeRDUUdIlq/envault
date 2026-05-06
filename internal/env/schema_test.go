package env

import (
	"strings"
	"testing"
)

func TestSchemaValidateRequiredMissing(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "DB_URL", Required: true}}}
	vars := map[string]string{}
	v := s.Validate(vars)
	if len(v) != 1 || v[0].Key != "DB_URL" {
		t.Fatalf("expected one violation for DB_URL, got %v", v)
	}
}

func TestSchemaValidateRequiredPresent(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "DB_URL", Required: true}}}
	vars := map[string]string{"DB_URL": "postgres://localhost"}
	if !s.IsValid(vars) {
		t.Fatal("expected valid")
	}
}

func TestSchemaValidatePatternMatch(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "PORT", Pattern: `^\d+$`}}}
	vars := map[string]string{"PORT": "8080"}
	if !s.IsValid(vars) {
		t.Fatal("expected valid")
	}
}

func TestSchemaValidatePatternMismatch(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "PORT", Pattern: `^\d+$`}}}
	vars := map[string]string{"PORT": "abc"}
	v := s.Validate(vars)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestSchemaValidateOptionalAbsent(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "LOG_LEVEL", Required: false}}}
	vars := map[string]string{}
	if !s.IsValid(vars) {
		t.Fatal("optional absent key should not violate")
	}
}

func TestParseSchemaBasic(t *testing.T) {
	input := `
# env schema
DB_URL required # database connection
PORT pattern=^\d+$
LOG_LEVEL
`
	s, err := ParseSchema(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(s.Fields))
	}
	if !s.Fields[0].Required {
		t.Error("DB_URL should be required")
	}
	if s.Fields[1].Pattern != `^\d+$` {
		t.Errorf("unexpected pattern: %q", s.Fields[1].Pattern)
	}
	if s.Fields[0].Description != "database connection" {
		t.Errorf("unexpected description: %q", s.Fields[0].Description)
	}
}

func TestParseSchemaUnknownToken(t *testing.T) {
	input := "KEY bogus\n"
	_, err := ParseSchema(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for unknown token")
	}
}

func TestParseSchemaEmpty(t *testing.T) {
	s, err := ParseSchema(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Fields) != 0 {
		t.Fatalf("expected 0 fields")
	}
}
