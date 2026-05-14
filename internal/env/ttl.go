package env

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// TTLEntry represents a key with an expiration time.
type TTLEntry struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TTLStore manages per-key expiration (time-to-live) for env vars.
type TTLStore struct {
	path    string
	entries map[string]TTLEntry
}

// NewTTLStore loads or creates a TTL store at the given path.
func NewTTLStore(path string) (*TTLStore, error) {
	s := &TTLStore{path: path, entries: make(map[string]TTLEntry)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	return s, json.Unmarshal(data, &s.entries)
}

// Set assigns a TTL duration to a key.
func (s *TTLStore) Set(key string, ttl time.Duration) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	if ttl <= 0 {
		return errors.New("ttl must be positive")
	}
	s.entries[key] = TTLEntry{Key: key, ExpiresAt: time.Now().Add(ttl)}
	return s.save()
}

// Get returns the TTLEntry for a key, or an error if not found.
func (s *TTLStore) Get(key string) (TTLEntry, error) {
	e, ok := s.entries[key]
	if !ok {
		return TTLEntry{}, errors.New("ttl entry not found")
	}
	return e, nil
}

// IsExpired reports whether the key's TTL has elapsed.
func (s *TTLStore) IsExpired(key string) bool {
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// Remove deletes the TTL entry for a key.
func (s *TTLStore) Remove(key string) error {
	if _, ok := s.entries[key]; !ok {
		return errors.New("ttl entry not found")
	}
	delete(s.entries, key)
	return s.save()
}

// Expired returns all keys whose TTL has elapsed.
func (s *TTLStore) Expired() []string {
	now := time.Now()
	var out []string
	for k, e := range s.entries {
		if now.After(e.ExpiresAt) {
			out = append(out, k)
		}
	}
	return out
}

func (s *TTLStore) save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}
