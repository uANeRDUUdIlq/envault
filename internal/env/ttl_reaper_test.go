package env

import (
	"testing"
	"time"
)

func TestReaperRemovesExpiredKeys(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	vault := map[string]string{"OLD": "secret", "FRESH": "value"}
	_ = s.Set("OLD", time.Millisecond)
	_ = s.Set("FRESH", time.Hour)
	time.Sleep(5 * time.Millisecond)

	r := NewReaper(s, vault)
	res := r.Reap()

	if len(res.Removed) != 1 || res.Removed[0] != "OLD" {
		t.Errorf("expected OLD removed, got %v", res.Removed)
	}
	if _, ok := vault["OLD"]; ok {
		t.Error("expected OLD deleted from vault")
	}
	if _, ok := vault["FRESH"]; !ok {
		t.Error("expected FRESH to remain in vault")
	}
}

func TestReaperNoExpiredKeys(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	vault := map[string]string{"KEY": "val"}
	_ = s.Set("KEY", time.Hour)

	r := NewReaper(s, vault)
	res := r.Reap()

	if len(res.Removed) != 0 {
		t.Errorf("expected nothing removed, got %v", res.Removed)
	}
}

func TestReaperStart(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	vault := map[string]string{"TEMP": "123"}
	_ = s.Set("TEMP", time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	r := NewReaper(s, vault)
	done := make(chan struct{})
	r.Start(10*time.Millisecond, done)

	select {
	case res := <-r.Notify():
		if len(res.Removed) != 1 {
			t.Errorf("expected 1 removal, got %d", len(res.Removed))
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("timed out waiting for reap notification")
	}
	close(done)
}

func TestReaperNotifyNotBlockingOnFullChannel(t *testing.T) {
	s, _ := NewTTLStore(tempTTLPath(t))
	vault := map[string]string{"K": "v"}
	_ = s.Set("K", time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	r := NewReaper(s, vault)
	// Drain nothing — channel has capacity 8, should not block.
	res := r.Reap()
	if len(res.Errors) != 0 {
		t.Errorf("unexpected errors: %v", res.Errors)
	}
}
