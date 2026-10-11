package main

import (
	"os"
	"testing"
)

func TestSmallExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("Failed to read example: %v", err)
	}

	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 10 {
		t.Errorf("part1: want 10, got %d", part1)
	}
	if part2 != 36 {
		t.Errorf("part2: want 36, got %d", part2)
	}
}

func TestMediumExample(t *testing.T) {
	input := `dc-end
HN-start
start-kj
dc-start
dc-HN
LN-dc
HN-end
kj-sa
kj-HN
kj-dc`

	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 19 {
		t.Errorf("part1: want 19, got %d", part1)
	}
	if part2 != 103 {
		t.Errorf("part2: want 103, got %d", part2)
	}
}

func TestLargeExample(t *testing.T) {
	input := `fs-end
he-DX
fs-he
start-DX
pj-DX
end-zg
zg-sl
zg-pj
pj-he
RW-he
fs-DX
pj-RW
zg-RW
start-pj
he-WI
zg-he
pj-fs
start-RW`

	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 226 {
		t.Errorf("part1: want 226, got %d", part1)
	}
	if part2 != 3509 {
		t.Errorf("part2: want 3509, got %d", part2)
	}
}

func TestSinglePath(t *testing.T) {
	input := "start-A\nA-end"

	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 1 {
		t.Errorf("part1: want 1, got %d", part1)
	}
	if part2 != 1 {
		t.Errorf("part2: want 1, got %d", part2)
	}
}

func TestInvalidEdgeFormat(t *testing.T) {
	input := "start-A-B"

	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for invalid edge format")
	}
}

func TestInvalidEdgeNoHyphen(t *testing.T) {
	input := "startA"

	_, _, err := solve(input)
	if err == nil {
		t.Error("expected error for edge without hyphen")
	}
}
