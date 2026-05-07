package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// PinEntry records a pinned snapshot for a named environment.
type PinEntry struct {
	Environment string    `json:"environment"`
	SnapshotID  string    `json:"snapshot_id"`
	PinnedBy    string    `json:"pinned_by"`
	PinnedAt    time.Time `json:"pinned_at"`
	Note        string    `json:"note,omitempty"`
}

// PinStore persists pinned snapshot references.
type PinStore struct {
	path string
	pins map[string]PinEntry // keyed by environment
}

// NewPinStore loads or creates a pin store at the given path.
func NewPinStore(path string) (*PinStore, error) {
	ps := &PinStore{path: path, pins: make(map[string]PinEntry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ps, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pin store read: %w", err)
	}
	if err := json.Unmarshal(data, &ps.pins); err != nil {
		return nil, fmt.Errorf("pin store parse: %w", err)
	}
	return ps, nil
}

// Pin sets the pinned snapshot for an environment.
func (ps *PinStore) Pin(env, snapshotID, pinnedBy, note string) error {
	if env == "" {
		return errors.New("environment name required")
	}
	if snapshotID == "" {
		return errors.New("snapshot ID required")
	}
	ps.pins[env] = PinEntry{
		Environment: env,
		SnapshotID:  snapshotID,
		PinnedBy:    pinnedBy,
		PinnedAt:    time.Now().UTC(),
		Note:        note,
	}
	return ps.save()
}

// Get returns the pin entry for an environment, or false if not pinned.
func (ps *PinStore) Get(env string) (PinEntry, bool) {
	e, ok := ps.pins[env]
	return e, ok
}

// Unpin removes the pin for an environment.
func (ps *PinStore) Unpin(env string) error {
	if _, ok := ps.pins[env]; !ok {
		return fmt.Errorf("no pin found for environment %q", env)
	}
	delete(ps.pins, env)
	return ps.save()
}

// List returns all current pin entries.
func (ps *PinStore) List() []PinEntry {
	out := make([]PinEntry, 0, len(ps.pins))
	for _, e := range ps.pins {
		out = append(out, e)
	}
	return out
}

func (ps *PinStore) save() error {
	data, err := json.MarshalIndent(ps.pins, "", "  ")
	if err != nil {
		return fmt.Errorf("pin store marshal: %w", err)
	}
	return os.WriteFile(ps.path, data, 0o600)
}
