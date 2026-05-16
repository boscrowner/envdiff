package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSnapshot_MissingSubcommand(t *testing.T) {
	err := runSnapshot([]string{})
	if err == nil {
		t.Fatal("expected error for missing subcommand")
	}
}

func TestRunSnapshot_UnknownSubcommand(t *testing.T) {
	err := runSnapshot([]string{"unknown"})
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestRunSnapshotTake_MissingArgs(t *testing.T) {
	err := runSnapshot([]string{"take"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunSnapshotTake_WritesFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("APP=test\nPORT=8080\n"), 0644); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "snap.json")
	err := runSnapshot([]string{"take", envPath, dest})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Error("expected snapshot file to exist")
	}
}

func TestRunSnapshotDiff_MissingArgs(t *testing.T) {
	err := runSnapshot([]string{"diff"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunSnapshotDiff_ShowsDiff(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("APP=test\nPORT=8080\n"), 0644); err != nil {
		t.Fatal(err)
	}
	snapPath := filepath.Join(dir, "snap.json")
	if err := runSnapshot([]string{"take", envPath, snapPath}); err != nil {
		t.Fatalf("take: %v", err)
	}
	// Modify the env file
	if err := os.WriteFile(envPath, []byte("APP=changed\nPORT=8080\nNEW=key\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := runSnapshot([]string{"diff", snapPath, envPath})
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
}
