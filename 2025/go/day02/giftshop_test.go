package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	example, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	tests := []struct {
		name         string
		input        string
		part1, part2 uint64
	}{
		{
			name:  "official example",
			input: string(example),
			part1: 1227775554,
			part2: 4174379265,
		},
		{
			name:  "part 2 counts sequences repeated more than twice",
			input: "123123123-123123123,1111111-1111111",
			part1: 0,
			part2: 123123123 + 1111111,
		},
		{
			name:  "single-digit ids never count as repeated",
			input: "5-5",
			part1: 0,
			part2: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part1 != tc.part1 {
				t.Errorf("part1 = %d, want %d", part1, tc.part1)
			}
			if part2 != tc.part2 {
				t.Errorf("part2 = %d, want %d", part2, tc.part2)
			}
		})
	}
}

func TestSolveMalformedInput(t *testing.T) {
	malformed := []string{
		"11-22,notarange",
		"22-11",
		"abc-5",
	}
	for _, input := range malformed {
		if _, _, err := solve(input); err == nil {
			t.Errorf("solve(%q): expected error, got nil", input)
		}
	}
}
