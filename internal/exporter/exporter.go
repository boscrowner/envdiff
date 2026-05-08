package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Format represents the output format for exported env data.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatExport Format = "export"
)

// Export writes the given entries to w in the specified format.
func Export(w io.Writer, entries []parser.Entry, format Format) error {
	sorted := make([]parser.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	switch format {
	case FormatDotEnv:
		return writeDotEnv(w, sorted)
	case FormatJSON:
		return writeJSON(w, sorted)
	case FormatExport:
		return writeExport(w, sorted)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func writeDotEnv(w io.Writer, entries []parser.Entry) error {
	for _, e := range entries {
		val := quoteIfNeeded(e.Value)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, entries []parser.Entry) error {
	if _, err := fmt.Fprintln(w, "{"); err != nil {
		return err
	}
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", e.Key, e.Value, comma); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, "}")
	return err
}

func writeExport(w io.Writer, entries []parser.Entry) error {
	for _, e := range entries {
		val := quoteIfNeeded(e.Value)
		if _, err := fmt.Fprintf(w, "export %s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func quoteIfNeeded(val string) string {
	if strings.ContainsAny(val, " \t#$") {
		return fmt.Sprintf("%q", val)
	}
	return val
}
