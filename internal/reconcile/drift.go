package reconcile

import (
	"fmt"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// DriftReport summarizes how far the current env has drifted from a snapshot.
type DriftReport struct {
	SnapshotFile string            `json:"snapshot_file"`
	CheckedAt    time.Time         `json:"checked_at"`
	Added        []string          `json:"added"`
	Removed      []string          `json:"removed"`
	Changed      []string          `json:"changed"`
	DriftScore   int               `json:"drift_score"`
	Entries      []diff.Entry      `json:"entries"`
}

// HasDrift returns true when any keys have been added, removed, or changed.
func (r *DriftReport) HasDrift() bool {
	return len(r.Added)+len(r.Removed)+len(r.Changed) > 0
}

// Summary returns a human-readable one-line summary of the drift report.
func (r *DriftReport) Summary() string {
	if !r.HasDrift() {
		return "no drift detected"
	}
	return fmt.Sprintf("drift detected: +%d added, -%d removed, ~%d changed (score: %d)",
		len(r.Added), len(r.Removed), len(r.Changed), r.DriftScore)
}

// DetectDrift compares a loaded snapshot against a live diff result and
// produces a DriftReport describing what has changed since the snapshot.
func DetectDrift(snapshotFile string, snap *Snapshot, result []diff.Entry) *DriftReport {
	report := &DriftReport{
		SnapshotFile: snapshotFile,
		CheckedAt:    time.Now().UTC(),
		Entries:      result,
	}

	snapshotKeys := make(map[string]string, len(snap.Keys))
	for _, sk := range snap.Keys {
		snapshotKeys[sk.Key] = sk.BaseValue
	}

	for _, e := range result {
		switch e.Status {
		case diff.Added:
			report.Added = append(report.Added, e.Key)
			report.DriftScore += 2
		case diff.Removed:
			report.Removed = append(report.Removed, e.Key)
			report.DriftScore += 3
		case diff.Changed:
			if prev, ok := snapshotKeys[e.Key]; ok && prev == e.BaseValue {
				// changed since snapshot
				report.Changed = append(report.Changed, e.Key)
				report.DriftScore++
			}
		}
	}

	return report
}
