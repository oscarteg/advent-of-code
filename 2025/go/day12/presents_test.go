package main

import (
	"os"
	"strings"
	"testing"
)

const (
	cPresent   = "0:\n###\n#..\n###\n\n"
	barPresent = "0:\n###\n...\n...\n\n"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	fitting, err := solve(string(input))

	if err != nil || fitting != 2 {
		t.Fatalf("got (%d, %v), want (2, nil)", fitting, err)
	}
}

func TestRegions(t *testing.T) {
	tests := []struct {
		name, input string
		fitting     int
		wantError   bool
	}{
		// The puzzle's own 4x4 example interlocks two of these presents.
		{"interlocking presents", cPresent + "4x4: 2", 1, false},
		{"area too small", cPresent + "3x3: 2", 0, false},
		{"one present per block", cPresent + "3x3: 1\n6x6: 4", 2, false},
		{"region without presents", cPresent + "1x1: 0", 1, false},
		{"bar needs a quarter turn", barPresent + "1x3: 1\n3x1: 1", 2, false},
		{"bar never fits 2x2", barPresent + "2x2: 1", 0, false},
		{"windows line endings", strings.ReplaceAll(cPresent+"4x4: 2", "\n", "\r\n"), 1, false},
		{"no presents", "4x4: 1", 0, true},
		{"empty input", "", 0, true},
		{"short present row", "0:\n##\n##.\n##.\n\n4x4: 1", 0, true},
		{"unknown present cell", "0:\n###\n#x.\n###\n\n4x4: 1", 0, true},
		{"present without cells", "0:\n...\n...\n...\n\n4x4: 1", 0, true},
		{"truncated present", "0:\n###\n###", 0, true},
		{"present out of order", "1:\n###\n###\n###\n\n4x4: 1", 0, true},
		{"present after region", cPresent + "4x4: 1\n1:\n###\n###\n###", 0, true},
		{"count mismatch", cPresent + "4x4: 1 2", 0, true},
		{"region size without x", cPresent + "4by4: 1", 0, true},
		{"zero sized region", cPresent + "0x4: 1", 0, true},
		{"unparsable count", cPresent + "4x4: x", 0, true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fitting, err := solve(testCase.input)

			if (err != nil) != testCase.wantError {
				t.Fatalf("error = %v, wantError = %v", err, testCase.wantError)
			}
			if !testCase.wantError && fitting != testCase.fitting {
				t.Fatalf("got %d, want %d", fitting, testCase.fitting)
			}
		})
	}
}

func TestPresentOrientations(t *testing.T) {
	tests := []struct {
		name     string
		rows     []string
		variants int
	}{
		{"square has one orientation", []string{"###", "###", "###"}, 1},
		{"bar has two", []string{"###", "...", "..."}, 2},
		{"h shape has two", []string{"###", ".#.", "###"}, 2},
		{"c shape has four", []string{"###", "#..", "###"}, 4},
		{"corner shape has eight", []string{"###", "##.", "##."}, 8},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			parsed, err := parsePresent(testCase.rows)
			if err != nil {
				t.Fatalf("parsePresent failed: %v", err)
			}
			if len(parsed.variants) != testCase.variants {
				t.Fatalf("got %d orientations, want %d", len(parsed.variants), testCase.variants)
			}
		})
	}
}
