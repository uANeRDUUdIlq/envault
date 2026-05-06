package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envault/envault/internal/audit"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.log")
}

func TestLogAndReadAll(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	entry := audit.Entry{
		Event:   audit.EventEncrypt,
		User:    "alice",
		File:    ".env",
		Success: true,
	}

	if err := logger.Log(entry); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != audit.EventEncrypt {
		t.Errorf("expected event %q, got %q", audit.EventEncrypt, entries[0].Event)
	}
	if entries[0].User != "alice" {
		t.Errorf("expected user alice, got %q", entries[0].User)
	}
}

func TestLogMultipleEntries(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	events := []audit.EventType{audit.EventPush, audit.EventPull, audit.EventRotate}
	for _, ev := range events {
		if err := logger.Log(audit.Entry{Event: ev, Success: true}); err != nil {
			t.Fatalf("Log(%s) error: %v", ev, err)
		}
	}

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestReadAllEmptyFile(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	entries, err := logger.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() on missing file error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestLogSetsTimestamp(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	before := time.Now().UTC()
	_ = logger.Log(audit.Entry{Event: audit.EventInit, Success: true})
	after := time.Now().UTC()

	entries, _ := logger.ReadAll()
	if len(entries) == 0 {
		t.Fatal("no entries")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of range [%v, %v]", ts, before, after)
	}
}

func TestLogPermissions(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)
	_ = logger.Log(audit.Entry{Event: audit.EventDecrypt, Success: false})

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}
