package reconcile

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// AuditEntry records a single change event with metadata.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Key       string    `json:"key"`
	Action    string    `json:"action"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
	Masked    bool      `json:"masked"`
}

// AuditLog holds a sequence of audit entries for a reconciliation run.
type AuditLog struct {
	CreatedAt time.Time    `json:"created_at"`
	BaseFile  string       `json:"base_file"`
	OtherFile string       `json:"other_file"`
	Entries   []AuditEntry `json:"entries"`
}

// NewAuditLog creates an AuditLog from a diff result.
func NewAuditLog(baseFile, otherFile string, results []diff.Result, masker *diff.Masker) *AuditLog {
	log := &AuditLog{
		CreatedAt: time.Now().UTC(),
		BaseFile:  baseFile,
		OtherFile: otherFile,
		Entries:   make([]AuditEntry, 0, len(results)),
	}
	for _, r := range results {
		entry := AuditEntry{
			Timestamp: log.CreatedAt,
			Key:       r.Key,
			Action:    string(r.Status),
			Masked:    masker.IsSensitive(r.Key),
		}
		if entry.Masked {
			entry.OldValue = maskIfSet(r.BaseValue)
			entry.NewValue = maskIfSet(r.OtherValue)
		} else {
			entry.OldValue = r.BaseValue
			entry.NewValue = r.OtherValue
		}
		log.Entries = append(log.Entries, entry)
	}
	return log
}

func maskIfSet(v string) string {
	if v == "" {
		return ""
	}
	return "***"
}

// Summary returns a human-readable summary of the audit log.
func (a *AuditLog) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Audit log: %s vs %s (%s)\n",
		a.BaseFile, a.OtherFile, a.CreatedAt.Format(time.RFC3339))
	for _, e := range a.Entries {
		fmt.Fprintf(&sb, "  [%s] %s", e.Action, e.Key)
		if e.OldValue != "" || e.NewValue != "" {
			fmt.Fprintf(&sb, " (%q -> %q)", e.OldValue, e.NewValue)
		}
		if e.Masked {
			fmt.Fprint(&sb, " [masked]")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
