package env

import (
	"testing"
)

func TestPipelineEmptyRun(t *testing.T) {
	p := NewPipeline()
	input := map[string]string{"FOO": "bar"}
	res, err := p.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Vars["FOO"])
	}
	if len(res.Stages) != 0 {
		t.Errorf("expected 0 stages, got %d", len(res.Stages))
	}
}

func TestPipelinePrefixStage(t *testing.T) {
	p := NewPipeline()
	p.AddStage(PrefixStage("APP_"))
	res, err := p.Run(map[string]string{"NAME": "envault"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["APP_NAME"]; !ok {
		t.Errorf("expected key APP_NAME")
	}
}

func TestPipelineUppercaseKeysStage(t *testing.T) {
	p := NewPipeline()
	p.AddStage(UppercaseKeysStage())
	res, err := p.Run(map[string]string{"foo": "1", "bar": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "1" || res.Vars["BAR"] != "2" {
		t.Errorf("uppercase failed: %v", res.Vars)
	}
}

func TestPipelineFilterKeysStage(t *testing.T) {
	p := NewPipeline()
	p.AddStage(FilterKeysStage([]string{"KEEP"}))
	res, err := p.Run(map[string]string{"KEEP": "yes", "DROP": "no"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["DROP"]; ok {
		t.Errorf("expected DROP to be filtered out")
	}
	if res.Vars["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes")
	}
}

func TestPipelineRequireKeysMissing(t *testing.T) {
	p := NewPipeline()
	p.AddStage(RequireKeysStage([]string{"MUST_EXIST"}))
	_, err := p.Run(map[string]string{"OTHER": "val"})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestPipelineMultipleStages(t *testing.T) {
	p := NewPipeline()
	p.AddStage(UppercaseKeysStage())
	p.AddStage(PrefixStage("X_"))
	p.AddStage(RequireKeysStage([]string{"X_FOO"}))
	res, err := p.Run(map[string]string{"foo": "42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["X_FOO"] != "42" {
		t.Errorf("expected X_FOO=42, got %q", res.Vars["X_FOO"])
	}
	if len(res.Stages) != 3 {
		t.Errorf("expected 3 stages recorded, got %d", len(res.Stages))
	}
}

func TestPipelineSummaryLines(t *testing.T) {
	p := NewPipeline()
	p.AddStage(UppercaseKeysStage())
	res, _ := p.Run(map[string]string{"a": "1", "b": "2"})
	lines := res.SummaryLines()
	if len(lines) < 2 {
		t.Errorf("expected at least 2 summary lines, got %d", len(lines))
	}
}
