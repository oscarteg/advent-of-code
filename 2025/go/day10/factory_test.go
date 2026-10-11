package main

import (
	"errors"
	"os"
	"testing"
)

func TestSolveOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 7 {
		t.Errorf("part1 = %d, want 7", part1)
	}
	if part2 != 33 {
		t.Errorf("part2 = %d, want 33", part2)
	}
}

func TestSolveTableDriven(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
	}{
		{
			name:  "first official machine",
			input: "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			part1: 2,
			part2: 10,
		},
		{
			name:  "second official machine",
			input: "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			part1: 3,
			part2: 12,
		},
		{
			name:  "third official machine",
			input: "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			part1: 2,
			part2: 11,
		},
		{
			name:  "configured machine with zero joltage needs no presses",
			input: "[..] (0) (1) {0,0}\n",
			part1: 0,
			part2: 0,
		},
		{
			name:  "only button pressed once for lights and three times for joltage",
			input: "[##] (0,1) {3,3}\n",
			part1: 1,
			part2: 3,
		},
		{
			name:  "surrounding whitespace and blank lines are ignored",
			input: "\n  [##] (0,1) {3,3}  \n\n",
			part1: 1,
			part2: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part1 != tt.part1 {
				t.Errorf("part1 = %d, want %d", part1, tt.part1)
			}
			if part2 != tt.part2 {
				t.Errorf("part2 = %d, want %d", part2, tt.part2)
			}
		})
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "missing diagram", input: "(0) (1) {1,1}\n"},
		{name: "missing joltage requirements", input: "[.#] (0) (1)\n"},
		{name: "missing buttons", input: "[.#] {1,1}\n"},
		{name: "unparseable light index", input: "[.#] (a) {1,1}\n"},
		{name: "unknown light state", input: "[.x] (0) {1,1}\n"},
		{name: "button wired past the last light", input: "[.#] (5) {1,1}\n"},
		{name: "negative joltage", input: "[.#] (0) {1,-1}\n"},
		{name: "joltage count does not match the lights", input: "[.#] (0) {1}\n"},
		{name: "no machines", input: "\n  \n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := solve(tt.input); err == nil {
				t.Errorf("solve(%q) = nil error, want error", tt.input)
			}
		})
	}
}

func TestSolveRejectsUnreachableMachines(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "lights the buttons cannot reach", input: "[#.] (0,1) {1,1}\n"},
		{name: "joltage the buttons cannot reach", input: "[#.] (0) {1,2}\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := solve(tt.input); !errors.Is(err, errUnreachable) {
				t.Errorf("solve(%q) error = %v, want %v", tt.input, err, errUnreachable)
			}
		})
	}
}
