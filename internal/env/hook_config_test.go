package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadHookConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hooks.json")

	cfg := &HookConfig{
		Hooks: []HookEntry{
			{Event: "pre-encrypt", Command: "echo pre"},
			{Event: "post-decrypt", Command: "echo post"},
		},
	}
	if err := SaveHookConfig(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadHookConfig(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Hooks) != 2 {
		t.Fatalf("expected 2 hooks, got %d", len(loaded.Hooks))
	}
	if loaded.Hooks[0].Command != "echo pre" {
		t.Errorf("unexpected command: %s", loaded.Hooks[0].Command)
	}
}

func TestLoadHookConfigMissingFile(t *testing.T) {
	cfg, err := LoadHookConfig("/nonexistent/hooks.json")
	if err != nil {
		t.Fatalf("expected empty config, got error: %v", err)
	}
	if len(cfg.Hooks) != 0 {
		t.Errorf("expected empty hooks")
	}
}

func TestLoadHookConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hooks.json")
	os.WriteFile(path, []byte("not-json"), 0600)
	if _, err := LoadHookConfig(path); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestToHooksConversion(t *testing.T) {
	cfg := &HookConfig{
		Hooks: []HookEntry{
			{Event: "pre-decrypt", Command: "make prepare"},
		},
	}
	hooks := cfg.ToHooks()
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(hooks))
	}
	if hooks[0].Event != HookPreDecrypt {
		t.Errorf("unexpected event: %s", hooks[0].Event)
	}
	if hooks[0].Command != "make prepare" {
		t.Errorf("unexpected command: %s", hooks[0].Command)
	}
}
