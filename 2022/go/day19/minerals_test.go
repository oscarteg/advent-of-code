package main

import (
	"errors"
	"os"
	"testing"
)

const secondBlueprint = "Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian."

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 33 || part2 != 3472 {
		t.Fatalf("got (%d, %d, %v), want (33, 3472, nil)", part1, part2, err)
	}
}

func TestBlueprintSchedules(t *testing.T) {
	tests := []struct {
		name, input  string
		part1, part2 int
	}{
		{
			// The example's second blueprint cracks 12 geodes in 24 minutes
			// and 62 in 32 minutes, so on its own it scores 2 * 12 and 62.
			name:  "quality level uses the blueprint id",
			input: secondBlueprint,
			part1: 24,
			part2: 62,
		},
		{
			name:  "unaffordable blueprint cracks nothing",
			input: "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 50 clay. Each geode robot costs 90 ore and 90 obsidian.",
			part1: 0,
			part2: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			part1, part2, err := solve(test.input)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if part1 != test.part1 || part2 != test.part2 {
				t.Errorf("got (%d, %d), want (%d, %d)", part1, part2, test.part1, test.part2)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct{ name, input string }{
		{name: "blank input", input: "\n  \n"},
		{name: "truncated blueprint", input: "Blueprint 1: Each ore robot costs 4 ore."},
		{name: "non numeric id", input: "Blueprint two: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian."},
		{name: "trailing text", input: secondBlueprint + " and a pony."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := solve(test.input); !errors.Is(err, errInvalidBlueprint) {
				t.Fatalf("got %v, want an %v", err, errInvalidBlueprint)
			}
		})
	}
}
