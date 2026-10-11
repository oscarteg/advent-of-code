package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

const corridor = "#####\n#S#E#\n#.#.#\n#...#\n#####\n"

func exampleMap(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	return string(input)
}

func TestOfficialExampleTwoPicosecondCheats(t *testing.T) {
	input := exampleMap(t)
	tests := []struct {
		minSaving, want int
	}{
		{64, 1},
		{40, 2},
		{20, 5},
		{2, 44},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("saves at least %d", tc.minSaving), func(t *testing.T) {
			part1, _, err := solveWith(input, tc.minSaving)
			if err != nil {
				t.Fatalf("solveWith error = %v", err)
			}
			if part1 != tc.want {
				t.Errorf("part 1 = %d, want %d", part1, tc.want)
			}
		})
	}
}

func TestOfficialExampleTwentyPicosecondCheats(t *testing.T) {
	input := exampleMap(t)
	tests := []struct {
		minSaving, want int
	}{
		{76, 3},
		{74, 7},
		{72, 29},
		{50, 285},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("saves at least %d", tc.minSaving), func(t *testing.T) {
			_, part2, err := solveWith(input, tc.minSaving)
			if err != nil {
				t.Fatalf("solveWith error = %v", err)
			}
			if part2 != tc.want {
				t.Errorf("part 2 = %d, want %d", part2, tc.want)
			}
		})
	}
}

func TestRealThresholdFindsNoCheatInTheExample(t *testing.T) {
	part1, part2, err := solve(exampleMap(t))

	if err != nil || part1 != 0 || part2 != 0 {
		t.Fatalf("got (%d, %d, %v), want (0, 0, nil)", part1, part2, err)
	}
}

func TestCorridorCheats(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		minSaving    int
		part1, part2 int
	}{
		{"every cheat that saves time", corridor, 2, 2, 4},
		{"only the best cheats", corridor, 4, 1, 1},
		{"windows line endings", strings.ReplaceAll(corridor, "\n", "\r\n"), 2, 2, 4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solveWith(tc.input, tc.minSaving)
			if err != nil {
				t.Fatalf("solveWith error = %v", err)
			}
			if part1 != tc.part1 || part2 != tc.part2 {
				t.Errorf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"empty input", ""},
		{"no start or end", "###\n#.#\n###\n"},
		{"no end", "###\n#S#\n###\n"},
		{"two starts", "#####\n#SSE#\n#####\n"},
		{"rows of different widths", "#####\n#S.E#\n####\n"},
		{"unknown map character", "#####\n#S?E#\n#####\n"},
		{"end walled off from the start", "#####\n#S#E#\n#####\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := solve(tc.input)

			if !errors.Is(err, errInvalidTrack) {
				t.Fatalf("error = %v, want one wrapping %v", err, errInvalidTrack)
			}
		})
	}
}
