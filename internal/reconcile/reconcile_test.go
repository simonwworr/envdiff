package reconcile_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reconcile"
)

func makeResult() diff.Result {
	return diff.Result{
		Added: []diff.Entry{
			{Key: "NEW_KEY", TargetValue: "new_val"},
		},
		Removed: []diff.Entry{
			{Key: "OLD_KEY", SourceValue: "old_val"},
		},
		Changed: []diff.Entry{
			{Key: "CHANGED_KEY", SourceValue: "v1", TargetValue: "v2"},
		},
		Unchanged: []diff.Entry{
			{Key: "SAME_KEY", SourceValue: "same", TargetValue: "same"},
		},
	}
}

func TestBuildPlan_HasChanges(t *testing.T) {
	plan := reconcile.BuildPlan(makeResult())
	if !plan.HasChanges() {
		t.Fatal("expected plan to have changes")
	}
	if len(plan.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(plan.Steps))
	}
}

func TestBuildPlan_StepActions(t *testing.T) {
	plan := reconcile.BuildPlan(makeResult())
	actions := map[string]reconcile.Action{}
	for _, s := range plan.Steps {
		actions[s.Key] = s.Action
	}
	if actions["NEW_KEY"] != reconcile.ActionAdd {
		t.Errorf("expected ADD for NEW_KEY, got %s", actions["NEW_KEY"])
	}
	if actions["OLD_KEY"] != reconcile.ActionRemove {
		t.Errorf("expected REMOVE for OLD_KEY, got %s", actions["OLD_KEY"])
	}
	if actions["CHANGED_KEY"] != reconcile.ActionUpdate {
		t.Errorf("expected UPDATE for CHANGED_KEY, got %s", actions["CHANGED_KEY"])
	}
}

func TestBuildPlan_EmptyResult(t *testing.T) {
	plan := reconcile.BuildPlan(diff.Result{})
	if plan.HasChanges() {
		t.Fatal("expected no changes for empty result")
	}
}

func TestPlan_Summary(t *testing.T) {
	plan := reconcile.BuildPlan(makeResult())
	summary := plan.Summary()
	if !strings.Contains(summary, "NEW_KEY") {
		t.Error("summary missing NEW_KEY")
	}
	if !strings.Contains(summary, "OLD_KEY") {
		t.Error("summary missing OLD_KEY")
	}
}

func TestApply(t *testing.T) {
	env := map[string]string{
		"OLD_KEY":     "old_val",
		"CHANGED_KEY": "v1",
		"SAME_KEY":    "same",
	}
	plan := reconcile.BuildPlan(makeResult())
	result := reconcile.Apply(env, plan)

	if _, ok := result["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if result["NEW_KEY"] != "new_val" {
		t.Errorf("expected NEW_KEY=new_val, got %s", result["NEW_KEY"])
	}
	if result["CHANGED_KEY"] != "v2" {
		t.Errorf("expected CHANGED_KEY=v2, got %s", result["CHANGED_KEY"])
	}
	if result["SAME_KEY"] != "same" {
		t.Error("SAME_KEY should be unchanged")
	}
}
