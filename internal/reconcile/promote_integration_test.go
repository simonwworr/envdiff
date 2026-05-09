package reconcile_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reconcile"
)

func TestPromoteIntegration_RoundTrip(t *testing.T) {
	// Simulate staging -> prod promotion with conflict avoidance.
	staging := map[string]string{
		"APP_ENV":    "staging",
		"LOG_LEVEL":  "debug",
		"DB_HOST":    "staging-db",
		"API_SECRET": "stg-secret",
	}
	prod := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "prod-db",
	}

	// Only promote Added/Changed keys from a diff.
	results := diff.Compare(staging, prod)
	delta := reconcile.PromoteFromDiff(results)

	opts := reconcile.PromoteOptions{
		OverwriteExisting: false,
		SkipSensitive:     true,
		SensitiveKeys:     map[string]bool{"API_SECRET": true},
	}

	promoted, pr := reconcile.Promote(delta, prod, opts)

	// LOG_LEVEL should be promoted (new key, not sensitive).
	if promoted["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", promoted["LOG_LEVEL"])
	}
	// APP_ENV should remain prod value (conflict, no overwrite).
	if promoted["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", promoted["APP_ENV"])
	}
	// API_SECRET should not appear.
	if _, ok := promoted["API_SECRET"]; ok {
		t.Error("API_SECRET should have been skipped")
	}
	if len(pr.Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(pr.Skipped))
	}
}

func TestPromoteIntegration_RenderOutput(t *testing.T) {
	src := map[string]string{"NEW_KEY": "value", "SHARED": "from-src"}
	dst := map[string]string{"SHARED": "from-dst"}

	promoted, _ := reconcile.Promote(src, dst, reconcile.PromoteOptions{OverwriteExisting: true})
	output := reconcile.RenderEnv(promoted)

	if !strings.Contains(output, "NEW_KEY=value") {
		t.Errorf("rendered output missing NEW_KEY: %q", output)
	}
	if !strings.Contains(output, "SHARED=from-src") {
		t.Errorf("rendered output should show overwritten SHARED: %q", output)
	}
}

func TestPromoteIntegration_SavePromotedFile(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}

	promoted, _ := reconcile.Promote(src, dst, reconcile.PromoteOptions{})
	output := reconcile.RenderEnv(promoted)

	tmp, err := os.CreateTemp(t.TempDir(), "promoted-*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.Close()

	if _, err := tmp.WriteString(output); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("saved file missing FOO=bar: %q", string(data))
	}
}
