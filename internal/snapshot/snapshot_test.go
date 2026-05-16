package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/snapshot"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestTake_PopulatesFields(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	before := time.Now().UTC()
	s, err := snapshot.Take(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Source != path {
		t.Errorf("source: got %q, want %q", s.Source, path)
	}
	if !s.CreatedAt.After(before.Add(-time.Second)) {
		t.Errorf("created_at not recent: %v", s.CreatedAt)
	}
	if len(s.Entries) != 2 {
		t.Errorf("entries: got %d, want 2", len(s.Entries))
	}
}

func TestTake_MissingFile(t *testing.T) {
	_, err := snapshot.Take("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nOTHER=123\n")
	s, err := snapshot.Take(path)
	if err != nil {
		t.Fatalf("take: %v", err)
	}
	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, dest); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Source != s.Source {
		t.Errorf("source mismatch: got %q, want %q", loaded.Source, s.Source)
	}
	if len(loaded.Entries) != len(s.Entries) {
		t.Errorf("entries count mismatch: got %d, want %d", len(loaded.Entries), len(s.Entries))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("not json")
	f.Close()
	_, err = snapshot.Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
