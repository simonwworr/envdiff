package reconcile

import (
	"strings"
	"testing"
)

func TestRenderEnv_BasicOutput(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	result := RenderEnv(env)
	if !strings.Contains(result, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", result)
	}
	if !strings.Contains(result, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", result)
	}
}

func TestRenderEnv_QuotesValuesWithSpaces(t *testing.T) {
	env := map[string]string{
		"GREETING": "hello world",
	}
	result := RenderEnv(env)
	if !strings.Contains(result, `GREETING="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", result)
	}
}

func TestRenderEnv_QuotesValuesWithSpecialChars(t *testing.T) {
	cases := map[string]string{
		"A": "val#comment",
		"B": "val=extra",
		"C": "val$dollar",
	}
	result := RenderEnv(cases)
	for key, val := range cases {
		expected := key + `="` + val + `"`
		if !strings.Contains(result, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, result)
		}
	}
}

func TestRenderEnv_EmptyValue(t *testing.T) {
	env := map[string]string{
		"EMPTY": "",
	}
	result := RenderEnv(env)
	if !strings.Contains(result, "EMPTY=") {
		t.Errorf("expected EMPTY= in output, got:\n%s", result)
	}
}

func TestRenderEnv_SortedKeys(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "1",
		"APPLE": "2",
		"MANGO": "3",
	}
	result := RenderEnv(env)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APPLE") {
		t.Errorf("expected first line to start with APPLE, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "MANGO") {
		t.Errorf("expected second line to start with MANGO, got %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected third line to start with ZEBRA, got %s", lines[2])
	}
}

func TestNeedsQuoting(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"simple", false},
		{"has space", true},
		{"has#hash", true},
		{"has=equals", true},
		{"has$dollar", true},
		{"", false},
	}
	for _, tc := range cases {
		got := needsQuoting(tc.input)
		if got != tc.expected {
			t.Errorf("needsQuoting(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}
