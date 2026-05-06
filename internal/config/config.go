package config

import (
	"errors"
	"strings"
)

// OutputFormat defines the output format for diff results.
type OutputFormat string

const (
	FormatText  OutputFormat = "text"
	FormatTable OutputFormat = "table"
	FormatJSON  OutputFormat = "json"
)

// Config holds the runtime configuration for envdiff.
type Config struct {
	BaseFile    string
	OtherFile   string
	Format      OutputFormat
	MaskSecrets bool
	ColorOutput bool
	IndentJSON  bool
	SecretKeys  []string
}

// Validate checks that the config is valid and returns an error if not.
func (c *Config) Validate() error {
	if c.BaseFile == "" {
		return errors.New("base file path is required")
	}
	if c.OtherFile == "" {
		return errors.New("other file path is required")
	}
	if !validFormat(c.Format) {
		return errors.New("invalid output format: must be one of text, table, json")
	}
	return nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Format:      FormatText,
		MaskSecrets: true,
		ColorOutput: true,
		IndentJSON:  false,
	}
}

// ParseFormat converts a string to an OutputFormat, normalising case.
func ParseFormat(s string) (OutputFormat, error) {
	f := OutputFormat(strings.ToLower(strings.TrimSpace(s)))
	if !validFormat(f) {
		return "", errors.New("unknown format: " + s)
	}
	return f, nil
}

func validFormat(f OutputFormat) bool {
	switch f {
	case FormatText, FormatTable, FormatJSON:
		return true
	}
	return false
}
