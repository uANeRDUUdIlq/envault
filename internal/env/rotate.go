package env

import (
	"fmt"
	"time"
)

// RotationRecord holds metadata about a single key rotation event.
type RotationRecord struct {
	Timestamp time.Time
	OldKeys   []string
	NewKeys   []string
	RemovedAt map[string]time.Time
}

// Rotator manages key rotation for environment variable sets.
type Rotator struct {
	history []RotationRecord
}

// NewRotator creates a new Rotator instance.
func NewRotator() *Rotator {
	return &Rotator{}
}

// Rotate takes an old and new set of env vars, records the rotation event,
// and returns a merged map where new values take precedence.
func (r *Rotator) Rotate(old, next map[string]string) (map[string]string, RotationRecord) {
	record := RotationRecord{
		Timestamp: time.Now().UTC(),
		RemovedAt: make(map[string]time.Time),
	}

	for k := range old {
		record.OldKeys = append(record.OldKeys, k)
		if _, exists := next[k]; !exists {
			record.RemovedAt[k] = record.Timestamp
		}
	}
	for k := range next {
		record.NewKeys = append(record.NewKeys, k)
	}

	merged := make(map[string]string, len(next))
	for k, v := range next {
		merged[k] = v
	}

	r.history = append(r.history, record)
	return merged, record
}

// History returns all rotation records.
func (r *Rotator) History() []RotationRecord {
	return r.history
}

// SummaryString returns a human-readable summary of a rotation record.
func SummaryString(rec RotationRecord) string {
	return fmt.Sprintf(
		"[%s] rotated: %d old keys -> %d new keys, %d removed",
		rec.Timestamp.Format(time.RFC3339),
		len(rec.OldKeys),
		len(rec.NewKeys),
		len(rec.RemovedAt),
	)
}
