package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Checkpoint represents a named point-in-time capture of env vars.
type Checkpoint struct {
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Vars      map[string]string `json:"vars"`
	Message   string            `json:"message,omitempty"`
}

// CheckpointStore persists named checkpoints to a JSON file.
type CheckpointStore struct {
	path string
	data map[string]Checkpoint
}

// NewCheckpointStore loads or initialises a checkpoint store at path.
func NewCheckpointStore(path string) (*CheckpointStore, error) {
	cs := &CheckpointStore{path: path, data: make(map[string]Checkpoint)}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cs, nil
		}
		return nil, fmt.Errorf("checkpoint: read %s: %w", path, err)
	}
	if err := json.Unmarshal(b, &cs.data); err != nil {
		return nil, fmt.Errorf("checkpoint: parse %s: %w", path, err)
	}
	return cs, nil
}

// Save creates or overwrites a named checkpoint.
func (cs *CheckpointStore) Save(name string, vars map[string]string, message string) error {
	if name == "" {
		return fmt.Errorf("checkpoint: name must not be empty")
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	cs.data[name] = Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Vars:      copy,
		Message:   message,
	}
	return cs.flush()
}

// Get retrieves a checkpoint by name.
func (cs *CheckpointStore) Get(name string) (Checkpoint, bool) {
	cp, ok := cs.data[name]
	return cp, ok
}

// List returns all checkpoint names in the store.
func (cs *CheckpointStore) List() []string {
	names := make([]string, 0, len(cs.data))
	for n := range cs.data {
		names = append(names, n)
	}
	return names
}

// Delete removes a checkpoint by name.
func (cs *CheckpointStore) Delete(name string) error {
	if _, ok := cs.data[name]; !ok {
		return fmt.Errorf("checkpoint: %q not found", name)
	}
	delete(cs.data, name)
	return cs.flush()
}

func (cs *CheckpointStore) flush() error {
	b, err := json.MarshalIndent(cs.data, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	if err := os.WriteFile(cs.path, b, 0600); err != nil {
		return fmt.Errorf("checkpoint: write %s: %w", cs.path, err)
	}
	return nil
}
