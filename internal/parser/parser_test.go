package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(env.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(env.Entries))
	}

	if env.Index["DB_HOST"].Value != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", env.Index["DB_HOST"].Value)
	}
}

func TestParseFile_IgnoresComments(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\nKEY=value\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(env.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"\n`)

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := env.Index["SECRET"].Value; v != "my secret value" {
		t.Errorf("expected unquoted value, got %q", v)
	}
}

func TestParseFile_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # default port\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := env.Index["PORT"]
	if e.Value != "8080" {
		t.Errorf("expected value 8080, got %q", e.Value)
	}
	if e.Comment != "default port" {
		t.Errorf("expected comment 'default port', got %q", e.Comment)
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
