// Package snapshot provides functionality for capturing, saving, and loading
// point-in-time snapshots of .env files.
//
// A Snapshot records the parsed entries of an env file along with metadata
// such as the source path and the time the snapshot was taken. Snapshots are
// persisted as JSON files, allowing teams to track environment state over time
// and compare historical configurations.
//
// Typical usage:
//
//	s, err := snapshot.Take(".env.production")
//	if err != nil { ... }
//	err = snapshot.Save(s, "snapshots/prod-2024-01-15.json")
//
//	// Later, load and compare:
//	old, err := snapshot.Load("snapshots/prod-2024-01-15.json")
package snapshot
