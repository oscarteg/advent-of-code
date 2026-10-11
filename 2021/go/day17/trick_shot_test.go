package main

import (
	"os"
	"testing"
)

func TestParseTarget(t *testing.T) {
	target, err := parseTarget("target area: x=20..30, y=-10..-5")
	if err != nil {
		t.Fatalf("parseTarget failed: %v", err)
	}
	if target.xMin != 20 || target.xMax != 30 || target.yMin != -10 || target.yMax != -5 {
		t.Errorf("parseTarget: got %v, want x=20..30, y=-10..-5", target)
	}
}

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 45 || part2 != 112 {
		t.Errorf("solve: got (%d, %d), want (45, 112)", part1, part2)
	}
}

func TestSpecificInput(t *testing.T) {
	tests := []struct {
		input string
		part1 int
		part2 int
	}{
		{"target area: x=20..30, y=-10..-5", 45, 112},
	}

	for _, tt := range tests {
		t.Run("inline", func(t *testing.T) {
			p1, p2, err := solve(tt.input)
			if err != nil {
				t.Fatalf("solve failed: %v", err)
			}
			if p1 != tt.part1 || p2 != tt.part2 {
				t.Errorf("solve: got (%d, %d), want (%d, %d)", p1, p2, tt.part1, tt.part2)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	_, _, err := solve("invalid input")
	if err == nil {
		t.Error("solve should reject malformed input")
	}
}

func TestSimulate(t *testing.T) {
	target := rect{20, 30, -10, -5}
	tests := []struct {
		vx        int
		vy        int
		shouldHit bool
	}{
		{7, 2, true},
		{6, 3, true},
		{9, 0, true},
		{0, 0, false},
	}

	for _, tt := range tests {
		_, hit := simulate(tt.vx, tt.vy, target)
		if hit != tt.shouldHit {
			t.Errorf("simulate(%d, %d) hit=%v, want %v", tt.vx, tt.vy, hit, tt.shouldHit)
		}
	}
}
