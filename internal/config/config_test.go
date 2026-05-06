package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envault/internal/config"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envault-config-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoad(t *testing.T) {
	dir := tempDir(t)
	cfg := &config.Config{
		BackendURL:   "https://sync.example.com",
		Project:      "myapp",
		Recipients:   []string{"age1abc123", "age1def456"},
		IdentityFile: "~/.config/envault/identity",
	}
	if err := config.Save(dir, cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	loaded, err := config.Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.Project != cfg.Project {
		t.Errorf("Project: got %q, want %q", loaded.Project, cfg.Project)
	}
	if loaded.BackendURL != cfg.BackendURL {
		t.Errorf("BackendURL: got %q, want %q", loaded.BackendURL, cfg.BackendURL)
	}
	if len(loaded.Recipients) != 2 {
		t.Errorf("Recipients: got %d, want 2", len(loaded.Recipients))
	}
}

func TestLoadNotFound(t *testing.T) {
	dir := tempDir(t)
	_, err := config.Load(dir)
	if err != config.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, ".envault.json")
	if err := os.WriteFile(path, []byte("not-json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := config.Load(dir)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestValidate(t *testing.T) {
	valid := &config.Config{
		BackendURL: "https://sync.example.com",
		Project:    "myapp",
		Recipients: []string{"age1abc"},
	}
	if err := config.Validate(valid); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	cases := []struct {
		name string
		cfg  *config.Config
	}{
		{"missing project", &config.Config{BackendURL: "https://x.com", Recipients: []string{"age1"}}},
		{"missing backend", &config.Config{Project: "p", Recipients: []string{"age1"}}},
		{"no recipients", &config.Config{Project: "p", BackendURL: "https://x.com"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := config.Validate(tc.cfg); err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}
