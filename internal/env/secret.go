package env

import (
	"fmt"
	"regexp"
	"strings"
)

// SecretPattern holds a compiled pattern and its label.
type SecretPattern struct {
	Label   string
	Pattern *regexp.Regexp
}

// SecretScanner scans env vars for keys or values that look like secrets.
type SecretScanner struct {
	patterns []SecretPattern
}

// SecretFinding represents a detected secret in an env map.
type SecretFinding struct {
	Key    string
	Reason string
}

var defaultPatterns = []SecretPattern{
	{Label: "high-entropy-value", Pattern: regexp.MustCompile(`[A-Za-z0-9+/]{32,}={0,2}`)},
	{Label: "aws-access-key", Pattern: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{Label: "private-key-header", Pattern: regexp.MustCompile(`-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----`)},
	{Label: "github-token", Pattern: regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36}`)},
}

var secretKeywords = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"private", "credential", "auth", "access_key", "signing",
}

// NewSecretScanner returns a scanner with default patterns.
func NewSecretScanner() *SecretScanner {
	return &SecretScanner{patterns: defaultPatterns}
}

// Scan checks the given vars for secret-like keys or values.
func (s *SecretScanner) Scan(vars map[string]string) []SecretFinding {
	var findings []SecretFinding
	for k, v := range vars {
		lower := strings.ToLower(k)
		for _, kw := range secretKeywords {
			if strings.Contains(lower, kw) {
				findings = append(findings, SecretFinding{
					Key:    k,
					Reason: fmt.Sprintf("key name contains sensitive keyword %q", kw),
				})
				break
			}
		}
		for _, p := range s.patterns {
			if p.Pattern.MatchString(v) {
				findings = append(findings, SecretFinding{
					Key:    k,
					Reason: fmt.Sprintf("value matches pattern %q", p.Label),
				})
				break
			}
		}
	}
	return findings
}

// ScanSummary returns a human-readable summary of findings.
func ScanSummary(findings []SecretFinding) string {
	if len(findings) == 0 {
		return "no secrets detected"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d potential secret(s) detected:\n", len(findings))
	for _, f := range findings {
		fmt.Fprintf(&sb, "  [%s] %s\n", f.Key, f.Reason)
	}
	return strings.TrimRight(sb.String(), "\n")
}
