package env

import (
	"encoding/json"
	"fmt"
	"os"
)

// HookConfig holds the serialisable hook definitions stored alongside a vault.
type HookConfig struct {
	Hooks []HookEntry `json:"hooks"`
}

// HookEntry is the JSON-serialisable representation of a single hook.
type HookEntry struct {
	Event   string `json:"event"`
	Command string `json:"command"`
}

// LoadHookConfig reads hook configuration from a JSON file.
// Returns an empty config if the file does not exist.
func LoadHookConfig(path string) (*HookConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &HookConfig{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read hook config: %w", err)
	}
	var cfg HookConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse hook config: %w", err)
	}
	return &cfg, nil
}

// SaveHookConfig writes hook configuration to a JSON file.
func SaveHookConfig(path string, cfg *HookConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal hook config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write hook config: %w", err)
	}
	return nil
}

// ToHooks converts a HookConfig into a slice of Hook values.
func (c *HookConfig) ToHooks() []Hook {
	out := make([]Hook, 0, len(c.Hooks))
	for _, e := range c.Hooks {
		out = append(out, Hook{Event: HookEvent(e.Event), Command: e.Command})
	}
	return out
}
