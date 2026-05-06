package reconciler

import (
	"fmt"
	"strings"

	"github.com/yourorg/envdiff/internal/differ"
	"github.com/yourorg/envdiff/internal/parser"
)

// Action represents a reconciliation action to apply.
type Action struct {
	Kind    string // "add", "remove", "update"
	Key     string
	Value   string
	Comment string
}

// Plan returns a list of Actions needed to make source match target.
func Plan(source, target []parser.Entry) []Action {
	diffs := differ.Diff(source, target)
	actions := make([]Action, 0, len(diffs))

	for _, d := range diffs {
		switch d.Kind {
		case differ.Added:
			actions = append(actions, Action{
				Kind:    "add",
				Key:     d.Key,
				Value:   d.TargetValue,
				Comment: fmt.Sprintf("add %s from target", d.Key),
			})
		case differ.Removed:
			actions = append(actions, Action{
				Kind:    "remove",
				Key:     d.Key,
				Comment: fmt.Sprintf("remove %s not present in target", d.Key),
			})
		case differ.Modified:
			actions = append(actions, Action{
				Kind:    "update",
				Key:     d.Key,
				Value:   d.TargetValue,
				Comment: fmt.Sprintf("update %s: %q -> %q", d.Key, d.SourceValue, d.TargetValue),
			})
		}
	}

	return actions
}

// Apply merges the given Actions into the source entries and returns the result.
func Apply(source []parser.Entry, actions []Action) []parser.Entry {
	entryMap := make(map[string]parser.Entry, len(source))
	order := make([]string, 0, len(source))

	for _, e := range source {
		entryMap[e.Key] = e
		order = append(order, e.Key)
	}

	for _, a := range actions {
		switch a.Kind {
		case "add":
			if _, exists := entryMap[a.Key]; !exists {
				order = append(order, a.Key)
			}
			entryMap[a.Key] = parser.Entry{Key: a.Key, Value: a.Value}
		case "update":
			if e, exists := entryMap[a.Key]; exists {
				e.Value = a.Value
				entryMap[a.Key] = e
			}
		case "remove":
			delete(entryMap, a.Key)
		}
	}

	result := make([]parser.Entry, 0, len(order))
	for _, k := range order {
		if e, ok := entryMap[k]; ok {
			result = append(result, e)
		}
	}

	return result
}

// Render serialises entries back to .env file content.
func Render(entries []parser.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
	}
	return sb.String()
}
