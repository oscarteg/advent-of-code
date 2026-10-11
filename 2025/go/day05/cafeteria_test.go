package main

import (
	"os"
	"testing"
)

func TestSolveExample(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(example))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 3 {
		t.Errorf("part1: got %d, want 3", part1)
	}
	if part2 != 14 {
		t.Errorf("part2: got %d, want 14", part2)
	}
}

func TestSolveSingleRange(t *testing.T) {
	input := "5-10\n\n7\n3\n12"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part1: got %d, want 1", part1)
	}
	if part2 != 6 {
		t.Errorf("part2: got %d, want 6", part2)
	}
}

func TestSolveNoRanges(t *testing.T) {
	input := "\n\n1\n2\n3"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part1: got %d, want 0", part1)
	}
	if part2 != 0 {
		t.Errorf("part2: got %d, want 0", part2)
	}
}

func TestSolveMalformedRange(t *testing.T) {
	input := "5-\n\n1"
	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for malformed range, got nil")
	}
}

func TestSolveAdjacentRanges(t *testing.T) {
	input := "1-3\n4-6\n\n2\n5"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 2 {
		t.Errorf("part1: got %d, want 2", part1)
	}
	if part2 != 6 {
		t.Errorf("part2: got %d, want 6", part2)
	}
}
