package reconcile

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// ValidationError represents a single validation issue found in a diff result.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationReport holds all errors found during validation.
type ValidationReport struct {
	Errors []ValidationError
}

func (r *ValidationReport) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationReport) Summary() string {
	if !r.HasErrors() {
		return "validation passed: no issues found"
	}
	lines := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		lines = append(lines, "  - "+e.Error())
	}
	return fmt.Sprintf("%d validation error(s):\n%s", len(r.Errors), strings.Join(lines, "\n"))
}

// Validate checks a diff.Result for common issues such as keys present in
// other but missing from base (added), which may indicate required variables
// that are not yet defined in the base environment.
func Validate(result diff.Result, requireAll bool) *ValidationReport {
	report := &ValidationReport{}

	for _, entry := range result.Entries {
		switch entry.Status {
		case diff.Added:
			if requireAll {
				report.Errors = append(report.Errors, ValidationError{
					Key:     entry.Key,
					Message: "key exists in other but is missing from base",
				})
			}
		case diff.Changed:
			if entry.BaseValue == "" {
				report.Errors = append(report.Errors, ValidationError{
					Key:     entry.Key,
					Message: "key changed but base value is empty",
				})
			}
		}
	}

	return report
}
