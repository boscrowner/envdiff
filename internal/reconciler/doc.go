// Package reconciler provides functionality to plan and apply reconciliation
// actions between two sets of .env entries.
//
// Workflow:
//
//  1. Parse source and target .env files using the parser package.
//  2. Call Plan(source, target) to obtain a slice of Actions describing what
//     must change in source to match target (adds, removals, updates).
//  3. Optionally inspect or filter the Actions before applying them.
//  4. Call Apply(source, actions) to produce a new slice of entries with all
//     actions merged in, preserving the original key order where possible.
//  5. Call Render(entries) to serialise the reconciled entries back to a
//     .env-formatted string suitable for writing to disk.
//
// Example:
//
//	src, _ := parser.ParseFile("staging.env")
//	tgt, _ := parser.ParseFile("production.env")
//	actions := reconciler.Plan(src, tgt)
//	reconciled := reconciler.Apply(src, actions)
//	fmt.Print(reconciler.Render(reconciled))
package reconciler
