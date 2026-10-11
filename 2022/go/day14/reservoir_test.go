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

	if err != nil || part1 != 24 || part2 != 93 {
		t.Fatalf("got (%d, %d, %v), want (24, 93, nil)", part1, part2, err)
	}
}

func TestReservoirBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		part1     uint64
		wantError bool
	}{
		{"flat ledge bounded pile", "498,3 -> 502,3", 4, false},
		{"missing comma", "498 4 -> 498 6", 0, true},
		{"diagonal segment", "498,4 -> 500,6", 0, true},
		{"empty input", "", 0, true},
		{"single point path", "498,4", 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, _, err := solve(tc.input)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && part1 != tc.part1 {
				t.Fatalf("got part1 = %d, want %d", part1, tc.part1)
			}
		})
	}
}
