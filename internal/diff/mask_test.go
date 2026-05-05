package diff

import "testing"

func TestMasker_IsSensitive(t *testing.T) {
	m := NewMasker()

	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"SECRET_KEY", true},
		{"DB_HOST", false},
		{"PORT", false},
		{"APP_NAME", false},
		{"PRIVATE_KEY", true},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := m.IsSensitive(tc.key)
			if got != tc.expected {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestMasker_MaskResult(t *testing.T) {
	m := NewMasker()

	r := &Result{
		Entries: []Entry{
			{Key: "DB_HOST", Status: StatusUnchanged, BaseVal: "localhost", OtherVal: "localhost"},
			{Key: "DB_PASSWORD", Status: StatusChanged, BaseVal: "hunter2", OtherVal: "s3cr3t"},
			{Key: "API_KEY", Status: StatusAdded, BaseVal: "", OtherVal: "abc123"},
		},
	}

	masked := m.MaskResult(r)

	if masked.Entries[0].BaseVal != "localhost" {
		t.Error("non-sensitive value should not be masked")
	}
	if masked.Entries[1].BaseVal != maskedValue {
		t.Errorf("expected masked BaseVal, got %s", masked.Entries[1].BaseVal)
	}
	if masked.Entries[1].OtherVal != maskedValue {
		t.Errorf("expected masked OtherVal, got %s", masked.Entries[1].OtherVal)
	}
	if masked.Entries[2].BaseVal != "" {
		t.Error("empty BaseVal should remain empty after masking")
	}
	if masked.Entries[2].OtherVal != maskedValue {
		t.Errorf("expected masked OtherVal for added key, got %s", masked.Entries[2].OtherVal)
	}
}

func TestMasker_CustomPatterns(t *testing.T) {
	m := NewMaskerWithPatterns([]string{"pin", "ssn"})

	if !m.IsSensitive("USER_PIN") {
		t.Error("USER_PIN should be sensitive with custom pattern")
	}
	if m.IsSensitive("DB_PASSWORD") {
		t.Error("DB_PASSWORD should NOT be sensitive with custom patterns")
	}
}
