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

	if err != nil || part1 != 13 || part2 != 43 {
		t.Fatalf("got (%d, %d, %v), want (13, 43, nil)", part1, part2, err)
	}
}

func TestPaperRolls(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantErr      bool
	}{
		{"isolated rolls are all accessible and removed in one pass", "@.@\n...\n@.@\n", 4, 4, false},
		{"fully packed block loses only its corners each round", "@@@@\n@@@@\n@@@@\n@@@@\n", 4, 4, false},
		{"single roll", "@", 1, 1, false},
		{"single empty cell", ".", 0, 0, false},
		{"inconsistent row width", "@@@\n@@\n", 0, 0, true},
		{"invalid character", "@.x\n@.@\n", 0, 0, true},
		{"empty input", "   \n  ", 0, 0, true},
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
