package env

import (
	"bytes"
	"strings"
	"testing"
)

// TestDriftWithParsedEnv ensures DetectDrift works end-to-end with Parse output.
func TestDriftWithParsedEnv(t *testing.T) {
	baselineRaw := "DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=secret\n"
	currentRaw := "DB_HOST=prod.db\nDB_PORT=5432\nNEW_FLAG=true\n"

	baseVars, err := Parse(strings.NewReader(baselineRaw))
	if err != nil {
		t.Fatalf("parse baseline: %v", err)
	}
	currVars, err := Parse(strings.NewReader(currentRaw))
	if err != nil {
		t.Fatalf("parse current: %v", err)
	}

	report := DetectDrift("production", baseVars, currVars)

	if !report.HasDrift() {
		t.Fatal("expected drift to be detected")
	}

	statuses := map[string]string{}
	for _, e := range report.Entries {
		statuses[e.Key] = e.Status
	}

	if statuses["DB_HOST"] != "changed" {
		t.Errorf("expected DB_HOST changed, got %q", statuses["DB_HOST"])
	}
	if statuses["API_KEY"] != "removed" {
		t.Errorf("expected API_KEY removed, got %q", statuses["API_KEY"])
	}
	if statuses["NEW_FLAG"] != "added" {
		t.Errorf("expected NEW_FLAG added, got %q", statuses["NEW_FLAG"])
	}
}

// TestDriftReporterRoundtrip verifies text and JSON reporters both succeed.
func TestDriftReporterRoundtrip(t *testing.T) {
	base := map[string]string{"SECRET": "old", "KEEP": "same"}
	curr := map[string]string{"KEEP": "same", "ADDED": "new"}
	report := DetectDrift("staging", base, curr)

	for _, fmt := range []string{"text", "json"} {
		var buf bytes.Buffer
		rep := NewDriftReporter(&buf, fmt)
		if err := rep.Write(report); err != nil {
			t.Errorf("format %s: Write error: %v", fmt, err)
		}
		if buf.Len() == 0 {
			t.Errorf("format %s: expected non-empty output", fmt)
		}
	}
}
