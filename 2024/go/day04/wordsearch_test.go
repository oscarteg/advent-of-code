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

	if err != nil || part1 != 18 || part2 != 9 {
		t.Fatalf("got (%d, %d, %v), want (18, 9, nil)", part1, part2, err)
	}
}

func TestWordSearch(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantErr      bool
	}{
		{"single diagonal x-mas", "M.S\n.A.\nM.S\n", 0, 1, false},
		{"windows line endings", "M.S\r\n.A.\r\nM.S\r\n", 0, 1, false},
		{"ragged grid", "XMAS\nXM\n", 0, 0, true},
		{"empty input", "", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
