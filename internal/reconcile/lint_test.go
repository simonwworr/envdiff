package reconcile

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeLintEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_NAME", BaseVal: "myapp", OtherVal: "myapp", Status: diff.Unchanged},
		{Key: "db_host", BaseVal: "localhost", OtherVal: "localhost", Status: diff.Unchanged},
		{Key: "EMPTY_KEY", BaseVal: "", OtherVal: "", Status: diff.Unchanged},
		{Key: "BAD KEY", BaseVal: "val", OtherVal: "val", Status: diff.Unchanged},
	}
}

func TestLint_NoIssuesForCleanEntries(t *testing.T) {
	entries := []diff.Entry{
		{Key: "APP_NAME", BaseVal: "myapp", OtherVal: "myapp", Status: diff.Unchanged},
		{Key: "PORT", BaseVal: "8080", OtherVal: "8080", Status: diff.Unchanged},
	}
	report := Lint("base.env", entries)
	if len(report.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(report.Issues))
	}
}

func TestLint_DetectsLowercaseKey(t *testing.T) {
	entries := []diff.Entry{
		{Key: "db_host", BaseVal: "localhost", OtherVal: "localhost", Status: diff.Unchanged},
	}
	report := Lint("base.env", entries)
	if len(report.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(report.Issues))
	}
	if report.Issues[0].Severity != LintWarn {
		t.Errorf("expected warn severity")
	}
	if !strings.Contains(report.Issues[0].Message, "uppercase") {
		t.Errorf("expected uppercase message, got: %s", report.Issues[0].Message)
	}
}

func TestLint_DetectsKeyWithSpaces(t *testing.T) {
	entries := []diff.Entry{
		{Key: "BAD KEY", BaseVal: "val", OtherVal: "val", Status: diff.Unchanged},
	}
	report := Lint("base.env", entries)
	if !report.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
	if report.Issues[0].Severity != LintError {
		t.Errorf("expected error severity for key with spaces")
	}
}

func TestLint_DetectsEmptyValue(t *testing.T) {
	entries := []diff.Entry{
		{Key: "EMPTY_KEY", BaseVal: "", OtherVal: "", Status: diff.Unchanged},
	}
	report := Lint("base.env", entries)
	if len(report.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(report.Issues))
	}
	if report.Issues[0].Severity != LintWarn {
		t.Errorf("expected warn for empty value")
	}
}

func TestLintReport_Summary_NoIssues(t *testing.T) {
	report := &LintReport{File: "prod.env", Issues: nil}
	summary := report.Summary()
	if !strings.Contains(summary, "no lint issues") {
		t.Errorf("expected no-issues message, got: %s", summary)
	}
}

func TestLintReport_Summary_WithIssues(t *testing.T) {
	report := Lint("test.env", makeLintEntries())
	summary := report.Summary()
	if !strings.Contains(summary, "test.env") {
		t.Errorf("expected filename in summary")
	}
	if !strings.Contains(summary, "[warn]") && !strings.Contains(summary, "[error]") {
		t.Errorf("expected severity tags in summary")
	}
}
