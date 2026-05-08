package reconcile

import (
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
)

func makeSnapshotForDrift(keys map[string]string) *Snapshot {
	snap := &Snapshot{
		CreatedAt: time.Now().UTC(),
		BaseFile:  "base.env",
		OtherFile: "other.env",
	}
	for k, v := range keys {
		snap.Keys = append(snap.Keys, SnapshotKey{Key: k, BaseValue: v, Status: string(diff.Unchanged)})
	}
	return snap
}

func TestDetectDrift_NoDrift(t *testing.T) {
	snap := makeSnapshotForDrift(map[string]string{"HOST": "localhost"})
	entries := []diff.Entry{
		{Key: "HOST", BaseValue: "localhost", OtherValue: "localhost", Status: diff.Unchanged},
	}
	report := DetectDrift("snap.json", snap, entries)
	if report.HasDrift() {
		t.Errorf("expected no drift, got: %s", report.Summary())
	}
	if report.DriftScore != 0 {
		t.Errorf("expected drift score 0, got %d", report.DriftScore)
	}
}

func TestDetectDrift_AddedKey(t *testing.T) {
	snap := makeSnapshotForDrift(map[string]string{})
	entries := []diff.Entry{
		{Key: "NEW_KEY", BaseValue: "", OtherValue: "value", Status: diff.Added},
	}
	report := DetectDrift("snap.json", snap, entries)
	if !report.HasDrift() {
		t.Error("expected drift to be detected")
	}
	if len(report.Added) != 1 || report.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in added, got %v", report.Added)
	}
	if report.DriftScore != 2 {
		t.Errorf("expected drift score 2, got %d", report.DriftScore)
	}
}

func TestDetectDrift_RemovedKey(t *testing.T) {
	snap := makeSnapshotForDrift(map[string]string{"OLD_KEY": "old"})
	entries := []diff.Entry{
		{Key: "OLD_KEY", BaseValue: "old", OtherValue: "", Status: diff.Removed},
	}
	report := DetectDrift("snap.json", snap, entries)
	if len(report.Removed) != 1 {
		t.Errorf("expected 1 removed key, got %v", report.Removed)
	}
	if report.DriftScore != 3 {
		t.Errorf("expected drift score 3, got %d", report.DriftScore)
	}
}

func TestDetectDrift_Summary_WithDrift(t *testing.T) {
	snap := makeSnapshotForDrift(map[string]string{"A": "1"})
	entries := []diff.Entry{
		{Key: "A", BaseValue: "1", OtherValue: "2", Status: diff.Changed},
		{Key: "B", BaseValue: "", OtherValue: "new", Status: diff.Added},
	}
	report := DetectDrift("snap.json", snap, entries)
	summary := report.Summary()
	if summary == "no drift detected" {
		t.Error("expected drift summary, got no-drift message")
	}
}

func TestDetectDrift_SnapshotFile(t *testing.T) {
	snap := makeSnapshotForDrift(map[string]string{})
	report := DetectDrift("my-snap.json", snap, nil)
	if report.SnapshotFile != "my-snap.json" {
		t.Errorf("expected snapshot file my-snap.json, got %s", report.SnapshotFile)
	}
}
