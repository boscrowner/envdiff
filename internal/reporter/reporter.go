// Package reporter formats and outputs diff results for human or machine consumption.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format defines the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// TextReport writes a human-readable diff report to w.
func TextReport(w io.Writer, diffs []differ.Diff) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	summary := differ.Summarize(diffs)
	_, err := fmt.Fprintf(w, "Summary: %d added, %d removed, %d modified, %d unchanged\n\n",
		summary[differ.KindAdded],
		summary[differ.KindRemoved],
		summary[differ.KindModified],
		summary[differ.KindUnchanged],
	)
	if err != nil {
		return err
	}

	for _, d := range diffs {
		line, err := formatDiffLine(d)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatDiffLine(d differ.Diff) (string, error) {
	var sb strings.Builder
	switch d.Kind {
	case differ.KindAdded:
		sb.WriteString(fmt.Sprintf("+ %-30s = %s", d.Key, d.NewValue))
	case differ.KindRemoved:
		sb.WriteString(fmt.Sprintf("- %-30s = %s", d.Key, d.OldValue))
	case differ.KindModified:
		sb.WriteString(fmt.Sprintf("~ %-30s : %s -> %s", d.Key, d.OldValue, d.NewValue))
	case differ.KindUnchanged:
		sb.WriteString(fmt.Sprintf("  %-30s = %s", d.Key, d.OldValue))
	default:
		return "", fmt.Errorf("unknown diff kind: %q", d.Kind)
	}
	return sb.String(), nil
}
