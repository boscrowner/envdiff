package exporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

var testEntries = []parser.Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "SECRET", Value: "my secret value"},
}

func TestExport_DotEnv(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(&buf, testEntries, exporter.FormatDotEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
}

func TestExport_DotEnv_QuotesSpaces(t *testing.T) {
	var buf strings.Builder
	entries := []parser.Entry{{Key: "MSG", Value: "hello world"}}
	_ = exporter.Export(&buf, entries, exporter.FormatDotEnv)
	out := buf.String()
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(&buf, testEntries, exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON to start with '{', got: %s", out)
	}
	if !strings.Contains(out, `"DB_HOST": "localhost"`) {
		t.Errorf("expected DB_HOST in JSON output, got:\n%s", out)
	}
}

func TestExport_ExportFormat(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(&buf, testEntries, exporter.FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected export statement, got:\n%s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	var buf strings.Builder
	err := exporter.Export(&buf, testEntries, exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_SortedOutput(t *testing.T) {
	var buf strings.Builder
	_ = exporter.Export(&buf, testEntries, exporter.FormatDotEnv)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected APP_ENV first (sorted), got: %s", lines[0])
	}
}
