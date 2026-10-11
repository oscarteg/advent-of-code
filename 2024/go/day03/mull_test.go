package main

import (
	"os"
	"testing"
)

func TestSimpleMul(t *testing.T) {
	part1, part2, err := solve("mul(2,3)mul(4,5)")
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 26 {
		t.Errorf("part1: expected 26, got %d", part1)
	}
	if part2 != 26 {
		t.Errorf("part2: expected 26, got %d", part2)
	}
}

func TestMulWithDoAndDont(t *testing.T) {
	part1, part2, err := solve("mul(2,3)don't()mul(4,5)do()mul(1,1)")
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 27 {
		t.Errorf("part1: expected 27, got %d", part1)
	}
	if part2 != 7 {
		t.Errorf("part2: expected 7, got %d", part2)
	}
}

func TestInvalidPatterns(t *testing.T) {
	part1, part2, err := solve("mul(1000,1000)mul(1,)mul(,1)mul(1 2)mul(1,2")
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 0 {
		t.Errorf("part1: expected 0, got %d", part1)
	}
	if part2 != 0 {
		t.Errorf("part2: expected 0, got %d", part2)
	}
}

func TestEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestExampleInput(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	part1, _, err := solve(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 161 {
		t.Errorf("part1: expected 161, got %d", part1)
	}
}

func TestExamplePart2Input(t *testing.T) {
	data, err := os.ReadFile("testdata/example2.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 161 || part2 != 48 {
		t.Errorf("got (%d, %d), want (161, 48)", part1, part2)
	}
}
