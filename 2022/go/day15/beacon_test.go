package main

import (
	"os"
	"testing"
)

const (
	exampleRow         = 10
	exampleSearchLimit = 20
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
	part1, part2, err := solveWith(readExample(t), exampleRow, exampleSearchLimit)
	if err != nil || part1 != 26 || part2 != 56000011 {
		t.Fatalf("got (%d, %d, %v), want (26, 56000011, nil)", part1, part2, err)
	}
}

func TestSolveUsesTheRealPuzzleRow(t *testing.T) {
	part1, _, err := solve(readExample(t))
	if err != nil || part1 != 0 {
		t.Fatalf("got (%d, %v), want (0, nil): no example sensor reaches row 2000000", part1, err)
	}
}

func TestBeaconCases(t *testing.T) {
	// Row 0 of this sensor covers x in [-2, 2]; its known beacon at x=2 stays a
	// valid position, and row 1 leaves x=2 uncovered.
	const oneSensor = "Sensor at x=0, y=0: closest beacon is at x=2, y=0"

	tests := []struct {
		name      string
		input     string
		row       int
		limit     int
		part1     int
		part2     int
		wantError bool
	}{
		{name: "single sensor rules out its diamond row", input: oneSensor, row: 0, limit: 2, part1: 4, part2: 8000001},
		{name: "row outside every diamond rules out nothing", input: oneSensor, row: 1000, limit: 2, part1: 0, part2: 8000001},
		{name: "fully covered search area", input: "Sensor at x=0, y=0: closest beacon is at x=10, y=0", row: 0, limit: 2, wantError: true},
		{name: "no readings", input: "\n\n", row: 0, limit: 2, wantError: true},
		{name: "truncated reading", input: "Sensor at x=2, y=18", row: 0, limit: 2, wantError: true},
		{name: "non-numeric coordinate", input: "Sensor at x=a, y=18: closest beacon is at x=-2, y=15", row: 0, limit: 2, wantError: true},
		{name: "missing sensor prefix", input: "at x=2, y=18: closest beacon is at x=-2, y=15", row: 0, limit: 2, wantError: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solveWith(tc.input, tc.row, tc.limit)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
