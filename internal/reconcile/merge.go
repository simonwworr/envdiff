package reconcile

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// PreferBase keeps the base file's value on conflict.
	PreferBase MergeStrategy = iota
	// PreferOther uses the other file's value on conflict.
	PreferOther
)

// MergeResult holds the merged key-value pairs and a log of decisions made.
type MergeResult struct {
	Env      map[string]string
	Conflicts []MergeConflict
}

// MergeConflict records a key where both files had differing values.
type MergeConflict struct {
	Key      string
	BaseVal  string
	OtherVal string
	Chosen   string
}

// Merge combines base and other env maps using the given strategy.
// Keys only in one map are always included. Conflicting keys are
// resolved according to strategy, and recorded in MergeResult.Conflicts.
func Merge(base, other map[string]string, strategy MergeStrategy) MergeResult {
	result := MergeResult{
		Env: make(map[string]string),
	}

	keys := unionKeys(base, other)
	sort.Strings(keys)

	for _, k := range keys {
		bVal, inBase := base[k]
		oVal, inOther := other[k]

		switch {
		case inBase && !inOther:
			result.Env[k] = bVal
		case !inBase && inOther:
			result.Env[k] = oVal
		case bVal == oVal:
			result.Env[k] = bVal
		default:
			var chosen string
			if strategy == PreferOther {
				chosen = oVal
			} else {
				chosen = bVal
			}
			result.Env[k] = chosen
			result.Conflicts = append(result.Conflicts, MergeConflict{
				Key:      k,
				BaseVal:  bVal,
				OtherVal: oVal,
				Chosen:   chosen,
			})
		}
	}

	return result
}

// unionKeys returns all unique keys from both maps.
// Reuses the helper already defined in internal/diff.
func unionKeys(a, b map[string]string) []string {
	return diff.UnionKeys(a, b)
}

// ConflictSummary returns a human-readable summary of merge conflicts.
func ConflictSummary(conflicts []MergeConflict) string {
	if len(conflicts) == 0 {
		return "no conflicts"
	}
	s := fmt.Sprintf("%d conflict(s):\n", len(conflicts))
	for _, c := range conflicts {
		s += fmt.Sprintf("  %s: base=%q other=%q -> chose=%q\n", c.Key, c.BaseVal, c.OtherVal, c.Chosen)
	}
	return s
}
