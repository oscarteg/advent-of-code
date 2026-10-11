package main

import (
	"os"
	"strings"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 405 || part2 != 400 {
		t.Fatalf("got (%d, %d, %v), want (405, 400, nil)", part1, part2, err)
	}
}

func TestSinglePatternExamples(t *testing.T) {
	tests := []struct {
		name  string
		grid  string
		part1 int
		part2 int
	}{
		{
			name:  "horizontal mirror at the edge",
			grid:  "#.##..##.\n..#.##.#.\n##......#\n##......#\n..#.##.#.\n..##..##.\n#.#.##.#.",
			part1: 5,
			part2: 300,
		},
		{
			name:  "vertical mirror in the middle",
			grid:  "#...##..#\n#....#..#\n..##..###\n#####.##.\n#####.##.\n..##..###\n#....#..#",
			part1: 400,
			part2: 100,
		},
		{
			name:  "vertical mirror at the edge",
			grid:  "##",
			part1: 1,
			part2: 0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, err := summarize(pattern(strings.Split(tc.grid, "\n")), 0)
			if err != nil {
				t.Fatalf("part 1: %v", err)
			}
			if part1 != tc.part1 {
				t.Fatalf("part 1: got %d, want %d", part1, tc.part1)
			}

			if tc.part2 == 0 {
				return
			}
			part2, err := summarize(pattern(strings.Split(tc.grid, "\n")), 1)
			if err != nil {
				t.Fatalf("part 2: %v", err)
			}
			if part2 != tc.part2 {
				t.Fatalf("part 2: got %d, want %d", part2, tc.part2)
			}
		})
	}
}

func TestMalformedInputReturnsErrors(t *testing.T) {
	for _, input := range []string{
		"",
		"#.#\n.#x",
		"##\n#",
		"##\n\n#",
	} {
		if _, _, err := solve(input); err == nil {
			t.Fatalf("expected error for input %q", input)
		}
	}
}

func TestNoSingleSmudgeFixReturnsError(t *testing.T) {
	if _, _, err := solve("#.\n.#"); err == nil {
		t.Fatal("expected error when no single-smudge mirror exists")
	}
}
