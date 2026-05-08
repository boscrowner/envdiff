package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

// runExport handles the `envdiff export` subcommand.
// Usage: envdiff export <file> [--format dotenv|json|export]
func runExport(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff export <file> [--format dotenv|json|export]")
	}

	filePath := args[0]
	format := exporter.FormatDotEnv

	for i := 1; i < len(args)-1; i++ {
		if args[i] == "--format" {
			format = exporter.Format(args[i+1])
		}
	}

	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", filePath, err)
	}

	if err := exporter.Export(os.Stdout, entries, format); err != nil {
		return fmt.Errorf("exporting: %w", err)
	}

	return nil
}
