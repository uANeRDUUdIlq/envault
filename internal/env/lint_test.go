package env

import (
	"strings"
	"testing"
)

func TestLintCleanVars(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"PORT":         "8080",
		"_INTERNAL":    "value",
	}
	result := Lint(vars)
	if result.HasIssues() {
		t.Errorf("expected no issues, got: %s", result.Summary())
	}
}

func TestLintLowercaseKey(t *testing.T) {
	vars := map[string]string{"database_url": "postgres://localhost"}
	result := Lint(vars)
	if !result.HasIssues() {
		t.Fatal("expected issue for lowercase key")
	}
	if !containsIssueFor(result, "database_url", "uppercase") {
		t.Errorf("expected uppercase warning, got: %s", result.Summary())
	}
}

func TestLintEmptyValue(t *testing.T) {
	vars := map[string]string{"API_KEY": ""}
	result := Lint(vars)
	if !result.HasIssues() {
		t.Fatal("expected issue for empty value")
	}
	if !containsIssueFor(result, "API_KEY", "empty") {
		t.Errorf("expected empty value warning, got: %s", result.Summary())
	}
}

func TestLintWhitespaceValue(t *testing.T) {
	vars := map[string]string{"SECRET": "  myvalue  "}
	result := Lint(vars)
	if !result.HasIssues() {
		t.Fatal("expected issue for whitespace-padded value")
	}
	if !containsIssueFor(result, "SECRET", "whitespace") {
		t.Errorf("expected whitespace warning, got: %s", result.Summary())
	}
}

func TestLintKeyStartsWithDigit(t *testing.T) {
	vars := map[string]string{"1INVALID": "value"}
	result := Lint(vars)
	if !result.HasIssues() {
		t.Fatal("expected issue for key starting with digit")
	}
	if !containsIssueFor(result, "1INVALID", "start") {
		t.Errorf("expected start-character warning, got: %s", result.Summary())
	}
}

func TestLintKeyWithInvalidChar(t *testing.T) {
	vars := map[string]string{"MY-KEY": "value"}
	result := Lint(vars)
	if !result.HasIssues() {
		t.Fatal("expected issue for key with hyphen")
	}
	if !containsIssueFor(result, "MY-KEY", "invalid character") {
		t.Errorf("expected invalid character warning, got: %s", result.Summary())
	}
}

func TestLintSummaryNoIssues(t *testing.T) {
	result := &LintResult{}
	if result.Summary() != "no issues found" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestLintSummaryWithIssues(t *testing.T) {
	vars := map[string]string{"bad-key": ""}
	result := Lint(vars)
	summary := result.Summary()
	if !strings.Contains(summary, "issue(s) found") {
		t.Errorf("expected issue count in summary, got: %s", summary)
	}
}

// containsIssueFor checks whether any issue for the given key contains the substring.
func containsIssueFor(r *LintResult, key, substr string) bool {
	for _, issue := range r.Issues {
		if issue.Key == key && strings.Contains(issue.Message, substr) {
			return true
		}
	}
	return false
}
