// Package linter provides static analysis rules for .env files.
//
// It inspects parsed entries and reports issues such as:
//   - keys that are not fully uppercase
//   - keys containing spaces
//   - values with leading or trailing whitespace
//   - values that appear to contain plaintext secrets
//
// Each finding is returned as an Issue with an associated Severity
// (warning or error), the offending key name, and the source line number.
//
// Usage:
//
//	issues, err := linter.Lint(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, iss := range issues {
//		fmt.Println(iss)
//	}
package linter
