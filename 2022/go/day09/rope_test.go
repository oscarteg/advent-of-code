package main

import (
	"os"
	"testing"
)

func TestSolveOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve returned error: %v", err)
	}
	if part1 != 13 {
		t.Errorf("part1 = %d, want 13", part1)
	}
	if part2 != 1 {
		t.Errorf("part2 = %d, want 1", part2)
	}
}

func TestSolveLargerExamplePart2(t *testing.T) {
	input := "R 5\nU 8\nL 8\nD 3\nR 17\nD 10\nL 25\nU 20\n"
	_, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve returned error: %v", err)
	}
	if part2 != 36 {
		t.Errorf("part2 = %d, want 36", part2)
	}
}

func TestSolveSingleStep(t *testing.T) {
	part1, part2, err := solve("R 1\n")
	if err != nil {
		t.Fatalf("solve returned error: %v", err)
	}
	if part1 != 1 || part2 != 1 {
		t.Errorf("solve(\"R 1\") = (%d, %d), want (1, 1)", part1, part2)
	}
}

func TestSolveEmptyInputIsRejected(t *testing.T) {
	if _, _, err := solve(""); err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestSolveMalformedDirectionIsRejected(t *testing.T) {
	if _, _, err := solve("X 4\n"); err == nil {
		t.Error("expected error for unknown direction, got nil")
	}
}

func TestSolveMalformedStepCountIsRejected(t *testing.T) {
	if _, _, err := solve("R four\n"); err == nil {
		t.Error("expected error for invalid step count, got nil")
	}
}
