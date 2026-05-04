package crypto_test

import (
	"strings"
	"testing"

	"filippo.io/age"

	"github.com/envault/envault/internal/crypto"
)

func generateKeyPair(t *testing.T) (pub, priv string) {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}
	return id.Recipient().String(), id.String()
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	pub, priv := generateKeyPair(t)

	enc, err := crypto.NewEncryptor([]string{pub}, []string{priv})
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}

	plaintext := []byte("SECRET_KEY=supersecret\nDB_PASS=hunter2")

	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if !strings.Contains(ciphertext, "-----BEGIN AGE ENCRYPTED FILE-----") {
		t.Error("expected ASCII-armored output")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("roundtrip mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptNoRecipients(t *testing.T) {
	enc, err := crypto.NewEncryptor(nil, nil)
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}
	_, err = enc.Encrypt([]byte("data"))
	if err == nil {
		t.Error("expected error when no recipients configured")
	}
}

func TestDecryptNoIdentities(t *testing.T) {
	pub, _ := generateKeyPair(t)

	enc, err := crypto.NewEncryptor([]string{pub}, nil)
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}

	ciphertext, err := enc.Encrypt([]byte("secret"))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = enc.Decrypt(ciphertext)
	if err == nil {
		t.Error("expected error when no identities configured")
	}
}

func TestInvalidPublicKey(t *testing.T) {
	_, err := crypto.NewEncryptor([]string{"not-a-valid-key"}, nil)
	if err == nil {
		t.Error("expected error for invalid public key")
	}
}
