// Package differ provides functionality to compare two parsed .env maps
// and produce a structured diff result.
//
// # Usage
//
// Parse two .env files using the parser package, then pass the resulting
// maps to Diff:
//
//	base, err := parser.ParseFile(".env.example")
//	if err != nil { ... }
//
//	other, err := parser.ParseFile(".env")
//	if err != nil { ... }
//
//	result := differ.Diff(base, other)
//	for _, entry := range result.Entries {
//		fmt.Printf("%s [%s]\n", entry.Key, entry.Kind)
//	}
//
// # Diff Kinds
//
// Each entry in the result has one of the following kinds:
//   - Added: key exists in other but not in base
//   - Removed: key exists in base but not in other
//   - Modified: key exists in both but values differ
//   - Unchanged: key exists in both with the same value
package differ
