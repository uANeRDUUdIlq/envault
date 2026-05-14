package env

import (
	"fmt"
	"os/exec"
	"strings"
)

// HookEvent represents the lifecycle event that triggers a hook.
type HookEvent string

const (
	HookPreEncrypt  HookEvent = "pre-encrypt"
	HookPostEncrypt HookEvent = "post-encrypt"
	HookPreDecrypt  HookEvent = "pre-decrypt"
	HookPostDecrypt HookEvent = "post-decrypt"
)

// Hook defines a shell command to run on a specific event.
type Hook struct {
	Event   HookEvent
	Command string
}

// HookRunner executes registered hooks for lifecycle events.
type HookRunner struct {
	hooks []Hook
}

// NewHookRunner creates a HookRunner with the given hooks.
func NewHookRunner(hooks []Hook) *HookRunner {
	return &HookRunner{hooks: hooks}
}

// Run executes all hooks registered for the given event.
// It stops and returns an error on the first failure.
func (r *HookRunner) Run(event HookEvent) error {
	for _, h := range r.hooks {
		if h.Event != event {
			continue
		}
		if err := r.exec(h.Command); err != nil {
			return fmt.Errorf("hook %q failed: %w", h.Command, err)
		}
	}
	return nil
}

// HasHooks reports whether any hooks are registered for the given event.
func (r *HookRunner) HasHooks(event HookEvent) bool {
	for _, h := range r.hooks {
		if h.Event == event {
			return true
		}
	}
	return false
}

func (r *HookRunner) exec(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty hook command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
