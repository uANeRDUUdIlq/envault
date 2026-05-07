package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempLockPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "vault.lock")
}

func TestAcquireAndRelease(t *testing.T) {
	ls := NewLockStore(tempLockPath(t))

	if err := ls.Acquire("alice", 60); err != nil {
		t.Fatalf("expected acquire to succeed: %v", err)
	}
	current := ls.Current()
	if current == nil {
		t.Fatal("expected current lock to be set")
	}
	if current.User != "alice" {
		t.Errorf("expected user alice, got %s", current.User)
	}
	if err := ls.Release("alice"); err != nil {
		t.Fatalf("expected release to succeed: %v", err)
	}
	if ls.Current() != nil {
		t.Error("expected no lock after release")
	}
}

func TestAcquireBlocksSecondUser(t *testing.T) {
	ls := NewLockStore(tempLockPath(t))

	if err := ls.Acquire("alice", 60); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ls.Acquire("bob", 60); err == nil {
		t.Error("expected error when acquiring lock held by alice")
	}
}

func TestReleaseWrongUser(t *testing.T) {
	ls := NewLockStore(tempLockPath(t))

	_ = ls.Acquire("alice", 60)
	if err := ls.Release("bob"); err == nil {
		t.Error("expected error releasing lock held by another user")
	}
}

func TestExpiredLockCanBeAcquired(t *testing.T) {
	path := tempLockPath(t)
	ls := NewLockStore(path)

	// Write an already-expired lock manually.
	entry := LockEntry{
		User:       "alice",
		AcquiredAt: time.Now().UTC().Add(-120 * time.Second),
		TTL:        10,
	}
	data := []byte(`{"user":"alice","acquired_at":"` +
		entry.AcquiredAt.Format(time.RFC3339) +
		`","ttl_seconds":10}`)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("failed to write expired lock: %v", err)
	}

	if err := ls.Acquire("bob", 60); err != nil {
		t.Errorf("expected expired lock to allow new acquisition: %v", err)
	}
}

func TestCurrentReturnsNilWhenNoLock(t *testing.T) {
	ls := NewLockStore(tempLockPath(t))
	if ls.Current() != nil {
		t.Error("expected nil when no lock file exists")
	}
}

func TestReleaseNoLock(t *testing.T) {
	ls := NewLockStore(tempLockPath(t))
	if err := ls.Release("alice"); err == nil {
		t.Error("expected error releasing non-existent lock")
	}
}
