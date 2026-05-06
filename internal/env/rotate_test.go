package env

import (
	"strings"
	"testing"
)

func TestRotateAddsNewKeys(t *testing.T) {
	r := NewRotator()
	old := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	next := map[string]string{"DB_HOST": "prod.db", "DB_PORT": "5432", "DB_NAME": "mydb"}

	merged, rec := r.Rotate(old, next)

	if merged["DB_HOST"] != "prod.db" {
		t.Errorf("expected DB_HOST=prod.db, got %s", merged["DB_HOST"])
	}
	if merged["DB_NAME"] != "mydb" {
		t.Errorf("expected DB_NAME=mydb, got %s", merged["DB_NAME"])
	}
	if len(rec.NewKeys) != 3 {
		t.Errorf("expected 3 new keys, got %d", len(rec.NewKeys))
	}
	if len(rec.RemovedAt) != 0 {
		t.Errorf("expected 0 removed keys, got %d", len(rec.RemovedAt))
	}
}

func TestRotateTracksRemovedKeys(t *testing.T) {
	r := NewRotator()
	old := map[string]string{"API_KEY": "secret", "OLD_FLAG": "true"}
	next := map[string]string{"API_KEY": "newsecret"}

	_, rec := r.Rotate(old, next)

	if _, removed := rec.RemovedAt["OLD_FLAG"]; !removed {
		t.Error("expected OLD_FLAG to be recorded as removed")
	}
	if len(rec.RemovedAt) != 1 {
		t.Errorf("expected 1 removed key, got %d", len(rec.RemovedAt))
	}
}

func TestRotateHistoryGrows(t *testing.T) {
	r := NewRotator()
	env1 := map[string]string{"A": "1"}
	env2 := map[string]string{"A": "2"}
	env3 := map[string]string{"A": "3"}

	r.Rotate(env1, env2)
	r.Rotate(env2, env3)

	if len(r.History()) != 2 {
		t.Errorf("expected 2 history records, got %d", len(r.History()))
	}
}

func TestSummaryString(t *testing.T) {
	r := NewRotator()
	old := map[string]string{"X": "1", "Y": "2"}
	next := map[string]string{"X": "10", "Z": "3"}

	_, rec := r.Rotate(old, next)
	summary := SummaryString(rec)

	if !strings.Contains(summary, "rotated:") {
		t.Errorf("expected summary to contain 'rotated:', got: %s", summary)
	}
	if !strings.Contains(summary, "1 removed") {
		t.Errorf("expected summary to mention 1 removed, got: %s", summary)
	}
}

func TestRotateEmptyOld(t *testing.T) {
	r := NewRotator()
	old := map[string]string{}
	next := map[string]string{"NEW_KEY": "value"}

	merged, rec := r.Rotate(old, next)

	if merged["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %s", merged["NEW_KEY"])
	}
	if len(rec.OldKeys) != 0 {
		t.Errorf("expected 0 old keys, got %d", len(rec.OldKeys))
	}
}
