package env

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergeStrategyOurs keeps the local value on conflict.
	MergeStrategyOurs MergeStrategy = iota
	// MergeStrategyTheirs keeps the remote value on conflict.
	MergeStrategyTheirs
)

// MergeResult holds the merged environment and metadata about the operation.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []MergeConflict
	Added     []string
	Removed   []string
}

// MergeConflict describes a key whose value differed between base and theirs.
type MergeConflict struct {
	Key   string
	Ours  string
	Theirs string
	Resolved string
}

// Merge performs a three-way merge of environment variable maps.
// base is the common ancestor, ours is the local version, theirs is the remote.
// The supplied strategy is used to resolve conflicts automatically.
func Merge(base, ours, theirs map[string]string, strategy MergeStrategy) MergeResult {
	result := MergeResult{
		Merged: make(map[string]string),
	}

	// Collect all keys across all three maps.
	allKeys := make(map[string]struct{})
	for k := range base {
		allKeys[k] = struct{}{}
	}
	for k := range ours {
		allKeys[k] = struct{}{}
	}
	for k := range theirs {
		allKeys[k] = struct{}{}
	}

	for key := range allKeys {
		baseVal, inBase := base[key]
		oursVal, inOurs := ours[key]
		theirsVal, inTheirs := theirs[key]

		switch {
		case !inOurs && !inTheirs:
			// Deleted on both sides — omit.
			result.Removed = append(result.Removed, key)

		case !inOurs:
			// We deleted it; theirs still has it — respect our deletion.
			if inBase {
				result.Removed = append(result.Removed, key)
			} else {
				// They added it, we never had it — take theirs.
				result.Merged[key] = theirsVal
				result.Added = append(result.Added, key)
			}

		case !inTheirs:
			// They deleted it; we still have it — keep ours.
			result.Merged[key] = oursVal

		case oursVal == theirsVal:
			// Both agree.
			result.Merged[key] = oursVal

		case oursVal == baseVal:
			// Only theirs changed — take theirs.
			result.Merged[key] = theirsVal

		case theirsVal == baseVal:
			// Only ours changed — keep ours.
			result.Merged[key] = oursVal

		default:
			// True conflict — apply strategy.
			resolved := oursVal
			if strategy == MergeStrategyTheirs {
				resolved = theirsVal
			}
			result.Merged[key] = resolved
			result.Conflicts = append(result.Conflicts, MergeConflict{
				Key:      key,
				Ours:     oursVal,
				Theirs:   theirsVal,
				Resolved: resolved,
			})
		}
		_ = baseVal
	}

	return result
}
