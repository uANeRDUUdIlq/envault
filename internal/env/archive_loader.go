package env

import (
	"fmt"
	"strings"
)

// ArchiveSummaryLine returns a human-readable summary line for an archive entry.
func ArchiveSummaryLine(e ArchiveEntry) string {
	note := e.Note
	if note == "" {
		note = "(no note)"
	}
	return fmt.Sprintf("%s — %d keys — %s — %s",
		e.Name,
		len(e.Vars),
		e.ArchivedAt.Format("2006-01-02 15:04:05"),
		note,
	)
}

// RestoreFromArchive applies the archived vars onto dst, optionally overwriting
// existing keys. Returns the number of keys applied.
func RestoreFromArchive(entry ArchiveEntry, dst map[string]string, overwrite bool) int {
	applied := 0
	for k, v := range entry.Vars {
		if _, exists := dst[k]; exists && !overwrite {
			continue
		}
		dst[k] = v
		applied++
	}
	return applied
}

// FilterArchive returns a new ArchiveEntry containing only the specified keys.
// If keys is empty, all vars are returned.
func FilterArchive(entry ArchiveEntry, keys []string) ArchiveEntry {
	if len(keys) == 0 {
		copy := make(map[string]string, len(entry.Vars))
		for k, v := range entry.Vars {
			copy[k] = v
		}
		return ArchiveEntry{Name: entry.Name, Vars: copy, ArchivedAt: entry.ArchivedAt, Note: entry.Note}
	}
	filtered := make(map[string]string)
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[strings.TrimSpace(k)] = struct{}{}
	}
	for k, v := range entry.Vars {
		if _, ok := set[k]; ok {
			filtered[k] = v
		}
	}
	return ArchiveEntry{Name: entry.Name, Vars: filtered, ArchivedAt: entry.ArchivedAt, Note: entry.Note}
}
