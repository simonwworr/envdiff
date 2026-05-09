package reconcile

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// PromoteArgs holds parsed CLI arguments for the promote sub-command.
type PromoteArgs struct {
	SourceFile        string
	TargetFile        string
	OutputFile        string
	OverwriteExisting bool
	SkipSensitive     bool
}

// ParsePromoteArgs parses promote sub-command flags from args.
func ParsePromoteArgs(args []string) (PromoteArgs, error) {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	var pa PromoteArgs
	fs.StringVar(&pa.OutputFile, "out", "", "write promoted env to file (default: stdout)")
	fs.BoolVar(&pa.OverwriteExisting, "overwrite", false, "overwrite existing keys in target")
	fs.BoolVar(&pa.SkipSensitive, "skip-sensitive", false, "skip keys matching common secret patterns")
	if err := fs.Parse(args); err != nil {
		return pa, err
	}
	rest := fs.Args()
	if len(rest) < 2 {
		return pa, fmt.Errorf("usage: promote [flags] <source.env> <target.env>")
	}
	pa.SourceFile = rest[0]
	pa.TargetFile = rest[1]
	return pa, nil
}

// RunPromote executes the promote workflow and writes the result.
func RunPromote(pa PromoteArgs, w io.Writer) error {
	src, err := parser.Parse(pa.SourceFile)
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}
	dst, err := parser.Parse(pa.TargetFile)
	if err != nil {
		return fmt.Errorf("target: %w", err)
	}

	masker := diff.NewMasker()
	sensitiveKeys := map[string]bool{}
	if pa.SkipSensitive {
		for k := range src {
			if masker.IsSensitive(k) {
				sensitiveKeys[k] = true
			}
		}
	}

	opts := PromoteOptions{
		OverwriteExisting: pa.OverwriteExisting,
		SkipSensitive:     pa.SkipSensitive,
		SensitiveKeys:     sensitiveKeys,
	}

	promoted, pr := Promote(src, dst, opts)
	fmt.Fprintln(w, PromoteSummary(pr))

	output := RenderEnv(promoted)
	if pa.OutputFile != "" {
		if err := os.WriteFile(pa.OutputFile, []byte(output), 0o644); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Fprintf(w, "wrote promoted env to %s\n", pa.OutputFile)
	} else {
		fmt.Fprint(w, output)
	}
	return nil
}
