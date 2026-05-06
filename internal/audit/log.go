package audit

import (
	"encoding/json"
	"os"
	"time"
)

// EventType represents the type of audit event.
type EventType string

const (
	EventEncrypt EventType = "encrypt"
	EventDecrypt EventType = "decrypt"
	EventPush    EventType = "push"
	EventPull    EventType = "pull"
	EventRotate  EventType = "rotate"
	EventInit    EventType = "init"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	User      string    `json:"user,omitempty"`
	File      string    `json:"file,omitempty"`
	Details   string    `json:"details,omitempty"`
	Success   bool      `json:"success"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Log appends an Entry to the audit log file.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(e)
}

// ReadAll reads all audit entries from the log file.
func (l *Logger) ReadAll() ([]Entry, error) {
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
