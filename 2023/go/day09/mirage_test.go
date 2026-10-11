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

	if err != nil || part1 != 114 || part2 != 2 {
		t.Fatalf("got (%d, %d, %v), want (114, 2, nil)", part1, part2, err)
	}
}

func TestExtrapolation(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		next, prev int64
		wantError  bool
	}{
		{"single history", "1 1 1", 1, 1, false},
		{"linear history", "0 3 6 9 12 15", 18, -3, false},
		{"single value", "5", 5, 5, false},
		{"blank lines ignored", "1 1 1\n\n2 2 2", 3, 3, false},
		{"empty input", "", 0, 0, true},
		{"malformed value", "1 2 x", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			next, prev, err := solve(tc.input)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (next != tc.next || prev != tc.prev) {
				t.Fatalf("got (%d, %d), want (%d, %d)", next, prev, tc.next, tc.prev)
			}
		})
	}
}
