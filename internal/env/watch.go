package env

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// ChangeEvent describes a detected change to a watched .env file.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watcher polls a file for changes and emits ChangeEvents.
type Watcher struct {
	path     string
	interval time.Duration
	lastHash string
	mu       sync.Mutex
	stop     chan struct{}
}

// NewWatcher creates a Watcher for the given file path and poll interval.
func NewWatcher(path string, interval time.Duration) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins polling and sends events on the returned channel.
// The caller must call Stop to release resources.
func (w *Watcher) Start() (<-chan ChangeEvent, error) {
	hash, err := hashFile(w.path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial hash: %w", err)
	}
	w.mu.Lock()
	w.lastHash = hash
	w.mu.Unlock()

	ch := make(chan ChangeEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				newHash, err := hashFile(w.path)
				if err != nil {
					continue
				}
				w.mu.Lock()
				prev := w.lastHash
				if newHash != prev {
					w.lastHash = newHash
					w.mu.Unlock()
					ch <- ChangeEvent{
						Path:    w.path,
						OldHash: prev,
						NewHash: newHash,
						At:      time.Now().UTC(),
					}
				} else {
					w.mu.Unlock()
				}
			}
		}
	}()
	return ch, nil
}

// Stop halts the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
