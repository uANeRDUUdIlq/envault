package env

import (
	"testing"
)

func TestRedactMasksSecretKeys(t *testing.T) {
	r := NewRedactor(DefaultRedactOptions())
	vars := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "envault",
	}
	out := r.Redact(vars)
	if out["DB_PASSWORD"] != "***********" {
		t.Errorf("expected masked password, got %q", out["DB_PASSWORD"])
	}
	if out["APP_NAME"] != "envault" {
		t.Errorf("expected plain value, got %q", out["APP_NAME"])
	}
}

func TestRedactPartialReveal(t *testing.T) {
	opts := DefaultRedactOptions()
	opts.PartialReveal = 3
	r := NewRedactor(opts)
	vars := map[string]string{"API_KEY": "abcdefgh"}
	out := r.Redact(vars)
	if out["API_KEY"] != "abc*****" {
		t.Errorf("expected partial reveal, got %q", out["API_KEY"])
	}
}

func TestRedactExplicitKeys(t *testing.T) {
	opts := DefaultRedactOptions()
	opts.Keys = []string{"MY_CUSTOM_VAR"}
	r := NewRedactor(opts)
	vars := map[string]string{
		"MY_CUSTOM_VAR": "hidden",
		"OTHER_VAR":     "visible",
	}
	out := r.Redact(vars)
	if out["MY_CUSTOM_VAR"] != "******" {
		t.Errorf("expected masked custom var, got %q", out["MY_CUSTOM_VAR"])
	}
	if out["OTHER_VAR"] != "visible" {
		t.Errorf("expected plain value, got %q", out["OTHER_VAR"])
	}
}

func TestRedactEmptyValue(t *testing.T) {
	r := NewRedactor(DefaultRedactOptions())
	vars := map[string]string{"SECRET_TOKEN": ""}
	out := r.Redact(vars)
	if out["SECRET_TOKEN"] != "" {
		t.Errorf("expected empty string for empty value, got %q", out["SECRET_TOKEN"])
	}
}

func TestRedactedKeys(t *testing.T) {
	r := NewRedactor(DefaultRedactOptions())
	vars := map[string]string{
		"AUTH_TOKEN": "tok",
		"PLAIN_VAR":  "val",
		"DB_SECRET":  "s3cr3t",
	}
	keys := r.RedactedKeys(vars)
	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}
	if !keySet["AUTH_TOKEN"] {
		t.Error("expected AUTH_TOKEN in redacted keys")
	}
	if !keySet["DB_SECRET"] {
		t.Error("expected DB_SECRET in redacted keys")
	}
	if keySet["PLAIN_VAR"] {
		t.Error("expected PLAIN_VAR not in redacted keys")
	}
}

func TestRedactCaseInsensitiveExplicit(t *testing.T) {
	opts := DefaultRedactOptions()
	opts.Keys = []string{"my_var"}
	r := NewRedactor(opts)
	vars := map[string]string{"MY_VAR": "value"}
	out := r.Redact(vars)
	if out["MY_VAR"] != "*****" {
		t.Errorf("expected masked value, got %q", out["MY_VAR"])
	}
}
