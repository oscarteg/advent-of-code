package main

import (
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
	if part1 != 18 {
		t.Errorf("part1 = %d, want 18", part1)
	}
	if part2 != 54 {
		t.Errorf("part2 = %d, want 54", part2)
	}
}

func TestSolveNoBlizzardsIsADirectWalk(t *testing.T) {
	input := "#.#\n#.#\n#.#\n#.#\n#.#\n"
	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 4 {
		t.Errorf("part1 = %d, want 4", part1)
	}
	if part2 <= part1 {
		t.Errorf("part2 = %d, want > part1 (%d)", part2, part1)
	}
}

func TestSolveMalformedInput(t *testing.T) {
	cases := map[string]string{
		"missing openings":   "###\n#.#\n#.#\n#.#\n###\n",
		"inconsistent width": "#.#\n#..#\n#.#\n#.#\n#.#\n",
		"unknown character":  "#.#\n#x#\n#.#\n#.#\n#.#\n",
		"too few lines":      "#.#\n#.#\n",
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, _, err := solve(input); err == nil {
				t.Errorf("solve(%q) = nil error, want error", name)
			}
		})
	}
}
