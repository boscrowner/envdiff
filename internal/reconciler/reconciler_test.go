package reconciler_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/reconciler"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestPlan_AddedKey(t *testing.T) {
	src := entries("A", "1")
	tgt := entries("A", "1", "B", "2")
	actions := reconciler.Plan(src, tgt)
	if len(actions) != 1 || actions[0].Kind != "add" || actions[0].Key != "B" {
		t.Fatalf("expected one add action for B, got %+v", actions)
	}
}

func TestPlan_RemovedKey(t *testing.T) {
	src := entries("A", "1", "B", "2")
	tgt := entries("A", "1")
	actions := reconciler.Plan(src, tgt)
	if len(actions) != 1 || actions[0].Kind != "remove" || actions[0].Key != "B" {
		t.Fatalf("expected one remove action for B, got %+v", actions)
	}
}

func TestPlan_UpdatedKey(t *testing.T) {
	src := entries("A", "old")
	tgt := entries("A", "new")
	actions := reconciler.Plan(src, tgt)
	if len(actions) != 1 || actions[0].Kind != "update" || actions[0].Value != "new" {
		t.Fatalf("expected one update action, got %+v", actions)
	}
}

func TestApply_AddAndUpdate(t *testing.T) {
	src := entries("A", "1")
	actions := []reconciler.Action{
		{Kind: "add", Key: "B", Value: "2"},
		{Kind: "update", Key: "A", Value: "99"},
	}
	result := reconciler.Apply(src, actions)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Value != "99" {
		t.Errorf("expected A=99, got %s", result[0].Value)
	}
	if result[1].Key != "B" || result[1].Value != "2" {
		t.Errorf("expected B=2, got %+v", result[1])
	}
}

func TestApply_Remove(t *testing.T) {
	src := entries("A", "1", "B", "2")
	actions := []reconciler.Action{{Kind: "remove", Key: "B"}}
	result := reconciler.Apply(src, actions)
	if len(result) != 1 || result[0].Key != "A" {
		t.Fatalf("expected only A, got %+v", result)
	}
}

func TestRender(t *testing.T) {
	entries := entries("FOO", "bar", "BAZ", "qux")
	out := reconciler.Render(entries)
	if !strings.Contains(out, "FOO=bar\n") || !strings.Contains(out, "BAZ=qux\n") {
		t.Errorf("unexpected render output: %q", out)
	}
}
