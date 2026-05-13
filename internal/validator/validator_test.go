package validator_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/validator"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestValidate_Clean(t *testing.T) {
	issues := validator.Validate(entries("APP_NAME", "myapp", "PORT", "8080"))
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_InvalidKey(t *testing.T) {
	issues := validator.Validate(entries("1INVALID", "value"))
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Kind != "invalid-key" {
		t.Errorf("expected kind invalid-key, got %s", issues[0].Kind)
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	issues := validator.Validate(entries("FOO", "bar", "FOO", "baz"))
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Kind != "duplicate-key" {
		t.Errorf("expected kind duplicate-key, got %s", issues[0].Kind)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	issues := validator.Validate(entries("SECRET", ""))
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Kind != "empty-value" {
		t.Errorf("expected kind empty-value, got %s", issues[0].Kind)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	// duplicate key + first occurrence has empty value
	issue := validator.Validate(entries("OK", "good", "BAD KEY", "v", "OK", "again"))
	if len(issue) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issue), issue)
	}
}

func TestIssue_String(t *testing.T) {
	i := validator.Issue{Line: 3, Key: "FOO", Kind: "empty-value", Msg: "value is empty"}
	s := i.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
}
