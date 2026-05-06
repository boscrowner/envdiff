// Package differ compares two parsed .env maps and produces a structured diff.
package differ

// DiffKind represents the type of change between two env files.
type DiffKind string

const (
	Added    DiffKind = "added"
	Removed  DiffKind = "removed"
	Modified DiffKind = "modified"
	Unchanged DiffKind = "unchanged"
)

// Entry represents a single diff entry for a key.
type Entry struct {
	Key      string
	Kind     DiffKind
	BaseVal  string
	OtherVal string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if there are any added, removed, or modified entries.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Diff compares base and other env maps and returns a Result.
// base is typically the reference environment (e.g. .env.example).
// other is the environment being checked (e.g. .env).
func Diff(base, other map[string]string) *Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, baseVal := range base {
		seen[k] = true
		otherVal, exists := other[k]
		switch {
		case !exists:
			entries = append(entries, Entry{Key: k, Kind: Removed, BaseVal: baseVal})
		case baseVal != otherVal:
			entries = append(entries, Entry{Key: k, Kind: Modified, BaseVal: baseVal, OtherVal: otherVal})
		default:
			entries = append(entries, Entry{Key: k, Kind: Unchanged, BaseVal: baseVal, OtherVal: otherVal})
		}
	}

	for k, otherVal := range other {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Kind: Added, OtherVal: otherVal})
		}
	}

	sortEntries(entries)
	return &Result{Entries: entries}
}

// sortEntries sorts diff entries alphabetically by key for deterministic output.
func sortEntries(entries []Entry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Key < entries[j-1].Key; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
