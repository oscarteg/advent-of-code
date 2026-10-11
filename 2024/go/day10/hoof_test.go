package main

import (
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	p1, p2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if p1 != 36 {
		t.Errorf("part 1: got %d, want 36", p1)
	}
	if p2 != 81 {
		t.Errorf("part 2: got %d, want 81", p2)
	}
}

func TestSingleTrail(t *testing.T) {
	input := "0123456789\n"
	p1, p2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if p1 != 1 {
		t.Errorf("part 1: got %d, want 1", p1)
	}
	if p2 != 1 {
		t.Errorf("part 2: got %d, want 1", p2)
	}
}

func TestVerticalTrail(t *testing.T) {
	input := "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n"
	p1, p2, err := solve(input)
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if p1 != 1 {
		t.Errorf("part 1: got %d, want 1", p1)
	}
	if p2 != 1 {
		t.Errorf("part 2: got %d, want 1", p2)
	}
}

func TestEmptyInput(t *testing.T) {
	_, _, err := solve("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestInvalidCharacter(t *testing.T) {
	_, _, err := solve("0a1\n")
	if err == nil {
		t.Error("expected error for invalid character")
	}
}

func TestIrregularGrid(t *testing.T) {
	_, _, err := solve("012\n01\n")
	if err == nil {
		t.Error("expected error for irregular grid")
	}
}
