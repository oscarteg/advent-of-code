package main

import (
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 6 {
		t.Errorf("part 1: expected 6, got %d", part1)
	}
	if part2 != 16 {
		t.Errorf("part 2: expected 16, got %d", part2)
	}
}

func TestSinglePattern(t *testing.T) {
	input := "a\na\naa\naaa"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 3 {
		t.Errorf("part 1: expected 3, got %d", part1)
	}
	if part2 != 3 {
		t.Errorf("part 2: expected 3, got %d", part2)
	}
}

func TestImpossibleDesign(t *testing.T) {
	input := "a, b\nc"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part 1: expected 0, got %d", part1)
	}
	if part2 != 0 {
		t.Errorf("part 2: expected 0, got %d", part2)
	}
}

func TestEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}
