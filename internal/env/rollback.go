package env

import (
	"errors"
	"fmt"
	"time"
)

// RollbackEntry records a single rollback operation.
type RollbackEntry struct {
	Timestamp   time.Time         `json:"timestamp"`
	FromVersion string            `json:"from_version"`
	ToVersion   string            `json:"to_version"`
	Namespace   string            `json:"namespace"`
	Keys        map[string]string `json:"keys"`
	RestoredBy  string            `json:"restored_by"`
}

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Entry      RollbackEntry
	Applied    []string
	Skipped    []string
	PrevValues map[string]string
}

// RollbackSummary returns a human-readable summary of the rollback.
func (r *RollbackResult) RollbackSummary() string {
	return fmt.Sprintf(
		"rollback to %s: applied=%d skipped=%d",
		r.Entry.ToVersion, len(r.Applied), len(r.Skipped),
	)
}

// Rollbacker applies snapshot-based rollbacks to a variable map.
type Rollbacker struct {
	snapshots *SnapshotStore
	vault     *Vault
}

// NewRollbacker creates a Rollbacker backed by the given SnapshotStore and Vault.
func NewRollbacker(snapshots *SnapshotStore, vault *Vault) *Rollbacker {
	return &Rollbacker{snapshots: snapshots, vault: vault}
}

// Rollback reverts the vault's variables to the state captured in the named snapshot.
// Keys present in the snapshot replace current values; keys absent from the snapshot
// are left untouched unless purge is true, in which case they are removed.
func (rb *Rollbacker) Rollback(snapshotID, user string, purge bool) (*RollbackResult, error) {
	if snapshotID == "" {
		return nil, errors.New("rollback: snapshot ID must not be empty")
	}
	if user == "" {
		return nil, errors.New("rollback: user must not be empty")
	}

	snap, err := rb.snapshots.Get(snapshotID)
	if err != nil {
		return nil, fmt.Errorf("rollback: snapshot %q not found: %w", snapshotID, err)
	}

	current := rb.vault.All()
	prev := make(map[string]string, len(current))
	for k, v := range current {
		prev[k] = v
	}

	applied := []string{}
	skipped := []string{}

	for k, v := range snap.Vars {
		rb.vault.Set(k, v)
		applied = append(applied, k)
	}

	if purge {
		for k := range current {
			if _, inSnap := snap.Vars[k]; !inSnap {
				rb.vault.Delete(k)
				applied = append(applied, k)
			}
		}
	} else {
		for k := range current {
			if _, inSnap := snap.Vars[k]; !inSnap {
				skipped = append(skipped, k)
			}
		}
	}

	entry := RollbackEntry{
		Timestamp:  time.Now().UTC(),
		ToVersion:  snapshotID,
		Namespace:  snap.Namespace,
		Keys:       snap.Vars,
		RestoredBy: user,
	}

	return &RollbackResult{
		Entry:      entry,
		Applied:    applied,
		Skipped:    skipped,
		PrevValues: prev,
	}, nil
}
