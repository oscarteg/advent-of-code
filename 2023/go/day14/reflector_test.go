package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		p1    int
		p2    int
	}{
		{
			name:  "single rock",
			input: "O\n.\n.",
			p1:    3,
			p2:    1,
		},
		{
			name:  "with blockers",
			input: "O.O\n###\nO.O",
			p1:    8,
			p2:    8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, p2, err := solve(tt.input)
			if err != nil {
				t.Fatalf("solve() error = %v", err)
			}
			if p1 != tt.p1 {
				t.Errorf("part 1: got %d, want %d", p1, tt.p1)
			}
			if p2 != tt.p2 {
				t.Errorf("part 2: got %d, want %d", p2, tt.p2)
			}
		})
	}
}

func TestExampleFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	p1, p2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve() error = %v", err)
	}

	if p1 != 136 {
		t.Errorf("part 1: got %d, want 136", p1)
	}
	if p2 != 64 {
		t.Errorf("part 2: got %d, want 64", p2)
	}
}

func TestMalformedInput(t *testing.T) {
	for _, input := range []string{"", "O.\n.", "O.\n.x"} {
		if _, _, err := solve(input); err == nil {
			t.Errorf("solve(%q): want an error", input)
		}
	}
}
