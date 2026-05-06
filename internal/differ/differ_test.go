package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestDiff_AddedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	other := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := differ.Diff(base, other)

	if !result.HasChanges() {
		t.Fatal("expected changes, got none")
	}
	found := findEntry(result, "NEW_KEY")
	if found == nil || found.Kind != differ.Added {
		t.Errorf("expected NEW_KEY to be Added, got %+v", found)
	}
}

func TestDiff_RemovedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar", "MISSING": "gone"}
	other := map[string]string{"FOO": "bar"}

	result := differ.Diff(base, other)

	found := findEntry(result, "MISSING")
	if found == nil || found.Kind != differ.Removed {
		t.Errorf("expected MISSING to be Removed, got %+v", found)
	}
}

func TestDiff_ModifiedKeys(t *testing.T) {
	base := map[string]string{"DB_URL": "old"}
	other := map[string]string{"DB_URL": "new"}

	result := differ.Diff(base, other)

	found := findEntry(result, "DB_URL")
	if found == nil || found.Kind != differ.Modified {
		t.Errorf("expected DB_URL to be Modified, got %+v", found)
	}
	if found.BaseVal != "old" || found.OtherVal != "new" {
		t.Errorf("unexpected values: base=%q other=%q", found.BaseVal, found.OtherVal)
	}
}

func TestDiff_UnchangedKeys(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	other := map[string]string{"PORT": "8080"}

	result := differ.Diff(base, other)

	if result.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	base := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	other := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}

	result := differ.Diff(base, other)

	keys := make([]string, len(result.Entries))
	for i, e := range result.Entries {
		keys[i] = e.Key
	}
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("expected key[%d]=%q, got %q", i, k, keys[i])
		}
	}
}

func findEntry(r *differ.Result, key string) *differ.Entry {
	for i := range r.Entries {
		if r.Entries[i].Key == key {
			return &r.Entries[i]
		}
	}
	return nil
}
