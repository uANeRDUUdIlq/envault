package env

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// LockEntry represents a lock held by a user on the vault.
type LockEntry struct {
	User      string    `json:"user"`
	AcquiredAt time.Time `json:"acquired_at"`
	TTL       int       `json:"ttl_seconds"`
}

// LockStore manages advisory locks for the vault.
type LockStore struct {
	path string
}

// NewLockStore creates a new LockStore backed by the given file path.
func NewLockStore(path string) *LockStore {
	return &LockStore{path: path}
}

// Acquire attempts to acquire the lock for the given user.
// Returns an error if the lock is already held and not expired.
func (ls *LockStore) Acquire(user string, ttlSeconds int) error {
	existing, err := ls.load()
	if err == nil {
		expiry := existing.AcquiredAt.Add(time.Duration(existing.TTL) * time.Second)
		if time.Now().Before(expiry) {
			return errors.New("vault is locked by " + existing.User)
		}
	}
	entry := LockEntry{
		User:       user,
		AcquiredAt: time.Now().UTC(),
		TTL:        ttlSeconds,
	}
	return ls.save(entry)
}

// Release removes the lock if it is held by the given user.
func (ls *LockStore) Release(user string) error {
	existing, err := ls.load()
	if err != nil {
		return errors.New("no lock to release")
	}
	if existing.User != user {
		return errors.New("lock is held by " + existing.User)
	}
	return os.Remove(ls.path)
}

// Current returns the active lock entry, or nil if no valid lock exists.
func (ls *LockStore) Current() *LockEntry {
	entry, err := ls.load()
	if err != nil {
		return nil
	}
	expiry := entry.AcquiredAt.Add(time.Duration(entry.TTL) * time.Second)
	if time.Now().After(expiry) {
		return nil
	}
	return entry
}

func (ls *LockStore) load() (*LockEntry, error) {
	data, err := os.ReadFile(ls.path)
	if err != nil {
		return nil, err
	}
	var entry LockEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (ls *LockStore) save(entry LockEntry) error {
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ls.path, data, 0600)
}
