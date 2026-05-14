package env

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"
)

// RollbackHistory persists a log of rollback operations to a JSON file.
type RollbackHistory struct {
	path    string
	entries []RollbackEntry
}

// NewRollbackHistory creates a RollbackHistory backed by the given file path.
func NewRollbackHistory(path string) *RollbackHistory {
	return &RollbackHistory{path: path}
}

func (h *RollbackHistory) load() error {
	data, err := os.ReadFile(h.path)
	if errors.Is(err, os.ErrNotExist) {
		h.entries = nil
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}

func (h *RollbackHistory) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o600)
}

// Append records a new rollback entry.
func (h *RollbackHistory) Append(entry RollbackEntry) error {
	if err := h.load(); err != nil {
		return err
	}
	h.entries = append(h.entries, entry)
	return h.save()
}

// All returns all rollback entries sorted by timestamp descending.
func (h *RollbackHistory) All() ([]RollbackEntry, error) {
	if err := h.load(); err != nil {
		return nil, err
	}
	out := make([]RollbackEntry, len(h.entries))
	copy(out, h.entries)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp.After(out[j].Timestamp)
	})
	return out, nil
}

// Since returns entries recorded after the given time.
func (h *RollbackHistory) Since(t time.Time) ([]RollbackEntry, error) {
	all, err := h.All()
	if err != nil {
		return nil, err
	}
	var out []RollbackEntry
	for _, e := range all {
		if e.Timestamp.After(t) {
			out = append(out, e)
		}
	}
	return out, nil
}
