package env

import (
	"strings"
	"testing"
)

// TestSecretScanWithParsedEnv verifies the scanner works end-to-end with the env parser.
func TestSecretScanWithParsedEnv(t *testing.T) {
	raw := `APP_ENV=staging
DB_PASSWORD=supersecret
PORT=5432
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
`
	vars, err := Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scanner := NewSecretScanner()
	findings := scanner.Scan(vars)

	if len(findings) < 2 {
		t.Errorf("expected at least 2 findings, got %d: %+v", len(findings), findings)
	}

	keys := make(map[string]bool)
	for _, f := range findings {
		keys[f.Key] = true
	}
	if !keys["DB_PASSWORD"] {
		t.Error("expected DB_PASSWORD in findings")
	}
	if !keys["AWS_ACCESS_KEY_ID"] {
		t.Error("expected AWS_ACCESS_KEY_ID in findings")
	}
}

// TestSecretScanSummaryIntegration verifies summary output from a realistic scan.
func TestSecretScanSummaryIntegration(t *testing.T) {
	raw := `GITHUB_TOKEN=ghp_aBcDeFgHiJkLmNoPqRsTuVwXyZ1234567890
LOG_LEVEL=debug
`
	vars, err := Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scanner := NewSecretScanner()
	findings := scanner.Scan(vars)
	summary := ScanSummary(findings)

	if !strings.Contains(summary, "potential secret") {
		t.Errorf("expected secrets in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "GITHUB_TOKEN") {
		t.Errorf("expected GITHUB_TOKEN in summary, got: %s", summary)
	}
}
