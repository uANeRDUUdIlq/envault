package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// DriftReporter formats and writes DriftReports.
type DriftReporter struct {
	format string // "text" | "json"
	out    io.Writer
}

// NewDriftReporter creates a reporter writing to out in the given format.
func NewDriftReporter(out io.Writer, format string) *DriftReporter {
	if format != "json" {
		format = "text"
	}
	return &DriftReporter{format: format, out: out}
}

// Write renders the report to the configured output.
func (r *DriftReporter) Write(report *DriftReport) error {
	switch r.format {
	case "json":
		return r.writeJSON(report)
	default:
		return r.writeText(report)
	}
}

func (r *DriftReporter) writeText(report *DriftReport) error {
	if !report.HasDrift() {
		_, err := fmt.Fprintf(r.out, "No drift detected in %s (checked at %s)\n",
			report.Environment, report.CheckedAt.Format("2006-01-02T15:04:05Z"))
		return err
	}
	lines := []string{
		fmt.Sprintf("Drift detected in %s (%d issue(s)):", report.Environment, len(report.Entries)),
	}
	lines = append(lines, report.SummaryLines()...)
	_, err := fmt.Fprintln(r.out, strings.Join(lines, "\n"))
	return err
}

func (r *DriftReporter) writeJSON(report *DriftReport) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
