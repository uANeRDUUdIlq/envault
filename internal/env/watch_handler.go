package env

import (
	"fmt"
	"io"
	"time"
)

// WatchHandler reacts to ChangeEvents by re-parsing the updated file
// and computing a diff against the previously known variables.
type WatchHandler struct {
	vars map[string]string
	out  io.Writer
}

// NewWatchHandler creates a WatchHandler with an initial set of variables
// and writes human-readable diff summaries to out.
func NewWatchHandler(initial map[string]string, out io.Writer) *WatchHandler {
	snap := make(map[string]string, len(initial))
	for k, v := range initial {
		snap[k] = v
	}
	return &WatchHandler{vars: snap, out: out}
}

// Handle processes a single ChangeEvent: it parses the new file content,
// computes the diff, prints a summary and updates the internal state.
func (h *WatchHandler) Handle(ev ChangeEvent) error {
	newVars, err := ParseFile(ev.Path)
	if err != nil {
		return fmt.Errorf("watch handler: parse %q: %w", ev.Path, err)
	}

	d := Diff(h.vars, newVars)
	if !HasChanges(d) {
		return nil
	}

	fmt.Fprintf(h.out, "[%s] %s changed\n", ev.At.Format(time.RFC3339), ev.Path)
	for _, e := range d {
		switch e.Op {
		case OpAdd:
			fmt.Fprintf(h.out, "  + %s=%s\n", e.Key, e.New)
		case OpRemove:
			fmt.Fprintf(h.out, "  - %s\n", e.Key)
		case OpUpdate:
			fmt.Fprintf(h.out, "  ~ %s: %s -> %s\n", e.Key, e.Old, e.New)
		}
	}

	h.vars = newVars
	return nil
}

// Current returns a copy of the most recently known variable set.
func (h *WatchHandler) Current() map[string]string {
	snap := make(map[string]string, len(h.vars))
	for k, v := range h.vars {
		snap[k] = v
	}
	return snap
}
