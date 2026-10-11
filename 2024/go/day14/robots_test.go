package main

import (
	"os"
	"testing"
)

func TestPart1Example(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, _, err := solveWith(string(example), 11, 7)
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 12 {
		t.Errorf("Part 1 example: got %d, want 12", part1)
	}
}

func TestInferredDimensionsMatchExample(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, _, err := solve(string(example))
	if err != nil {
		t.Fatal(err)
	}
	if part1 != 12 {
		t.Errorf("inferred-dimension part 1: got %d, want 12", part1)
	}
}

func TestPictureDetectorFindsImmediateNonOverlap(t *testing.T) {
	robots, err := parseRobots("p=0,0 v=0,0\np=1,0 v=0,0\n")
	if err != nil {
		t.Fatal(err)
	}
	t_, err := findPictureSecond(robots, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if t_ != 0 {
		t.Errorf("got %d, want 0", t_)
	}
}

func TestPictureDetectorFindsLaterNonOverlap(t *testing.T) {
	robots, err := parseRobots("p=0,0 v=1,0\np=0,0 v=2,0\n")
	if err != nil {
		t.Fatal(err)
	}
	second, err := findPictureSecond(robots, 5, 1)
	if err != nil {
		t.Fatal(err)
	}
	if second != 1 {
		t.Errorf("got %d, want 1", second)
	}
}

func TestMalformedInput(t *testing.T) {
	if _, _, err := solve(""); err == nil {
		t.Error("empty input should error")
	}
	if _, _, err := solve("garbage line\n"); err == nil {
		t.Error("garbage line should error")
	}
}
