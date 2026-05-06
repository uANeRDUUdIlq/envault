package env

import "fmt"

// RestoreResult describes the outcome of restoring a snapshot.
type RestoreResult struct {
	SnapshotID string
	Author     string
	Message    string
	Diff       DiffResult
}

// Restorer applies a historical snapshot to produce a new variable map.
type Restorer struct {
	store *SnapshotStore
}

// NewRestorer creates a Restorer backed by the given SnapshotStore.
func NewRestorer(store *SnapshotStore) *Restorer {
	return &Restorer{store: store}
}

// Restore applies the snapshot identified by id on top of current,
// returning the restored variable map and a diff describing what changed.
func (r *Restorer) Restore(id string, current map[string]string) (map[string]string, RestoreResult, error) {
	snap, err := r.store.Get(id)
	if err != nil {
		return nil, RestoreResult{}, fmt.Errorf("restore: %w", err)
	}

	restored := make(map[string]string, len(snap.Vars))
	for k, v := range snap.Vars {
		restored[k] = v
	}

	diff := Diff(current, restored)

	result := RestoreResult{
		SnapshotID: snap.ID,
		Author:     snap.Author,
		Message:    snap.Message,
		Diff:       diff,
	}
	return restored, result, nil
}

// SummaryString returns a human-readable description of a RestoreResult.
func (rr RestoreResult) SummaryString() string {
	return fmt.Sprintf(
		"restored snapshot %s (by %s: %q) — added:%d removed:%d updated:%d",
		rr.SnapshotID, rr.Author, rr.Message,
		len(rr.Diff.Added), len(rr.Diff.Removed), len(rr.Diff.Updated),
	)
}
