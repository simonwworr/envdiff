package reconcile

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeValidateResult(entries []diff.Entry) diff.Result {
	return diff.Result{Entries: entries}
}

func TestValidate_NoErrors(t *testing.T) {
	result := makeValidateResult([]diff.Entry{
		{Key: "HOST", Status: diff.Unchanged, BaseValue: "localhost", OtherValue: "localhost"},
		{Key: "PORT", Status: diff.Removed, BaseValue: "8080", OtherValue: ""},
	})
	report := Validate(result, false)
	if report.HasErrors() {
		t.Errorf("expected no errors, got: %s", report.Summary())
	}
}

func TestValidate_AddedRequireAll(t *testing.T) {
	result := makeValidateResult([]diff.Entry{
		{Key: "NEW_KEY", Status: diff.Added, BaseValue: "", OtherValue: "value"},
	})
	report := Validate(result, true)
	if !report.HasErrors() {
		t.Fatal("expected errors for added key with requireAll=true")
	}
	if len(report.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(report.Errors))
	}
	if report.Errors[0].Key != "NEW_KEY" {
		t.Errorf("unexpected error key: %s", report.Errors[0].Key)
	}
}

func TestValidate_AddedNoRequireAll(t *testing.T) {
	result := makeValidateResult([]diff.Entry{
		{Key: "NEW_KEY", Status: diff.Added, BaseValue: "", OtherValue: "value"},
	})
	report := Validate(result, false)
	if report.HasErrors() {
		t.Errorf("expected no errors when requireAll=false, got: %s", report.Summary())
	}
}

func TestValidate_ChangedEmptyBase(t *testing.T) {
	result := makeValidateResult([]diff.Entry{
		{Key: "DB_PASS", Status: diff.Changed, BaseValue: "", OtherValue: "secret"},
	})
	report := Validate(result, false)
	if !report.HasErrors() {
		t.Fatal("expected error for changed entry with empty base value")
	}
	if report.Errors[0].Key != "DB_PASS" {
		t.Errorf("unexpected error key: %s", report.Errors[0].Key)
	}
}

func TestValidationReport_Summary_NoErrors(t *testing.T) {
	report := &ValidationReport{}
	got := report.Summary()
	if got != "validation passed: no issues found" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestValidationReport_Summary_WithErrors(t *testing.T) {
	report := &ValidationReport{
		Errors: []ValidationError{
			{Key: "FOO", Message: "some issue"},
			{Key: "BAR", Message: "another issue"},
		},
	}
	summary := report.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	for _, key := range []string{"FOO", "BAR"} {
		if !containsStr(summary, key) {
			t.Errorf("expected summary to contain key %q", key)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
