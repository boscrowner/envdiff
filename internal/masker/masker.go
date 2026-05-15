// Package masker provides utilities for masking sensitive values
// in .env entries before display or logging.
package masker

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how values are masked.
type Strategy int

const (
	// StrategyRedact replaces the entire value with "[REDACTED]".
	StrategyRedact Strategy = iota
	// StrategyPartial reveals the first and last character, masking the middle.
	StrategyPartial
)

// defaultSensitivePatterns are substrings checked (case-insensitive) against
// key names to determine whether a value should be masked.
var defaultSensitivePatterns = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"private",
	"credential",
	"auth",
}

// IsSensitive reports whether the given key name looks sensitive based on
// the default pattern list.
func IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range defaultSensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// MaskValue masks a single value string using the given strategy.
// If the value is empty it is returned unchanged.
func MaskValue(value string, strategy Strategy) string {
	if value == "" {
		return value
	}
	switch strategy {
	case StrategyPartial:
		if len(value) <= 2 {
			return strings.Repeat("*", len(value))
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	default:
		return "[REDACTED]"
	}
}

// MaskEntries returns a new slice of entries where sensitive values have been
// masked according to the provided strategy. Non-sensitive entries are copied
// unchanged.
func MaskEntries(entries []parser.Entry, strategy Strategy) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if IsSensitive(e.Key) {
			out[i] = parser.Entry{Key: e.Key, Value: MaskValue(e.Value, strategy)}
		} else {
			out[i] = e
		}
	}
	return out
}
