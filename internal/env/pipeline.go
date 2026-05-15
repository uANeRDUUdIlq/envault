package env

import "fmt"

// PipelineStage represents a single transformation step applied to env vars.
type PipelineStage struct {
	Name    string
	Transform func(vars map[string]string) (map[string]string, error)
}

// PipelineResult holds the output of a pipeline run.
type PipelineResult struct {
	Vars   map[string]string
	Stages []string // names of stages executed
}

// Pipeline applies a sequence of transformation stages to env vars.
type Pipeline struct {
	stages []PipelineStage
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage appends a transformation stage to the pipeline.
func (p *Pipeline) AddStage(stage PipelineStage) {
	p.stages = append(p.stages, stage)
}

// Run executes all stages in order, passing vars through each.
// Execution stops and returns an error if any stage fails.
func (p *Pipeline) Run(input map[string]string) (*PipelineResult, error) {
	current := copyVars(input)
	executed := make([]string, 0, len(p.stages))

	for _, stage := range p.stages {
		result, err := stage.Transform(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q failed: %w", stage.Name, err)
		}
		current = result
		executed = append(executed, stage.Name)
	}

	return &PipelineResult{
		Vars:   current,
		Stages: executed,
	}, nil
}

// SummaryLines returns a human-readable summary of the pipeline result.
func (r *PipelineResult) SummaryLines() []string {
	lines := make([]string, 0, len(r.Stages)+1)
	for _, s := range r.Stages {
		lines = append(lines, fmt.Sprintf("  [stage] %s", s))
	}
	lines = append(lines, fmt.Sprintf("  [result] %d keys", len(r.Vars)))
	return lines
}

func copyVars(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
