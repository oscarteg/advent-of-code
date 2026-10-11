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

	if err != nil || part1 != 2 || part2 != 4 {
		t.Fatalf("got (%d, %d, %v), want (2, 4, nil)", part1, part2, err)
	}
}

func TestReports(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantError    bool
	}{
		{"single level is safe for both parts", "5\n", 1, 1, false},
		{"windows lines", "7 6 4 2 1\r\n1 2 7 8 9\r\n", 1, 1, false},
		{"malformed level", "1 two 3\n", 0, 0, true},
		{"two unsafe levels dampened to one", "10 4\n", 0, 1, false},
		{"whitespace-only line", "   \n", 0, 0, true},
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
