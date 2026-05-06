package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestFilterByKind(t *testing.T) {
	base := map[string]string{
		"REMOVED_KEY": "val",
		"SAME_KEY":    "same",
		"MOD_KEY":     "old",
	}
	other := map[string]string{
		"ADDED_KEY": "new",
		"SAME_KEY":  "same",
		"MOD_KEY":   "new",
	}

	result := differ.Diff(base, other)

	added := result.FilterByKind(differ.Added)
	if len(added.Entries) != 1 || added.Entries[0].Key != "ADDED_KEY" {
		t.Errorf("expected 1 added entry, got %+v", added.Entries)
	}

	removed := result.FilterByKind(differ.Removed)
	if len(removed.Entries) != 1 || removed.Entries[0].Key != "REMOVED_KEY" {
		t.Errorf("expected 1 removed entry, got %+v", removed.Entries)
	}

	changes := result.FilterByKind(differ.Added, differ.Removed, differ.Modified)
	if len(changes.Entries) != 3 {
		t.Errorf("expected 3 changed entries, got %d", len(changes.Entries))
	}
}

func TestSummarize(t *testing.T) {
	base := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	other := map[string]string{
		"A": "1",
		"B": "changed",
		"D": "4",
	}

	result := differ.Diff(base, other)
	s := result.Summarize()

	if s.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", s.Unchanged)
	}
	if s.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", s.Modified)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Added != 1 {
		t.Errorf("expected 1 added, got %d", s.Added)
	}
}
