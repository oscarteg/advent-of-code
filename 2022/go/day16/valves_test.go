package main

import (
	"os"
	"testing"
)

func readExample(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	return string(input)
}

func TestOfficialExample(t *testing.T) {
	part1, part2, err := solve(readExample(t))
	if err != nil || part1 != 1651 || part2 != 1707 {
		t.Fatalf("got (%d, %d, %v), want (1651, 1707, nil)", part1, part2, err)
	}
}

func TestValveCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		part1     int
		part2     int
		wantError bool
	}{
		{
			// A minute to walk to BB and a minute to open it leaves 28 of the 30
			// minutes at flow 10, and 24 of the 26 when sharing the work.
			name: "lone valve is opened as early as possible",
			input: "Valve AA has flow rate=0; tunnel leads to valve BB\n" +
				"Valve BB has flow rate=10; tunnel leads to valve AA\n",
			part1: 280,
			part2: 240,
		},
		{
			name: "valves without flow release nothing",
			input: "Valve AA has flow rate=0; tunnels lead to valves BB, CC\n" +
				"Valve BB has flow rate=0; tunnel leads to valve AA\n" +
				"Valve CC has flow rate=0; tunnel leads to valve AA\n",
		},
		{
			name: "valves without a tunnel route stay closed",
			input: "Valve AA has flow rate=0; tunnel leads to valve BB\n" +
				"Valve BB has flow rate=0; tunnel leads to valve AA\n" +
				"Valve CC has flow rate=5; tunnel leads to valve DD\n" +
				"Valve DD has flow rate=0; tunnel leads to valve CC\n",
		},
		{name: "no valves", input: "\n\n", wantError: true},
		{name: "truncated valve", input: "Valve AA has flow rate=0\n", wantError: true},
		{name: "non-numeric flow rate", input: "Valve AA has flow rate=x; tunnel leads to valve BB\n", wantError: true},
		{name: "malformed tunnel list", input: "Valve AA has flow rate=0; leads to BB\n", wantError: true},
		{name: "missing valve prefix", input: "AA has flow rate=0; tunnel leads to valve AA\n", wantError: true},
		{name: "unknown valve in tunnel list", input: "Valve AA has flow rate=0; tunnel leads to valve ZZ\n", wantError: true},
		{name: "no start valve", input: "Valve BB has flow rate=1; tunnel leads to valve BB\n", wantError: true},
		{
			name: "valve listed twice",
			input: "Valve AA has flow rate=0; tunnel leads to valve AA\n" +
				"Valve AA has flow rate=1; tunnel leads to valve AA\n",
			wantError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
