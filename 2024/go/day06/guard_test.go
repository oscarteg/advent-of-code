package main

import (
	"os"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, part2, err := solve(string(input))

	if err != nil || part1 != 41 || part2 != 6 {
		t.Fatalf("got (%d, %d, %v), want (41, 6, nil)", part1, part2, err)
	}
}

func TestGuardPatrol(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantErr      bool
	}{
		{"straight line leaves immediately", "...\n.^.\n...\n", 2, 0, false},
		{"empty input", "", 0, 0, true},
		{"irregular grid", "..\n...\n", 0, 0, true},
		{"invalid character", "..\n.?\n", 0, 0, true},
		{"missing guard", "...\n...\n", 0, 0, true},
		{"multiple guards", "^.\n.^\n", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)

			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tc.wantErr)
			}
			if !tc.wantErr && (part1 != tc.part1 || part2 != tc.part2) {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
