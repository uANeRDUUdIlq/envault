package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcherDetectsChange(t *testing.T) {
	path := writeTempEnv(t, "KEY=original\n")
	w := NewWatcher(path, 20*time.Millisecond)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
		if ev.At.IsZero() {
			t.Error("expected non-zero timestamp")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcherNoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnv(t, "KEY=stable\n")
	w := NewWatcher(path, 20*time.Millisecond)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer w.Stop()

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(150 * time.Millisecond):
		// expected: no event
	}
}

func TestWatcherStartMissingFile(t *testing.T) {
	w := NewWatcher("/nonexistent/.env", 20*time.Millisecond)
	_, err := w.Start()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWatcherStopClosesChannel(t *testing.T) {
	path := writeTempEnv(t, "KEY=val\n")
	w := NewWatcher(path, 20*time.Millisecond)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	w.Stop()

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed after Stop")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for channel close")
	}
}
