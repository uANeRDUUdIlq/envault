package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// SecretReport holds the result of a secret scan for reporting purposes.
type SecretReport struct {
	ScannedAt time.Time       `json:"scanned_at"`
	TotalKeys int             `json:"total_keys"`
	Findings  []SecretFinding `json:"findings"`
	Clean     bool            `json:"clean"`
}

// NewSecretReport builds a report from a set of vars using the default scanner.
func NewSecretReport(vars map[string]string) *SecretReport {
	scanner := NewSecretScanner()
	findings := scanner.Scan(vars)
	return &SecretReport{
		ScannedAt: time.Now().UTC(),
		TotalKeys: len(vars),
		Findings:  findings,
		Clean:     len(findings) == 0,
	}
}

// WriteText writes a human-readable report to w.
func (r *SecretReport) WriteText(w io.Writer) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Secret Scan Report — %s\n", r.ScannedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "Keys scanned : %d\n", r.TotalKeys)
	fmt.Fprintf(&sb, "Findings     : %d\n", len(r.Findings))
	if r.Clean {
		sb.WriteString("Status       : CLEAN\n")
	} else {
		sb.WriteString("Status       : WARNINGS\n")
		for _, f := range r.Findings {
			fmt.Fprintf(&sb, "  [%s] %s\n", f.Key, f.Reason)
		}
	}
	_, err := io.WriteString(w, sb.String())
	return err
}

// WriteJSON writes the report as JSON to w.
func (r *SecretReport) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
