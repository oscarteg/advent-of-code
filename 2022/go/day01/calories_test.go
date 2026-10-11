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

	if err != nil || part1 != 24000 || part2 != 45000 {
		t.Fatalf("got (%d, %d, %v), want (24000, 45000, nil)", part1, part2, err)
	}
}

func TestCalorieCounting(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantError    bool
	}{
		{"single elf", "100\n200", 300, 300, false},
		{"fewer than three elves", "100\n\n200", 200, 300, false},
		{"malformed calorie value", "1000\nnot-a-number\n\n2000", 0, 0, true},
		{"empty input", "", 0, 0, true},
		{"blank line only", "\n\n", 0, 0, true},
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
