package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
	ErrStoreNotFound  = errors.New("store file not found")
)

// SecretEntry holds an encrypted secret and its metadata.
type SecretEntry struct {
	Name      string    `json:"name"`
	Ciphertext string   `json:"ciphertext"` // base64-encoded age ciphertext
	UpdatedAt time.Time `json:"updated_at"`
}

// Store represents the on-disk secrets store.
type Store struct {
	Version int                     `json:"version"`
	Secrets map[string]*SecretEntry `json:"secrets"`
	path    string
}

// New loads an existing store from path, or creates a new empty one.
func New(path string) (*Store, error) {
	s := &Store{
		Version: 1,
		Secrets: make(map[string]*SecretEntry),
		path:    path,
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	s.path = path
	return s, nil
}

// Put inserts or updates a secret entry.
func (s *Store) Put(name, ciphertext string) {
	s.Secrets[name] = &SecretEntry{
		Name:       name,
		Ciphertext: ciphertext,
		UpdatedAt:  time.Now().UTC(),
	}
}

// Get retrieves a secret entry by name.
func (s *Store) Get(name string) (*SecretEntry, error) {
	entry, ok := s.Secrets[name]
	if !ok {
		return nil, ErrSecretNotFound
	}
	return entry, nil
}

// Delete removes a secret by name.
func (s *Store) Delete(name string) error {
	if _, ok := s.Secrets[name]; !ok {
		return ErrSecretNotFound
	}
	delete(s.Secrets, name)
	return nil
}

// Save persists the store to disk as JSON.
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
