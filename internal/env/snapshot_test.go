package env

import (
	"encoding/json"
	"testing"
)

func sampleVars() map[string]string {
	return map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
}

func TestSnapshotAdd(t *testing.T) {
	store := NewSnapshotStore()
	snap := store.Add("alice", "initial", sampleVars())

	if snap.Author != "alice" {
		t.Errorf("expected author alice, got %s", snap.Author)
	}
	if snap.Message != "initial" {
		t.Errorf("expected message initial, got %s", snap.Message)
	}
	if snap.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost")
	}
	if snap.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestSnapshotList(t *testing.T) {
	store := NewSnapshotStore()
	store.Add("alice", "first", sampleVars())
	store.Add("bob", "second", sampleVars())

	list := store.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(list))
	}
	if list[0].Author != "alice" || list[1].Author != "bob" {
		t.Error("unexpected snapshot order or authors")
	}
}

func TestSnapshotGet(t *testing.T) {
	store := NewSnapshotStore()
	snap := store.Add("alice", "test", sampleVars())

	found, err := store.Get(snap.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.ID != snap.ID {
		t.Errorf("ID mismatch: %s vs %s", found.ID, snap.ID)
	}
}

func TestSnapshotGetNotFound(t *testing.T) {
	store := NewSnapshotStore()
	_, err := store.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestSnapshotVarsAreCopied(t *testing.T) {
	store := NewSnapshotStore()
	vars := sampleVars()
	store.Add("alice", "copy-test", vars)
	vars["DB_HOST"] = "mutated"

	list := store.List()
	if list[0].Vars["DB_HOST"] == "mutated" {
		t.Error("snapshot vars should be a copy, not a reference")
	}
}

func TestSnapshotJSONRoundtrip(t *testing.T) {
	store := NewSnapshotStore()
	store.Add("alice", "json-test", sampleVars())

	data, err := store.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	store2 := NewSnapshotStore()
	if err := json.Unmarshal(data, store2); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(store2.List()) != 1 {
		t.Errorf("expected 1 snapshot after roundtrip, got %d", len(store2.List()))
	}
}
