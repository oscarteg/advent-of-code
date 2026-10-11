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

	if err != nil || part1 != 1320 || part2 != 145 {
		t.Fatalf("got (%d, %d, %v), want (1320, 145, nil)", part1, part2, err)
	}
}

func TestHashKnownValues(t *testing.T) {
	tests := []struct {
		step string
		want uint32
	}{
		{"HASH", 52},
		{"rn=1", 30},
		{"cm-", 253},
		{"qp=3", 97},
	}
	for _, tc := range tests {
		if got := hash(tc.step); got != tc.want {
			t.Errorf("hash(%q) = %d, want %d", tc.step, got, tc.want)
		}
	}
}

func TestSolveEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		part1, part2 int
		wantErr      bool
	}{
		{"empty input", "", 0, 0, false},
		{"reinsert replaces focal length", "rn=1,rn=2", int(hash("rn=1") + hash("rn=2")), 2, false},
		{"remove clears slot", "rn=1,rn-", int(hash("rn=1") + hash("rn-")), 0, false},
		{"malformed step missing operator", "rn", 0, 0, true},
		{"malformed step empty label", "=5", 0, 0, true},
		{"malformed step non numeric focal length", "rn=abc", 0, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2, err := solve(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if part1 != tc.part1 || part2 != tc.part2 {
				t.Fatalf("got (%d, %d), want (%d, %d)", part1, part2, tc.part1, tc.part2)
			}
		})
	}
}
