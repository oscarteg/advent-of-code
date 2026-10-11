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

	if err != nil || part1 != 1656 || part2 != 195 {
		t.Fatalf("got (%d, %d, %v), want (1656, 195, nil)", part1, part2, err)
	}
}

func TestOctopusBoundaries(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		input       string
		steps       int
		wantFlashes int
		wantError   bool
	}{
		{
			name:        "official example after ten steps",
			input:       string(example),
			steps:       10,
			wantFlashes: 204,
		},
		{
			name:      "not a grid",
			input:     "not a grid",
			wantError: true,
		},
		{
			name:      "wrong row count",
			input:     "123\n456",
			wantError: true,
		},
		{
			name:      "invalid digit",
			input:     "x483143223\n2745854711\n5264556173\n6141336146\n6357385478\n4167524645\n2176841721\n6882881134\n4846848554\n5283751526",
			wantError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flashes, _, err := solveWith(tc.input, tc.steps)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && flashes != tc.wantFlashes {
				t.Fatalf("got %d flashes, want %d", flashes, tc.wantFlashes)
			}
		})
	}
}
