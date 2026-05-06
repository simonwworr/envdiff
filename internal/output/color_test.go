package output

import (
	"strings"
	"testing"
)

func TestColorizer_Disabled(t *testing.T) {
	c := NewColorizer(false)

	tests := []struct {
		name string
		fn   func(string) string
		input string
	}{
		{"Added", c.Added, "hello"},
		{"Removed", c.Removed, "world"},
		{"Changed", c.Changed, "foo"},
		{"Unchanged", c.Unchanged, "bar"},
		{"Key", c.Key, "MY_KEY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			if got != tt.input {
				t.Errorf("expected %q, got %q", tt.input, got)
			}
		})
	}
}

func TestColorizer_Enabled_ContainsANSI(t *testing.T) {
	c := NewColorizer(true)

	tests := []struct {
		name  string
		fn    func(string) string
		color string
	}{
		{"Added", c.Added, colorGreen},
		{"Removed", c.Removed, colorRed},
		{"Changed", c.Changed, colorYellow},
		{"Unchanged", c.Unchanged, colorGray},
		{"Key", c.Key, colorCyan},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn("test")
			if !strings.Contains(got, tt.color) {
				t.Errorf("expected ANSI code %q in output %q", tt.color, got)
			}
			if !strings.Contains(got, colorReset) {
				t.Errorf("expected reset code in output %q", got)
			}
			if !strings.Contains(got, "test") {
				t.Errorf("expected original text in output %q", got)
			}
		})
	}
}

func TestColorizer_Enabled_PreservesText(t *testing.T) {
	c := NewColorizer(true)
	input := "VALUE=123"
	got := c.Added(input)
	if !strings.Contains(got, input) {
		t.Errorf("original text %q not found in %q", input, got)
	}
}
