package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]Entry
}

// ParseFile reads and parses a .env file at the given path.
func ParseFile(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		entry, ok := parseLine(raw, lineNum)
		if !ok {
			continue
		}
		env.Entries = append(env.Entries, entry)
		env.Index[entry.Key] = entry
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", path, err)
	}

	return env, nil
}

// parseLine parses a single line into an Entry.
func parseLine(line string, lineNum int) (Entry, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return Entry{}, false
	}

	var comment string
	if idx := strings.Index(trimmed, " #"); idx != -1 {
		comment = strings.TrimSpace(trimmed[idx+2:])
		trimmed = strings.TrimSpace(trimmed[:idx])
	}

	parts := strings.SplitN(trimmed, "=", 2)
	if len(parts) != 2 {
		return Entry{}, false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

	return Entry{
		Key:     key,
		Value:   value,
		Comment: comment,
		Line:    lineNum,
	}, true
}
