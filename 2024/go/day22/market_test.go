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

	part1, _, err := solve(string(input))
	if err != nil || part1 != 37327623 {
		t.Fatalf("got (%d, %v), want (37327623, nil)", part1, err)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int64
		wantError    bool
	}{
		{"part two best change sequence", "1\n2\n3\n2024\n", 0, 23, false},
		{"single seed matches the official per-seed example", "1\n", 8685429, 0, false},
		{"empty input", "", 0, 0, true},
		{"non-numeric seed", "not-a-number\n", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if tc.wantError {
				return
			}
			if tc.name == "single seed matches the official per-seed example" && part1 != tc.part1 {
				t.Fatalf("part1 = %d, want %d", part1, tc.part1)
			}
			if tc.name == "part two best change sequence" && part2 != tc.part2 {
				t.Fatalf("part2 = %d, want %d", part2, tc.part2)
			}
		})
	}
}
