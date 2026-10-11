package main

import (
	"errors"
	"os"
	"testing"
)

func exampleCodes(t *testing.T) string {
	t.Helper()
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}
	return string(input)
}

func TestOfficialExample(t *testing.T) {
	part1, part2, err := solve(exampleCodes(t))
	if err != nil {
		t.Fatalf("solve error = %v", err)
	}
	if part1 != 126384 {
		t.Errorf("part 1 = %d, want 126384", part1)
	}
	if part2 != 154115708116294 {
		t.Errorf("part 2 = %d, want 154115708116294", part2)
	}
}

func TestPublishedSequenceLengths(t *testing.T) {
	tests := []struct {
		code   string
		length int64
	}{
		{"029A", 68},
		{"980A", 60},
		{"179A", 68},
		{"456A", 64},
		{"379A", 64},
	}
	for _, tc := range tests {
		t.Run(tc.code, func(t *testing.T) {
			presses, err := newSolver().typeOn(numericPad, tc.code, part1Robots)
			if err != nil {
				t.Fatalf("typeOn error = %v", err)
			}
			if presses != tc.length {
				t.Errorf("presses = %d, want %d", presses, tc.length)
			}
		})
	}
}

func TestCodesAreIndependent(t *testing.T) {
	firstPart1, firstPart2, err := solve("029A\n")
	if err != nil {
		t.Fatalf("solve error = %v", err)
	}
	secondPart1, secondPart2, err := solve("980A\n")
	if err != nil {
		t.Fatalf("solve error = %v", err)
	}
	bothPart1, bothPart2, err := solve("029A\n980A\n")
	if err != nil {
		t.Fatalf("solve error = %v", err)
	}

	if bothPart1 != firstPart1+secondPart1 || bothPart2 != firstPart2+secondPart2 {
		t.Errorf("got (%d, %d), want (%d, %d)", bothPart1, bothPart2,
			firstPart1+secondPart1, firstPart2+secondPart2)
	}
}

func TestSingleCodes(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int64
	}{
		{"one code", "029A\n", 68 * 29, 82050061710 * 29},
		{"leading zeroes leave no complexity", "000A\n", 0, 0},
		{"carriage returns and blank lines", "029A\r\n\r\n", 68 * 29, 82050061710 * 29},
		{"no trailing newline", "029A", 68 * 29, 82050061710 * 29},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if err != nil {
				t.Fatalf("solve error = %v", err)
			}
			if part1 != tc.part1 || part2 != tc.part2 {
				t.Errorf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}

func TestMalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"empty input", ""},
		{"only blank lines", "   \n\n"},
		{"no A to submit the code", "029\n"},
		{"not a digit", "02XA\n"},
		{"no digits", "A\n"},
		{"signed numeric part", "-29A\n"},
		{"lowercase submit key", "029a\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := solve(tc.input)

			if !errors.Is(err, errInvalidCode) {
				t.Fatalf("error = %v, want one wrapping %v", err, errInvalidCode)
			}
		})
	}
}
