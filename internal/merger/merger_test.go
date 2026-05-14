package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
	"github.com/user/envdiff/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMerge_NoConflict(t *testing.T) {
	base := entries("A", "1", "B", "2")
	override := entries("C", "3")
	res, err := merger.Merge(base, override, merger.PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.AddedKeys) != 1 || res.AddedKeys[0] != "C" {
		t.Errorf("expected C added, got %v", res.AddedKeys)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_PreferBase(t *testing.T) {
	base := entries("A", "base")
	override := entries("A", "override")
	res, err := merger.Merge(base, override, merger.PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Entries[0].Value != "base" {
		t.Errorf("expected base value preserved, got %q", res.Entries[0].Value)
	}
}

func TestMerge_PreferOverride(t *testing.T) {
	base := entries("A", "base")
	override := entries("A", "override")
	res, err := merger.Merge(base, override, merger.PreferOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "override" {
		t.Errorf("expected override value, got %q", res.Entries[0].Value)
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	base := entries("A", "1")
	override := entries("A", "2")
	_, err := merger.Merge(base, override, merger.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_IdenticalValues_NoConflict(t *testing.T) {
	base := entries("A", "same")
	override := entries("A", "same")
	res, err := merger.Merge(base, override, merger.ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("identical values should not be a conflict")
	}
}

func TestMerge_PreservesOrder(t *testing.T) {
	base := entries("Z", "last", "A", "first")
	override := entries("NEW", "added")
	res, err := merger.Merge(base, override, merger.PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Key != "Z" || res.Entries[1].Key != "A" || res.Entries[2].Key != "NEW" {
		t.Errorf("order not preserved: %v", res.Entries)
	}
}
