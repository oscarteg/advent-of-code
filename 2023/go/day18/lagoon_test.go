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

	if err != nil || part1 != 62 || part2 != 952408144115 {
		t.Fatalf("got (%d, %d, %v), want (62, 952408144115, nil)", part1, part2, err)
	}
}

func TestLagoonArea(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		useColour bool
		want      int64
	}{
		{"4x4 square", "R 4 (#000000)\nD 4 (#000000)\nL 4 (#000000)\nU 4 (#000000)", false, 25},
		{"1x1 square", "R 1 (#000000)\nD 1 (#000000)\nL 1 (#000000)\nU 1 (#000000)", false, 4},
		{"decoded 4x4 square", "R 0 (#000040)\nD 0 (#000041)\nL 0 (#000042)\nU 0 (#000043)", true, 25},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			steps, err := parseDigPlan(tc.input, tc.useColour)
			if err != nil {
				t.Fatal(err)
			}
			if got := lagoonArea(steps); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	for _, input := range []string{
		"",
		"X 4 (#000000)",
		"R x (#000000)",
		"R 4",
		"RR 4 (#000000)",
	} {
		if _, _, err := solve(input); err == nil {
			t.Fatalf("solve(%q): want error, got nil", input)
		}
	}

	for _, input := range []string{
		"R 4 (#000000",
		"R 4 (#00000)",
		"R 4 (#00000g)",
		"R 4 000000",
	} {
		if _, err := parseDigPlan(input, true); err == nil {
			t.Fatalf("parseDigPlan(%q, true): want error, got nil", input)
		}
	}
}
