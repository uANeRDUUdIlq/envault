package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

var validAliasName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// Alias maps a short name to a canonical environment key.
type Alias struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// AliasStore persists key aliases to a JSON file.
type AliasStore struct {
	path    string
	aliases map[string]string // name -> canonical key
}

// NewAliasStore loads or creates an alias store at the given path.
func NewAliasStore(path string) (*AliasStore, error) {
	s := &AliasStore{path: path, aliases: make(map[string]string)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("alias store: read: %w", err)
	}
	if err := json.Unmarshal(data, &s.aliases); err != nil {
		return nil, fmt.Errorf("alias store: parse: %w", err)
	}
	return s, nil
}

// Set registers an alias name pointing to a canonical key.
func (s *AliasStore) Set(name, key string) error {
	if !validAliasName.MatchString(name) {
		return fmt.Errorf("alias store: invalid alias name %q", name)
	}
	if key == "" {
		return errors.New("alias store: key must not be empty")
	}
	s.aliases[name] = key
	return s.save()
}

// Get resolves an alias name to its canonical key.
func (s *AliasStore) Get(name string) (string, bool) {
	k, ok := s.aliases[name]
	return k, ok
}

// Delete removes an alias by name.
func (s *AliasStore) Delete(name string) error {
	if _, ok := s.aliases[name]; !ok {
		return fmt.Errorf("alias store: alias %q not found", name)
	}
	delete(s.aliases, name)
	return s.save()
}

// List returns all registered aliases.
func (s *AliasStore) List() []Alias {
	out := make([]Alias, 0, len(s.aliases))
	for n, k := range s.aliases {
		out = append(out, Alias{Name: n, Key: k})
	}
	return out
}

// Resolve returns the canonical key for name, falling back to name itself.
func (s *AliasStore) Resolve(name string) string {
	if k, ok := s.aliases[name]; ok {
		return k
	}
	return name
}

func (s *AliasStore) save() error {
	data, err := json.MarshalIndent(s.aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("alias store: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return fmt.Errorf("alias store: write: %w", err)
	}
	return nil
}
