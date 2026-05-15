package masker_test

import (
	"testing"

	"github.com/user/envdiff/internal/masker"
	"github.com/user/envdiff/internal/parser"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	sensitiveKeys := []string{
		"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN", "PRIVATE_KEY",
		"STRIPE_API_KEY", "USER_CREDENTIAL", "PASSWD",
	}
	for _, k := range sensitiveKeys {
		if !masker.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_IgnoresSafeKeys(t *testing.T) {
	safeKeys := []string{"APP_ENV", "PORT", "LOG_LEVEL", "DATABASE_HOST"}
	for _, k := range safeKeys {
		if masker.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestMaskValue_Redact(t *testing.T) {
	got := masker.MaskValue("supersecret", masker.StrategyRedact)
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMaskValue_Partial(t *testing.T) {
	got := masker.MaskValue("supersecret", masker.StrategyPartial)
	if got != "s*********t" {
		t.Errorf("unexpected partial mask: %q", got)
	}
}

func TestMaskValue_PartialShort(t *testing.T) {
	got := masker.MaskValue("ab", masker.StrategyPartial)
	if got != "**" {
		t.Errorf("expected **, got %q", got)
	}
}

func TestMaskValue_Empty(t *testing.T) {
	got := masker.MaskValue("", masker.StrategyRedact)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMaskEntries_MasksSensitiveOnly(t *testing.T) {
	input := []parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
		{Key: "PORT", Value: "8080"},
		{Key: "API_SECRET", Value: "abc123"},
	}

	result := masker.MaskEntries(input, masker.StrategyRedact)

	if result[0].Value != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", result[0].Value)
	}
	if result[1].Value != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result[1].Value)
	}
	if result[2].Value != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result[2].Value)
	}
	if result[3].Value != "[REDACTED]" {
		t.Errorf("API_SECRET should be redacted, got %q", result[3].Value)
	}
}

func TestMaskEntries_DoesNotMutateOriginal(t *testing.T) {
	input := []parser.Entry{
		{Key: "DB_PASSWORD", Value: "original"},
	}
	masker.MaskEntries(input, masker.StrategyRedact)
	if input[0].Value != "original" {
		t.Error("original slice was mutated")
	}
}
