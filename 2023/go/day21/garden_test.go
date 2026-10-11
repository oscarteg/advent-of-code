package main

import (
	"os"
	"strings"
	"testing"
)

// structuredMap is square with an open border and an open row and column
// through the start, the structure the real input has, so its tile-aligned
// counts are exactly quadratic.
const structuredMap = `...........
.#....#....
..#....#...
...#....#..
....#....#.
.....S.....
.#....#....
..#....#...
...#....#..
....#....#.
...........
`

const openMap = `.....
.....
..S..
.....
.....
`

func exampleMap(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	return string(input)
}

func mustParse(t *testing.T, input string) *garden {
	t.Helper()
	g, err := parseGarden(input)
	if err != nil {
		t.Fatalf("parse garden: %v", err)
	}
	return g
}

func TestSolveWithExampleStepCounts(t *testing.T) {
	part1, part2, err := solveWith(exampleMap(t), 6, 10)
	if err != nil {
		t.Fatalf("solveWith: %v", err)
	}
	if part1 != 16 {
		t.Errorf("part 1: got %d, want 16", part1)
	}
	if part2 != 50 {
		t.Errorf("part 2: got %d, want 50", part2)
	}
}

func TestSolveUsesThePuzzleStepCounts(t *testing.T) {
	// Without rocks every plot of the diamond with the right parity is
	// reachable: 13 of the 25 plots within the 5x5 map, and (steps+1)^2 on the
	// tiled map.
	part1, part2, err := solve(openMap)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 13 {
		t.Errorf("part 1 at %d steps: got %d, want 13", part1Steps, part1)
	}
	if want := (part2Steps + 1) * (part2Steps + 1); part2 != want {
		t.Errorf("part 2 at %d steps: got %d, want %d", part2Steps, part2, want)
	}
}

func TestReachable(t *testing.T) {
	tests := []struct {
		name  string
		steps int
		want  int
	}{
		{name: "start only", steps: 0, want: 1},
		{name: "one step", steps: 1, want: 2},
		{name: "two steps", steps: 2, want: 4},
		{name: "three steps", steps: 3, want: 6},
		{name: "official example", steps: 6, want: 16},
		{name: "saturated well before 64 steps", steps: part1Steps, want: 42},
	}
	g := mustParse(t, exampleMap(t))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := g.reachable(tt.steps); got != tt.want {
				t.Errorf("reachable(%d): got %d, want %d", tt.steps, got, tt.want)
			}
		})
	}
}

func TestReachableTiled(t *testing.T) {
	tests := []struct {
		steps int
		want  int
	}{
		{steps: 6, want: 16},
		{steps: 10, want: 50},
		{steps: 50, want: 1594},
		{steps: 100, want: 6536},
		{steps: 500, want: 167004},
	}
	g := mustParse(t, exampleMap(t))
	for _, tt := range tests {
		if got := g.reachableTiled(tt.steps); got != tt.want {
			t.Errorf("reachableTiled(%d): got %d, want %d", tt.steps, got, tt.want)
		}
	}
}

func TestReachableExtrapolatedMatchesTheFlood(t *testing.T) {
	// The quadratic is fitted through 5, 16 and 27 steps, so these counts are
	// predictions rather than fitted points.
	tests := []struct {
		steps int
		want  int
	}{
		{steps: 38, want: 1309},
		{steps: 49, want: 2156},
		{steps: 60, want: 3213},
	}
	g := mustParse(t, structuredMap)
	for _, tt := range tests {
		if got := g.reachableTiled(tt.steps); got != tt.want {
			t.Errorf("reachableTiled(%d): got %d, want %d", tt.steps, got, tt.want)
		}
		got, err := g.reachableExtrapolated(tt.steps)
		if err != nil {
			t.Fatalf("reachableExtrapolated(%d): %v", tt.steps, err)
		}
		if got != tt.want {
			t.Errorf("reachableExtrapolated(%d): got %d, want %d", tt.steps, got, tt.want)
		}
	}
}

func TestReachableExtrapolatedFloodsSmallStepCounts(t *testing.T) {
	g := mustParse(t, exampleMap(t))
	got, err := g.reachableExtrapolated(6)
	if err != nil {
		t.Fatalf("reachableExtrapolated(6): %v", err)
	}
	if got != 16 {
		t.Errorf("reachableExtrapolated(6): got %d, want 16", got)
	}
}

func TestWalledInStartCannotStep(t *testing.T) {
	g := mustParse(t, "###\n#S#\n###")
	if got := g.reachable(0); got != 1 {
		t.Errorf("reachable(0): got %d, want 1", got)
	}
	if got := g.reachableTiled(6); got != 0 {
		t.Errorf("reachableTiled(6): got %d, want 0", got)
	}
}

func TestExtrapolationNeedsASquareMap(t *testing.T) {
	_, _, err := solve(".....\n..S..")
	if err == nil {
		t.Fatal("solve: got no error for a non-square map")
	}
	if !strings.Contains(err.Error(), "square") {
		t.Errorf("solve: got %v, want an error about a square map", err)
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "blank lines only", input: "\n\n"},
		{name: "no start", input: "...\n...\n..."},
		{name: "two starts", input: ".S.\n...\n.S."},
		{name: "ragged rows", input: "...\n..\n..."},
		{name: "unknown tile", input: "...\n.X.\n..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := solve(tt.input); err == nil {
				t.Errorf("solve(%q): got no error", tt.input)
			}
		})
	}
}
