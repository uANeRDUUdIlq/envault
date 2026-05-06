package audit

import "time"

// Summary holds aggregated statistics over a set of audit log entries.
type Summary struct {
	Total      int
	Succeeded  int
	Failed     int
	ByEvent    map[string]int
	ByUser     map[string]int
	FirstEntry *time.Time
	LastEntry  *time.Time
}

// Summarize computes aggregate statistics from the provided entries.
func Summarize(entries []Entry) Summary {
	s := Summary{
		ByEvent: make(map[string]int),
		ByUser:  make(map[string]int),
	}

	for i := range entries {
		e := &entries[i]
		s.Total++

		if e.Success {
			s.Succeeded++
		} else {
			s.Failed++
		}

		s.ByEvent[e.Event]++
		if e.User != "" {
			s.ByUser[e.User]++
		}

		t := e.Timestamp
		if s.FirstEntry == nil || t.Before(*s.FirstEntry) {
			s.FirstEntry = &t
		}
		if s.LastEntry == nil || t.After(*s.LastEntry) {
			s.LastEntry = &t
		}
	}

	return s
}

// SuccessRate returns the ratio of succeeded entries to total entries as a
// value between 0.0 and 1.0. Returns 0 if there are no entries.
func (s Summary) SuccessRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Succeeded) / float64(s.Total)
}
