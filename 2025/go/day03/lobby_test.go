package main

import (
	"os"
	"testing"
)

func TestExampleInput(t *testing.T) {
	data, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatalf("failed to read example: %v", err)
	}

	part1, part2, err := solve(string(data))
	if err != nil {
		t.Fatalf("solve failed: %v", err)
	}

	if part1 != 357 {
		t.Errorf("part1: got %d, want 357", part1)
	}
	if part2 != 3121910778619 {
		t.Errorf("part2: got %d, want 3121910778619", part2)
	}
}

func TestPickKDigits(t *testing.T) {
	tests := []struct {
		name    string
		bank    string
		k       int
		want    uint64
		wantErr bool
	}{
		{"pick 2 from 987654321111111", "987654321111111", 2, 98, false},
		{"pick 2 from 811111111111119", "811111111111119", 2, 89, false},
		{"pick 12 from 987654321111111", "987654321111111", 12, 987654321111, false},
		{"pick 12 from 811111111111119", "811111111111119", 12, 811111111119, false},
		{"short bank", "54321", 2, 54, false},
		{"not enough digits", "123", 5, 0, true},
		{"all same", "1111111111111111", 2, 11, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pickKDigits(tt.bank, tt.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("pickKDigits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("pickKDigits() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPart1 uint64
		wantPart2 uint64
		wantErr   bool
	}{
		{"single bank", "98\n", 9, 0, true},
		{"valid input", "987654321111111\n", 98, 987654321111, false},
		{"empty lines skipped", "98\n\n99\n", 197, 0, true},
		{"invalid character", "12a34567890123456\n", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := solve(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("solve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got1 != tt.wantPart1 {
					t.Errorf("solve() part1 = %d, want %d", got1, tt.wantPart1)
				}
				if got2 != tt.wantPart2 {
					t.Errorf("solve() part2 = %d, want %d", got2, tt.wantPart2)
				}
			}
		})
	}
}
