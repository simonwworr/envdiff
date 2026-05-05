package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]*Entry
}

// Parse reads and parses a .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		line := strings.TrimSpace(raw)

		// Skip blank lines
		if line == "" {
			continue
		}

		// Comment-only lines
		if strings.HasPrefix(line, "#") {
			continue
		}

		key, value, comment, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		entry := Entry{
			Key:     key,
			Value:   value,
			Comment: comment,
			Line:    lineNum,
		}
		env.Entries = append(env.Entries, entry)
		env.Index[key] = &env.Entries[len(env.Entries)-1]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return env, nil
}

// parseLine splits a raw env line into key, value, and inline comment.
func parseLine(line string) (key, value, comment string, err error) {
	// Strip export prefix
	line = strings.TrimPrefix(line, "export ")

	eqIdx := strings.IndexByte(line, '=')
	if eqIdx < 0 {
		return "", "", "", fmt.Errorf("missing '=' in line: %q", line)
	}

	key = strings.TrimSpace(line[:eqIdx])
	rest := line[eqIdx+1:]

	// Handle quoted values
	if len(rest) > 0 && (rest[0] == '"' || rest[0] == '\'') {
		quote := rest[0]
		close := strings.IndexByte(rest[1:], quote)
		if close < 0 {
			return "", "", "", fmt.Errorf("unclosed quote in value for key %q", key)
		}
		value = rest[1 : close+1]
		return key, value, "", nil
	}

	// Unquoted: split on first ' #'
	if idx := strings.Index(rest, " #"); idx >= 0 {
		value = strings.TrimSpace(rest[:idx])
		comment = strings.TrimSpace(rest[idx+2:])
	} else {
		value = strings.TrimSpace(rest)
	}

	return key, value, comment, nil
}
