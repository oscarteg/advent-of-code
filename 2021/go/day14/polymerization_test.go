package main

import (
	"os"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 1588 || part2 != 2188189693529 {
		t.Errorf("got (%d, %d), want (1588, 2188189693529)", part1, part2)
	}
}

func TestMissingTemplate(t *testing.T) {
	input := "\n\nAB -> C"
	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestMissingRules(t *testing.T) {
	input := "ABC"
	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for missing rules")
	}
}

func TestInvalidRuleFormat(t *testing.T) {
	input := "ABC\n\nAB : C"
	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for invalid rule format")
	}
}

func TestEdgeCase(t *testing.T) {
	input := "AB\n\nAA -> X\nAB -> Y\nAX -> Z\nAY -> W\nAZ -> Q\nAW -> K\nYA -> M\nYB -> V"
	_, _, err := solve(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
