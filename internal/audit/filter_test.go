package audit_test

import (
	"testing"
	"time"

	"github.com/envault/envault/internal/audit"
)

func boolPtr(b bool) *bool { return &b }

func sampleEntries() []audit.Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []audit.Entry{
		{Timestamp: base, Event: audit.EventEncrypt, User: "alice", File: ".env", Success: true},
		{Timestamp: base.Add(time.Hour), Event: audit.EventPush, User: "bob", File: ".env.prod", Success: true},
		{Timestamp: base.Add(2 * time.Hour), Event: audit.EventDecrypt, User: "alice", File: ".env", Success: false},
		{Timestamp: base.Add(3 * time.Hour), Event: audit.EventPull, User: "carol", File: ".env.prod", Success: true},
	}
}

func TestFilterByEvent(t *testing.T) {
	f := audit.Filter{Event: audit.EventEncrypt}
	got := f.Apply(sampleEntries())
	if len(got) != 1 || got[0].User != "alice" {
		t.Errorf("unexpected results: %+v", got)
	}
}

func TestFilterByUser(t *testing.T) {
	f := audit.Filter{User: "alice"}
	got := f.Apply(sampleEntries())
	if len(got) != 2 {
		t.Errorf("expected 2 entries for alice, got %d", len(got))
	}
}

func TestFilterBySuccess(t *testing.T) {
	f := audit.Filter{Success: boolPtr(false)}
	got := f.Apply(sampleEntries())
	if len(got) != 1 || got[0].Event != audit.EventDecrypt {
		t.Errorf("unexpected results: %+v", got)
	}
}

func TestFilterBySince(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f := audit.Filter{Since: base.Add(2 * time.Hour)}
	got := f.Apply(sampleEntries())
	if len(got) != 2 {
		t.Errorf("expected 2 entries since +2h, got %d", len(got))
	}
}

func TestFilterNoMatch(t *testing.T) {
	f := audit.Filter{User: "unknown"}
	got := f.Apply(sampleEntries())
	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}

func TestQueryIntegration(t *testing.T) {
	path := tempLogPath(t)
	logger := audit.NewLogger(path)

	_ = logger.Log(audit.Entry{Event: audit.EventPush, User: "alice", Success: true})
	_ = logger.Log(audit.Entry{Event: audit.EventPull, User: "bob", Success: true})
	_ = logger.Log(audit.Entry{Event: audit.EventPush, User: "alice", Success: false})

	results, err := logger.Query(audit.Filter{User: "alice"})
	if err != nil {
		t.Fatalf("Query() error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
