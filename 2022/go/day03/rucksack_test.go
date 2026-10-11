package main

import (
	"os"
	"testing"
)

func readExample(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	return string(data)
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "example", input: readExample(t), want: 157},
		{name: "single rucksack", input: "vJrwpWtwJgWrhcsFMMfFFhFp\n", want: 16},
		{name: "odd length rucksack", input: "abc\n", wantErr: true},
		{name: "invalid item type", input: "ab1ab1\n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := part1(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("part1(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("part1(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("part1(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "example", input: readExample(t), want: 70},
		{name: "not multiple of three", input: "abc\ndef\n", wantErr: true},
		{name: "empty input", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := part2(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("part2(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("part2(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("part2(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestSolve(t *testing.T) {
	p1, p2, err := solve(readExample(t))
	if err != nil {
		t.Fatalf("solve: unexpected error: %v", err)
	}
	if p1 != 157 || p2 != 70 {
		t.Errorf("solve = (%d, %d), want (157, 70)", p1, p2)
	}
}
