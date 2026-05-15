package reconcile

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// LintSeverity indicates the severity of a lint issue.
type LintSeverity string

const (
	LintWarn  LintSeverity = "warn"
	LintError LintSeverity = "error"
)

// LintIssue represents a single lint finding for a key.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintReport holds all issues found during linting.
type LintReport struct {
	Issues []LintIssue
	File   string
}

// HasErrors returns true if any issue has error severity.
func (r *LintReport) HasErrors() bool {
	for _, iss := range r.Issues {
		if iss.Severity == LintError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of the lint report.
func (r *LintReport) Summary() string {
	if len(r.Issues) == 0 {
		return fmt.Sprintf("%s: no lint issues found", r.File)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: %d issue(s)\n", r.File, len(r.Issues)))
	for _, iss := range r.Issues {
		sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n", iss.Severity, iss.Key, iss.Message))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Lint inspects a set of diff results and reports style/quality issues.
func Lint(file string, entries []diff.Entry) *LintReport {
	report := &LintReport{File: file}
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		if strings.ToUpper(e.Key) != e.Key {
			report.Issues = append(report.Issues, LintIssue{
				Key:      e.Key,
				Message:  "key should be uppercase",
				Severity: LintWarn,
			})
		}
		if strings.Contains(e.Key, " ") {
			report.Issues = append(report.Issues, LintIssue{
				Key:      e.Key,
				Message:  "key contains spaces",
				Severity: LintError,
			})
		}
		if e.BaseVal == "" && e.Status != diff.Added {
			report.Issues = append(report.Issues, LintIssue{
				Key:      e.Key,
				Message:  "key has empty value",
				Severity: LintWarn,
			})
		}
	}
	return report
}
