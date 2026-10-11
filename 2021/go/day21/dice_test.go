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

	if err != nil || part1 != 739785 || part2 != 444356092776315 {
		t.Fatalf("got (%d, %d, %v), want (739785, 444356092776315, nil)", part1, part2, err)
	}
}

func TestAdvanceWrapsAroundTheTrack(t *testing.T) {
	tests := []struct {
		position uint8
		roll     uint32
		want     uint8
	}{
		{7, 5, 2},
		{10, 1, 1},
		{4, 6, 10},
	}
	for _, tc := range tests {
		if got := advance(tc.position, tc.roll); got != tc.want {
			t.Errorf("advance(%d, %d) = %d, want %d", tc.position, tc.roll, got, tc.want)
		}
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"not a game", "not a valid game"},
		{"missing second player", "Player 1 starting position: 4\n"},
		{"out of range position", "Player 1 starting position: 0\nPlayer 2 starting position: 8\n"},
		{"bad number", "Player 1 starting position: x\nPlayer 2 starting position: 8\n"},
		{"wrong player order", "Player 2 starting position: 4\nPlayer 1 starting position: 8\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("expected error for input %q", tc.input)
			}
		})
	}
}
