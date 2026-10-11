package main

import (
	"errors"
	"os"
	"strconv"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if want := "4,6,3,5,6,3,5,2,1,0"; part1 != want {
		t.Errorf("part1 = %q, want %q", part1, want)
	}
	if part2 != 0 {
		t.Errorf("part2 = %d, want 0: no A makes this program print itself", part2)
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		part1     string
		part2     uint64
		wantError bool
	}{
		{
			name:  "program that prints itself",
			input: "Register A: 2024\nRegister B: 0\nRegister C: 0\n\nProgram: 0,3,5,4,3,0\n",
			part1: "5,7,3,0",
			part2: 117440,
		},
		{
			name:  "literal and register operands",
			input: "Register A: 10\nRegister B: 0\nRegister C: 0\n\nProgram: 5,0,5,1,5,4\n",
			part1: "0,1,2",
		},
		{
			name:  "shifting program",
			input: "Register A: 2024\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1,5,4,3,0\n",
			part1: "4,2,5,6,7,7,7,7,3,1,0",
		},
		{
			name:  "program without an out instruction prints nothing",
			input: "Register A: 0\nRegister B: 29\nRegister C: 0\n\nProgram: 1,7\n",
			part1: "",
		},
		{name: "empty input", input: "", wantError: true},
		{
			name:      "missing program line",
			input:     "Register A: 1\nRegister B: 0\nRegister C: 0\n",
			wantError: true,
		},
		{
			name:      "non-numeric register",
			input:     "Register A: x\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1\n",
			wantError: true,
		},
		{
			name:      "program value wider than three bits",
			input:     "Register A: 1\nRegister B: 0\nRegister C: 0\n\nProgram: 0,8\n",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if tc.wantError {
				if err == nil {
					t.Fatalf("got (%q, %d, nil), want an error", part1, part2)
				}
				return
			}
			if err != nil {
				t.Fatalf("solve: %v", err)
			}
			if part1 != tc.part1 {
				t.Errorf("part1 = %q, want %q", part1, tc.part1)
			}
			if part2 != tc.part2 {
				t.Errorf("part2 = %d, want %d", part2, tc.part2)
			}
		})
	}
}

func TestParseRejectsOutOfRangeProgramValue(t *testing.T) {
	_, err := parse("Register A: 1\nRegister B: 0\nRegister C: 0\n\nProgram: 9\n")
	if !errors.Is(err, strconv.ErrRange) {
		t.Fatalf("err = %v, want it to wrap strconv.ErrRange", err)
	}
}
