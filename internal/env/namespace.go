package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

var validNamespace = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{0,63}$`)

// Namespace represents a named environment scope (e.g. "production", "staging").
type Namespace struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
}

// NamespaceStore manages a collection of namespaces persisted to disk.
type NamespaceStore struct {
	path       string
	namespaces map[string]Namespace
}

// NewNamespaceStore loads or initialises a namespace store at the given path.
func NewNamespaceStore(path string) (*NamespaceStore, error) {
	ns := &NamespaceStore{path: path, namespaces: make(map[string]Namespace)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ns, nil
	}
	if err != nil {
		return nil, fmt.Errorf("namespace store: read: %w", err)
	}
	if err := json.Unmarshal(data, &ns.namespaces); err != nil {
		return nil, fmt.Errorf("namespace store: parse: %w", err)
	}
	return ns, nil
}

// Add registers a new namespace. Returns an error if the name is invalid or already exists.
func (s *NamespaceStore) Add(ns Namespace) error {
	if !validNamespace.MatchString(ns.Name) {
		return fmt.Errorf("invalid namespace name %q: must match %s", ns.Name, validNamespace)
	}
	if _, exists := s.namespaces[ns.Name]; exists {
		return fmt.Errorf("namespace %q already exists", ns.Name)
	}
	s.namespaces[ns.Name] = ns
	return s.save()
}

// Get returns the namespace with the given name.
func (s *NamespaceStore) Get(name string) (Namespace, bool) {
	ns, ok := s.namespaces[name]
	return ns, ok
}

// List returns all registered namespaces in insertion-stable order.
func (s *NamespaceStore) List() []Namespace {
	out := make([]Namespace, 0, len(s.namespaces))
	for _, ns := range s.namespaces {
		out = append(out, ns)
	}
	return out
}

// Delete removes a namespace by name.
func (s *NamespaceStore) Delete(name string) error {
	if _, ok := s.namespaces[name]; !ok {
		return fmt.Errorf("namespace %q not found", name)
	}
	delete(s.namespaces, name)
	return s.save()
}

func (s *NamespaceStore) save() error {
	data, err := json.MarshalIndent(s.namespaces, "", "  ")
	if err != nil {
		return fmt.Errorf("namespace store: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return fmt.Errorf("namespace store: write: %w", err)
	}
	return nil
}
