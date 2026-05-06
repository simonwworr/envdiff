package output

import "fmt"

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// ColorMode controls whether ANSI color codes are emitted.
type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

// Colorizer applies ANSI color codes to strings based on the configured mode.
type Colorizer struct {
	enabled bool
}

// NewColorizer creates a Colorizer. When enabled is true, output will include
// ANSI escape sequences.
func NewColorizer(enabled bool) *Colorizer {
	return &Colorizer{enabled: enabled}
}

// Added returns the string styled for an added value (green).
func (c *Colorizer) Added(s string) string {
	return c.wrap(colorGreen, s)
}

// Removed returns the string styled for a removed value (red).
func (c *Colorizer) Removed(s string) string {
	return c.wrap(colorRed, s)
}

// Changed returns the string styled for a changed value (yellow).
func (c *Colorizer) Changed(s string) string {
	return c.wrap(colorYellow, s)
}

// Unchanged returns the string styled for an unchanged value (gray).
func (c *Colorizer) Unchanged(s string) string {
	return c.wrap(colorGray, s)
}

// Key returns the string styled for a key name (cyan).
func (c *Colorizer) Key(s string) string {
	return c.wrap(colorCyan, s)
}

func (c *Colorizer) wrap(code, s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", code, s, colorReset)
}
