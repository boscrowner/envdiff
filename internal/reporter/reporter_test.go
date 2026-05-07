package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reporter"
)

func TestTextReport_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.TextReport(&buf, []differ.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected 'No differences found', got: %s", buf.String())
	}
}

func TestTextReport_Summary(t *testing.T) {
	diffs := []differ.Diff{
		{Key: "NEW_KEY", Kind: differ.KindAdded, NewValue: "val1"},
		{Key: "OLD_KEY", Kind: differ.KindRemoved, OldValue: "val2"},
		{Key: "MOD_KEY", Kind: differ.KindModified, OldValue: "old", NewValue: "new"},
		{Key: "SAME_KEY", Kind: differ.KindUnchanged, OldValue: "same"},
	}

	var buf bytes.Buffer
	err := reporter.TextReport(&buf, diffs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "1 added") {
		t.Errorf("expected '1 added' in output, got: %s", out)
	}
	if !strings.Contains(out, "1 removed") {
		t.Errorf("expected '1 removed' in output, got: %s", out)
	}
	if !strings.Contains(out, "1 modified") {
		t.Errorf("expected '1 modified' in output, got: %s", out)
	}
}

func TestTextReport_DiffLineFormats(t *testing.T) {
	diffs := []differ.Diff{
		{Key: "A", Kind: differ.KindAdded, NewValue: "1"},
		{Key: "B", Kind: differ.KindRemoved, OldValue: "2"},
		{Key: "C", Kind: differ.KindModified, OldValue: "x", NewValue: "y"},
		{Key: "D", Kind: differ.KindUnchanged, OldValue: "z"},
	}

	var buf bytes.Buffer
	if err := reporter.TextReport(&buf, diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"+ A", "- B", "~ C", "  D"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output\n%s", want, out)
		}
	}
}
