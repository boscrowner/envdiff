package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestRunExport_MissingArgs(t *testing.T) {
	err := runExport([]string{})
	if err == nil {
		t.Error("expected error for missing args, got nil")
	}
}

func TestRunExport_FileNotFound(t *testing.T) {
	err := runExport([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunExport_DefaultFormat(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=value\nFOO=bar\n")
	err := runExport([]string{path})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunExport_JSONFormat(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=value\n")
	err := runExport([]string{path, "--format", "json"})
	if err != nil {
		t.Errorf("unexpected error with json format: %v", err)
	}
}

func TestRunExport_ExportFormat(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=value\n")
	err := runExport([]string{path, "--format", "export"})
	if err != nil {
		t.Errorf("unexpected error with export format: %v", err)
	}
}
