package store

import (
	"encoding/base64"
	"fmt"

	"github.com/envault/envault/internal/crypto"
)

// Manager combines a Store with an Encryptor to provide
// high-level put/get operations on plaintext secrets.
type Manager struct {
	store     *Store
	encryptor *crypto.Encryptor
}

// NewManager creates a Manager backed by the given store path and encryptor.
func NewManager(storePath string, enc *crypto.Encryptor) (*Manager, error) {
	s, err := New(storePath)
	if err != nil {
		return nil, fmt.Errorf("store: open: %w", err)
	}
	return &Manager{store: s, encryptor: enc}, nil
}

// SetSecret encrypts plaintext and stores it under name.
func (m *Manager) SetSecret(name, plaintext string) error {
	ciphertext, err := m.encryptor.Encrypt([]byte(plaintext))
	if err != nil {
		return fmt.Errorf("store: encrypt %q: %w", name, err)
	}
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	m.store.Put(name, encoded)
	return m.store.Save()
}

// GetSecret retrieves and decrypts the secret stored under name.
func (m *Manager) GetSecret(name string) (string, error) {
	entry, err := m.store.Get(name)
	if err != nil {
		return "", fmt.Errorf("store: get %q: %w", name, err)
	}
	ciphertext, err := base64.StdEncoding.DecodeString(entry.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("store: decode %q: %w", name, err)
	}
	plaintext, err := m.encryptor.Decrypt(ciphertext)
	if err != nil {
		return "", fmt.Errorf("store: decrypt %q: %w", name, err)
	}
	return string(plaintext), nil
}

// RemoveSecret deletes the secret and persists the store.
func (m *Manager) RemoveSecret(name string) error {
	if err := m.store.Delete(name); err != nil {
		return fmt.Errorf("store: delete %q: %w", name, err)
	}
	return m.store.Save()
}

// ListNames returns all stored secret names.
func (m *Manager) ListNames() []string {
	names := make([]string, 0, len(m.store.Secrets))
	for k := range m.store.Secrets {
		names = append(names, k)
	}
	return names
}
