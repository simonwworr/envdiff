package config_test

import (
	"testing"

	"github.com/user/envdiff/internal/config"
)

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.Format != config.FormatText {
		t.Errorf("expected default format text, got %s", cfg.Format)
	}
	if !cfg.MaskSecrets {
		t.Error("expected MaskSecrets to be true by default")
	}
	if !cfg.ColorOutput {
		t.Error("expected ColorOutput to be true by default")
	}
}

func TestValidate_MissingBase(t *testing.T) {
	cfg := config.Default()
	cfg.OtherFile = "other.env"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing base file")
	}
}

func TestValidate_MissingOther(t *testing.T) {
	cfg := config.Default()
	cfg.BaseFile = "base.env"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing other file")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	cfg := config.Default()
	cfg.BaseFile = "base.env"
	cfg.OtherFile = "other.env"
	cfg.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := config.Default()
	cfg.BaseFile = "base.env"
	cfg.OtherFile = "other.env"
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected config.OutputFormat
	}{
		{"text", config.FormatText},
		{"TABLE", config.FormatTable},
		{" json ", config.FormatJSON},
	}
	for _, tc := range cases {
		f, err := config.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc.input, err)
		}
		if f != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, f)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := config.ParseFormat("yaml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
