package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, part2, err := solve(string(example))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 4361 {
		t.Errorf("part1: got %d, want 4361", part1)
	}
	if part2 != 467835 {
		t.Errorf("part2: got %d, want 467835", part2)
	}
}

func TestSolveEdgeCase_NoSymbols(t *testing.T) {
	input := "123\n456\n789"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part1: got %d, want 0", part1)
	}
	if part2 != 0 {
		t.Errorf("part2: got %d, want 0", part2)
	}
}

func TestSolveEdgeCase_GearWithOneNumber(t *testing.T) {
	input := "..*.\n123."
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 123 {
		t.Errorf("part1: got %d, want 123", part1)
	}
	if part2 != 0 {
		t.Errorf("part2: got %d, want 0 (only one number)", part2)
	}
}

func TestSolveRejectsRaggedRows(t *testing.T) {
	if _, _, err := solve("467..\n...*\n"); err == nil {
		t.Error("want an error for rows of different lengths")
	}
}

func TestSolveEdgeCase_EmptyInput(t *testing.T) {
	input := ""
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part1: got %d, want 0", part1)
	}
	if part2 != 0 {
		t.Errorf("part2: got %d, want 0", part2)
	}
}
