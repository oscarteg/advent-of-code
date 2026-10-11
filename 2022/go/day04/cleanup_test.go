package main

import (
	"os"
	"testing"
)

const exampleInput = `2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8`

func TestSolveExample(t *testing.T) {
	part1, part2, err := solve(exampleInput)
	if err != nil {
		t.Fatalf("solve() returned error: %v", err)
	}
	if part1 != 2 {
		t.Errorf("part1 = %d, want 2", part1)
	}
	if part2 != 4 {
		t.Errorf("part2 = %d, want 4", part2)
	}
}

func TestSolveExampleFile(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve() returned error: %v", err)
	}
	if part1 != 2 {
		t.Errorf("part1 = %d, want 2", part1)
	}
	if part2 != 4 {
		t.Errorf("part2 = %d, want 4", part2)
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		input string
		want  Range
		err   bool
	}{
		{"2-4", Range{2, 4}, false},
		{"6-8", Range{6, 8}, false},
		{"invalid", Range{}, true},
		{"2-a", Range{}, true},
	}

	for _, tt := range tests {
		got, err := parseRange(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("parseRange(%q) error = %v, want error %v", tt.input, err, tt.err)
		}
		if !tt.err && got != tt.want {
			t.Errorf("parseRange(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		a    Range
		b    Range
		want bool
	}{
		{Range{2, 8}, Range{3, 7}, true},
		{Range{3, 7}, Range{2, 8}, false},
		{Range{4, 6}, Range{6, 6}, true},
		{Range{6, 6}, Range{4, 6}, false},
	}

	for _, tt := range tests {
		got := contains(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("contains(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestOverlaps(t *testing.T) {
	tests := []struct {
		a    Range
		b    Range
		want bool
	}{
		{Range{5, 7}, Range{7, 9}, true},
		{Range{2, 8}, Range{3, 7}, true},
		{Range{6, 6}, Range{4, 6}, true},
		{Range{2, 4}, Range{6, 8}, false},
	}

	for _, tt := range tests {
		got := overlaps(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("overlaps(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []string{
		"2-4",
		"2-4,invalid",
		"a-4,6-8",
	}

	for _, input := range tests {
		_, _, err := solve(input)
		if err == nil {
			t.Errorf("solve(%q) expected error, got nil", input)
		}
	}
}
