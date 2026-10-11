package main

import (
	"os"
	"strings"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	first, last, err := solve(string(input))

	if err != nil || first != 4512 || last != 1924 {
		t.Fatalf("got (%d, %d, %v), want (4512, 1924, nil)", first, last, err)
	}
}

func TestBingoBoundaries(t *testing.T) {
	const grid = "0 1 2 3 4\n5 6 7 8 9\n10 11 12 13 14\n15 16 17 18 19\n20 21 22 23 24"
	tests := []struct {
		name, input string
		score       uint64
		wantError   bool
	}{
		{"repeated draws", "0,0,1,2,3,4\n\n" + grid, 1160, false},
		{"column", "0,5,10,15,20\n\n" + grid, 5000, false},
		{"diagonal", "0,6,12,18,24\n\n" + grid, 0, true},
		{"short board", "1,2,3,4,5\n\n1 2 3 4 5", 0, true},
		{"bad draw", "x\n\n" + grid, 0, true},
		{"empty draws", "\n\n" + grid, 0, true},
		{"empty draw field", "1,,2\n\n" + grid, 0, true},
		{"bad board number", "0,1,2,3,4\n\n1 2 x 4 5\n6 7 8 9 10\n11 12 13 14 15\n16 17 18 19 20\n21 22 23 24 25", 0, true},
		{"draws end before win", "0\n\n" + grid, 0, true},
		{"windows lines", strings.ReplaceAll("0,1,2,3,4\n\n"+grid, "\n", "\r\n"), 1160, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			first, last, err := solve(tc.input)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (first != tc.score || last != tc.score) {
				t.Fatalf("got (%d, %d), want (%d, %d)", first, last, tc.score, tc.score)
			}
		})
	}
}
