package env

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestHookRoundtripAndRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell integration test on windows")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "hooks.json")
	marker := filepath.Join(dir, "ran")

	cfg := &HookConfig{
		Hooks: []HookEntry{
			{Event: string(HookPostEncrypt), Command: "touch " + marker},
		},
	}
	if err := SaveHookConfig(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadHookConfig(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	runner := NewHookRunner(loaded.ToHooks())
	if err := runner.Run(HookPostEncrypt); err != nil {
		t.Fatalf("run: %v", err)
	}

	if _, err := os.Stat(marker); os.IsNotExist(err) {
		t.Fatal("hook did not create marker file")
	}
}

func TestHookIntegrationMultipleEvents(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell integration test on windows")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "hooks.json")

	cfg := &HookConfig{
		Hooks: []HookEntry{
			{Event: string(HookPreEncrypt), Command: "echo pre-encrypt"},
			{Event: string(HookPostEncrypt), Command: "echo post-encrypt"},
		},
	}
	if err := SaveHookConfig(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, _ := LoadHookConfig(path)
	runner := NewHookRunner(loaded.ToHooks())

	for _, ev := range []HookEvent{HookPreEncrypt, HookPostEncrypt} {
		if err := runner.Run(ev); err != nil {
			t.Errorf("event %s failed: %v", ev, err)
		}
	}
}
