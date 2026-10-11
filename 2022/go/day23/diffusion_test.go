package main

import (
	"os"
	"testing"
)

func TestPart1Example(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}
	part1, _, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 110 {
		t.Errorf("part1: got %d, want 110", part1)
	}
}

func TestPart2Example(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}
	_, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part2 != 20 {
		t.Errorf("part2: got %d, want 20", part2)
	}
}

func TestSingleElf(t *testing.T) {
	input := "#"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 0 {
		t.Errorf("part1: got %d, want 0", part1)
	}
	if part2 != 1 {
		t.Errorf("part2: got %d, want 1", part2)
	}
}

func TestMalformedInput(t *testing.T) {
	input := "#\n$\n#"
	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for malformed input, got nil")
	}
}

func TestNoElves(t *testing.T) {
	input := "...\n..."
	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for no elves, got nil")
	}
}
