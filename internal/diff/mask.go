package diff

import (
	"strings"
)

const maskedValue = "***"

// defaultSecretPatterns are substrings that, when found in a key name
// (case-insensitive), cause the value to be masked.
var defaultSecretPatterns = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"auth",
	"private",
	"credential",
}

// Masker masks sensitive values in diff entries.
type Masker struct {
	patterns []string
}

// NewMasker creates a Masker using the default secret patterns.
func NewMasker() *Masker {
	return &Masker{patterns: defaultSecretPatterns}
}

// NewMaskerWithPatterns creates a Masker with custom key patterns.
func NewMaskerWithPatterns(patterns []string) *Masker {
	return &Masker{patterns: patterns}
}

// IsSensitive returns true if the key matches any secret pattern.
func (m *Masker) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range m.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// MaskResult returns a copy of the Result with sensitive values replaced.
func (m *Masker) MaskResult(r *Result) *Result {
	masked := make([]Entry, len(r.Entries))
	for i, e := range r.Entries {
		if m.IsSensitive(e.Key) {
			e.BaseVal = maskIfNonEmpty(e.BaseVal)
			e.OtherVal = maskIfNonEmpty(e.OtherVal)
		}
		masked[i] = e
	}
	return &Result{Entries: masked}
}

func maskIfNonEmpty(v string) string {
	if v == "" {
		return ""
	}
	return maskedValue
}
