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

	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 26397 {
		t.Errorf("part1 = %d, want 26397", part1)
	}
	if part2 != 288957 {
		t.Errorf("part2 = %d, want 288957", part2)
	}
}

func TestSingleCorrupted(t *testing.T) {
	line := "{([(<{}[<>[]}>{[]{[(<()>"
	status, err := analyzeLine(line)
	if err != nil {
		t.Fatalf("analyzeLine failed: %v", err)
	}

	if !status.corrupted {
		t.Errorf("expected corrupted line, got incomplete")
	}
	if status.errorScore != 1197 {
		t.Errorf("errorScore = %d, want 1197", status.errorScore)
	}
}

func TestIncompleteLine(t *testing.T) {
	line := "[({(<(())[]>[[{[]{<()<>>"
	status, err := analyzeLine(line)
	if err != nil {
		t.Fatalf("analyzeLine failed: %v", err)
	}

	if !status.incomplete {
		t.Errorf("expected incomplete line, got corrupted or error")
	}
	if len(status.stackToClose) == 0 {
		t.Errorf("expected non-empty stack, got empty")
	}
}

func TestMalformedEmptyStack(t *testing.T) {
	line := ")"
	status, err := analyzeLine(line)
	if err != nil {
		t.Fatalf("analyzeLine failed: %v", err)
	}

	if !status.corrupted {
		t.Errorf("expected corrupted line")
	}
	if status.errorScore != 3 {
		t.Errorf("errorScore = %d, want 3", status.errorScore)
	}
}

func TestInvalidCharacter(t *testing.T) {
	line := "@ []"
	_, err := analyzeLine(line)
	if err == nil {
		t.Errorf("expected error for invalid character, got nil")
	}
}
