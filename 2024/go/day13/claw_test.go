package main

import (
	"os"
	"testing"
)

func TestPart1Example(t *testing.T) {
	example, _ := os.ReadFile("testdata/example.txt")
	part1, _, _ := solve(string(example))
	if part1 != 480 {
		t.Errorf("Part 1 example: got %d, want 480", part1)
	}
}

func TestPart2Example(t *testing.T) {
	example, _ := os.ReadFile("testdata/example.txt")
	_, part2, _ := solve(string(example))
	if part2 != 875318608908 {
		t.Errorf("Part 2 example: got %d, want 875318608908", part2)
	}
}

func TestUnsolvable(t *testing.T) {
	input := "Button A: X+1, Y+1\nButton B: X+2, Y+2\nPrize: X=3, Y=4\n"
	part1, _, _ := solve(input)
	if part1 != 0 {
		t.Errorf("Unsolvable machine: got %d, want 0", part1)
	}
}

func TestSolvable(t *testing.T) {
	input := "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3, Y=2\n"
	part1, _, _ := solve(input)
	if part1 != 11 {
		t.Errorf("Solvable machine: got %d, want 11", part1)
	}
}

func TestEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("Empty input should error")
	}
}
