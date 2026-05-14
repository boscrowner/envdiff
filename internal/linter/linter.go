package linter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Severity represents the level of a lint issue.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue represents a single lint finding for an env entry.
type Issue struct {
	Key      string
	Line     int
	Severity Severity
	Message  string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s [line %d] %s: %s", i.Severity, i.Line, i.Key, i.Message)
}

// Lint runs all lint rules against the parsed entries from the given file path
// and returns any issues found.
func Lint(path string) ([]Issue, error) {
	entries, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("linter: parse %q: %w", path, err)
	}

	var issues []Issue
	for _, e := range entries {
		issues = append(issues, runRules(e)...)
	}
	return issues, nil
}

func runRules(e parser.Entry) []Issue {
	var issues []Issue

	if strings.ToUpper(e.Key) != e.Key {
		issues = append(issues, Issue{
			Key:      e.Key,
			Line:     e.Line,
			Severity: SeverityWarning,
			Message:  "key should be uppercase",
		})
	}

	if strings.Contains(e.Key, " ") {
		issues = append(issues, Issue{
			Key:      e.Key,
			Line:     e.Line,
			Severity: SeverityError,
			Message:  "key must not contain spaces",
		})
	}

	if strings.HasPrefix(e.Value, " ") || strings.HasSuffix(e.Value, " ") {
		issues = append(issues, Issue{
			Key:      e.Key,
			Line:     e.Line,
			Severity: SeverityWarning,
			Message:  "value has leading or trailing whitespace",
		})
	}

	if strings.Contains(e.Value, "password") || strings.Contains(e.Value, "secret") {
		issues = append(issues, Issue{
			Key:      e.Key,
			Line:     e.Line,
			Severity: SeverityWarning,
			Message:  "value appears to contain a plaintext secret",
		})
	}

	return issues
}
