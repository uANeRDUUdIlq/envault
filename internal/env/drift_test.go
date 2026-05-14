package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestDriftNoChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "1", "B": "2"}
	r := DetectDrift("test", base, curr)
	if r.HasDrift() {
		t.Fatalf("expected no drift, got %+v", r.Entries)
	}
}

func TestDriftChanged(t *testing.T) {
	base := map[string]string{"A": "old"}
	curr := map[string]string{"A": "new"}
	r := DetectDrift("test", base, curr)
	if len(r.Entries) != 1 || r.Entries[0].Status != "changed" {
		t.Fatalf("expected 1 changed entry, got %+v", r.Entries)
	}
	if r.Entries[0].Baseline != "old" || r.Entries[0].Current != "new" {
		t.Errorf("unexpected values: %+v", r.Entries[0])
	}
}

func TestDriftAdded(t *testing.T) {
	base := map[string]string{}
	curr := map[string]string{"NEW_KEY": "val"}
	r := DetectDrift("test", base, curr)
	if len(r.Entries) != 1 || r.Entries[0].Status != "added" {
		t.Fatalf("expected 1 added entry, got %+v", r.Entries)
	}
}

func TestDriftRemoved(t *testing.T) {
	base := map[string]string{"GONE": "val"}
	curr := map[string]string{}
	r := DetectDrift("test", base, curr)
	if len(r.Entries) != 1 || r.Entries[0].Status != "removed" {
		t.Fatalf("expected 1 removed entry, got %+v", r.Entries)
	}
}

func TestDriftSummaryLines(t *testing.T) {
	base := map[string]string{"A": "1", "B": "old"}
	curr := map[string]string{"B": "new", "C": "3"}
	r := DetectDrift("prod", base, curr)
	lines := r.SummaryLines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 summary lines, got %d: %v", len(lines), lines)
	}
}

func TestDriftReporterText(t *testing.T) {
	base := map[string]string{"X": "a"}
	curr := map[string]string{"X": "b"}
	report := DetectDrift("staging", base, curr)
	var buf bytes.Buffer
	rep := NewDriftReporter(&buf, "text")
	if err := rep.Write(report); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected env name in output, got: %s", out)
	}
	if !strings.Contains(out, "changed") || !strings.Contains(out, "X") {
		t.Errorf("expected key info in output, got: %s", out)
	}
}

func TestDriftReporterJSON(t *testing.T) {
	base := map[string]string{}
	curr := map[string]string{}
	report := DetectDrift("dev", base, curr)
	var buf bytes.Buffer
	rep := NewDriftReporter(&buf, "json")
	if err := rep.Write(report); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"Environment\"") {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestDriftReporterNoDrift(t *testing.T) {
	report := DetectDrift("ci", map[string]string{"K": "v"}, map[string]string{"K": "v"})
	var buf bytes.Buffer
	rep := NewDriftReporter(&buf, "text")
	_ = rep.Write(report)
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}
