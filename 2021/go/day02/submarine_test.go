package main

import (
	"os"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 150 || part2 != 900 {
		t.Fatalf("got (%d, %d, %v), want (150, 900, nil)", part1, part2, err)
	}
}

func TestUpClampsAtZero(t *testing.T) {
	part1, part2, err := solve("up 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if part1 != 0 || part2 != 0 {
		t.Fatalf("got (%d, %d), want (0, 0)", part1, part2)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unknown direction", "sideways 5"},
		{"missing distance", "forward"},
		{"non-numeric distance", "forward abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := solve(tt.input); err == nil {
				t.Fatalf("solve(%q): expected error, got nil", tt.input)
			}
		})
	}
}
