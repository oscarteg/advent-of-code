package main

import (
	"errors"
	"os"
	"testing"
)

const (
	secondExample = `#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################
`

	// Facing east at the start, so the end due north costs one turn plus two
	// steps.
	corridor = `###
#E#
#.#
#S#
###
`

	adjacentEnd = `####
#SE#
####
`
)

func TestSolveExamples(t *testing.T) {
	firstExample, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	cases := []struct {
		name      string
		input     string
		wantPart1 int
		wantPart2 int
	}{
		{"first example", string(firstExample), 7036, 45},
		{"second example", secondExample, 11048, 64},
		{"corridor", corridor, 1002, 3},
		{"end next to the start", adjacentEnd, 1, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part1 != tc.wantPart1 {
				t.Errorf("part1 = %d, want %d", part1, tc.wantPart1)
			}
			if part2 != tc.wantPart2 {
				t.Errorf("part2 = %d, want %d", part2, tc.wantPart2)
			}
		})
	}
}

func TestSolveMalformedInput(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"empty input", "", errInvalidMaze},
		{"ragged rows", "####\n#SE#\n###\n", errInvalidMaze},
		{"unknown tile", "####\n#SX#\n####\n", errInvalidMaze},
		{"no start tile", "####\n#.E#\n####\n", errInvalidMaze},
		{"no end tile", "####\n#S.#\n####\n", errInvalidMaze},
		{"two start tiles", "####\n#SS#\n####\n", errInvalidMaze},
		{"walled off end", "#####\n#S#E#\n#####\n", errNoPath},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); !errors.Is(err, tc.wantErr) {
				t.Fatalf("solve(%q) error = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}
