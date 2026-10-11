package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 374 {
		t.Errorf("part1: got %d, want 374", part1)
	}
	if part2 != 82000210 {
		t.Errorf("part2: got %d, want 82000210", part2)
	}
}

func TestSolveWithExpansionFactor10(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	galaxies, err := parseGalaxies(string(input))
	if err != nil {
		t.Fatalf("parseGalaxies failed: %v", err)
	}
	result := solveWithExpansion(galaxies, 10)
	if result != 1030 {
		t.Errorf("expansion factor 10: got %d, want 1030", result)
	}
}

func TestSolveWithExpansionFactor100(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	galaxies, err := parseGalaxies(string(input))
	if err != nil {
		t.Fatalf("parseGalaxies failed: %v", err)
	}
	result := solveWithExpansion(galaxies, 100)
	if result != 8410 {
		t.Errorf("expansion factor 100: got %d, want 8410", result)
	}
}

func TestMalformedInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestTwoGalaxiesWithEmptyCol(t *testing.T) {
	galaxies, err := parseGalaxies("#.#")
	if err != nil {
		t.Fatalf("parseGalaxies failed: %v", err)
	}
	result := solveWithExpansion(galaxies, 2)
	if result != 3 {
		t.Errorf("two galaxies: got %d, want 3", result)
	}
}
