package env

import (
	"testing"
	"time"
)

// TestTTLReaperRoundtrip exercises the full lifecycle:
// set TTL → wait for expiry → reap → confirm vault and store are clean.
func TestTTLReaperRoundtrip(t *testing.T) {
	path := tempTTLPath(t)
	store, err := NewTTLStore(path)
	if err != nil {
		t.Fatal(err)
	}

	vault := map[string]string{
		"SECRET_A": "aaa",
		"SECRET_B": "bbb",
	}

	_ = store.Set("SECRET_A", time.Millisecond)
	_ = store.Set("SECRET_B", time.Hour)
	time.Sleep(10 * time.Millisecond)

	reaper := NewReaper(store, vault)
	res := reaper.Reap()

	if len(res.Removed) != 1 || res.Removed[0] != "SECRET_A" {
		t.Fatalf("expected SECRET_A removed, got %v", res.Removed)
	}
	if _, ok := vault["SECRET_A"]; ok {
		t.Error("SECRET_A should be absent from vault")
	}
	if vault["SECRET_B"] != "bbb" {
		t.Error("SECRET_B should remain intact")
	}

	// Reload store and confirm SECRET_A is gone.
	store2, err := NewTTLStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store2.Get("SECRET_A"); err == nil {
		t.Error("expected SECRET_A absent after reload")
	}
	if _, err := store2.Get("SECRET_B"); err != nil {
		t.Error("expected SECRET_B present after reload")
	}
}

// TestTTLMultipleExpirations ensures multiple keys expire and are reaped correctly.
func TestTTLMultipleExpirations(t *testing.T) {
	store, _ := NewTTLStore(tempTTLPath(t))
	vault := map[string]string{"A": "1", "B": "2", "C": "3"}

	_ = store.Set("A", time.Millisecond)
	_ = store.Set("B", time.Millisecond)
	_ = store.Set("C", time.Hour)
	time.Sleep(10 * time.Millisecond)

	reaper := NewReaper(store, vault)
	res := reaper.Reap()

	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removals, got %d: %v", len(res.Removed), res.Removed)
	}
	if len(vault) != 1 {
		t.Errorf("expected vault length 1, got %d", len(vault))
	}
}
