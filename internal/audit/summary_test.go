package audit

import (
	"testing"
	"time"
)

func makeSampleEntries() []Entry {
	base := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	return []Entry{
		{Timestamp: base, Event: "encrypt", User: "alice", Success: true},
		{Timestamp: base.Add(time.Hour), Event: "decrypt", User: "bob", Success: true},
		{Timestamp: base.Add(2 * time.Hour), Event: "encrypt", User: "alice", Success: false},
		{Timestamp: base.Add(3 * time.Hour), Event: "push", User: "bob", Success: true},
		{Timestamp: base.Add(4 * time.Hour), Event: "pull", User: "alice", Success: false},
	}
}

func TestSummarizeTotals(t *testing.T) {
	entries := makeSampleEntries()
	s := Summarize(entries)

	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Succeeded != 3 {
		t.Errorf("expected Succeeded=3, got %d", s.Succeeded)
	}
	if s.Failed != 2 {
		t.Errorf("expected Failed=2, got %d", s.Failed)
	}
}

func TestSummarizeByEvent(t *testing.T) {
	s := Summarize(makeSampleEntries())

	if s.ByEvent["encrypt"] != 2 {
		t.Errorf("expected encrypt=2, got %d", s.ByEvent["encrypt"])
	}
	if s.ByEvent["push"] != 1 {
		t.Errorf("expected push=1, got %d", s.ByEvent["push"])
	}
}

func TestSummarizeByUser(t *testing.T) {
	s := Summarize(makeSampleEntries())

	if s.ByUser["alice"] != 3 {
		t.Errorf("expected alice=3, got %d", s.ByUser["alice"])
	}
	if s.ByUser["bob"] != 2 {
		t.Errorf("expected bob=2, got %d", s.ByUser["bob"])
	}
}

func TestSummarizeTimestamps(t *testing.T) {
	entries := makeSampleEntries()
	s := Summarize(entries)

	expectedFirst := entries[0].Timestamp
	expectedLast := entries[len(entries)-1].Timestamp

	if s.FirstEntry == nil || !s.FirstEntry.Equal(expectedFirst) {
		t.Errorf("unexpected FirstEntry: %v", s.FirstEntry)
	}
	if s.LastEntry == nil || !s.LastEntry.Equal(expectedLast) {
		t.Errorf("unexpected LastEntry: %v", s.LastEntry)
	}
}

func TestSummarizeEmpty(t *testing.T) {
	s := Summarize([]Entry{})

	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if s.FirstEntry != nil || s.LastEntry != nil {
		t.Error("expected nil timestamps for empty input")
	}
}
