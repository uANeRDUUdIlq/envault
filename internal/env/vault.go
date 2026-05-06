package env

import (
	"fmt"

	"github.com/envault/envault/internal/crypto"
)

// Vault manages encryption and decryption of .env file contents.
type Vault struct {
	encryptor *crypto.Encryptor
}

// NewVault creates a Vault backed by the given Encryptor.
func NewVault(enc *crypto.Encryptor) *Vault {
	return &Vault{encryptor: enc}
}

// Encrypt parses the raw env content, serializes it canonically, then encrypts it.
func (v *Vault) Encrypt(raw string) ([]byte, error) {
	parsed, err := Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("env encrypt: parse: %w", err)
	}
	normalized := Serialize(parsed)
	ciphertext, err := v.encryptor.Encrypt([]byte(normalized))
	if err != nil {
		return nil, fmt.Errorf("env encrypt: %w", err)
	}
	return ciphertext, nil
}

// Decrypt decrypts ciphertext and returns the parsed env map.
func (v *Vault) Decrypt(ciphertext []byte) (map[string]string, error) {
	plaintext, err := v.encryptor.Decrypt(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("env decrypt: %w", err)
	}
	parsed, err := Parse(string(plaintext))
	if err != nil {
		return nil, fmt.Errorf("env decrypt: parse: %w", err)
	}
	return parsed, nil
}

// DecryptWithDiff decrypts ciphertext and returns the env map along with a diff
// against the provided previous state.
func (v *Vault) DecryptWithDiff(ciphertext []byte, previous map[string]string) (map[string]string, []Change, error) {
	current, err := v.Decrypt(ciphertext)
	if err != nil {
		return nil, nil, err
	}
	changes := Diff(previous, current)
	return current, changes, nil
}
