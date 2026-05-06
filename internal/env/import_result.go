package env

import (
	"fmt"
	"sort"
	"strings"
)

// Summary returns a human-readable summary of the import operation.
func (r *ImportResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Imported: %d variable(s)", r.Imported)
	if len(r.Skipped) > 0 {
		sort.Strings(r.Skipped)
		fmt.Fprintf(&sb, ", Skipped: %d (%s)", len(r.Skipped), strings.Join(r.Skipped, ", "))
	}
	return sb.String()
}

// Merge applies the imported vars onto the provided base map and returns
// a new map containing both existing and newly imported keys.
func (r *ImportResult) Merge(base map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(r.Vars))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range r.Vars {
		out[k] = v
	}
	return out
}

// Keys returns a sorted list of all successfully imported key names.
func (r *ImportResult) Keys() []string {
	keys := make([]string, 0, len(r.Vars))
	for k := range r.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
