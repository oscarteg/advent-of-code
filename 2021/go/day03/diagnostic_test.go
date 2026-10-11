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

	if part1 != 198 {
		t.Errorf("part1: got %d, want 198", part1)
	}
	if part2 != 230 {
		t.Errorf("part2: got %d, want 230", part2)
	}
}

func TestSolveSingleLine(t *testing.T) {
	part1, part2, err := solve("10101")
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 210 {
		t.Errorf("part1: got %d, want 210", part1)
	}
	if part2 != 441 {
		t.Errorf("part2: got %d, want 441", part2)
	}
}

func TestSolveAlternatingBits(t *testing.T) {
	part1, part2, err := solve("10\n01")
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part1: got %d, want 0", part1)
	}
	if part2 != 2 {
		t.Errorf("part2: got %d, want 2", part2)
	}
}

func TestSolveEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("empty input should fail")
	}
}

func TestSolveNonBinary(t *testing.T) {
	_, _, err := solve("10102")
	if err == nil {
		t.Error("non-binary character should fail")
	}
}

func TestSolveUnevenLines(t *testing.T) {
	_, _, err := solve("101\n10")
	if err == nil {
		t.Error("uneven line lengths should fail")
	}
}
