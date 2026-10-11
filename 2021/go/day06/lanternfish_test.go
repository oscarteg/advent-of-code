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

	if err != nil || part1 != 5934 || part2 != 26984457539 {
		t.Fatalf("got (%d, %d, %v), want (5934, 26984457539, nil)", part1, part2, err)
	}
}

func TestSimulateBoundaries(t *testing.T) {
	tests := []struct {
		name  string
		input string
		days  int
		want  uint64
	}{
		{"single fish about to spawn, one day", "0", 1, 2},
		{"single fish about to spawn, ten days", "0", 10, 4},
		{"zero days returns input count", "3,4,3,1,2", 0, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			timers, err := parseTimers(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := simulate(timers, tc.days); got != tc.want {
				t.Fatalf("simulate(%q, %d) = %d, want %d", tc.input, tc.days, got, tc.want)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"non-numeric timer", "3,4,x,1,2"},
		{"out of range timer", "3,4,9,1,2"},
		{"empty input", ""},
		{"negative timer", "3,-1,2"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want error", tc.input)
			}
		})
	}
}
