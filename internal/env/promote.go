package env

import "fmt"

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	SourceNamespace string
	TargetNamespace string
	Copied          []string
	Skipped         []string
	Overwritten     []string
}

// SummaryString returns a human-readable summary of the promotion.
func (r *PromoteResult) SummaryString() string {
	return fmt.Sprintf(
		"promote %s -> %s: copied=%d skipped=%d overwritten=%d",
		r.SourceNamespace, r.TargetNamespace,
		len(r.Copied), len(r.Skipped), len(r.Overwritten),
	)
}

// PromoteOptions controls promotion behaviour.
type PromoteOptions struct {
	// Overwrite existing keys in the target namespace.
	Overwrite bool
	// Keys restricts promotion to specific keys; empty means all keys.
	Keys []string
}

// Promoter copies env vars from one namespace to another.
type Promoter struct {
	ns *NamespaceStore
}

// NewPromoter creates a Promoter backed by the given NamespaceStore.
func NewPromoter(ns *NamespaceStore) *Promoter {
	return &Promoter{ns: ns}
}

// Promote copies variables from src namespace into dst namespace.
func (p *Promoter) Promote(src, dst string, opts PromoteOptions) (*PromoteResult, error) {
	srcVars, err := p.ns.GetVars(src)
	if err != nil {
		return nil, fmt.Errorf("promote: read source %q: %w", src, err)
	}
	dstVars, err := p.ns.GetVars(dst)
	if err != nil {
		return nil, fmt.Errorf("promote: read target %q: %w", dst, err)
	}

	filter := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		filter[k] = struct{}{}
	}

	result := &PromoteResult{SourceNamespace: src, TargetNamespace: dst}

	for k, v := range srcVars {
		if len(filter) > 0 {
			if _, ok := filter[k]; !ok {
				continue
			}
		}
		_, exists := dstVars[k]
		if exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		dstVars[k] = v
		if exists {
			result.Overwritten = append(result.Overwritten, k)
		} else {
			result.Copied = append(result.Copied, k)
		}
	}

	if err := p.ns.SetVars(dst, dstVars); err != nil {
		return nil, fmt.Errorf("promote: write target %q: %w", dst, err)
	}
	return result, nil
}
