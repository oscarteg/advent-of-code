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

	if part1 != 3 {
		t.Errorf("part1 = %d, want 3", part1)
	}
	if part2 != 6 {
		t.Errorf("part2 = %d, want 6", part2)
	}
}

func TestSingleRotationEndingAtZero(t *testing.T) {
	input := "R50\n"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part1 = %d, want 1", part1)
	}
	if part2 != 1 {
		t.Errorf("part2 = %d, want 1", part2)
	}
}

func TestLeftRotationFromZero(t *testing.T) {
	input := "R50\nL50\n"
	part1, _, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part1 = %d, want 1", part1)
	}
}

func TestMultipleCrossings(t *testing.T) {
	input := "R250\n"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part1 = %d, want 1", part1)
	}
	if part2 != 3 {
		t.Errorf("part2 = %d, want 3", part2)
	}
}

func TestMalformedDirectionReturnsError(t *testing.T) {
	input := "X50\n"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for invalid direction, got nil")
	}
}

func TestMissingDistanceReturnsError(t *testing.T) {
	input := "L\n"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for missing distance, got nil")
	}
}

func TestEmptyLineHandling(t *testing.T) {
	input := "R50\n\nL50\n"
	part1, _, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part1 = %d, want 1", part1)
	}
}
