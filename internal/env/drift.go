package env

import (
	"fmt"
	"sort"
	"time"
)

// DriftEntry describes a single variable that has drifted from its baseline.
type DriftEntry struct {
	Key      string
	Baseline string
	Current  string
	Status   string // "changed", "added", "removed"
}

// DriftReport holds the result of a drift detection run.
type DriftReport struct {
	Environment string
	CheckedAt   time.Time
	Entries     []DriftEntry
}

// HasDrift returns true if any entries were detected.
func (r *DriftReport) HasDrift() bool {
	return len(r.Entries) > 0
}

// SummaryLines returns human-readable lines describing the drift.
func (r *DriftReport) SummaryLines() []string {
	lines := make([]string, 0, len(r.Entries))
	for _, e := range r.Entries {
		switch e.Status {
		case "changed":
			lines = append(lines, fmt.Sprintf("~ %s: baseline=%q current=%q", e.Key, e.Baseline, e.Current))
		case "added":
			lines = append(lines, fmt.Sprintf("+ %s (not in baseline)", e.Key))
		case "removed":
			lines = append(lines, fmt.Sprintf("- %s (missing from current)", e.Key))
		}
	}
	return lines
}

// DetectDrift compares current vars against a baseline snapshot.
// baseline and current are maps of key→value.
func DetectDrift(env string, baseline, current map[string]string) *DriftReport {
	report := &DriftReport{
		Environment: env,
		CheckedAt:   time.Now().UTC(),
	}

	allKeys := make(map[string]struct{})
	for k := range baseline {
		allKeys[k] = struct{}{}
	}
	for k := range current {
		allKeys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, k := range sorted {
		bv, inBase := baseline[k]
		cv, inCurr := current[k]
		switch {
		case inBase && inCurr && bv != cv:
			report.Entries = append(report.Entries, DriftEntry{Key: k, Baseline: bv, Current: cv, Status: "changed"})
		case !inBase && inCurr:
			report.Entries = append(report.Entries, DriftEntry{Key: k, Current: cv, Status: "added"})
		case inBase && !inCurr:
			report.Entries = append(report.Entries, DriftEntry{Key: k, Baseline: bv, Status: "removed"})
		}
	}

	return report
}
