package reconcile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Action represents what should be done to reconcile an env file.
type Action string

const (
	ActionAdd    Action = "ADD"
	ActionRemove Action = "REMOVE"
	ActionUpdate Action = "UPDATE"
)

// Step describes a single reconciliation step.
type Step struct {
	Action Action
	Key    string
	Value  string // target value (empty for REMOVE)
}

// Plan holds all steps needed to bring source in line with target.
type Plan struct {
	Steps []Step
}

// HasChanges reports whether the plan contains any steps.
func (p *Plan) HasChanges() bool {
	return len(p.Steps) > 0
}

// Summary returns a human-readable summary of the plan.
func (p *Plan) Summary() string {
	if !p.HasChanges() {
		return "No changes required."
	}
	var sb strings.Builder
	for _, s := range p.Steps {
		switch s.Action {
		case ActionAdd:
			fmt.Fprintf(&sb, "+ %s=%s\n", s.Key, s.Value)
		case ActionRemove:
			fmt.Fprintf(&sb, "- %s\n", s.Key)
		case ActionUpdate:
			fmt.Fprintf(&sb, "~ %s=%s\n", s.Key, s.Value)
		}
	}
	return sb.String()
}

// BuildPlan creates a reconciliation plan from a diff result.
// The plan describes changes needed to make source match target.
func BuildPlan(result diff.Result) Plan {
	var steps []Step

	for _, r := range result.Added {
		steps = append(steps, Step{Action: ActionAdd, Key: r.Key, Value: r.TargetValue})
	}
	for _, r := range result.Removed {
		steps = append(steps, Step{Action: ActionRemove, Key: r.Key})
	}
	for _, r := range result.Changed {
		steps = append(steps, Step{Action: ActionUpdate, Key: r.Key, Value: r.TargetValue})
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Key < steps[j].Key
	})

	return Plan{Steps: steps}
}

// Apply merges the plan steps into the provided env map (in-place).
func Apply(env map[string]string, plan Plan) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}
	for _, s := range plan.Steps {
		switch s.Action {
		case ActionAdd, ActionUpdate:
			result[s.Key] = s.Value
		case ActionRemove:
			delete(result, s.Key)
		}
	}
	return result
}
