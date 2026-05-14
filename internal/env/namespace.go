package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

var validNamespaceName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

type namespaceEntry struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

type namespaceData struct {
	Namespaces []namespaceEntry `json:"namespaces"`
}

// NamespaceStore manages named env namespaces persisted to a JSON file.
type NamespaceStore struct {
	path string
	data namespaceData
}

// NewNamespaceStore loads or initialises a NamespaceStore at path.
func NewNamespaceStore(path string) (*NamespaceStore, error) {
	s := &NamespaceStore{path: path}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

func (s *NamespaceStore) load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *NamespaceStore) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0600)
}

// Add registers a new namespace.
func (s *NamespaceStore) Add(name string) error {
	if !validNamespaceName.MatchString(name) {
		return fmt.Errorf("namespace: invalid name %q", name)
	}
	for _, ns := range s.data.Namespaces {
		if ns.Name == name {
			return fmt.Errorf("namespace: %q already exists", name)
		}
	}
	s.data.Namespaces = append(s.data.Namespaces, namespaceEntry{Name: name, Vars: map[string]string{}})
	return s.save()
}

// List returns all namespace names.
func (s *NamespaceStore) List() []string {
	names := make([]string, len(s.data.Namespaces))
	for i, ns := range s.data.Namespaces {
		names[i] = ns.Name
	}
	return names
}

// GetVars returns a copy of the vars for a namespace.
func (s *NamespaceStore) GetVars(name string) (map[string]string, error) {
	for _, ns := range s.data.Namespaces {
		if ns.Name == name {
			copy := make(map[string]string, len(ns.Vars))
			for k, v := range ns.Vars {
				copy[k] = v
			}
			return copy, nil
		}
	}
	return nil, fmt.Errorf("namespace: %q not found", name)
}

// SetVars replaces the vars for a namespace.
func (s *NamespaceStore) SetVars(name string, vars map[string]string) error {
	for i, ns := range s.data.Namespaces {
		if ns.Name == name {
			s.data.Namespaces[i].Vars = vars
			return s.save()
		}
	}
	return fmt.Errorf("namespace: %q not found", name)
}

// Delete removes a namespace.
func (s *NamespaceStore) Delete(name string) error {
	for i, ns := range s.data.Namespaces {
		if ns.Name == name {
			s.data.Namespaces = append(s.data.Namespaces[:i], s.data.Namespaces[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("namespace: %q not found", name)
}
