package main

import (
	"os"
	"strings"
	"testing"
)

func readExample(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	return string(input)
}

func TestOfficialExample(t *testing.T) {
	part1, part2, err := solve(readExample(t))

	if err != nil || part1 != 35 || part2 != 3351 {
		t.Fatalf("got (%d, %d, %v), want (35, 3351, nil)", part1, part2, err)
	}
}

func TestSingleStepMatchesKnownIntermediateCount(t *testing.T) {
	part1, part2, err := solveWith(readExample(t), 1, 1)

	if err != nil || part1 != 24 || part2 != 24 {
		t.Fatalf("got (%d, %d, %v), want (24, 24, nil)", part1, part2, err)
	}
}

func TestMalformedInput(t *testing.T) {
	validAlgorithm := strings.Repeat("#", 512)
	tests := []struct {
		name  string
		input string
	}{
		{"missing blank line", "...\n#.."},
		{"wrong algorithm length", "...\n\n#.."},
		{"invalid algorithm character", strings.Repeat("#", 511) + "x\n\n#..\n.#.\n..#"},
		{"invalid image pixel", validAlgorithm + "\n\n#.x\n.#.\n..#"},
		{"uneven image rows", validAlgorithm + "\n\n#.\n.#."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}
