package main

import (
	"os"
	"testing"
)

func TestSolveWithExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	part1, part2, err := solveWith(string(input), 7, 12)
	if err != nil {
		t.Fatalf("solveWith: %v", err)
	}
	if part1 != 22 {
		t.Errorf("part1 = %d, want 22", part1)
	}
	if part2 != "6,1" {
		t.Errorf("part2 = %q, want %q", part2, "6,1")
	}
}

func TestSolveWithTwoByteGrid(t *testing.T) {
	part1, part2, err := solveWith("1,0\n0,1\n", 2, 0)
	if err != nil {
		t.Fatalf("solveWith: %v", err)
	}
	if part1 != 2 {
		t.Errorf("part1 = %d, want 2", part1)
	}
	if part2 != "0,1" {
		t.Errorf("part2 = %q, want %q", part2, "0,1")
	}
}

func TestSolveWithMalformedInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
		size  int
		bytes int
	}{
		{"empty input", "", 7, 0},
		{"missing comma", "not a coordinate\n", 7, 1},
		{"not enough bytes", "0,0\n", 7, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solveWith(tc.input, tc.size, tc.bytes); err == nil {
				t.Fatalf("solveWith(%q) = nil error, want error", tc.input)
			}
		})
	}
}
