package env

import (
	"encoding/json"
	"fmt"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Author    string            `json:"author"`
	Message   string            `json:"message"`
	Vars      map[string]string `json:"vars"`
}

// SnapshotStore manages a list of snapshots for a vault.
type SnapshotStore struct {
	snapshots []Snapshot
}

// NewSnapshotStore creates an empty SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{}
}

// Add records a new snapshot with the given metadata and variable map.
func (s *SnapshotStore) Add(author, message string, vars map[string]string) Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	snap := Snapshot{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Author:    author,
		Message:   message,
		Vars:      copy,
	}
	s.snapshots = append(s.snapshots, snap)
	return snap
}

// List returns all snapshots in chronological order.
func (s *SnapshotStore) List() []Snapshot {
	out := make([]Snapshot, len(s.snapshots))
	copy(out, s.snapshots)
	return out
}

// Get returns the snapshot with the given ID, or an error if not found.
func (s *SnapshotStore) Get(id string) (Snapshot, error) {
	for _, snap := range s.snapshots {
		if snap.ID == id {
			return snap, nil
		}
	}
	return Snapshot{}, fmt.Errorf("snapshot %q not found", id)
}

// MarshalJSON serialises the store's snapshots.
func (s *SnapshotStore) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.snapshots)
}

// UnmarshalJSON restores snapshots from JSON.
func (s *SnapshotStore) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.snapshots)
}
