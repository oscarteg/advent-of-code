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
	if err != nil || part1 != 13 || part2 != 30 {
		t.Fatalf("got (%d, %d, %v), want (13, 30, nil)", part1, part2, err)
	}
}

func TestCardMatches(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantError    bool
	}{
		{"single card no match", "Card 1: 1 2 3 | 4 5 6", 0, 1, false},
		{"single card one match", "Card 1: 1 2 3 | 3 5 6", 1, 1, false},
		{"matches do not cascade past the last card", "Card 1: 1 2 | 1 2\nCard 2: 1 2 | 3 4", 2, 3, false},
		{"windows line endings", "Card 1: 1 2 | 1 2\r\nCard 2: 1 2 | 3 4", 2, 3, false},
		{"empty input", "", 0, 0, true},
		{"missing colon", "Card 1 1 2 | 3 4", 0, 0, true},
		{"missing pipe", "Card 1: 1 2 3 4", 0, 0, true},
		{"bad winning number", "Card 1: x 2 | 3 4", 0, 0, true},
		{"bad your number", "Card 1: 1 2 | x 4", 0, 0, true},
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
