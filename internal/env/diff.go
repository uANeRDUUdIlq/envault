package env

import "fmt"

// ChangeType represents the kind of change in a diff.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeUpdated ChangeType = "updated"
)

// Change represents a single key-level difference between two env maps.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// String returns a human-readable representation of the change.
func (c Change) String() string {
	switch c.Type {
	case ChangeAdded:
		return fmt.Sprintf("+ %s=%s", c.Key, c.NewVal)
	case ChangeRemoved:
		return fmt.Sprintf("- %s=%s", c.Key, c.OldVal)
	case ChangeUpdated:
		return fmt.Sprintf("~ %s: %s -> %s", c.Key, c.OldVal, c.NewVal)
	}
	return ""
}

// Diff computes the difference between two parsed env maps (old vs new).
// It returns a slice of Change entries describing what was added, removed, or updated.
func Diff(old, new map[string]string) []Change {
	var changes []Change

	for key, newVal := range new {
		oldVal, exists := old[key]
		if !exists {
			changes = append(changes, Change{Key: key, Type: ChangeAdded, NewVal: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: key, Type: ChangeUpdated, OldVal: oldVal, NewVal: newVal})
		}
	}

	for key, oldVal := range old {
		if _, exists := new[key]; !exists {
			changes = append(changes, Change{Key: key, Type: ChangeRemoved, OldVal: oldVal})
		}
	}

	return changes
}

// HasChanges returns true if there are any differences between old and new.
func HasChanges(old, new map[string]string) bool {
	return len(Diff(old, new)) > 0
}
