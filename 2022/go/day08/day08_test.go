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
		t.Fatalf("solve returned error: %v", err)
	}
	if part1 != 21 {
		t.Errorf("Part 1: got %d, want 21", part1)
	}
}

func TestPart2Example(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}
	_, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve returned error: %v", err)
	}
	if part2 != 8 {
		t.Errorf("Part 2: got %d, want 8", part2)
	}
}

func TestSingleTree(t *testing.T) {
	input := "5\n"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve returned error: %v", err)
	}
	if part1 != 1 {
		t.Errorf("Part 1: got %d, want 1", part1)
	}
	if part2 != 0 {
		t.Errorf("Part 2: got %d, want 0", part2)
	}
}

func TestMalformedInput(t *testing.T) {
	input := "12a34\n"
	_, _, err := solve(input)
	if err == nil {
		t.Error("Expected error for malformed input")
	}
}

func TestInconsistentRowLength(t *testing.T) {
	input := "123\n12\n"
	_, _, err := solve(input)
	if err == nil {
		t.Error("Expected error for inconsistent row length")
	}
}
