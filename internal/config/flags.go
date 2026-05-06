package config

import "flag"

// FlagSet binds CLI flags to a Config and returns the FlagSet for parsing.
func FlagSet(cfg *Config) *flag.FlagSet {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	fs.StringVar(&cfg.BaseFile, "base", "", "path to the base .env file (required)")
	fs.StringVar(&cfg.OtherFile, "other", "", "path to the other .env file to compare (required)")

	var format string
	fs.StringVar(&format, "format", string(FormatText), "output format: text, table, json")

	fs.BoolVar(&cfg.MaskSecrets, "mask", true, "mask sensitive values in output")
	fs.BoolVar(&cfg.ColorOutput, "color", true, "enable colored output")
	fs.BoolVar(&cfg.IndentJSON, "indent", false, "pretty-print JSON output")

	// Post-parse hook stored via a wrapper; callers must call ApplyFormat after Parse.
	_ = format

	return fs
}

// FromArgs parses os.Args-style arguments into a Config.
// It returns the populated Config and any parse error.
func FromArgs(args []string) (*Config, error) {
	cfg := Default()
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	var format string
	fs.StringVar(&cfg.BaseFile, "base", "", "path to the base .env file (required)")
	fs.StringVar(&cfg.OtherFile, "other", "", "path to the other .env file to compare (required)")
	fs.StringVar(&format, "format", string(FormatText), "output format: text, table, json")
	fs.BoolVar(&cfg.MaskSecrets, "mask", true, "mask sensitive values in output")
	fs.BoolVar(&cfg.ColorOutput, "color", true, "enable colored output")
	fs.BoolVar(&cfg.IndentJSON, "indent", false, "pretty-print JSON output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	f, err := ParseFormat(format)
	if err != nil {
		return nil, err
	}
	cfg.Format = f

	return cfg, nil
}
