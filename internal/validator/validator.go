package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Issue represents a validation problem found in an env file.
type Issue struct {
	Line int
	Key  string
	Kind string
	Msg  string
}

func (i Issue) String() string {
	return fmt.Sprintf("line %d [%s] %s: %s", i.Line, i.Kind, i.Key, i.Msg)
}

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks a parsed env file for common problems and returns a list
// of issues. An empty slice means the file is clean.
func Validate(entries []parser.Entry) []Issue {
	var issues []Issue
	seen := make(map[string]int)

	for idx, e := range entries {
		lineNum := idx + 1

		// Check key naming convention.
		if !validKeyRe.MatchString(e.Key) {
			issues = append(issues, Issue{
				Line: lineNum,
				Key:  e.Key,
				Kind: "invalid-key",
				Msg:  "key contains invalid characters (must match [A-Za-z_][A-Za-z0-9_]*)",
			})
		}

		// Check for duplicate keys.
		if prev, ok := seen[e.Key]; ok {
			issues = append(issues, Issue{
				Line: lineNum,
				Key:  e.Key,
				Kind: "duplicate-key",
				Msg:  fmt.Sprintf("duplicate of key first seen on line %d", prev),
			})
		} else {
			seen[e.Key] = lineNum
		}

		// Warn about empty values.
		if strings.TrimSpace(e.Value) == "" {
			issues = append(issues, Issue{
				Line: lineNum,
				Key:  e.Key,
				Kind: "empty-value",
				Msg:  "value is empty",
			})
		}
	}

	return issues
}
