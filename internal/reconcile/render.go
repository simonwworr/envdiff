package reconcile

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// RenderEnv writes the env map as a sorted .env file to w.
func RenderEnv(w io.Writer, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		line := formatLine(k, v)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("render: write key %q: %w", k, err)
		}
	}
	return nil
}

// formatLine produces a properly quoted key=value line.
func formatLine(key, value string) string {
	if needsQuoting(value) {
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		return fmt.Sprintf(`%s="%s"`, key, escaped)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// needsQuoting returns true when a value contains whitespace, quotes, or
// special shell characters that require double-quoting.
func needsQuoting(v string) bool {
	specials := " \t\n#$'\\"
	return strings.ContainsAny(v, specials)
}
