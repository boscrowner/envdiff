// Package sorter provides utilities for sorting and grouping .env file entries
// by key prefix, alphabetical order, or custom ordering strategies.
package sorter

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Strategy defines how entries should be sorted.
type Strategy string

const (
	// Alphabetical sorts all keys in A-Z order.
	Alphabetical Strategy = "alpha"
	// ByPrefix groups keys by their prefix (e.g. DB_, AWS_) then sorts within groups.
	ByPrefix Strategy = "prefix"
)

// Result holds the sorted entries and any grouping metadata.
type Result struct {
	Entries []parser.Entry
	Groups  map[string][]parser.Entry // only populated for ByPrefix strategy
}

// Sort returns a Result with entries ordered according to the given strategy.
func Sort(entries []parser.Entry, strategy Strategy) Result {
	copied := make([]parser.Entry, len(entries))
	copy(copied, entries)

	switch strategy {
	case ByPrefix:
		return sortByPrefix(copied)
	default:
		return sortAlphabetical(copied)
	}
}

func sortAlphabetical(entries []parser.Entry) Result {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return Result{Entries: entries}
}

func sortByPrefix(entries []parser.Entry) Result {
	groups := make(map[string][]parser.Entry)

	for _, e := range entries {
		prefix := extractPrefix(e.Key)
		groups[prefix] = append(groups[prefix], e)
	}

	// Sort within each group
	for prefix := range groups {
		sort.Slice(groups[prefix], func(i, j int) bool {
			return groups[prefix][i].Key < groups[prefix][j].Key
		})
	}

	// Collect sorted prefix names
	prefixes := make([]string, 0, len(groups))
	for p := range groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	// Flatten into ordered slice
	result := make([]parser.Entry, 0, len(entries))
	for _, p := range prefixes {
		result = append(result, groups[p]...)
	}

	return Result{Entries: result, Groups: groups}
}

// extractPrefix returns the portion of a key before the first underscore,
// or the full key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
