package env

import (
	"testing"
)

const sampleTemplate = `
# Database connection URL
DATABASE_URL=

# Optional port with default
PORT=8080

SECRET_KEY=
DEBUG=false
`

func TestParseTemplateBasic(t *testing.T) {
	tmpl, err := ParseTemplate(sampleTemplate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Vars) != 4 {
		t.Fatalf("expected 4 vars, got %d", len(tmpl.Vars))
	}
}

func TestParseTemplateRequired(t *testing.T) {
	tmpl, err := ParseTemplate(sampleTemplate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requiredKeys := map[string]bool{}
	for _, v := range tmpl.Vars {
		if v.Required {
			requiredKeys[v.Key] = true
		}
	}
	if !requiredKeys["DATABASE_URL"] {
		t.Error("expected DATABASE_URL to be required")
	}
	if !requiredKeys["SECRET_KEY"] {
		t.Error("expected SECRET_KEY to be required")
	}
	if requiredKeys["PORT"] {
		t.Error("PORT should not be required")
	}
}

func TestParseTemplateComment(t *testing.T) {
	tmpl, err := ParseTemplate(sampleTemplate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Vars[0].Comment != "Database connection URL" {
		t.Errorf("unexpected comment: %q", tmpl.Vars[0].Comment)
	}
}

func TestParseTemplateInvalidLine(t *testing.T) {
	_, err := ParseTemplate("BADLINE")
	if err == nil {
		t.Fatal("expected error for line missing '='")
	}
}

func TestTemplateValidateMissingKeys(t *testing.T) {
	tmpl, _ := ParseTemplate(sampleTemplate)
	env := map[string]string{"PORT": "9090", "DEBUG": "true"}
	missing := tmpl.Validate(env)
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing, got %d: %v", len(missing), missing)
	}
}

func TestTemplateValidateAllPresent(t *testing.T) {
	tmpl, _ := ParseTemplate(sampleTemplate)
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"SECRET_KEY":   "s3cr3t",
	}
	missing := tmpl.Validate(env)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got: %v", missing)
	}
}

func TestTemplateFillDefaults(t *testing.T) {
	tmpl, _ := ParseTemplate(sampleTemplate)
	env := map[string]string{"DATABASE_URL": "postgres://localhost/db", "SECRET_KEY": "abc"}
	filled := tmpl.Fill(env)
	if filled["PORT"] != "8080" {
		t.Errorf("expected PORT=8080 from default, got %q", filled["PORT"])
	}
	if filled["DATABASE_URL"] != "postgres://localhost/db" {
		t.Errorf("env override not applied")
	}
}
