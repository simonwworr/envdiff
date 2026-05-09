package reconcile

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestPromote_AppliesNewKeys(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	dst := map[string]string{"EXISTING": "yes"}

	out, pr := Promote(src, dst, PromoteOptions{})

	if len(pr.Applied) != 2 {
		t.Fatalf("expected 2 applied, got %d", len(pr.Applied))
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" || out["EXISTING"] != "yes" {
		t.Error("unexpected output map contents")
	}
}

func TestPromote_ConflictWithoutOverwrite(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}

	out, pr := Promote(src, dst, PromoteOptions{OverwriteExisting: false})

	if len(pr.Conflicts) != 1 || pr.Conflicts[0] != "FOO" {
		t.Errorf("expected conflict on FOO, got %v", pr.Conflicts)
	}
	if out["FOO"] != "old" {
		t.Error("conflict key should retain target value")
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := map[string]string{"FOO": "new"}
	dst := map[string]string{"FOO": "old"}

	out, pr := Promote(src, dst, PromoteOptions{OverwriteExisting: true})

	if len(pr.Applied) != 1 {
		t.Fatalf("expected 1 applied, got %d", len(pr.Applied))
	}
	if out["FOO"] != "new" {
		t.Error("expected overwritten value")
	}
}

func TestPromote_SkipSensitive(t *testing.T) {
	src := map[string]string{"SECRET_KEY": "abc", "APP_NAME": "myapp"}
	dst := map[string]string{}
	opts := PromoteOptions{
		SkipSensitive: true,
		SensitiveKeys: map[string]bool{"SECRET_KEY": true},
	}

	out, pr := Promote(src, dst, opts)

	if len(pr.Skipped) != 1 || pr.Skipped[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY skipped, got %v", pr.Skipped)
	}
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("sensitive key should not appear in output")
	}
	if out["APP_NAME"] != "myapp" {
		t.Error("non-sensitive key should be promoted")
	}
}

func TestPromoteSummary_Format(t *testing.T) {
	pr := PromoteResult{
		Source:    "staging",
		Target:    "prod",
		Applied:   []string{"A", "B"},
		Skipped:   []string{"C"},
		Conflicts: []string{},
	}
	s := PromoteSummary(pr)
	if !strings.Contains(s, "staging -> prod") {
		t.Errorf("summary missing source/target: %q", s)
	}
	if !strings.Contains(s, "applied=2") {
		t.Errorf("summary missing applied count: %q", s)
	}
}

func TestPromoteFromDiff_FiltersStatuses(t *testing.T) {
	results := []diff.Result{
		{Key: "ADDED_KEY", Status: diff.Added, BaseValue: "v1"},
		{Key: "CHANGED_KEY", Status: diff.Changed, BaseValue: "v2"},
		{Key: "REMOVED_KEY", Status: diff.Removed, BaseValue: "v3"},
		{Key: "SAME_KEY", Status: diff.Unchanged, BaseValue: "v4"},
	}

	out := PromoteFromDiff(results)

	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out["ADDED_KEY"] != "v1" || out["CHANGED_KEY"] != "v2" {
		t.Error("unexpected promote map values")
	}
}
