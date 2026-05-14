// Package merger provides functionality for combining two sets of parsed env
// entries into a single unified result.
//
// # Overview
//
// When managing configuration across environments it is common to maintain a
// base .env file (e.g. .env.defaults) and one or more override files
// (e.g. .env.local). The merger package handles combining these two sources
// while giving the caller explicit control over conflict resolution.
//
// # Strategies
//
// Three strategies are available:
//
//   - PreferBase        – keep the base value when the same key appears in both files.
//   - PreferOverride    – use the override value on conflict (classic "override" semantics).
//   - ErrorOnConflict   – return an error immediately if any key differs between files.
//
// # Usage
//
//	base, _ := parser.ParseFile(".env.defaults")
//	local, _ := parser.ParseFile(".env.local")
//	result, err := merger.Merge(base, local, merger.PreferOverride)
//
The Result value exposes the merged Entries slice, a list of Conflicts
// encountered, and the keys that were newly added from the override file.
package merger
