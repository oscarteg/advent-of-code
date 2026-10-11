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

	if err != nil || part1 != 15 || part2 != 1134 {
		t.Fatalf("got (%d, %d, %v), want (15, 1134, nil)", part1, part2, err)
	}
}

func TestHeightMapBoundaries(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 uint32
		wantError    bool
	}{
		{"flat grid with no low points is rejected", "11\n11", 0, 0, true},
		{"three small basins are multiplied", "109\n919\n091", 4, 3, false},
		{"ragged rows are rejected", "123\n12", 0, 0, true},
		{"invalid digit is rejected", "12a\n123", 0, 0, true},
		{"empty input is rejected", "", 0, 0, true},
		{"fewer than three basins is rejected", "12\n12", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || uint32(part2) != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
