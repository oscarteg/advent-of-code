package main

import (
	"os"
	"testing"
)

var exampleInput = `199
200
208
210
200
207
240
269
260
263`

func TestSolveExample(t *testing.T) {
	part1, part2, err := solve(exampleInput)
	if err != nil {
		t.Fatalf("solve example failed: %v", err)
	}
	if part1 != 7 {
		t.Errorf("part1 example: got %d, want 7", part1)
	}
	if part2 != 5 {
		t.Errorf("part2 example: got %d, want 5", part2)
	}
}

func TestSolveRealInput(t *testing.T) {
	data, err := os.ReadFile("../../rust/input/day1.txt")
	if err != nil {
		t.Fatal(err)
	}
	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve real input failed: %v", err)
	}
	if part1 != 1692 {
		t.Errorf("part1 real: got %d, want 1692", part1)
	}
	if part2 != 1724 {
		t.Errorf("part2 real: got %d, want 1724", part2)
	}
}

func TestSolveEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		part1   int
		part2   int
		wantErr bool
	}{
		{
			name:    "single measurement",
			input:   "100",
			part1:   0,
			part2:   0,
			wantErr: false,
		},
		{
			name:    "two measurements",
			input:   "100\n200",
			part1:   1,
			part2:   0,
			wantErr: false,
		},
		{
			name:    "three measurements",
			input:   "100\n200\n300",
			part1:   2,
			part2:   0,
			wantErr: false,
		},
		{
			name:    "four measurements",
			input:   "100\n200\n150\n300",
			part1:   2,
			part2:   1,
			wantErr: false,
		},
		{
			name:    "all decreasing",
			input:   "300\n200\n100",
			part1:   0,
			part2:   0,
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			part1:   0,
			part2:   0,
			wantErr: true,
		},
		{
			name:    "malformed input",
			input:   "100\nabc\n200",
			part1:   0,
			part2:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solve(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, want error %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if part1 != tt.part1 {
					t.Errorf("part1: got %d, want %d", part1, tt.part1)
				}
				if part2 != tt.part2 {
					t.Errorf("part2: got %d, want %d", part2, tt.part2)
				}
			}
		})
	}
}
