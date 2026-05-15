package env

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ArchiveEntry holds a named archive of env vars with metadata.
type ArchiveEntry struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	ArchivedAt time.Time        `json:"archived_at"`
	Note      string            `json:"note,omitempty"`
}

// ArchiveStore persists named archives to disk.
type ArchiveStore struct {
	path    string
	entries map[string]ArchiveEntry
}

// NewArchiveStore loads or creates an archive store at the given path.
func NewArchiveStore(path string) (*ArchiveStore, error) {
	s := &ArchiveStore{path: path, entries: make(map[string]ArchiveEntry)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("archive: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, fmt.Errorf("archive: parse %s: %w", path, err)
	}
	return s, nil
}

// Save archives the given vars under name with an optional note.
func (s *ArchiveStore) Save(name string, vars map[string]string, note string) error {
	if name == "" {
		return fmt.Errorf("archive: name must not be empty")
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	s.entries[name] = ArchiveEntry{
		Name:       name,
		Vars:       copy,
		ArchivedAt: time.Now().UTC(),
		Note:       note,
	}
	return s.flush()
}

// Get retrieves an archive by name.
func (s *ArchiveStore) Get(name string) (ArchiveEntry, bool) {
	e, ok := s.entries[name]
	return e, ok
}

// List returns all archive names sorted alphabetically.
func (s *ArchiveStore) List() []string {
	names := make([]string, 0, len(s.entries))
	for n := range s.entries {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Delete removes an archive by name.
func (s *ArchiveStore) Delete(name string) error {
	if _, ok := s.entries[name]; !ok {
		return fmt.Errorf("archive: %q not found", name)
	}
	delete(s.entries, name)
	return s.flush()
}

func (s *ArchiveStore) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("archive: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0600); err != nil {
		return fmt.Errorf("archive: write %s: %w", s.path, err)
	}
	return nil
}
