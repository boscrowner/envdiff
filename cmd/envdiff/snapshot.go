package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourorg/envdiff/internal/differ"
	"github.com/yourorg/envdiff/internal/reporter"
	"github.com/yourorg/envdiff/internal/snapshot"
)

// runSnapshot handles the `envdiff snapshot` subcommand.
// Usage:
//
//	envdiff snapshot take <env-file> <output.json>
//	envdiff snapshot diff <snapshot.json> <env-file>
func runSnapshot(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff snapshot <take|diff> [args...]")
	}
	switch args[0] {
	case "take":
		return runSnapshotTake(args[1:])
	case "diff":
		return runSnapshotDiff(args[1:])
	default:
		return fmt.Errorf("unknown snapshot subcommand %q", args[0])
	}
}

func runSnapshotTake(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff snapshot take <env-file> [output.json]")
	}
	envFile := args[0]
	dest := args[1]
	if len(args) < 2 {
		base := filepath.Base(envFile)
		dest = fmt.Sprintf("%s-%s.json", base, time.Now().UTC().Format("20060102T150405Z"))
	}
	s, err := snapshot.Take(envFile)
	if err != nil {
		return err
	}
	if err := snapshot.Save(s, dest); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "snapshot saved: %s (%d entries)\n", dest, len(s.Entries))
	return nil
}

func runSnapshotDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff snapshot diff <snapshot.json> <env-file>")
	}
	s, err := snapshot.Load(args[0])
	if err != nil {
		return err
	}
	current, err := snapshot.Take(args[1])
	if err != nil {
		return err
	}
	diffs := differ.Diff(s.Entries, current.Entries)
	reporter.TextReport(os.Stdout, diffs)
	return nil
}
