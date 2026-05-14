package env_test

import (
	"testing"

	"github.com/envault/envault/internal/env"
)

func TestRedactRoundtripWithParser(t *testing.T) {
	raw := `APP_NAME=envault
DB_PASSWORD=supersecret
API_KEY=abc123xyz
PORT=8080
`
	vars, err := env.Parse(raw)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	r := env.NewRedactor(env.DefaultRedactOptions())
	redacted := r.Redact(vars)

	if redacted["APP_NAME"] != "envault" {
		t.Errorf("APP_NAME should not be redacted, got %q", redacted["APP_NAME"])
	}
	if redacted["PORT"] != "8080" {
		t.Errorf("PORT should not be redacted, got %q", redacted["PORT"])
	}
	if redacted["DB_PASSWORD"] == "supersecret" {
		t.Error("DB_PASSWORD should be redacted")
	}
	if redacted["API_KEY"] == "abc123xyz" {
		t.Error("API_KEY should be redacted")
	}
}

func TestRedactSerializeDoesNotLeakSecrets(t *testing.T) {
	vars := map[string]string{
		"SECRET_KEY": "topsecret",
		"REGION":     "us-east-1",
	}

	opts := env.DefaultRedactOptions()
	opts.PartialReveal = 2
	r := env.NewRedactor(opts)
	redacted := r.Redact(vars)

	serialised := env.Serialize(redacted)

	if contains(serialised, "topsecret") {
		t.Error("serialized output should not contain plain secret")
	}
	if !contains(serialised, "us-east-1") {
		t.Error("serialized output should contain plain non-secret value")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
