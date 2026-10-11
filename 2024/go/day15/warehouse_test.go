package main

import (
	"errors"
	"os"
	"testing"
)

const (
	smallExample = `########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<
`
	verticalPushExample = `#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######

<vv<<^^<<^^
`
	blockedPush = `#####
#.O@#
#####

<<
`
)

func TestSolveExamples(t *testing.T) {
	largeExample, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	cases := []struct {
		name      string
		input     string
		wantPart1 int
		wantPart2 int
	}{
		{"large example", string(largeExample), 10092, 9021},
		// Edge case: the second push runs into the wall and leaves the narrow
		// map untouched, while the widened map has one floor tile more and
		// takes both pushes.
		{"push into a wall", blockedPush, 101, 102},
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

// The puzzle states only one part's sum for each of these maps.
func TestSolveSinglePartExamples(t *testing.T) {
	cases := []struct {
		name  string
		input string
		part  int
		want  int
	}{
		{"smaller example", smallExample, 1, 2028},
		{"vertical push of half-aligned boxes", verticalPushExample, 2, 618},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			got := part1
			if tc.part == 2 {
				got = part2
			}
			if got != tc.want {
				t.Errorf("part%d = %d, want %d", tc.part, got, tc.want)
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
		{"empty input", "", errInvalidWarehouse},
		{"ragged rows", "#####\n###\n\n<\n", errInvalidWarehouse},
		{"unknown tile", "#####\n#.X@#\n#####\n\n<\n", errInvalidWarehouse},
		{"no robot", "#####\n#...#\n#####\n\n<\n", errInvalidWarehouse},
		{"two robots", "#####\n#@@.#\n#####\n\n<\n", errInvalidWarehouse},
		{"map not enclosed", "#####\n#.@..\n#####\n\n<\n", errInvalidWarehouse},
		{"no moves", "#####\n#.@.#\n#####\n", errInvalidMoves},
		{"unknown move", "#####\n#.@.#\n#####\n\n<x\n", errInvalidMoves},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); !errors.Is(err, tc.wantErr) {
				t.Fatalf("solve(%q) error = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}
