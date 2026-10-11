package main

import (
	"os"
	"testing"
)

func TestExamplePart1(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, _, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 143 {
		t.Errorf("part1 = %d, want 143", part1)
	}
}

func TestExamplePart2(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	_, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part2 != 123 {
		t.Errorf("part2 = %d, want 123", part2)
	}
}

func TestSinglePage(t *testing.T) {
	input := "1|2\n\n1\n"
	part1, _, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 1 {
		t.Errorf("single page part1 = %d, want 1", part1)
	}
}

func TestInvalidRule(t *testing.T) {
	input := "invalid\n\n1,2\n"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for invalid rule, got nil")
	}
}

func TestInvalidPageNumber(t *testing.T) {
	input := "1|abc\n\n1,2\n"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for invalid page number, got nil")
	}
}
