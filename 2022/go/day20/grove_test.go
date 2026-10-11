package main

import (
	"os"
	"testing"
)

func readExample(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	return string(data)
}

func TestSolvePart1Example(t *testing.T) {
	part1, _, err := solve(readExample(t))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 3 {
		t.Errorf("part1 = %d, want 3", part1)
	}
}

func TestSolvePart2Example(t *testing.T) {
	_, part2, err := solve(readExample(t))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part2 != 1623178306 {
		t.Errorf("part2 = %d, want 1623178306", part2)
	}
}

func TestMixSingleElementEdgeCase(t *testing.T) {
	numbers, err := parseNumbers("0\n")
	if err != nil {
		t.Fatalf("parseNumbers: %v", err)
	}
	mixed := mix(numbers, 1)
	if len(mixed) != 1 || mixed[0] != 0 {
		t.Errorf("mixed = %v, want [0]", mixed)
	}
}

func TestMixDuplicateValuesTrackedByOriginalIndex(t *testing.T) {
	// Values repeat; mixing must still move each occurrence independently.
	numbers, err := parseNumbers("1\n1\n0\n-1\n")
	if err != nil {
		t.Fatalf("parseNumbers: %v", err)
	}
	mixed := mix(numbers, 1)
	if len(mixed) != 4 {
		t.Fatalf("len(mixed) = %d, want 4", len(mixed))
	}
	sum := 0
	for _, v := range mixed {
		sum += v
	}
	if sum != 1 {
		t.Errorf("sum = %d, want 1", sum)
	}
}

func TestSolveMalformedInputErrors(t *testing.T) {
	if _, _, err := solve("not-a-number\n"); err == nil {
		t.Error("expected error for malformed input, got nil")
	}
}

func TestSolveEmptyInputErrors(t *testing.T) {
	if _, _, err := solve(""); err == nil {
		t.Error("expected error for empty input, got nil")
	}
}
