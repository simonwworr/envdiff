// Package output provides formatting utilities for displaying diff results
// in human-readable and machine-readable formats.
package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format controls the output style for diff results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// WriteDiff writes a human-readable diff of the results to w.
// Added lines are prefixed with '+', removed with '-', changed with '~'.
func WriteDiff(w io.Writer, results []diff.Result, masked bool) {
	keys := make([]string, 0, len(results))
	byKey := make(map[string]diff.Result, len(results))
	for _, r := range results {
		keys = append(keys, r.Key)
		byKey[r.Key] = r
	}
	sort.Strings(keys)

	for _, k := range keys {
		r := byKey[k]
		switch r.Status {
		case diff.Added:
			fmt.Fprintf(w, "+ %s=%s\n", r.Key, displayValue(r.ValueB, masked, r.Masked))
		case diff.Removed:
			fmt.Fprintf(w, "- %s=%s\n", r.Key, displayValue(r.ValueA, masked, r.Masked))
		case diff.Changed:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", r.Key,
				displayValue(r.ValueA, masked, r.Masked),
				displayValue(r.ValueB, masked, r.Masked))
		case diff.Unchanged:
			fmt.Fprintf(w, "  %s=%s\n", r.Key, displayValue(r.ValueA, masked, r.Masked))
		}
	}
}

// Summary returns a one-line summary of the diff results.
func Summary(results []diff.Result) string {
	var added, removed, changed, unchanged int
	for _, r := range results {
		switch r.Status {
		case diff.Added:
			added++
		case diff.Removed:
			removed++
		case diff.Changed:
			changed++
		case diff.Unchanged:
			unchanged++
		}
	}
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}
	if unchanged > 0 {
		parts = append(parts, fmt.Sprintf("%d unchanged", unchanged))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}

func displayValue(val string, masked bool, isSensitive bool) string {
	if masked && isSensitive {
		return "***"
	}
	return val
}
