package main

import (
	"os"
	"testing"
)

func TestSolveWithExampleInput(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example file: %v", err)
	}

	part1, part2, err := solve(string(example))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 26 {
		t.Errorf("part 1: got %d, want 26", part1)
	}
	if part2 != 61229 {
		t.Errorf("part 2: got %d, want 61229", part2)
	}
}

func TestSolveSingleEntry(t *testing.T) {
	input := "acedgfb cdfbe gcdfa fbcad dab cefabd cdfgeb eafb cagedb ab | cdfeb fcadb cdfeb cdbaf"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part 1: got %d, want 0 (no easy digits in output)", part1)
	}
	if part2 != 5353 {
		t.Errorf("part 2: got %d, want 5353", part2)
	}
}

func TestSolveEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Errorf("expected error for empty input")
	}
}

func TestSolveInvalidFormat(t *testing.T) {
	input := "abc def"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for invalid format (missing pipe)")
	}
}

func TestSolveTooFewPatterns(t *testing.T) {
	input := "abc def | ghi jkl"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for too few patterns")
	}
}

func TestSolveWrongOutputCount(t *testing.T) {
	input := "abc def ghi jkl mno pqr stu vwx yz ab | cd ef"
	_, _, err := solve(input)
	if err == nil {
		t.Errorf("expected error for wrong output count")
	}
}
