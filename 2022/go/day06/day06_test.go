package main

import (
	"os"
	"testing"
)

func TestPart1Examples(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 7},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 5},
		{"nppdvjthqldpwncqszvftbrmjlhg", 6},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 10},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 11},
	}

	for _, tt := range tests {
		got, err := findMarker(tt.input, 4)
		if err != nil {
			t.Errorf("findMarker(%q, 4) error: %v", tt.input, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("findMarker(%q, 4) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestPart2Examples(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 19},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 23},
		{"nppdvjthqldpwncqszvftbrmjlhg", 23},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 29},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 26},
	}

	for _, tt := range tests {
		got, err := findMarker(tt.input, 14)
		if err != nil {
			t.Errorf("findMarker(%q, 14) error: %v", tt.input, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("findMarker(%q, 14) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestExampleFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("reading example file: %v", err)
	}
	input := string(data)

	part1, part2, err := solve(input)
	if err != nil {
		t.Fatalf("solve error: %v", err)
	}

	if part1 != 7 {
		t.Errorf("part 1: got %v, want 7", part1)
	}
	if part2 != 19 {
		t.Errorf("part 2: got %v, want 19", part2)
	}
}

func TestEdgeCases(t *testing.T) {
	got, err := findMarker("abcd", 4)
	if err != nil {
		t.Errorf("findMarker(\"abcd\", 4) error: %v", err)
	}
	if got != 4 {
		t.Errorf("findMarker(\"abcd\", 4) = %v, want 4", got)
	}

	got, err = findMarker("abcde", 4)
	if err != nil {
		t.Errorf("findMarker(\"abcde\", 4) error: %v", err)
	}
	if got != 4 {
		t.Errorf("findMarker(\"abcde\", 4) = %v, want 4", got)
	}
}

func TestMalformedInput(t *testing.T) {
	if _, err := findMarker("ab", 4); err == nil {
		t.Errorf("findMarker(\"ab\", 4) expected error")
	}
}

func TestNoMarkerFound(t *testing.T) {
	if _, err := findMarker("aaaaaaaaaa", 4); err == nil {
		t.Errorf("findMarker(\"aaaaaaaaaa\", 4) expected error")
	}
}
