package main

import (
	"os"
	"strings"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 5 || part2 != 7 {
		t.Fatalf("got (%d, %d, %v), want (5, 7, nil)", part1, part2, err)
	}
}

func TestSettledStacks(t *testing.T) {
	const column = "0,0,1~0,0,1\n0,0,2~0,0,2\n0,0,3~0,0,3"
	const sharedLoad = "0,0,1~0,0,1\n1,0,1~1,0,1\n0,0,2~1,0,2"
	tests := []struct {
		name, input  string
		part1, part2 int
	}{
		{"single cube", "0,0,1~0,0,1", 1, 0},
		{"cube falls to the ground", "0,0,9~0,0,9", 1, 0},
		{"column", column, 1, 3},
		{"shared load", sharedLoad, 3, 0},
		{"rests on the highest brick", "0,0,1~0,0,1\n1,0,1~1,0,2\n0,0,3~1,0,3", 2, 1},
		{"cascade stops at a second supporter", sharedLoad + "\n0,0,3~0,0,3", 3, 1},
		{"windows line endings", strings.ReplaceAll(column, "\n", "\r\n"), 1, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if err != nil || part1 != tc.part1 || part2 != tc.part2 {
				t.Fatalf("got (%d, %d, %v), want (%d, %d, nil)", part1, part2, err, tc.part1, tc.part2)
			}
		})
	}
}

func TestSnapshotOrderIsIrrelevant(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Fields(strings.TrimSpace(string(input)))
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	part1, part2, err := solve(strings.Join(lines, "\n"))

	if err != nil || part1 != 5 || part2 != 7 {
		t.Fatalf("got (%d, %d, %v), want (5, 7, nil)", part1, part2, err)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty", ""},
		{"missing separator", "1,0,1"},
		{"short corner", "1,0,1~1,0"},
		{"long corner", "1,0,1~1,0,1,1"},
		{"on the ground", "1,0,0~1,2,0"},
		{"not a straight line", "0,0,1~1,1,1"},
		{"not a number", "x,0,1~1,0,1"},
		{"negative coordinate", "0,0,1~0,0,-1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want an error", tc.input)
			}
		})
	}
}
