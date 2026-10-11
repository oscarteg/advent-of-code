package main

import (
	"errors"
	"os"
	"testing"
)

func TestSolveOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if part1 != 50 {
		t.Errorf("part1 = %d, want 50", part1)
	}
	if part2 != 24 {
		t.Errorf("part2 = %d, want 24", part2)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
	}{
		{
			name:  "square loop covers every tile it spans",
			input: "0,0\n4,0\n4,4\n0,4\n",
			part1: 25,
			part2: 25,
		},
		{
			name:  "L-shaped loop rules out rectangles reaching into the notch",
			input: "0,0\n4,0\n4,2\n2,2\n2,4\n0,4\n",
			part1: 25,
			part2: 15,
		},
		{
			name:  "comb loop keeps rectangles out of its gaps",
			input: "0,0\n9,0\n9,6\n7,6\n7,2\n5,2\n5,6\n3,6\n3,2\n1,2\n1,6\n0,6\n",
			part1: 70,
			part2: 27,
		},
		{
			name:  "surrounding whitespace is ignored",
			input: "  0,0 \n\n 4,0\n4,4\n 0,4 \n\n",
			part1: 25,
			part2: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part1 != tt.part1 {
				t.Errorf("part1 = %d, want %d", part1, tt.part1)
			}
			if part2 != tt.part2 {
				t.Errorf("part2 = %d, want %d", part2, tt.part2)
			}
		})
	}
}

func TestSolveRejectsMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{name: "missing coordinate", input: "0,0\n4,0\n4\n0,4\n", want: errInvalidTile},
		{name: "third coordinate", input: "0,0\n4,0,1\n4,4\n0,4\n", want: errInvalidTile},
		{name: "non-numeric coordinate", input: "0,0\n4,x\n4,4\n0,4\n", want: errInvalidTile},
		{name: "fewer than four tiles", input: "1,2\n3,4\n", want: errNotALoop},
		{name: "only blank lines", input: "  \n \n", want: errNotALoop},
		{name: "diagonal step between tiles", input: "0,0\n4,1\n4,4\n0,4\n", want: errNotALoop},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := solve(tt.input); !errors.Is(err, tt.want) {
				t.Errorf("solve(%q) error = %v, want %v", tt.input, err, tt.want)
			}
		})
	}
}
