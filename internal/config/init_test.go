package config_test

import (
	"testing"

	"github.com/your-org/envault/internal/config"
)

func TestInitCreatesConfig(t *testing.T) {
	dir := tempDir(t)
	opts := config.InitOptions{
		Dir:        dir,
		Project:    "demo",
		BackendURL: "https://sync.example.com",
		Recipients: []string{"age1pubkey1"},
	}
	cfg, err := config.Init(opts)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}
	if cfg.Project != "demo" {
		t.Errorf("Project: got %q, want %q", cfg.Project, "demo")
	}

	// Verify it was actually persisted.
	loaded, err := config.Load(dir)
	if err != nil {
		t.Fatalf("Load() after Init() error: %v", err)
	}
	if loaded.Project != "demo" {
		t.Errorf("loaded Project: got %q, want %q", loaded.Project, "demo")
	}
}

func TestInitFailsIfExists(t *testing.T) {
	dir := tempDir(t)
	opts := config.InitOptions{
		Dir:        dir,
		Project:    "demo",
		BackendURL: "https://sync.example.com",
		Recipients: []string{"age1pubkey1"},
	}
	if _, err := config.Init(opts); err != nil {
		t.Fatalf("first Init() error: %v", err)
	}
	_, err := config.Init(opts)
	if err != config.ErrAlreadyExists {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestInitOverwrite(t *testing.T) {
	dir := tempDir(t)
	opts := config.InitOptions{
		Dir:        dir,
		Project:    "original",
		BackendURL: "https://sync.example.com",
		Recipients: []string{"age1pubkey1"},
	}
	if _, err := config.Init(opts); err != nil {
		t.Fatalf("first Init() error: %v", err)
	}
	opts.Project = "updated"
	opts.Overwrite = true
	if _, err := config.Init(opts); err != nil {
		t.Fatalf("overwrite Init() error: %v", err)
	}
	loaded, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Project != "updated" {
		t.Errorf("Project after overwrite: got %q, want %q", loaded.Project, "updated")
	}
}

func TestInitValidationError(t *testing.T) {
	dir := tempDir(t)
	opts := config.InitOptions{
		Dir:        dir,
		Project:    "", // missing — should fail validation
		BackendURL: "https://sync.example.com",
		Recipients: []string{"age1pubkey1"},
	}
	_, err := config.Init(opts)
	if err == nil {
		t.Error("expected validation error for empty project, got nil")
	}
}
