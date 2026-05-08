package reconcile

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
)

func makeSnapshotResult() diff.Result {
	return diff.Result{
		Entries: []diff.Entry{
			{Key: "HOST", BaseValue: "localhost", OtherValue: "prod.example.com", Status: diff.Changed},
			{Key: "PORT", BaseValue: "", OtherValue: "8080", Status: diff.Added},
			{Key: "OLD_KEY", BaseValue: "val", OtherValue: "", Status: diff.Removed},
			{Key: "DEBUG", BaseValue: "true", OtherValue: "true", Status: diff.Unchanged},
		},
	}
}

func TestTakeSnapshot_Summary(t *testing.T) {
	result := makeSnapshotResult()
	snap := TakeSnapshot("base.env", "prod.env", result)

	if snap.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", snap.Summary.Added)
	}
	if snap.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", snap.Summary.Removed)
	}
	if snap.Summary.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", snap.Summary.Changed)
	}
	if snap.Summary.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", snap.Summary.Unchanged)
	}
}

func TestTakeSnapshot_Metadata(t *testing.T) {
	before := time.Now().UTC()
	snap := TakeSnapshot("a.env", "b.env", makeSnapshotResult())
	after := time.Now().UTC()

	if snap.BaseFile != "a.env" {
		t.Errorf("expected base a.env, got %s", snap.BaseFile)
	}
	if snap.OtherFile != "b.env" {
		t.Errorf("expected other b.env, got %s", snap.OtherFile)
	}
	if snap.CreatedAt.Before(before) || snap.CreatedAt.After(after) {
		t.Errorf("unexpected CreatedAt: %v", snap.CreatedAt)
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := TakeSnapshot("base.env", "other.env", makeSnapshotResult())
	if err := SaveSnapshot(path, orig); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if loaded.BaseFile != orig.BaseFile {
		t.Errorf("BaseFile mismatch: got %s", loaded.BaseFile)
	}
	if loaded.Summary.Added != orig.Summary.Added {
		t.Errorf("Summary.Added mismatch")
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("Entries length mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveSnapshot_BadPath(t *testing.T) {
	snap := TakeSnapshot("a.env", "b.env", makeSnapshotResult())
	err := SaveSnapshot("/nonexistent/dir/snap.json", snap)
	if err == nil {
		t.Error("expected error for bad path")
	}
	_ = os.Remove("/nonexistent/dir/snap.json")
}
