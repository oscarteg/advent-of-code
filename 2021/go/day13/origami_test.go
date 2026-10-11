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

	wantPart2 := "#####\n#...#\n#...#\n#...#\n#####"
	if err != nil || part1 != 17 || part2 != wantPart2 {
		t.Fatalf("got (%d, %q, %v), want (17, %q, nil)", part1, part2, err, wantPart2)
	}
}

func TestOrigamiBoundaries(t *testing.T) {
	tests := []struct {
		name, input string
		part1       int
		part2       string
		wantError   bool
	}{
		{"single fold", "0,0\n3,0\n\nfold along x=2", 2, "##", false},
		{"overlapping dots count once", "0,0\n2,0\n1,0\n\nfold along x=1", 2, "##", false},
		{"missing folds", "0,0\n1,1", 0, "", true},
		{"malformed dot", "not-a-dot\n\nfold along x=1", 0, "", true},
		{"unknown axis", "0,0\n\nfold along z=1", 0, "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %q), want (%d, %q)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
