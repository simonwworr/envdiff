package config_test

import (
	"testing"

	"github.com/user/envdiff/internal/config"
)

func TestFromArgs_Defaults(t *testing.T) {
	cfg, err := config.FromArgs([]string{"-base", "a.env", "-other", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != config.FormatText {
		t.Errorf("expected default format text, got %s", cfg.Format)
	}
	if !cfg.MaskSecrets {
		t.Error("expected MaskSecrets true")
	}
	if !cfg.ColorOutput {
		t.Error("expected ColorOutput true")
	}
}

func TestFromArgs_FormatJSON(t *testing.T) {
	cfg, err := config.FromArgs([]string{"-base", "a.env", "-other", "b.env", "-format", "json", "-indent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != config.FormatJSON {
		t.Errorf("expected json format, got %s", cfg.Format)
	}
	if !cfg.IndentJSON {
		t.Error("expected IndentJSON true")
	}
}

func TestFromArgs_FormatTable(t *testing.T) {
	cfg, err := config.FromArgs([]string{"-base", "a.env", "-other", "b.env", "-format", "table"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != config.FormatTable {
		t.Errorf("expected table format, got %s", cfg.Format)
	}
}

func TestFromArgs_InvalidFormat(t *testing.T) {
	_, err := config.FromArgs([]string{"-base", "a.env", "-other", "b.env", "-format", "csv"})
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestFromArgs_NoMask(t *testing.T) {
	cfg, err := config.FromArgs([]string{"-base", "a.env", "-other", "b.env", "-mask=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaskSecrets {
		t.Error("expected MaskSecrets false")
	}
}
