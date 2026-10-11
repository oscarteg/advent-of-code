package main

import (
	"os"
	"testing"
)

func TestSolve(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		input        string
		tiles        int
		part1, part2 uint32
		expectErr    bool
	}{
		{name: "official example", input: string(input), tiles: 5, part1: 40, part2: 315},
		{name: "single cell grid has zero risk", input: "5", tiles: 1, part1: 0, part2: 0},
		{name: "tile count of one matches untiled grid", input: string(input), tiles: 1, part1: 40, part2: 40},
		{name: "uneven rows are rejected", input: "123\n12", expectErr: true},
		{name: "non digit cells are rejected", input: "12a\n345", expectErr: true},
		{name: "empty input is rejected", input: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2, err := solveWith(tt.input, tt.tiles)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected an error, got part1=%d part2=%d", part1, part2)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if part1 != tt.part1 || part2 != tt.part2 {
				t.Errorf("solveWith() = (%d, %d), want (%d, %d)", part1, part2, tt.part1, tt.part2)
			}
		})
	}
}
