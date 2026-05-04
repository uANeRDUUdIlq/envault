package crypto

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// Encryptor handles age-based encryption and decryption of secrets.
type Encryptor struct {
	recipients []age.Recipient
	identities []age.Identity
}

// NewEncryptor creates an Encryptor from PEM-encoded public keys (recipients)
// and optional private keys (identities).
func NewEncryptor(publicKeys []string, privateKeys []string) (*Encryptor, error) {
	var recipients []age.Recipient
	for _, pub := range publicKeys {
		r, err := age.ParseX25519Recipient(strings.TrimSpace(pub))
		if err != nil {
			return nil, fmt.Errorf("invalid public key %q: %w", pub, err)
		}
		recipients = append(recipients, r)
	}

	var identities []age.Identity
	for _, priv := range privateKeys {
		id, err := age.ParseX25519Identity(strings.TrimSpace(priv))
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		identities = append(identities, id)
	}

	return &Encryptor{recipients: recipients, identities: identities}, nil
}

// Encrypt encrypts plaintext using the configured recipients and returns
// an ASCII-armored ciphertext string.
func (e *Encryptor) Encrypt(plaintext []byte) (string, error) {
	if len(e.recipients) == 0 {
		return "", fmt.Errorf("no recipients configured")
	}

	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)

	w, err := age.Encrypt(armorWriter, e.recipients...)
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return "", fmt.Errorf("write plaintext: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close writer: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return "", fmt.Errorf("close armor: %w", err)
	}

	return buf.String(), nil
}

// Decrypt decrypts an ASCII-armored ciphertext string using the configured
// identities and returns the original plaintext.
func (e *Encryptor) Decrypt(ciphertext string) ([]byte, error) {
	if len(e.identities) == 0 {
		return nil, fmt.Errorf("no identities configured for decryption")
	}

	armorReader := armor.NewReader(strings.NewReader(ciphertext))
	r, err := age.Decrypt(armorReader, e.identities...)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read plaintext: %w", err)
	}

	return plaintext, nil
}
