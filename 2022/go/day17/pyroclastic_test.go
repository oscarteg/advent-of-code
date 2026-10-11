package main

import (
	"os"
	"testing"
)

func TestSolveOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("reading example: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 3068 {
		t.Errorf("part1 = %d, want 3068", part1)
	}
	if part2 != 1514285714288 {
		t.Errorf("part2 = %d, want 1514285714288", part2)
	}
}

func TestTowerHeightSingleJetPushesRockToTheWall(t *testing.T) {
	jets, err := parseJets("<")
	if err != nil {
		t.Fatalf("parseJets: %v", err)
	}
	if got := towerHeight(jets, 1); got != 1 {
		t.Errorf("towerHeight = %d, want 1", got)
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	cases := []string{"abc", "", "   "}
	for _, input := range cases {
		if _, _, err := solve(input); err == nil {
			t.Errorf("solve(%q) = nil error, want error", input)
		}
	}
}
