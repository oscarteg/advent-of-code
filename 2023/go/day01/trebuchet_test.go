package main

import (
	"os"
	"testing"
)

func TestPart1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"example", "1abc2\npqr3stu8vwx\na1b2c3d4e5f\ntreb7uchet", 142},
		{"single digit", "5", 55},
		{"many digits", "1234567", 17},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _, err := solve(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"example", "two1nine\neightwothree\nabcone2threexyz\nxtwone3four\n4nineeightseven2\nzoneight234\n7pqrstsixteen", 281},
		{"spelled only", "one", 11},
		{"overlapping", "twone", 21},
		{"eightwo overlap", "eightwo", 82},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, result, err := solve(tt.input)
			if err != nil {
				t.Fatalf("unexpected error on %q: %v", tt.name, err)
			}
			if result != tt.expected {
				t.Errorf("got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestProcessLinePart2Direct(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"one", 11},
		{"eightwothree", 83},
		{"twone", 21},
	}

	for _, tt := range tests {
		result, err := processLinePart2(tt.input)
		if err != nil {
			t.Errorf("processLinePart2(%q): unexpected error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("processLinePart2(%q): got %d, want %d", tt.input, result, tt.expected)
		}
	}
}

func TestOfficialExampleFixtures(t *testing.T) {
	tests := []struct {
		file     string
		part     int
		expected int
	}{
		{"testdata/example.txt", 1, 142},
		{"testdata/example2.txt", 2, 281},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			input, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatal(err)
			}
			part1, part2, err := solve(string(input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := []int{part1, part2}[tt.part-1]; got != tt.expected {
				t.Errorf("part %d = %d, want %d", tt.part, got, tt.expected)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	if _, _, err := solve("1abc2\nxyz"); err == nil {
		t.Error("want an error for a line without any digit")
	}
}
