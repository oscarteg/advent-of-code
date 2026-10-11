package main

import (
	"os"
	"testing"
)

func TestTinyExample(t *testing.T) {
	input := "1,1,1\n2,1,1"
	part1, _, err := solve(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if part1 != 10 {
		t.Errorf("expected part1=10, got %d", part1)
	}
}

func TestExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if part1 != 64 {
		t.Errorf("expected part1=64, got %d", part1)
	}
	if part2 != 58 {
		t.Errorf("expected part2=58, got %d", part2)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []string{
		"1,1",
		"1,1,a",
		"",
	}

	for _, input := range tests {
		_, _, err := solve(input)
		if err == nil {
			t.Errorf("expected error for input %q", input)
		}
	}
}
