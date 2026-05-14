package merger

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how conflicts are resolved when merging two env files.
type Strategy int

const (
	// PreferBase keeps the base file's value on conflict.
	PreferBase Strategy = iota
	// PreferOverride uses the override file's value on conflict.
	PreferOverride
	// ErrorOnConflict returns an error if any key exists in both files with different values.
	ErrorOnConflict
)

// Result holds the merged entries and metadata about the merge operation.
type Result struct {
	Entries    []parser.Entry
	Conflicts  []Conflict
	AddedKeys  []string
}

// Conflict describes a key that existed in both files with differing values.
type Conflict struct {
	Key       string
	BaseValue string
	OverrideValue string
}

// Merge combines base and override env entries according to the given strategy.
// Keys present only in override are appended. Comments and blank lines from
// base are preserved; override comments are discarded.
func Merge(base, override []parser.Entry, strategy Strategy) (Result, error) {
	index := make(map[string]int, len(base))
	result := make([]parser.Entry, len(base))
	copy(result, base)

	for i, e := range result {
		if e.Key != "" {
			index[e.Key] = i
		}
	}

	var conflicts []Conflict
	var added []string

	for _, oe := range override {
		if oe.Key == "" {
			continue
		}
		if idx, exists := index[oe.Key]; exists {
			if result[idx].Value == oe.Value {
				continue
			}
			conflicts = append(conflicts, Conflict{
				Key:           oe.Key,
				BaseValue:     result[idx].Value,
				OverrideValue: oe.Value,
			})
			switch strategy {
			case ErrorOnConflict:
				return Result{}, fmt.Errorf("merger: conflict on key %q", oe.Key)
			case PreferOverride:
				result[idx].Value = oe.Value
			// PreferBase: do nothing
			}
		} else {
			result = append(result, oe)
			added = append(added, oe.Key)
		}
	}

	return Result{Entries: result, Conflicts: conflicts, AddedKeys: added}, nil
}
