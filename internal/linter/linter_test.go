package linter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return path
}

func TestLint_Clean(t *testing.T) {
	path := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	issues, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	path := writeTempEnv(t, "app_host=localhost\n")
	issues, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
}

func TestLint_WhitespaceValue(t *testing.T) {
	path := writeTempEnv(t, "APP_HOST= localhost \n")
	issues, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
}

func TestLint_PlaintextSecret(t *testing.T) {
	path := writeTempEnv(t, "DB_PASSWORD=mysecretpassword\n")
	issues, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for plaintext secret")
	}
	found := false
	for _, iss := range issues {
		if iss.Severity == linter.SeverityWarning && iss.Key == "DB_PASSWORD" {
			found = true
		}
	}
	if !found {
		t.Error("expected plaintext secret warning for DB_PASSWORD")
	}
}

func TestLint_FileNotFound(t *testing.T) {
	_, err := linter.Lint("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
