package reconcile

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/config"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunValidate parses base and other env files, computes a diff, then runs
// validation. It writes a human-readable report to w and returns a non-zero
// exit code if validation fails.
func RunValidate(cfg *config.Config, requireAll bool, w io.Writer) int {
	base, err := parser.Parse(cfg.BaseFile)
	if err != nil {
		fmt.Fprintf(w, "error reading base file %q: %v\n", cfg.BaseFile, err)
		return 1
	}

	other, err := parser.Parse(cfg.OtherFile)
	if err != nil {
		fmt.Fprintf(w, "error reading other file %q: %v\n", cfg.OtherFile, err)
		return 1
	}

	result := diff.Compare(base, other)
	report := Validate(result, requireAll)

	fmt.Fprintln(w, report.Summary())

	if report.HasErrors() {
		return 2
	}
	return 0
}

// RunValidateFromArgs is a convenience wrapper that builds config from CLI args
// and calls RunValidate, writing output to stdout.
func RunValidateFromArgs(args []string, requireAll bool) int {
	cfg, err := config.FromArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid arguments: %v\n", err)
		return 1
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		return 1
	}
	return RunValidate(cfg, requireAll, os.Stdout)
}
