// Package exporter provides functionality for serializing parsed .env entries
// into multiple output formats.
//
// Supported formats:
//
//   - dotenv  — standard KEY=VALUE format, compatible with most tools
//   - json    — JSON object mapping keys to values
//   - export  — shell-compatible format using `export KEY=VALUE` syntax
//
// Values containing spaces, tabs, hash signs, or dollar signs are automatically
// quoted in dotenv and export formats to ensure safe parsing.
//
// Entries are always written in lexicographic key order for deterministic output.
//
// Example:
//
//	var buf bytes.Buffer
//	err := exporter.Export(&buf, entries, exporter.FormatJSON)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(buf.String())
package exporter
