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

	if err != nil || part1 != 8 || part2 != 2286 {
		t.Fatalf("got (%d, %d, %v), want (8, 2286, nil)", part1, part2, err)
	}
}

func TestGamePossibility(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{"within bag", "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green", true},
		{"exceeds red", "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green", false},
		{"exact limits", "Game 9: 12 red, 13 green, 14 blue", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g, err := parseGame(tc.line)
			if err != nil {
				t.Fatal(err)
			}
			if got := g.isPossible(bagRed, bagGreen, bagBlue); got != tc.want {
				t.Fatalf("isPossible() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGamePower(t *testing.T) {
	tests := []struct {
		name string
		line string
		want uint64
	}{
		{"three colors", "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green", 48},
		{"single color", "Game 1: 3 red; 5 red", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g, err := parseGame(tc.line)
			if err != nil {
				t.Fatal(err)
			}
			if got := g.power(); got != tc.want {
				t.Fatalf("power() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []string{
		"",
		"Game 1",
		"Game: 3 red",
		"Foo 1: 3 red",
		"Game x: 3 red",
		"Game 1: 3 red; ",
		"Game 1: three red",
		"Game 1: 3 purple",
		"Game 1: 3",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, _, err := solve(input); err == nil {
				t.Fatalf("solve(%q) succeeded, want error", input)
			}
		})
	}
}
