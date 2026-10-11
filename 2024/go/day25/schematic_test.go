package main

import (
	"os"
	"strings"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	part1, err := solve(string(input))

	if err != nil || part1 != 3 {
		t.Fatalf("got (%d, %v), want (3, nil)", part1, err)
	}
}

func TestFittingPairs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      int
		wantError bool
	}{
		{
			name: "no pairs fit when every pin collides",
			input: "#####\n#####\n#####\n#####\n#####\n#####\n.....\n" +
				"\n" +
				".....\n#####\n#####\n#####\n#####\n#####\n#####",
			want: 0,
		},
		{
			name: "windows line endings",
			input: strings.ReplaceAll(
				"#####\n.####\n.####\n.####\n.#.#.\n.#...\n.....\n"+
					"\n"+
					".....\n#....\n#....\n#...#\n#.#.#\n#.###\n#####",
				"\n", "\r\n"),
			want: 0,
		},
		{name: "not a schematic", input: "not a schematic", wantError: true},
		{name: "wrong row count", input: "#####\n.####\n.....", wantError: true},
		{name: "neither lock nor key", input: "#####\n#####\n#####\n#####\n#####\n#####\n#####", wantError: true},
		{name: "only locks", input: "#####\n#####\n#####\n#####\n#####\n#####\n.....", wantError: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := solve(tc.input)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
