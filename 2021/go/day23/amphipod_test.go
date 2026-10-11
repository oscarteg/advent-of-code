package main

import (
	"os"
	"strings"
	"testing"
)

const (
	organized  = "#############\n#...........#\n###A#B#C#D###\n  #A#B#C#D#\n  #########\n"
	oneRoomRow = "#############\n#...........#\n###A#B#C#D###\n  #########\n"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 12521 || part2 != 44169 {
		t.Fatalf("got (%d, %d, %v), want (12521, 44169, nil)", part1, part2, err)
	}
}

func TestLeastEnergy(t *testing.T) {
	tests := []struct {
		name, input string
		want        int
	}{
		{"already organized", organized, 0},
		{"windows newlines", strings.ReplaceAll(organized, "\n", "\r\n"), 0},
		// Amber leaves room B for hallway cell 5 (2 steps), Bronze crosses into room B
		// along its shortest path (4 steps) and Amber walks on into room A (4 steps):
		// 6 energy for Amber, 40 for Bronze.
		{"swapped rooms", "#############\n#...........#\n###B#A#C#D###\n  #########\n", 46},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, err := parseBurrow(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			got, err := leastEnergy(start)

			if err != nil || got != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", got, err, tc.want)
			}
		})
	}
}

func TestOnlyTwoRoomRowsUnfold(t *testing.T) {
	if _, err := parseBurrow(oneRoomRow); err != nil {
		t.Fatalf("parseBurrow = %v, want nil", err)
	}

	if _, _, err := solve(oneRoomRow); err == nil {
		t.Fatal("solve = nil error, want an error")
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty", ""},
		{"only a wall row", "#############\n"},
		{"occupied hallway", "#############\n#..B........#\n###A#B#C#D###\n  #A#B#C#D#\n  #########\n"},
		{"unknown amphipod", "#############\n#...........#\n###A#B#X#D###\n  #A#B#C#D#\n  #########\n"},
		{"short room row", "#############\n#...........#\n###A#B#C###\n  #A#B#C#D#\n  #########\n"},
		{"unbalanced kinds", "#############\n#...........#\n###A#B#C#D###\n  #A#B#C#A#\n  #########\n"},
		{"no room rows", "#############\n#...........#\n  #########\n"},
		{"too many room rows", "#############\n#...........#\n###A#B#C#D###\n  #A#B#C#D#\n  #A#B#C#D#\n  #A#B#C#D#\n  #A#B#C#D#\n  #########\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want an error", tc.input)
			}
		})
	}
}
