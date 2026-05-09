package reconcile

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// PromoteResult holds the outcome of a promotion between environments.
type PromoteResult struct {
	Source      string
	Target      string
	Applied     []string
	Skipped     []string
	Conflicts   []string
}

// PromoteOptions controls how promotion behaves.
type PromoteOptions struct {
	// OverwriteExisting replaces keys that already exist in target.
	OverwriteExisting bool
	// SkipSensitive avoids promoting keys matched by the masker.
	SkipSensitive bool
	// SensitiveKeys is a set of key names to skip when SkipSensitive is true.
	SensitiveKeys map[string]bool
}

// Promote copies keys from source env map into target env map according to opts.
// It returns a PromoteResult describing what was applied, skipped, or conflicted.
func Promote(source, target map[string]string, opts PromoteOptions) (map[string]string, PromoteResult) {
	result := map[string]string{}
	for k, v := range target {
		result[k] = v
	}

	pr := PromoteResult{}

	for k, v := range source {
		if opts.SkipSensitive && opts.SensitiveKeys[k] {
			pr.Skipped = append(pr.Skipped, k)
			continue
		}
		if _, exists := target[k]; exists && !opts.OverwriteExisting {
			pr.Conflicts = append(pr.Conflicts, k)
			continue
		}
		result[k] = v
		pr.Applied = append(pr.Applied, k)
	}

	return result, pr
}

// PromoteSummary returns a human-readable summary of a PromoteResult.
func PromoteSummary(pr PromoteResult) string {
	return fmt.Sprintf(
		"promote %s -> %s: applied=%d skipped=%d conflicts=%d",
		pr.Source, pr.Target,
		len(pr.Applied), len(pr.Skipped), len(pr.Conflicts),
	)
}

// PromoteFromDiff builds a source map from diff entries that are Added or Changed
// so callers can promote only the delta from a comparison.
func PromoteFromDiff(results []diff.Result) map[string]string {
	out := map[string]string{}
	for _, r := range results {
		if r.Status == diff.Added || r.Status == diff.Changed {
			out[r.Key] = r.BaseValue
		}
	}
	return out
}
