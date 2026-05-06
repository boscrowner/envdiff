package differ

// FilterByKind returns a new Result containing only entries matching the given kinds.
func (r *Result) FilterByKind(kinds ...DiffKind) *Result {
	kindSet := make(map[DiffKind]bool, len(kinds))
	for _, k := range kinds {
		kindSet[k] = true
	}

	var filtered []Entry
	for _, e := range r.Entries {
		if kindSet[e.Kind] {
			filtered = append(filtered, e)
		}
	}
	return &Result{Entries: filtered}
}

// Summary returns counts of each diff kind.
type Summary struct {
	Added     int
	Removed   int
	Modified  int
	Unchanged int
}

// Summarize returns a Summary of the diff result.
func (r *Result) Summarize() Summary {
	var s Summary
	for _, e := range r.Entries {
		switch e.Kind {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Modified:
			s.Modified++
		case Unchanged:
			s.Unchanged++
		}
	}
	return s
}
