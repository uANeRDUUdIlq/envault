package env

import (
	"strings"
	"testing"
)

func TestScanDetectsSecretKeyword(t *testing.T) {
	scanner := NewSecretScanner()
	vars := map[string]string{
		"DB_PASSWORD": "hunter2",
	}
	findings := scanner.Scan(vars)
	if len(findings) == 0 {
		t.Fatal("expected findings for DB_PASSWORD, got none")
	}
	if findings[0].Key != "DB_PASSWORD" {
		t.Errorf("expected key DB_PASSWORD, got %s", findings[0].Key)
	}
}

func TestScanDetectsAWSKey(t *testing.T) {
	scanner := NewSecretScanner()
	vars := map[string]string{
		"CLOUD_ID": "AKIAIOSFODNN7EXAMPLE",
	}
	findings := scanner.Scan(vars)
	if len(findings) == 0 {
		t.Fatal("expected findings for AWS-like key value")
	}
	if !strings.Contains(findings[0].Reason, "aws-access-key") {
		t.Errorf("unexpected reason: %s", findings[0].Reason)
	}
}

func TestScanCleanVars(t *testing.T) {
	scanner := NewSecretScanner()
	vars := map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
		"PORT":     "8080",
	}
	findings := scanner.Scan(vars)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestScanDetectsHighEntropyValue(t *testing.T) {
	scanner := NewSecretScanner()
	vars := map[string]string{
		"RANDOM_BLOB": "dGhpcyBpcyBhIHZlcnkgbG9uZyBiYXNlNjQgc3RyaW5n",
	}
	findings := scanner.Scan(vars)
	if len(findings) == 0 {
		t.Fatal("expected findings for high-entropy base64 value")
	}
}

func TestScanSummaryNoFindings(t *testing.T) {
	summary := ScanSummary(nil)
	if summary != "no secrets detected" {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestScanSummaryWithFindings(t *testing.T) {
	findings := []SecretFinding{
		{Key: "API_KEY", Reason: "key name contains sensitive keyword \"api_key\""},
		{Key: "TOKEN", Reason: "key name contains sensitive keyword \"token\""},
	}
	summary := ScanSummary(findings)
	if !strings.Contains(summary, "2 potential secret(s)") {
		t.Errorf("unexpected summary: %s", summary)
	}
	if !strings.Contains(summary, "API_KEY") {
		t.Errorf("expected API_KEY in summary")
	}
}

func TestScanMultipleKeywords(t *testing.T) {
	scanner := NewSecretScanner()
	vars := map[string]string{
		"STRIPE_SECRET_KEY": "sk_test_abc123",
	}
	findings := scanner.Scan(vars)
	// Should detect at least one finding (keyword match on "secret" or "key")
	if len(findings) == 0 {
		t.Fatal("expected at least one finding for STRIPE_SECRET_KEY")
	}
}
