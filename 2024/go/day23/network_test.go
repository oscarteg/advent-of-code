package main

import (
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(example))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 7 {
		t.Errorf("part 1: got %d, want 7", part1)
	}
	if part2 != "co,de,ka,ta" {
		t.Errorf("part 2: got %s, want co,de,ka,ta", part2)
	}
}

func TestSingleTriangle(t *testing.T) {
	input := "a-b\nb-c\na-c"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part 1: got %d, want 0", part1)
	}
	if part2 != "a,b,c" {
		t.Errorf("part 2: got %s, want a,b,c", part2)
	}
}

func TestTriangleWithT(t *testing.T) {
	input := "a-b\nb-t\na-t"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part 1: got %d, want 1", part1)
	}
	if part2 != "a,b,t" {
		t.Errorf("part 2: got %s, want a,b,t", part2)
	}
}

func TestInvalidFormat(t *testing.T) {
	input := "invalid"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for invalid format")
	}
}
