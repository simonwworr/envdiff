package reconcile

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// RedactedEntry holds a key with its redacted representation.
type RedactedEntry struct {
	Key      string
	Original string
	Redacted string
	WasMasked bool
}

// RedactReport contains all redacted entries and a summary.
type RedactReport struct {
	Entries      []RedactedEntry
	TotalKeys    int
	RedactedKeys int
}

// Summary returns a human-readable summary of the redaction report.
func (r *RedactReport) Summary() string {
	return fmt.Sprintf("%d/%d keys redacted", r.RedactedKeys, r.TotalKeys)
}

// Redact takes a diff result and a masker, and returns a RedactReport
// where sensitive values are replaced with a redacted placeholder.
func Redact(results []diff.Result, masker *diff.Masker) *RedactReport {
	const placeholder = "[REDACTED]"

	report := &RedactReport{}

	for _, r := range results {
		key := r.Key
		value := pickValue(r)

		entry := RedactedEntry{
			Key:      key,
			Original: value,
			Redacted: value,
			WasMasked: false,
		}

		if masker != nil && masker.IsSensitive(key) && strings.TrimSpace(value) != "" {
			entry.Redacted = placeholder
			entry.WasMasked = true
			report.RedactedKeys++
		}

		report.Entries = append(report.Entries, entry)
		report.TotalKeys++
	}

	return report
}

// pickValue selects the most relevant value from a diff result.
// Prefers BaseValue; falls back to OtherValue for added entries.
func pickValue(r diff.Result) string {
	if r.BaseValue != "" {
		return r.BaseValue
	}
	return r.OtherValue
}
