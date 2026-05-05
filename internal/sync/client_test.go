package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(t *testing.T) (*httptest.Server, *map[string]Payload) {
	t.Helper()
	store := make(map[string]Payload)
	mux := http.NewServeMux()
	mux.HandleFunc("/vaults/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/vaults/"):]
		switch r.Method {
		case http.MethodPut:
			var p Payload
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			store[name] = p
			w.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			p, ok := store[name]
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p) //nolint:errcheck
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return httptest.NewServer(mux), &store
}

func TestPushAndPull(t *testing.T) {
	srv, _ := newTestServer(t)
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	want := Payload{
		Name:       "production",
		Ciphertext: "AGE-ENCRYPTED-DATA",
		UpdatedAt:  time.Now().UTC().Truncate(time.Second),
	}

	if err := client.Push(want); err != nil {
		t.Fatalf("Push() error = %v", err)
	}

	got, err := client.Pull(want.Name)
	if err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
	if got.Name != want.Name || got.Ciphertext != want.Ciphertext {
		t.Errorf("Pull() = %+v, want %+v", got, want)
	}
}

func TestPullNotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	defer srv.Close()

	client := NewClient(srv.URL, "")
	_, err := client.Pull("nonexistent")
	if err == nil {
		t.Fatal("Pull() expected error for missing vault, got nil")
	}
}

func TestPushServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "key")
	err := client.Push(Payload{Name: "test", Ciphertext: "data"})
	if err == nil {
		t.Fatal("Push() expected error on server 500, got nil")
	}
}
