package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempTTLPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ttl.json")
}

func TestTTLSetAndGet(t *testing.T) {
	s, err := NewTTLStore(tempTTLPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("API_KEY", time.Hour); err != nil {
		t.Fatal(err)
	}
	e, err := s.Get("API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if e.Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", e.Key)
	}
	if time.Until(e.ExpiresAt) <= 0 {
		t.Error("expected future expiry")
	}
}

func TestTTLPersists(t *testing.T) {
	path := tempTTLPath(t)
	s, _ := NewTTLStore(path)
	_ = s.Set("DB_PASS", time.Hour)

	s2, err := NewTTLStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s2.Get("DB_PASS"); err != nil {
		t.Error("expected persisted entry")
	}
}

func TestTTLIsExpired(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	_ = s.Set("OLD_KEY", time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !s.IsExpired("OLD_KEY") {
		t.Error("expected key to be expired")
	}
}

func TestTTLNotExpired(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	_ = s.Set("FRESH_KEY", time.Hour)
	if s.IsExpired("FRESH_KEY") {
		t.Error("expected key to not be expired")
	}
}

func TestTTLRemove(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	_ = s.Set("TMP", time.Hour)
	if err := s.Remove("TMP"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("TMP"); err == nil {
		t.Error("expected entry to be removed")
	}
}

func TestTTLExpiredList(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	_ = s.Set("STALE", time.Millisecond)
	_ = s.Set("FRESH", time.Hour)
	time.Sleep(5 * time.Millisecond)
	expired := s.Expired()
	if len(expired) != 1 || expired[0] != "STALE" {
		t.Errorf("unexpected expired list: %v", expired)
	}
}

func TestTTLInvalidInputs(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	if err := s.Set("", time.Hour); err == nil {
		t.Error("expected error for empty key")
	}
	if err := s.Set("K", -time.Second); err == nil {
		t.Error("expected error for non-positive ttl")
	}
	if err := s.Remove("MISSING"); err == nil {
		t.Error("expected error removing missing key")
	}
}

func TestTTLMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no.json")
	s, err := NewTTLStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Error("expected non-nil store")
	}
	_ = os.Remove(path)
}
