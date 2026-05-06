package env

import (
	"fmt"

	"github.com/envault/envault/internal/crypto"
)

// Vault manages encryption and decryption of environment variable maps.
type Vault struct {
	enc *crypto.Encryptor
}

// NewVault creates a Vault backed by the given Encryptor.
func NewVault(enc *crypto.Encryptor) *Vault {
	return &Vault{enc: enc}
}

// Seal serializes and encrypts an env map, returning ciphertext bytes.
func (v *Vault) Seal(vars map[string]string) ([]byte, error) {
	raw := Serialize(vars)
	ciphertext, err := v.enc.Encrypt([]byte(raw))
	if err != nil {
		return nil, fmt.Errorf("vault seal: %w", err)
	}
	return ciphertext, nil
}

// Open decrypts ciphertext and parses it back into an env map.
func (v *Vault) Open(ciphertext []byte) (map[string]string, error) {
	plaintext, err := v.enc.Decrypt(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("vault open: %w", err)
	}
	vars, err := Parse(string(plaintext))
	if err != nil {
		return nil, fmt.Errorf("vault open parse: %w", err)
	}
	return vars, nil
}

// SealAndRotate rotates old vars into next vars, then seals the result.
func (v *Vault) SealAndRotate(old, next map[string]string, rotator *Rotator) ([]byte, RotationRecord, error) {
	merged, record := rotator.Rotate(old, next)
	ciphertext, err := v.Seal(merged)
	if err != nil {
		return nil, RotationRecord{}, fmt.Errorf("seal after rotate: %w", err)
	}
	return ciphertext, record, nil
}
