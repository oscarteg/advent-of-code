package main

import (
	"os"
	"testing"
)

func TestSolveExample(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	part1, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve example failed: %v", err)
	}
	if part1 != "2=-1=0" {
		t.Errorf("part1 example: got %q, want %q", part1, "2=-1=0")
	}
}

func TestSolveSingleNumber(t *testing.T) {
	part1, err := solve("1=-0-2")
	if err != nil {
		t.Fatalf("solve single number failed: %v", err)
	}
	if part1 != "1=-0-2" {
		t.Errorf("part1: got %q, want %q", part1, "1=-0-2")
	}
}

func TestSolveMalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"blank lines only", "   \n  \n"},
		{"invalid digit", "1=-0-9"},
		{"non-snafu character", "1a2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := solve(tt.input); err == nil {
				t.Errorf("expected an error for input %q", tt.input)
			}
		})
	}
}

func TestParseSnafuDigits(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"1", 1},
		{"2", 2},
		{"1=", 3},
		{"1-", 4},
		{"10", 5},
	}
	for _, tt := range tests {
		got, err := parseSnafu(tt.input)
		if err != nil {
			t.Fatalf("parseSnafu(%q) failed: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("parseSnafu(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestToSnafuRoundTrip(t *testing.T) {
	values := []int64{1, 2, 3, 4, 5, 10, 2022, 4890}
	for _, value := range values {
		snafu := toSnafu(value)
		got, err := parseSnafu(snafu)
		if err != nil {
			t.Fatalf("parseSnafu(%q) failed: %v", snafu, err)
		}
		if got != value {
			t.Errorf("round trip for %d: got %d via %q", value, got, snafu)
		}
	}
}
