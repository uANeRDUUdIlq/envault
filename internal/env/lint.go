package env

import (
	"fmt"
	"strings"
	"unicode"
)

// LintIssue represents a single linting problem found in a set of env vars.
type LintIssue struct {
	Key     string
	Message string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// LintResult holds all issues found during a lint pass.
type LintResult struct {
	Issues []LintIssue
}

// HasIssues returns true when at least one issue was found.
func (r *LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

// Summary returns a human-readable summary of all issues.
func (r *LintResult) Summary() string {
	if !r.HasIssues() {
		return "no issues found"
	}
	lines := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		lines = append(lines, "  - "+issue.String())
	}
	return fmt.Sprintf("%d issue(s) found:\n%s", len(r.Issues), strings.Join(lines, "\n"))
}

// Lint checks a map of env vars for common problems and returns a LintResult.
func Lint(vars map[string]string) *LintResult {
	result := &LintResult{}

	for key, value := range vars {
		// Check for empty keys (shouldn't happen after parsing, but defensive)
		if key == "" {
			result.Issues = append(result.Issues, LintIssue{Key: "(empty)", Message: "key must not be empty"})
			continue
		}

		// Keys should be uppercase
		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, LintIssue{Key: key, Message: "key should be uppercase"})
		}

		// Keys should start with a letter or underscore
		if first := rune(key[0]); !unicode.IsLetter(first) && first != '_' {
			result.Issues = append(result.Issues, LintIssue{Key: key, Message: "key should start with a letter or underscore"})
		}

		// Keys should only contain letters, digits, and underscores
		for _, ch := range key {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
				result.Issues = append(result.Issues, LintIssue{Key: key, Message: "key contains invalid character: " + string(ch)})
				break
			}
		}

		// Warn on empty values
		if strings.TrimSpace(value) == "" {
			result.Issues = append(result.Issues, LintIssue{Key: key, Message: "value is empty"})
		}

		// Warn on values with leading/trailing whitespace
		if value != strings.TrimSpace(value) {
			result.Issues = append(result.Issues, LintIssue{Key: key, Message: "value has leading or trailing whitespace"})
		}
	}

	return result
}
