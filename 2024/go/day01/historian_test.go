package main

import (
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 11 {
		t.Errorf("part1 got %d, want 11", part1)
	}
	if part2 != 31 {
		t.Errorf("part2 got %d, want 31", part2)
	}
}

func TestSinglePair(t *testing.T) {
	part1, part2, err := solve("5 10\n")
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 5 {
		t.Errorf("part1 got %d, want 5", part1)
	}
	if part2 != 0 {
		t.Errorf("part2 got %d, want 0", part2)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []string{
		"5\n",
		"abc 123\n",
	}
	for _, input := range tests {
		_, _, err := solve(input)
		if err == nil {
			t.Errorf("expected error for input %q", input)
		}
	}
}
