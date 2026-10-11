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

	if err != nil || part1 != 21 || part2 != 40 {
		t.Fatalf("got (%d, %d, %v), want (21, 40, nil)", part1, part2, err)
	}
}

func TestBeamSimulation(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 uint64
		wantError    bool
	}{
		{"no splitters leaves a single straight beam", ".S.\n...\n...\n", 0, 1, false},
		{"a splitter at the grid edge loses both offshoots", "S\n^\n", 1, 0, false},
		{"start on the bottom row never splits", "S\n", 0, 1, false},
		{"missing start", "...\n^..\n", 0, 0, true},
		{"more than one start", "S.S\n...\n", 0, 0, true},
		{"inconsistent row width", "S..\n^.\n", 0, 0, true},
		{"invalid character", "S.x\n^..\n", 0, 0, true},
		{"empty input", "   \n  ", 0, 0, true},
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
