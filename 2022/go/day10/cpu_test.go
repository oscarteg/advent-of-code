package main

import (
	"os"
	"testing"
)

func TestExamplePart1(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("could not read example: %v", err)
	}

	part1, _, err := solve(string(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if part1 != 13140 {
		t.Errorf("part 1: got %d, want 13140", part1)
	}
}

func TestExamplePart2(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("could not read example: %v", err)
	}

	_, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "##..##..##..##..##..##..##..##..##..##..\n" +
		"###...###...###...###...###...###...###.\n" +
		"####....####....####....####....####....\n" +
		"#####.....#####.....#####.....#####.....\n" +
		"######......######......######......####\n" +
		"#######.......#######.......#######....."

	if part2 != expected {
		t.Errorf("part 2 mismatch:\ngot:\n%s\nwant:\n%s", part2, expected)
	}
}

func TestNoopOnly(t *testing.T) {
	part1, _, err := solve("noop\nnoop\nnoop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if part1 != 0 {
		t.Errorf("part 1: got %d, want 0", part1)
	}
}

func TestMalformedInstruction(t *testing.T) {
	_, _, err := solve("invalid")
	if err == nil {
		t.Errorf("expected error for invalid instruction")
	}
}

func TestMalformedAddx(t *testing.T) {
	_, _, err := solve("addx abc")
	if err == nil {
		t.Errorf("expected error for invalid addx value")
	}
}
