package env

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/user/envault/internal/crypto"
)

// Vault combines env parsing with age encryption so that .env files
// can be encrypted before being handed off to the store.
type Vault struct {
	enc *crypto.Encryptor
}

// NewVault creates a Vault backed by the given Encryptor.
func NewVault(enc *crypto.Encryptor) *Vault {
	return &Vault{enc: enc}
}

// Encrypt parses src as a .env file, serialises it, then returns the
// age-encrypted ciphertext.
func (v *Vault) Encrypt(src string) ([]byte, error) {
	entries, err := Parse(strings.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("parsing env: %w", err)
	}

	plaintext := Serialize(entries)

	var buf bytes.Buffer
	if err := v.enc.Encrypt(&buf, strings.NewReader(plaintext)); err != nil {
		return nil, fmt.Errorf("encrypting env: %w", err)
	}

	return buf.Bytes(), nil
}

// Decrypt decrypts ciphertext and returns the parsed entries.
func (v *Vault) Decrypt(ciphertext []byte) ([]Entry, error) {
	var buf bytes.Buffer
	if err := v.enc.Decrypt(&buf, bytes.NewReader(ciphertext)); err != nil {
		return nil, fmt.Errorf("decrypting env: %w", err)
	}

	entries, err := Parse(&buf)
	if err != nil {
		return nil, fmt.Errorf("parsing decrypted env: %w", err)
	}

	return entries, nil
}

// DecryptToString decrypts ciphertext and returns the .env file as a string.
func (v *Vault) DecryptToString(ciphertext []byte) (string, error) {
	entries, err := v.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return Serialize(entries), nil
}
