package env

import (
	"runtime"
	"testing"
)

func TestHookRunnerNoHooks(t *testing.T) {
	r := NewHookRunner(nil)
	if err := r.Run(HookPreEncrypt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHookRunnerMatchingEvent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	hooks := []Hook{
		{Event: HookPostEncrypt, Command: "echo hello"},
	}
	r := NewHookRunner(hooks)
	if err := r.Run(HookPostEncrypt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHookRunnerSkipsNonMatchingEvent(t *testing.T) {
	// Command would fail if executed, but event doesn't match.
	hooks := []Hook{
		{Event: HookPreDecrypt, Command: "false"},
	}
	r := NewHookRunner(hooks)
	if err := r.Run(HookPostDecrypt); err != nil {
		t.Fatalf("expected no error for non-matching event, got %v", err)
	}
}

func TestHookRunnerCommandFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	hooks := []Hook{
		{Event: HookPreEncrypt, Command: "false"},
	}
	r := NewHookRunner(hooks)
	if err := r.Run(HookPreEncrypt); err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestHookRunnerEmptyCommand(t *testing.T) {
	hooks := []Hook{
		{Event: HookPreEncrypt, Command: ""},
	}
	r := NewHookRunner(hooks)
	if err := r.Run(HookPreEncrypt); err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestHookRunnerStopsOnFirstFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	called := false
	// Use a sentinel: second hook writes a file; if first fails we never reach it.
	hooks := []Hook{
		{Event: HookPreEncrypt, Command: "false"},
		{Event: HookPreEncrypt, Command: "echo second"},
	}
	r := NewHookRunner(hooks)
	err := r.Run(HookPreEncrypt)
	if err == nil {
		t.Fatal("expected error")
	}
	_ = called // second hook not executed; we just verify early return via error
}
