package env

import (
	"time"
)

// ReapResult summarises a reaping run.
type ReapResult struct {
	Removed []string
	Errors  []error
}

// Reaper periodically removes expired TTL keys from a variable map.
type Reaper struct {
	store  *TTLStore
	vault  map[string]string
	notify chan ReapResult
}

// NewReaper creates a Reaper that will operate on the given vault map.
func NewReaper(store *TTLStore, vault map[string]string) *Reaper {
	return &Reaper{
		store:  store,
		vault:  vault,
		notify: make(chan ReapResult, 8),
	}
}

// Notify returns a channel that receives a ReapResult after each reap cycle.
func (r *Reaper) Notify() <-chan ReapResult {
	return r.notify
}

// Reap performs a single reap pass and returns the result.
func (r *Reaper) Reap() ReapResult {
	expired := r.store.Expired()
	result := ReapResult{}
	for _, key := range expired {
		delete(r.vault, key)
		if err := r.store.Remove(key); err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		result.Removed = append(result.Removed, key)
	}
	return result
}

// Start launches a background goroutine that reaps every interval.
// It stops when the done channel is closed.
func (r *Reaper) Start(interval time.Duration, done <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				res := r.Reap()
				if len(res.Removed) > 0 || len(res.Errors) > 0 {
					select {
					case r.notify <- res:
					default:
					}
				}
			case <-done:
				return
			}
		}
	}()
}
