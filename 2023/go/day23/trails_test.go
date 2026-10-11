package main

import (
	"os"
	"strings"
	"testing"
)

const (
	corridor      = "#.#\n#.#\n#.#\n"
	downhillSlope = "#.#\n#v#\n#.#\n"
	uphillSlope   = "#.#\n#^#\n#.#\n"
	deadEndSpur   = "#.##\n#..#\n#.##\n"
	// Two ways around one loop: three steps along the left side, five along the
	// right, whose last tile is a slope pointing back west.
	loopWithSlope = "#.###\n#.<.#\n#.#.#\n#...#\n##.##\n"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 94 || part2 != 154 {
		t.Fatalf("got (%d, %d, %v), want (94, 154, nil)", part1, part2, err)
	}
}

func TestSmallMaps(t *testing.T) {
	tests := []struct {
		name, input  string
		part1, part2 int
	}{
		{"straight corridor", corridor, 2, 2},
		{"downhill slope", downhillSlope, 2, 2},
		{"dead end spur", deadEndSpur, 2, 2},
		{"slope blocks the long way round", loopWithSlope, 5, 7},
		{"entrance is the exit", "#.#\n", 0, 0},
		{"windows line endings", strings.ReplaceAll(corridor, "\n", "\r\n"), 2, 2},
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

func TestUphillSlopeBlocksTheOnlyRoute(t *testing.T) {
	trails, err := parseMap(uphillSlope)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := trails.longestHike(icy); err == nil {
		t.Fatal("longestHike(icy) = nil error, want an error")
	}
	steps, err := trails.longestHike(dry)
	if err != nil || steps != 2 {
		t.Fatalf("longestHike(dry) = (%d, %v), want (2, nil)", steps, err)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty", ""},
		{"ragged rows", "#.#\n#.##\n#.#\n"},
		{"unknown tile", "#.#\n#x#\n#.#\n"},
		{"no entrance", "###\n#.#\n#.#\n"},
		{"no exit", "#.#\n#.#\n###\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("solve(%q) = nil error, want an error", tc.input)
			}
		})
	}
}
