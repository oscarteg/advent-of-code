package main

import (
	"os"
	"testing"
)

func TestPart1Example(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	part1, _, err := solve(string(example))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part1 != 15 {
		t.Errorf("part1 = %d, want 15", part1)
	}
}

func TestPart2Example(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}
	_, part2, err := solve(string(example))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}
	if part2 != 12 {
		t.Errorf("part2 = %d, want 12", part2)
	}
}

func TestCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 int
		want2 int
	}{
		{
			name:  "part1 win",
			input: "A Y",
			want1: 8,
			want2: 4,
		},
		{
			name:  "part1 loss",
			input: "B X",
			want1: 1,
			want2: 1,
		},
		{
			name:  "part1 draw",
			input: "C Z",
			want1: 6,
			want2: 7,
		},
		{
			name:  "invalid move",
			input: "D X",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if tt.name == "invalid move" {
				if err == nil {
					t.Errorf("solve should have failed for invalid move")
				}
				return
			}
			if err != nil {
				t.Fatalf("solve failed: %v", err)
			}
			if part1 != tt.want1 {
				t.Errorf("part1 = %d, want %d", part1, tt.want1)
			}
			if part2 != tt.want2 {
				t.Errorf("part2 = %d, want %d", part2, tt.want2)
			}
		})
	}
}
