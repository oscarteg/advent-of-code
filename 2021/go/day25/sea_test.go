package main

import (
	"os"
	"reflect"
	"testing"
)

func TestExampleInput(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	result, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if result != 58 {
		t.Errorf("expected 58, got %d", result)
	}
}

func TestPackedGridStopsOnFirstStep(t *testing.T) {
	result, err := solve(">v\nv>")
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}
}

func TestEastHerdMovesOnlyIntoFreeCells(t *testing.T) {
	grid, err := parseGrid("...>>>>>...")
	if err != nil {
		t.Fatal(err)
	}
	want, err := parseGrid("...>>>>.>..")
	if err != nil {
		t.Fatal(err)
	}

	moveEast(grid, 1, 11)

	if !reflect.DeepEqual(grid, want) {
		t.Errorf("after one east move got %v, want %v", grid, want)
	}
}

func TestInvalidChar(t *testing.T) {
	input := "x"
	_, err := solve(input)
	if err == nil {
		t.Error("expected error for invalid character")
	}
}

func TestEmptyInput(t *testing.T) {
	_, err := solve("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}
