// Package parser provides functionality for reading and parsing .env files
// into structured representations suitable for diffing and reconciliation.
//
// A .env file is expected to follow the common KEY=VALUE format, with support
// for:
//   - Blank lines (ignored)
//   - Full-line comments starting with '#' (ignored)
//   - Inline comments after a value using ' #' as delimiter
//   - Double-quoted values (quotes are stripped)
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range env.Entries {
//		fmt.Printf("%s = %s\n", entry.Key, entry.Value)
//	}
package parser
