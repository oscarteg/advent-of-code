package main

import (
	"os"
	"strings"
	"testing"
)

func TestSolve(t *testing.T) {
	exampleBytes, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	example := string(exampleBytes)

	tests := []struct {
		name          string
		input         string
		part1, part2  int64
		wantErrSubstr string
	}{
		{name: "official example", input: example, part1: 4277556, part2: 3263827},
		{name: "single problem edge case", input: "5\n+\n", part1: 5, part2: 5},
		{name: "empty input", input: "", wantErrSubstr: "parse worksheet"},
		{name: "missing operator", input: "12\n34\n", wantErrSubstr: "find operator"},
		{name: "invalid digit", input: "1x\n+ \n", wantErrSubstr: "solve by rows"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if tt.wantErrSubstr != "" {
				if err == nil {
					t.Fatalf("solve(%q) = (%d, %d, nil), want error containing %q", tt.input, part1, part2, tt.wantErrSubstr)
				}
				if !strings.Contains(err.Error(), tt.wantErrSubstr) {
					t.Fatalf("solve(%q) error = %q, want containing %q", tt.input, err.Error(), tt.wantErrSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("solve(%q) unexpected error: %v", tt.input, err)
			}
			if part1 != tt.part1 || part2 != tt.part2 {
				t.Errorf("solve(%q) = (%d, %d), want (%d, %d)", tt.input, part1, part2, tt.part1, tt.part2)
			}
		})
	}
}
