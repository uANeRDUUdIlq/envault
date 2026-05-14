package env

import "sort"

// CompareResult holds the result of comparing two sets of env vars across environments.
type CompareResult struct {
	OnlyInA    map[string]string // keys present in A but not B
	OnlyInB    map[string]string // keys present in B but not A
	Different  map[string][2]string // keys in both but with different values: [A, B]
	Identical  []string // keys with identical values in both
}

// Compare performs a side-by-side comparison of two env var maps (e.g. staging vs production).
func Compare(a, b map[string]string) CompareResult {
	result := CompareResult{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Different: make(map[string][2]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				result.Identical = append(result.Identical, k)
			} else {
				result.Different[k] = [2]string{va, vb}
			}
		} else {
			result.OnlyInA[k] = va
		}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			result.OnlyInB[k] = vb
		}
	}

	sort.Strings(result.Identical)
	return result
}

// IsIdentical returns true when both env maps are completely equal.
func (r CompareResult) IsIdentical() bool {
	return len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 && len(r.Different) == 0
}

// SummaryLines returns a human-readable slice of summary lines for the comparison.
func (r CompareResult) SummaryLines() []string {
	var lines []string

	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	for _, k := range keys(r.OnlyInA) {
		lines = append(lines, "< "+k+"="+r.OnlyInA[k])
	}
	for _, k := range keys(r.OnlyInB) {
		lines = append(lines, "> "+k+"="+r.OnlyInB[k])
	}

	diffKeys := make([]string, 0, len(r.Different))
	for k := range r.Different {
		diffKeys = append(diffKeys, k)
	}
	sort.Strings(diffKeys)
	for _, k := range diffKeys {
		pair := r.Different[k]
		lines = append(lines, "~ "+k+": "+pair[0]+" -> "+pair[1])
	}

	return lines
}
