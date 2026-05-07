package reconcile

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// PatchMode controls how missing keys are handled during patch generation.
type PatchMode int

const (
	// PatchAddMissing includes keys present in other but missing from base.
	PatchAddMissing PatchMode = iota
	// PatchChangedOnly includes only keys whose values differ.
	PatchChangedOnly
)

// Patch represents a set of changes to apply to an env file.
type Patch struct {
	Lines []string
	Mode  PatchMode
}

// GeneratePatch builds a patch from a diff result that can be written
// as an env fragment to bring base in line with other.
func GeneratePatch(result diff.Result, mode PatchMode) Patch {
	var lines []string

	for _, entry := range result.Entries {
		switch entry.Status {
		case diff.Added:
			if mode == PatchAddMissing {
				lines = append(lines, formatLine(entry.Key, entry.OtherValue))
			}
		case diff.Changed:
			lines = append(lines, formatLine(entry.Key, entry.OtherValue))
		}
	}

	return Patch{Lines: lines, Mode: mode}
}

// String renders the patch as env file content.
func (p Patch) String() string {
	if len(p.Lines) == 0 {
		return ""
	}
	header := fmt.Sprintf("# envdiff patch (mode=%s)\n", p.modeLabel())
	return header + strings.Join(p.Lines, "\n") + "\n"
}

// Empty reports whether the patch contains no changes.
func (p Patch) Empty() bool {
	return len(p.Lines) == 0
}

func (p Patch) modeLabel() string {
	if p.Mode == PatchChangedOnly {
		return "changed-only"
	}
	return "add-missing"
}
