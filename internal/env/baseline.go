package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved reference state of env vars for drift detection.
type Baseline struct {
	CreatedAt time.Time         `json:"created_at"`
	Vars      map[string]string `json:"vars"`
}

// BaselineStore persists and retrieves baselines by name.
type BaselineStore struct {
	path string
	data map[string]Baseline
}

// NewBaselineStore creates or loads a baseline store at the given path.
func NewBaselineStore(path string) (*BaselineStore, error) {
	bs := &BaselineStore{path: path, data: make(map[string]Baseline)}
	if err := bs.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("baseline store load: %w", err)
	}
	return bs, nil
}

// Set saves a baseline under the given name.
func (bs *BaselineStore) Set(name string, vars map[string]string) error {
	if name == "" {
		return errors.New("baseline name must not be empty")
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	bs.data[name] = Baseline{CreatedAt: time.Now().UTC(), Vars: copy}
	return bs.save()
}

// Get retrieves a baseline by name.
func (bs *BaselineStore) Get(name string) (Baseline, error) {
	b, ok := bs.data[name]
	if !ok {
		return Baseline{}, fmt.Errorf("baseline %q not found", name)
	}
	return b, nil
}

// Delete removes a baseline by name.
func (bs *BaselineStore) Delete(name string) error {
	if _, ok := bs.data[name]; !ok {
		return fmt.Errorf("baseline %q not found", name)
	}
	delete(bs.data, name)
	return bs.save()
}

// List returns all baseline names.
func (bs *BaselineStore) List() []string {
	names := make([]string, 0, len(bs.data))
	for k := range bs.data {
		names = append(names, k)
	}
	return names
}

// Drift compares current vars against a named baseline and returns a DiffResult.
func (bs *BaselineStore) Drift(name string, current map[string]string) (DiffResult, error) {
	b, err := bs.Get(name)
	if err != nil {
		return DiffResult{}, err
	}
	return Diff(b.Vars, current), nil
}

func (bs *BaselineStore) load() error {
	f, err := os.Open(bs.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&bs.data)
}

func (bs *BaselineStore) save() error {
	f, err := os.Create(bs.path)
	if err != nil {
		return fmt.Errorf("baseline save: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(bs.data)
}
