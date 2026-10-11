package main

import (
	"os"
	"testing"
)

func TestSolveExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 46 {
		t.Errorf("part1: got %d, want 46", part1)
	}

	if part2 != 51 {
		t.Errorf("part2: got %d, want 51", part2)
	}
}

func TestSolveEmpty(t *testing.T) {
	part1, part2, err := solve("")
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 || part2 != 0 {
		t.Errorf("empty input: got (%d, %d), want (0, 0)", part1, part2)
	}
}

func TestSolveSingleDot(t *testing.T) {
	part1, part2, err := solve(".")
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("single dot part1: got %d, want 1", part1)
	}

	if part2 != 1 {
		t.Errorf("single dot part2: got %d, want 1", part2)
	}
}

func TestSolveRejectsMalformedGrid(t *testing.T) {
	for _, input := range []string{".|.\n..", ".x.\n..."} {
		if _, _, err := solve(input); err == nil {
			t.Errorf("solve(%q): want an error", input)
		}
	}
}
