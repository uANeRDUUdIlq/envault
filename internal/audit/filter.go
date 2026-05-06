package audit

import "time"

// Filter holds optional criteria for querying audit entries.
type Filter struct {
	Event   EventType
	User    string
	File    string
	Since   time.Time
	Success *bool
}

// Apply returns entries from the provided slice that match all non-zero
// fields of the Filter.
func (f Filter) Apply(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if f.Event != "" && e.Event != f.Event {
			continue
		}
		if f.User != "" && e.User != f.User {
			continue
		}
		if f.File != "" && e.File != f.File {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		if f.Success != nil && e.Success != *f.Success {
			continue
		}
		out = append(out, e)
	}
	return out
}

// Query reads all entries from the logger and applies the filter.
func (l *Logger) Query(f Filter) ([]Entry, error) {
	all, err := l.ReadAll()
	if err != nil {
		return nil, err
	}
	return f.Apply(all), nil
}
