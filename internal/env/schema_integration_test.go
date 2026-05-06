package env

import (
	"strings"
	"testing"
)

// TestSchemaRoundtrip parses a schema, then validates a full env map against it.
func TestSchemaRoundtrip(t *testing.T) {
	schemaText := `
DB_URL required # main database
PORT required pattern=^\d+$
DEBUG
SECRET_KEY required pattern=^[A-Za-z0-9]{16,}$
`
	s, err := ParseSchema(strings.NewReader(schemaText))
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	good := map[string]string{
		"DB_URL":     "postgres://localhost/mydb",
		"PORT":       "5432",
		"SECRET_KEY": "abcdefghij123456",
	}
	if !s.IsValid(good) {
		v := s.Validate(good)
		t.Fatalf("expected valid, got violations: %v", v)
	}
}

// TestSchemaMultipleViolations ensures all violations are returned at once.
func TestSchemaMultipleViolations(t *testing.T) {
	schemaText := "A required\nB required\nC required\n"
	s, err := ParseSchema(strings.NewReader(schemaText))
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	vars := map[string]string{}
	v := s.Validate(vars)
	if len(v) != 3 {
		t.Fatalf("expected 3 violations, got %d: %v", len(v), v)
	}
}

// TestSchemaWithParsedEnv uses the env parser together with schema validation.
func TestSchemaWithParsedEnv(t *testing.T) {
	envText := "DB_URL=postgres://localhost\nPORT=5432\n"
	vars, err := Parse(strings.NewReader(envText))
	if err != nil {
		t.Fatalf("parse env: %v", err)
	}

	s := &Schema{Fields: []SchemaField{
		{Key: "DB_URL", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	}}

	if !s.IsValid(vars) {
		v := s.Validate(vars)
		t.Fatalf("expected valid, got: %v", v)
	}
}
