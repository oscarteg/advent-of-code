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

func TestSolve(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantPart1    string
		wantPart2    string
		wantParseErr bool
	}{
		{
			name:      "official example",
			input:     readExample(t),
			wantPart1: "CMZ",
			wantPart2: "MCD",
		},
		{
			name:      "single crate single move",
			input:     "[A]\n 1 \n\nmove 1 from 1 to 1",
			wantPart1: "A",
			wantPart2: "A",
		},
		{
			name:         "malformed input missing blank line",
			input:        "[A]\n 1 \nmove 1 from 1 to 1",
			wantParseErr: true,
		},
		{
			name:         "malformed move line",
			input:        "[A]\n 1 \n\nmove one from 1 to 1",
			wantParseErr: true,
		},
		{
			name:         "move more crates than available",
			input:        "[A]\n 1 \n\nmove 5 from 1 to 1",
			wantParseErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if tt.wantParseErr {
				if err == nil {
					t.Fatalf("solve(%q) expected error, got none", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("solve(%q) unexpected error: %v", tt.input, err)
			}
			if part1 != tt.wantPart1 {
				t.Errorf("part1 = %q, want %q", part1, tt.wantPart1)
			}
			if part2 != tt.wantPart2 {
				t.Errorf("part2 = %q, want %q", part2, tt.wantPart2)
			}
		})
	}
}
