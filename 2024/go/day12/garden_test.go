package main

import (
	"os"
	"testing"
)

const (
	fiveRegions = `AAAA
BBCD
BBCC
EEEC
`
	nestedRegions = `OOOOO
OXOXO
OOOOO
OXOXO
OOOOO
`
	eShape = `EEEEE
EXXXX
EEEEE
EXXXX
EEEEE
`
	diagonallyTouching = `AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
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
		{"five regions", fiveRegions, 140, 80},
		{"nested regions", nestedRegions, 772, 436},
		{"large example", string(largeExample), 1930, 1206},
		{"single plot", "A\n", 4, 4},
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

// The puzzle states only the side-counting price for these two maps.
func TestSolveSideCountingExamples(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantPart2 int
	}{
		{"e shape", eShape, 236},
		{"diagonally touching regions", diagonallyTouching, 368},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part2 != tc.wantPart2 {
				t.Errorf("part2 = %d, want %d", part2, tc.wantPart2)
			}
		})
	}
}

func TestSolveMalformedInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"ragged rows", "AAAA\nBB\n"},
		{"plot is not a letter", "AA\nA1\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want error", tc.input)
			}
		})
	}
}
