package main

import (
	"os"
	"testing"
)

func TestExampleInput(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 14 {
		t.Errorf("part 1: expected 14, got %d", part1)
	}
	if part2 != 34 {
		t.Errorf("part 2: expected 34, got %d", part2)
	}
}

func TestEmptyGrid(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("expected error for empty grid")
	}
}

func TestIrregularGrid(t *testing.T) {
	_, _, err := solve("...\n..\n...")
	if err == nil {
		t.Error("expected error for irregular grid")
	}
}
