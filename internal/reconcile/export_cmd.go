package reconcile

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// ExportArgs holds parsed CLI arguments for the export command.
type ExportArgs struct {
	BaseFile    string
	OtherFile   string
	Format      ExportFormat
	OnlyChanged bool
	NoMask      bool
	OutputFile  string
}

// ParseExportArgs parses CLI flags for the export subcommand.
func ParseExportArgs(args []string) (ExportArgs, error) {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	format := fs.String("format", "dotenv", "export format: dotenv, shell, docker")
	onlyChanged := fs.Bool("only-changed", false, "export only changed keys")
	noMask := fs.Bool("no-mask", false, "disable secret masking")
	output := fs.String("out", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return ExportArgs{}, err
	}
	positional := fs.Args()
	if len(positional) < 2 {
		return ExportArgs{}, fmt.Errorf("export: requires <base> <other> arguments")
	}
	return ExportArgs{
		BaseFile:    positional[0],
		OtherFile:   positional[1],
		Format:      ExportFormat(*format),
		OnlyChanged: *onlyChanged,
		NoMask:      *noMask,
		OutputFile:  *output,
	}, nil
}

// RunExport executes the export command with the given arguments.
func RunExport(args []string) error {
	opts, err := ParseExportArgs(args)
	if err != nil {
		return err
	}

	base, err := parser.Parse(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("export: parse base: %w", err)
	}
	other, err := parser.Parse(opts.OtherFile)
	if err != nil {
		return fmt.Errorf("export: parse other: %w", err)
	}

	result := diff.Compare(base, other)

	var masker *diff.Masker
	if !opts.NoMask {
		masker = diff.NewMasker()
	}

	exportOpts := ExportOptions{
		Format:      opts.Format,
		OnlyChanged: opts.OnlyChanged,
		Masker:      masker,
	}

	if opts.OutputFile != "" {
		return ExportToFile(opts.OutputFile, result.Entries, exportOpts)
	}
	return Export(os.Stdout, result.Entries, exportOpts)
}
