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
	if err != nil || part1 != 3749 || part2 != 11387 {
		t.Fatalf("got (%d, %d, %v), want (3749, 11387, nil)", part1, part2, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int64
		wantError    bool
	}{
		{"single operand is trivially solvable", "5: 5\n", 5, 5, false},
		{"unsolvable equation contributes nothing", "10: 2 2\n", 0, 0, false},
		{"concatenation only solvable in part two", "156: 15 6\n", 0, 156, false},
		{"empty input", "", 0, 0, true},
		{"missing colon", "190 10 19\n", 0, 0, true},
		{"no operands", "190:\n", 0, 0, true},
		{"non-numeric target", "abc: 1 2\n", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
