package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/sorter"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestSort_Alphabetical(t *testing.T) {
	in := entries("ZEBRA", "z", "APPLE", "a", "MANGO", "m")
	res := sorter.Sort(in, sorter.Alphabetical)

	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, e := range res.Entries {
		if e.Key != want[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Key, want[i])
		}
	}
}

func TestSort_AlphabeticalDoesNotMutateOriginal(t *testing.T) {
	in := entries("Z_KEY", "z", "A_KEY", "a")
	origFirst := in[0].Key
	sorter.Sort(in, sorter.Alphabetical)
	if in[0].Key != origFirst {
		t.Errorf("original slice was mutated: got %q, want %q", in[0].Key, origFirst)
	}
}

func TestSort_ByPrefix_GroupsCorrectly(t *testing.T) {
	in := entries(
		"DB_HOST", "localhost",
		"AWS_SECRET", "s3cr3t",
		"DB_PORT", "5432",
		"AWS_KEY", "key",
		"APP_ENV", "prod",
	)
	res := sorter.Sort(in, sorter.ByPrefix)

	if res.Groups == nil {
		t.Fatal("expected Groups to be populated for ByPrefix strategy")
	}

	if len(res.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(res.Groups["DB"]))
	}
	if len(res.Groups["AWS"]) != 2 {
		t.Errorf("expected 2 AWS entries, got %d", len(res.Groups["AWS"]))
	}
	if len(res.Groups["APP"]) != 1 {
		t.Errorf("expected 1 APP entry, got %d", len(res.Groups["APP"]))
	}
}

func TestSort_ByPrefix_OrderedOutput(t *testing.T) {
	in := entries(
		"DB_PORT", "5432",
		"APP_ENV", "prod",
		"DB_HOST", "localhost",
	)
	res := sorter.Sort(in, sorter.ByPrefix)

	// APP comes before DB alphabetically
	want := []string{"APP_ENV", "DB_HOST", "DB_PORT"}
	for i, e := range res.Entries {
		if e.Key != want[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Key, want[i])
		}
	}
}

func TestSort_ByPrefix_NoPrefixKey(t *testing.T) {
	in := entries("STANDALONE", "val", "DB_HOST", "localhost")
	res := sorter.Sort(in, sorter.ByPrefix)

	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	// STANDALONE has no underscore, uses full key as prefix
	if _, ok := res.Groups["STANDALONE"]; !ok {
		t.Error("expected STANDALONE to appear as its own group")
	}
}
