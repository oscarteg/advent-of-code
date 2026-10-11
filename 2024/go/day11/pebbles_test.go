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
	if err != nil || part1 != 55312 || part2 != 65601038650482 {
		t.Fatalf("got (%d, %d, %v), want (55312, 65601038650482, nil)", part1, part2, err)
	}
}

func TestSolveWith(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		blinks1   int
		blinks2   int
		part1     int64
		wantError bool
	}{
		{"example after six blinks", "125 17\n", 6, 6, 22, false},
		{"single zero stone becomes one", "0\n", 1, 1, 1, false},
		{"empty input", "", 1, 1, 0, true},
		{"non-numeric stone", "abc\n", 1, 1, 0, true},
		{"negative stone", "-5\n", 1, 1, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, _, err := solveWith(tc.input, tc.blinks1, tc.blinks2)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && part1 != tc.part1 {
				t.Fatalf("got part1 = %d, want %d", part1, tc.part1)
			}
		})
	}
}
