package main

import (
	"os"
	"strings"
	"testing"
)

func TestSolveOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 37 || part2 != 168 {
		t.Fatalf("solve(example) = (%d, %d), want (37, 168)", part1, part2)
	}
}

func TestSolveSinglePosition(t *testing.T) {
	part1, part2, err := solve("8")
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 0 || part2 != 0 {
		t.Fatalf("solve(single position) = (%d, %d), want (0, 0)", part1, part2)
	}
}

func TestSolveHandlesEvenPositionCount(t *testing.T) {
	part1, part2, err := solve("0,10")
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 10 || part2 != 30 {
		t.Fatalf("solve(even positions) = (%d, %d), want (10, 30)", part1, part2)
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	for _, input := range []string{"", "1,,2", "1,not-a-number", "-1,2", "1,2\n3"} {
		t.Run(strings.ReplaceAll(input, ",", "_"), func(t *testing.T) {
			if _, _, err := solve(input); err == nil {
				t.Fatalf("solve(%q) succeeded, want an error", input)
			}
		})
	}
}
