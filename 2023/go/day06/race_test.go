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

	if err != nil || part1 != 288 || part2 != 71503 {
		t.Fatalf("got (%d, %d, %v), want (288, 71503, nil)", part1, part2, err)
	}
}

func TestWinningHoldTimes(t *testing.T) {
	tests := []struct {
		name     string
		r        race
		expected uint64
	}{
		{"time 7 distance 9", race{time: 7, distance: 9}, 4},
		{"time 15 distance 40", race{time: 15, distance: 40}, 8},
		{"time 30 distance 200", race{time: 30, distance: 200}, 9},
		{"unbeatable record", race{time: 1, distance: 100}, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.r.winningHoldTimes(); got != tc.expected {
				t.Fatalf("got %d, want %d", got, tc.expected)
			}
		})
	}
}

func TestJoinedDigitsIgnoreSpaces(t *testing.T) {
	r, err := parseSingleRace("Time: 1 2\nDistance: 3 4")
	if err != nil {
		t.Fatal(err)
	}
	if r.time != 12 || r.distance != 34 {
		t.Fatalf("got %+v, want time=12 distance=34", r)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"missing distance line", "Time: 1"},
		{"mismatched columns", "Time: 1 2\nDistance: 1"},
		{"non-numeric time", "Time: x\nDistance: 1"},
		{"reversed sections", "Distance: 1\nTime: 1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("expected error for input %q", tc.input)
			}
		})
	}
}
