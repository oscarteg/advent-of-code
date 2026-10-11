package main

import (
	"os"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	first, _, err := solve(string(input))
	if err != nil {
		t.Fatal(err)
	}
	if first != 32000000 {
		t.Fatalf("part1 = %d, want 32000000", first)
	}
}

func TestSecondExample(t *testing.T) {
	const input = `broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output
`
	first, _, err := solve(input)
	if err != nil {
		t.Fatal(err)
	}
	if first != 11687500 {
		t.Fatalf("part1 = %d, want 11687500", first)
	}
}

// Synthetic module graph (not a puzzle example): two binary counters of
// period 2 and 4 feed single-input conjunctions inv1/inv2, which feed con,
// the sole input of the undefined sink rx. inv1 first sends a high pulse at
// press 2, inv2 at press 4, so the answer is lcm(2, 4) = 4.
const syntheticCycleGraph = `broadcaster -> a, b
%a -> inv1
%b -> c
%c -> inv2
&inv1 -> con
&inv2 -> con
&con -> rx
`

func TestPart2SyntheticCycle(t *testing.T) {
	_, second, err := solve(syntheticCycleGraph)
	if err != nil {
		t.Fatal(err)
	}
	if second != 4 {
		t.Fatalf("part2 = %d, want 4", second)
	}
}

func TestInvalidConfigurations(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"missing arrow", "broadcaster a, b"},
		{"empty input", ""},
		{"no destinations", "broadcaster ->"},
		{"ambiguous sink", "broadcaster -> a, b\n%a -> rx\n%b -> rx2"},
		{"target with two inputs", "broadcaster -> a, b\n%a -> rx\n%b -> rx"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) succeeded, want error", tc.input)
			}
		})
	}
}

func TestNoImplicitSinkSkipsPart2(t *testing.T) {
	_, second, err := solve("broadcaster -> a\n%a -> broadcaster")
	if err != nil {
		t.Fatal(err)
	}
	if second != 0 {
		t.Fatalf("part2 = %d, want 0 when there is no rx-style target", second)
	}
}
