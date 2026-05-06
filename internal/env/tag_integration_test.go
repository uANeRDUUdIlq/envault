package env

import (
	"testing"
)

// TestTagFilterRoundtrip verifies that tagging keys and filtering a parsed env
// file returns exactly the expected subset of variables.
func TestTagFilterRoundtrip(t *testing.T) {
	raw := `DB_HOST=localhost
DB_PORT=5432
DB_PASS=supersecret
APP_ENV=production
APP_PORT=8080
`
	vars, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	ts := NewTagStore()
	if err := ts.Add("database", []string{"DB_HOST", "DB_PORT", "DB_PASS"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := ts.Add("app", []string{"APP_ENV", "APP_PORT"}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	dbVars, err := ts.Filter("database", vars)
	if err != nil {
		t.Fatalf("Filter database: %v", err)
	}
	if len(dbVars) != 3 {
		t.Errorf("expected 3 database vars, got %d", len(dbVars))
	}
	if dbVars["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST mismatch: %q", dbVars["DB_HOST"])
	}

	appVars, err := ts.Filter("app", vars)
	if err != nil {
		t.Fatalf("Filter app: %v", err)
	}
	if len(appVars) != 2 {
		t.Errorf("expected 2 app vars, got %d", len(appVars))
	}
}

// TestTagSerializeSubset verifies that a filtered tag subset can be serialized
// back to a valid .env string.
func TestTagSerializeSubset(t *testing.T) {
	vars := map[string]string{
		"SECRET_KEY": "abc123",
		"PUBLIC_URL": "https://example.com",
		"INTERNAL": "only-internal",
	}

	ts := NewTagStore()
	_ = ts.Add("public", []string{"PUBLIC_URL"})

	subset, err := ts.Filter("public", vars)
	if err != nil {
		t.Fatalf("Filter: %v", err)
	}

	out := Serialize(subset)
	if out == "" {
		t.Fatal("expected non-empty serialized output")
	}

	reparsed, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	if reparsed["PUBLIC_URL"] != "https://example.com" {
		t.Errorf("PUBLIC_URL mismatch after roundtrip: %q", reparsed["PUBLIC_URL"])
	}
	if _, ok := reparsed["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should not appear in public subset")
	}
}
