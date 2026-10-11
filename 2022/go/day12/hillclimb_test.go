package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("reading example: %v", err)
	}

	tests := []struct {
		name         string
		input        string
		part1, part2 int
		expectErr    bool
	}{
		{name: "official example", input: string(example), part1: 31, part2: 29},
		{name: "single row staircase", input: "SbcdefghijklmnopqrstuvwxyE", part1: 25, part2: 25},
		{name: "ragged rows are rejected", input: "Sab\nabcE", expectErr: true},
		{name: "missing start is rejected", input: "abz\nabE", expectErr: true},
		{name: "invalid elevation character is rejected", input: "S1\naE", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("solve(%q) = (%d, %d, nil), want error", tt.input, part1, part2)
				}
				return
			}
			if err != nil {
				t.Fatalf("solve(%q) returned unexpected error: %v", tt.input, err)
			}
			if part1 != tt.part1 || part2 != tt.part2 {
				t.Errorf("solve(%q) = (%d, %d), want (%d, %d)", tt.input, part1, part2, tt.part1, tt.part2)
			}
		})
	}
}
