package env

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a label attached to a set of env var keys.
type Tag struct {
	Name string
	Keys []string
}

// TagStore manages named tags that group env variable keys.
type TagStore struct {
	tags map[string]*Tag
}

// NewTagStore creates an empty TagStore.
func NewTagStore() *TagStore {
	return &TagStore{tags: make(map[string]*Tag)}
}

// Add creates or updates a tag with the given keys.
func (ts *TagStore) Add(name string, keys []string) error {
	if name == "" {
		return fmt.Errorf("tag name must not be empty")
	}
	if len(keys) == 0 {
		return fmt.Errorf("tag %q must contain at least one key", name)
	}
	deduped := dedupKeys(keys)
	ts.tags[name] = &Tag{Name: name, Keys: deduped}
	return nil
}

// Get returns the Tag for the given name, or an error if not found.
func (ts *TagStore) Get(name string) (*Tag, error) {
	t, ok := ts.tags[name]
	if !ok {
		return nil, fmt.Errorf("tag %q not found", name)
	}
	return t, nil
}

// Delete removes a tag by name.
func (ts *TagStore) Delete(name string) error {
	if _, ok := ts.tags[name]; !ok {
		return fmt.Errorf("tag %q not found", name)
	}
	delete(ts.tags, name)
	return nil
}

// List returns all tag names in sorted order.
func (ts *TagStore) List() []string {
	names := make([]string, 0, len(ts.tags))
	for n := range ts.tags {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Filter returns only the vars whose keys belong to the given tag.
func (ts *TagStore) Filter(name string, vars map[string]string) (map[string]string, error) {
	tag, err := ts.Get(name)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, k := range tag.Keys {
		if v, ok := vars[k]; ok {
			result[k] = v
		}
	}
	return result, nil
}

func dedupKeys(keys []string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if _, exists := seen[k]; !exists {
			seen[k] = struct{}{}
			out = append(out, k)
		}
	}
	return out
}
