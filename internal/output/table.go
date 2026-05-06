package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// TableWriter renders diff results as an aligned table.
type TableWriter struct {
	w       io.Writer
	masker  *diff.Masker
	padding int
}

// NewTableWriter creates a TableWriter with optional masking.
func NewTableWriter(w io.Writer, masker *diff.Masker) *TableWriter {
	return &TableWriter{w: w, masker: masker, padding: 2}
}

// Write renders the diff result as a formatted table.
func (t *TableWriter) Write(result diff.Result) error {
	if len(result.Entries) == 0 {
		fmt.Fprintln(t.w, "No differences found.")
		return nil
	}

	keyWidth := t.maxKeyWidth(result)
	statusWidth := 9 // len("UNCHANGED")

	header := fmt.Sprintf("%-*s  %-*s  %-20s  %-20s",
		keyWidth, "KEY",
		statusWidth, "STATUS",
		"BASE",
		"TARGET",
	)
	fmt.Fprintln(t.w, header)
	fmt.Fprintln(t.w, strings.Repeat("-", len(header)))

	for _, e := range result.Entries {
		base := displayValue(e.BaseValue, t.masker, e.Key)
		target := displayValue(e.TargetValue, t.masker, e.Key)
		line := fmt.Sprintf("%-*s  %-*s  %-20s  %-20s",
			keyWidth, e.Key,
			statusWidth, e.Status,
			truncate(base, 20),
			truncate(target, 20),
		)
		fmt.Fprintln(t.w, line)
	}
	return nil
}

func (t *TableWriter) maxKeyWidth(result diff.Result) int {
	max := len("KEY")
	for _, e := range result.Entries {
		if len(e.Key) > max {
			max = len(e.Key)
		}
	}
	return max + t.padding
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
