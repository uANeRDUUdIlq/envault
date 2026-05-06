package env

import (
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for exported variables.
type ExportFormat string

const (
	FormatShell  ExportFormat = "shell"
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
)

// ExportOptions controls how variables are exported.
type ExportOptions struct {
	Format ExportFormat
	Export bool // prefix with `export` in shell format
	Keys   []string // if non-empty, only export these keys
}

// Exporter writes env vars to a destination in a given format.
type Exporter struct {
	opts ExportOptions
}

// NewExporter creates an Exporter with the given options.
func NewExporter(opts ExportOptions) *Exporter {
	if opts.Format == "" {
		opts.Format = FormatDotenv
	}
	return &Exporter{opts: opts}
}

// Export writes the vars to the given file path, or stdout if path is "-".
func (e *Exporter) Export(vars map[string]string, path string) error {
	var w *os.File
	var err error
	if path == "-" {
		w = os.Stdout
	} else {
		w, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("export: create file: %w", err)
		}
		defer w.Close()
	}

	filtered := e.filterKeys(vars)

	switch e.opts.Format {
	case FormatShell:
		_, err = fmt.Fprint(w, e.renderShell(filtered))
	case FormatJSON:
		_, err = fmt.Fprint(w, e.renderJSON(filtered))
	default:
		_, err = fmt.Fprint(w, Serialize(filtered))
	}
	return err
}

func (e *Exporter) filterKeys(vars map[string]string) map[string]string {
	if len(e.opts.Keys) == 0 {
		return vars
	}
	out := make(map[string]string, len(e.opts.Keys))
	for _, k := range e.opts.Keys {
		if v, ok := vars[k]; ok {
			out[k] = v
		}
	}
	return out
}

func (e *Exporter) renderShell(vars map[string]string) string {
	var sb strings.Builder
	prefix := ""
	if e.opts.Export {
		prefix = "export "
	}
	for k, v := range vars {
		fmt.Fprintf(&sb, "%s%s=%q\n", prefix, k, v)
	}
	return sb.String()
}

func (e *Exporter) renderJSON(vars map[string]string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	i := 0
	for k, v := range vars {
		comma := ","
		if i == len(vars)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, v, comma)
		i++
	}
	sb.WriteString("}\n")
	return sb.String()
}
