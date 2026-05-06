package env

import (
	"strings"
	"testing"
)

func TestRenderWithComments(t *testing.T) {
	input := "# DB url\nDATABASE_URL=\nPORT=3000\n"
	tmpl, err := ParseTemplate(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	out := tmpl.Render(nil, RenderOptions{IncludeComments: true})
	if !strings.Contains(out, "# DB url") {
		t.Errorf("expected comment in output, got:\n%s", out)
	}
}

func TestRenderMasksSecrets(t *testing.T) {
	input := "API_KEY=\nHOST=localhost\n"
	tmpl, _ := ParseTemplate(input)
	env := map[string]string{"API_KEY": "supersecret", "HOST": "prod.example.com"}
	out := tmpl.Render(env, RenderOptions{MaskSecrets: true})
	if strings.Contains(out, "supersecret") {
		t.Error("secret value should be masked")
	}
	if !strings.Contains(out, "prod.example.com") {
		t.Error("non-secret value should be visible")
	}
}

func TestRenderUsesDefaults(t *testing.T) {
	input := "PORT=8080\nDEBUG=false\n"
	tmpl, _ := ParseTemplate(input)
	out := tmpl.Render(nil, RenderOptions{})
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected default PORT, got:\n%s", out)
	}
}

func TestRenderOverridesDefault(t *testing.T) {
	input := "PORT=8080\n"
	tmpl, _ := ParseTemplate(input)
	out := tmpl.Render(map[string]string{"PORT": "9090"}, RenderOptions{})
	if !strings.Contains(out, "PORT=9090") {
		t.Errorf("expected overridden PORT, got:\n%s", out)
	}
}

func TestExtraKeys(t *testing.T) {
	input := "PORT=8080\nDEBUG=false\n"
	tmpl, _ := ParseTemplate(input)
	env := map[string]string{"PORT": "8080", "UNKNOWN": "x", "ALSO_UNKNOWN": "y"}
	extra := tmpl.ExtraKeys(env)
	if len(extra) != 2 {
		t.Errorf("expected 2 extra keys, got %d: %v", len(extra), extra)
	}
}

func TestIsSecretCustomSuffixes(t *testing.T) {
	if !isSecret("MY_CREDENTIAL", []string{"CREDENTIAL"}) {
		t.Error("expected MY_CREDENTIAL to be secret with custom suffix")
	}
	if isSecret("HOST", []string{"CREDENTIAL"}) {
		t.Error("HOST should not be secret")
	}
}
