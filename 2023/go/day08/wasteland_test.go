package main

import (
	"os"
	"strings"
	"testing"
)

const (
	repeatedInstructionsExample = "LLR\n\nAAA = (BBB, BBB)\nBBB = (AAA, ZZZ)\nZZZ = (ZZZ, ZZZ)"
	ghostExample                = "LR\n\n11A = (11B, XXX)\n11B = (XXX, 11Z)\n11Z = (11B, XXX)\n22A = (22B, XXX)\n22B = (22C, 22C)\n22C = (22Z, 22Z)\n22Z = (22B, 22B)\nXXX = (XXX, XXX)"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	camel, ghosts, err := solve(string(input))

	if err != nil || camel != 2 || ghosts != 2 {
		t.Fatalf("got (%d, %d, %v), want (2, 2, nil)", camel, ghosts, err)
	}
}

func TestCamelSteps(t *testing.T) {
	tests := []struct {
		name, input string
		want        uint64
	}{
		{"repeated instructions", repeatedInstructionsExample, 6},
		{"arrives mid instruction list", "LR\n\nAAA = (ZZZ, AAA)\nZZZ = (ZZZ, ZZZ)", 1},
		{"windows line endings", strings.ReplaceAll(repeatedInstructionsExample, "\n", "\r\n"), 6},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			camel, _, err := solve(tc.input)

			if err != nil || camel != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", camel, err, tc.want)
			}
		})
	}
}

func TestGhostSteps(t *testing.T) {
	tests := []struct {
		name, input string
		want        uint64
	}{
		{"official example", ghostExample, 6},
		{"single ghost matches the camel", repeatedInstructionsExample, 6},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			docs, err := parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			ghosts, err := docs.ghostSteps()

			if err != nil || ghosts != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", ghosts, err, tc.want)
			}
		})
	}
}

func TestLeastCommonMultiple(t *testing.T) {
	tests := []struct {
		name      string
		a, b      uint64
		want      uint64
		wantError bool
	}{
		{"coprime", 2, 3, 6, false},
		{"shared factor", 4, 6, 12, false},
		{"equal", 7, 7, 7, false},
		{"identity", 1, 13, 13, false},
		{"overflow", 1 << 63, 3, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := leastCommonMultiple(tc.a, tc.b)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMalformedDocuments(t *testing.T) {
	tests := []struct{ name, input string }{
		{"empty", ""},
		{"only instructions", "LR"},
		{"unknown instruction", "LX\n\nAAA = (AAA, ZZZ)\nZZZ = (ZZZ, ZZZ)"},
		{"node without assignment", "LR\n\nAAA (BBB, CCC)"},
		{"node without parentheses", "LR\n\nAAA = BBB, CCC"},
		{"node with one target", "LR\n\nAAA = (AAA)"},
		{"unknown target", "LR\n\nAAA = (BBB, CCC)"},
		{"duplicate node", "LR\n\nAAA = (AAA, AAA)\nAAA = (AAA, AAA)"},
		{"missing start node", "LR\n\nBBB = (BBB, BBB)"},
		{"unreachable destination", "L\n\nAAA = (AAA, AAA)\nZZZ = (ZZZ, ZZZ)"},
		{"ghost without destination", "L\n\nAAA = (ZZZ, ZZZ)\n11A = (11A, 11A)\nZZZ = (ZZZ, ZZZ)"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := solve(tc.input); err == nil {
				t.Fatalf("got nil error for %q", tc.input)
			}
		})
	}
}
