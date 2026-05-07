package reconcile

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeEntry(key string, status diff.Status, base, other string) diff.Entry {
	return diff.Entry{
		Key:        key,
		Status:     status,
		BaseValue:  base,
		OtherValue: other,
	}
}

func makeResultForPatch(entries []diff.Entry) diff.Result {
	return diff.Result{Entries: entries}
}

func TestGeneratePatch_ChangedOnly(t *testing.T) {
	result := makeResultForPatch([]diff.Entry{
		makeEntry("HOST", diff.Changed, "localhost", "prod.example.com"),
		makeEntry("PORT", diff.Unchanged, "5432", "5432"),
		makeEntry("NEW_KEY", diff.Added, "", "newval"),
	})

	patch := GeneratePatch(result, PatchChangedOnly)

	if patch.Empty() {
		t.Fatal("expected non-empty patch")
	}
	if len(patch.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(patch.Lines))
	}
	if !strings.Contains(patch.Lines[0], "HOST") {
		t.Errorf("expected HOST in patch line, got %q", patch.Lines[0])
	}
}

func TestGeneratePatch_AddMissing(t *testing.T) {
	result := makeResultForPatch([]diff.Entry{
		makeEntry("HOST", diff.Changed, "localhost", "prod.example.com"),
		makeEntry("NEW_KEY", diff.Added, "", "newval"),
	})

	patch := GeneratePatch(result, PatchAddMissing)

	if len(patch.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(patch.Lines))
	}
}

func TestGeneratePatch_Empty(t *testing.T) {
	result := makeResultForPatch([]diff.Entry{
		makeEntry("HOST", diff.Unchanged, "localhost", "localhost"),
	})

	patch := GeneratePatch(result, PatchChangedOnly)

	if !patch.Empty() {
		t.Error("expected empty patch for unchanged result")
	}
	if patch.String() != "" {
		t.Errorf("expected empty string, got %q", patch.String())
	}
}

func TestPatch_String_ContainsHeader(t *testing.T) {
	result := makeResultForPatch([]diff.Entry{
		makeEntry("DB", diff.Changed, "old", "new"),
	})

	patch := GeneratePatch(result, PatchAddMissing)
	out := patch.String()

	if !strings.HasPrefix(out, "# envdiff patch") {
		t.Errorf("expected header comment, got %q", out)
	}
	if !strings.Contains(out, "add-missing") {
		t.Errorf("expected mode label in header, got %q", out)
	}
}

func TestPatch_ModeLabel_ChangedOnly(t *testing.T) {
	p := Patch{Mode: PatchChangedOnly, Lines: []string{"X=1"}}
	if !strings.Contains(p.String(), "changed-only") {
		t.Errorf("expected changed-only label, got %q", p.String())
	}
}
