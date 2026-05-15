package reconcile

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/parser"
)

// ParseLintArgs parses CLI arguments for the lint subcommand.
func ParseLintArgs(args []string) (file string, errorsOnly bool, err error) {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	fs.Bool("errors-only", false, "only report error-level issues")
	if err = fs.Parse(args); err != nil {
		return
	}
	errorsOnly = fs.Lookup("errors-only").Value.String() == "true"
	if fs.NArg() < 1 {
		err = fmt.Errorf("usage: envdiff lint <file.env>")
		return
	}
	file = fs.Arg(0)
	return
}

// RunLint executes the lint command, writing results to w.
// Returns a non-nil error if lint errors are found or the file cannot be read.
func RunLint(file string, errorsOnly bool, w io.Writer) error {
	env, err := parser.Parse(file)
	if err != nil {
		return fmt.Errorf("lint: cannot parse %s: %w", file, err)
	}

	// Convert parsed map to entries for linting.
	entries := mapToEntries(env)
	report := Lint(file, entries)

	for _, iss := range report.Issues {
		if errorsOnly && iss.Severity != LintError {
			continue
		}
		fmt.Fprintf(w, "[%s] %s: %s\n", iss.Severity, iss.Key, iss.Message)
	}

	if report.HasErrors() {
		return fmt.Errorf("lint: %s has error-level issues", file)
	}
	return nil
}

// RunLintFromArgs is the CLI entry point for the lint subcommand.
func RunLintFromArgs(args []string) int {
	file, errorsOnly, err := ParseLintArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	if err := RunLint(file, errorsOnly, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
