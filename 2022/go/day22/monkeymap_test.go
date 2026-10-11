package main

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// A 1-4-1 net with one tile per face: a layout nothing like the example's.
const (
	unitNet     = " .\n....\n  .\n\n2"
	unitNetWall = " .\n....\n  #\n\n2"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 6032 || part2 != 5031 {
		t.Fatalf("got (%d, %d, %v), want (6032, 5031, nil)", part1, part2, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
	}{
		// Flat wrapping never leaves the top row; on the cube the walk crosses
		// onto the right face and then onto the front face.
		{"single tile faces", unitNet, 1008, 3013},
		{"wall stops the walk after a crossing", unitNetWall, 1008, 2013},
		{"windows line endings", strings.ReplaceAll(unitNet, "\n", "\r\n"), 1008, 3013},
		{"path ending in a turn", " .\n....\n  .\n\n2R", 1009, 3014},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if err != nil {
				t.Fatalf("solve() error = %v", err)
			}
			if part1 != tc.part1 || part2 != tc.part2 {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}

func TestSolveRejectsMalformedNotes(t *testing.T) {
	tests := []struct{ name, input string }{
		{"missing path", "....\n...."},
		{"missing map", ""},
		{"map character", " .\n..x.\n  .\n\n2"},
		{"path character", " .\n....\n  .\n\n2X"},
		{"tile count is not six faces", "....\n....\n\n1"},
		{"not a cube net", "......\n......\n\n1"},
		{"no open start tile", " #\n####\n  #\n\n2"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); !errors.Is(err, errInvalidNotes) {
				t.Fatalf("solve() error = %v, want %v", err, errInvalidNotes)
			}
		})
	}
}
