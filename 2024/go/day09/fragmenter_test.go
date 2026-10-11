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
	if err != nil || part1 != 1928 || part2 != 2858 {
		t.Fatalf("got (%d, %d, %v), want (1928, 2858, nil)", part1, part2, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int64
		wantError    bool
	}{
		{"single file with no free space", "5", 0, 0, false},
		{"trailing free space is never filled", "12", 0, 0, false},
		{"single block moves into a leading gap", "112", 3, 5, false},
		{"empty input", "", 0, 0, true},
		{"non-numeric character", "12a45", 0, 0, true},
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
