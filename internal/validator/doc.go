// Package validator provides static analysis for .env files.
//
// It inspects a slice of parsed entries and returns a list of Issues
// describing problems such as:
//
//   - invalid-key: the key does not follow the POSIX variable naming
//     convention ([A-Za-z_][A-Za-z0-9_]*).
//
//   - duplicate-key: the same key appears more than once in the file,
//     which leads to ambiguous behaviour across tools.
//
//   - empty-value: the key is present but its value is blank, which is
//     often a sign of an incomplete configuration.
//
// Usage:
//
//	entries, err := parser.ParseFile("staging.env")
//	if err != nil { ... }
//	issues := validator.Validate(entries)
//	for _, iss := range issues {
//		fmt.Println(iss)
//	}
package validator
