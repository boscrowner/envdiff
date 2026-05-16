package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/envdiff/internal/parser"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	CreatedAt time.Time              `json:"created_at"`
	Source    string                 `json:"source"`
	Entries   []parser.Entry         `json:"entries"`
}

// Take creates a Snapshot from a parsed env file path.
func Take(path string) (*Snapshot, error) {
	entries, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: parse %q: %w", path, err)
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Source:    path,
		Entries:   entries,
	}, nil
}

// Save writes a Snapshot to a JSON file at dest.
func Save(s *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", dest, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a Snapshot from a JSON file at path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
