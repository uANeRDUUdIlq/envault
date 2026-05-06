package env

import (
	"strings"
	"testing"
)

// TestTemplateRoundtrip verifies that a template can be parsed, filled with
// env values, rendered, and then re-parsed to produce the same keys.
func TestTemplateRoundtrip(t *testing.T) {
	raw := `# App secret
SECRET_KEY=
# Server port
PORT=8080
DEBUG=false
`
	tmpl, err := ParseTemplate(raw)
	if err != nil {
		t.Fatalf("ParseTemplate: %v", err)
	}

	env := map[string]string{"SECRET_KEY": "abc123", "PORT": "9090"}
	filled := tmpl.Fill(env)

	rendered := tmpl.Render(filled, RenderOptions{IncludeComments: false})

	// Re-parse as a plain .env
	parsed, err := Parse(strings.NewReader(rendered))
	if err != nil {
		t.Fatalf("Parse rendered output: %v", err)
	}
	if parsed["SECRET_KEY"] != "abc123" {
		t.Errorf("SECRET_KEY mismatch: %q", parsed["SECRET_KEY"])
	}
	if parsed["PORT"] != "9090" {
		t.Errorf("PORT mismatch: %q", parsed["PORT"])
	}
	if parsed["DEBUG"] != "false" {
		t.Errorf("DEBUG mismatch: %q", parsed["DEBUG"])
	}
}

// TestTemplateValidateAndFillCycle exercises the full validate → fill cycle.
func TestTemplateValidateAndFillCycle(t *testing.T) {
	raw := "DB_URL=\nCACHE_URL=redis://localhost\nTIMEOUT=30\n"
	tmpl, _ := ParseTemplate(raw)

	// Initially missing DB_URL
	missing := tmpl.Validate(map[string]string{})
	if len(missing) != 1 || missing[0] != "DB_URL" {
		t.Fatalf("expected [DB_URL] missing, got %v", missing)
	}

	// After providing DB_URL, no missing keys
	env := map[string]string{"DB_URL": "postgres://localhost/mydb"}
	missing = tmpl.Validate(env)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}

	filled := tmpl.Fill(env)
	if filled["CACHE_URL"] != "redis://localhost" {
		t.Errorf("expected default CACHE_URL, got %q", filled["CACHE_URL"])
	}
	if filled["TIMEOUT"] != "30" {
		t.Errorf("expected default TIMEOUT=30, got %q", filled["TIMEOUT"])
	}
}
