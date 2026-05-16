package env

import (
	"os"
	"testing"
)

func TestInjectBasic(t *testing.T) {
	os.Unsetenv("APP_FOO")
	os.Unsetenv("APP_BAR")
	inj := NewInjector(InjectOptions{Overwrite: true, Prefix: "APP_"})
	vars := map[string]string{"FOO": "hello", "BAR": "world"}
	res, err := inj.Inject(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(res.Injected))
	}
	if os.Getenv("APP_FOO") != "hello" {
		t.Errorf("APP_FOO not set")
	}
}

func TestInjectSkipsExisting(t *testing.T) {
	os.Setenv("MYKEY", "original")
	defer os.Unsetenv("MYKEY")
	inj := NewInjector(InjectOptions{Overwrite: false})
	res, err := inj.Inject(map[string]string{"MYKEY": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if os.Getenv("MYKEY") != "original" {
		t.Errorf("MYKEY should not have been overwritten")
	}
}

func TestInjectFilterKeys(t *testing.T) {
	os.Unsetenv("ONLY")
	os.Unsetenv("SKIP")
	defer os.Unsetenv("ONLY")
	inj := NewInjector(InjectOptions{Overwrite: true, Keys: []string{"ONLY"}})
	res, err := inj.Inject(map[string]string{"ONLY": "yes", "SKIP": "no"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Injected) != 1 || res.Injected[0] != "ONLY" {
		t.Errorf("expected only ONLY injected, got %v", res.Injected)
	}
	if _, ok := os.LookupEnv("SKIP"); ok {
		t.Errorf("SKIP should not be set")
	}
}

func TestEject(t *testing.T) {
	os.Setenv("EJECT_A", "1")
	os.Setenv("EJECT_B", "2")
	inj := NewInjector(InjectOptions{Prefix: "EJECT_"})
	removed := inj.Eject(map[string]string{"A": "1", "B": "2"})
	if len(removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(removed))
	}
	if _, ok := os.LookupEnv("EJECT_A"); ok {
		t.Errorf("EJECT_A should have been removed")
	}
}

func TestInjectSummaryString(t *testing.T) {
	os.Unsetenv("S_X")
	defer os.Unsetenv("S_X")
	inj := NewInjector(InjectOptions{Overwrite: true, Prefix: "S_"})
	res, _ := inj.Inject(map[string]string{"X": "v"})
	s := res.SummaryString()
	if s == "" {
		t.Errorf("summary string should not be empty")
	}
}

func TestInjectOverwrite(t *testing.T) {
	os.Setenv("OW_KEY", "old")
	defer os.Unsetenv("OW_KEY")
	inj := NewInjector(InjectOptions{Overwrite: true, Prefix: "OW_"})
	_, err := inj.Inject(map[string]string{"KEY": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if os.Getenv("OW_KEY") != "new" {
		t.Errorf("expected OW_KEY=new, got %s", os.Getenv("OW_KEY"))
	}
}
