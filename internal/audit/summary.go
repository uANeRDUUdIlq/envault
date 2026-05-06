package audit

import "fmt"

// Summary holds aggregated statistics over a set of audit entries.
type Summary struct {
	Total      int
	Succeeded  int
	Failed     int
	ByEvent    map[EventType]int
	ByUser     map[string]int
}

// Summarize computes a Summary from the provided entries.
func Summarize(entries []Entry) Summary {
	s := Summary{
		ByEvent: make(map[EventType]int),
		ByUser:  make(map[string]int),
	}
	for _, e := range entries {
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
	}
	return s
}

// String returns a human-readable summary string.
func (s Summary) String() string {
	return fmt.Sprintf(
		"total=%d succeeded=%d failed=%d events=%d users=%d",
		s.Total, s.Succeeded, s.Failed,
		len(s.ByEvent), len(s.ByUser),
	)
}
