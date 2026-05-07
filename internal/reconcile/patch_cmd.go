package reconcile

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// PatchOptions configures patch generation behaviour.
type PatchOptions struct {
	Mode    PatchMode
	DryRun  bool
	Verbose bool
}

// RunPatch generates a patch from the given diff result and writes it to w.
// If DryRun is set, no patch content is written but a summary is printed.
func RunPatch(w io.Writer, result diff.Result, opts PatchOptions) error {
	patch := GeneratePatch(result, opts.Mode)

	if patch.Empty() {
		_, err := fmt.Fprintln(w, "# no changes to patch")
		return err
	}

	if opts.DryRun {
		_, err := fmt.Fprintf(w, "# dry-run: %d line(s) would be written\n", len(patch.Lines))
		return err
	}

	if opts.Verbose {
		_, err := fmt.Fprintf(w, "# generating patch with %d change(s)\n", len(patch.Lines))
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w, patch.String())
	return err
}
